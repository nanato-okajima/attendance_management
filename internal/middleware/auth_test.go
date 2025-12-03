package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/nanato-okajima/attendance_management/config"
)

func TestGenerateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expiration  string
		checkClaims bool
	}{
		{
			name:        "トークン生成成功",
			expiration:  "1h",
			checkClaims: true,
		},
		{
			name:        "無効な有効期限（デフォルト24時間適用）",
			expiration:  "invalid",
			checkClaims: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := &config.JWTConfig{
				Secret:     "test-secret",
				Expiration: tt.expiration,
			}

			token, err := GenerateToken(1, "test@example.com", "employee", cfg)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			if tt.checkClaims {
				// トークンをパースして検証
				claims, err := ParseToken(token, cfg)
				assert.NoError(t, err)
				assert.Equal(t, 1, claims.EmployeeNumber)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, "employee", claims.Role)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	t.Parallel()

	cfg := &config.JWTConfig{
		Secret:     "test-secret",
		Expiration: "1h",
	}

	tests := []struct {
		name        string
		setupToken  func() string
		expectError bool
		checkClaims func(*testing.T, *Claims)
	}{
		{
			name: "パース成功",
			setupToken: func() string {
				token, _ := GenerateToken(1, "test@example.com", "admin", cfg)
				return token
			},
			expectError: false,
			checkClaims: func(t *testing.T, claims *Claims) {
				assert.Equal(t, 1, claims.EmployeeNumber)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, "admin", claims.Role)
			},
		},
		{
			name: "無効なトークン",
			setupToken: func() string {
				return "invalid.token.here"
			},
			expectError: true,
			checkClaims: nil,
		},
		{
			name: "期限切れトークン",
			setupToken: func() string {
				// 既に期限切れのトークンを手動で作成
				expirationTime := time.Now().Add(-1 * time.Hour) // 1時間前
				claims := &Claims{
					EmployeeNumber: 1,
					Email:          "test@example.com",
					Role:           "employee",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(expirationTime),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.Secret))
				return tokenString
			},
			expectError: true,
			checkClaims: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tokenString := tt.setupToken()
			claims, err := ParseToken(tokenString, cfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				if tt.checkClaims != nil {
					tt.checkClaims(t, claims)
				}
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	cfg := &config.JWTConfig{
		Secret:     "test-secret",
		Expiration: "1h",
	}

	tests := []struct {
		name           string
		setupRequest   func(*gin.Context)
		expectedStatus int
		checkContext   func(*testing.T, *gin.Context)
		expectAborted  bool
	}{
		{
			name: "認証成功",
			setupRequest: func(c *gin.Context) {
				token, _ := GenerateToken(1, "test@example.com", "employee", cfg)
				c.Request.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				assert.Equal(t, 1, c.GetInt("employee_number"))
				assert.Equal(t, "test@example.com", c.GetString("email"))
				assert.Equal(t, "employee", c.GetString("role"))
			},
			expectAborted: false,
		},
		{
			name: "認証ヘッダーなし",
			setupRequest: func(c *gin.Context) {
				// ヘッダーを設定しない
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   nil,
			expectAborted:  true,
		},
		{
			name: "無効なトークン",
			setupRequest: func(c *gin.Context) {
				c.Request.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   nil,
			expectAborted:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			tt.setupRequest(c)

			middleware := AuthMiddleware(cfg)
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkContext != nil {
				tt.checkContext(t, c)
			}
			assert.Equal(t, tt.expectAborted, c.IsAborted())
		})
	}
}

func TestAdminOnly(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		role           string
		setRole        bool
		expectedStatus int
		expectAborted  bool
	}{
		{
			name:           "管理者権限あり",
			role:           "admin",
			setRole:        true,
			expectedStatus: http.StatusOK,
			expectAborted:  false,
		},
		{
			name:           "管理者権限なし（一般ユーザー）",
			role:           "employee",
			setRole:        true,
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
		{
			name:           "ロール未設定",
			role:           "",
			setRole:        false,
			expectedStatus: http.StatusForbidden,
			expectAborted:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setRole {
				c.Set("role", tt.role)
			}

			middleware := AdminOnly()
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectAborted, c.IsAborted())
		})
	}
}
