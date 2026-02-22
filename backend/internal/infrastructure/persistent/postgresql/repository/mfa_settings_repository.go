package postgresql

import (
	"database/sql"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
)

type mfaSettingsRepo struct {
	db *sql.DB
}

func NewMFASettingsRepo(db *sql.DB) user.MFASettingsRepository {
	return &mfaSettingsRepo{db: db}
}

func (r *mfaSettingsRepo) Get() (*user.MFASettings, error) {
	var s user.MFASettings
	var updatedBy sql.NullInt64
	err := r.db.QueryRow("SELECT id, enforce_2fa, updated_by, updated_at FROM mfa_settings ORDER BY id LIMIT 1").Scan(&s.ID, &s.Enforce2FA, &updatedBy, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default settings if no row exists
			return &user.MFASettings{Enforce2FA: false}, nil
		}
		return nil, err
	}
	if updatedBy.Valid {
		s.UpdatedBy = updatedBy.Int64
	}
	return &s, nil
}

func (r *mfaSettingsRepo) Upsert(settings user.MFASettings) error {
	_, err := r.db.Exec(`
		INSERT INTO mfa_settings (id, enforce_2fa, updated_by, updated_at) 
		VALUES (1, $1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET enforce_2fa = $1, updated_by = $2, updated_at = $3
	`, settings.Enforce2FA, settings.UpdatedBy, time.Now())
	return err
}
