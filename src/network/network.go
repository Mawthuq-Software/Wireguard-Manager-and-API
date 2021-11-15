package network

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/vishvananda/netlink"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

func SetupWG() {
	log.Println("Info - Setting up WG interface")
	db.WGStart()
	wg0, errLink := netlink.LinkByName("wg0")
	if errLink != nil {
		log.Fatal("Error - Failed to get link to wireguard interface")
	}
	ipCheck(wg0)
}

func addIP(wg0 netlink.Link, ipAddr *netlink.Addr) {
	ipAddErr := netlink.AddrAdd(wg0, ipAddr)
	if ipAddErr != nil {
		fmt.Println("Error - Failed to add IP address ", ipAddErr)
		log.Fatal("Error - Failed to add IP address ", ipAddErr)
	}
	log.Println("Info - Added IP address to interface")
}

func ipCheck(wg0 netlink.Link) {
	log.Println("Info - Checking if IPs exist")

	IPs, err := netlink.AddrList(wg0, 0) //list of IP addresses in system, equivalent to: `ip addr show`
	if err != nil {
		fmt.Println("Error - Failed to get find wireguard interface")
		log.Fatal("Error - Failed to get find wireguard interface")
	}

	ipv4Check := false //variables for checks
	ipv6Check := false

	wgIPv4 := viper.GetString("INSTANCES.wg0.IP.LOCAL.IPV4.ADDRESS") //IPv4 in config
	ipv4Subnet := viper.GetString("INSTANCES.wg0.IP.LOCAL.IPV4.SUBNET")
	wgIPv6 := viper.GetString("INSTANCES.wg0.IP.LOCAL.IPV6.ADDRESS") //IPv6 in config

	ipv4Addr, errParsev4 := netlink.ParseAddr(wgIPv4 + ipv4Subnet) //add subnet of 16 to IP
	if errParsev4 != nil {
		log.Fatal("Error - Failed to get parse IPv4 Address")
	}

	if wgIPv6 != "-" { //if IPv6 is not set to - in config
		ipv6Subnet := viper.GetString("INSTANCES.wg0.IP.LOCAL.IPV6.SUBNET")
		ipv6Addr, errParsev6 := netlink.ParseAddr(wgIPv6 + ipv6Subnet)
		if errParsev6 != nil {
			log.Fatal("Error - Failed to get parse IPv6 Address")
		}
		for i := 0; i < len(IPs); i++ { //checks if IPs wanted exist
			if IPs[i].Equal(*ipv4Addr) { //Check if IPv4 address wanted is already present
				ipv4Check = true
			} else if IPs[i].Equal(*ipv6Addr) { //Check if IPv6 address wanted is already present
				ipv6Check = true
			}
		}
		if !ipv6Check {
			addIP(wg0, ipv6Addr) //add IPv6 to system
		}
	} else {
		for i := 0; i < len(IPs); i++ { //checks if IPs wanted exist
			if IPs[i].Equal(*ipv4Addr) { //Check if IPv4 address wanted is already present
				ipv4Check = true
			}
		}
	}
	if !ipv4Check {
		addIP(wg0, ipv4Addr) //add IPv4 to system
	}
}
