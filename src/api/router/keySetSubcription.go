package router

import (
	"log"
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keySetSubJSON struct {
	KeyID     string `json:"keyID"`
	BWLimit   int64  `json:"bwLimit"`
	SubExpiry string `json:"subExpiry"`
	BWReset   bool   `json:"bwReset"`
}

func keySetSubscription(res http.ResponseWriter, req *http.Request) {
	var incomingJson keySetSubJSON

	err := parseResponse(req, &incomingJson) //parse JSON
	if err != nil {
		log.Println("Error - Parsing request", err)
		sentStandardRes(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if incomingJson.KeyID == "" {
		sentStandardRes(res, map[string]string{"response": "Bad Request, keyID must be filled"}, http.StatusBadRequest)
		return
	}

	boolRes, mapRes := db.SetSubscription(incomingJson.KeyID, incomingJson.BWLimit, incomingJson.SubExpiry, incomingJson.BWReset) //add key to db
	if !boolRes {
		sentStandardRes(res, mapRes, http.StatusBadRequest)
	} else {
		sentStandardRes(res, mapRes, http.StatusAccepted)
	}
}
