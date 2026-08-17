package user_test

import (
	"errors"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMFASettings_GetSettings(t *testing.T) {
	repo := new(MockMFASettingsRepository)
	u := uc.NewMFASettingsUsecase(repo)

	repo.On("Get").Return(&user.MFASettings{ID: 1, Enforce2FA: true}, nil).Once()

	got, err := u.GetSettings()
	assert.NoError(t, err)
	assert.True(t, got.Enforce2FA)
	repo.AssertExpectations(t)
}

func TestMFASettings_UpdateSettings(t *testing.T) {
	repo := new(MockMFASettingsRepository)
	u := uc.NewMFASettingsUsecase(repo)

	repo.On("Upsert", mock.MatchedBy(func(s user.MFASettings) bool {
		return s.Enforce2FA && s.UpdatedBy == 7
	})).Return(nil).Once()

	assert.NoError(t, u.UpdateSettings(true, 7))
	repo.AssertExpectations(t)
}

func TestMFASettings_UpdateSettings_Error(t *testing.T) {
	repo := new(MockMFASettingsRepository)
	u := uc.NewMFASettingsUsecase(repo)

	repo.On("Upsert", mock.Anything).Return(errors.New("db")).Once()

	assert.Error(t, u.UpdateSettings(false, 1))
	repo.AssertExpectations(t)
}
