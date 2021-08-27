package router

import (
	"log"
	"net/http"
	"os"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keyDisableJSON struct {
	KeyID string `json:"keyID"`
}

func keyDisable(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyDisableJSON

	jsonResponse := make(map[string]string)

	err := parseResponse(req, &incomingJson)
	if err != nil {
		log.Println(err)
		sendResponse(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if os.Getenv("AUTH") != "-" { //check AUTH
		authHeader := req.Header.Get("Authorization")
		if os.Getenv("AUTH") != authHeader {
			sendResponse(res, map[string]string{"response": "Authentication key is not valid"}, http.StatusBadRequest)
			return
		}
	}

	if incomingJson.KeyID == "" {
		jsonResponse["response"] = "Bad Request, keyID needs to be filled"
		sendResponse(res, map[string]string{"response": "Bad Request, keyID needs to be filled"}, http.StatusBadRequest)
		return
	}

	keyID := incomingJson.KeyID

	boolRes, mapRes := db.DisableKey(keyID)
	if !boolRes {
		sendResponse(res, mapRes, http.StatusBadRequest)
	} else {
		sendResponse(res, mapRes, http.StatusAccepted)
	}
}
