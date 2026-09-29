package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/afandimsr/cashbook-backend/internal/domain/telegram"
	"github.com/lib/pq"
)

// Postgres error codes this repository handles explicitly instead of letting
// them surface as raw 500s — see https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pqUniqueViolation     = "23505"
	pqForeignKeyViolation = "23503"
)

type telegramQueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type telegramRepo struct {
	db telegramQueryExecutor
}

func NewTelegramRepo(db *sql.DB) telegram.Repository {
	return &telegramRepo{db: db}
}

func (r *telegramRepo) FindLinkByChatID(ctx context.Context, chatID int64) (telegram.UserTelegramLink, error) {
	var l telegram.UserTelegramLink
	var username sql.NullString
	err := r.db.QueryRow(
		"SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links WHERE telegram_chat_id = $1 AND is_active = true",
		chatID,
	).Scan(&l.ID, &l.UserID, &l.TelegramChatID, &username, &l.IsActive, &l.LinkedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return telegram.UserTelegramLink{}, telegram.ErrLinkNotFound
	}
	if err != nil {
		return telegram.UserTelegramLink{}, err
	}
	l.TelegramUsername = username.String
	return l, nil
}

func (r *telegramRepo) FindLinkByUserID(ctx context.Context, userID int64) (telegram.UserTelegramLink, error) {
	var l telegram.UserTelegramLink
	var username sql.NullString
	err := r.db.QueryRow(
		"SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links WHERE user_id = $1 AND is_active = true",
		userID,
	).Scan(&l.ID, &l.UserID, &l.TelegramChatID, &username, &l.IsActive, &l.LinkedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return telegram.UserTelegramLink{}, telegram.ErrLinkNotFound
	}
	if err != nil {
		return telegram.UserTelegramLink{}, err
	}
	l.TelegramUsername = username.String
	return l, nil
}

func (r *telegramRepo) ListActiveLinks(ctx context.Context) ([]telegram.UserTelegramLink, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links WHERE is_active = true",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]telegram.UserTelegramLink, 0)
	for rows.Next() {
		var l telegram.UserTelegramLink
		var username sql.NullString
		if err := rows.Scan(&l.ID, &l.UserID, &l.TelegramChatID, &username, &l.IsActive, &l.LinkedAt); err != nil {
			return nil, err
		}
		l.TelegramUsername = username.String
		links = append(links, l)
	}
	return links, nil
}

// SaveLink upserts on user_id so relinking (e.g. a new phone) just replaces
// the existing row. The telegram_chat_id unique constraint can still fire —
// that means this exact chat is already bound to a DIFFERENT user, which must
// stay a hard error (never silently reassigned).
func (r *telegramRepo) SaveLink(ctx context.Context, link *telegram.UserTelegramLink) error {
	err := r.db.QueryRow(
		`INSERT INTO user_telegram_links (user_id, telegram_chat_id, telegram_username, is_active)
		 VALUES ($1, $2, $3, true)
		 ON CONFLICT (user_id) DO UPDATE
		   SET telegram_chat_id = EXCLUDED.telegram_chat_id,
		       telegram_username = EXCLUDED.telegram_username,
		       is_active = true,
		       linked_at = now()
		 RETURNING id, linked_at`,
		link.UserID, link.TelegramChatID, link.TelegramUsername,
	).Scan(&link.ID, &link.LinkedAt)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == pqUniqueViolation {
		return telegram.ErrChatAlreadyLinked
	}
	return err
}

func (r *telegramRepo) Unlink(ctx context.Context, userID int64) error {
	_, err := r.db.Exec("UPDATE user_telegram_links SET is_active = false WHERE user_id = $1", userID)
	return err
}

func (r *telegramRepo) CreateLinkCode(ctx context.Context, code *telegram.LinkCode) error {
	err := r.db.QueryRow(
		"INSERT INTO telegram_link_codes(user_id, code, expires_at) VALUES($1, $2, $3) RETURNING id",
		code.UserID, code.Code, code.ExpiresAt,
	).Scan(&code.ID)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == pqForeignKeyViolation {
		return telegram.ErrUserNotFound
	}
	return err
}

func (r *telegramRepo) FindActiveLinkCode(ctx context.Context, code string) (telegram.LinkCode, error) {
	var lc telegram.LinkCode
	// FOR UPDATE: this always runs inside WithTransaction (see LinkAccount), so
	// it row-locks the code for the transaction's duration — closes a TOCTOU
	// race where two concurrent redemptions of the same code could both read
	// "unused" before either commits.
	err := r.db.QueryRow(
		"SELECT id, user_id, code, expires_at, used_at FROM telegram_link_codes WHERE code = $1 AND used_at IS NULL AND expires_at > now() FOR UPDATE",
		code,
	).Scan(&lc.ID, &lc.UserID, &lc.Code, &lc.ExpiresAt, &lc.UsedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return telegram.LinkCode{}, telegram.ErrLinkCodeInvalid
	}
	if err != nil {
		return telegram.LinkCode{}, err
	}
	return lc, nil
}

func (r *telegramRepo) FindActiveLinkCodeByUserID(ctx context.Context, userID int64) (telegram.LinkCode, error) {
	var lc telegram.LinkCode
	err := r.db.QueryRow(
		"SELECT id, user_id, code, expires_at, used_at FROM telegram_link_codes WHERE user_id = $1 AND used_at IS NULL AND expires_at > now() ORDER BY id DESC LIMIT 1",
		userID,
	).Scan(&lc.ID, &lc.UserID, &lc.Code, &lc.ExpiresAt, &lc.UsedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return telegram.LinkCode{}, telegram.ErrLinkCodeInvalid
	}
	if err != nil {
		return telegram.LinkCode{}, err
	}
	return lc, nil
}

func (r *telegramRepo) MarkLinkCodeUsed(ctx context.Context, id int64) error {
	_, err := r.db.Exec("UPDATE telegram_link_codes SET used_at = now() WHERE id = $1", id)
	return err
}

func (r *telegramRepo) WithTransaction(ctx context.Context, fn func(repo telegram.Repository) error) error {
	db, ok := r.db.(*sql.DB)
	if !ok {
		return fmt.Errorf("repository must be initialized with *sql.DB to start a transaction")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txRepo := &telegramRepo{db: tx}
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
