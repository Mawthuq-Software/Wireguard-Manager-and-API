package db

import (
	"strconv"
	"time"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

func SetSubscription(keyID string, bwLimit int64, subExpiry string, bwReset bool) (bool, map[string]string) {
	var subStructModify Subscription
	db := DBSystem
	combinedLogger := logger.GetCombinedLogger()

	responseMap := make(map[string]string)
	keyIDInt, _ := strconv.Atoi(keyID) //convert to int
	resultSub := db.Where("key_id = ?", keyIDInt).First(&subStructModify)
	if resultSub.Error != nil {
		combinedLogger.Error("Error Finding subscription in db " + resultSub.Error.Error())
		responseMap["response"] = "Error - Finding subscription"
		return false, responseMap
	}
	if bwLimit >= 0 {
		subStructModify.BandwidthAllotted = bwLimit
	}
	if subExpiry != "-1" {
		_, subErr := time.Parse("2006-Jan-02 03:04:05 PM", subExpiry)
		if !logger.ErrorHandler("Error - Parsing stored time ", subErr) {
			responseMap["response"] = "Error - Parsing time"
			return false, responseMap
		} else {
			subStructModify.SubscriptionEnd = subExpiry
		}
	}
	if bwReset {
		subStructModify.BandwidthUsed = 0
	}
	db.Where("key_id = ?", keyIDInt).Save(&subStructModify)
	responseMap["response"] = "Updated successfully"
	return true, responseMap
}
