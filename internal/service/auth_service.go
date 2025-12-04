package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/domain/user"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
	"github.com/nanato-okajima/attendance_management/internal/middleware"
)

type AuthService interface {
	Login(email, password string) (string, *entity.User, error)
	Register(user *entity.User, password string) error
}

type authService struct {
	userRepo user.Repository
	cfg      *config.Config
}

func NewAuthService(userRepo user.Repository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Login(email, password string) (string, *entity.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New(errors.ErrCodeInvalidCredentials, "Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New(errors.ErrCodeInvalidCredentials, "Invalid credentials")
	}

	token, err := middleware.GenerateToken(user.EmployeeNumber, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return "", nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	return token, user, nil
}

func (s *authService) Register(user *entity.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Create(user)
}
