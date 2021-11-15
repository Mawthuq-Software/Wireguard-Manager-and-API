package router

import (
	"net/http"

	"github.com/spf13/viper"
)

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		auth := viper.GetString("SERVER.AUTH")
		if auth != "-" { //check AUTH
			authHeader := req.Header.Get("Authorization")
			if auth != authHeader {
				sendResponse(res, map[string]string{"response": "Authentication key is not valid"}, http.StatusBadRequest)
				return
			}
		}

		handler.ServeHTTP(res, req)
	})
}
