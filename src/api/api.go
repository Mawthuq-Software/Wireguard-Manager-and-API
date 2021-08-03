package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

type keyAdd struct {
	PublicKey    string `json:"publicKey"`
	PresharedKey string `json:"presharedKey"`
	KeyID        string `json:"keyID"`
}
type keyDelete struct {
	KeyID string `json:"keyID"`
}

type authStruct struct {
	Active bool `json:"active"`
}

func API() {
	router := mux.NewRouter()                                       //Router for routes
	router.HandleFunc("/manager/keys", keyCreate).Methods("POST")   //post route for adding keys
	router.HandleFunc("/manager/keys", keyRemove).Methods("DELETE") //delete route for removing keys

	serverDev := os.Getenv("SERVER_SECURITY")
	if serverDev == "disabled" {
		port := os.Getenv("PORT")
		fmt.Printf("Info - HTTP about to listen on %s.", port)
		log.Printf("Info - HTTP about to listen on %s.", port)

		err := http.ListenAndServe(":"+port, router)
		log.Fatal("Error - Startup of API server", err)
	} else {
		port := os.Getenv("PORT")
		fullchainCert := os.Getenv("FULLCHAIN_CERT")
		privKeyCert := os.Getenv("PK_CERT")

		log.Printf("HTTPS about to listen on %s.", port)
		fmt.Printf("HTTPS about to listen on %s.", port)

		err := http.ListenAndServeTLS(":"+port,
			fullchainCert, //fullchain
			privKeyCert, router)
		log.Fatal("Error - Startup of API server", err)
	}

}

func parseResponse(req *http.Request, schema interface{}) error {
	var unmarshalErr *json.UnmarshalTypeError

	headerContentTtype := req.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		return errors.New("content type is not application/json")
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields() //throws error if uneeded JSON is added
	err := decoder.Decode(schema)   //decodes the incoming JSON into the struct
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return errors.New("Bad Request. Wrong Type provided for field " + unmarshalErr.Field)
		} else {
			return errors.New("Bad Request " + err.Error())
		}
	}
	return nil
}

func keyCreate(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyAdd

	err := parseResponse(req, &incomingJson) //parse JSON
	if err != nil {
		log.Println("Error - Parsing request", err)
		response(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if incomingJson.PresharedKey == "" || incomingJson.PublicKey == "" {
		response(res, map[string]string{"response": "Bad Request, presharedKey and publicKey must be filled"}, http.StatusBadRequest)
		return
	}
	if os.Getenv("AUTH") != "-" { //check AUTH
		authHeader := req.Header.Get("Authorization")
		if os.Getenv("AUTH") != authHeader {
			response(res, map[string]string{"response": "Authentication key is not valid"}, http.StatusBadRequest)
			return
		}
	}

	boolRes, mapRes := db.CreateKey(incomingJson.PublicKey, incomingJson.PresharedKey) //add key to db
	if !boolRes {
		response(res, mapRes, http.StatusBadRequest)
	} else {
		response(res, mapRes, http.StatusAccepted)
	}
}

func keyRemove(res http.ResponseWriter, req *http.Request) {
	var incomingJson keyDelete

	jsonResponse := make(map[string]string)

	err := parseResponse(req, &incomingJson)
	if err != nil {
		log.Println(err)
		response(res, map[string]string{"response": err.Error()}, http.StatusBadRequest)
		return
	}

	if incomingJson.KeyID == "" {
		jsonResponse["response"] = "Bad Request, keyID needs to be filled"
		response(res, map[string]string{"response": "Bad Request, keyID needs to be filled"}, http.StatusBadRequest)
		return
	}

	if os.Getenv("AUTH") != "-" { //check AUTH
		authHeader := req.Header.Get("Authorization")
		if os.Getenv("AUTH") != authHeader {
			response(res, map[string]string{"response": "Authentication key is not valid"}, http.StatusBadRequest)
			return
		}
	}

	boolRes, mapRes := db.DeleteKey(incomingJson.KeyID)
	if !boolRes {
		response(res, mapRes, http.StatusBadRequest)
	} else {
		//need to add in pubKey of server, dns, allowedIP
		response(res, mapRes, http.StatusAccepted)
	}
}

func response(res http.ResponseWriter, responseMap map[string]string, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(responseMap)
	res.Write(jsonResp)
}
