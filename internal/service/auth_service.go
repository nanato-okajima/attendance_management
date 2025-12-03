package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/domain/service"
	"github.com/nanato-okajima/attendance_management/internal/middleware"
	"github.com/nanato-okajima/attendance_management/internal/models"
)

type AuthService interface {
	Login(email, password string) (string, *models.User, error)
	Register(user *models.User, password string) error
}

type authService struct {
	userRepo service.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo service.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := middleware.GenerateToken(user.EmployeeNumber, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return "", nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	return token, user, nil
}

func (s *authService) Register(user *models.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Create(user)
}
