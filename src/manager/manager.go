package manager

import (
	"log"
	"net"
	"time"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Link interface { //https://github.com/vishvananda/netlink/blob/master/link.go
	Attrs() *netlink.LinkAttrs
	Type() string
}
type wgInterface struct {
	Attributes *netlink.LinkAttrs
	TypeName   string
}

func createInstance() (*wgctrl.Client, error) {
	log.Println("Info - Creating wg device client")
	return wgctrl.New()
}
func closeInstance(client *wgctrl.Client) error {
	log.Println("Info - Closing wg device client")
	return client.Close()
}
func AddPeersInterface(interfaceName string, pk string, listenPort int, peers []wgtypes.PeerConfig) { //adds client peers to wg interface
	client, errClient := createInstance()
	if errClient != nil {
		log.Fatal("Error - Creating new client", errClient)
	}
	devices, errDev := client.Devices() //get all wireguard devices
	if errDev != nil {
		log.Fatal("Error - Retrieving devices", errDev)
	}
	wg0Found := false
	for i := 0; i < len(devices); i++ { //find if wg0 interface exists
		if devices[i].Name == "wg0" {
			wg0Found = true
		}
	}

	if !wg0Found {
		log.Println("Info - Adding new wg interface, none was found")
		pkEncoded, pkErr := wgtypes.ParseKey(pk)
		if pkErr != nil {
			log.Fatal("Error - Parsing private key", pkErr)
		}

		handle, errHandle := netlink.NewHandle() //create new handle to add link
		if errHandle != nil {
			log.Fatal("Error - Creating new handle", errHandle)
		}

		linkA := netlink.NewLinkAttrs() //create new interface attributes
		linkA.Name = "wg0"              //set the name for the interface
		linkA.MTU = 1420
		linkA.TxQLen = 1000

		linkWG := wgInterface{Attributes: &linkA, TypeName: "wireguard"} //create new interface

		realLink := returnNewLink(linkWG)      //get new link
		errAddLink := handle.LinkAdd(realLink) //add new link
		if errAddLink != nil {
			log.Fatal("Error - Creating new link", errAddLink)
		}

		wgConfig := wgtypes.Config{ //setup wireguard interface
			PrivateKey:   &pkEncoded,
			ListenPort:   &listenPort,
			ReplacePeers: false,
			Peers:        peers,
		}
		errConfDev := client.ConfigureDevice(interfaceName, wgConfig)
		if errConfDev != nil {
			log.Fatal("Error - Configuring device", errConfDev)
		}
	} else {
		log.Println("Info - Interface exists, adding peers")
		wgConfig := wgtypes.Config{
			ListenPort:   &listenPort,
			ReplacePeers: true,
			Peers:        peers,
		}

		errConfDev := client.ConfigureDevice(interfaceName, wgConfig)
		if errConfDev != nil {
			log.Fatal("Error - Configuring device", errConfDev)
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
		log.Fatal("Error - Creating new handle", errHandle)
	}
	realLinkTwo, errUp := handles.LinkByName("wg0")
	if errUp != nil {
		log.Fatal("Error - Link Up", errUp)
	}
	errLinkUp := handles.LinkSetUp(realLinkTwo)
	if errLinkUp != nil {
		log.Fatal("Error - Link Up", errLinkUp)
	}
	closeInstance(client)
}
func (attr wgInterface) Attrs() *netlink.LinkAttrs {
	return attr.Attributes
}
func (t wgInterface) Type() string {
	return t.TypeName
}
func returnNewLink(l Link) netlink.Link {
	return l
}
func AddKey(interfaceName string, ipv4Address string, ipv6Address string, publicStr string, presharedStr string) (bool, string) {
	var ipAddresses []net.IPNet
	var zeroTime time.Duration
	var arrayConfig []wgtypes.PeerConfig

	if ipv4Address != "-" { //check for unneeded IP
		ipv4Address = ipv4Address + "/32"     //add subnet to IP
		ipv4, errIpv4 := ParseIP(ipv4Address) //Parse IP into readable form
		if !errorHandler("Error - Parsing IPv4 Address", errIpv4) {
			return false, "An error has occurred when parsing the IPv4 Address"
		}

		ipAddresses = append(ipAddresses, *ipv4) //add IP to array
	}
	if ipv6Address != "-" { //check for unneeded IP
		ipv6Address = ipv6Address + "/128" //add subnet to IP
		ipv6, errIpv6 := ParseIP(ipv6Address)
		if !errorHandler("Error - Parsing IPv4 Address", errIpv6) {
			return false, "An error has occurred when parsing the IPv6 Address"
		}
		ipAddresses = append(ipAddresses, *ipv6) //add IP to array
	}

	publicKey, errPubParse := ParseKey(publicStr) //parse string into readable form
	if !errorHandler("Error - Parsing public key", errPubParse) {
		return false, "An error has occurred when parsing the server public key"
	}
	presharedKey, errPreParse := ParseKey(presharedStr) //parse string into readable form
	if !errorHandler("Error - Parsing preshared key", errPreParse) {
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

	if !errorHandler("Configuring device on add key", errConfigure) {
		return false, "An error has occurred when configuring the device"
	}
	closeInstance(client) //release resources used by client
	return true, "Successfully added key"
}
func DeleteKey(interfaceName string, publicStr string) (bool, string) {
	var arrayConfig []wgtypes.PeerConfig //array of config to be removed

	client, errInstance := createInstance() //create new communication wg device
	if errInstance != nil {
		log.Println("Create instance", errInstance)
		return false, "An error has occurred when creating a WG instance"
	}

	publicKey, err := wgtypes.ParseKey(publicStr)
	if !errorHandler("Parsing public key on delete key", err) {
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
	if !errorHandler("Configuring device on delete key", err) {
		return false, "An error has occurred when configuring the device"
	}
	closeInstance(client) //close and release resources from communication device
	return true, "Removed the key successfully"
}
func AddRemovePeer() bool { //add and readds clients based on if a time to last connection has been reached
	client, errInstance := createInstance()
	if errInstance != nil {
		log.Println("Create instance", errInstance)
		return false
	}
	getInterfaces, err := client.Devices()
	if !errorHandler("Info - Finding interfaces", err) {
		return false
	}
	for interfaces := 0; interfaces < len(getInterfaces); interfaces++ { //get interfaces
		for peer := 0; peer < len(getInterfaces[interfaces].Peers); peer++ { //get each peer in the wg interface
			currentPeer := getInterfaces[interfaces].Peers[peer]          //get the current peer in for loop
			lastConnection := currentPeer.LastHandshakeTime               //get handshake time of peer
			currentTime := time.Now().UTC()                               //get current time
			zeroTime, errTime := time.Parse("2006-Jan-02", "0001-Jan-01") //parse time into format
			if errTime != nil {
				continue //on error pass over to next client
			}

			handshakeAfterDeadline := lastConnection.Add(time.Minute * 5)                     //add 5 minutes to last connection
			if currentTime.After(handshakeAfterDeadline) && !lastConnection.Equal(zeroTime) { //compare now and 5 minutes after last connection
				allowedIPs := currentPeer.AllowedIPs          //get IPs of client
				publicKey := currentPeer.PublicKey            //get public key of client
				presharedKey := currentPeer.PresharedKey      //get preshared key of client
				wgInterface := getInterfaces[interfaces].Name //get interface name client is on

				userConfigDel := wgtypes.PeerConfig{ //remove key
					PublicKey: publicKey,
					Remove:    true,
				}
				var arrayConfigDel []wgtypes.PeerConfig
				arrayConfigDel = append(arrayConfigDel, userConfigDel) //add config to array

				err = client.ConfigureDevice(wgInterface, wgtypes.Config{ //configure wg device
					Peers: arrayConfigDel,
				})
				if !errorHandler("Error - Configuring device on key deletion", err) {
					continue //continue even on error
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

				err = client.ConfigureDevice(wgInterface, wgtypes.Config{ //configure wg device
					Peers: arrayConfigAdd,
				})
				if !errorHandler("Error - Configuring device on add key", err) {
					continue //continue on error
				}
			} else {
				continue
			}
		}
	}
	closeInstance(client) //release resources and close instance
	return true
}
func ParseKey(key string) (parsedKey wgtypes.Key, err error) { //parses string into key
	parsedKey, err = wgtypes.ParseKey(key)
	if err != nil {
		return
	}
	return
}
func ParseIP(address string) (ipvX *net.IPNet, err error) { //parses string into IP address
	_, ipvX, err = net.ParseCIDR(address)
	if err != nil {
		return
	}
	return
}
func errorHandler(message string, err error) bool { //error handler
	if err != nil {
		log.Println(message, err)
		return false
	}
	return true
}
