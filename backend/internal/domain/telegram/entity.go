package telegram

import (
	"context"
	"errors"
	"time"
)

// ErrLinkNotFound is returned when a telegram_chat_id has no linked user,
// or a user has no linked Telegram account.
var ErrLinkNotFound = errors.New("telegram link not found")

// ErrLinkCodeInvalid is returned when a link code is unknown, expired, or already used.
var ErrLinkCodeInvalid = errors.New("link code is invalid, expired, or already used")

// ErrChatAlreadyLinked is returned when the telegram_chat_id being linked is
// already bound to a different user. Never silently reassigned — relinking
// the SAME user to a new chat is fine (see SaveLink's upsert), but a chat
// already owned by someone else must stay a hard error.
var ErrChatAlreadyLinked = errors.New("this telegram account is already linked to another user")

// ErrUserNotFound is returned when an operation references a user_id that
// doesn't exist (e.g. an admin generating a link code for a bad user id).
var ErrUserNotFound = errors.New("user not found")

// UserTelegramLink maps a CashBook user to the Telegram chat that may act on their behalf.
type UserTelegramLink struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	TelegramChatID   int64     `json:"telegram_chat_id"`
	TelegramUsername string    `json:"telegram_username,omitempty"`
	IsActive         bool      `json:"is_active"`
	LinkedAt         time.Time `json:"linked_at"`
}

// LinkCode is a short-lived one-time code a user generates in the app and sends
// to the bot to prove ownership of the Telegram chat being linked.
type LinkCode struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Code      string     `json:"code"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

type Repository interface {
	FindLinkByChatID(ctx context.Context, chatID int64) (UserTelegramLink, error)
	FindLinkByUserID(ctx context.Context, userID int64) (UserTelegramLink, error)
	ListActiveLinks(ctx context.Context) ([]UserTelegramLink, error)
	SaveLink(ctx context.Context, link *UserTelegramLink) error
	Unlink(ctx context.Context, userID int64) error

	CreateLinkCode(ctx context.Context, code *LinkCode) error
	FindActiveLinkCode(ctx context.Context, code string) (LinkCode, error)
	FindActiveLinkCodeByUserID(ctx context.Context, userID int64) (LinkCode, error)
	MarkLinkCodeUsed(ctx context.Context, id int64) error

	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
}
