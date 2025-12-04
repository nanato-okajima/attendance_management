package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/entity"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmployeeNumber(employeeNumber int) (*entity.User, error) {
	args := m.Called(employeeNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: "24h",
		},
	}

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		email         string
		password      string
		setupMock     func(*MockUserRepository)
		hasError      bool
		expectedError string
		checkToken    bool
		checkUser     bool
	}{
		{
			name:     "ログイン成功",
			email:    "test@example.com",
			password: password,
			setupMock: func(mockRepo *MockUserRepository) {
				expectedUser := &entity.User{
					ID:             1,
					EmployeeNumber: 1,
					Email:          "test@example.com",
					PasswordHash:   string(hashedPassword),
					Role:           "employee",
					IsActive:       true,
				}
				mockRepo.On("FindByEmail", "test@example.com").Return(expectedUser, nil)
				mockRepo.On("UpdateLastLogin", uint(1)).Return(nil)
			},
			hasError:   false,
			checkToken: true,
			checkUser:  true,
		},
		{
			name:     "メールアドレス不正",
			email:    "wrong@example.com",
			password: password,
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("FindByEmail", "wrong@example.com").Return(nil, errors.New("user not found"))
			},
			hasError:      true,
			expectedError: "[AUTH_004] Invalid credentials",
			checkToken:    false,
			checkUser:     false,
		},
		{
			name:     "パスワード不正",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func(mockRepo *MockUserRepository) {
				expectedUser := &entity.User{
					ID:             1,
					EmployeeNumber: 1,
					Email:          "test@example.com",
					PasswordHash:   string(hashedPassword),
					Role:           "employee",
					IsActive:       true,
				}
				mockRepo.On("FindByEmail", "test@example.com").Return(expectedUser, nil)
			},
			hasError:      true,
			expectedError: "[AUTH_004] Invalid credentials",
			checkToken:    false,
			checkUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			service := NewAuthService(mockRepo, cfg)

			token, user, err := service.Login(tt.email, tt.password)

			if tt.hasError {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Nil(t, user)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if tt.checkToken {
					assert.NotEmpty(t, token)
				}
				if tt.checkUser {
					assert.NotNil(t, user)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: "24h",
		},
	}

	tests := []struct {
		name          string
		user          *entity.User
		password      string
		setupMock     func(*MockUserRepository)
		hasError      bool
		expectedError string
		checkPassword bool
	}{
		{
			name: "登録成功",
			user: &entity.User{
				EmployeeNumber: 1,
				Email:          "test@example.com",
				Role:           "employee",
			},
			password: "password123",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)
			},
			hasError:      false,
			checkPassword: true,
		},
		{
			name: "リポジトリエラー",
			user: &entity.User{
				EmployeeNumber: 1,
				Email:          "test@example.com",
				Role:           "employee",
			},
			password: "password123",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))
			},
			hasError:      true,
			expectedError: "database error",
			checkPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			service := NewAuthService(mockRepo, cfg)

			err := service.Register(tt.user, tt.password)

			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if tt.checkPassword {
					assert.NotEmpty(t, tt.user.PasswordHash)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
