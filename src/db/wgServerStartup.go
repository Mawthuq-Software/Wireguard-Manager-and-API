package db

import (
	"log"

	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func WGStart() {
	log.Println("Info - Starting up wg interface")

	var wgInterface WireguardInterface
	db := DBSystem
	result := db.Where("interface_name = ?", "wg0").First(&wgInterface) //find interface in sqlite db

	if result.Error != nil { //if an interface is not found, create one
		pkServer, errPk := wgtypes.GeneratePrivateKey()
		if errPk != nil {
			log.Fatal("Error - Generating new private key", errPk)
		}

		pubServer := pkServer.PublicKey()

		ipv4Addr := viper.GetString("INSTANCE.IP.LOCAL.IPV4.ADDRESS")
		ipv4Subnet := viper.GetString("INSTANCE.IP.LOCAL.IPV4.SUBNET")
		ipv6Addr := viper.GetString("INSTANCE.IP.LOCAL.IPV6.ADDRESS")
		ipv6Enabled := viper.GetBool("INSTANCE.IP.LOCAL.IPV6.ENABLED")
		if ipv6Enabled {
			ipv6Subnet := viper.GetString("INSTANCE.IP.LOCAL.IPV6.SUBNET")
			createWG(pkServer.String(), pubServer.String(), 51820, ipv4Addr+ipv4Subnet, ipv6Addr+ipv6Subnet)
		} else {
			createWG(pkServer.String(), pubServer.String(), 51820, ipv4Addr+ipv4Subnet, "-")
		}

		peers := generatePeerArray()
		manager.AddPeersInterface("wg0", pkServer.String(), 51820, peers)
		return
	}

	peers := generatePeerArray()
	manager.AddPeersInterface("wg0", wgInterface.PrivateKey, wgInterface.ListenPort, peers)
}
