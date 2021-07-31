package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func LoggerSetup() {
	currentTime := time.Now()                                   //get current time
	timeStr := currentTime.Format("02-01-2006 15:04:05")        //set time format
	errCreateDir := os.MkdirAll("/opt/wgManagerAPI/logs", 0755) //create log directory if does not exist

	if errCreateDir != nil { //if an error on creating directory
		fmt.Println("Error on creating directory for logger \n", errCreateDir)
		os.Exit(1) //exit program
	}

	file, errLog := os.OpenFile("/opt/wgManagerAPI"+timeStr+" log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755) //open file
	if errLog != nil {
		fmt.Println("Error when opening log file \n", errLog)
	} else {
		log.SetOutput(file)
	}
}
