package db

import (
	"errors"
	"log"
	"strconv"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"gorm.io/gorm"
)

func DeleteKey(keyID string) (bool, map[string]string) {
	var ipStruct IP
	var keyStruct Key
	var subStruct Subscription

	responseMap := make(map[string]string)
	db := DBSystem

	keyIDInt, _ := strconv.Atoi(keyID)                              //convert key string to int
	resultKey := db.Where("key_id = ?", keyIDInt).First(&keyStruct) //find key in database
	if errors.Is(resultKey.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}

	pubKey := keyStruct.PublicKey                              //set pub key
	ipv4 := keyStruct.IPv4Address                              //set ipv4 address
	delKey := db.Where("key_id = ?", keyID).Delete(&keyStruct) //delete key from db
	if delKey.Error != nil {
		log.Println("Finding key in DB", delKey.Error)
		responseMap["response"] = "Error occurred when finding the key in database"
		return false, responseMap
	}

	resultIP := db.Where("ipv4_address = ?", ipv4).First(&ipStruct) //find IP in db
	if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}

	ipStruct.InUse = "false"       //set IP back to unused
	ipUpdate := db.Save(&ipStruct) //save data
	if ipUpdate.Error != nil {
		responseMap["response"] = "Error in updating IP"
		return false, responseMap
	}

	delSub := db.Where("key_id = ?", keyID).Delete(&subStruct) //delete subcription from db
	if delSub.Error != nil {
		log.Println("Finding key in DB", delSub.Error)
		responseMap["response"] = "Error occurred when finding the subscription in database"
		return false, responseMap
	}
	boolRes, stringRes := manager.DeleteKey("wg0", pubKey) //delete key from wg interface
	responseMap["response"] = stringRes
	return boolRes, responseMap
}
