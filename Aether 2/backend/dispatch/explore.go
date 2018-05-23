// Backend > Routines > Explore
// This file contains the explore routine in dispatch. Explore is the routine that helps us discover new nodes, and it also makes sure that we occasionally check our static nodes and bootstrappers as well.

package dispatch

import (
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
)

// Explore is a function that reaches further into the network than our neighbourhood watch. At every call, it finds a live regular remote that is new to us and does a sync, and at the end it inserts it into our neighbourhood. At every 6 ticks, it syncs with our statics, and at every 36 ticks, it syncs with bootstrappers. An explore tick is not always a minute, it depends on the schedule, but it is something around 10 minutes. So that means refreshing statics every hour, and hitting bootstrappers every 6 hours.
func Explore() {
	ticker := globals.BackendTransientConfig.ExplorerTick
	if ticker%36 == 0 && ticker != 0 {
		globals.BackendTransientConfig.ExplorerTick = 0
		// call 3 live and 1 static bootstrap nodes and sync with them.
		liveBs, err := persistence.ReadAddresses("", "", 0, 0, 0, 3, 0, 3, "limit")
		if err != nil {
			logging.Logf(1, "There was an error when we tried to read live bootstrapper addresses for Explore schedule. Error: %#v", err)
		}
		staticBs, err2 := persistence.ReadAddresses("", "", 0, 0, 0, 1, 0, 254, "limit")
		if err2 != nil {
			logging.Logf(1, "There was an error when we tried to read static bootstrapper addresses for Explore schedule. Error: %#v", err)
		}
		addrs := append(liveBs, staticBs...)
		for key, _ := range addrs {
			Sync(addrs[key], []string{})
		}
	} else if ticker%6 == 0 && ticker != 0 {
		// go through all statics to see if there are any updates.
		statics, err := persistence.ReadAddresses("", "", 0, 0, 0, 4, 0, 255, "limit") // get all static nodes we know of.
		if err != nil {
			logging.Logf(1, "There was an error when we tried to read static addresses for Explore schedule. Error: %#v", err)
		}
		for key, _ := range statics {
			Sync(statics[key], []string{})
		}
	} else {
		// find a new node that we haven't synced before, and sync with it.
		Scout()
	}
	globals.BackendTransientConfig.ExplorerTick++
}
