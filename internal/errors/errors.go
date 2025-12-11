package errors

import "fmt"

// AppError はアプリケーション全体で使用する統一エラー型
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap は元のエラーを返す（errors.Is/As対応）
func (e *AppError) Unwrap() error {
	return e.Err
}

// New は新しいAppErrorを作成
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap は既存のエラーをAppErrorでラップ
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewDBError はデータベースエラーを作成
func NewDBError(code ErrorCode, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: "database error",
		Err:     err,
	}
}

// NewValidationError はバリデーションエラーを作成
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}
