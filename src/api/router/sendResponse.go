package router

import (
	"encoding/json"
	"net/http"
)

func sendResponse(res http.ResponseWriter, responseMap map[string]string, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(responseMap)
	res.Write(jsonResp)
}
