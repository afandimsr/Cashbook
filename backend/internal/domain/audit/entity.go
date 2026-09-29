// Package audit is a small, deliberately narrow admin-action log — currently
// only wired into the Telegram admin endpoints (AdminGenerateLinkCode,
// AdminUnlink). Other admin actions (ResetPassword, DeleteUser, ...) are
// equally good candidates for the same helper, but extending it there is a
// separate follow-up, not bundled into this pass.
package audit

import (
	"context"
	"time"
)

type Entry struct {
	ID           int64          `json:"id"`
	AdminUserID  int64          `json:"admin_user_id"`
	Action       string         `json:"action"`
	TargetUserID *int64         `json:"target_user_id,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

type Repository interface {
	Save(ctx context.Context, entry *Entry) error
}
