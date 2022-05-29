package manager

import (
	"fmt"
	"time"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func AddRemovePeer(peer wgtypes.Peer, wgIntName string) bool { //add and readds clients based on if a time to last connection has been reached
	client, errInstance := createInstance()
	combinedLogger := logger.GetCombinedLogger()

	if errInstance != nil {
		combinedLogger.Error(fmt.Sprintf("Create instance %s", errInstance))
		return false
	}
	//get each peer in the wg interface
	lastConnection := peer.LastHandshakeTime                      //get handshake time of peer
	currentTime := time.Now().UTC()                               //get current time
	zeroTime, errTime := time.Parse("2006-Jan-02", "0001-Jan-01") //parse time into format
	if errTime != nil {
		return false
	}

	handshakeAfterDeadline := lastConnection.Add(time.Minute * 1)                     //add 5 minutes to last connection
	if currentTime.After(handshakeAfterDeadline) && !lastConnection.Equal(zeroTime) { //compare now and 5 minutes after last connection
		allowedIPs := peer.AllowedIPs        //get IPs of client
		publicKey := peer.PublicKey          //get public key of client
		presharedKey := peer.PresharedKey    //get preshared key of client
		wgInterface := wgIntName             //get interface name client is on
		userConfigDel := wgtypes.PeerConfig{ //remove key
			PublicKey: publicKey,
			Remove:    true,
		}
		var arrayConfigDel []wgtypes.PeerConfig
		arrayConfigDel = append(arrayConfigDel, userConfigDel) //add config to array

		errConf := client.ConfigureDevice(wgInterface, wgtypes.Config{ //configure wg device
			Peers: arrayConfigDel,
		})
		if !logger.ErrorHandler("Error - Configuring device on key deletion", errConf) {
			return false
		}

		var zeroTime time.Duration        //nil time
		userConfig := wgtypes.PeerConfig{ //config to add back into device
			PublicKey:                   publicKey,
			PresharedKey:                &presharedKey,
			PersistentKeepaliveInterval: &zeroTime,
			AllowedIPs:                  allowedIPs,
		}
		var arrayConfigAdd []wgtypes.PeerConfig
		arrayConfigAdd = append(arrayConfigAdd, userConfig) //add config into array

		errWG := client.ConfigureDevice(wgInterface, wgtypes.Config{ //configure wg device
			Peers: arrayConfigAdd,
		})
		if !logger.ErrorHandler("Error - Configuring device on add key", errWG) {
			return false
		}

	} else {
		return false
	}
	closeInstance(client) //release resources and close instance
	return true
}
