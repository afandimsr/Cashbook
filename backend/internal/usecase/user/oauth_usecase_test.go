package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/auth"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockOauthStateRepository struct{ mock.Mock }

func (m *MockOauthStateRepository) Save(s user.OauthState) error { return m.Called(s).Error(0) }
func (m *MockOauthStateRepository) FindByState(state string) (*user.OauthState, error) {
	args := m.Called(state)
	if v := args.Get(0); v != nil {
		return v.(*user.OauthState), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockOauthStateRepository) Update(s user.OauthState) error { return m.Called(s).Error(0) }

type MockGoogleAuth struct{ mock.Mock }

func (m *MockGoogleAuth) GetAuthURL(state string) string { return m.Called(state).String(0) }
func (m *MockGoogleAuth) ExchangeCode(code string) (*oauth2.Token, error) {
	args := m.Called(code)
	if v := args.Get(0); v != nil {
		return v.(*oauth2.Token), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockGoogleAuth) GetUserData(token *oauth2.Token) (*auth.GoogleUser, error) {
	args := m.Called(token)
	if v := args.Get(0); v != nil {
		return v.(*auth.GoogleUser), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestOAuth_GetGoogleAuthURL(t *testing.T) {
	userRepo := new(MockUserRepository)
	stateRepo := new(MockOauthStateRepository)
	ga := new(MockGoogleAuth)
	u := uc.NewOAuthUsecase(userRepo, stateRepo, ga)

	stateRepo.On("Save", mock.Anything).Return(nil).Once()
	ga.On("GetAuthURL", mock.AnythingOfType("string")).Return("https://accounts.google.com/auth").Once()

	url, err := u.GetGoogleAuthURL("1.2.3.4", "agent")
	assert.NoError(t, err)
	assert.Equal(t, "https://accounts.google.com/auth", url)
	stateRepo.AssertExpectations(t)
}

func validState() *user.OauthState {
	return &user.OauthState{
		State:     "s",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		UsedAt:    nil,
		// Empty IP/UA hashes so binding checks are skipped.
	}
}

func TestOAuth_Callback_ExistingUser(t *testing.T) {
	jwt.SetSecret("test-secret")
	userRepo := new(MockUserRepository)
	stateRepo := new(MockOauthStateRepository)
	ga := new(MockGoogleAuth)
	u := uc.NewOAuthUsecase(userRepo, stateRepo, ga)

	stateRepo.On("FindByState", "s").Return(validState(), nil).Once()
	stateRepo.On("Update", mock.Anything).Return(nil).Once()
	ga.On("ExchangeCode", "code").Return(&oauth2.Token{AccessToken: "tok"}, nil).Once()
	ga.On("GetUserData", mock.Anything).Return(&auth.GoogleUser{ID: "g1", Email: "u@e.com", Name: "U"}, nil).Once()
	userRepo.On("FindByGoogleID", "g1").Return(user.User{ID: 3, Email: "u@e.com", Name: "U", Roles: []string{"USER"}}, nil).Once()

	token, err := u.HandleGoogleCallback("code", "s", "1.2.3.4", "agent")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	stateRepo.AssertExpectations(t)
}

func TestOAuth_Callback_InvalidState(t *testing.T) {
	u := uc.NewOAuthUsecase(new(MockUserRepository), func() *MockOauthStateRepository {
		m := new(MockOauthStateRepository)
		m.On("FindByState", "s").Return(nil, errors.New("not found")).Once()
		return m
	}(), new(MockGoogleAuth))

	_, err := u.HandleGoogleCallback("code", "s", "ip", "ua")
	assert.Error(t, err)
}

func TestOAuth_Callback_Expired(t *testing.T) {
	stateRepo := new(MockOauthStateRepository)
	expired := &user.OauthState{State: "s", ExpiresAt: time.Now().Add(-time.Minute)}
	stateRepo.On("FindByState", "s").Return(expired, nil).Once()

	u := uc.NewOAuthUsecase(new(MockUserRepository), stateRepo, new(MockGoogleAuth))
	_, err := u.HandleGoogleCallback("code", "s", "ip", "ua")
	assert.Error(t, err)
}

func TestOAuth_Callback_AlreadyUsed(t *testing.T) {
	stateRepo := new(MockOauthStateRepository)
	used := time.Now()
	st := &user.OauthState{State: "s", ExpiresAt: time.Now().Add(time.Minute), UsedAt: &used}
	stateRepo.On("FindByState", "s").Return(st, nil).Once()

	u := uc.NewOAuthUsecase(new(MockUserRepository), stateRepo, new(MockGoogleAuth))
	_, err := u.HandleGoogleCallback("code", "s", "ip", "ua")
	assert.Error(t, err)
}

func TestOAuth_Callback_NewUser(t *testing.T) {
	jwt.SetSecret("test-secret")
	userRepo := new(MockUserRepository)
	stateRepo := new(MockOauthStateRepository)
	ga := new(MockGoogleAuth)
	u := uc.NewOAuthUsecase(userRepo, stateRepo, ga)

	stateRepo.On("FindByState", "s").Return(validState(), nil).Once()
	stateRepo.On("Update", mock.Anything).Return(nil).Once()
	ga.On("ExchangeCode", "code").Return(&oauth2.Token{AccessToken: "tok"}, nil).Once()
	ga.On("GetUserData", mock.Anything).Return(&auth.GoogleUser{ID: "g9", Email: "new@e.com", Name: "New"}, nil).Once()

	// Not found by google id nor email → create, then re-fetch by google id.
	userRepo.On("FindByGoogleID", "g9").Return(user.User{}, errors.New("nf")).Once()
	userRepo.On("FindByEmail", "new@e.com").Return(user.User{}, errors.New("nf")).Once()
	userRepo.On("Save", mock.MatchedBy(func(usr user.User) bool { return usr.GoogleID == "g9" })).Return(nil).Once()
	userRepo.On("FindByGoogleID", "g9").Return(user.User{ID: 10, Email: "new@e.com", Name: "New", Roles: []string{"USER"}}, nil).Once()

	token, err := u.HandleGoogleCallback("code", "s", "1.2.3.4", "agent")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	userRepo.AssertExpectations(t)
}
