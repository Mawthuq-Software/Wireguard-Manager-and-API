package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fileLogger *zap.Logger
var fileLoggerCore zapcore.Core

func GetFileLogger() *zap.Logger {
	if fileLogger == nil {
		core := zapcore.NewTee(GetFileLoggerCore())
		logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		defer logger.Sync()

		fileLogger = logger
	}

	return fileLogger
}

func GetFileLoggerCore() zapcore.Core {
	if fileLoggerCore == nil {
		config := zap.NewProductionEncoderConfig()
		config.EncodeLevel = zapcore.CapitalLevelEncoder
		config.EncodeCaller = nil
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoder := zapcore.NewConsoleEncoder(config)
		logFile, _ := os.OpenFile("/opt/wgManagerAPI/logs/"+time.Now().Format("02-01-2006 15:04:05")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		core := zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.DebugLevel)

		fileLoggerCore = core
	}

	return fileLoggerCore
}
