package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.SetConfigType("json")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/opt/wgManagerAPI/") // path to look for the config file in

	//DEFAULTS
	viper.SetDefault("SERVER.PORT", "8443")                          //PORT
	viper.SetDefault("SERVER.SECURITY", true)                        //SECURITY
	viper.SetDefault("SERVER.AUTH", "ABCDEFG")                       //AUTHORISATION
	viper.SetDefault("INSTANCE.IP.GLOBAL.ALLOWED", "0.0.0.0/0, ::0") //ALLOWED_IPs wireguard

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic("Cannot read env file!")
	}
}
