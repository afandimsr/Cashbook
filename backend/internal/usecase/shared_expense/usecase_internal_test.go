package shared_expense

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	se "github.com/afandimsr/cashbook-backend/internal/domain/shared_expense"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func i64(v int64) *int64 { return &v }

// --- fakeRepo: a hand-written se.Repository for orchestration tests ---

type fakeRepo struct {
	bill              se.SplitBill
	participantsOnAll []se.SplitParticipant // returned by FindParticipantsByBillID

	savedBill         *se.SplitBill
	savedParticipants []*se.SplitParticipant
	savedItems        []*se.SplitItem
	itemLinks         [][2]string
	savedTx           []*transaction.Transaction
	settledIDs        []string
	updatedStatus     se.SplitBillStatus
	pCount            int
	itemCount         int
}

func (f *fakeRepo) FindAllByUserID(ctx context.Context, userID int64, filter se.Filter) ([]se.SplitBill, error) {
	return []se.SplitBill{f.bill}, nil
}
func (f *fakeRepo) FindByID(ctx context.Context, id string) (se.SplitBill, error) {
	return f.bill, nil
}
func (f *fakeRepo) Save(ctx context.Context, b *se.SplitBill) error {
	b.ID = "bill-1"
	f.savedBill = b
	return nil
}
func (f *fakeRepo) UpdateStatus(ctx context.Context, id string, status se.SplitBillStatus) error {
	f.updatedStatus = status
	return nil
}
func (f *fakeRepo) SaveParticipant(ctx context.Context, p *se.SplitParticipant) error {
	f.pCount++
	p.ID = fmt.Sprintf("p-%d", f.pCount)
	f.savedParticipants = append(f.savedParticipants, p)
	return nil
}
func (f *fakeRepo) UpdateParticipantSettlement(ctx context.Context, participantID string, isPaid bool, settledAt *time.Time, transactionID *int64) error {
	f.settledIDs = append(f.settledIDs, participantID)
	return nil
}
func (f *fakeRepo) FindParticipantsByBillID(ctx context.Context, billID string) ([]se.SplitParticipant, error) {
	return f.participantsOnAll, nil
}
func (f *fakeRepo) SaveItem(ctx context.Context, item *se.SplitItem) error {
	f.itemCount++
	item.ID = fmt.Sprintf("item-%d", f.itemCount)
	f.savedItems = append(f.savedItems, item)
	return nil
}
func (f *fakeRepo) AddItemParticipant(ctx context.Context, itemID, participantID string) error {
	f.itemLinks = append(f.itemLinks, [2]string{itemID, participantID})
	return nil
}
func (f *fakeRepo) FindItemsByBillID(ctx context.Context, billID string) ([]se.SplitItem, error) {
	return f.bill.Items, nil
}
func (f *fakeRepo) GetSummary(ctx context.Context, userID int64) (se.SplitSummary, error) {
	return se.SplitSummary{}, nil
}
func (f *fakeRepo) SaveTransaction(ctx context.Context, t *transaction.Transaction) error {
	f.savedTx = append(f.savedTx, t)
	return nil
}
func (f *fakeRepo) WithTransaction(ctx context.Context, fn func(repo se.Repository) error) error {
	return fn(f)
}

// sumShares returns the sum of every participant's final share.
func sumShares(bill *se.SplitBill) float64 {
	var s float64
	for i := range bill.Participants {
		s += bill.Participants[i].ShareAmount
	}
	return s
}

// --- computeShares (pure logic) ---

func TestComputeShares_Equal(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodEqual, Subtotal: 100000,
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
	}
	err := computeShares(&bill)
	assert.NoError(t, err)
	assert.Equal(t, 100000.0, bill.TotalAmount)
	assert.Equal(t, 50000.0, bill.Participants[0].ShareAmount)
	assert.Equal(t, 50000.0, bill.Participants[1].ShareAmount)
	assert.Equal(t, bill.TotalAmount, sumShares(&bill))
}

func TestComputeShares_Equal_PayerAbsorbsRemainder(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodEqual, Subtotal: 100001,
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
	}
	err := computeShares(&bill)
	assert.NoError(t, err)
	assert.Equal(t, 100001.0, sumShares(&bill))
	assert.Equal(t, 50001.0, bill.Participants[0].ShareAmount) // payer absorbs the extra rupiah
	assert.Equal(t, 50000.0, bill.Participants[1].ShareAmount)
}

func TestComputeShares_Exact(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodExact, Subtotal: 100,
		Participants: []se.SplitParticipant{{UserID: i64(1), BaseShare: 30}, {UserID: i64(2), BaseShare: 70}},
	}
	err := computeShares(&bill)
	assert.NoError(t, err)
	assert.Equal(t, 30.0, bill.Participants[0].ShareAmount)
	assert.Equal(t, 70.0, bill.Participants[1].ShareAmount)
}

func TestComputeShares_Exact_Mismatch(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodExact, Subtotal: 100,
		Participants: []se.SplitParticipant{{UserID: i64(1), BaseShare: 30}, {UserID: i64(2), BaseShare: 60}},
	}
	err := computeShares(&bill)
	assert.ErrorIs(t, err, se.ErrInvalidSplit)
}

func TestComputeShares_Percentage(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodPercentage, Subtotal: 200,
		Participants: []se.SplitParticipant{{UserID: i64(1), Percentage: 25}, {UserID: i64(2), Percentage: 75}},
	}
	err := computeShares(&bill)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, bill.Participants[0].ShareAmount)
	assert.Equal(t, 150.0, bill.Participants[1].ShareAmount)
}

func TestComputeShares_Percentage_NotHundred(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodPercentage, Subtotal: 200,
		Participants: []se.SplitParticipant{{UserID: i64(1), Percentage: 25}, {UserID: i64(2), Percentage: 70}},
	}
	assert.ErrorIs(t, computeShares(&bill), se.ErrInvalidSplit)
}

func TestComputeShares_Item(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodItem,
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
		Items: []se.SplitItem{
			{Name: "Nasi", Price: 100, Quantity: 1, ParticipantIndexes: []int{0, 1}},
			{Name: "Es Teh", Price: 40, Quantity: 1, ParticipantIndexes: []int{1}},
		},
	}
	err := computeShares(&bill)
	assert.NoError(t, err)
	assert.Equal(t, 140.0, bill.Subtotal) // computed from items
	assert.Equal(t, 50.0, bill.Participants[0].ShareAmount)
	assert.Equal(t, 90.0, bill.Participants[1].ShareAmount)
	assert.Equal(t, 140.0, sumShares(&bill))
}

func TestComputeShares_Item_NoItems(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodItem,
		Participants: []se.SplitParticipant{{UserID: i64(1)}},
	}
	assert.ErrorIs(t, computeShares(&bill), se.ErrInvalidSplit)
}

func TestComputeShares_Item_ItemWithoutParticipants(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodItem,
		Participants: []se.SplitParticipant{{UserID: i64(1)}},
		Items:        []se.SplitItem{{Name: "X", Price: 10, Quantity: 1, ParticipantIndexes: nil}},
	}
	assert.ErrorIs(t, computeShares(&bill), se.ErrInvalidSplit)
}

func TestComputeShares_ChargesDistributedProportionally(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, SplitMethod: se.MethodEqual, Subtotal: 100000,
		TaxAmount: 10000, ServiceCharge: 5000, DiscountAmount: 5000, OtherCharge: 0,
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
	}
	err := computeShares(&bill)
	assert.NoError(t, err)
	// net = 100000 - 5000 + 10000 + 5000 = 110000
	assert.Equal(t, 110000.0, bill.TotalAmount)
	assert.Equal(t, 55000.0, bill.Participants[0].ShareAmount)
	assert.Equal(t, 55000.0, bill.Participants[1].ShareAmount)
	assert.Equal(t, bill.TotalAmount, sumShares(&bill))
}

func TestComputeShares_Validations(t *testing.T) {
	base := func() se.SplitBill {
		return se.SplitBill{
			PayerID: 1, SplitMethod: se.MethodEqual, Subtotal: 1000,
			Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
		}
	}

	t.Run("DiscountExceedsSubtotal", func(t *testing.T) {
		b := base()
		b.DiscountAmount = 2000
		assert.ErrorIs(t, computeShares(&b), se.ErrInvalidSplit)
	})
	t.Run("NegativeCharge", func(t *testing.T) {
		b := base()
		b.TaxAmount = -1
		assert.ErrorIs(t, computeShares(&b), se.ErrInvalidSplit)
	})
	t.Run("ZeroSubtotal", func(t *testing.T) {
		b := base()
		b.Subtotal = 0
		assert.ErrorIs(t, computeShares(&b), se.ErrInvalidSplit)
	})
	t.Run("UnknownMethod", func(t *testing.T) {
		b := base()
		b.SplitMethod = "WEIRD"
		assert.ErrorIs(t, computeShares(&b), se.ErrInvalidSplit)
	})
}

// --- payerExpenses (ledger split by category) ---

func TestPayerExpenses_SingleForNonItem(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, CategoryID: 5, Title: "T", SplitMethod: se.MethodEqual,
		Subtotal: 100000, TotalAmount: 100000,
	}
	txs := payerExpenses(&bill)
	assert.Len(t, txs, 1)
	assert.Equal(t, int64(5), txs[0].CategoryID)
	assert.Equal(t, 100000.0, txs[0].Amount)
	assert.Equal(t, "Split: T", txs[0].Note)
}

func TestPayerExpenses_ItemSingleCategoryListsItems(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, CategoryID: 5, Title: "WARUNK", SplitMethod: se.MethodItem,
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
		Items: []se.SplitItem{
			{Name: "Nasi", Price: 50000, Quantity: 1, ParticipantIndexes: []int{0, 1}},
			{Name: "Ayam", Price: 30000, Quantity: 1, ParticipantIndexes: []int{0}},
		},
	}
	assert.NoError(t, computeShares(&bill))
	txs := payerExpenses(&bill)
	assert.Len(t, txs, 1)
	assert.Equal(t, int64(5), txs[0].CategoryID)
	assert.Equal(t, "Split: WARUNK (Nasi, Ayam)", txs[0].Note)
	assert.Equal(t, bill.TotalAmount, txs[0].Amount)
}

func TestPayerExpenses_ItemMultiCategory(t *testing.T) {
	bill := se.SplitBill{
		PayerID: 1, CategoryID: 5, Title: "WARUNK", SplitMethod: se.MethodItem,
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
		Items: []se.SplitItem{
			{Name: "Nasi", Price: 50000, Quantity: 1, CategoryID: i64(10), ParticipantIndexes: []int{0, 1}},
			{Name: "Es Teh", Price: 10000, Quantity: 1, CategoryID: i64(20), ParticipantIndexes: []int{1}},
		},
	}
	assert.NoError(t, computeShares(&bill))
	txs := payerExpenses(&bill)
	assert.Len(t, txs, 2)

	var total float64
	byCat := map[int64]*transaction.Transaction{}
	for _, tx := range txs {
		total += tx.Amount
		byCat[tx.CategoryID] = tx
	}
	assert.Equal(t, bill.TotalAmount, total) // allocation sums to the net
	assert.Contains(t, byCat[10].Note, "(Nasi)")
	assert.Contains(t, byCat[20].Note, "(Es Teh)")
}

// --- CreateSplitBill (orchestration) ---

func TestCreateSplitBill_HappyPath(t *testing.T) {
	f := &fakeRepo{}
	u := New(f)

	err := u.CreateSplitBill(context.Background(), se.SplitBill{
		PayerID: 1, CategoryID: 5, Title: "Lunch", SplitMethod: se.MethodEqual,
		Subtotal: 100000, Date: time.Now(),
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
	})
	assert.NoError(t, err)
	assert.Len(t, f.savedParticipants, 2)
	assert.True(t, f.savedParticipants[0].IsPaid, "payer is marked paid on create")
	assert.False(t, f.savedParticipants[1].IsPaid)
	assert.Len(t, f.savedTx, 1)
	assert.Equal(t, 100000.0, f.savedTx[0].Amount)
	assert.Equal(t, "expense", f.savedTx[0].Type)
}

func TestCreateSplitBill_DefaultsMethodAndLinksItems(t *testing.T) {
	f := &fakeRepo{}
	u := New(f)

	err := u.CreateSplitBill(context.Background(), se.SplitBill{
		PayerID: 1, CategoryID: 5, Title: "Cafe", SplitMethod: se.MethodItem, Date: time.Now(),
		Participants: []se.SplitParticipant{{UserID: i64(1)}, {UserID: i64(2)}},
		Items: []se.SplitItem{
			{Name: "Nasi", Price: 40000, Quantity: 1, ParticipantIndexes: []int{0, 1}},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, f.savedItems, 1)
	assert.Len(t, f.itemLinks, 2) // item linked to both participants
}

func TestCreateSplitBill_ValidationError(t *testing.T) {
	f := &fakeRepo{}
	u := New(f)

	err := u.CreateSplitBill(context.Background(), se.SplitBill{
		PayerID: 1, Title: "X", SplitMethod: se.MethodEqual, Subtotal: 0,
		Participants: []se.SplitParticipant{{UserID: i64(1)}},
	})
	assert.ErrorIs(t, err, se.ErrInvalidSplit)
	assert.Empty(t, f.savedTx)
}

func TestCreateSplitBill_NoParticipants(t *testing.T) {
	f := &fakeRepo{}
	u := New(f)
	err := u.CreateSplitBill(context.Background(), se.SplitBill{PayerID: 1, SplitMethod: se.MethodEqual, Subtotal: 100})
	assert.ErrorIs(t, err, se.ErrInvalidSplit)
}

// --- SettleParticipant ---

func settleBill() se.SplitBill {
	return se.SplitBill{
		ID: "bill-1", PayerID: 1, CategoryID: 5, Title: "WARUNK",
		SplitMethod: se.MethodItem,
		Participants: []se.SplitParticipant{
			{ID: "p1", UserID: i64(1), ShareAmount: 25000, IsPaid: true},
			{ID: "p2", UserID: i64(2), ShareAmount: 35000, IsPaid: false},
		},
		Items: []se.SplitItem{
			{Name: "Nasi", ParticipantIDs: []string{"p1", "p2"}},
			{Name: "Es Teh", ParticipantIDs: []string{"p2"}},
		},
	}
}

func TestSettleParticipant_RecordsTransactionsWithItemNames(t *testing.T) {
	f := &fakeRepo{bill: settleBill()}
	// After settling p2, all participants are paid → bill becomes SETTLED.
	f.participantsOnAll = []se.SplitParticipant{
		{ID: "p1", IsPaid: true}, {ID: "p2", IsPaid: true},
	}
	u := New(f)

	err := u.SettleParticipant(context.Background(), "bill-1", "p2", 1)
	assert.NoError(t, err)
	assert.Len(t, f.savedTx, 2) // income for payer + expense for the registered participant
	assert.Equal(t, "income", f.savedTx[0].Type)
	assert.Contains(t, f.savedTx[0].Note, "Settlement from")
	assert.Contains(t, f.savedTx[0].Note, "(Nasi, Es Teh)")
	assert.Equal(t, "expense", f.savedTx[1].Type)
	assert.Equal(t, []string{"p2"}, f.settledIDs)
	assert.Equal(t, se.StatusSettled, f.updatedStatus)
}

func TestSettleParticipant_NotAllPaidKeepsStatusOpen(t *testing.T) {
	f := &fakeRepo{bill: settleBill()}
	f.participantsOnAll = []se.SplitParticipant{
		{ID: "p1", IsPaid: true}, {ID: "p2", IsPaid: false},
	}
	u := New(f)

	err := u.SettleParticipant(context.Background(), "bill-1", "p2", 1)
	assert.NoError(t, err)
	assert.Equal(t, se.SplitBillStatus(""), f.updatedStatus) // UpdateStatus not called
}

func TestSettleParticipant_OnlyPayerCanSettle(t *testing.T) {
	f := &fakeRepo{bill: settleBill()}
	u := New(f)

	err := u.SettleParticipant(context.Background(), "bill-1", "p2", 999)
	assert.Error(t, err)
	assert.Empty(t, f.savedTx)
}

func TestSettleParticipant_AlreadySettled(t *testing.T) {
	bill := settleBill()
	bill.Participants[1].IsPaid = true
	f := &fakeRepo{bill: bill}
	u := New(f)

	err := u.SettleParticipant(context.Background(), "bill-1", "p2", 1)
	assert.Error(t, err)
}

func TestSettleParticipant_NotFound(t *testing.T) {
	f := &fakeRepo{bill: settleBill()}
	u := New(f)

	err := u.SettleParticipant(context.Background(), "bill-1", "nope", 1)
	assert.Error(t, err)
}

func TestSettleParticipant_ShadowNoExpense(t *testing.T) {
	bill := settleBill()
	bill.Participants[1] = se.SplitParticipant{ID: "p2", UserID: nil, ShadowName: "Guest", ShareAmount: 35000}
	f := &fakeRepo{bill: bill}
	f.participantsOnAll = []se.SplitParticipant{{ID: "p1", IsPaid: true}, {ID: "p2", IsPaid: true}}
	u := New(f)

	err := u.SettleParticipant(context.Background(), "bill-1", "p2", 1)
	assert.NoError(t, err)
	assert.Len(t, f.savedTx, 1) // only the payer's income; no expense for a shadow user
	assert.True(t, strings.HasPrefix(f.savedTx[0].Note, "Settlement from Guest"))
}

// ensure errors package stays referenced even if assertions change
var _ = errors.Is
