package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/totp"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	otptotp "github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockMFABackupCodeRepository struct {
	mock.Mock
}

func (m *MockMFABackupCodeRepository) SaveBatch(userID int64, hashes []string) error {
	return m.Called(userID, hashes).Error(0)
}
func (m *MockMFABackupCodeRepository) FindByUserID(userID int64) ([]user.MFABackupCode, error) {
	args := m.Called(userID)
	return args.Get(0).([]user.MFABackupCode), args.Error(1)
}
func (m *MockMFABackupCodeRepository) MarkUsed(id int64) error   { return m.Called(id).Error(0) }
func (m *MockMFABackupCodeRepository) DeleteByUserID(id int64) error { return m.Called(id).Error(0) }

func TestTwoFA_Setup(t *testing.T) {
	userRepo := new(MockUserRepository)
	backupRepo := new(MockMFABackupCodeRepository)
	u := uc.NewTwoFAUsecase(userRepo, backupRepo)

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, Email: "a@b.com", TOTPEnabled: false}, nil).Once()
	userRepo.On("UpdateTOTPSecret", mock.MatchedBy(func(usr user.User) bool {
		return usr.TOTPSecret != ""
	})).Return(nil).Once()

	resp, err := u.Setup(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Secret)
	assert.NotEmpty(t, resp.QRCode)
	userRepo.AssertExpectations(t)
}

func TestTwoFA_Setup_AlreadyEnabled(t *testing.T) {
	userRepo := new(MockUserRepository)
	u := uc.NewTwoFAUsecase(userRepo, new(MockMFABackupCodeRepository))

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPEnabled: true}, nil).Once()

	_, err := u.Setup(1)
	assert.Error(t, err)
}

func TestTwoFA_VerifySetup(t *testing.T) {
	secret, _, err := totp.GenerateSecret("a@b.com")
	assert.NoError(t, err)
	code, err := otptotp.GenerateCode(secret, time.Now())
	assert.NoError(t, err)

	userRepo := new(MockUserRepository)
	u := uc.NewTwoFAUsecase(userRepo, new(MockMFABackupCodeRepository))

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPSecret: secret}, nil).Once()
	userRepo.On("EnableTOTP", mock.Anything).Return(nil).Once()

	assert.NoError(t, u.VerifySetup(1, code))
	userRepo.AssertExpectations(t)
}

func TestTwoFA_VerifySetup_InvalidCode(t *testing.T) {
	secret, _, _ := totp.GenerateSecret("a@b.com")
	userRepo := new(MockUserRepository)
	u := uc.NewTwoFAUsecase(userRepo, new(MockMFABackupCodeRepository))

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPSecret: secret}, nil).Once()

	assert.Error(t, u.VerifySetup(1, "000000"))
	userRepo.AssertNotCalled(t, "EnableTOTP", mock.Anything)
}

func TestTwoFA_VerifySetup_NotInitiated(t *testing.T) {
	userRepo := new(MockUserRepository)
	u := uc.NewTwoFAUsecase(userRepo, new(MockMFABackupCodeRepository))

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPSecret: ""}, nil).Once()

	assert.Error(t, u.VerifySetup(1, "123456"))
}

func TestTwoFA_Disable(t *testing.T) {
	userRepo := new(MockUserRepository)
	backupRepo := new(MockMFABackupCodeRepository)
	u := uc.NewTwoFAUsecase(userRepo, backupRepo)

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPEnabled: true, TOTPSecret: "x"}, nil).Once()
	userRepo.On("DisableTOTP", mock.MatchedBy(func(usr user.User) bool {
		return !usr.TOTPEnabled && usr.TOTPSecret == ""
	})).Return(nil).Once()
	backupRepo.On("DeleteByUserID", int64(1)).Return(nil).Once()

	assert.NoError(t, u.Disable(1))
	userRepo.AssertExpectations(t)
	backupRepo.AssertExpectations(t)
}

func TestTwoFA_VerifyLogin(t *testing.T) {
	jwt.SetSecret("test-secret")
	secret, _, _ := totp.GenerateSecret("a@b.com")
	code, _ := otptotp.GenerateCode(secret, time.Now())
	tempToken, _ := jwt.GenerateTempToken(1, "a@b.com", "verify")

	userRepo := new(MockUserRepository)
	u := uc.NewTwoFAUsecase(userRepo, new(MockMFABackupCodeRepository))
	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, Email: "a@b.com", TOTPSecret: secret}, nil).Once()

	token, err := u.VerifyLogin(tempToken, code)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestTwoFA_VerifyLogin_BadToken(t *testing.T) {
	jwt.SetSecret("test-secret")
	u := uc.NewTwoFAUsecase(new(MockUserRepository), new(MockMFABackupCodeRepository))

	_, err := u.VerifyLogin("garbage", "123456")
	assert.Error(t, err)
}

func TestTwoFA_GenerateBackupCodes(t *testing.T) {
	userRepo := new(MockUserRepository)
	backupRepo := new(MockMFABackupCodeRepository)
	u := uc.NewTwoFAUsecase(userRepo, backupRepo)

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPEnabled: true}, nil).Once()
	backupRepo.On("DeleteByUserID", int64(1)).Return(nil).Once()
	backupRepo.On("SaveBatch", int64(1), mock.MatchedBy(func(h []string) bool {
		return len(h) == 10
	})).Return(nil).Once()

	codes, err := u.GenerateBackupCodes(1)
	assert.NoError(t, err)
	assert.Len(t, codes, 10)
	backupRepo.AssertExpectations(t)
}

func TestTwoFA_GenerateBackupCodes_NotEnabled(t *testing.T) {
	userRepo := new(MockUserRepository)
	u := uc.NewTwoFAUsecase(userRepo, new(MockMFABackupCodeRepository))

	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, TOTPEnabled: false}, nil).Once()

	_, err := u.GenerateBackupCodes(1)
	assert.Error(t, err)
}

func TestTwoFA_VerifyBackupCode(t *testing.T) {
	jwt.SetSecret("test-secret")
	tempToken, _ := jwt.GenerateTempToken(1, "a@b.com", "verify")
	plain := "abcd-1234"
	hash, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)

	userRepo := new(MockUserRepository)
	backupRepo := new(MockMFABackupCodeRepository)
	u := uc.NewTwoFAUsecase(userRepo, backupRepo)

	backupRepo.On("FindByUserID", int64(1)).Return([]user.MFABackupCode{
		{ID: 5, CodeHash: string(hash), UsedAt: nil},
	}, nil).Once()
	backupRepo.On("MarkUsed", int64(5)).Return(nil).Once()
	userRepo.On("FindByID", int64(1)).Return(user.User{ID: 1, Email: "a@b.com"}, nil).Once()

	token, err := u.VerifyBackupCode(tempToken, plain)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	backupRepo.AssertExpectations(t)
}

func TestTwoFA_VerifyBackupCode_NoMatch(t *testing.T) {
	jwt.SetSecret("test-secret")
	tempToken, _ := jwt.GenerateTempToken(1, "a@b.com", "verify")
	hash, _ := bcrypt.GenerateFromPassword([]byte("other-code"), bcrypt.DefaultCost)

	backupRepo := new(MockMFABackupCodeRepository)
	u := uc.NewTwoFAUsecase(new(MockUserRepository), backupRepo)

	backupRepo.On("FindByUserID", int64(1)).Return([]user.MFABackupCode{
		{ID: 5, CodeHash: string(hash), UsedAt: nil},
	}, nil).Once()

	_, err := u.VerifyBackupCode(tempToken, "wrong-code")
	assert.Error(t, err)
	backupRepo.AssertNotCalled(t, "MarkUsed", mock.Anything)
}

func TestTwoFA_VerifyBackupCode_SkipsUsed(t *testing.T) {
	jwt.SetSecret("test-secret")
	tempToken, _ := jwt.GenerateTempToken(1, "a@b.com", "verify")
	plain := "abcd-1234"
	hash, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	used := time.Now()

	backupRepo := new(MockMFABackupCodeRepository)
	u := uc.NewTwoFAUsecase(new(MockUserRepository), backupRepo)

	backupRepo.On("FindByUserID", int64(1)).Return([]user.MFABackupCode{
		{ID: 5, CodeHash: string(hash), UsedAt: &used}, // already used → skipped
	}, nil).Once()

	_, err := u.VerifyBackupCode(tempToken, plain)
	assert.Error(t, err)
}

var _ = errors.Is
