package db

import (
	"errors"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Key struct {
	KeyID        int    `gorm:"primaryKey;autoIncrement"`
	PublicKey    string `gorm:"unique"`
	PresharedKey string `gorm:"unique"`
	IPv4Address  string `gorm:"foreignKey:IPv4Address"`
}
type IP struct {
	IPv4Address string `gorm:"primaryKey"`
	IPv6Address string `gorm:"unique"`
	InUse       string
	WGInterface string
}
type WireguardInterface struct {
	InterfaceName string `gorm:"primaryKey"`
	PrivateKey    string `gorm:"unique"`
	PublicKey     string `gorm:"unique"`
	ListenPort    int    `gorm:"unique"`
	IPv4Address   string
	IPv6Address   string
}

var DBSystem *gorm.DB

func DBStart() {
	log.Println("Info - Database connection starting")
	errCreateDir := os.MkdirAll("/opt/wgManagerAPI/wg", 0755) //create dir if not exist
	if errCreateDir != nil {
		log.Fatal("Error - Creating wg directory", errCreateDir)
	}
	db, err := gorm.Open(sqlite.Open("/opt/wgManagerAPI/wg/wireguardPeers.db"), &gorm.Config{}) //open sqlite db
	if err != nil {
		panic("Error - Failed to connect database")
	}

	DBSystem = db //set global variable up
	// Migrate the schema

	errMigrate := db.AutoMigrate(&Key{}, &IP{}, &WireguardInterface{}) //Migrate tables to sqlite
	if errMigrate != nil {
		log.Fatal("Error - Migrating database", errMigrate)
	} else {
		log.Println("Info - Successfully migrated db")
	}
	generateIPs()
}

func WGStart() {
	log.Println("Info - Starting up wg interface")

	var wgInterface WireguardInterface
	db := DBSystem
	result := db.Where("interface_name = ?", "wg0").First(&wgInterface) //find nterface in sqlite db

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

func createWG(PrivateKey string, PublicKey string, ListenPort int, IPv4Address string, IPv6Address string) {
	log.Println("Info - Creating wireguard interface")
	db := DBSystem
	var wgCreate WireguardInterface

	if IPv6Address != "-" {
		wgCreate = WireguardInterface{ //with IPv6
			InterfaceName: "wg0",
			PrivateKey:    PrivateKey,
			PublicKey:     PublicKey,
			ListenPort:    ListenPort,
			IPv4Address:   IPv4Address,
			IPv6Address:   IPv6Address,
		}
	} else {
		wgCreate = WireguardInterface{ //without IPv6
			InterfaceName: "wg0",
			PrivateKey:    PrivateKey,
			PublicKey:     PublicKey,
			ListenPort:    ListenPort,
			IPv4Address:   IPv4Address,
		}
	}

	err := db.Create(wgCreate).Error
	if err != nil {
		log.Fatal("Creating interface", err)
	}
	log.Println("Successfully created interface")
}

func CreateKey(pubKey string, preKey string) (bool, map[string]string) {
	var ipStruct IP
	var wgStruct WireguardInterface
	responseMap := make(map[string]string)
	db := DBSystem

	resultIP := db.Where("in_use = ?", "false").First(&ipStruct) //find IP not in use
	if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "No available IPs on the server"
		return false, responseMap
	}

	keyStructCreate := Key{PublicKey: pubKey, PresharedKey: preKey, IPv4Address: ipStruct.IPv4Address} //create Key object
	resultKeyCreate := db.Create(&keyStructCreate)                                                     //add object to db
	if resultKeyCreate.Error != nil {
		log.Println("Error - Adding key to db", resultKeyCreate.Error)
		responseMap["response"] = "Error when adding key to database"
		return false, responseMap
	}
	ipStruct.InUse = "true"                                                                                             //set ip to in use
	db.Save(&ipStruct)                                                                                                  //update IP in db
	keyIDStr := strconv.Itoa(keyStructCreate.KeyID)                                                                     //convert keyID to string
	boolRes, strRes := manager.AddKey(ipStruct.WGInterface, ipStruct.IPv4Address, ipStruct.IPv6Address, pubKey, preKey) //add key to wg interface
	if !boolRes {                                                                                                       //if an error occurred
		responseMap["response"] = strRes
		return boolRes, responseMap
	} else {
		responseMap["response"] = "Added key successfully"
		responseMap["ipv4Address"] = ipStruct.IPv4Address + "/32"
		if ipStruct.IPv6Address != "-" {
			responseMap["ipv6Address"] = ipStruct.IPv6Address + "/128"
		}
		responseMap["ipAddress"] = os.Getenv("IP_ADDRESS")
		responseMap["dns"] = os.Getenv("DNS")
		responseMap["allowedIPs"] = os.Getenv("ALLOWED_IP")
		responseMap["keyID"] = keyIDStr
	}

	resultWG := db.Where("interface_name = ?", "wg0").First(&wgStruct) //get wireguard server info

	if resultWG.Error != nil {
		responseMap["response"] = "Issue in finding a key for the server"
		return false, responseMap
	} else {
		responseMap["publicKey"] = wgStruct.PublicKey                 //return back wg server pub key
		responseMap["listenPort"] = strconv.Itoa(wgStruct.ListenPort) //return back wg server listenPort
		return boolRes, responseMap
	}
}

func DeleteKey(keyID string) (bool, map[string]string) {
	var ipStruct IP
	var keyStruct Key
	responseMap := make(map[string]string)
	db := DBSystem

	keyIDInt, _ := strconv.Atoi(keyID)                              //convert key string to int
	resultKey := db.Where("key_id = ?", keyIDInt).First(&keyStruct) //find key in database
	if errors.Is(resultKey.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}

	pubKey := keyStruct.PublicKey                                 //set pub key
	ipv4 := keyStruct.IPv4Address                                 //set ipv4 address
	resultDel := db.Where("key_id = ?", keyID).Delete(&keyStruct) //delete key from db
	if resultDel.Error != nil {
		log.Println("Finding key in DB", resultDel.Error)
		responseMap["response"] = "Error occurred when finding the key in database"
		return false, responseMap
	}

	resultIP := db.Where("ipv4_address = ?", ipv4).First(&ipStruct) //find IP in db
	if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
		responseMap["response"] = "Key was not found on the server"
		return false, responseMap
	}

	ipStruct.InUse = "false"       //set IP back to unused
	ipUpdate := db.Save(&ipStruct) //save data
	if ipUpdate.Error != nil {
		responseMap["response"] = "Error in updating IP"
		return false, responseMap
	}
	boolRes, stringRes := manager.DeleteKey("wg0", pubKey) //delete key from wg interface
	responseMap["response"] = stringRes
	return boolRes, responseMap
}

func generateIPs() {
	log.Println("Info - Generating IPs")

	var availableIPStruct IP
	db := DBSystem

	maxIPStr := os.Getenv("MAX_IP")         //get maximum IPs in env
	maxIPInt, err := strconv.Atoi(maxIPStr) //convert to int

	if err != nil {
		log.Fatal("Unable to convert IP to int", err)
	}

	ipv4Addr := os.Getenv("WG_IPV4")                //IPv4 Subnet Address
	ipv4Splice := strings.SplitAfter(ipv4Addr, ".") //split str at decimal
	ipv4Query := ipv4Splice[0] + ipv4Splice[1]      //get first two subnet

	ipv6Addr := os.Getenv("WG_IPV6") //IPv6 Subnet Address

	if ipv6Addr != "-" {
		ipv6Splice := strings.SplitAfter(ipv6Addr, ":") //split at str at colon
		ipv6Query := ipv6Splice[0] + ipv6Splice[1] + ipv6Splice[2] + ":"
		result := db.Where("ipv6_address = ?", ipv6Query+maxIPStr).First(&availableIPStruct) //find if any IP has been generated
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("Info - IP Address not found, generating")
			pointer := 3
			octet := 0 //so that we start at 10.6.x.3

			for j := 3; j < maxIPInt+1; j++ { //for loop to iterate and create IP in db
				if pointer >= 245 {
					pointer = 3
					octet++
				}
				ipv4 := ipv4Query + strconv.Itoa(octet) + "." + strconv.Itoa(pointer)
				ipv6 := ipv6Query + strconv.Itoa(j)
				currentIP := IP{IPv4Address: ipv4, IPv6Address: ipv6, InUse: "false", WGInterface: "wg0"}
				db.Create(currentIP)
				pointer++
			}
			log.Println("Info - Generated IPs successfully")
		}
	} else {
		modulus := maxIPInt % 242

		division := maxIPInt / 242
		division = int(math.Floor(float64(division)))

		thirdOctetInt, _ := strconv.Atoi(ipv4Splice[2])
		thirdOctetInt = thirdOctetInt + division - 1 //add octet to calculated octet and subtract one
		thirdOctetStr := strconv.Itoa(thirdOctetInt) //convert to string

		fourthOctetInt, _ := strconv.Atoi(ipv4Splice[3])
		fourthOctetInt = fourthOctetInt + modulus - 1  //add octet to calculated octet and subtract one
		fourthOctetStr := strconv.Itoa(fourthOctetInt) //convert to string

		ipv4Search := ipv4Query + "." + thirdOctetStr + "." + fourthOctetStr
		result := db.Where("ipv4_address = ?", ipv4Search).First(&availableIPStruct) //find if any IP has been generated
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {                         //if not found
			log.Println("Info - IP Address not found, generating")
			pointer := 3
			octet := 0 //so that we start at 10.6.x.3

			for j := 3; j < maxIPInt+1; j++ { //for loop to iterate and create IP in db
				if pointer >= 245 {
					pointer = 3
					octet++
				}
				ipv4 := ipv4Query + strconv.Itoa(octet) + "." + strconv.Itoa(pointer)
				currentIP := IP{IPv4Address: ipv4, IPv6Address: "-", InUse: "false", WGInterface: "wg0"}
				db.Create(currentIP)
				pointer++
			}
			log.Println("Info - Generated IPs successfully")
		}
	}
}
func generatePeerArray() []wgtypes.PeerConfig {
	var keyStruct []Key               //key struct
	var keyArray []wgtypes.PeerConfig //peers (clients)
	db := DBSystem

	resultKey := db.Find(&keyStruct)
	if errors.Is(resultKey.Error, gorm.ErrRecordNotFound) {
		return keyArray
	} else if resultKey.Error != nil {
		log.Println("Error - Finding keys", resultKey.Error)
	}

	for i := 0; i < len(keyStruct); i++ { //loop over all clients in db
		var ipStruct IP
		resultIP := db.Where("ipv4_address = ?", keyStruct[i].IPv4Address).First(&ipStruct)
		if errors.Is(resultIP.Error, gorm.ErrRecordNotFound) {
			log.Println("Cant find IPs", keyStruct[i].IPv4Address)
			continue //continue even on error
		} else if resultIP.Error != nil {
			log.Println("Error - Finding IPs", keyStruct[i].IPv4Address, resultKey.Error)
		}

		pubKey, pubErr := manager.ParseKey(keyStruct[i].PublicKey)
		preKey, preErr := manager.ParseKey(keyStruct[i].PresharedKey)
		if pubErr != nil || preErr != nil {
			log.Fatal("Error - Unable to parse keys on generate array")
		}

		var ipAddresses []net.IPNet
		ipv4, errIPv4 := manager.ParseIP(ipStruct.IPv4Address + "/32")
		if errIPv4 != nil {
			log.Fatal("Error - Parsing IPv4 Address", errIPv4)
		}
		ipAddresses = append(ipAddresses, *ipv4)

		if ipStruct.IPv6Address != "-" {
			ipv6, errIPv6 := manager.ParseIP(ipStruct.IPv6Address + "/128")
			if errIPv6 != nil {
				log.Fatal("Error - Parsing IPv6 Address", errIPv6)
			}
			ipAddresses = append(ipAddresses, *ipv6)
		}

		var zeroTime time.Duration
		userConfig := wgtypes.PeerConfig{
			PublicKey:                   pubKey,
			PresharedKey:                &preKey,
			PersistentKeepaliveInterval: &zeroTime,
			AllowedIPs:                  ipAddresses,
		}
		keyArray = append(keyArray, userConfig) //add config to client array
	}
	return keyArray
}
