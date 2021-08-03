package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/api"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/autocheck"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/network"
)

func main() {
	fmt.Println("WG MANAGER AND API STARTING UP")

	fmt.Println("Env file loading - 1/6")
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		fmt.Println("Env failed to load - FAILED")
		os.Exit(1)
	}

	fmt.Println("Logger starting up - 2/6")
	logger.LoggerSetup()

	fmt.Println("Starting database - 3/6")
	db.DBStart()

	fmt.Println("Starting of network - 4/6")
	network.SetupWG()

	if os.Getenv("autocheck") == "enabled" {
		fmt.Println("Starting autochecker - 5/6")
		autocheck.AutoStart()
	} else {
		fmt.Println("Skipped autochecker - 5/6")
	}

	fmt.Println("Starting API - 6/6")
	api.API()
}
