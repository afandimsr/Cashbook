package postgresql_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	tg "github.com/afandimsr/cashbook-backend/internal/domain/telegram"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTelegramRepo_SaveLink(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO user_telegram_links")).
		WithArgs(int64(1), int64(555), "afandi").
		WillReturnRows(sqlmock.NewRows([]string{"id", "linked_at"}).AddRow(int64(10), now))

	link := tg.UserTelegramLink{UserID: 1, TelegramChatID: 555, TelegramUsername: "afandi"}
	err := r.SaveLink(context.Background(), &link)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), link.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_SaveLink_ChatAlreadyLinkedToAnotherUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO user_telegram_links")).
		WithArgs(int64(2), int64(555), "someone-else").
		WillReturnError(&pq.Error{Code: "23505"})

	link := tg.UserTelegramLink{UserID: 2, TelegramChatID: 555, TelegramUsername: "someone-else"}
	err := r.SaveLink(context.Background(), &link)
	assert.ErrorIs(t, err, tg.ErrChatAlreadyLinked)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_Unlink(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE user_telegram_links SET is_active = false WHERE user_id = $1")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := r.Unlink(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_FindLinkByChatID_Found(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links")).
		WithArgs(int64(555)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "telegram_chat_id", "telegram_username", "is_active", "linked_at"}).
			AddRow(int64(10), int64(1), int64(555), "afandi", true, now))

	link, err := r.FindLinkByChatID(context.Background(), 555)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), link.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_FindLinkByChatID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links")).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "telegram_chat_id", "telegram_username", "is_active", "linked_at"}))

	_, err := r.FindLinkByChatID(context.Background(), 999)
	assert.ErrorIs(t, err, tg.ErrLinkNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_ListActiveLinks(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links WHERE is_active = true")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "telegram_chat_id", "telegram_username", "is_active", "linked_at"}).
			AddRow(int64(10), int64(1), int64(555), "afandi", true, now).
			AddRow(int64(11), int64(2), int64(556), "", true, now))

	links, err := r.ListActiveLinks(context.Background())
	assert.NoError(t, err)
	assert.Len(t, links, 2)
	assert.Equal(t, int64(555), links[0].TelegramChatID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_ListActiveLinks_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, telegram_chat_id, telegram_username, is_active, linked_at FROM user_telegram_links WHERE is_active = true")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "telegram_chat_id", "telegram_username", "is_active", "linked_at"}))

	links, err := r.ListActiveLinks(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, links)
	assert.Len(t, links, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_CreateLinkCode(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	expires := time.Now().Add(10 * time.Minute)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO telegram_link_codes")).
		WithArgs(int64(1), "ABC123", expires).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(7)))

	code := tg.LinkCode{UserID: 1, Code: "ABC123", ExpiresAt: expires}
	err := r.CreateLinkCode(context.Background(), &code)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), code.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_FindActiveLinkCode_Invalid(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, code, expires_at, used_at FROM telegram_link_codes")).
		WithArgs("EXPIRED1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "code", "expires_at", "used_at"}))

	_, err := r.FindActiveLinkCode(context.Background(), "EXPIRED1")
	assert.ErrorIs(t, err, tg.ErrLinkCodeInvalid)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_FindActiveLinkCodeByUserID_Found(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	expires := time.Now().Add(5 * time.Minute)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, code, expires_at, used_at FROM telegram_link_codes WHERE user_id = $1")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "code", "expires_at", "used_at"}).
			AddRow(int64(5), int64(1), "EXIST1", expires, nil))

	lc, err := r.FindActiveLinkCodeByUserID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "EXIST1", lc.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_FindActiveLinkCodeByUserID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, code, expires_at, used_at FROM telegram_link_codes WHERE user_id = $1")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "code", "expires_at", "used_at"}))

	_, err := r.FindActiveLinkCodeByUserID(context.Background(), 1)
	assert.ErrorIs(t, err, tg.ErrLinkCodeInvalid)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_CreateLinkCode_UserNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	expires := time.Now().Add(10 * time.Minute)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO telegram_link_codes")).
		WithArgs(int64(999999), "ABC123", expires).
		WillReturnError(&pq.Error{Code: "23503"})

	code := tg.LinkCode{UserID: 999999, Code: "ABC123", ExpiresAt: expires}
	err := r.CreateLinkCode(context.Background(), &code)
	assert.ErrorIs(t, err, tg.ErrUserNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTelegramRepo_MarkLinkCodeUsed(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTelegramRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE telegram_link_codes SET used_at = now() WHERE id = $1")).
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := r.MarkLinkCodeUsed(context.Background(), 7)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
