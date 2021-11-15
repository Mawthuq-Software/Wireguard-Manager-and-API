package router

import (
	"net/http"
	"os"
)

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if os.Getenv("AUTH") != "-" { //check AUTH
			authHeader := req.Header.Get("Authorization")
			if os.Getenv("AUTH") != authHeader {
				sendResponse(res, map[string]string{"response": "Authentication key is not valid"}, http.StatusBadRequest)
				return
			}
		}

		handler.ServeHTTP(res, req)
	})
}
