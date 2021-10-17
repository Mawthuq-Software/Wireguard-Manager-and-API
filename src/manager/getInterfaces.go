package manager

import (
	"errors"
	"log"

	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GetInterfaces() ([]*wgtypes.Device, error) { //get interfaces
	client, errInstance := createInstance()
	if errInstance != nil {
		log.Println("Create instance", errInstance)
		return nil, errors.New(errInstance.Error())
	}
	getInterfaces, err := client.Devices()
	if !logger.ErrorHandler("Info - Finding interfaces", err) {
		return nil, errors.New(err.Error())
	}

	closeInstance(client) //release resources and close instance
	return getInterfaces, nil
}
