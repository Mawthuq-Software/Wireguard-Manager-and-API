package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var combinedLogger *zap.Logger

func GetCombinedLogger() *zap.Logger {
	if combinedLogger == nil {
		fileLoggerCore := GetFileLoggerCore()
		consoleLoggerCore := GetConsoleLoggerCore()

		core := zapcore.NewTee(
			fileLoggerCore,
			consoleLoggerCore,
		)

		combinedLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return combinedLogger
}
