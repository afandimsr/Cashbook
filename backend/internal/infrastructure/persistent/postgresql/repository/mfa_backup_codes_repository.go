package postgresql

import (
	"database/sql"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
)

type mfaBackupCodeRepo struct {
	db *sql.DB
}

func NewMFABackupCodeRepo(db *sql.DB) user.MFABackupCodeRepository {
	return &mfaBackupCodeRepo{db: db}
}

func (r *mfaBackupCodeRepo) SaveBatch(userID int64, codeHashes []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, hash := range codeHashes {
		_, err := tx.Exec("INSERT INTO mfa_backup_codes (user_id, code_hash, created_at) VALUES ($1, $2, $3)", userID, hash, time.Now())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *mfaBackupCodeRepo) FindByUserID(userID int64) ([]user.MFABackupCode, error) {
	rows, err := r.db.Query("SELECT id, user_id, code_hash, used_at, created_at FROM mfa_backup_codes WHERE user_id = $1 ORDER BY created_at", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []user.MFABackupCode
	for rows.Next() {
		var c user.MFABackupCode
		if err := rows.Scan(&c.ID, &c.UserID, &c.CodeHash, &c.UsedAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		codes = append(codes, c)
	}
	return codes, nil
}

func (r *mfaBackupCodeRepo) MarkUsed(id int64) error {
	now := time.Now()
	_, err := r.db.Exec("UPDATE mfa_backup_codes SET used_at = $1 WHERE id = $2", now, id)
	return err
}

func (r *mfaBackupCodeRepo) DeleteByUserID(userID int64) error {
	_, err := r.db.Exec("DELETE FROM mfa_backup_codes WHERE user_id = $1", userID)
	return err
}
