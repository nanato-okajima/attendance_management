# Attendance Management System

勤怠管理システム - Clean Architecture / DDD を採用した Go 製 REST API

## 📋 目次

- [概要](#概要)
- [アーキテクチャ](#アーキテクチャ)
- [ディレクトリ構成](#ディレクトリ構成)
- [技術スタック](#技術スタック)
- [セットアップ](#セットアップ)
- [開発](#開発)
- [エラーハンドリング](#エラーハンドリング)
- [API仕様](#api仕様)

---

## 概要

従業員の勤怠管理を行うREST APIシステムです。出退勤の記録、勤務時間の計算、月次レポートの生成などの機能を提供します。

### 主な機能

- ✅ 出勤・退勤の打刻
- ✅ 勤務時間の自動計算
- ✅ 遅刻・早退の自動判定
- ✅ 月次勤怠レポート
- ✅ JWT認証
- ✅ 位置情報記録（オプション）

---

## アーキテクチャ

### Clean Architecture / DDD

本プロジェクトは **Clean Architecture** と **Domain-Driven Design (DDD)** の原則に基づいて設計されています。

```
┌─────────────────────────────────────────┐
│          Handler Layer (API)            │
│  - HTTP リクエスト/レスポンス処理        │
│  - バリデーション                        │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Service Layer (UseCase)         │
│  - ビジネスロジック                      │
│  - トランザクション制御                  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Domain Layer (Entity)            │
│  - ドメインロジック                      │
│  - リポジトリインターフェース            │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Repository Layer (Database)        │
│  - データ永続化                          │
│  - GORM による DB アクセス               │
└─────────────────────────────────────────┘
```

### レイヤー分離の原則

- **依存性逆転の原則 (DIP)**: 上位レイヤーは下位レイヤーに依存しない
- **単一責任の原則 (SRP)**: 各レイヤーは明確な責任を持つ
- **テスタビリティ**: モックを使用した単体テストが容易

---

## ディレクトリ構成

```
attendance_management/
├── main.go                    # エントリーポイント
├── Taskfile.yml              # タスクランナー設定
├── go.mod                    # Go モジュール定義
├── docs/                     # ドキュメント
├── build/                    # Docker 設定
│   └── Dockerfile
└── internal/
    ├── config/               # 設定管理
    │   ├── config.go
    │   ├── database.go
    │   ├── jwt.go
    │   ├── server.go
    │   └── validator.go
    ├── domain/               # ドメイン層（集約ごと）
    │   ├── attendance/
    │   │   ├── attendance.go           # ドメインロジック
    │   │   ├── repository.go           # リポジトリIF
    │   │   └── work_hours_calculator.go
    │   └── user/
    │       └── repository.go
    ├── entity/               # エンティティ定義
    │   ├── attendance.go
    │   ├── user.go
    │   └── ...
    ├── errors/               # エラーハンドリング
    │   ├── errors.go         # AppError型
    │   ├── codes.go          # エラーコード定義
    │   └── handler.go        # HTTPエラーハンドラー
    ├── handler/              # ハンドラー層
    │   ├── attendance.go
    │   └── auth_handler.go
    ├── middleware/           # ミドルウェア
    │   ├── auth.go           # JWT認証
    │   ├── cors.go           # CORS設定
    │   ├── logger.go         # リクエストログ
    │   └── recovery.go       # パニックリカバリー
    ├── repository/           # リポジトリ実装
    │   └── database/
    │       ├── attendance.go
    │       ├── connection.go
    │       └── user.go
    ├── router/               # ルーティング
    │   └── router.go
    ├── service/              # サービス層
    │   ├── attendance.go
    │   └── auth_service.go
    └── mock/                 # モック（自動生成）
        ├── mock_attendance_repository.go
        └── mock_user_repository.go
```

---

## 技術スタック

### バックエンド

- **言語**: Go 1.24
- **フレームワーク**: Gin (HTTP)
- **ORM**: GORM
- **データベース**: MySQL 8.0
- **認証**: JWT (golang-jwt/jwt)

### 開発ツール

- **タスクランナー**: Task
- **コンテナ**: Docker / Docker Compose
- **リンター**: golangci-lint
- **モック生成**: mockgen (uber-go/mock)
- **ホットリロード**: Air

---

## セットアップ

### 前提条件

- Docker & Docker Compose
- Task (タスクランナー)

### 環境変数

`.env`ファイルを作成してください:

```env
# Database
DB_USER=root
DB_PASSWORD=password
DB_HOST=db
DB_NAME=attendance_db

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Work Hours
WORK_START_HOUR=9
WORK_START_MINUTE=0
WORK_END_HOUR=18
WORK_END_MINUTE=0
BREAK_HOURS=1
STANDARD_WORK_HOURS=8
```

### 起動

```bash
# コンテナ起動
task up

# マイグレーション実行
task migrate

# 動作確認
curl http://localhost:8080/health
```

---

## 開発

### よく使うコマンド

```bash
# 開発サーバー起動（ホットリロード）
task dev

# テスト実行
task test

# テストカバレッジ確認
task test-coverage

# Mock生成
task generate-mocks

# Lint実行
task lint

# Lint自動修正
task lint-fix

# コンテナ停止
task down
```

### テスト

```bash
# 全テスト実行
task test

# 特定パッケージのテスト
go test ./internal/service/... -v

# カバレッジレポート生成
task test-coverage
# coverage.html が生成されます
```

### Mock生成

リポジトリインターフェースのモックは自動生成されます:

```bash
task generate-mocks
```

モックは`internal/mock/`に生成されます。

---

## エラーハンドリング

### エラーコード体系

すべてのエラーには一意のエラーコードが付与されます:

| カテゴリ | コード | 説明 |
|---------|--------|------|
| **認証** | AUTH_001 | 認証失敗 |
| | AUTH_002 | 無効なトークン |
| | AUTH_003 | トークン期限切れ |
| | AUTH_004 | 認証情報不正 |
| **勤怠** | ATT_001 | 既に出勤済み |
| | ATT_002 | 出勤記録なし |
| | ATT_003 | 既に退勤済み |
| **DB** | DB_001 | レコード未検出 |
| | DB_002 | 重複エラー |

### エラーレスポンス形式

```json
{
  "error": "Already clocked in today",
  "code": "ATT_001"
}
```

---

## API仕様

### 認証

#### ログイン

```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**レスポンス:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "type": "Bearer",
  "user": {
    "id": 1,
    "employee_number": 1001,
    "email": "user@example.com",
    "role": "employee"
  }
}
```

#### ユーザー登録

```http
POST /v1/auth/register
Content-Type: application/json

{
  "employee_number": 1001,
  "email": "user@example.com",
  "password": "password123",
  "role": "employee"
}
```

### 勤怠管理

**認証が必要です。** `Authorization: Bearer <token>` ヘッダーを付与してください。

#### 出勤打刻

```http
POST /v1/attendances/clock-in
Authorization: Bearer <token>
Content-Type: application/json

{
  "latitude": 35.6812,
  "longitude": 139.7671
}
```

**レスポンス:**

```json
{
  "message": "clocked in successfully",
  "attendance": {
    "attendance_id": 1,
    "employee_id": 1001,
    "target_date": "2025-12-05",
    "opening_time": "2025-12-05T09:00:00Z",
    "attendance_status": 1
  }
}
```

#### 退勤打刻

```http
POST /v1/attendances/clock-out
Authorization: Bearer <token>
Content-Type: application/json

{
  "latitude": 35.6812,
  "longitude": 139.7671
}
```

#### 当日の勤怠取得

```http
GET /v1/attendances/today
Authorization: Bearer <token>
```

#### 月次勤怠取得

```http
GET /v1/attendances/monthly?year=2025&month=12
Authorization: Bearer <token>
```

---

## ライセンス

MIT License

---

## 開発者向け情報

### アーキテクチャの詳細

詳細なアーキテクチャドキュメントは `docs/` ディレクトリを参照してください:

- `docs/requirements_specification.md` - 要件定義書
- `docs/implementation_plan.md` - 実装計画書

### コントリビューション

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### コーディング規約

- Go の標準的なコーディング規約に従う
- `golangci-lint` でリントエラーがないこと
- テストカバレッジ 80% 以上を維持
- すべてのpublic関数にコメントを記載
