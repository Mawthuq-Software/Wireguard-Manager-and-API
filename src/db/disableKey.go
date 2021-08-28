package db

import (
	"errors"
	"strconv"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"gorm.io/gorm"
)

func DisableKey(keyID string) (bool, map[string]string) {
	var keyStruct Key
	responseMap := make(map[string]string)
	db := DBSystem

	keyIDInt, _ := strconv.Atoi(keyID)                              //convert key string to int
	resultKey := db.Where("key_id = ?", keyIDInt).First(&keyStruct) //find key in database
	if errors.Is(resultKey.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}

	keyStruct.Enabled = "false"      //set IP back to unused
	keyUpdate := db.Save(&keyStruct) //save data
	if keyUpdate.Error != nil {
		responseMap["response"] = "Error in updating client key"
		return false, responseMap
	}
	pubKey := keyStruct.PublicKey
	boolRes, stringRes := manager.DeleteKey("wg0", pubKey) //delete key from wg interface
	if boolRes {
		responseMap["response"] = "Disabled key successfully"
	} else {
		responseMap["response"] = stringRes
	}
	return boolRes, responseMap
}
