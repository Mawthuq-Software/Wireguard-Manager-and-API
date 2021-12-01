package router

import (
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

func getKeys(res http.ResponseWriter, req *http.Request) {
	boolRes, mapRes := db.ReturnKeys() //add key to db
	if !boolRes {
		sendKeysRes(res, mapRes, http.StatusBadRequest)
	} else {
		sendKeysRes(res, mapRes, http.StatusAccepted)
	}
}
