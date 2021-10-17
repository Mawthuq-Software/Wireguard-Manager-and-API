package autocheck

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gitlab.com/raspberry.tech/wireguard-manager-and-api/src/db"
)

func AutoStart() {
	log.Println("Info - AutoStart running")
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Minute().Do(checkPeer)
	s.Every(1).Minute().Do(checkBW)
	s.StartAsync()
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
