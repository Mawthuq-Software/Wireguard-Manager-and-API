package manager

import (
	"github.com/vishvananda/netlink"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func AddPeersInterface(interfaceName string, pk string, listenPort int, peers []wgtypes.PeerConfig) { //adds client peers to wg interface
	combinedLogger := logger.GetCombinedLogger()

	client, errClient := createInstance()
	if errClient != nil {
		combinedLogger.Error("Creating new client " + errClient.Error())
	}
	devices, errDev := client.Devices() //get all wireguard devices
	if errDev != nil {
		combinedLogger.Error("Retrieving devices " + errDev.Error())
	}
	wg0Found := false
	for i := 0; i < len(devices); i++ { //find if wg0 interface exists
		if devices[i].Name == "wg0" {
			wg0Found = true
		}
	}

	if !wg0Found {
		combinedLogger.Info("Adding new wg interface, none was found")
		pkEncoded, pkErr := wgtypes.ParseKey(pk)
		if pkErr != nil {
			combinedLogger.Error("Parsing private key " + pkErr.Error())
		}

		handle, errHandle := netlink.NewHandle() //create new handle to add link
		if errHandle != nil {
			combinedLogger.Error("Creating new handle " + errHandle.Error())
		}

		linkA := netlink.NewLinkAttrs() //create new interface attributes
		linkA.Name = "wg0"              //set the name for the interface
		linkA.MTU = 1420
		linkA.TxQLen = 1000

		linkWG := wgInterface{Attributes: &linkA, TypeName: "wireguard"} //create new interface

		realLink := returnNewLink(linkWG)      //get new link
		errAddLink := handle.LinkAdd(realLink) //add new link
		if errAddLink != nil {
			combinedLogger.Error("Creating new link " + errAddLink.Error())
		}

		wgConfig := wgtypes.Config{ //setup wireguard interface
			PrivateKey:   &pkEncoded,
			ListenPort:   &listenPort,
			ReplacePeers: false,
			Peers:        peers,
		}
		errConfDev := client.ConfigureDevice(interfaceName, wgConfig)
		if errConfDev != nil {
			combinedLogger.Error("Configuring device " + errConfDev.Error())
		}
	} else {
		combinedLogger.Info("Interface exists, adding peers")

		wgConfig := wgtypes.Config{
			ListenPort:   &listenPort,
			ReplacePeers: true,
			Peers:        peers,
		}

		errConfDev := client.ConfigureDevice(interfaceName, wgConfig)
		if errConfDev != nil {
			combinedLogger.Error("Configuring device " + errConfDev.Error())
		}
	}
	//code for listing all clients on interfaces
	/*getInterfaces, errInt := client.Devices()
	if errInt != nil {
		log.Println(errInt)
	}
	for interfaces := 0; interfaces < len(getInterfaces); interfaces++ { //get interfaces
		for peer := 0; peer < len(getInterfaces[interfaces].Peers); peer++ { //get each peer in the interface
			currentPeer := getInterfaces[interfaces].Peers[peer]
			log.Println(currentPeer)
		}
	}*/
	handles, errHandle := netlink.NewHandle() //create new handle to add link
	if errHandle != nil {
		combinedLogger.Error("PCreating new handle " + errHandle.Error())
	}
	realLinkTwo, errUp := handles.LinkByName("wg0")
	if errUp != nil {
		combinedLogger.Error("Link Up " + errUp.Error())
	}
	errLinkUp := handles.LinkSetUp(realLinkTwo)
	if errLinkUp != nil {
		combinedLogger.Error("Link Up " + errLinkUp.Error())
	}
	closeInstance(client)
}
