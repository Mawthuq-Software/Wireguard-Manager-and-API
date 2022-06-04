package autocheck

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/logger"
)

func AutoStart() {
	combinedLogger := logger.GetCombinedLogger()
	autocheckBool := viper.GetBool("SERVER.AUTOCHECK")

	if autocheckBool {
		combinedLogger.Info("Setting up wireguard network interface")
		combinedLogger.Info("AutoStart running")

		s := gocron.NewScheduler(time.UTC)
		s.Every(5).Minute().Do(checkPeer)
		s.Every(1).Minute().Do(checkBW)
		s.StartAsync()

	} else {
		combinedLogger.Info("Skipping autochecker as SERVER.AUTOCHECK is set to false")
	}
}

var checkPeer = func() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Running check wg check peers")

	for i := 0; i < 2; i++ {
		boolWG := db.AddRemovePeers()
		if !boolWG {
			combinedLogger.Error("When AddRemovePeer was run")

		} else {
			combinedLogger.Info("Successfully ran AddRemovePeer")
			break
		}
	}
}

var checkBW = func() {
	combinedLogger := logger.GetCombinedLogger()
	combinedLogger.Info("Running check bandwidth")

	db.BWPeerCheck()
}
