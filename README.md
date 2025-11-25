# Attendance Management System

勤怠管理システムのバックエンドAPI

## 技術スタック

- **言語**: Go 1.24
- **フレームワーク**: Gin
- **データベース**: MySQL 8.4
- **ORM**: GORM
- **認証**: JWT
- **開発ツール**: Air (ホットリロード), golangci-lint

## セットアップ

### 前提条件

- Docker Desktop がインストールされていること
- Task がインストールされていること ([Installation](https://taskfile.dev/installation/))

### 環境構築

1. リポジトリをクローン

```bash
git clone https://github.com/nanato-okajima/attendance_management.git
cd attendance_management
```

2. 環境変数を設定

`.env` ファイルを編集し、必要な環境変数を設定します：

```bash
# DATABASE
DB_USER=root
DB_PASSWORD=admin
DB_HOST=db
DB_NAME=attendance_management

# JWT
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION=24h

# LINE (Optional)
LINE_CHANNEL_SECRET=
LINE_CHANNEL_TOKEN=
```

> **⚠️ 重要**: 本番環境では `JWT_SECRET` を強力なランダム文字列に変更してください。

3. Docker コンテナを起動

```bash
task up
```

4. アプリケーションの動作確認

```bash
curl localhost:8080/health
# {"status":"ok"}
```

## 開発ワークフロー

### 利用可能な Task コマンド

```bash
# コンテナを起動
task up

# コンテナを停止
task down

# コンテナを再起動
task reboot

# ログを表示
task logs

# コンテナの状態を確認
task ps

# Lint を実行 (ローカル環境)
task lint

# Lint を実行 (Docker コンテナ内)
task lint-docker
```

### コード品質チェック

このプロジェクトでは [golangci-lint](https://github.com/golangci/golangci-lint) を使用してコード品質を管理しています。

#### Lint の実行方法

**Docker コンテナ内で実行（推奨）:**

```bash
task lint-docker
```

**ローカル環境で実行:**

golangci-lint がインストールされている場合：

```bash
task lint
```

または直接実行：

```bash
golangci-lint run ./...
```

#### 有効化されているリンター

- `gofmt` - コードフォーマットのチェック
- `goimports` - import文の整理
- `govet` - Go の静的解析
- `errcheck` - エラーハンドリングのチェック
- `staticcheck` - 静的解析
- `gosec` - セキュリティ問題の検出
- その他多数（詳細は `.golangci.yml` を参照）

#### Lint 設定のカスタマイズ

`.golangci.yml` ファイルを編集することで、リンターの設定をカスタマイズできます。

### ホットリロード

開発中は Air によるホットリロードが有効になっています。ファイルを保存すると自動的にアプリケーションが再起動されます。

## プロジェクト構造

```
.
├── build/              # Docker関連ファイル
│   ├── Dockerfile
│   └── database/       # データベース設定・初期化SQL
├── config/             # 設定管理
├── database/           # データベース接続
├── docs/               # ドキュメント
├── handler/            # HTTPハンドラー
├── logger/             # ロガー設定
├── middleware/         # ミドルウェア
├── models/             # データモデル
├── service/            # ビジネスロジック
│   └── repository/     # データアクセス層
├── validator/          # バリデーション
├── main.go             # エントリーポイント
├── route.go            # ルーティング設定
├── .golangci.yml       # golangci-lint設定
├── compose.yml         # Docker Compose設定
└── Taskfile.yml        # Taskコマンド定義
```

## API エンドポイント

### ヘルスチェック

```bash
GET /health
```

レスポンス:
```json
{"status":"ok"}
```

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。
