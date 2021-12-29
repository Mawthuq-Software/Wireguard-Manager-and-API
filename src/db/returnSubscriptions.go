package db

import (
	"errors"

	"gorm.io/gorm"
)

type ReturnSubsRes struct {
	Response      string
	Subscriptions []Subscription
}

func ReturnSubscriptions() (bool, ReturnSubsRes) {
	var allSubsStruct []Subscription
	var responseMap ReturnSubsRes
	db := DBSystem

	resultKeys := db.Find(&allSubsStruct)
	if errors.Is(resultKeys.Error, gorm.ErrRecordNotFound) {
		responseMap.Response = "Subscriptions were not found on the server"
		return false, responseMap
	}

	responseMap.Response = "All subscriptions successfully parsed"
	responseMap.Subscriptions = allSubsStruct
	return true, responseMap
}
