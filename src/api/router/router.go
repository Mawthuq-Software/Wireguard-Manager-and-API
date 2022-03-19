package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter() //Router for routes
	router.Use(setHeader)     //need to allow CORS and OPTIONS
	router.Use(authMiddleware)

	manager := router.PathPrefix("/manager").Subrouter() //main subrouter

	keys := manager.PathPrefix("/key").Subrouter() //specific subrouter
	keys.HandleFunc("", getKeys).Methods("GET")
	keys.HandleFunc("", keyCreate).Methods("POST")          //post route for adding keys
	keys.HandleFunc("", keyRemove).Methods("DELETE")        //delete route for removing keys
	keys.HandleFunc("/enable", keyEnable).Methods("POST")   //post route for enabling key
	keys.HandleFunc("/disable", keyDisable).Methods("POST") //post route for disabling key

	subscriptions := manager.PathPrefix("/subscription").Subrouter() //specific subrouter
	subscriptions.HandleFunc("/all", getSubscriptions).Methods("GET")
	subscriptions.HandleFunc("/edit", keySetSubscription).Methods("POST") //for editing subscription
	subscriptions.HandleFunc("", getKeySub).Methods("GET")

	router.MethodNotAllowedHandler = http.HandlerFunc(setCorsHeader) //if method is not found allow OPTIONS
	return router
}
