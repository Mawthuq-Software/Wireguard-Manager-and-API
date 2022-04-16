package db

import (
	"log"
	"strconv"
)

func GetUserSubscription(keyID string) (bool, map[string]string) {
	var subStructModify Subscription
	db := DBSystem
	responseMap := make(map[string]string)
	keyIDInt, _ := strconv.Atoi(keyID) //convert to int
	resultSub := db.Where("key_id = ?", keyIDInt).First(&subStructModify)
	if resultSub.Error != nil {
		log.Println("Error - Finding subscription in db", resultSub.Error)
		responseMap["response"] = "Error - Finding subscription"
		return false, responseMap
	}

	responseMap["response"] = "Queried successfully"
	responseMap["bandwidthUsed"] = strconv.Itoa(int(subStructModify.BandwidthUsed))
	responseMap["bandwidthLimit"] = strconv.Itoa(int(subStructModify.BandwidthAllotted))
	responseMap["subscriptionEnd"] = subStructModify.SubscriptionEnd
	return true, responseMap
}
