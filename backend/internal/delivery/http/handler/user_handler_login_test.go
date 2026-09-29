package handler_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/config"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/pkg/ratelimit"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// stubUserRepo backs a real *uc.Usecase (UserHandler's usecase field is a
// concrete struct, not an interface) so the login-lockout logic under test
// runs through the actual Login implementation, not a mock of it.
type stubUserRepo struct {
	user user.User
}

func (s *stubUserRepo) FindAll(limit, offset int) ([]user.User, error) { return nil, nil }
func (s *stubUserRepo) FindByID(id int64) (user.User, error)           { return user.User{}, nil }
func (s *stubUserRepo) FindByEmail(email string) (user.User, error)    { return s.user, nil }
func (s *stubUserRepo) FindByGoogleID(googleID string) (user.User, error) {
	return user.User{}, nil
}
func (s *stubUserRepo) Save(u user.User) error             { return nil }
func (s *stubUserRepo) Update(u user.User) error           { return nil }
func (s *stubUserRepo) Delete(id int64) error              { return nil }
func (s *stubUserRepo) DisableTOTP(u user.User) error      { return nil }
func (s *stubUserRepo) UpdateTOTPSecret(u user.User) error { return nil }
func (s *stubUserRepo) EnableTOTP(u user.User) error       { return nil }

func newLoginHandler(t *testing.T, maxAttempts int) *handler.UserHandler {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	repo := &stubUserRepo{user: user.User{
		ID:          1,
		Email:       "a@b.com",
		Password:    string(hashed),
		IsActive:    true,
		TOTPSecret:  "already-set-up", // skip the "must set up 2FA" branch
		TOTPEnabled: false,            // skip the "must verify 2FA" branch
	}}
	usecase := uc.New(repo, nil)
	limiter := ratelimit.New(maxAttempts, time.Minute)
	return handler.New(&config.Config{}, usecase, nil, limiter)
}

func TestLogin_Success(t *testing.T) {
	h := newLoginHandler(t, 3)
	r := newRouter()
	r.POST("/login", h.Login)

	w := doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"correct-password"}`)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin_WrongPassword_Rejected(t *testing.T) {
	h := newLoginHandler(t, 3)
	r := newRouter()
	r.POST("/login", h.Login)

	w := doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"wrong"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_LockedOutAfterMaxFailedAttempts(t *testing.T) {
	h := newLoginHandler(t, 3)
	r := newRouter()
	r.POST("/login", h.Login)

	for i := 0; i < 3; i++ {
		w := doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"wrong"}`)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}

	// 4th attempt, even with the CORRECT password, must be blocked by the
	// per-account lockout before the usecase is even consulted.
	w := doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"correct-password"}`)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestLogin_SuccessResetsFailureCount(t *testing.T) {
	h := newLoginHandler(t, 2)
	r := newRouter()
	r.POST("/login", h.Login)

	w := doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"wrong"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"correct-password"}`)
	assert.Equal(t, http.StatusOK, w.Code)

	// Failure count should have reset on success, so this single failure
	// must not trip the 2-attempt lockout.
	w = doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"wrong"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_LockoutIsPerAccount(t *testing.T) {
	h := newLoginHandler(t, 1)
	r := newRouter()
	r.POST("/login", h.Login)

	w := doJSON(r, http.MethodPost, "/login", `{"email":"a@b.com","password":"wrong"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// A different email is unaffected by the first account's lockout.
	w = doJSON(r, http.MethodPost, "/login", `{"email":"someone-else@b.com","password":"wrong"}`)
	assert.NotEqual(t, http.StatusTooManyRequests, w.Code)
}
