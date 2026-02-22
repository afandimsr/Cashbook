package user

import (
	"log"

	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	repo            user.UserRepository
	authService     user.AuthService
	mfaSettingsRepo user.MFASettingsRepository
}

func New(repo user.UserRepository, authService user.AuthService) *Usecase {
	return &Usecase{
		repo:        repo,
		authService: authService,
	}
}

func (u *Usecase) SetMFASettingsRepo(repo user.MFASettingsRepository) {
	u.mfaSettingsRepo = repo
}

func (u *Usecase) GetAll(page, limit int) ([]user.User, error) {
	offset := (page - 1) * limit
	return u.repo.FindAll(limit, offset)
}

func (u *Usecase) GetByID(id int64) (user.User, error) {
	return u.repo.FindByID(id)
}

func (u *Usecase) Create(newUser user.User) error {
	if newUser.Email == "" {
		return apperror.BadRequest("email is required", nil)
	}
	if newUser.Password == "" {
		return apperror.BadRequest("password is required", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Internal(err)
	}
	newUser.Password = string(hashedPassword)

	if err := u.repo.Save(newUser); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) Update(id int64, updatedUser user.User) error {
	if updatedUser.Email == "" {
		return apperror.BadRequest("email is required", nil)
	}

	// Check if user exists
	existingUser, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	existingUser.Name = updatedUser.Name
	existingUser.Email = updatedUser.Email
	existingUser.IsActive = updatedUser.IsActive
	existingUser.Roles = updatedUser.Roles

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return apperror.Internal(err)
		}
		existingUser.Password = string(hashedPassword)
	}

	if err := u.repo.Update(existingUser); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) Delete(id int64) error {
	// Check if user exists
	if _, err := u.repo.FindByID(id); err != nil {
		return err
	}

	if err := u.repo.Delete(id); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) Login(email, password string) (*user.LoginResponse, error) {
	// 1. Find user by email
	existingUser, err := u.repo.FindByEmail(email)
	if err != nil {
		log.Printf("Login: user not found by email=%s err=%v", email, err)
		return nil, apperror.Unauthorized("invalid credentials [1]", nil)
	}
	log.Printf("Login: found user id=%d email=%s roles=%v google_id=%s", existingUser.ID, existingUser.Email, existingUser.Roles, existingUser.GoogleID)

	if !existingUser.IsActive {
		return nil, apperror.Unauthorized("invalid credentials [2]", nil)
	}

	// 2. Authenticate
	authenticated := false
	if u.authService != nil {
		isAuth, err := u.authService.Login(email, password)
		if err == nil && isAuth {
			authenticated = true
		}
	}

	if !authenticated {
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
			log.Printf("Login: bcrypt compare failed for user id=%d err=%v", existingUser.ID, err)
			return nil, apperror.Unauthorized("invalid credentials [2]", nil)
		}
		log.Printf("Login: password verified for user id=%d", existingUser.ID)
	}

	// 3. Check if totp secret null return to register 2fa
	if existingUser.TOTPSecret == "" {
		tempToken, err := jwt.GenerateTempToken(existingUser.ID, existingUser.Email, "setup")
		if err != nil {
			return nil, apperror.Internal(err)
		}
		return &user.LoginResponse{
			Requires2FA: true,
			TempToken:   tempToken,
		}, nil
	}

	// 4. Check if 2FA is enabled for this user
	if existingUser.TOTPEnabled {
		tempToken, err := jwt.GenerateTempToken(existingUser.ID, existingUser.Email, "verify")
		if err != nil {
			return nil, apperror.Internal(err)
		}
		return &user.LoginResponse{
			Requires2FA: true,
			TempToken:   tempToken,
		}, nil
	}

	// 4. Check if MFA is enforced system-wide (user hasn't set up 2FA yet)
	if u.mfaSettingsRepo != nil {
		settings, err := u.mfaSettingsRepo.Get()
		if err == nil && settings != nil && settings.Enforce2FA && !existingUser.TOTPEnabled {
			// Return a token but signal that 2FA setup is required
			token, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
			if err != nil {
				return nil, apperror.Internal(err)
			}
			return &user.LoginResponse{
				Token: token,
			}, nil
		}
	}

	// 5. Generate full Token
	token, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &user.LoginResponse{
		Token: token,
	}, nil
}

func (u *Usecase) ResetPassword(id int64, newPassword string) error {
	if err := u.validatePassword(newPassword); err != nil {
		return err
	}

	existingUser, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Internal(err)
	}

	existingUser.Password = string(hashedPassword)

	if err := u.repo.Update(existingUser); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) validatePassword(password string) error {
	if len(password) < 8 {
		return apperror.BadRequest("password must be at least 8 characters long", nil)
	}

	var hasUpper, hasLower, hasNumber, hasSymbol bool
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*' || char == '(' || char == ')' || char == '-' || char == '_' || char == '+' || char == '=':
			hasSymbol = true
		}
	}

	if !hasUpper {
		return apperror.BadRequest("password must contain at least one uppercase letter", nil)
	}
	if !hasLower {
		return apperror.BadRequest("password must contain at least one lowercase letter", nil)
	}
	if !hasNumber {
		return apperror.BadRequest("password must contain at least one number", nil)
	}
	if !hasSymbol {
		return apperror.BadRequest("password must contain at least one special character", nil)
	}

	return nil
}
