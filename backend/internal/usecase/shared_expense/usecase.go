package shared_expense

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/shared_expense"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type Usecase interface {
	CreateSplitBill(ctx context.Context, bill shared_expense.SplitBill) error
	GetByID(ctx context.Context, id string) (shared_expense.SplitBill, error)
	GetByUserID(ctx context.Context, userID int64, filter shared_expense.Filter) ([]shared_expense.SplitBill, error)
	SettleParticipant(ctx context.Context, billID, participantID string, userID int64) error
	GetSummary(ctx context.Context, userID int64) (shared_expense.SplitSummary, error)
}

type usecase struct {
	repo shared_expense.Repository
}

func New(repo shared_expense.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateSplitBill(ctx context.Context, bill shared_expense.SplitBill) error {
	if len(bill.Participants) == 0 {
		return fmt.Errorf("%w: at least one participant is required", shared_expense.ErrInvalidSplit)
	}

	if bill.SplitMethod == "" {
		bill.SplitMethod = shared_expense.MethodEqual
	}

	// Compute base + final shares (and the subtotal for the ITEM method),
	// distributing tax and discount proportionally.
	if err := computeShares(&bill); err != nil {
		return err
	}

	bill.Status = shared_expense.StatusOpen

	// Use transaction for atomic split bill + standard transaction creation
	return u.repo.WithTransaction(ctx, func(txRepo shared_expense.Repository) error {
		if err := txRepo.Save(ctx, &bill); err != nil {
			return err
		}

		for i := range bill.Participants {
			p := &bill.Participants[i]
			p.SplitBillID = bill.ID

			// If the participant is the payer, mark as paid immediately
			if p.UserID != nil && *p.UserID == bill.PayerID {
				p.IsPaid = true
				now := time.Now()
				p.SettledAt = &now
			}

			if err := txRepo.SaveParticipant(ctx, p); err != nil {
				return err
			}
		}

		// Persist line items and their participant assignments (ITEM method).
		for i := range bill.Items {
			item := &bill.Items[i]
			item.SplitBillID = bill.ID
			if err := txRepo.SaveItem(ctx, item); err != nil {
				return err
			}
			for _, idx := range item.ParticipantIndexes {
				if idx < 0 || idx >= len(bill.Participants) {
					continue
				}
				if err := txRepo.AddItemParticipant(ctx, item.ID, bill.Participants[idx].ID); err != nil {
					return err
				}
			}
		}

		// Record the payer's expense within the same DB transaction so the split
		// and the ledger stay consistent. For per-item categories this becomes
		// one transaction per category.
		for _, t := range payerExpenses(&bill) {
			if err := txRepo.SaveTransaction(ctx, t); err != nil {
				return err
			}
		}
		return nil
	})
}

// payerExpenses builds the payer's ledger expense(s) for the full bill total.
// For the ITEM method with multiple distinct categories, the total (including
// tax, service, other charges, and discount) is allocated proportionally across
// those categories so budgets/reports see the right split. Otherwise it is a
// single expense under the bill's category.
func payerExpenses(bill *shared_expense.SplitBill) []*transaction.Transaction {
	baseNote := fmt.Sprintf("Split: %s", bill.Title)
	net := bill.TotalAmount

	// Non-ITEM methods (or no items) → single expense under the bill category.
	if bill.SplitMethod != shared_expense.MethodItem || bill.Subtotal <= 0 || len(bill.Items) == 0 {
		return []*transaction.Transaction{{
			UserID:     bill.PayerID,
			CategoryID: bill.CategoryID,
			Amount:     net,
			Note:       baseNote,
			Date:       bill.Date,
			Type:       "expense",
		}}
	}

	// Group item totals and names by effective category (item category, else bill).
	catTotal := map[int64]float64{}
	catNames := map[int64][]string{}
	var order []int64
	for i := range bill.Items {
		item := &bill.Items[i]
		cid := bill.CategoryID
		if item.CategoryID != nil {
			cid = *item.CategoryID
		}
		if _, ok := catTotal[cid]; !ok {
			order = append(order, cid)
		}
		qty := item.Quantity
		if qty <= 0 {
			qty = 1
		}
		catTotal[cid] += item.Price * float64(qty)
		catNames[cid] = append(catNames[cid], item.Name)
	}

	// Single category → one expense listing every item.
	if len(order) == 1 {
		cid := order[0]
		return []*transaction.Transaction{{
			UserID:     bill.PayerID,
			CategoryID: cid,
			Amount:     net,
			Note:       baseNote + itemsParen(catNames[cid]),
			Date:       bill.Date,
			Type:       "expense",
		}}
	}

	// Multiple categories → one expense per category, each listing its items.
	var txs []*transaction.Transaction
	var allocated float64
	for k, cid := range order {
		var amt float64
		if k == len(order)-1 {
			amt = math.Round(net - allocated) // last category absorbs the remainder
		} else {
			amt = math.Round(catTotal[cid] / bill.Subtotal * net)
			allocated += amt
		}
		if amt == 0 {
			continue
		}
		txs = append(txs, &transaction.Transaction{
			UserID:     bill.PayerID,
			CategoryID: cid,
			Amount:     amt,
			Note:       baseNote + itemsParen(catNames[cid]),
			Date:       bill.Date,
			Type:       "expense",
		})
	}
	return txs
}

// computeShares fills BaseShare and ShareAmount for every participant based on
// the split method, then distributes tax and discount proportionally over the
// base shares. All amounts are whole rupiah; the payer absorbs the rounding
// remainder so the participant shares sum exactly to the total.
func computeShares(bill *shared_expense.SplitBill) error {
	n := len(bill.Participants)

	switch bill.SplitMethod {
	case shared_expense.MethodEqual:
		if bill.Subtotal <= 0 {
			return fmt.Errorf("%w: subtotal must be greater than zero", shared_expense.ErrInvalidSplit)
		}
		base := math.Floor(bill.Subtotal / float64(n))
		for i := range bill.Participants {
			bill.Participants[i].BaseShare = base
		}

	case shared_expense.MethodExact:
		if bill.Subtotal <= 0 {
			return fmt.Errorf("%w: subtotal must be greater than zero", shared_expense.ErrInvalidSplit)
		}
		var sum float64
		for i := range bill.Participants {
			sum += bill.Participants[i].BaseShare
		}
		if math.Round(sum) != math.Round(bill.Subtotal) {
			return fmt.Errorf("%w: exact shares (%.0f) must equal subtotal (%.0f)", shared_expense.ErrInvalidSplit, sum, bill.Subtotal)
		}

	case shared_expense.MethodPercentage:
		if bill.Subtotal <= 0 {
			return fmt.Errorf("%w: subtotal must be greater than zero", shared_expense.ErrInvalidSplit)
		}
		var pct float64
		for i := range bill.Participants {
			pct += bill.Participants[i].Percentage
		}
		if math.Abs(pct-100) > 0.01 {
			return fmt.Errorf("%w: percentages must sum to 100 (got %.2f)", shared_expense.ErrInvalidSplit, pct)
		}
		for i := range bill.Participants {
			bill.Participants[i].BaseShare = math.Floor(bill.Subtotal * bill.Participants[i].Percentage / 100)
		}

	case shared_expense.MethodItem:
		if len(bill.Items) == 0 {
			return fmt.Errorf("%w: at least one item is required for the ITEM method", shared_expense.ErrInvalidSplit)
		}
		for i := range bill.Participants {
			bill.Participants[i].BaseShare = 0
		}
		var subtotal float64
		for it := range bill.Items {
			item := &bill.Items[it]
			if item.Quantity <= 0 {
				item.Quantity = 1
			}
			if item.Price < 0 {
				return fmt.Errorf("%w: item %q has a negative price", shared_expense.ErrInvalidSplit, item.Name)
			}
			itemTotal := item.Price * float64(item.Quantity)
			subtotal += itemTotal

			idxs := item.ParticipantIndexes
			if len(idxs) == 0 {
				return fmt.Errorf("%w: item %q has no participants assigned", shared_expense.ErrInvalidSplit, item.Name)
			}
			per := math.Floor(itemTotal / float64(len(idxs)))
			rem := math.Round(itemTotal - per*float64(len(idxs)))
			for k, pi := range idxs {
				if pi < 0 || pi >= n {
					return fmt.Errorf("%w: item %q references an invalid participant", shared_expense.ErrInvalidSplit, item.Name)
				}
				share := per
				if k == 0 { // first assigned participant absorbs the item's remainder
					share += rem
				}
				bill.Participants[pi].BaseShare += share
			}
		}
		bill.Subtotal = subtotal

	default:
		return fmt.Errorf("%w: unknown split method %q", shared_expense.ErrInvalidSplit, bill.SplitMethod)
	}

	if bill.TaxAmount < 0 || bill.DiscountAmount < 0 || bill.ServiceCharge < 0 || bill.OtherCharge < 0 {
		return fmt.Errorf("%w: charges and discount must not be negative", shared_expense.ErrInvalidSplit)
	}
	if bill.DiscountAmount > bill.Subtotal {
		return fmt.Errorf("%w: discount cannot exceed subtotal", shared_expense.ErrInvalidSplit)
	}
	if bill.Subtotal <= 0 {
		return fmt.Errorf("%w: subtotal must be greater than zero", shared_expense.ErrInvalidSplit)
	}

	net := math.Round(bill.Subtotal - bill.DiscountAmount + bill.TaxAmount + bill.ServiceCharge + bill.OtherCharge)
	bill.TotalAmount = net

	// Distribute the net proportionally to each base share.
	payerIdx := 0
	for i := range bill.Participants {
		if bill.Participants[i].UserID != nil && *bill.Participants[i].UserID == bill.PayerID {
			payerIdx = i
			break
		}
	}

	var allocated float64
	for i := range bill.Participants {
		if i == payerIdx {
			continue
		}
		final := math.Round(bill.Participants[i].BaseShare / bill.Subtotal * net)
		bill.Participants[i].ShareAmount = final
		allocated += final
	}
	bill.Participants[payerIdx].ShareAmount = math.Round(net - allocated)

	return nil
}

func (u *usecase) GetByID(ctx context.Context, id string) (shared_expense.SplitBill, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *usecase) GetByUserID(ctx context.Context, userID int64, filter shared_expense.Filter) ([]shared_expense.SplitBill, error) {
	return u.repo.FindAllByUserID(ctx, userID, filter)
}

func (u *usecase) SettleParticipant(ctx context.Context, billID, participantID string, userID int64) error {
	bill, err := u.repo.FindByID(ctx, billID)
	if err != nil {
		return err
	}

	if bill.PayerID != userID {
		return fmt.Errorf("only the payer can mark a participant as paid")
	}

	var targetParticipant *shared_expense.SplitParticipant
	for i := range bill.Participants {
		if bill.Participants[i].ID == participantID {
			targetParticipant = &bill.Participants[i]
			break
		}
	}

	if targetParticipant == nil {
		return fmt.Errorf("participant not found")
	}

	if targetParticipant.IsPaid {
		return fmt.Errorf("participant is already settled")
	}

	now := time.Now()

	// Include the participant's assigned item names for clarity, e.g.
	// "WARUNK UBNORMAL (Nasi Goreng, Es Teh)".
	billLabel := bill.Title + participantItemsSuffix(&bill, targetParticipant)

	return u.repo.WithTransaction(ctx, func(txRepo shared_expense.Repository) error {
		// 1. Record Income for Payer
		tPayer := &transaction.Transaction{
			UserID:     bill.PayerID,
			CategoryID: bill.CategoryID,
			Amount:     targetParticipant.ShareAmount,
			Note:       fmt.Sprintf("Settlement from %s: %s", u.getParticipantName(targetParticipant), billLabel),
			Date:       now,
			Type:       "income",
		}
		err = txRepo.SaveTransaction(ctx, tPayer)
		if err != nil {
			return err
		}

		// 2. Record Expense for Participant (if registered)
		if targetParticipant.UserID != nil {
			tPart := &transaction.Transaction{
				UserID:     *targetParticipant.UserID,
				CategoryID: bill.CategoryID,
				Amount:     targetParticipant.ShareAmount,
				Note:       fmt.Sprintf("Settlement to Payer: %s", billLabel),
				Date:       now,
				Type:       "expense",
			}
			err = txRepo.SaveTransaction(ctx, tPart)
			if err != nil {
				return err
			}
		}

		// 3. Update Participant status
		err = txRepo.UpdateParticipantSettlement(ctx, participantID, true, &now, &tPayer.ID)
		if err != nil {
			return err
		}

		// 4. Check if all paid and update bill status
		participants, err := txRepo.FindParticipantsByBillID(ctx, billID)
		if err != nil {
			return err
		}

		allPaid := true
		for _, p := range participants {
			if !p.IsPaid {
				allPaid = false
				break
			}
		}

		if allPaid {
			err = txRepo.UpdateStatus(ctx, billID, shared_expense.StatusSettled)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (u *usecase) GetSummary(ctx context.Context, userID int64) (shared_expense.SplitSummary, error) {
	return u.repo.GetSummary(ctx, userID)
}

func (u *usecase) getParticipantName(p *shared_expense.SplitParticipant) string {
	if p.ShadowName != "" {
		return p.ShadowName
	}
	if p.UserName != "" {
		return p.UserName
	}
	return "User"
}

// participantItemsSuffix returns " (item1, item2)" listing the items assigned to
// the participant, or an empty string when there are none (non-ITEM methods).
func participantItemsSuffix(bill *shared_expense.SplitBill, p *shared_expense.SplitParticipant) string {
	var names []string
	for i := range bill.Items {
		item := &bill.Items[i]
		for _, pid := range item.ParticipantIDs {
			if pid == p.ID {
				names = append(names, item.Name)
				break
			}
		}
	}
	return itemsParen(names)
}

// itemsParen renders item names as " (name1, name2)", or "" when the list is empty.
func itemsParen(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return " (" + strings.Join(names, ", ") + ")"
}
