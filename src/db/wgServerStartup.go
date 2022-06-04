package db

import (
	"fmt"

	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func WGStart() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Starting your wireguard network interface")

	var wgInterface WireguardInterface
	db := DBSystem
	result := db.Where("interface_name = ?", "wg0").First(&wgInterface) //find interface in sqlite db

	if result.Error != nil { //if an interface is not found, create one
		pkServer, errPk := wgtypes.GeneratePrivateKey()
		if errPk != nil {
			combinedLogger.Fatal("Generating new private key " + errPk.Error())
		}

		pubServer := pkServer.PublicKey()

		ipv4Addr := viper.GetString("INSTANCE.IP.LOCAL.IPV4.ADDRESS")
		ipv4Subnet := viper.GetString("INSTANCE.IP.LOCAL.IPV4.SUBNET")
		ipv6Addr := viper.GetString("INSTANCE.IP.LOCAL.IPV6.ADDRESS")
		ipv6Enabled := viper.GetBool("INSTANCE.IP.LOCAL.IPV6.ENABLED")

		wgPort := viper.GetInt("INSTANCE.PORT")

		if ipv6Enabled {
			ipv6Subnet := viper.GetString("INSTANCE.IP.LOCAL.IPV6.SUBNET")
			createWG(pkServer.String(), pubServer.String(), wgPort, ipv4Addr+ipv4Subnet, ipv6Addr+ipv6Subnet)
		} else {
			createWG(pkServer.String(), pubServer.String(), wgPort, ipv4Addr+ipv4Subnet, "-")
		}

		peers := generatePeerArray()
		manager.AddPeersInterface("wg0", pkServer.String(), wgPort, peers)
		combinedLogger.Info("Created wireguard instance on port " + fmt.Sprint(wgPort))
		return
	} else {
		combinedLogger.Info("Wireguard instance in database was found - overriding some values.")
		combinedLogger.Info("Created wireguard instance on port " + fmt.Sprint(wgInterface.ListenPort))
	}

	peers := generatePeerArray()
	manager.AddPeersInterface("wg0", wgInterface.PrivateKey, wgInterface.ListenPort, peers)
}
