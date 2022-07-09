package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func SetupLogger() {
	logger, _ := zap.NewDevelopment()
	Logger = logger.Sugar()
}
