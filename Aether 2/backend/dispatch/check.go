// Backend > Routines > Check
// This file contains the dispatch routines that dispatch uses to deal with certain cases such as dealing with an update, encountering a new node, etc.

package dispatch

import (
	// "aether-core/backend/responsegenerator"
	"aether-core/io/api"
	// "aether-core/io/persistence"
	"aether-core/services/globals"
	// "aether-core/services/logging"
	// tb "aether-core/services/toolbox"
	// "aether-core/services/verify"
	"errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/fatih/color"
	"net"
	// "strconv"
	"strings"
	"time"
)

// Check is the short routine that reaches out to a node to see if it is online, and if so, pull the node data. This returns an updated api.Address object. Sync logic uses check as a starting point.
func Check(a api.Address) (api.Address, bool, api.ApiResponse, error) {
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check is being called for: %s:%d\n", a.Location, a.Port)
	// }
	NODE_STATIC := false
	/*
	   - Status GET to check if the node is online.
	*/
	_, err := api.Fetch(string(a.Location), string(a.Sublocation), a.Port, "status", "GET", []byte{})
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check error: %#v\n", err)
	// }
	if err != nil && (strings.Contains(err.Error(), "Client.Timeout exceeded") || strings.Contains(err.Error(), "i/o timeout")) {
		// CASE: NO RESPONSE
		// The node is offline. It can actually be offline or just too slow, but for our purposes it's the same. The timeout can be set from globals.
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, err //
	} else if err != nil {
		// CASE: CONNECTION REFUSED
		// This is where 'connection refused' would go.
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, err
	}
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through first pass for: %s:%d\n", a.Location, a.Port)
	// }
	/*
	   - The node is online. Ask for node data.
	   (This is a legitimate user of GetPageRaw because the other entities that use Check sometimes need NodeId and other fields within it.)
	*/
	apiResp, err2 := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "GET", []byte{})
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check error: %#v\n", err2)
	// }
	if err2 != nil {
		return api.Address{}, NODE_STATIC, apiResp, err2
	}
	if apiResp.NodeId == api.Fingerprint(globals.BackendConfig.GetNodeId()) {
		/*
		   This node is using the same NodeId as we do. This is, in most cases, a node connecting to itself over a loopback interface. Most router will not allow their own address to be pinged from within network, but in testing and in other rare occasions this can happen.
		*/
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, errors.New(fmt.Sprintf("This node appears to have found itself through a loopback interface, or via calling its own IP. IP: %s:%d", a.Location, a.Port))
	}
	if apiResp.Address.Type == 255 || apiResp.Address.Type == 254 {
		NODE_STATIC = true
		// 255: static node
		// 254: static bootstrapper node
	}
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through GET for: %s:%d\n", a.Location, a.Port)
	// }
	/*
	   - If the node is not static, present yourself.
	*/
	var postApiResp api.ApiResponse
	if !NODE_STATIC {
		// apiReq := responsegenerator.GeneratePrefilledApiResponse()
		apiReq := api.ApiResponse{}
		apiReq.Prefill()
		signingErr := apiReq.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
		if signingErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, signingErr
		}
		reqAsJson, jsonErr := apiReq.ToJSON()
		// reqAsJson, jsonErr := responsegenerator.ConvertApiResponseToJson(&apiReq)
		if jsonErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, jsonErr
		}
		var err3 error
		postApiResp, err3 = api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "POST", reqAsJson) // Raw call instead of regular one because we need access to the inbound remote timestamp.
		// if a.Location == "127.0.0.1" {
		// 	fmt.Printf("Check error: %#v\n", err3)
		// }
		if err3 != nil {
			// Mind that this can fail for verification failure also (if 4 entities in a page fails verification, the page fails verification. This is a page that actually has entities.)
			return api.Address{}, NODE_STATIC, apiResp, errors.New(fmt.Sprintf("Getting POST Endpoint in Check() routine for this entity type failed. Endpoint type: %s, Error: %s", "node", err3))
		}
	}
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through POST for: %s:%d\n", a.Location, a.Port)
	// }
	/*
	   - Collect the newly built address data.
	*/
	var addr api.Address
	var lastSuccessfulPing api.Timestamp
	if NODE_STATIC {
		addr = apiResp.Address
		lastSuccessfulPing = apiResp.Timestamp
	} else {
		addr = postApiResp.Address // addr is what comes from remote, a is local.
		lastSuccessfulPing = api.Timestamp(time.Now().Unix())
	}
	addr = *insertFirstPartyAddressData(&addr, &a, lastSuccessfulPing)
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through all. Resulting address: %#v\n", addr)
	// }
	return addr, NODE_STATIC, apiResp, nil
}

/*
//////////
Internal functions
//////////
*/

// insertFirstPartyAddressData inserts the first-party data we know about this address.
func insertFirstPartyAddressData(inboundAddrPtr *api.Address, localAddrPtr *api.Address, lastSuccessfulPing api.Timestamp) *api.Address {
	addr := *inboundAddrPtr
	addr.Location = localAddrPtr.Location // We know this to be true, because we just connected to it through a.Location. If the remote says it's a different IP, it's lying.
	addr.Sublocation = localAddrPtr.Sublocation
	// Determine IP type from the local address we just used to connect to this remote.
	ipAddrAsIP := net.ParseIP(string(localAddrPtr.Location))
	ipV4Test := ipAddrAsIP.To4()
	if ipV4Test == nil {
		// This is an IpV6 address
		addr.LocationType = 6
	} else {
		addr.LocationType = 4
	}
	addr.Port = localAddrPtr.Port // Because we just connected to this port and it worked. If the remote says it's a different port, it's lying.
	addr.LastSuccessfulPing = lastSuccessfulPing
	addr.EntityVersion = globals.BackendTransientConfig.EntityVersions.Address
	return &addr
}
