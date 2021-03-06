package router

import (
	"log"
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keyRemoveJSON struct {
	KeyID string `json:"keyID"`
}

func keyRemove(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyRemoveJSON

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

	boolRes, mapRes := db.DeleteKey(incomingJson.KeyID)
	if !boolRes {
		sentStandardRes(res, mapRes, http.StatusBadRequest)
	} else {
		//need to add in pubKey of server, dns, allowedIP
		sentStandardRes(res, mapRes, http.StatusAccepted)
	}
}
