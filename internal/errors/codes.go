package errors

// ErrorCode はエラーを一意に識別するコード
type ErrorCode string

const (
	// 認証関連エラー
	ErrCodeUnauthorized       ErrorCode = "AUTH_001"
	ErrCodeInvalidToken       ErrorCode = "AUTH_002"
	ErrCodeExpiredToken       ErrorCode = "AUTH_003"
	ErrCodeInvalidCredentials ErrorCode = "AUTH_004"

	// 勤怠関連エラー
	ErrCodeAlreadyClockedIn  ErrorCode = "ATT_001"
	ErrCodeNotClockedIn      ErrorCode = "ATT_002"
	ErrCodeAlreadyClockedOut ErrorCode = "ATT_003"
	ErrCodeInvalidDate       ErrorCode = "ATT_004"

	// バリデーションエラー
	ErrCodeValidation     ErrorCode = "VAL_001"
	ErrCodeInvalidRequest ErrorCode = "VAL_002"

	// データベースエラー
	ErrCodeNotFound      ErrorCode = "DB_001"
	ErrCodeDuplicate     ErrorCode = "DB_002"
	ErrCodeDatabaseError ErrorCode = "DB_003"

	// 一般エラー
	ErrCodeInternal   ErrorCode = "INT_001"
	ErrCodeBadRequest ErrorCode = "INT_002"
)
