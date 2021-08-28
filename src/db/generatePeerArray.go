package db

import (
	"errors"
	"log"
	"net"
	"time"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"
)

func generatePeerArray() []wgtypes.PeerConfig {
	var keyStruct []Key               //key struct
	var keyArray []wgtypes.PeerConfig //peers (clients)
	db := DBSystem

	resultKey := db.Find(&keyStruct)
	if errors.Is(resultKey.Error, gorm.ErrRecordNotFound) {
		return keyArray
	} else if resultKey.Error != nil {
		log.Println("Error - Finding keys", resultKey.Error)
	}

	for i := 0; i < len(keyStruct); i++ { //loop over all clients in db
		var ipStruct IP
		resultIP := db.Where("ipv4_address = ?", keyStruct[i].IPv4Address).First(&ipStruct)
		if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
			log.Println("Cant find IPs", keyStruct[i].IPv4Address)
			continue //continue even on error
		} else if resultIP.Error != nil {
			log.Println("Error - Finding IPs", keyStruct[i].IPv4Address, resultKey.Error)
		} else if keyStruct[i].Enabled == "true" { //checks if key is enabled
			pubKey, pubErr := manager.ParseKey(keyStruct[i].PublicKey)
			preKey, preErr := manager.ParseKey(keyStruct[i].PresharedKey)
			if pubErr != nil || preErr != nil {
				log.Fatal("Error - Unable to parse keys on generate array")
			}

			var ipAddresses []net.IPNet
			ipv4, errIPv4 := manager.ParseIP(ipStruct.IPv4Address + "/32")
			if errIPv4 != nil {
				log.Fatal("Error - Parsing IPv4 Address", errIPv4)
			}
			ipAddresses = append(ipAddresses, *ipv4)

			if ipStruct.IPv6Address != "-" {
				ipv6, errIPv6 := manager.ParseIP(ipStruct.IPv6Address + "/128")
				if errIPv6 != nil {
					log.Fatal("Error - Parsing IPv6 Address", errIPv6)
				}
				ipAddresses = append(ipAddresses, *ipv6)
			}

			var zeroTime time.Duration
			userConfig := wgtypes.PeerConfig{
				PublicKey:                   pubKey,
				PresharedKey:                &preKey,
				PersistentKeepaliveInterval: &zeroTime,
				AllowedIPs:                  ipAddresses,
			}
			keyArray = append(keyArray, userConfig) //add config to client array
		}
	}
	return keyArray
}
