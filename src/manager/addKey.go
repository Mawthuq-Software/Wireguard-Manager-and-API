package manager

import (
	"log"
	"net"
	"time"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func AddKey(interfaceName string, ipv4Address string, ipv6Address string, publicStr string, presharedStr string) (bool, string) {
	var ipAddresses []net.IPNet
	var zeroTime time.Duration
	var arrayConfig []wgtypes.PeerConfig

	if ipv4Address != "-" { //check for unneeded IP
		ipv4Address = ipv4Address + "/32"     //add subnet to IP
		ipv4, errIpv4 := ParseIP(ipv4Address) //Parse IP into readable form
		if !logger.ErrorHandler("Error - Parsing IPv4 Address", errIpv4) {
			return false, "An error has occurred when parsing the IPv4 Address"
		}

		ipAddresses = append(ipAddresses, *ipv4) //add IP to array
	}
	if ipv6Address != "-" { //check for unneeded IP
		ipv6Address = ipv6Address + "/128" //add subnet to IP
		ipv6, errIpv6 := ParseIP(ipv6Address)
		if !logger.ErrorHandler("Error - Parsing IPv4 Address", errIpv6) {
			return false, "An error has occurred when parsing the IPv6 Address"
		}
		ipAddresses = append(ipAddresses, *ipv6) //add IP to array
	}

	publicKey, errPubParse := ParseKey(publicStr) //parse string into readable form
	if !logger.ErrorHandler("Error - Parsing public key", errPubParse) {
		return false, "An error has occurred when parsing the server public key"
	}
	presharedKey, errPreParse := ParseKey(presharedStr) //parse string into readable form
	if !logger.ErrorHandler("Error - Parsing preshared key", errPreParse) {
		return false, "An error has occurred when parsing the preshared key"
	}

	userConfig := wgtypes.PeerConfig{ //setup client config for server
		PublicKey:                   publicKey,
		PresharedKey:                &presharedKey,
		PersistentKeepaliveInterval: &zeroTime,
		AllowedIPs:                  ipAddresses,
	}
	arrayConfig = append(arrayConfig, userConfig) //add client config to array of configs

	client, errInstance := createInstance() //new client to communicate with wireguard device
	if errInstance != nil {
		log.Println("Error - Creating instance", errInstance)
		return false, "An error has occurred when creating a WG instance"
	}

	errConfigure := client.ConfigureDevice(interfaceName, wgtypes.Config{ //add new peers to wg interface
		Peers:        arrayConfig,
		ReplacePeers: false,
	})

	if !logger.ErrorHandler("Configuring device on add key", errConfigure) {
		return false, "An error has occurred when configuring the device"
	}
	closeInstance(client) //release resources used by client
	return true, "Successfully added key"
}
