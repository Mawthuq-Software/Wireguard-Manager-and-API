package network

import (
	"fmt"
	"log"
	"os"

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

	wgIPv4 := os.Getenv("WG_IPV4") //IPv4 in env
	wgIPv6 := os.Getenv("WG_IPV6") //IPv6 in env

	ipv4Addr, errParsev4 := netlink.ParseAddr(wgIPv4 + "/16") //add subnet of 16 to IP
	if errParsev4 != nil {
		log.Fatal("Error - Failed to get parse IPv4 Address")
	}

	if wgIPv6 != "-" { //if IPv6 is not set to - in env
		ipv6Addr, errParsev6 := netlink.ParseAddr(os.Getenv("WG_IPV6") + "/64")
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
