package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nanato-okajima/attendance_management/internal/entity"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string) (string, *entity.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*entity.User), args.Error(2)
}

func (m *MockAuthService) Register(user *entity.User, password string) error {
	args := m.Called(user, password)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "ログイン成功",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mockService *MockAuthService) {
				expectedUser := &entity.User{
					ID:             1,
					EmployeeNumber: 1,
					Email:          "test@example.com",
					Role:           "employee",
				}
				mockService.On("Login", "test@example.com", "password123").
					Return("test-token", expectedUser, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "test-token", response["token"])
				assert.Equal(t, "Bearer", response["type"])
			},
		},
		{
			name: "認証情報不正",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mockService *MockAuthService) {
				mockService.On("Login", "test@example.com", "wrongpassword").
					Return("", nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name:           "不正なJSON",
			requestBody:    "invalid json",
			setupMock:      func(mockService *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var jsonBody []byte
			if str, ok := tt.requestBody.(string); ok {
				jsonBody = []byte(str)
			} else {
				jsonBody, _ = json.Marshal(tt.requestBody)
			}

			c.Request = httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    RegisterRequest
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "登録成功",
			requestBody: RegisterRequest{
				EmployeeNumber: 1,
				Email:          "test@example.com",
				Password:       "password123",
				Role:           "employee",
			},
			setupMock: func(mockService *MockAuthService) {
				mockService.On("Register", mock.AnythingOfType("*entity.User"), "password123").
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user registered successfully", response["message"])
			},
		},
		{
			name: "サービスエラー",
			requestBody: RegisterRequest{
				EmployeeNumber: 1,
				Email:          "test@example.com",
				Password:       "password123",
				Role:           "employee",
			},
			setupMock: func(mockService *MockAuthService) {
				mockService.On("Register", mock.AnythingOfType("*entity.User"), "password123").
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
		{
			name: "パスワード不正",
			requestBody: RegisterRequest{
				EmployeeNumber: 1,
				Email:          "test@example.com",
				Password:       "short", // 8文字未満
				Role:           "employee",
			},
			setupMock:      func(mockService *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "ロール不正",
			requestBody: RegisterRequest{
				EmployeeNumber: 1,
				Email:          "test@example.com",
				Password:       "password123",
				Role:           "invalid_role", // admin/employee以外
			},
			setupMock:      func(mockService *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.requestBody)

			c.Request = httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
			mockService.AssertExpectations(t)
		})
	}
}
