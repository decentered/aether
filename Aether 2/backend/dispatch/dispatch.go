// Backend > Dispatch
// This file is the subsystem that decides on which remotes to connect to.

package dispatch

import (
	"aether-core/io/api"
	// "aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	// "aether-core/services/safesleep"
	// "errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"strings"
	// tb "aether-core/services/toolbox"
	// "net"
	"time"
)

/*
Dispatcher is the big thing here.
One thing to keep thinking about, this behaviour of the dispatch to get one online node that is not excluded, might actually create 'islands' that only connect to each other.
To be able to diagnose this, I might need to build a tool that visualises the connections between the nodes.. Just to make sure that there are no islands.
*/

const (
	maxScoutRetriesIfRemoteTooBusy = 4
	// ^ This is how many times scout will try another node if the node found is too busy. If zero no retries are made, but the main attempt will still go through.
)

// NeighbourWatch keeps in sync with all our neighbours.
func NeighbourWatch() {
	logging.Log(2, "NeighbourWatch triggers.")
	loc, subloc, port := globals.BackendTransientConfig.NeighboursList.Pop()
	a := api.Address{
		Location:    api.Location(loc),
		Sublocation: api.Location(subloc),
		Port:        port,
	}
	// The pop was blank
	if isBlank(a) {
		Scout(nil)
		return
	}
	// The pop was not blank...
	err := connect(a)
	// ...but it failed
	if err != nil {
		// if it failed because the remote was too busy,
		if strings.Contains(err.Error(), "Received status code: 429") {
			// add the failed node as exclusion, and attempt to find another.
			excl := []api.Address{a}
			Scout(&excl)
			return
		}
		// if it failed because for any other reason,
		logging.Logf(1, "NeighbourWatch: Connect failed. Error: %#v", err)
		// log and bail.
		return
	}
	// our connect succeeded, all is well. return.
	return
}

/*
Scout finds nodes that are online and that we have not connected to, and connects to them - thus adding them to the neighbourhood.

If we receive an error that is not 'too busy', we stop and wait for the next cycle. If we get a too busy error, we try to find a new node to connect to up to maxNeighbourWatchRetriesIfRemoteTooBusy times.

Rationale #1 is that if you get a too busy response, it means that the response returned fairly quick. If you get something else, you might have already spent considerable time in sync, and it might be better to just let go and wait for the next pop.

Rationale #2 is that if a sync failure happened, except in the case of 'remote too busy', it has a chance that it might be because of the local machine being under too much stress, such as DB writes timing out. In most cases, waiting a bit to reduce the pressure might help. In any case, breathlessly retrying again and again will definitely won't.

So we retry when we 100% know it's because of the remote (remote too busy case), but we let goo and let the next tick attempt a new sync. A tick is usually around a minute, so it's a good way to just wait a bit, give a little breathing room to the node, and start fresh.
*/
func Scout(excl *[]api.Address) error {
	logging.Log(2, "Scout triggers.")
	addrs := GetUnconnAddr(1, excl)
	if len(addrs) == 0 {
		logging.Log(2, "Scout got no unconnected addresses. Bailing.")
		return errors.New("Scout got no unconnected addresses. Bailing.")
	}
	err := connect(addrs[0])
	// if connect succeeded, all is well.
	if err == nil {
		return nil
	}
	// if connect failed for a reason that *wasn't* remote is too busy, bail.
	if !strings.Contains(err.Error(), "Received status code: 429") {
		logging.Logf(1, "Scout: Connect failed. Error: %#v", err)
		return errors.New(fmt.Sprintf("Scout: Connect failed. Error: %#v", err))
	}
	// if connect failed and it failed because remote was too busy,
	logging.Logf(1, "Scout: Connect failed. Attempted remote was too busy. We'll try up to %v times to sync with different remotes.", maxScoutRetriesIfRemoteTooBusy)
	// Attempt to reconnect to different remotes. maxScoutRetriesIfRemoteTooBusy times.
	// Count down frommaxScoutRetriesIfRemoteTooBusy,
	for i := maxScoutRetriesIfRemoteTooBusy - 1; i >= 0; i-- {
		logging.Logf(1, "Scout is starting a retry attempt. Attempts left: %v", i)
		// If excl wasn't given, init
		if excl == nil {
			e := []api.Address{}
			excl = &e
		}
		// Add the remote we just tried to exclusions,
		*excl = append(*excl, addrs[0])
		// Get a new address that is not one of those we tried,
		addrs = GetUnconnAddr(1, excl)
		// (And if we get nothing, bail)
		if len(addrs) == 0 {
			logging.Log(2, "Scout got no unconnected addresses in the retry process. Bailing.")
			return errors.New("Scout got no unconnected addresses in the retry process. Bailing.")
		}
		// and attempt to connect to the address.
		err := connect(addrs[0])
		// If succeeds, all is well.
		if err == nil {
			return nil
		}
		// if connect failed for a reason that *wasn't* remote is too busy, bail.
		if !strings.Contains(err.Error(), "Received status code: 429") {
			logging.Logf(1, "Scout: Connect failed. Error: %#v", err)
			return errors.New(fmt.Sprintf("Scout: Connect failed. Error: %#v", err))
		}
		// If it failed because of a 'remote too busy, log the failure,
		logging.Logf(1, "We encountered an error in Scout retry attempt cycle. Retry attempts left: %v, Error: %v", i, err)
		// and add the address we just tried to exclusions, so that we won't try to pick it again in the next retry attempt. The next iteration of the loop will retry with a different remote.
		*excl = append(*excl, addrs[0])
	}
	// If we've come here without returning, scout tried its best, but could not connect to anything. All errors we got were 'remote too busy'.
	allNodesBusyError := errors.New(fmt.Sprintf("Scout failed because all nodes we have tried responded with 'too busy'. We tried %v different nodes.", maxScoutRetriesIfRemoteTooBusy+1))
	logging.Logf(1, "Scout: Connect failed. Error: %#v", allNodesBusyError)
	return errors.New(fmt.Sprintf("Scout: Connect failed. Error: %#v", allNodesBusyError))
}

// InboundConnectionWatch takes a look at how many inbound connections we have received in the past 3 minutes. If the number is zero, it triggers a reverse connection open request to a node.
func InboundConnectionWatch() {
	nt := globals.BackendConfig.GetNodeType()
	if nt != 2 {
		// If not a live node, we don't request reverse opens.
		return
	}
	logging.Log(2, "Inbound connection watch triggers.")
	pastInboundConns := globals.BackendTransientConfig.Bouncer.GetInboundsInLastXMinutes(3)
	if len(pastInboundConns) < 2 || globals.BackendTransientConfig.NewContentCommitted {
		// Request reverse connect to a node we think is online.
		online, err := findOnlineNodes(1, -1, nil)
		if err != nil {
			logging.Logf(1, "Find online nodes for InboundConnectionWatch failed. Error: %v", err)
		}
		if len(online) > 0 {
			api.RequestInboundSync(string(online[0].Location),
				string(online[0].Sublocation),
				online[0].Port)
		}
	}
	globals.BackendTransientConfig.NewContentCommitted = false
}

/*
//////////
Internal functions
//////////
*/

func isBlank(a api.Address) bool {
	return len(a.Location) == 0 &&
		len(a.Sublocation) == 0 &&
		a.Port == 0
}

func connect(a api.Address) error {
	// sync
	err := Sync(a, []string{}, nil)
	if err != nil {
		logging.Log(2, fmt.Sprintf("Sync failed. Address: %#v, Error: %#v", a, err))
		return errors.Wrapf(err, "Sync failed. Address: %#v", a)
	}
	now := time.Now()
	// Add to exclusions for a while
	addrIface := interface{}(a)
	globals.BackendTransientConfig.DispatcherExclusions[&addrIface] = now
	globals.BackendTransientConfig.NeighboursList.Push(string(a.Location), string(a.Sublocation), a.Port)
	return nil
}

// sameAddress checks if the addresses given are the same
func sameAddress(a1 *api.Address, a2 *api.Address) bool {
	if a1.Location == a2.Location && a1.Sublocation == a2.Sublocation && a1.Port == a2.Port {
		return true
	}
	return false
}

// addrsInGivenSlice checks if the address is extant in a given slice.
func addrsInGivenSlice(addr *api.Address, slc *[]api.Address) bool {
	address := *addr
	slice := *slc
	for i, _ := range slice {
		if sameAddress(&address, &slice[i]) {
			return true
		}
	}
	return false
}

// eliminateExcludedAddressesFromList returns a clean address list that is devoid of any address in the exclusions list. It also checks whether the address is a given type.
func eliminateExcludedAddressesFromList(addrs *[]api.Address, excls *[]api.Address, addressType uint8) []api.Address {
	addresses := *addrs
	exclusions := *excls
	var cleanList []api.Address
	for i, _ := range addresses {
		if !addrsInGivenSlice(&addresses[i], &exclusions) && addresses[i].Type == addressType {
			// If address is not in the exclusions list and in the type we want
			cleanList = append(cleanList, addresses[i])
		}
	}
	return cleanList
}

// processExclusions processes the exclusions in the Dispatcher, and it returns a slice of Addresses. It can also differentiate between different types of addresses (static vs live) and apply different exclusion expiry durations.
func processExclusions(excl *map[*interface{}]time.Time) []api.Address {
	liveExpiry := globals.BackendConfig.GetDispatchExclusionExpiryForLiveAddress()
	staticExpiry := globals.BackendConfig.GetDispatchExclusionExpiryForStaticAddress()
	exclusionsList := *excl
	excludedAddressesToReturn := []api.Address{}
	newExclusionsList := make(map[*interface{}]time.Time)
	for untypedKeyPt, value := range exclusionsList {
		untypedKey := *untypedKeyPt
		switch typedKey := untypedKey.(type) {
		case api.Address:
			if typedKey.Type == 255 { // Static node
				if time.Since(value) < staticExpiry {
					/*
						If the time that has passed is less than expiry,
						a) add it to the current exclusions list to be returned
						b) Add it to the new exclusions list (based on interface{}) to be set back to the original location for future cycles.
					*/
					excludedAddressesToReturn = append(excludedAddressesToReturn, typedKey)
					newExclusionsList[untypedKeyPt] = value // value is a time.Time object.
				}
			} else { // Not static node (Probably live)
				if time.Since(value) < liveExpiry {
					/*
						If the time that has passed is less than expiry,
						a) add it to the current exclusions list to be returned
						b) Add it to the new exclusions list (based on interface{}) to be set back to the original location for future cycles.
					*/
					excludedAddressesToReturn = append(excludedAddressesToReturn, typedKey)
					newExclusionsList[untypedKeyPt] = value // value is a time.Time object.
				}
			}

		default:
			// Basically, if this happens, it will be guaranteed to not persist as it won't go through this sieve.
			logging.Log(1, fmt.Sprintf("processExclusions encountered an object in the exclusions map that was not an address. Object: %#v", untypedKey))
		}
	}
	// After processing, set the updated exclusions list back to its original location, and return the created processed list.
	excl = &newExclusionsList
	return excludedAddressesToReturn
}
