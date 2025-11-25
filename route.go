package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nanato-okajima/attendance_management/config"
)

func setupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	// TODO: 将来的にエンドポイントを追加
	// v1 := r.Group("/v1")
	// {
	// 	// 認証不要エンドポイント
	// 	// auth := v1.Group("/auth")
	// 	// {
	// 	//     auth.POST("/login", authHandler.Login)
	// 	//     auth.POST("/register", authHandler.Register)
	// 	// }

	// 	// 認証必要エンドポイント
	// 	// protected := v1.Group("")
	// 	// protected.Use(middleware.AuthMiddleware(&cfg.JWT))
	// 	// {
	// 	//     // 従業員エンドポイント
	// 	//     // 勤怠エンドポイント
	// 	//     // 休暇エンドポイント
	// 	// }
	// }

	return r
}
