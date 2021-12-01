package router

import (
	"encoding/json"
	"net/http"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

func sentStandardRes(res http.ResponseWriter, responseMap map[string]string, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(responseMap)
	res.Write(jsonResp)
}

func sendKeysRes(res http.ResponseWriter, responseStruct db.ReturnKeysRes, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(responseStruct)
	res.Write(jsonResp)
}
