package router

import (
	"net/http"
)

func setHeader(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		res.Header().Set("Access-Control-Allow-Origin", "*")

		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		}
		hand.ServeHTTP(res, req)
	})
}

func setCorsHeader(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusOK)
		return
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}
