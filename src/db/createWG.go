package db

import (
	"log"
)

func createWG(PrivateKey string, PublicKey string, ListenPort int, IPv4Address string, IPv6Address string) {
	log.Println("Info - Creating wireguard interface")
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
		log.Fatal("Creating interface", err)
	}
	log.Println("Successfully created interface")
}
