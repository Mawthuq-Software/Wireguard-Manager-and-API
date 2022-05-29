package api

import (
	"fmt"
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
		combinedLogger.Info(fmt.Sprintf("HTTP about to listen on %s.", port))

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.Serve(resolveTCP, newRouter)
		combinedLogger.Error(fmt.Sprintf("Failed to startup of API server %s", errServer))
	} else {
		port := viper.GetString("SERVER.PORT")
		fullchainCert := viper.GetString("SERVER.CERT.FULLCHAIN")
		privKeyCert := viper.GetString("SERVER.CERT.PK")

		combinedLogger.Info(fmt.Sprintf("HTTPS about to listen on %s.", port))

		resolve, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:"+port)
		resolveTCP, _ := net.ListenTCP("tcp4", resolve)

		errServer := http.ServeTLS(resolveTCP, newRouter, fullchainCert, privKeyCert)
		combinedLogger.Error(fmt.Sprintf("Failed to startup of API server %s", errServer))
	}
}
