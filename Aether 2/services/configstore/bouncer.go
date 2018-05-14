// Services > ConfigStore > Bouncer
// This module controls how many remotes can connect to this computer at the same time, and how many outbound syncs can be happening at the same time.

package configstore

import (
	// "fmt"
	"sync"
	"time"
)

// These two are local variables that only affect this specific library. Since there is no reason to modify them from the outside, this is not brought into the main settings JSON.
const (
	LeaseDurationSeconds        = 60
	MinimumFlushIntervalSeconds = 60
	// MaxInboundConns = 10 // moved to main permanent config store.
	// MaxOutboundConns            = 1
)

type Bouncer struct {
	lock      sync.Mutex
	Inbounds  []ActiveConnection
	Outbounds []ActiveConnection // as of now, not used
	LastFlush Timestamp
}

type ActiveConnection struct {
	Location    string
	Sublocation string
	Port        uint16
	LastAccess  Timestamp
}

func (n *ActiveConnection) equal(c ActiveConnection) bool {
	return n.Location == c.Location && n.Sublocation == n.Sublocation && n.Port == c.Port
}
func (n *ActiveConnection) hasActiveLease() bool {
	cutoff := Timestamp(time.Now().Add(-(time.Duration(LeaseDurationSeconds) * time.Second)).Unix())
	if n.LastAccess > cutoff {
		return true
	} else {
		return false
	}

}

func (n *Bouncer) indexOf(direction string, loc, subloc string, port uint16) int {
	switch direction {
	case "inbound":
		for key, _ := range n.Inbounds {
			if n.Inbounds[key].equal(ActiveConnection{Location: loc, Sublocation: subloc, Port: port}) {
				return key
			}
		}
		return -1
	case "outbound":
		for key, _ := range n.Outbounds {
			if n.Inbounds[key].equal(ActiveConnection{Location: loc, Sublocation: subloc, Port: port}) {
				return key
			}
		}
		return -1
	default:
		return -1
	}
}

func (n *Bouncer) insert(direction string, loc, subloc string, port uint16) {
	entry := ActiveConnection{Location: loc, Sublocation: subloc, Port: port, LastAccess: Timestamp(time.Now().Unix())}
	switch direction {
	case "inbound":
		n.Inbounds = append(n.Inbounds, entry)
	case "outbound":
		n.Inbounds = append(n.Outbounds, entry)
	}
}

func (n *Bouncer) removeItem(direction string, i int) {
	finalList := []ActiveConnection{}
	switch direction {
	case "inbound":
		finalList = append(n.Inbounds[0:i], n.Inbounds[i+1:len(n.Inbounds)]...)
		n.Inbounds = finalList
	case "outbound":
		finalList = append(n.Outbounds[0:i], n.Outbounds[i+1:len(n.Outbounds)]...)
		n.Outbounds = finalList
	}
}

func (n *Bouncer) flush() {
	// If there's been a flush in the past 10 minutes, ignore flush. This is because flush is in a hot path, we want to avoid unnecessary repeats.
	if n.LastFlush > Timestamp(time.Now().Add(-(time.Duration(MinimumFlushIntervalSeconds) * time.Second)).Unix()) {
		return
	}
	// Set lastflush to now if the gate above passes.
	n.LastFlush = Timestamp(time.Now().Add(-(time.Duration(MinimumFlushIntervalSeconds) * time.Second)).Unix())
	for i := len(n.Inbounds) - 1; i >= 0; i-- {
		if !n.Inbounds[i].hasActiveLease() {
			n.removeItem("inbound", i)
		}
	}
	for i := len(n.Outbounds) - 1; i >= 0; i-- {
		if !n.Outbounds[i].hasActiveLease() {
			n.removeItem("outbound", i)
		}
	}
}

func (n *Bouncer) RequestInboundLease(loc, subloc string, port uint16) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.flush()
	// fmt.Println("An inbound lease was requested.")
	// fmt.Printf("This is our Inbound list after flush: %#v\n", n.Inbounds)
	direction := "inbound"
	leaseIndex := n.indexOf(direction, loc, subloc, port)
	if leaseIndex != -1 && n.Inbounds[leaseIndex].hasActiveLease() {
		// fmt.Println("Lease was renewed.")
		n.Inbounds[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
		return true
	} else {
		if len(n.Inbounds) < bc.GetMaxInboundConns() {
			n.insert(direction, loc, subloc, port)
			// fmt.Println("Lease was granted.")
			// fmt.Printf("This is our Inbound list after insert: %#v\n", n.Inbounds)
			return true
		} else {
			// fmt.Println("A lease was denied.")
			return false
		}
	}
}

// Probably works but untested. We'll use it if we end up having to gate outbound connections.
// func (n *Bouncer) RequestOutboundLease(loc, subloc string, port uint16) bool {
// 	n.lock.Lock()
// 	defer n.lock.Unlock()
// 	n.flush()
// 	direction := "outbound"
// 	leaseIndex := n.indexOf(direction, loc, subloc, port)
// 	if leaseIndex != -1 && n.Outbounds[leaseIndex].hasActiveLease() {
// 		n.Outbounds[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
// 		return true
// 	} else {
// 		if len(n.Outbounds) < MaxOutboundConns {
// 			n.insert(direction, loc, subloc, port)
// 			return true
// 		} else {
// 			return false
// 		}
// 	}
// }
