package db

import (
	"errors"

	"gorm.io/gorm"
)

type ReturnKeysRes struct {
	Response string
	Keys     []Key
}

func ReturnKeys() (bool, ReturnKeysRes) {
	var allKeysStruct []Key
	var responseMap ReturnKeysRes
	db := DBSystem

	resultKeys := db.Find(&allKeysStruct)
	if errors.Is(resultKeys.Error, gorm.ErrRecordNotFound) {
		responseMap.Response = "Keys were not found on the server"
		return false, responseMap
	}

	for i := 0; i < len(allKeysStruct); i++ {
		allKeysStruct[i].PresharedKey = "REDACTED"
	}

	responseMap.Response = "All key successfully parsed"
	responseMap.Keys = allKeysStruct
	return true, responseMap
}
