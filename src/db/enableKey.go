package db

import (
	"errors"
	"strconv"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"gorm.io/gorm"
)

func EnableKey(keyID string) (bool, map[string]string) {
	var keyStruct Key
	var ipStruct IP
	responseMap := make(map[string]string)
	db := DBSystem

	keyIDInt, _ := strconv.Atoi(keyID)                              //convert key string to int
	resultKey := db.Where("key_id = ?", keyIDInt).First(&keyStruct) //find key in database
	if errors.Is(resultKey.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}

	keyStruct.Enabled = "true"       //set IP back to unused
	keyUpdate := db.Save(&keyStruct) //save data
	if keyUpdate.Error != nil {
		responseMap["response"] = "Error in updating client key"
		return false, responseMap
	}

	ipv4Addr := keyStruct.IPv4Address
	pubKey := keyStruct.PublicKey
	preKey := keyStruct.PresharedKey

	resultIP := db.Where("ipv4_address = ?", ipv4Addr).First(&ipStruct) //find IP in db
	if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}
	ipv6Addr := ipStruct.IPv6Address
	boolRes, stringRes := manager.AddKey("wg0", ipv4Addr, ipv6Addr, pubKey, preKey) //add key to wg interface
	if boolRes == true {
		responseMap["response"] = "Enabled key successfully"
	} else {
		responseMap["response"] = stringRes
	}
	return boolRes, responseMap
}
