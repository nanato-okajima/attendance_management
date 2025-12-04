package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nanato-okajima/attendance_management/internal/logger"
)

// RequestLogger はHTTPリクエストの詳細をログに記録するミドルウェア
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// ステータスコードに応じてログレベルを変更
		if statusCode >= 500 {
			logger.Error("Request failed",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency,
				"ip", clientIP,
			)
		} else if statusCode >= 400 {
			logger.Warn("Client error",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency,
				"ip", clientIP,
			)
		} else {
			logger.Info("Request",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency,
				"ip", clientIP,
			)
		}
	}
}
