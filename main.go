package main

import (
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/api"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/autocheck"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/config"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/network"
)

func main() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Firing up Wireguard Manager & API")

	config.LoadConfig()
	db.DBStart()
	network.SetupWG()
	autocheck.AutoStart()
	api.API()
}
