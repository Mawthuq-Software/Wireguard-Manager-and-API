package router

import (
	"log"
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keyCreateJSON struct {
	PublicKey    string `json:"publicKey"`
	PresharedKey string `json:"presharedKey"`
	BWLimit      int64  `json:"bwLimit"`
	SubExpiry    string `json:"subExpiry"`
}

func keyCreate(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyCreateJSON

	err := parseResponse(req, &incomingJson) //parse JSON
	if err != nil {
		log.Println("Error - Parsing request", err)
		sentStandardRes(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if incomingJson.PresharedKey == "" || incomingJson.PublicKey == "" {
		sentStandardRes(res, map[string]string{"response": "Bad Request, presharedKey and publicKey must be filled"}, http.StatusBadRequest)
		return
	} else if incomingJson.BWLimit < 0 {
		sentStandardRes(res, map[string]string{"response": "Bad Request, bandwidth cannot be negative"}, http.StatusBadRequest)
		return
	} else if incomingJson.SubExpiry == "" {
		sentStandardRes(res, map[string]string{"response": "Bad Request, subscription expiry must be filled"}, http.StatusBadRequest)
		return
	}

	boolRes, mapRes := db.CreateKey(incomingJson.PublicKey, incomingJson.PresharedKey, incomingJson.BWLimit, incomingJson.SubExpiry) //add key to db
	if !boolRes {
		sentStandardRes(res, mapRes, http.StatusBadRequest)
	} else {
		sentStandardRes(res, mapRes, http.StatusAccepted)
	}
}
