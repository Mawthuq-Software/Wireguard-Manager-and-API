package router

import (
	"log"
	"net/http"
	"os"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keyCreateJSON struct {
	PublicKey    string `json:"publicKey"`
	PresharedKey string `json:"presharedKey"`
	KeyID        string `json:"keyID"`
}

func keyCreate(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyCreateJSON

	err := parseResponse(req, &incomingJson) //parse JSON
	if err != nil {
		log.Println("Error - Parsing request", err)
		sendResponse(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if incomingJson.PresharedKey == "" || incomingJson.PublicKey == "" {
		sendResponse(res, map[string]string{"response": "Bad Request, presharedKey and publicKey must be filled"}, http.StatusBadRequest)
		return
	}
	if os.Getenv("AUTH") != "-" { //check AUTH
		authHeader := req.Header.Get("Authorization")
		if os.Getenv("AUTH") != authHeader {
			sendResponse(res, map[string]string{"response": "Authentication key is not valid"}, http.StatusBadRequest)
			return
		}
	}

	boolRes, mapRes := db.CreateKey(incomingJson.PublicKey, incomingJson.PresharedKey) //add key to db
	if !boolRes {
		sendResponse(res, mapRes, http.StatusBadRequest)
	} else {
		sendResponse(res, mapRes, http.StatusAccepted)
	}
}