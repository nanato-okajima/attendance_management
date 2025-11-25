# 従業員用出勤管理システム 実装計画書

**文書バージョン**: 1.0  
**作成日**: 2025-11-25  
**対応要件仕様書**: v3.0

---

## 1. 実装方針

### 1.1 基本方針

- **段階的実装**: 機能を優先度に応じてフェーズ分けし、段階的にリリース
- **MVP優先**: 最小限の機能で動作するシステムを早期に構築
- **テスト駆動**: 各機能実装時にユニットテスト・統合テストを作成
- **継続的インテグレーション**: CI/CDパイプラインを早期に構築

### 1.2 開発環境

- **言語**: Go 1.25
- **フレームワーク**: Gin
- **データベース**: MySQL 8.4.7
- **コンテナ**: Docker / Docker Compose
- **バージョン管理**: Git

### 1.3 アーキテクチャ

クリーンアーキテクチャに基づく3層構造:
- **Handler層**: HTTPリクエスト処理
- **Service層**: ビジネスロジック
- **Repository層**: データアクセス

---

## 2. フェーズ分け

### Phase 1: 基盤構築（Week 1-2）
**目標**: 開発環境とコア機能の基盤を構築

- データベース設計・構築
- プロジェクト構造の整備
- 基本的なCRUD機能
- 認証基盤

### Phase 2: コア勤怠機能（Week 3-4）
**目標**: 基本的な勤怠管理機能を実装

- 打刻機能
- 勤怠記録管理
- 従業員管理

### Phase 3: 休暇管理・承認フロー（Week 5-6）
**目標**: 休暇申請と承認ワークフローを実装

- 休暇申請機能
- 承認・却下機能
- 通知機能

### Phase 4: レポート・集計機能（Week 7-8）
**目標**: 勤怠データの集計とレポート機能を実装

- 月次レポート
- 部署別サマリー
- データエクスポート

### Phase 5: LINE連携（Week 9-11）
**目標**: LINE Botによる勤怠管理を実装

- LINE連携基盤
- LINE打刻機能
- LINE休暇申請
- LINEプッシュ通知

### Phase 6: 最適化・拡張（Week 12-）
**目標**: パフォーマンス最適化と追加機能

- パフォーマンスチューニング
- LIFF アプリ
- 追加機能

---

## 3. Phase 1: 基盤構築（Week 1-2）

### 3.1 環境セットアップ

#### Step 1.1: プロジェクト初期化
```bash
# Goモジュール初期化
go mod init github.com/nanato-okajima/attendance_management

# 必要なパッケージのインストール
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u go.uber.org/zap
go get -u github.com/go-playground/validator/v10
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt
```

**成果物**:
- `go.mod`
- `go.sum`

#### Step 1.2: ディレクトリ構造作成
```bash
mkdir -p {handler,service,database,logger,validator,myerrors,middleware,config,models}
mkdir -p build/database/{sql,data}
mkdir -p docs
```

**ディレクトリ構成**:
```
attendance_management/
├── main.go
├── route.go
├── config/           # 設定管理
├── handler/          # HTTPハンドラー
├── service/          # ビジネスロジック
│   └── repository/   # リポジトリインターフェース
├── database/         # データアクセス層
├── models/           # データモデル
├── middleware/       # ミドルウェア（認証等）
├── validator/        # カスタムバリデーター
├── logger/           # ロギング
├── myerrors/         # カスタムエラー
├── build/
│   ├── Dockerfile
│   └── database/
│       ├── sql/      # マイグレーションSQL
│       └── data/     # データ永続化
├── docs/             # ドキュメント
└── tests/            # テストコード
```

#### Step 1.3: Docker環境構築

**ファイル**: `docker-compose.yml`
```yaml
version: '3.8'
services:
  app:
    container_name: 'attendance_app'
    build:
      context: .
      dockerfile: ./build/Dockerfile
    env_file:
      - .env
    volumes:
      - .:/go/app/
    ports:
      - 8080:8080
    depends_on:
      - db

  db:
    container_name: 'attendance_db'
    image: mysql:8.4.7
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: admin
      MYSQL_DATABASE: attendance_management
    ports:
      - 3306:3306
    volumes:
      - ./build/database/data:/var/lib/mysql
      - ./build/database/sql:/docker-entrypoint-initdb.d
      - ./build/database/my.cnf:/etc/mysql/conf.d/my.cnf
```

**ファイル**: `build/Dockerfile`
```dockerfile
FROM golang:1.25-alpine

WORKDIR /go/app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]
```

**ファイル**: `.env`
```bash
DB_USER=root
DB_PASSWORD=admin
DB_HOST=db
DB_NAME=attendance_management

JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h
```

**成果物**:
- `docker-compose.yml`
- `build/Dockerfile`
- `.env`
- `.air.toml` (ホットリロード設定)

### 3.2 データベース設計・構築

#### Step 1.4: マイグレーションSQL作成

**ファイル**: `build/database/sql/01_create_tables.sql`
```sql
-- 従業員マスタ
CREATE TABLE employees (
    employee_number INT UNSIGNED AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    name_kana VARCHAR(255) NOT NULL,
    birthday DATE NOT NULL,
    gender_cd TINYINT UNSIGNED NOT NULL COMMENT '1:男性, 2:女性, 3:その他',
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    department_id INT UNSIGNED,
    position_id INT UNSIGNED,
    hire_date DATE NOT NULL,
    employment_type TINYINT UNSIGNED NOT NULL COMMENT '1:正社員, 2:契約, 3:パート, 4:アルバイト',
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (employee_number),
    INDEX idx_email (email),
    INDEX idx_is_deleted (is_deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ユーザーアカウント
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('admin', 'employee') NOT NULL DEFAULT 'employee',
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 出勤記録
CREATE TABLE attendances (
    attendance_id INT UNSIGNED AUTO_INCREMENT,
    employee_id INT UNSIGNED NOT NULL,
    target_date DATE NOT NULL,
    opening_time DATETIME,
    closing_time DATETIME,
    attendance_status TINYINT UNSIGNED NOT NULL COMMENT '1:出勤, 2:欠勤, 3:遅刻, 4:早退',
    clock_source TINYINT UNSIGNED DEFAULT 1 COMMENT '1:Web, 2:モバイル, 3:LINE',
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    work_hours DECIMAL(5,2),
    overtime_hours DECIMAL(5,2),
    note TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (attendance_id),
    FOREIGN KEY (employee_id) REFERENCES employees(employee_number) ON DELETE CASCADE,
    UNIQUE KEY uk_employee_date (employee_id, target_date),
    INDEX idx_target_date (target_date),
    INDEX idx_employee_date (employee_id, target_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 休暇申請
CREATE TABLE leave_requests (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    leave_type TINYINT UNSIGNED NOT NULL COMMENT '1:有給, 2:特別休暇, 3:欠勤, 4:半休',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    half_day_type TINYINT UNSIGNED COMMENT '1:午前半休, 2:午後半休',
    reason TEXT NOT NULL,
    approval_status TINYINT UNSIGNED DEFAULT 1 COMMENT '1:承認待ち, 2:承認済み, 3:却下',
    approver_id INT UNSIGNED,
    approved_at DATETIME,
    approval_comment TEXT,
    reject_reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_employee (employee_number),
    INDEX idx_status (approval_status),
    INDEX idx_dates (start_date, end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 打刻修正申請
CREATE TABLE attendance_corrections (
    id INT UNSIGNED AUTO_INCREMENT,
    attendance_id INT UNSIGNED NOT NULL,
    employee_number INT UNSIGNED NOT NULL,
    correction_type TINYINT UNSIGNED NOT NULL COMMENT '1:出勤修正, 2:退勤修正, 3:両方',
    corrected_opening_time DATETIME,
    corrected_closing_time DATETIME,
    reason TEXT NOT NULL,
    approval_status TINYINT UNSIGNED DEFAULT 1 COMMENT '1:承認待ち, 2:承認済み, 3:却下',
    approver_id INT UNSIGNED,
    approved_at DATETIME,
    reject_reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (attendance_id) REFERENCES attendances(attendance_id) ON DELETE CASCADE,
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_employee (employee_number),
    INDEX idx_status (approval_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 部署マスタ
CREATE TABLE departments (
    id INT UNSIGNED AUTO_INCREMENT,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    parent_id INT UNSIGNED,
    manager_employee_number INT UNSIGNED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY (manager_employee_number) REFERENCES employees(employee_number) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 役職マスタ
CREATE TABLE positions (
    id INT UNSIGNED AUTO_INCREMENT,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    level TINYINT UNSIGNED NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- LINE連携ユーザー
CREATE TABLE line_users (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    line_user_id VARCHAR(255) NOT NULL UNIQUE,
    linked_at DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    UNIQUE KEY uk_employee (employee_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- LINE連携コード
CREATE TABLE line_linking_codes (
    id INT UNSIGNED AUTO_INCREMENT,
    employee_number INT UNSIGNED NOT NULL,
    code VARCHAR(8) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    used_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_number) REFERENCES employees(employee_number) ON DELETE CASCADE,
    INDEX idx_expires (expires_at),
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**ファイル**: `build/database/sql/02_insert_master_data.sql`
```sql
-- 部署マスタ初期データ
INSERT INTO departments (code, name, parent_id, manager_employee_number) VALUES
('DEV', '開発部', NULL, NULL),
('SALES', '営業部', NULL, NULL),
('HR', '人事部', NULL, NULL);

-- 役職マスタ初期データ
INSERT INTO positions (code, name, level) VALUES
('MGR', 'マネージャー', 5),
('LEAD', 'リーダー', 4),
('SR', 'シニア', 3),
('MID', 'ミドル', 2),
('JR', 'ジュニア', 1);

-- テスト用従業員データ
INSERT INTO employees (name, name_kana, birthday, gender_cd, email, department_id, position_id, hire_date, employment_type) VALUES
('山田太郎', 'ヤマダタロウ', '1990-01-01', 1, 'yamada@example.com', 1, 3, '2020-04-01', 1),
('佐藤花子', 'サトウハナコ', '1992-05-15', 2, 'sato@example.com', 2, 2, '2021-04-01', 1);

-- テスト用ユーザーアカウント（パスワード: password123）
INSERT INTO users (employee_number, email, password_hash, role) VALUES
(1, 'yamada@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'),
(2, 'sato@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'employee');
```

**成果物**:
- `build/database/sql/01_create_tables.sql`
- `build/database/sql/02_insert_master_data.sql`

#### Step 1.5: データベース接続設定

**ファイル**: `config/config.go`
```go
package config

import (
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    Database DatabaseConfig
    JWT      JWTConfig
    LINE     LINEConfig
}

type DatabaseConfig struct {
    User     string `envconfig:"DB_USER" required:"true"`
    Password string `envconfig:"DB_PASSWORD" required:"true"`
    Host     string `envconfig:"DB_HOST" required:"true"`
    Name     string `envconfig:"DB_NAME" required:"true"`
}

type JWTConfig struct {
    Secret     string `envconfig:"JWT_SECRET" required:"true"`
    Expiration string `envconfig:"JWT_EXPIRATION" default:"24h"`
}

type LINEConfig struct {
    ChannelSecret string `envconfig:"LINE_CHANNEL_SECRET"`
    ChannelToken  string `envconfig:"LINE_CHANNEL_TOKEN"`
}

func Load() (*Config, error) {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

**ファイル**: `database/database.go`
```go
package database

import (
    "fmt"
    "github.com/nanato-okajima/attendance_management/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func SetupDB(cfg *config.DatabaseConfig) error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Name)
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return err
    }
    
    DB = db
    return nil
}

func GetDB() *gorm.DB {
    return DB
}
```

**成果物**:
- `config/config.go`
- `database/database.go`

### 3.3 ロギング・エラーハンドリング基盤

#### Step 1.6: ロガー設定

**ファイル**: `logger/logger.go`
```go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func SetupLogger() error {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    logger, err := config.Build()
    if err != nil {
        return err
    }
    
    Logger = logger.Sugar()
    return nil
}

func Info(args ...interface{}) {
    Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
    Logger.Infof(template, args...)
}

func Error(args ...interface{}) {
    Logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
    Logger.Errorf(template, args...)
}
```

**成果物**:
- `logger/logger.go`

#### Step 1.7: カスタムエラー定義

**ファイル**: `myerrors/errors.go`
```go
package myerrors

import "fmt"

type BadRequestError struct {
    Err error
}

func (e *BadRequestError) Error() string {
    return fmt.Sprintf("bad request: %v", e.Err)
}

type UnauthorizedError struct {
    Err error
}

func (e *UnauthorizedError) Error() string {
    return fmt.Sprintf("unauthorized: %v", e.Err)
}

type ForbiddenError struct {
    Err error
}

func (e *ForbiddenError) Error() string {
    return fmt.Sprintf("forbidden: %v", e.Err)
}

type NotFoundError struct {
    Err error
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("not found: %v", e.Err)
}

type InternalServerError struct {
    Err error
}

func (e *InternalServerError) Error() string {
    return fmt.Sprintf("internal server error: %v", e.Err)
}
```

**成果物**:
- `myerrors/errors.go`

### 3.4 認証基盤

#### Step 1.8: JWTユーティリティ

**ファイル**: `middleware/auth.go`
```go
package middleware

import (
    "errors"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/nanato-okajima/attendance_management/config"
    "github.com/nanato-okajima/attendance_management/myerrors"
)

type Claims struct {
    EmployeeNumber int    `json:"employee_number"`
    Email          string `json:"email"`
    Role           string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(employeeNumber int, email, role string, cfg *config.JWTConfig) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        EmployeeNumber: employeeNumber,
        Email:          email,
        Role:           role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.Secret))
}

func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}
        
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.Secret), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        c.Set("employee_number", claims.EmployeeNumber)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != "admin" {
            c.JSON(403, gin.H{"error": "admin access required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**成果物**:
- `middleware/auth.go`

#### Step 1.9: メインエントリーポイント

**ファイル**: `main.go`
```go
package main

import (
    "log"
    
    "github.com/nanato-okajima/attendance_management/config"
    "github.com/nanato-okajima/attendance_management/database"
    "github.com/nanato-okajima/attendance_management/logger"
)

func main() {
    // ロガー初期化
    if err := logger.SetupLogger(); err != nil {
        log.Fatal("Failed to setup logger:", err)
    }
    
    // 設定読み込み
    cfg, err := config.Load()
    if err != nil {
        logger.Error("Failed to load config:", err)
        log.Fatal(err)
    }
    
    // データベース接続
    if err := database.SetupDB(&cfg.Database); err != nil {
        logger.Error("Failed to setup database:", err)
        log.Fatal(err)
    }
    
    // ルーター設定
    router := setupRouter(cfg)
    
    // サーバー起動
    logger.Info("Server starting on :8080")
    if err := router.Run(":8080"); err != nil {
        logger.Error("Failed to start server:", err)
        log.Fatal(err)
    }
}
```

**ファイル**: `route.go`
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/nanato-okajima/attendance_management/config"
    "github.com/nanato-okajima/attendance_management/middleware"
)

func setupRouter(cfg *config.Config) *gin.Engine {
    r := gin.Default()
    
    // ヘルスチェック
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1
    v1 := r.Group("/v1")
    {
        // 認証不要エンドポイント
        // auth := v1.Group("/auth")
        // {
        //     auth.POST("/login", authHandler.Login)
        //     auth.POST("/register", authHandler.Register)
        // }
        
        // 認証必要エンドポイント
        // protected := v1.Group("")
        // protected.Use(middleware.AuthMiddleware(&cfg.JWT))
        // {
        //     // 従業員エンドポイント
        //     // 勤怠エンドポイント
        //     // 休暇エンドポイント
        // }
    }
    
    return r
}
```

**成果物**:
- `main.go`
- `route.go`

---

## 4. Phase 2: コア勤怠機能（Week 3-4）

### 4.1 認証機能実装

#### Step 2.1: ユーザーモデル定義

**ファイル**: `models/user.go`
```go
package models

import "time"

type User struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    EmployeeNumber int       `gorm:"not null" json:"employee_number"`
    Email          string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash   string    `gorm:"not null" json:"-"`
    Role           string    `gorm:"type:enum('admin','employee');default:'employee'" json:"role"`
    IsActive       bool      `gorm:"default:true" json:"is_active"`
    LastLoginAt    *time.Time `json:"last_login_at"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type Employee struct {
    EmployeeNumber   int       `gorm:"primaryKey;autoIncrement" json:"employee_number"`
    Name             string    `gorm:"not null" json:"name"`
    NameKana         string    `gorm:"not null" json:"name_kana"`
    Birthday         time.Time `gorm:"type:date;not null" json:"birthday"`
    GenderCd         int       `gorm:"not null" json:"gender_cd"`
    Email            string    `gorm:"uniqueIndex;not null" json:"email"`
    Phone            string    `json:"phone"`
    DepartmentID     *int      `json:"department_id"`
    PositionID       *int      `json:"position_id"`
    HireDate         time.Time `gorm:"type:date;not null" json:"hire_date"`
    EmploymentType   int       `gorm:"not null" json:"employment_type"`
    IsDeleted        bool      `gorm:"default:false" json:"is_deleted"`
    DeletedAt        *time.Time `json:"deleted_at"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

**成果物**:
- `models/user.go`

#### Step 2.2: 認証サービス実装

**ファイル**: `service/repository/user_repository.go`
```go
package repository

import "github.com/nanato-okajima/attendance_management/models"

type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    FindByEmployeeNumber(employeeNumber int) (*models.User, error)
    UpdateLastLogin(id uint) error
}
```

**ファイル**: `database/user_repository.go`
```go
package database

import (
    "github.com/nanato-okajima/attendance_management/models"
    "github.com/nanato-okajima/attendance_management/service/repository"
    "time"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
    return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
    return DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := DB.Where("email = ? AND is_active = ?", email, true).First(&user).Error
    return &user, err
}

func (r *userRepository) FindByEmployeeNumber(employeeNumber int) (*models.User, error) {
    var user models.User
    err := DB.Where("employee_number = ? AND is_active = ?", employeeNumber, true).First(&user).Error
    return &user, err
}

func (r *userRepository) UpdateLastLogin(id uint) error {
    now := time.Now()
    return DB.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}
```

**ファイル**: `service/auth_service.go`
```go
package service

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/nanato-okajima/attendance_management/config"
    "github.com/nanato-okajima/attendance_management/middleware"
    "github.com/nanato-okajima/attendance_management/models"
    "github.com/nanato-okajima/attendance_management/service/repository"
)

type AuthService interface {
    Login(email, password string) (string, error)
    Register(user *models.User, password string) error
}

type authService struct {
    userRepo repository.UserRepository
    cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
    return &authService{
        userRepo: userRepo,
        cfg:      cfg,
    }
}

func (s *authService) Login(email, password string) (string, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    
    token, err := middleware.GenerateToken(user.EmployeeNumber, user.Email, user.Role, &s.cfg.JWT)
    if err != nil {
        return "", err
    }
    
    s.userRepo.UpdateLastLogin(user.ID)
    
    return token, nil
}

func (s *authService) Register(user *models.User, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    user.PasswordHash = string(hashedPassword)
    return s.userRepo.Create(user)
}
```

**成果物**:
- `service/repository/user_repository.go`
- `database/user_repository.go`
- `service/auth_service.go`

#### Step 2.3: 認証ハンドラー実装

**ファイル**: `handler/auth_handler.go`
```go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/nanato-okajima/attendance_management/models"
    "github.com/nanato-okajima/attendance_management/service"
)

type AuthHandler interface {
    Login(c *gin.Context)
    Register(c *gin.Context)
}

type authHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
    return &authHandler{
        authService: authService,
    }
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    EmployeeNumber int    `json:"employee_number" binding:"required"`
    Email          string `json:"email" binding:"required,email"`
    Password       string `json:"password" binding:"required,min=8"`
    Role           string `json:"role" binding:"required,oneof=admin employee"`
}

func (h *authHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }
    
    c.JSON(200, gin.H{
        "token": token,
        "type":  "Bearer",
    })
}

func (h *authHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user := &models.User{
        EmployeeNumber: req.EmployeeNumber,
        Email:          req.Email,
        Role:           req.Role,
    }
    
    if err := h.authService.Register(user, req.Password); err != nil {
        c.JSON(500, gin.H{"error": "failed to register user"})
        return
    }
    
    c.JSON(201, gin.H{"message": "user registered successfully"})
}
```

**成果物**:
- `handler/auth_handler.go`

### 4.2 打刻機能実装

#### Step 2.4: 勤怠モデル定義

**ファイル**: `models/attendance.go`
```go
package models

import "time"

type Attendance struct {
    AttendanceID     uint       `gorm:"primaryKey;autoIncrement" json:"attendance_id"`
    EmployeeID       int        `gorm:"not null" json:"employee_id"`
    TargetDate       time.Time  `gorm:"type:date;not null" json:"target_date"`
    OpeningTime      *time.Time `json:"opening_time"`
    ClosingTime      *time.Time `json:"closing_time"`
    AttendanceStatus int        `gorm:"not null" json:"attendance_status"`
    ClockSource      int        `gorm:"default:1" json:"clock_source"`
    Latitude         *float64   `json:"latitude"`
    Longitude        *float64   `json:"longitude"`
    WorkHours        *float64   `json:"work_hours"`
    OvertimeHours    *float64   `json:"overtime_hours"`
    Note             string     `json:"note"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
}
```

**成果物**:
- `models/attendance.go`

#### Step 2.5: 打刻サービス実装

**ファイル**: `service/repository/attendance_repository.go`
```go
package repository

import (
    "github.com/nanato-okajima/attendance_management/models"
    "time"
)

type AttendanceRepository interface {
    Create(attendance *models.Attendance) error
    FindByEmployeeAndDate(employeeID int, date time.Time) (*models.Attendance, error)
    Update(attendance *models.Attendance) error
    FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]models.Attendance, error)
}
```

**ファイル**: `database/attendance_repository.go`
```go
package database

import (
    "github.com/nanato-okajima/attendance_management/models"
    "github.com/nanato-okajima/attendance_management/service/repository"
    "time"
)

type attendanceRepository struct{}

func NewAttendanceRepository() repository.AttendanceRepository {
    return &attendanceRepository{}
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
    return DB.Create(attendance).Error
}

func (r *attendanceRepository) FindByEmployeeAndDate(employeeID int, date time.Time) (*models.Attendance, error) {
    var attendance models.Attendance
    err := DB.Where("employee_id = ? AND target_date = ?", employeeID, date.Format("2006-01-02")).First(&attendance).Error
    return &attendance, err
}

func (r *attendanceRepository) Update(attendance *models.Attendance) error {
    return DB.Save(attendance).Error
}

func (r *attendanceRepository) FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]models.Attendance, error) {
    var attendances []models.Attendance
    err := DB.Where("employee_id = ? AND target_date BETWEEN ? AND ?", 
        employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
        Order("target_date ASC").
        Find(&attendances).Error
    return attendances, err
}
```

**ファイル**: `service/attendance_service.go`
```go
package service

import (
    "errors"
    "time"
    "github.com/nanato-okajima/attendance_management/models"
    "github.com/nanato-okajima/attendance_management/service/repository"
)

type AttendanceService interface {
    ClockIn(employeeID int, latitude, longitude *float64, clockSource int) error
    ClockOut(employeeID int, latitude, longitude *float64) error
    GetTodayAttendance(employeeID int) (*models.Attendance, error)
    GetMonthlyAttendances(employeeID int, year, month int) ([]models.Attendance, error)
}

type attendanceService struct {
    attendanceRepo repository.AttendanceRepository
}

func NewAttendanceService(attendanceRepo repository.AttendanceRepository) AttendanceService {
    return &attendanceService{
        attendanceRepo: attendanceRepo,
    }
}

func (s *attendanceService) ClockIn(employeeID int, latitude, longitude *float64, clockSource int) error {
    today := time.Now()
    
    // 既に打刻済みかチェック
    existing, _ := s.attendanceRepo.FindByEmployeeAndDate(employeeID, today)
    if existing != nil && existing.OpeningTime != nil {
        return errors.New("already clocked in today")
    }
    
    now := time.Now()
    attendance := &models.Attendance{
        EmployeeID:       employeeID,
        TargetDate:       today,
        OpeningTime:      &now,
        AttendanceStatus: 1, // 通常出勤
        ClockSource:      clockSource,
        Latitude:         latitude,
        Longitude:        longitude,
    }
    
    // TODO: 遅刻判定ロジック
    
    if existing != nil {
        attendance.AttendanceID = existing.AttendanceID
        return s.attendanceRepo.Update(attendance)
    }
    
    return s.attendanceRepo.Create(attendance)
}

func (s *attendanceService) ClockOut(employeeID int, latitude, longitude *float64) error {
    today := time.Now()
    
    attendance, err := s.attendanceRepo.FindByEmployeeAndDate(employeeID, today)
    if err != nil {
        return errors.New("no clock-in record found")
    }
    
    if attendance.ClosingTime != nil {
        return errors.New("already clocked out")
    }
    
    now := time.Now()
    attendance.ClosingTime = &now
    
    if latitude != nil {
        attendance.Latitude = latitude
    }
    if longitude != nil {
        attendance.Longitude = longitude
    }
    
    // TODO: 勤務時間計算
    // TODO: 早退判定
    
    return s.attendanceRepo.Update(attendance)
}

func (s *attendanceService) GetTodayAttendance(employeeID int) (*models.Attendance, error) {
    return s.attendanceRepo.FindByEmployeeAndDate(employeeID, time.Now())
}

func (s *attendanceService) GetMonthlyAttendances(employeeID int, year, month int) ([]models.Attendance, error) {
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
    endDate := startDate.AddDate(0, 1, -1)
    return s.attendanceRepo.FindByEmployeeAndDateRange(employeeID, startDate, endDate)
}
```

**成果物**:
- `service/repository/attendance_repository.go`
- `database/attendance_repository.go`
- `service/attendance_service.go`

#### Step 2.6: 打刻ハンドラー実装

**ファイル**: `handler/attendance_handler.go`
```go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/nanato-okajima/attendance_management/service"
)

type AttendanceHandler interface {
    ClockIn(c *gin.Context)
    ClockOut(c *gin.Context)
    GetToday(c *gin.Context)
    GetMonthly(c *gin.Context)
}

type attendanceHandler struct {
    attendanceService service.AttendanceService
}

func NewAttendanceHandler(attendanceService service.AttendanceService) AttendanceHandler {
    return &attendanceHandler{
        attendanceService: attendanceService,
    }
}

type ClockRequest struct {
    Latitude  *float64 `json:"latitude"`
    Longitude *float64 `json:"longitude"`
}

func (h *attendanceHandler) ClockIn(c *gin.Context) {
    employeeNumber := c.GetInt("employee_number")
    
    var req ClockRequest
    c.ShouldBindJSON(&req)
    
    if err := h.attendanceService.ClockIn(employeeNumber, req.Latitude, req.Longitude, 1); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, gin.H{"message": "clocked in successfully"})
}

func (h *attendanceHandler) ClockOut(c *gin.Context) {
    employeeNumber := c.GetInt("employee_number")
    
    var req ClockRequest
    c.ShouldBindJSON(&req)
    
    if err := h.attendanceService.ClockOut(employeeNumber, req.Latitude, req.Longitude); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "clocked out successfully"})
}

func (h *attendanceHandler) GetToday(c *gin.Context) {
    employeeNumber := c.GetInt("employee_number")
    
    attendance, err := h.attendanceService.GetTodayAttendance(employeeNumber)
    if err != nil {
        c.JSON(404, gin.H{"error": "no attendance record found"})
        return
    }
    
    c.JSON(200, attendance)
}

func (h *attendanceHandler) GetMonthly(c *gin.Context) {
    employeeNumber := c.GetInt("employee_number")
    
    var query struct {
        Year  int `form:"year" binding:"required"`
        Month int `form:"month" binding:"required,min=1,max=12"`
    }
    
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    attendances, err := h.attendanceService.GetMonthlyAttendances(employeeNumber, query.Year, query.Month)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to fetch attendances"})
        return
    }
    
    c.JSON(200, attendances)
}
```

**成果物**:
- `handler/attendance_handler.go`

### 4.3 ルーティング統合

#### Step 2.7: ルーター更新

**ファイル**: `route.go` (更新)
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/nanato-okajima/attendance_management/config"
    "github.com/nanato-okajima/attendance_management/database"
    "github.com/nanato-okajima/attendance_management/handler"
    "github.com/nanato-okajima/attendance_management/middleware"
    "github.com/nanato-okajima/attendance_management/service"
)

func setupRouter(cfg *config.Config) *gin.Engine {
    r := gin.Default()
    
    // リポジトリ初期化
    userRepo := database.NewUserRepository()
    attendanceRepo := database.NewAttendanceRepository()
    
    // サービス初期化
    authService := service.NewAuthService(userRepo, cfg)
    attendanceService := service.NewAttendanceService(attendanceRepo)
    
    // ハンドラー初期化
    authHandler := handler.NewAuthHandler(authService)
    attendanceHandler := handler.NewAttendanceHandler(attendanceService)
    
    // ヘルスチェック
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1
    v1 := r.Group("/v1")
    {
        // 認証エンドポイント
        auth := v1.Group("/auth")
        {
            auth.POST("/login", authHandler.Login)
            auth.POST("/register", authHandler.Register)
        }
        
        // 認証必要エンドポイント
        protected := v1.Group("")
        protected.Use(middleware.AuthMiddleware(&cfg.JWT))
        {
            // 打刻エンドポイント
            attendance := protected.Group("/attendance")
            {
                attendance.POST("/clock-in", attendanceHandler.ClockIn)
                attendance.POST("/clock-out", attendanceHandler.ClockOut)
                attendance.GET("/today", attendanceHandler.GetToday)
                attendance.GET("/monthly", attendanceHandler.GetMonthly)
            }
        }
    }
    
    return r
}
```

**成果物**:
- `route.go` (更新)

---

## 5. Phase 3: 休暇管理・承認フロー（Week 5-6）

### 5.1 休暇申請機能

#### Step 3.1: 休暇モデル・リポジトリ実装

**実装内容**:
- `models/leave_request.go`: 休暇申請モデル
- `service/repository/leave_repository.go`: リポジトリインターフェース
- `database/leave_repository.go`: リポジトリ実装
- `service/leave_service.go`: 休暇申請ビジネスロジック
- `handler/leave_handler.go`: 休暇申請ハンドラー

**主要機能**:
- 休暇申請作成
- 休暇申請一覧取得
- 有給残日数計算
- 申請期間の重複チェック

### 5.2 承認フロー実装

#### Step 3.2: 承認機能実装

**実装内容**:
- `service/approval_service.go`: 承認ビジネスロジック
- `handler/approval_handler.go`: 承認ハンドラー

**主要機能**:
- 休暇申請承認
- 休暇申請却下
- 打刻修正申請承認
- 承認待ち一覧取得

### 5.3 通知機能基盤

#### Step 3.3: 通知サービス実装

**実装内容**:
- `service/notification_service.go`: 通知ビジネスロジック
- バッチ処理用スクリプト

**主要機能**:
- メール通知（将来実装）
- 承認待ち通知
- 承認結果通知

---

## 6. Phase 4: レポート・集計機能（Week 7-8）

### 6.1 月次レポート

#### Step 4.1: レポートサービス実装

**実装内容**:
- `service/report_service.go`: レポート生成ロジック
- `handler/report_handler.go`: レポートハンドラー

**主要機能**:
- 月次勤怠集計
- 出勤日数・欠勤日数計算
- 総勤務時間・残業時間計算
- 遅刻・早退回数集計

### 6.2 データエクスポート

#### Step 4.2: エクスポート機能実装

**実装内容**:
- CSV生成機能
- Excel生成機能（オプション）

---

## 7. Phase 5: LINE連携（Week 9-11）

### 7.1 LINE Bot基盤構築

#### Step 5.1: LINE SDK導入

```bash
go get -u github.com/line/line-bot-sdk-go/v8/linebot
```

**実装内容**:
- `config/config.go`: LINE設定追加
- `handler/line_webhook_handler.go`: Webhookハンドラー
- `service/line_service.go`: LINE連携サービス

### 7.2 アカウント連携

#### Step 5.2: 連携コード機能実装

**実装内容**:
- `models/line_user.go`: LINEユーザーモデル
- `models/line_linking_code.go`: 連携コードモデル
- `service/line_linking_service.go`: 連携サービス
- `handler/line_linking_handler.go`: 連携ハンドラー

**主要機能**:
- 連携コード発行
- LINE User ID紐付け
- 連携解除

### 7.3 LINE打刻機能

#### Step 5.3: LINE打刻実装

**実装内容**:
- Webhookイベント処理
- テキストメッセージハンドリング
- 位置情報メッセージハンドリング
- クイックリプライ実装

**主要機能**:
- 「出勤」コマンド処理
- 「退勤」コマンド処理
- 位置情報検証
- 打刻結果通知

### 7.4 LINE休暇申請

#### Step 5.4: 対話型申請フロー実装

**実装内容**:
- ステートマシン管理
- カレンダー表示（Flex Message）
- 確認画面表示

### 7.5 LINEプッシュ通知

#### Step 5.5: プッシュ通知実装

**実装内容**:
- 打刻忘れ通知バッチ
- 承認待ち通知
- 承認結果通知
- 残業アラート通知

### 7.6 リッチメニュー

#### Step 5.6: リッチメニュー設定

**実装内容**:
- 従業員用リッチメニュー画像作成
- 管理者用リッチメニュー画像作成
- リッチメニュー設定API実装

---

## 8. Phase 6: 最適化・拡張（Week 12-）

### 8.1 テスト実装

#### Step 6.1: ユニットテスト

**実装内容**:
- サービス層のテスト
- リポジトリ層のテスト（モック使用）
- ハンドラー層のテスト

#### Step 6.2: 統合テスト

**実装内容**:
- API統合テスト
- データベーステスト

### 8.2 パフォーマンス最適化

#### Step 6.3: 最適化実装

**実装内容**:
- データベースインデックス最適化
- クエリ最適化
- キャッシュ導入（Redis）

### 8.3 CI/CD構築

#### Step 6.4: パイプライン構築

**実装内容**:
- GitHub Actions設定
- 自動テスト実行
- 自動デプロイ

---

## 9. 開発ガイドライン

### 9.1 コーディング規約

- Goの標準的な命名規則に従う
- `gofmt`でフォーマット
- `golint`でリント
- エラーは明示的にハンドリング
- ログは構造化ロギング（Zap）を使用

### 9.2 Git運用

- ブランチ戦略: Git Flow
- コミットメッセージ: Conventional Commits
- プルリクエスト必須
- コードレビュー実施

### 9.3 テスト方針

- ユニットテストカバレッジ80%以上
- 統合テスト必須
- E2Eテスト（重要フロー）

---

## 10. リスク管理

### 10.1 技術的リスク

| リスク | 影響度 | 対策 |
|--------|--------|------|
| LINE API仕様変更 | 中 | SDK最新版を使用、変更通知を監視 |
| データベース性能問題 | 高 | 早期に性能テスト実施 |
| セキュリティ脆弱性 | 高 | 定期的なセキュリティ監査 |

### 10.2 スケジュールリスク

| リスク | 影響度 | 対策 |
|--------|--------|------|
| 要件変更 | 中 | アジャイル開発、柔軟な対応 |
| 技術的難易度 | 中 | 早期にPoCを実施 |
| リソース不足 | 高 | 優先度を明確化、MVP重視 |

---

## 11. 成果物チェックリスト

### Phase 1
- [ ] Docker環境構築完了
- [ ] データベーステーブル作成完了
- [ ] 基本的なプロジェクト構造完成
- [ ] ロガー・エラーハンドリング実装完了
- [ ] JWT認証基盤完成

### Phase 2
- [ ] ログイン・ユーザー登録機能完成
- [ ] 出勤・退勤打刻機能完成
- [ ] 勤怠記録取得機能完成
- [ ] 基本的なバリデーション実装完了

### Phase 3
- [ ] 休暇申請機能完成
- [ ] 承認・却下機能完成
- [ ] 通知機能基盤完成

### Phase 4
- [ ] 月次レポート機能完成
- [ ] データエクスポート機能完成

### Phase 5
- [ ] LINE Bot基盤完成
- [ ] LINE連携機能完成
- [ ] LINE打刻機能完成
- [ ] LINE休暇申請完成
- [ ] LINEプッシュ通知完成

### Phase 6
- [ ] ユニットテスト実装完了
- [ ] 統合テスト実装完了
- [ ] CI/CD構築完了

---

**文書バージョン**: 1.0  
**作成日**: 2025-11-25  
**最終更新日**: 2025-11-25
