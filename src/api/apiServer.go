package api

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/api/router"
)

type authStruct struct {
	Active bool `json:"active"`
}

func API() {
	newRouter := router.NewRouter()

	serverDev := viper.GetBool("SERVER.SECURITY")
	if !serverDev {
		port := viper.GetString("SERVER.PORT")
		fmt.Printf("Info - HTTP about to listen on %s.", port)
		log.Printf("Info - HTTP about to listen on %s.", port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.Serve(resolveTCP, newRouter)
		log.Fatal("Error - Startup of API server", errServer)
	} else {
		port := viper.GetString("SERVER.PORT")
		fullchainCert := viper.GetString("SERVER.CERT.FULLCHAIN")
		privKeyCert := viper.GetString("SERVER.CERT.PK")

		log.Printf("HTTPS about to listen on %s.", port)
		fmt.Printf("HTTPS about to listen on %s.", port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.ServeTLS(resolveTCP, newRouter, fullchainCert, privKeyCert)
		log.Fatal("Error - Startup of API server", errServer)
	}
}
