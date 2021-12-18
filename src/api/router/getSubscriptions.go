package router

import (
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

func getSubscriptions(res http.ResponseWriter, req *http.Request) {
	boolRes, mapRes := db.ReturnSubscriptions() //add key to db
	if !boolRes {
		sendSubsRes(res, mapRes, http.StatusBadRequest)
	} else {
		sendSubsRes(res, mapRes, http.StatusAccepted)
	}
}
