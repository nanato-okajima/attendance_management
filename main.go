package main

import (
	"log"

	"github.com/nanato-okajima/attendance_management/config"
	"github.com/nanato-okajima/attendance_management/database"
	"github.com/nanato-okajima/attendance_management/logger"
)

func main() {
	// ロガー初期化
	if err := logger.SetupLogger(); err != nil {
		log.Fatal("Failed to setup logger:", err)
	}

	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config:", err)
		log.Fatal(err)
	}

	// データベース接続
	if err := database.SetupDB(&cfg.Database); err != nil {
		logger.Error("Failed to setup database:", err)
		log.Fatal(err)
	}

	// ルーター設定
	router := setupRouter(cfg)

	// サーバー起動
	logger.Info("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("Failed to start server:", err)
		log.Fatal(err)
	}
}
