package postgresql_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/afandimsr/cashbook-backend/internal/domain/audit"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuditRepo_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewAuditRepo(db)

	now := time.Now()
	targetID := int64(7)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO admin_audit_log")).
		WithArgs(int64(1), "telegram.unlink", &targetID, []byte(`{"chat_id":555}`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(int64(1), now))

	entry := audit.Entry{
		AdminUserID:  1,
		Action:       "telegram.unlink",
		TargetUserID: &targetID,
		Metadata:     map[string]any{"chat_id": float64(555)},
	}
	err := r.Save(context.Background(), &entry)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), entry.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditRepo_Save_NoMetadata(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewAuditRepo(db)

	now := time.Now()
	targetID := int64(7)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO admin_audit_log")).
		WithArgs(int64(1), "telegram.generate_link_code", &targetID, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(int64(2), now))

	entry := audit.Entry{AdminUserID: 1, Action: "telegram.generate_link_code", TargetUserID: &targetID}
	err := r.Save(context.Background(), &entry)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
