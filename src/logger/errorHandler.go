package logger

import (
	"fmt"
)

func ErrorHandler(message string, err error) bool { //error handler
	combinedLogger := GetCombinedLogger()

	if err != nil {
		combinedLogger.Error(fmt.Sprintf(message+" %s", err))
		return false
	}
	return true
}
