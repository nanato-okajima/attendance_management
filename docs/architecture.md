# アーキテクチャドキュメント

最終更新: 2025-12-05

## 概要

本システムは **Clean Architecture** と **Domain-Driven Design (DDD)** の原則に基づいて設計された勤怠管理システムです。

## アーキテクチャ原則

### 1. レイヤー分離

```
┌─────────────────────────────────────────┐
│          Handler Layer                  │  ← HTTPリクエスト/レスポンス
│  - リクエストバリデーション                 │
│  - DTOマッピング                          │
│  - エラーハンドリング                      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Service Layer                   │  ← ビジネスロジック
│  - ユースケース実装                      │
│  - トランザクション制御                  │
│  - ドメインロジック呼び出し              │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Domain Layer                     │  ← ドメインロジック
│  - エンティティ                          │
│  - ドメインサービス                      │
│  - リポジトリインターフェース            │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Repository Layer                   │  ← データ永続化
│  - GORM実装                              │
│  - DB接続管理                            │
└─────────────────────────────────────────┘
```

### 2. 依存性逆転の原則 (DIP)

- Service層はDomain層のインターフェースに依存
- Repository層はDomain層のインターフェースを実装
- 上位レイヤーは下位レイヤーの実装に依存しない

### 3. 集約 (Aggregate) による整理

Domain層は集約ごとにディレクトリを分割:

```
domain/
├── attendance/          # 勤怠集約
│   ├── attendance.go    # ドメインロジック
│   ├── repository.go    # Repository IF
│   └── work_hours_calculator.go
└── user/                # ユーザー集約
    └── repository.go
```

## ディレクトリ構成

現在のディレクトリ構成 (2025-12-05時点):

```
internal/
├── config/              # 設定管理
│   ├── config.go        # メイン設定
│   ├── database.go      # DB設定
│   ├── jwt.go           # JWT設定
│   ├── server.go        # サーバー設定
│   └── validator.go     # 設定バリデーション
├── domain/              # ドメイン層（集約ごと）
│   ├── attendance/
│   │   ├── attendance.go
│   │   ├── repository.go
│   │   └── work_hours_calculator.go
│   └── user/
│       └── repository.go
├── entity/              # エンティティ定義
│   ├── attendance.go
│   ├── user.go
│   ├── employee.go
│   ├── leave_request.go
│   ├── department.go
│   ├── position.go
│   ├── line_user.go
│   └── line_linking_code.go
├── errors/              # エラーハンドリング
│   ├── errors.go        # AppError型
│   ├── codes.go         # エラーコード定義
│   └── handler.go       # HTTPエラーハンドラー
├── handler/             # ハンドラー層
│   ├── attendance.go
│   └── auth_handler.go
├── middleware/          # ミドルウェア
│   ├── auth.go          # JWT認証
│   ├── cors.go          # CORS設定
│   ├── logger.go        # リクエストログ
│   └── recovery.go      # パニックリカバリー
├── repository/          # リポジトリ実装
│   └── database/
│       ├── attendance.go
│       ├── connection.go
│       └── user.go
├── router/              # ルーティング
│   └── router.go
├── service/             # サービス層
│   ├── attendance.go
│   ├── auth_service.go
│   └── time_provider.go
└── mock/                # モック（自動生成）
    ├── mock_attendance_repository.go
    └── mock_user_repository.go
```

## エラーハンドリング

### AppError型

すべてのエラーは`AppError`型で統一:

```go
type AppError struct {
    Code    ErrorCode  // エラーコード
    Message string     // エラーメッセージ
    Err     error      // 元のエラー
}
```

### エラーコード体系

| カテゴリ | プレフィックス | 例 |
|---------|---------------|-----|
| 認証 | AUTH_ | AUTH_001, AUTH_002 |
| 勤怠 | ATT_ | ATT_001, ATT_002 |
| DB | DB_ | DB_001, DB_002 |
| バリデーション | VAL_ | VAL_001 |
| 内部エラー | INT_ | INT_001 |

### HTTPステータスコードマッピング

- `AUTH_*` → 401 Unauthorized
- `ATT_*` → 400 Bad Request
- `DB_001` (Not Found) → 404 Not Found
- `DB_002` (Duplicate) → 409 Conflict
- その他 → 500 Internal Server Error

## ミドルウェア

### 1. Recovery

パニックをキャッチしてサーバークラッシュを防ぐ。スタックトレースをログに記録。

### 2. RequestLogger

すべてのHTTPリクエストをログに記録。ステータスコードに応じてログレベルを変更:
- 500番台 → Error
- 400番台 → Warn
- 200番台 → Info

### 3. CORS

Cross-Origin Resource Sharingを設定。フロントエンドからのアクセスを許可。

### 4. Auth

JWT認証。トークンの検証とクレーム抽出。

## テスト戦略

### テストカバレッジ

- Domain層: 100%
- Service層: 83.3%
- Handler層: 30.2%
- Middleware層: 55.2%

### Mock生成

`go:generate`ディレクティブを使用して自動生成:

```go
//go:generate mockgen -source=repository.go -destination=../../mock/mock_attendance_repository.go -package=mock -mock_names Repository=MockAttendanceRepository
```

### テストパッケージ構成

Service層などのテストでは、循環参照を回避し、ブラックボックステストを推奨するために `package service_test` のように別パッケージとしてテストを記述します。

```go
package service_test

import (
    "testing"
    "github.com/nanato-okajima/attendance_management/internal/service"
    "github.com/nanato-okajima/attendance_management/internal/mock"
)
```

## 今後の拡張

### Phase 3C: Graceful Shutdown

処理中のリクエストを完了してから安全にシャットダウン。

### 将来的な機能拡張

- Rate Limiting
- Request ID追跡
- メトリクス収集 (Prometheus)
- 分散トレーシング
- キャッシング (Redis)
