package db

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	errMigrate := db.AutoMigrate(&Key{}, &IP{}, &WireguardInterface{}, &Subscription{}) //Migrate tables to sqlite
	if errMigrate != nil {
		log.Fatal("Error - Migrating database", errMigrate)
	} else {
		log.Println("Info - Successfully migrated db")
	}
	generateIPs()
}

func generateIPs() {
	log.Println("Info - Generating IPs")

	var availableIPStruct IP
	db := DBSystem

	maxIPStr := viper.GetString("SERVER.MAX_IP") //get maximum IPs in config.json
	maxIPInt, err := strconv.Atoi(maxIPStr)      //convert to int

	if err != nil {
		log.Fatal("Unable to convert IP to int", err)
	}

	ipv4Addr := viper.GetString("INSTANCE.IP.LOCAL.IPV4.ADDRESS")  //IPv4 Subnet Address
	ipv4Splice := strings.SplitAfter(ipv4Addr, ".")                //split str at decimal
	ipv4Query := ipv4Splice[0] + ipv4Splice[1]                     //get first two subnet
	ipv6Addr := viper.GetString("INSTANCE.IP.LOCAL.IPv6.ADDRESS")  //IPv6 Subnet Address
	ipv6Enabled := viper.GetBool("INSTANCE.IP.LOCAL.IPV6.ENABLED") //Check if IPV6 is enabled

	if ipv6Enabled {
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
