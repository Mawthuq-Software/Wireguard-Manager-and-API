package api

import (
	"net"
	"net/http"

	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/api/router"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

type authStruct struct {
	Active bool `json:"active"`
}

func API() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Starting web server")

	newRouter := router.NewRouter()

	serverDev := viper.GetBool("SERVER.SECURITY")
	if !serverDev {
		port := viper.GetString("SERVER.PORT")
		combinedLogger.Info("HTTP about to listen on " + port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.Serve(resolveTCP, newRouter)
		combinedLogger.Error("Failed to startup of API server " + errServer.Error())
	} else {
		port := viper.GetString("SERVER.PORT")
		fullchainCert := viper.GetString("SERVER.CERT.FULLCHAIN")
		privKeyCert := viper.GetString("SERVER.CERT.PK")

		combinedLogger.Info("HTTP about to listen on " + port)

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.ServeTLS(resolveTCP, newRouter, fullchainCert, privKeyCert)
		combinedLogger.Error("Failed to startup of API server " + errServer.Error())
	}
}
