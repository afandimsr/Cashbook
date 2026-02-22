package user

import (
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
)

type MFASettingsUsecase struct {
	repo user.MFASettingsRepository
}

func NewMFASettingsUsecase(repo user.MFASettingsRepository) *MFASettingsUsecase {
	return &MFASettingsUsecase{repo: repo}
}

func (u *MFASettingsUsecase) GetSettings() (*user.MFASettings, error) {
	return u.repo.Get()
}

func (u *MFASettingsUsecase) UpdateSettings(enforce2FA bool, adminID int64) error {
	settings := user.MFASettings{
		Enforce2FA: enforce2FA,
		UpdatedBy:  adminID,
	}
	if err := u.repo.Upsert(settings); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
