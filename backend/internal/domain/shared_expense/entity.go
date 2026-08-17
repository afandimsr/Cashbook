package shared_expense

import (
	"context"
	"errors"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

// ErrInvalidSplit wraps validation failures so the delivery layer can map them
// to a 400 Bad Request instead of a 500.
var ErrInvalidSplit = errors.New("invalid split bill")

type SplitBillStatus string

const (
	StatusOpen             SplitBillStatus = "OPEN"
	StatusPartiallySettled SplitBillStatus = "PARTIALLY_SETTLED"
	StatusSettled          SplitBillStatus = "SETTLED"
)

type SplitMethod string

const (
	MethodEqual      SplitMethod = "EQUAL"
	MethodExact      SplitMethod = "EXACT"
	MethodPercentage SplitMethod = "PERCENTAGE"
	MethodItem       SplitMethod = "ITEM"
)

type SplitBill struct {
	ID             string             `json:"id"`
	CreatorID      int64              `json:"creator_id"`
	PayerID        int64              `json:"payer_id"`
	CategoryID     int64              `json:"category_id"`
	Title          string             `json:"title"`
	Subtotal       float64            `json:"subtotal"`
	TaxAmount      float64            `json:"tax_amount"`
	ServiceCharge  float64            `json:"service_charge"`
	OtherCharge    float64            `json:"other_charge"`
	DiscountAmount float64            `json:"discount_amount"`
	TotalAmount    float64            `json:"total_amount"`
	Date           time.Time          `json:"date"`
	Status         SplitBillStatus    `json:"status"`
	SplitMethod    SplitMethod        `json:"split_method"`
	Participants   []SplitParticipant `json:"participants,omitempty"`
	Items          []SplitItem        `json:"items,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type SplitParticipant struct {
	ID            string     `json:"id"`
	SplitBillID   string     `json:"split_bill_id"`
	UserID        *int64     `json:"user_id"`     // Nullable for Shadow Users
	ShadowName    string     `json:"shadow_name"` // Used if UserID is null
	UserName      string     `json:"user_name"`   // Name of the registered user (empty for shadow users)
	BaseShare     float64    `json:"base_share"`  // Portion before tax/discount
	ShareAmount   float64    `json:"share_amount"`
	Percentage    float64    `json:"percentage,omitempty"` // Input-only, used by the PERCENTAGE method
	IsPaid        bool       `json:"is_paid"`
	SettledAt     *time.Time `json:"settled_at"`
	TransactionID *int64     `json:"transaction_id"`
	CreatedAt     time.Time  `json:"created_at"`
}

type SplitItem struct {
	ID          string    `json:"id"`
	SplitBillID string    `json:"split_bill_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CategoryID  *int64    `json:"category_id"` // Falls back to the bill category when null
	// ParticipantIndexes maps to the SplitBill.Participants slice on create (input only).
	ParticipantIndexes []int `json:"participant_indexes,omitempty"`
	// ParticipantIDs lists the participants sharing this item (populated on read).
	ParticipantIDs []string  `json:"participant_ids,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type SplitSummary struct {
	TotalOwedToMe float64 `json:"total_owed_to_me"`
	TotalIOwe     float64 `json:"total_i_owe"`
}

type Filter struct {
	UserID    int64
	Status    SplitBillStatus
	StartDate time.Time
	EndDate   time.Time
}

type Repository interface {
	FindAllByUserID(ctx context.Context, userID int64, filter Filter) ([]SplitBill, error)
	FindByID(ctx context.Context, id string) (SplitBill, error)
	Save(ctx context.Context, bill *SplitBill) error
	UpdateStatus(ctx context.Context, id string, status SplitBillStatus) error

	SaveParticipant(ctx context.Context, participant *SplitParticipant) error
	UpdateParticipantSettlement(ctx context.Context, participantID string, isPaid bool, settledAt *time.Time, transactionID *int64) error
	FindParticipantsByBillID(ctx context.Context, billID string) ([]SplitParticipant, error)

	SaveItem(ctx context.Context, item *SplitItem) error
	AddItemParticipant(ctx context.Context, itemID, participantID string) error
	FindItemsByBillID(ctx context.Context, billID string) ([]SplitItem, error)

	GetSummary(ctx context.Context, userID int64) (SplitSummary, error)

	// SaveTransaction records a standard income/expense transaction using the
	// same executor as the repository, so it participates in WithTransaction.
	SaveTransaction(ctx context.Context, t *transaction.Transaction) error

	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
}
