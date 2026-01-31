package postgresql

import (
	"database/sql"
	"errors"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
)

type oauthStateRepo struct {
	db *sql.DB
}

func NewOauthStateRepo(db *sql.DB) user.OauthStateRepository {
	return &oauthStateRepo{
		db: db,
	}
}

func (r *oauthStateRepo) Save(state user.OauthState) error {
	query := `INSERT INTO oauth_states (id, state, provider, client_id, redirect_uri, ip_hash, user_agent_hash, expires_at, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query,
		state.ID,
		state.State,
		state.Provider,
		state.ClientID,
		state.RedirectURI,
		state.IPHash,
		state.UserAgentHash,
		state.ExpiresAt,
		state.CreatedAt,
	)
	return err
}

func (r *oauthStateRepo) FindByState(state string) (*user.OauthState, error) {
	query := `SELECT id, state, provider, client_id, redirect_uri, ip_hash, user_agent_hash, created_at, expires_at, used_at
			  FROM oauth_states WHERE state = $1`

	var s user.OauthState
	var clientID, redirectURI, ipHash, userAgentHash sql.NullString
	var usedAt sql.NullTime

	err := r.db.QueryRow(query, state).Scan(
		&s.ID,
		&s.State,
		&s.Provider,
		&clientID,
		&redirectURI,
		&ipHash,
		&userAgentHash,
		&s.CreatedAt,
		&s.ExpiresAt,
		&usedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("state not found")
		}
		return nil, err
	}

	if clientID.Valid {
		s.ClientID = clientID.String
	}
	if redirectURI.Valid {
		s.RedirectURI = redirectURI.String
	}
	if ipHash.Valid {
		s.IPHash = ipHash.String
	}
	if userAgentHash.Valid {
		s.UserAgentHash = userAgentHash.String
	}
	if usedAt.Valid {
		t := usedAt.Time
		s.UsedAt = &t
	}

	return &s, nil
}

func (r *oauthStateRepo) Update(state user.OauthState) error {
	query := `UPDATE oauth_states 
			  SET is_used = $2, used_at = $3
			  WHERE id = $1`

	// The plan had "is_used", but the migration uses "used_at" as a timestamp which implies usage.
	// But let's check the migration again. The user asked for "used_at TIMESTAMP NULL".
	// The implementation plan had "is_used" in the first draft but the final one removed it in favor of "used_at".
	// However, I should make sure I match the migration.

	// Re-reading migration: `used_at TIMESTAMP NULL`. There is NO `is_used` column in the final migration I wrote.
	// So I should only update `used_at`.

	query = `UPDATE oauth_states SET used_at = $2 WHERE id = $1`

	_, err := r.db.Exec(query, state.ID, state.UsedAt)
	return err
}
