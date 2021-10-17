package router

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter() //Router for routes

	manager := router.PathPrefix("/manager").Subrouter() //main subrouter

	keys := manager.PathPrefix("/keys").Subrouter()         //specific subrouter
	keys.HandleFunc("", keyCreate).Methods("POST")          //post route for adding keys
	keys.HandleFunc("", keyRemove).Methods("DELETE")        //delete route for removing keys
	keys.HandleFunc("/enable", keyEnable).Methods("POST")   //post route for enabling key
	keys.HandleFunc("/disable", keyDisable).Methods("POST") //post route for disabling key

	subscriptions := manager.PathPrefix("/subscriptions").Subrouter()     //specific subrouter
	subscriptions.HandleFunc("/edit", keySetSubscription).Methods("POST") //for editing subscription
	return router
}
