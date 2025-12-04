package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError はAppErrorを適切なHTTPレスポンスに変換
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		statusCode := getStatusCode(appErr.Code)
		c.JSON(statusCode, gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		})
		return
	}

	// AppErrorでない場合は500エラー
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
		"code":  ErrCodeInternal,
	})
}

// getStatusCode はエラーコードに対応するHTTPステータスコードを返す
func getStatusCode(code ErrorCode) int {
	switch code {
	// 認証エラー -> 401
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeExpiredToken, ErrCodeInvalidCredentials:
		return http.StatusUnauthorized

	// Not Found -> 404
	case ErrCodeNotFound:
		return http.StatusNotFound

	// バリデーションエラー -> 400
	case ErrCodeValidation, ErrCodeInvalidRequest, ErrCodeBadRequest:
		return http.StatusBadRequest

	// 勤怠関連エラー -> 400
	case ErrCodeAlreadyClockedIn, ErrCodeNotClockedIn, ErrCodeAlreadyClockedOut, ErrCodeInvalidDate:
		return http.StatusBadRequest

	// 重複エラー -> 409
	case ErrCodeDuplicate:
		return http.StatusConflict

	// その他 -> 500
	default:
		return http.StatusInternalServerError
	}
}
