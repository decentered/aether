// Backend > Dispatch
// This file is the subsystem that decides on which remotes to connect to.

package dispatch

import (
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"fmt"
	// "strings"
	// "errors"
	"time"
)

// processExclusions processes the exclusions in the Dispatcher, and it returns a slice of Addresses. It can also differentiate between different types of addresses (static vs live) and apply different exclusion expiry durations.
func processExclusions(excl *map[*interface{}]time.Time) []api.Address {
	liveExpiry := globals.DispatcherExclusionsExpiryLiveAddress
	staticExpiry := globals.DispatcherExclusionsExpiryStaticAddress
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

/*
Dispatcher is the big thing here.
One thing to keep thinking about, this behaviour of the dispatch to get one online node that is not excluded, might actually create 'islands' that only connect to each other.
To be able to diagnose this, I might need to build a tool that visualises the connections between the nodes.. Just to make sure that there are no islands.
*/

// Dispatcher is the loop that controls the outbound connections.
func Dispatcher(addressType uint8) {
	logging.Log(1, fmt.Sprintf("Dispatch for AddressType: %d has started.", addressType))
	defer logging.Log(1, fmt.Sprintf("Dispatch for AddressType: %d is complete.", addressType))
	/*
		Set up the mutexes.
	*/
	if addressType == 2 {
		globals.LiveDispatchRunning = true
	} else if addressType == 255 {
		globals.StaticDispatchRunning = true
	}
	/*
		Check the exclusions list and clean out the expired exclusions.
	*/
	exclSlice := processExclusions(&globals.DispatcherExclusions)
	/*
		Ask for one online node.
	*/
	onlineAddresses, err := GetOnlineAddresses(1, exclSlice, addressType)
	if err != nil {
		logging.Log(1, err)
	}
	if len(onlineAddresses) > 0 {
		/*
			If there are any online addresses, connect to the first one.
		*/
		err2 := Sync(onlineAddresses[0])
		if err2 != nil {
			logging.Log(1, fmt.Sprintf("Sync call from Dispatcher failed. Address: %#v, Error: %#v", onlineAddresses[0], err2))
		}
		/*
			After the sync is complete, add it to the exclusions list.
		*/
		addrsAsIface := interface{}(onlineAddresses[0])
		globals.DispatcherExclusions[&addrsAsIface] = time.Now()
		// Set the last live / static node connection timestamps.
		now := time.Now()
		if addressType == 2 {
			globals.LastLiveNodeConnectionTs = now.Unix()
		} else if addressType == 255 {
			globals.LastStaticNodeConnectionTs = now.Unix()
		}
	} else {
		logging.Log(1, "Dispatcher could not find any online addresses. It will a)trigger the AddressScanner so it can convert more addresses to known addresses, rendering them eligible to be used by Dispatcher in the next iteration, and b) Connect to a static node if a static node hasn't been connected in the last 10 minutes.")
		if addressType == 2 {
			// Specifically for the live node check, if there are no active nodes found, check if the last static check was more than ten minutes ago. If so, run a static update before filing a request with tha AddressScanner.
			TenMinAgo := int64(time.Now().Add(-10 * time.Minute).Unix())
			if globals.LastStaticNodeConnectionTs < TenMinAgo {
				logging.Log(1, "Live node dispatcher is firing a static node dispatcher because it's been more than ten minutes after a connection to a static node.")
				if !globals.StaticDispatchRunning {
					Dispatcher(255)
				} else {
					logging.Log(1, "Live node dispatcher attempted to fire a static dispatcher because it found no live nodes, but it was preempted by an existing static dispatcher currently running.")
				}
			}
		}
		AddressScanner()
	}
	/*
		Clear the mutexes.
	*/
	if addressType == 2 {
		globals.LiveDispatchRunning = false
	} else if addressType == 255 {
		globals.StaticDispatchRunning = false
	}
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

// eliminateExcludedAddressesFromList returns a clean address list that is devoid of any address in the exclusions list.
func eliminateExcludedAddressesFromList(addrs *[]api.Address, excls *[]api.Address) []api.Address {
	addresses := *addrs
	exclusions := *excls
	var cleanList []api.Address
	for i, _ := range addresses {
		if !addrsInGivenSlice(&addresses[i], &exclusions) {
			// If address is not in the exclusions list
			cleanList = append(cleanList, addresses[i])
		}
	}
	return cleanList
}

// TODO: We need tests for this. Kind of hard to mock as it requires actually online nodes. But it does have a few things that can end up a bit hairy.
// GetOnlineAddresses goes through the addresses database and finds the requested amount of live nodes and provides it back. It provides a useful feature with exclusions, in that you can provide a list of addresses that you want to exclude (perhaps, addresses you connected recently, and that you don't want to connect for a while).
func GetOnlineAddresses(noOfOnlineAddressesRequested int, exclude []api.Address, addressType uint8) ([]api.Address, error) {
	logging.Log(1, fmt.Sprintf("SEEK START for %d online addresses with type %d in the DB with %d addresses excluded.", noOfOnlineAddressesRequested, addressType, len(exclude)))
	var onlineAddresses []api.Address
	PAGESIZE := globals.OnlineAddressFinderPageSize
	offset := 0
	// Until the number of online addresses found is equal to or more than addresses requested,
	for len(onlineAddresses) < noOfOnlineAddressesRequested {
		// Read addresses from the database,
		resp, err := persistence.ReadAddresses(
			"", "", 0, 0, 0, PAGESIZE, offset, addressType)
		if err != nil {
			return []api.Address{}, err
		}
		if len(resp) == 0 {
			// We ran out of addresses in the database.
			logging.Log(1, "We ran out of items in the database while trying to do GetOnlineAddresses.")
			return onlineAddresses, nil
		}
		// Put the read addresses into Pinger to extract the live addresses,
		updatedAddresses := Pinger(resp)
		// (And commit the newly found addressed to DB, just in case)
		persistence.InsertOrUpdateAddresses(&updatedAddresses)
		// Check for the exclusions, so that the address we have isn't what we want to exclude.
		cleanedUpdatedAddresses := eliminateExcludedAddressesFromList(&updatedAddresses, &exclude)
		// Add the found online addresses to the result set,
		onlineAddresses = append(onlineAddresses, cleanedUpdatedAddresses...)
		// Set the offset by the page size, so you get the next 'page' from the database
		offset = offset + PAGESIZE
		logging.Log(2, fmt.Sprintf("Number of online addresses in this GetOnlineAddress page: %d", len(onlineAddresses)))
	}
	logging.Log(1, fmt.Sprintf("SEEK END for %d online addresses in the DB. We found %d online addresses.", noOfOnlineAddressesRequested, len(onlineAddresses)))
	// If we arrived here, either we ended up with enough (or more than enough nodes, or we ran out of addresses to check in the DB.)
	return onlineAddresses, nil
}

func unconnectedAddressSearch(days int) ([]api.Address, error) {
	now := time.Now()
	pastTs := api.Timestamp(now.AddDate(0, 0, -days).Unix())
	// Get me all addresses that was inputted up to two weeks ago
	resp, err := persistence.ReadAddresses(
		"", "", 0, pastTs, 0, 0, 0, 0)
	// resp, err := persistence.ReadAddresses(
	// 	"", "", 0, 0, 0, 1000,0)
	if err != nil {
		return []api.Address{}, err
	}
	return resp, nil
}

// AddressScanner goes through all the **prior-unconnected** addresses that were provided by other remotes up to two weeks ago, and it goes through them. If it finds any online nodes that are able to connect, it will mark them as such, and set the appropriate address type, rendering them **known**. Setting the node type renders the address eligible to be connected via dispatch. This method will be called by Dispatch if it ends up finding no nodes to connect to, and in 6-hour intervals.
func AddressScanner() error {
	if globals.AddressesScannerActive {
		logging.Log(1, "AddressScanner is already running right now. Skipping this call. (This happens when a Dispatch runs out of items and calls AddressScanner on its own while it's already running on a scheduled time)")
		return nil
	}
	globals.AddressesScannerActive = true
	defer func() { globals.AddressesScannerActive = false }()
	logging.Log(1, "SEEK START for prior-unconnected addresses.")
	defer logging.Log(1, "SEEK END for prior-unconnected addresses.")
	resp, err := unconnectedAddressSearch(14)
	if err != nil {
		return err
	}
	if len(resp) == 0 {
		logging.Log(1, "AddressesScanner could not find any prior-unconnected addresses in the last 14 days. Expanding search to last 30 days.")
		resp2, err2 := unconnectedAddressSearch(30)
		if err2 != nil {
			return err2
		}
		if len(resp2) == 0 {
			logging.Log(1, "AddressesScanner could not find any prior-unconnected addresses in the last 30 days. Expanding search to last 90 days.")
			resp3, err3 := unconnectedAddressSearch(90)
			if err3 != nil {
				return err3
			}
			if len(resp3) == 0 {
				logging.Log(1, "AddressesScanner could not find any prior-unconnected addresses in the last 90 days. Expanding search to all past addresses captured at any time.")
				resp4, err4 := persistence.ReadAddresses(
					"", "", 0, 0, 0, 0, 0, 0)
				if err4 != nil {
					return err4
				}
				// Assign to a scope out so it'll trickle down to the resp eventually.
				resp3 = resp4
			}
			// Assign to a scope out so it'll trickle down to the resp eventually.
			resp2 = resp3
		}
		// Assign to a scope out so it'll trickle down to the resp eventually.
		resp = resp2
	}
	logging.Log(1, fmt.Sprintf("We have this many prior-unconnected addresses to check: %d", len(resp)))
	updatedAddresses := Pinger(resp)
	persistence.InsertOrUpdateAddresses(&updatedAddresses)
	return nil
}

// Pinger goes through the list of extant nodes, and pings the status endpoints to see if they are online. If no response is provided in X seconds, the node is offline. It returns a set of online nodes.
// We need to do this in batches of 100. Otherwise we end up with "socket: too many open files" error.
func Pinger(fullAddressesSlice []api.Address) []api.Address {
	// Paginate addresses first. We batch these into pages of 100, because it's very easy to run into too many open files error if you just dump it through.
	var pages [][]api.Address
	dataSet := fullAddressesSlice
	PAGESIZE := globals.PingerPageSize
	numPages := len(dataSet)/PAGESIZE + 1
	var allUpdatedAddresses []api.Address
	// The division above is floored.
	for i := 0; i < numPages; i++ {
		beg := i * PAGESIZE
		var end int
		// This is to protect from 'slice bounds out of range'
		if (i+1)*PAGESIZE > len(dataSet) {
			end = len(dataSet)
		} else {
			end = (i + 1) * PAGESIZE
		}
		pageData := dataSet[beg:end]
		var page []api.Address
		page = pageData
		pages = append(pages, page)
	}
	// For every page,
	for i, _ := range pages {
		// Run the core logic.
		addrs := pages[i]
		outputChan := make(chan api.Address)
		for j, _ := range addrs {
			logging.Log(2, fmt.Sprintf("Pinging the address at %#s:%d", addrs[j].Location, addrs[j].Port))
			go Ping(addrs[j], outputChan)
		}
		var updatedAddresses []api.Address
		// We will receive as many addresses as answers. Every time something is put into a channel, this will fire, if the channel is empty, it will block.
		for i := 0; i < len(addrs); i++ {
			var a api.Address
			a = <-outputChan
			updatedAddresses = append(updatedAddresses, a)
		}
		allUpdatedAddresses = append(allUpdatedAddresses, updatedAddresses...)
	}
	// Clean blanks.
	logging.Log(2, fmt.Sprintf("All updated addresses count (this should be the same as goroutine count: %d", len(allUpdatedAddresses)))
	var cleanedAllUpdatedAddresses []api.Address
	for i, _ := range allUpdatedAddresses {
		if allUpdatedAddresses[i].Location != "" {
			// The location is not blank. This is an actual updated address.
			cleanedAllUpdatedAddresses = append(cleanedAllUpdatedAddresses, allUpdatedAddresses[i])
		}
	}
	logging.Log(2, fmt.Sprintf("Cleaned addresses count (this should be the same as the online addresses count: %d", len(cleanedAllUpdatedAddresses)))

	return cleanedAllUpdatedAddresses
}

// Ping runs a Check and returns the result. If there is an error, it returns a blank address.
func Ping(addr api.Address, processedAddresses chan<- api.Address) {
	logging.Log(2, fmt.Sprintf("Connection attempt started: %v:%v", addr.Location, addr.Port))
	var blankAddr api.Address
	updatedAddr, _, _, err := Check(addr)
	if err != nil {
		updatedAddr = blankAddr
		logging.Log(2, err)
	}
	processedAddresses <- updatedAddr
}
