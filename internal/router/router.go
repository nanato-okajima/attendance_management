package router

import (
	"github.com/gin-gonic/gin"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/handler"
	"github.com/nanato-okajima/attendance_management/internal/middleware"
	"github.com/nanato-okajima/attendance_management/internal/repository/database"
	"github.com/nanato-okajima/attendance_management/internal/service"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.New() // gin.Default()の代わりにgin.New()を使用

	// グローバルミドルウェアを適用
	r.Use(middleware.Recovery())      // パニックリカバリー
	r.Use(middleware.RequestLogger()) // リクエストログ
	r.Use(middleware.CORS())          // CORS設定

	// リポジトリ初期化
	db := database.GetDB()
	userRepo := database.NewUserRepository()
	attendanceRepo := database.NewAttendanceRepository()
	leaveRepo := database.NewLeaveRepository(db)
	correctionRepo := database.NewCorrectionRepository(db)

	// サービス初期化
	authService := service.NewAuthService(userRepo, cfg)
	attendanceService := service.NewAttendanceService(attendanceRepo, &cfg.WorkHours)
	notificationService := service.NewNotificationService()
	leaveService := service.NewLeaveService(leaveRepo)
	correctionService := service.NewAttendanceCorrectionService(correctionRepo, attendanceRepo, notificationService)
	approvalService := service.NewApprovalService(leaveRepo, notificationService)

	// ハンドラー初期化
	authHandler := handler.NewAuthHandler(authService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	leaveHandler := handler.NewLeaveHandler(leaveService)
	correctionHandler := handler.NewAttendanceCorrectionHandler(correctionService)
	approvalHandler := handler.NewApprovalHandler(approvalService, correctionService)

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

		// 認証必要エンドポイント
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(&cfg.JWT))
		{
			// 勤怠エンドポイント
			attendances := protected.Group("/attendances")
			{
				attendances.POST("/clock-in", attendanceHandler.ClockIn)
				attendances.POST("/clock-out", attendanceHandler.ClockOut)
				attendances.GET("/today", attendanceHandler.GetToday)
				attendances.GET("/monthly", attendanceHandler.GetMonthly)
			}

			// 休暇申請エンドポイント
			leave := protected.Group("/leave")
			{
				leave.POST("/requests", leaveHandler.CreateLeaveRequest)
				leave.GET("/requests", leaveHandler.GetLeaveRequests)
				leave.GET("/remaining-days", leaveHandler.GetRemainingPaidLeaveDays)
			}

			// 打刻修正申請エンドポイント
			corrections := protected.Group("/corrections")
			{
				corrections.POST("/requests", correctionHandler.CreateCorrectionRequest)
				corrections.GET("/requests", correctionHandler.GetCorrectionRequests)
				corrections.GET("/requests/:id", correctionHandler.GetCorrectionRequest)
			}

			// 承認エンドポイント（管理者のみ）
			approvals := protected.Group("/approvals")
			approvals.Use(middleware.AdminOnly())
			{
				approvals.GET("/pending", approvalHandler.GetPendingApprovals)
				approvals.POST("/leave/:id/approve", approvalHandler.ApproveLeaveRequest)
				approvals.POST("/leave/:id/reject", approvalHandler.RejectLeaveRequest)
				approvals.POST("/corrections/:id/approve", approvalHandler.ApproveCorrectionRequest)
				approvals.POST("/corrections/:id/reject", approvalHandler.RejectCorrectionRequest)
			}
		}
	}

	return r
}
