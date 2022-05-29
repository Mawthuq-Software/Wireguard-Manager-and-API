package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var consoleLogger *zap.Logger
var consoleLoggerCore zapcore.Core

func GetConsoleLogger() *zap.Logger {
	if consoleLogger == nil {
		core := zapcore.NewTee(GetConsoleLoggerCore())
		logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		defer logger.Sync()

		consoleLogger = logger
	}

	return consoleLogger
}

func GetConsoleLoggerCore() zapcore.Core {
	if consoleLoggerCore == nil {
		config := zap.NewProductionEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncodeCaller = nil
		config.EncodeTime = nil
		consoleEncoder := zapcore.NewConsoleEncoder(config)
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

		consoleLoggerCore = core
	}

	return consoleLoggerCore
}
