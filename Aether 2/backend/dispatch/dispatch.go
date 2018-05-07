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
	// "strings"
	// tb "aether-core/services/toolbox"
	// "net"
	"time"
)

/*
Dispatcher is the big thing here.
One thing to keep thinking about, this behaviour of the dispatch to get one online node that is not excluded, might actually create 'islands' that only connect to each other.
To be able to diagnose this, I might need to build a tool that visualises the connections between the nodes.. Just to make sure that there are no islands.
*/

// // Dispatcher is the loop that controls the outbound connections.
// func Dispatcher(addressType uint8) {
// 	logging.Log(1, fmt.Sprintf("Dispatch for type %d has started.", addressType))
// 	defer logging.Log(1, fmt.Sprintf("Dispatch for type %d has exited.", addressType))
// 	// Set up the mutexes.
// 	if addressType == 2 {
// 		globals.BackendTransientConfig.LiveDispatchRunning = true
// 	} else if addressType == 255 {
// 		globals.BackendTransientConfig.StaticDispatchRunning = true
// 	}
// 	//	Check the exclusions list and clean out the expired exclusions.
// 	exclSlice := processExclusions(&globals.BackendTransientConfig.DispatcherExclusions)
// 	//	Ask for one online node.
// 	onlineAddresses := []api.Address{}
// 	err := errors.New("")
// 	onlineAddresses, err = GetOnlineAddresses(1, exclSlice, addressType, false)
// 	if err != nil {
// 		logging.Log(1, err)
// 	}
// 	if len(onlineAddresses) == 0 {
// 		logging.Log(1, fmt.Sprintf("We've found no addresses online that we've connected before. We'll now check addresses that we haven't connected before."))
// 		onlineAddresses, err = GetOnlineAddresses(1, exclSlice, addressType, true)
// 		if err != nil {
// 			logging.Log(1, err)
// 		}
// 	}
// 	// At this point, we've both checked prior-connected and non-prior connected addresses.
// 	if len(onlineAddresses) > 0 {
// 		// If there are any online addresses, connect to the first one.
// 		err2 := Sync(onlineAddresses[0])
// 		if err2 != nil {
// 			logging.Log(1, fmt.Sprintf("Sync call from Dispatcher failed. Address: %#v, Error: %#v", onlineAddresses[0], err2))
// 		}
// 		// After the sync is complete, add it to the exclusions list.
// 		addrsAsIface := interface{}(onlineAddresses[0])
// 		globals.BackendTransientConfig.DispatcherExclusions[&addrsAsIface] = time.Now()
// 		// Set the last live / static node connection timestamps.
// 		now := time.Now()
// 		if addressType == 2 {
// 			globals.BackendConfig.SetLastLiveAddressConnectionTimestamp(now.Unix())
// 		} else if addressType == 255 {
// 			globals.BackendConfig.SetLastStaticAddressConnectionTimestamp(now.Unix())
// 		}
// 	} else {
// 		// If we have no nodes that we have connected prior,
// 		logging.Log(1, "Dispatcher could not find any online addresses.")
// 	}
// 	/*
// 		Clear the mutexes.
// 	*/
// 	if addressType == 2 {
// 		globals.BackendTransientConfig.LiveDispatchRunning = false
// 	} else if addressType == 255 {
// 		globals.BackendTransientConfig.StaticDispatchRunning = false
// 	}
// }

// NeighbourWatch keeps in sync with all our neighbours.
func NeighbourWatch() {
	logging.Log(2, "NeighbourWatch triggers.")
	loc, subloc, port := globals.BackendTransientConfig.NeighboursList.Pop()
	a := api.Address{
		Location:    api.Location(loc),
		Sublocation: api.Location(subloc),
		Port:        port,
	}
	if !isBlank(a) {
		err := connect(a)
		if err != nil {
			logging.Logf(1, "NeighbourWatch: Connect failed. Error: %#v", err)
		}
	} else {
		logging.Log(2, "NeighbourWatch ended up with an empty pop. Triggering Scout.")
		Scout()
	}
}

// Scout finds nodes that are online and that we have not connected to, and connects to them - thus adding them to the neighbourhood.
func Scout() {
	logging.Log(2, "Scout triggers.")
	addrs := GetUnconnAddr(1)
	if len(addrs) == 0 {
		// fmt.Println("Scout got no unconnected addresses, bailing.")
		logging.Log(2, "Scout got no unconnected addresses. Bailing.")
		return
	}
	err := connect(addrs[0])
	if err != nil {
		logging.Logf(1, "Scout: Connect failed. Error: %#v", err)
		return
	}
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
	// Set mutex
	globals.BackendTransientConfig.ActiveOutbound.Lock()
	defer globals.BackendTransientConfig.ActiveOutbound.Unlock()
	// sync
	err := Sync(a)
	if err != nil {
		logging.Log(1, fmt.Sprintf("Sync failed. Address: %#v, Error: %#v", a, err))
		return errors.Wrapf(err, "Sync failed. Address: %#v", a)
	}
	now := time.Now()
	// Add to exclusions for a while
	addrIface := interface{}(a)
	globals.BackendTransientConfig.DispatcherExclusions[&addrIface] = now
	// Mark
	switch a.Type {
	case 2:
		globals.BackendConfig.SetLastLiveAddressConnectionTimestamp(now.Unix())
	case 3:
		//globals.BackendConfig.SetLastLiveAddressConnectionTimestamp(now.Unix())
	case 255:
		globals.BackendConfig.SetLastStaticAddressConnectionTimestamp(now.Unix())
	}
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
