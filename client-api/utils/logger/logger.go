package logger

import (
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func NewLogger() {
	config := zap.NewProductionConfig()
	config.Level = getLevel()
	buildConfig, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = buildConfig
}

func getLevel() zap.AtomicLevel {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "DEBUG":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "WARN":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "ERROR":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "INFO", "":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}
