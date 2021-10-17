package manager

import (
	"log"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func DeleteKey(interfaceName string, publicStr string) (bool, string) {
	var arrayConfig []wgtypes.PeerConfig //array of config to be removed

	client, errInstance := createInstance() //create new communication wg device
	if errInstance != nil {
		log.Println("Create instance", errInstance)
		return false, "An error has occurred when creating a WG instance"
	}

	publicKey, err := wgtypes.ParseKey(publicStr)
	if !logger.ErrorHandler("Parsing public key on delete key", err) {
		return false, "An error has occurred when parsing the public key"
	}
	userConfig := wgtypes.PeerConfig{ //create config object
		PublicKey: publicKey,
		Remove:    true,
	}
	arrayConfig = append(arrayConfig, userConfig) //add user config to array to be parsed

	err = client.ConfigureDevice(interfaceName, wgtypes.Config{
		Peers: arrayConfig,
	})
	if !logger.ErrorHandler("Configuring device on delete key", err) {
		return false, "An error has occurred when configuring the device"
	}
	closeInstance(client) //close and release resources from communication device
	return true, "Removed the key successfully"
}
