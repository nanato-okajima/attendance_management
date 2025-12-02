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
	attendanceService := service.NewAttendanceService(attendanceRepo, &cfg.WorkHours)

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
		// 認証エンドポイント（認証不要）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// 勤怠エンドポイント（認証必要）
		attendances := v1.Group("/attendances")
		attendances.Use(middleware.AuthMiddleware(&cfg.JWT))
		{
			attendances.POST("/clock-in", attendanceHandler.ClockIn)
			attendances.POST("/clock-out", attendanceHandler.ClockOut)
			attendances.GET("/today", attendanceHandler.GetToday)
			attendances.GET("/monthly", attendanceHandler.GetMonthly)
		}
	}

	return r
}
