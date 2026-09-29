package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/afandimsr/cashbook-backend/internal/domain/audit"
)

type auditRepo struct {
	db *sql.DB
}

func NewAuditRepo(db *sql.DB) audit.Repository {
	return &auditRepo{db: db}
}

func (r *auditRepo) Save(ctx context.Context, entry *audit.Entry) error {
	// Declared as `any`, not `[]byte`: a nil []byte boxed into an interface{}
	// is a typed non-nil interface (concrete type []byte, nil value), which
	// the driver would NOT treat as SQL NULL. Leaving this as a genuinely nil
	// interface when there's no metadata is what maps to NULL correctly.
	var metadataJSON any
	if entry.Metadata != nil {
		b, err := json.Marshal(entry.Metadata)
		if err != nil {
			return err
		}
		metadataJSON = b
	}

	return r.db.QueryRowContext(ctx,
		"INSERT INTO admin_audit_log(admin_user_id, action, target_user_id, metadata) VALUES($1, $2, $3, $4) RETURNING id, created_at",
		entry.AdminUserID, entry.Action, entry.TargetUserID, metadataJSON,
	).Scan(&entry.ID, &entry.CreatedAt)
}
