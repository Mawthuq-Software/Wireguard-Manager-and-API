package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

func main() {
	fmt.Println("WG MANAGER AND API STARTING UP")

	fmt.Println("Env file loading - 1/")
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		fmt.Println("Env failed to load - FAILED")
		os.Exit(1)
	}

	fmt.Println("Logger starting up - 2/")
	logger.LoggerSetup()

}
