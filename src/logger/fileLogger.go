package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fileLogger *zap.Logger

func GetFileLogger() *zap.Logger {
	if fileLogger == nil {
		config := zap.NewProductionEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncodeCaller = nil
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewJSONEncoder(config)
		logFile, _ := os.OpenFile(time.Now().Format("02-01-2006 15:04:05")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(logFile), zapcore.DebugLevel)
		logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel))
		defer logger.Sync()

		fileLogger = logger
	}

	return fileLogger
}
