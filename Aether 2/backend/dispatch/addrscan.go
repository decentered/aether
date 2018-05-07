// Backend > Dispatch > AddrScan
// This file is the subsystem that decides on which remotes to connect to.

package dispatch

import (
	"aether-core/io/api"
	pers "aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	// "aether-core/services/safesleep"
	// "errors"
	"fmt"
	"github.com/pkg/errors"
	// "strings"
	"aether-core/services/toolbox"
	// "net"
	"github.com/davecgh/go-spew/spew"
	"time"
)

// TODO: We need tests for this. Kind of hard to mock as it requires actually online nodes. But it does have a few things that can end up a bit hairy.
/*
GetOnlineAddresses goes through the addresses database and finds the requested amount of live nodes and provides it back. It provides a useful feature with exclusions, in that you can provide a list of addresses that you want to exclude (perhaps, addresses you connected recently, and that you don't want to connect for a while).

forceUnconnected: This will make the address attempt to find online nodes of the given type that we have not connected before. This is useful in the case that our first pre-connected check fails. It's a step below full-db connected nodes scan. It only scans for the top 300 prior-unconnecteds by localArrival.
*/
// func GetOnlineAddressesOld(
// 	noOfOnlineAddressesRequested int,
// 	exclude []api.Address,
// 	addressType uint8,
// 	forceUnconnected bool,
// ) (
// 	[]api.Address, error,
// ) {
// 	logging.Log(1, fmt.Sprintf("SEEK START for %d online addresses with type %d in the DB with %d addresses excluded. Force unconnected: %v", noOfOnlineAddressesRequested, addressType, len(exclude), forceUnconnected))
// 	var onlineAddresses []api.Address
// 	PAGESIZE := globals.BackendConfig.GetOnlineAddressFinderPageSize()
// 	offset := 0
// 	// Until the number of online addresses found is equal to or more than addresses requested,
// 	for len(onlineAddresses) < noOfOnlineAddressesRequested {
// 		// Before we do anything, if forceUnconnected is enabled, we break at the 3 pages mark.
// 		if offset >= 297 { // 3 pages.
// 			break
// 		}
// 		// Read addresses from the database,
// 		resp := []api.Address{}
// 		err := errors.New("")
// 		if !forceUnconnected {
// 			resp, err = pers.ReadAddresses(
// 				"", "", 0, 0, 0, PAGESIZE, offset, addressType, "limit") // Get only addresses with addresstype = we've connected to them in the past.
// 		} else {
// 			resp, err = pers.ReadAddresses(
// 				"", "", 0, 0, 0, PAGESIZE, offset, 0, "limit") // Get only addresses we have *NOT* connected in the past. Mind the zero in the addressType.
// 		}
// 		if err != nil {
// 			return []api.Address{}, err
// 		}
// 		if len(resp) == 0 {
// 			// We ran out of addresses in the database.
// 			logging.Log(1, fmt.Sprintf("We ran out of items in the database while trying to do GetOnlineAddresses. Force unconnected: %v", forceUnconnected))
// 			return onlineAddresses, nil
// 		}
// 		// THINK: should we change this so that it prefers the addresses that it has connected before?
// 		// Put the read addresses into Pinger to extract the live addresses,
// 		updatedAddresses := Pinger(resp)
// 		// (And commit the newly found addressed to DB, just in case.) Even if the address found is not the type we want, it is still saved as a  prior-connected (above) for future use.
// 		errs := pers.InsertOrUpdateAddresses(&updatedAddresses)
// 		if len(errs) > 0 {
// 			logging.Log(1, fmt.Sprintf("These errors were encountered on InsertOrUpdateAddress attempt: %s", errs))
// 			continue
// 		}
// 		// Check for the exclusions, so that the address we have isn't what we want to exclude. This also enforces the address type.
// 		cleanedUpdatedAddresses := eliminateExcludedAddressesFromList(&updatedAddresses, &exclude, addressType)
// 		// Add the found online addresses to the result set,
// 		onlineAddresses = append(onlineAddresses, cleanedUpdatedAddresses...)
// 		// Set the offset by the page size, so you get the next 'page' from the database
// 		offset = offset + PAGESIZE
// 		logging.Log(2, fmt.Sprintf("Number of online addresses in this GetOnlineAddress page: %d", len(onlineAddresses)))
// 	}
// 	if forceUnconnected {
// 		logging.Log(1, fmt.Sprintf("SEEK END for prior-unconnected type %d online addresses in the DB. Wanted: %d. Found %d.", addressType, noOfOnlineAddressesRequested, len(onlineAddresses)))
// 	} else {
// 		logging.Log(1, fmt.Sprintf("SEEK END for type %d online addresses in the DB. Wanted: %d. Found %d.", addressType, noOfOnlineAddressesRequested, len(onlineAddresses)))
// 	}

// 	// If we arrived here, either we ended up with enough (or more than enough nodes, or we ran out of addresses to check in the DB.)
// 	return onlineAddresses, nil
// }

// AddressScanner goes through all the **prior-unconnected** addresses that were provided by other remotes up to two weeks ago, and it goes through them. If it finds any online nodes that are able to connect, it will mark them as such, and set the appropriate address type, rendering them **known**. Setting the node type renders the address eligible to be connected via dispatch. This method will be called by Dispatch if it ends up finding no nodes to connect to, and in 6-hour intervals.
// func AddressScannerOld() error {
// 	if globals.BackendTransientConfig.AddressesScannerActive {
// 		logging.Log(1, "AddressScanner is already running right now. Skipping this call. (This happens when a Dispatch runs out of items and calls AddressScanner on its own while it's already running on a scheduled time)")
// 		return nil
// 	}
// 	globals.BackendTransientConfig.AddressesScannerActive = true
// 	defer func() { globals.BackendTransientConfig.AddressesScannerActive = false }()
// 	logging.Log(1, "SEEK START for prior-unconnected addresses.")
// 	defer logging.Log(1, "SEEK END for prior-unconnected addresses.")
// 	resp, err := unconnectedAddressSearch(14)
// 	if err != nil {
// 		return err
// 	}
// 	if len(resp) == 0 {
// 		logging.Log(1, "AddressesScanner could not find any prior-unconnected addresses in the last 14 days. Expanding search to last 30 days.")
// 		resp2, err2 := unconnectedAddressSearch(30)
// 		if err2 != nil {
// 			return err2
// 		}
// 		if len(resp2) == 0 {
// 			logging.Log(1, "AddressesScanner could not find any prior-unconnected addresses in the last 30 days. Expanding search to last 90 days.")
// 			resp3, err3 := unconnectedAddressSearch(90)
// 			if err3 != nil {
// 				return err3
// 			}
// 			if len(resp3) == 0 {
// 				logging.Log(1, "AddressesScanner could not find any prior-unconnected addresses in the last 90 days. Expanding search to all past addresses captured at any time.")
// 				resp4, err4 := pers.ReadAddresses(
// 					"", "", 0, 0, 0, 0, 0, 0, "timerange_all")
// 				if err4 != nil {
// 					return err4
// 				}
// 				// Assign to a scope out so it'll trickle down to the resp eventually.
// 				resp3 = resp4
// 			}
// 			// Assign to a scope out so it'll trickle down to the resp eventually.
// 			resp2 = resp3
// 		}
// 		// Assign to a scope out so it'll trickle down to the resp eventually.
// 		resp = resp2
// 	}
// 	logging.Log(1, fmt.Sprintf("We have this many prior-unconnected addresses to check: %d", len(resp)))
// 	updatedAddresses := Pinger(resp)
// 	errs := pers.InsertOrUpdateAddresses(&updatedAddresses)
// 	if len(errs) > 0 {
// 		logging.Log(1, fmt.Sprintf("These errors were encountered on InsertOrUpdateAddress attempt: %s", errs))
// 	}
// 	return nil
// }

/*
//////////
Internal functions
//////////
*/

// func unconnectedAddressSearch(days int) ([]api.Address, error) {
// 	now := time.Now()
// 	pastTs := api.Timestamp(now.AddDate(0, 0, -days).Unix())
// 	// Get me all addresses that was inputted up to two weeks ago
// 	resp, err := pers.ReadAddresses(
// 		"", "", 0, pastTs, 0, 0, 0, 0, "timerange_all")
// 	if err != nil {
// 		return []api.Address{}, err
// 	}
// 	return resp, nil
// }

func getAllAddresses(isDesc bool) (*[]api.Address, error) {
	searchType := ""
	if isDesc {
		searchType = "all_desc"
	} else {
		searchType = "all_asc"
	}
	resp, err := pers.ReadAddresses("", "", 0, 0, 0, 0, 0, 0, searchType)
	if err != nil {
		errors.Wrap(err, "getAllAddresses in AddressScanner failed.")
	}
	return &resp, nil
}

func filterByLastSuccessfulPing(addrs *[]api.Address, cutoff api.Timestamp) *[]api.Address {
	live := []api.Address{}
	for key, _ := range *addrs {
		if (*addrs)[key].LastSuccessfulPing >= cutoff {
			live = append(live, (*addrs)[key])
		}
	}
	return &live
}

func filterByType(addrType int, addrs *[]api.Address) (*[]api.Address, *[]api.Address) {
	if addrType <= -1 {
		return addrs, &[]api.Address{}
	}
	filteredAddrs := []api.Address{}
	remainder := []api.Address{}
	for key, _ := range *addrs {
		if (*addrs)[key].Type == uint8(addrType) {
			filteredAddrs = append(filteredAddrs, (*addrs)[key])
		} else {
			remainder = append(remainder, (*addrs)[key])
		}
	}
	return &filteredAddrs, &remainder
}

func removeAddr(addr api.Address, addrs *[]api.Address) *[]api.Address {
	for key, _ := range *addrs {
		if addr.Location == (*addrs)[key].Location &&
			addr.Sublocation == (*addrs)[key].Sublocation &&
			addr.Port == (*addrs)[key].Port {
			first := (*addrs)[:key]
			second := (*addrs)[key+1 : len(*addrs)]
			*addrs = append(first, second...)
		}
	}
	return addrs
}

func updateAddrs(addrs *[]api.Address) (*[]api.Address, error) {
	updatedAddrs := Pinger(*addrs)
	err := pers.AddrTrustedInsert(&updatedAddrs)
	if err != nil {
		return &[]api.Address{}, errors.Wrap(err, "findOnlineNodes encountered an error in AddrTrustedInsert.")
	}
	return &updatedAddrs, nil
}

// if count == 0, we do a full-range search and return all live nodes.
// addrType == -1 : Give me ANY nodes
// addrType == -2 : Attempt to give me the nodes I've not synced before. If not, anything works.
// Why the negative numbers? because anything 0<= is used as a real value in addrType.
func findOnlineNodes(count int, addrType int, excl *[]api.Address) ([]api.Address, error) {
	start := api.Timestamp(time.Now().Unix())
	// filteredAddrs := &[]api.Address{}
	// unfilteredAddrs := &[]api.Address{}
	addrs, err := getAllAddresses(true) // desc - last synced first
	if err != nil {
		return []api.Address{}, errors.Wrap(err, "findOnlineNodes: getAllAddresses within this function failed.")
	}
	logging.Logf(1, "All addresses: %s", Dbg_convertAddrSliceToNameSlice(*addrs))
	addrs, _ = filterByType(addrType, addrs)
	logging.Logf(1, "Filtered addresses: %s", Dbg_convertAddrSliceToNameSlice(*addrs))
	// If filter produces nothing from the DB, we revert to standard type 2 addresses.
	// if len(*filteredAddrs) > 0 {
	// 	addrs = filteredAddrs
	// } else {
	// 	standardAddrs, _ := filterByType(2, unfilteredAddrs)
	// 	addrs = standardAddrs
	// }
	if excl != nil {
		for _, addr := range *excl {
			addrs = removeAddr(addr, addrs)
		}
	}
	updatedAddrs, err := updateAddrs(addrs)
	logging.Logf(1, "Updated addresses: %s", Dbg_convertAddrSliceToNameSlice(*updatedAddrs))
	if err != nil {
		errors.Wrap(err, "findOnlineNodes: updateAddress within this function failed.")
	}
	liveNodes := filterByLastSuccessfulPing(updatedAddrs, start)
	logging.Logf(1, "Live addresses: %s", Dbg_convertAddrSliceToNameSlice(*updatedAddrs))
	// if len(*liveNodes) == 0 {
	// 	logging.Logf(1, "Our first pass at livenodes at this type (%#v) ended up yielding no usable nodes. reverting back to Address.Type=2 search.", addrType)
	// 	// We scanned this type and found no live nodes of given type.
	// 	// Revert to standard type 2 node, and scan. Return from those.
	// 	standardAddrs, _ := filterByType(2, unfilteredAddrs)
	// 	if excl != nil {
	// 		for _, addr := range *excl {
	// 			standardAddrs = removeAddr(addr, standardAddrs)
	// 		}
	// 	}
	// 	err := errors.New("")
	// 	updatedAddrs, err = updateAddrs(standardAddrs)
	// 	if err != nil {
	// 		errors.Wrap(err, "findOnlineNodes: updateAddress within this function failed.")
	// 	}
	// 	liveNodes = filterByLastSuccessfulPing(updatedAddrs, start)
	// 	if len(*liveNodes) == 0 { // If it's still zero after flipping to 2, bail.
	// 		return *liveNodes, errors.New("This database has no addresses online.")
	// 	}
	// }
	if count == 0 { // count == 0: return everything found.
		return *liveNodes, nil
	}
	if addrType == -2 {
		// logging.Logf(1, "Live nodes are these. Live nodes: %s", Dbg_convertAddrSliceToNameSlice(*liveNodes))
		logging.Log(1, "AddrType = -2, we are looking for nonconnected addrs.")
		nonconnected := pickUnconnectedAddrs(liveNodes)
		if len(*nonconnected) != 0 {
			logging.Logf(1, "AddrType = -2, we found some nonconnected onlines. Let's pull from those first. Found: %s", Dbg_convertAddrSliceToNameSlice(*nonconnected))
			liveNodes = nonconnected
		}
	}
	if len(*liveNodes) == 0 { // If zero, bail.
		return *liveNodes, errors.New("This database has no addresses online.")
	}
	rands := toolbox.GetInsecureRands(len(*liveNodes), count)
	selected := []api.Address{}
	// fmt.Println(rands)
	// fmt.Println(len(*liveNodes))
	// spew.Dump(*liveNodes)
	for _, val := range rands {
		selected = append(selected, (*liveNodes)[val])
	}
	return selected, nil
}

func pickUnconnectedAddrs(addrs *[]api.Address) *[]api.Address {
	nonconnecteds := []api.Address{}
	for key, _ := range *addrs {
		if (*addrs)[key].LastSuccessfulSync == 0 {
			nonconnecteds = append(nonconnecteds, (*addrs)[key])
		}
	}
	return &nonconnecteds
}

func RefreshAddresses() error {
	addrs, err := getAllAddresses(false) // asc - the oldest unconnected first
	if err != nil {
		return errors.Wrap(err, "RefreshAddresses: getAllAddresses within this function failed.")
	}
	updateAddrs(addrs)
	return nil
}

func GetOnlineAddresses(
	noOfOnlineAddressesRequested int,
	exclude []api.Address,
	addressType uint8,
	forceUnconnected bool,
) (
	[]api.Address, error,
) {
	ln, err := findOnlineNodes(noOfOnlineAddressesRequested, int(addressType), &exclude)
	spew.Dump(ln)
	return ln, err
}

func AddressScanner() {
	globals.BackendTransientConfig.AddressesScannerActive.Lock()
	defer globals.BackendTransientConfig.AddressesScannerActive.Unlock()
	err := RefreshAddresses()
	if err != nil {
		logging.Log(1, fmt.Sprintf("AddressScanner failed. Error: %#v", err))
		// return errors.Wrap(err, "AddressScanner failed.")
	}
}

func GetUnconnAddr(count int) []api.Address {
	fmt.Println("GetUnconnAddr hits")
	addrs, err := findOnlineNodes(count, -2, nil)
	// fmt.Println(len(addrs))
	if err != nil {
		logging.Log(1, fmt.Sprintf("Unconnected address search failed. Error: %#v", err))
		return []api.Address{}
	}
	return addrs
}

func Dbg_convertAddrSliceToNameSlice(nodes []api.Address) []string {
	names := []string{}
	for _, val := range nodes {
		names = append(names, val.Client.ClientName)
	}
	return names
}
