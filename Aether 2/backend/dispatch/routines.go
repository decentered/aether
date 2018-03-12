// Backend > Routines
// This file contains the dispatch routines that dispatch uses to deal with certain cases such as dealing with an update, encountering a new node, etc.

package dispatch

import (
	"aether-core/backend/responsegenerator"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/logging"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// Sync is the core logic of a single connection. It pulls updates from a remote node and patches it to the current node.
func Sync(a api.Address) error {
	// --------------------
	// Steps
	// - Fetch /status GET to see if the node is online.
	// - Fetch /status POST to see if the node data is valid. Save the Address with the updated last online.
	// - Check if there is a record of the node in the nodes table. If not, create it.
	// - For every entity endpoint, hit the caches that have a later end date than the timestamp for that entity endpoint.
	// - If the node is not static, do posts requets with the update timestamp of that endpoint. The remote will automatically filter the response down to entities that came after the end of the last cache.
	// - At the completion of every endpoint (get + post), save the timestamp.
	// --------------------
	// // - Status GET to check if the node is online.
	// _, err := api.Fetch(string(a.Location), string(a.Sublocation), a.Port, "status", "GET", []byte{})
	// if err != nil {
	// 	return err
	// }
	// // Reach out to node endpoint to see the node type. We want to see if the node is static, because if so we cannot make post requests.
	// NODE_STATIC := false
	// apiResp, err2 := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "GET", []byte{})
	// if err2 != nil {
	// 	return err2
	// }
	// if apiResp.Address.Type == 255 {
	// 	NODE_STATIC = true
	// }
	// // - Commit the oncoming address into the database, after adding in the basic values.
	// addr := apiResp.Address    // addr is what comes from remote, a is local.
	// addr.Location = a.Location // We know this to be true, because we just connected to it through a.Location (i.e. this is an outbound connection, if it made it to here, a.Location is correct by definition.)
	// addr.Sublocation = a.Sublocation
	// // Determine IP type from the local address we just used to connect to this remote.
	// ipAddrAsIP := net.ParseIP(string(a.Location))
	// ipV4Test := ipAddrAsIP.To4()
	// if ipV4Test == nil {
	// 	// This is an IpV6 address
	// 	addr.LocationType = 6
	// } else {
	// 	addr.LocationType = 4
	// }
	// addr.Port = a.Port
	// addr.LastOnline = api.Timestamp(time.Now().Unix())
	logging.Log(1, fmt.Sprintf("SYNC STARTED with node: %s:%d", a.Location, a.Port))
	defer logging.Log(1, fmt.Sprintf("SYNC COMPLETE with node: %s:%d", a.Location, a.Port))
	addr, NODE_STATIC, apiResp, err := Check(a)
	if err != nil {
		return err
	}
	// FULLY TRUSTED ADDRESS ENTRY
	// Anything here will be committed in and will write over existing data, since all of this data is either coming from a first-party remote, or from the client.
	addrs := []api.Address{addr}
	persistence.InsertOrUpdateAddresses(&addrs)

	// - Check if there is a record of this node in the nodes table. If not so, create and commit.
	var n persistence.DbNode
	var err4 error
	n, err4 = persistence.ReadNode(apiResp.NodeId)
	if err4 != nil && strings.Contains(err4.Error(), "The node you have asked for could not be found") {
		// Node does not exist in the DB. Create it and commit it to DB.
		n.Fingerprint = apiResp.NodeId
		err5 := persistence.InsertNode(n)
		if err5 != nil {
			// DB commit error.
			return err5
		}
	} else if err4 != nil {
		// We have an error in node query and it's not 'node not found'
		return err4
	}
	// For every endpoint, hit the caches. If the node is not static, hit the POSTs too.
	endpoints := map[string]api.Timestamp{
		"boards":      n.BoardsLastCheckin,
		"threads":     n.ThreadsLastCheckin,
		"posts":       n.PostsLastCheckin,
		"votes":       n.VotesLastCheckin,
		"addresses":   n.AddressesLastCheckin,
		"keys":        n.KeysLastCheckin,
		"truststates": n.TruststatesLastCheckin}
	logging.Log(1, fmt.Sprintf("SYNC:PULL STARTED with data from node: %s:%d", a.Location, a.Port))
	logging.Log(1, fmt.Sprintf("Endpoints: %#v", endpoints))

	for key, val := range endpoints {
		// // GET
		// Do an endpoint GET with the timestamp. (Mind that the timestamp is being provided into the GetEndpoint, it will only fetch stuff after that timestamp.)
		logging.Log(1, fmt.Sprintf("Asking for entity type: %s", key))
		resp, err6 := api.GetEndpoint(string(a.Location), string(a.Sublocation), a.Port, key, val)
		if err6 != nil {
			logging.Log(1, fmt.Sprintf("Getting GET Endpoint for the entity type '%s' failed. Error: %s, Address: %#v", key, err6, a))
		}
		logging.Log(2, fmt.Sprintf("Response to be moved to the interface pack: %#v", resp))
		// Move the objects into an interface to prepare them to be committed.
		iface := moveEntitiesToInterfacePack(&resp)
		// Save the response to the database.
		persistence.BatchInsert(*iface)
		// Set the last checkin timestamp for each entity type to the beginning of this process. (We will update this later before committing the node checkin set based on the POST response receipts, if any)
		endpoints[key] = apiResp.Timestamp
		// GET portion of this sync is done. Now on to POST requests.

		// // POST
		// POST requests can have two types of responses. If the results of that POST request is few enough, the data might just be provided as a response to the post request directly. Or, if there are many pages of results, the remote saves these into a folder that is available for the next half hour or so, and sends back the link to that folder. The two cases below deal with this.
		if !NODE_STATIC {
			// Generate the POST request.
			// POST request is essentially an ApiResponse converted to JSON. This can have fields like:
			// "filters": [
			//  {"type":"timestamp", "values": ["0", "1483641920"]}
			//  ]
			// which allows us to filter. But if you create an empty request for POST to an entity endpoint, it will give you all the entities for that endpoint since the last cache generation, automatically. There are no filters required for that kind of query.

			// But before anything, we need to create the mapping for the endpoint URLs.
			endpointsMap := map[string]string{
				"boards":      "c0/boards",
				"threads":     "c0/threads",
				"posts":       "c0/posts",
				"votes":       "c0/votes",
				"keys":        "c0/keys",
				"truststates": "c0/truststates",
				"addresses":   "addresses",
			}

			apiReq := responsegenerator.GeneratePrefilledApiResponse()
			reqAsJson, jsonErr := responsegenerator.ConvertApiResponseToJson(apiReq)
			if jsonErr != nil {
				return jsonErr
			}
			postApiResp, err7 := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, endpointsMap[key], "POST", reqAsJson) // Raw call instead of regular one because we need access to the inbound remote timestamp.
			if err7 != nil {
				return errors.New(fmt.Sprintf("Getting POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", key, err7))
			}
			var postResp api.Response
			postResp = api.InsertApiResponseToResponse(postResp, postApiResp)
			// Now, check if this is an one-page response, or links to another location for a cache hit.
			if len(postResp.CacheLinks) > 0 { // This response needed more than one page, so the remote split it into multiple pages, and saved it to a cache.
				fmt.Println("THese are the cache links we received.")
				fmt.Println(postResp.CacheLinks)
				// We're adding /responses/ because that's where the singular responses will be.
				postResultResp, err8 := api.GetCache(string(a.Location), string(a.Sublocation), a.Port, fmt.Sprintf("responses/%s", postResp.CacheLinks[0].ResponseUrl)) // There is only one if it's a prepared request.
				if err8 != nil {
					return errors.New(fmt.Sprintf("Getting Multi page POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", key, err8))
				}
				postresultIface := moveEntitiesToInterfacePack(&postResultResp)
				persistence.BatchInsert(*postresultIface)
			} else {
				// This response is one page, so the result is embedded into the POST response itself. Simple.
				postIface := moveEntitiesToInterfacePack(&postResp)
				persistence.BatchInsert(*postIface)
			}
			endpoints[key] = postApiResp.Timestamp
		}
	}
	logging.Log(1, fmt.Sprintf("SYNC:PULL COMPLETE with data from node: %s:%d", a.Location, a.Port))
	// Both POST and GETs are committed into the database. We now need to save the Node LastCheckin timestamps into the database.
	n.BoardsLastCheckin = endpoints["boards"]
	n.ThreadsLastCheckin = endpoints["threads"]
	n.PostsLastCheckin = endpoints["posts"]
	n.VotesLastCheckin = endpoints["votes"]
	n.AddressesLastCheckin = endpoints["addresses"]
	n.KeysLastCheckin = endpoints["keys"]
	n.TruststatesLastCheckin = endpoints["truststates"]
	err9 := persistence.InsertNode(n)
	if err9 != nil {
		return err9
	}
	return nil // TODO: This could return something more informative, about the status of the sync that was just completed.
}

func moveEntitiesToInterfacePack(r *api.Response) *[]interface{} {
	resp := *r
	var carrier []interface{}
	for i, _ := range resp.Boards {
		carrier = append(carrier, resp.Boards[i])
	}
	for i, _ := range resp.Threads {
		carrier = append(carrier, resp.Threads[i])
	}
	for i, _ := range resp.Posts {
		carrier = append(carrier, resp.Posts[i])
	}
	for i, _ := range resp.Votes {
		carrier = append(carrier, resp.Votes[i])
	}
	for i, _ := range resp.Addresses {
		carrier = append(carrier, resp.Addresses[i])
	}
	for i, _ := range resp.Keys {
		carrier = append(carrier, resp.Keys[i])
	}
	for i, _ := range resp.Truststates {
		carrier = append(carrier, resp.Truststates[i])
	}
	return &carrier
}

// Check is the short routine that reaches out to a node to see if it is online, and if so, pull the node data. This returns an updated api.Address object. Sync logic uses check as a starting point.
func Check(a api.Address) (api.Address, bool, api.ApiResponse, error) {
	NODE_STATIC := false
	/*
		- Status GET to check if the node is online.
	*/
	_, err := api.Fetch(string(a.Location), string(a.Sublocation), a.Port, "status", "GET", []byte{})
	if err != nil && (strings.Contains(err.Error(), "Client.Timeout exceeded") || strings.Contains(err.Error(), "i/o timeout")) {
		// CASE: NO RESPONSE
		// The node is offline. It can actually be offline or just too slow, but for our purposes it's the same. The timeout can be set from globals.
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, err //
	} else if err != nil {
		// CASE: CONNECTION REFUSED
		// This is where 'connection refused' would go.
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, err
	}
	/*
		- The node is online. Ask for node data.
	*/
	apiResp, err2 := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "GET", []byte{})
	if err2 != nil {
		return api.Address{}, NODE_STATIC, apiResp, err2
	}
	if apiResp.Address.Type == 255 {
		NODE_STATIC = true
	}
	/*
		- If the node is not static, present yourself.
	*/
	var postApiResp api.ApiResponse
	if !NODE_STATIC {
		apiReq := responsegenerator.GeneratePrefilledApiResponse()
		reqAsJson, jsonErr := responsegenerator.ConvertApiResponseToJson(apiReq)
		if jsonErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, jsonErr
		}
		var err3 error
		postApiResp, err3 = api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "POST", reqAsJson) // Raw call instead of regular one because we need access to the inbound remote timestamp.
		if err3 != nil {
			return api.Address{}, NODE_STATIC, apiResp, errors.New(fmt.Sprintf("Getting POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", "node", err3))
		}
	}
	/*
		- Collect the newly built address data.
	*/
	var addr api.Address
	if NODE_STATIC {
		addr = apiResp.Address
	} else {
		addr = postApiResp.Address // addr is what comes from remote, a is local.
	}
	// logging.LogCrash(fmt.Sprintf("%#v", postApiResp.Address))
	addr.Location = a.Location // We know this to be true, because we just connected to it through a.Location (i.e. this is an outbound connection, if it made it to here, a.Location is correct by definition.)
	addr.Sublocation = a.Sublocation
	// Determine IP type from the local address we just used to connect to this remote.
	ipAddrAsIP := net.ParseIP(string(a.Location))
	ipV4Test := ipAddrAsIP.To4()
	if ipV4Test == nil {
		// This is an IpV6 address
		addr.LocationType = 6
	} else {
		addr.LocationType = 4
	}
	addr.Port = a.Port
	// todo: lastonline should have its own logic for static
	addr.LastOnline = api.Timestamp(time.Now().Unix())
	// fmt.Printf("Resulting address at the end of the check process %#v", addr)
	// Addr is the container for the newly obtained address data.
	return addr, NODE_STATIC, apiResp, nil
}
