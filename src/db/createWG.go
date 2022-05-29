package db

import (
	"fmt"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

func createWG(PrivateKey string, PublicKey string, ListenPort int, IPv4Address string, IPv6Address string) {
	combinedLogger := logger.GetCombinedLogger()

	combinedLogger.Info("Creating wireguard interface")
	db := DBSystem
	var wgCreate WireguardInterface

	if IPv6Address != "-" {
		wgCreate = WireguardInterface{ //with IPv6
			InterfaceName: "wg0",
			PrivateKey:    PrivateKey,
			PublicKey:     PublicKey,
			ListenPort:    ListenPort,
			IPv4Address:   IPv4Address,
			IPv6Address:   IPv6Address,
		}
	} else {
		wgCreate = WireguardInterface{ //without IPv6
			InterfaceName: "wg0",
			PrivateKey:    PrivateKey,
			PublicKey:     PublicKey,
			ListenPort:    ListenPort,
			IPv4Address:   IPv4Address,
		}
	}

	err := db.Create(wgCreate).Error
	if err != nil {
		combinedLogger.Error(fmt.Sprintf("Creating interface %s", err))
	}
	combinedLogger.Info("Successfully created interface")
}
