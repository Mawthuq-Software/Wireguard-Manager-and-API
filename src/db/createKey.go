package db

import (
	"errors"
	"log"
	"os"
	"strconv"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"gorm.io/gorm"
)

func CreateKey(pubKey string, preKey string) (bool, map[string]string) {
	var ipStruct IP
	var wgStruct WireguardInterface
	responseMap := make(map[string]string)
	db := DBSystem

	resultIP := db.Where("in_use = ?", "false").First(&ipStruct) //find IP not in use
	if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "No available IPs on the server"
		return false, responseMap
	}

	keyStructCreate := Key{PublicKey: pubKey, PresharedKey: preKey, IPv4Address: ipStruct.IPv4Address, Enabled: "true"} //create Key object
	resultKeyCreate := db.Create(&keyStructCreate)                                                                      //add object to db
	if resultKeyCreate.Error != nil {
		log.Println("Error - Adding key to db", resultKeyCreate.Error)
		responseMap["response"] = "Error when adding key to database"
		return false, responseMap
	}
	ipStruct.InUse = "true"                                                                                             //set ip to in use
	db.Save(&ipStruct)                                                                                                  //update IP in db
	keyIDStr := strconv.Itoa(keyStructCreate.KeyID)                                                                     //convert keyID to string
	boolRes, strRes := manager.AddKey(ipStruct.WGInterface, ipStruct.IPv4Address, ipStruct.IPv6Address, pubKey, preKey) //add key to wg interface
	if !boolRes {                                                                                                       //if an error occurred
		responseMap["response"] = strRes
		return boolRes, responseMap
	} else {
		responseMap["response"] = "Added key successfully"
		responseMap["ipv4Address"] = ipStruct.IPv4Address + "/32"
		if ipStruct.IPv6Address != "-" {
			responseMap["ipv6Address"] = ipStruct.IPv6Address + "/128"
		}
		responseMap["ipAddress"] = os.Getenv("IP_ADDRESS")
		responseMap["dns"] = os.Getenv("DNS")
		responseMap["allowedIPs"] = os.Getenv("ALLOWED_IP")
		responseMap["keyID"] = keyIDStr
	}

	resultWG := db.Where("interface_name = ?", "wg0").First(&wgStruct) //get wireguard server info

	if resultWG.Error != nil {
		responseMap["response"] = "Issue in finding a key for the server"
		return false, responseMap
	} else {
		responseMap["publicKey"] = wgStruct.PublicKey                 //return back wg server pub key
		responseMap["listenPort"] = strconv.Itoa(wgStruct.ListenPort) //return back wg server listenPort
		return boolRes, responseMap
	}
}
