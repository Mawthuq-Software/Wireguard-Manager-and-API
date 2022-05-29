package db

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"gorm.io/gorm"
)

func CreateKey(pubKey string, preKey string, bwLimit int64, subEnd string, ipIndex int) (bool, map[string]string) {
	var ipStruct IP
	var wgStruct WireguardInterface
	responseMap := make(map[string]string)
	db := DBSystem
	combinedLogger := logger.GetCombinedLogger()

	//Must be first to prevent uneccessary key additions to DB
	ipMap := viper.GetStringSlice("INSTANCE.IP.GLOBAL.ADDRESS.IPV4")
	if len(ipMap)-1 < ipIndex {
		responseMap["response"] = "IP index does not exist"
		return false, responseMap
	}
	ipSelected := ipMap[ipIndex]

	_, subErr := time.Parse("2006-Jan-02 03:04:05 PM", subEnd)
	if !logger.ErrorHandler("Error - Parsing stored time ", subErr) {
		responseMap["response"] = "Error when parsing time"
		return false, responseMap
	}

	resultIP := db.Where("in_use = ?", "false").First(&ipStruct) //find IP not in use
	if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "No available IPs on the server"
		return false, responseMap
	}

	keyStructCreate := Key{PublicKey: pubKey, PresharedKey: preKey, IPv4Address: ipStruct.IPv4Address, Enabled: "true"} //create Key object
	resultKeyCreate := db.Create(&keyStructCreate)                                                                      //add object to db
	if resultKeyCreate.Error != nil {
		combinedLogger.Error(fmt.Sprintf("Adding key to db %s", resultKeyCreate.Error))
		responseMap["response"] = "Error when adding key to database"
		return false, responseMap
	}
	ipStruct.InUse = "true"                         //set ip to in use
	db.Save(&ipStruct)                              //update IP in db
	keyIDStr := strconv.Itoa(keyStructCreate.KeyID) //convert keyID to string

	subStructCreate := Subscription{KeyID: keyStructCreate.KeyID, PublicKey: pubKey, BandwidthUsed: 0, BandwidthAllotted: bwLimit, SubscriptionEnd: subEnd}
	resultSub := db.Create(&subStructCreate)
	if resultSub.Error != nil {
		combinedLogger.Error(fmt.Sprintf("Adding subscription to db %s", resultKeyCreate.Error))
		responseMap["response"] = "Error when adding subscription to database"
		return false, responseMap
	}

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

		responseMap["ipAddress"] = ipSelected
		responseMap["dns"] = viper.GetString("INSTANCE.IP.GLOBAL.DNS")
		responseMap["allowedIPs"] = viper.GetString("INSTANCE.IP.GLOBAL.ALLOWED")
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
