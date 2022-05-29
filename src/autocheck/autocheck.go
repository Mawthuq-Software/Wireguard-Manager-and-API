package autocheck

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

func AutoStart() {
	consoleLogger := logger.GetConsoleLogger()
	autocheckBool := viper.GetBool("SERVER.AUTOCHECK")

	if autocheckBool {
		consoleLogger.Info("Setting up wireguard network interface")

		log.Println("Info - AutoStart running")
		s := gocron.NewScheduler(time.UTC)
		s.Every(5).Minute().Do(checkPeer)
		s.Every(1).Minute().Do(checkBW)
		s.StartAsync()

	} else {
		consoleLogger.Info("Skipping autochecker as SERVER.AUTOCHECK is set to false")
	}
}

var checkPeer = func() {
	log.Println("Info - Running check wg check peers")
	for i := 0; i < 2; i++ {
		boolWG := db.AddRemovePeers()
		if !boolWG {
			log.Println("Error - When AddRemovePeer was run")

		} else {
			log.Println("Info - Successfully ran AddRemovePeer")
			break
		}
	}
}

var checkBW = func() {
	log.Println("Info - Running check bandwidth")
	db.BWPeerCheck()
}
