package db

type Key struct {
	KeyID        int    `gorm:"primaryKey;autoIncrement"`
	PublicKey    string `gorm:"unique"`
	PresharedKey string `gorm:"unique"`
	IPv4Address  string `gorm:"foreignKey:IPv4Address"`
	Enabled      string
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
