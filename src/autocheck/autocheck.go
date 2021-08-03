package autocheck

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/manager"
)

func AutoStart() {
	log.Println("Info - AutoStart running")
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Minutes().Do(checkPeer)
	s.StartAsync()
}

var checkPeer = func() {
	log.Println("Info - Running check wg check peers")
	for i := 0; i < 2; i++ {
		boolWG := manager.AddRemovePeer()
		if !boolWG {
			log.Println("Error - When AddRemovePeer was run")

		} else {
			log.Println("Info - Successfully ran AddRemovePeer")
			break
		}
	}
}
