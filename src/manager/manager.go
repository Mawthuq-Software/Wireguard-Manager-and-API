package manager

import (
	"net"

	"github.com/vishvananda/netlink"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"go.uber.org/zap"
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

var combinedLogger *zap.Logger = logger.GetCombinedLogger()

func createInstance() (*wgctrl.Client, error) {
	combinedLogger.Info("Creating wg device client")
	return wgctrl.New()
}
func closeInstance(client *wgctrl.Client) error {
	combinedLogger.Info("Closing wg device client")
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
