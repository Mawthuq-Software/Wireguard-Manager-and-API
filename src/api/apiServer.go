package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/api/router"
)

type authStruct struct {
	Active bool `json:"active"`
}

func API() {
	newRouter := router.NewRouter()

	serverDev := os.Getenv("SERVER_SECURITY")
	if serverDev == "disabled" {
		port := os.Getenv("PORT")
		fmt.Printf("Info - HTTP about to listen on %s.", port)
		log.Printf("Info - HTTP about to listen on %s.", port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:8443")
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.Serve(resolveTCP, newRouter)
		log.Fatal("Error - Startup of API server", errServer)
	} else {
		port := os.Getenv("PORT")
		fullchainCert := os.Getenv("FULLCHAIN_CERT")
		privKeyCert := os.Getenv("PK_CERT")

		log.Printf("HTTPS about to listen on %s.", port)
		fmt.Printf("HTTPS about to listen on %s.", port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:8443")
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.ServeTLS(resolveTCP, newRouter, fullchainCert, privKeyCert)
		log.Fatal("Error - Startup of API server", errServer)
	}
}
