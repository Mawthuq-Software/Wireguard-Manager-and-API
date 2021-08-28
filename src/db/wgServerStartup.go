package db

import (
	"log"
	"os"

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

		ipv4Addr := os.Getenv("WG_IPV4")
		ipv6Addr := os.Getenv("WG_IPV6")

		if ipv6Addr != "-" {
			createWG(pkServer.String(), pubServer.String(), 51820, ipv4Addr+"/16", ipv6Addr+"/64")
		} else {
			createWG(pkServer.String(), pubServer.String(), 51820, ipv4Addr+"/16", "-")
		}

		peers := generatePeerArray()
		manager.AddPeersInterface("wg0", pkServer.String(), 51820, peers)
		return
	}

	peers := generatePeerArray()
	manager.AddPeersInterface("wg0", wgInterface.PrivateKey, wgInterface.ListenPort, peers)
}
