package logger

import "log"

func ErrorHandler(message string, err error) bool { //error handler
	if err != nil {
		log.Println(message, err)
		return false
	}
	return true
}
