package user

type UserRepository interface {
	FindAll(limit, offset int) ([]User, error)
	FindByID(id int64) (User, error)
	FindByEmail(email string) (User, error)
	FindByGoogleID(googleID string) (User, error)
	Save(user User) error
	Update(user User) error
	Delete(id int64) error
	DisableTOTP(user User) error
	UpdateTOTPSecret(user User) error
	EnableTOTP(user User) error
}

type AuthService interface {
	Login(email, password string) (bool, error)
}

type OauthStateRepository interface {
	Save(state OauthState) error
	FindByState(state string) (*OauthState, error)
	Update(state OauthState) error
}

type MFASettingsRepository interface {
	Get() (*MFASettings, error)
	Upsert(settings MFASettings) error
}

type MFABackupCodeRepository interface {
	SaveBatch(userID int64, codeHashes []string) error
	FindByUserID(userID int64) ([]MFABackupCode, error)
	MarkUsed(id int64) error
	DeleteByUserID(userID int64) error
}
