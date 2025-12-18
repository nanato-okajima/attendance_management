# リファクタリング記録: Service層のテストパッケージ変更

## 背景

`internal/service` パッケージのテストにおいて、`internal/mock` パッケージを使用する際に循環参照エラー (`import cycle not allowed in test`) が発生しました。

### 循環参照の原因

1. `internal/service` パッケージ: ビジネスロジックとインターフェース定義
2. `internal/mock` パッケージ: `service` パッケージのインターフェースをモック化（`service` パッケージをインポート）
3. `internal/service` 内のテスト (`package service`): テストコードが `mock` パッケージをインポート

依存関係: `service` (test) -> `mock` -> `service` -> 循環

## 対応内容

Service層のテストファイルを `package service` から **`package service_test`** に変更しました。

### 変更点

- テストファイルのパッケージ宣言を `package service_test` に変更
- テスト対象の `service` パッケージをインポート (`github.com/nanato-okajima/attendance_management/internal/service`)
- テスト内で使用する `service` パッケージの型や関数に `service.` プレフィックスを付与

### メリット

1. **循環参照の解消**: テストパッケージ (`service_test`) は `service` パッケージとは別のパッケージとして扱われるため、`mock` パッケージを経由しても循環しません。
   - `service_test` -> `mock` -> `service`
   - `service_test` -> `service`
2. **ブラックボックステストの強制**: 公開API（Exportされた関数・型）のみを使用したテストとなり、内部実装の詳細に依存しない堅牢なテストになります。
3. **Go標準の推奨構成**: Goの標準ライブラリや多くのOSSで採用されている一般的な構成です。

## 影響範囲

以下のテストファイルが変更されました:
- `internal/service/approval_service_test.go`
- `internal/service/attendance_correction_service_test.go`
- `internal/service/attendance_service_test.go`
- `internal/service/auth_service_test.go`
- `internal/service/leave_service_test.go`
