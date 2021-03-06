package router

import (
	"log"
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keyEnableJSON struct {
	KeyID string `json:"keyID"`
}

func keyEnable(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyEnableJSON

	jsonResponse := make(map[string]string)

	err := parseResponse(req, &incomingJson)
	if err != nil {
		log.Println(err)
		sentStandardRes(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if incomingJson.KeyID == "" {
		jsonResponse["response"] = "Bad Request, keyID needs to be filled"
		sentStandardRes(res, map[string]string{"response": "Bad Request, keyID needs to be filled"}, http.StatusBadRequest)
		return
	}

	keyID := incomingJson.KeyID

	boolRes, mapRes := db.EnableKey(keyID)
	if !boolRes {
		sentStandardRes(res, mapRes, http.StatusBadRequest)
	} else {
		sentStandardRes(res, mapRes, http.StatusAccepted)
	}
}
