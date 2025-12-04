package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/logger"
	"github.com/nanato-okajima/attendance_management/internal/repository/database"
	"github.com/nanato-okajima/attendance_management/internal/router"
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
	r := router.Setup(cfg)

	// HTTPサーバー設定
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// サーバーを別ゴルーチンで起動
	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error:", err)
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// シャットダウンタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	// 処理中のリクエストを完了してからシャットダウン
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown:", err)
		log.Fatal(err)
	}

	logger.Info("Server exited")
}
