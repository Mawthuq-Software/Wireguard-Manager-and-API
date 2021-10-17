package manager

import (
	"log"
	"net"

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

func (attr wgInterface) Attrs() *netlink.LinkAttrs {
	return attr.Attributes
}
func (t wgInterface) Type() string {
	return t.TypeName
}
func returnNewLink(l Link) netlink.Link {
	return l
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
