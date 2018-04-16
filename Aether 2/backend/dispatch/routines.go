// Backend > Routines
// This file contains the dispatch routines that dispatch uses to deal with certain cases such as dealing with an update, encountering a new node, etc.

package dispatch

import (
	"aether-core/backend/responsegenerator"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	// "aether-core/services/verify"
	"errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"net"
	"strconv"
	"strings"
	"time"
)

// Metrics container for every sync.
type CurrentOutboundSyncMetrics struct {
	BoardsReceived                 int
	BoardsSinglePage               bool
	BoardsGETNetworkTime           float64
	BoardsPOSTNetworkTime          float64
	BoardsPOSTTimeToFirstResponse  float64
	BoardsDBCommitTime             float64
	BoardOwnerDBCommitTime         float64
	BoardOwnerDeletionDBCommitTime float64

	ThreadsReceived                int
	ThreadsSinglePage              bool
	ThreadsGETNetworkTime          float64
	ThreadsPOSTNetworkTime         float64
	ThreadsPOSTTimeToFirstResponse float64
	ThreadsDBCommitTime            float64

	PostsReceived                int
	PostsSinglePage              bool
	PostsGETNetworkTime          float64
	PostsPOSTNetworkTime         float64
	PostsPOSTTimeToFirstResponse float64
	PostsDBCommitTime            float64

	VotesReceived                int
	VotesSinglePage              bool
	VotesGETNetworkTime          float64
	VotesPOSTNetworkTime         float64
	VotesPOSTTimeToFirstResponse float64
	VotesDBCommitTime            float64

	KeysReceived                int
	KeysSinglePage              bool
	KeysGETNetworkTime          float64
	KeysPOSTNetworkTime         float64
	KeysPOSTTimeToFirstResponse float64
	KeysDBCommitTime            float64

	TruststatesReceived                int
	TruststatesSinglePage              bool
	TruststatesGETNetworkTime          float64
	TruststatesPOSTNetworkTime         float64
	TruststatesPOSTTimeToFirstResponse float64
	TruststatesDBCommitTime            float64

	AddressesReceived                int
	AddressesSinglePage              bool
	AddressesGETNetworkTime          float64
	AddressesPOSTNetworkTime         float64
	AddressesPOSTTimeToFirstResponse float64
	AddressesDBCommitTime            float64

	LocalIp                 string
	LocalPort               int
	RemoteIp                string
	RemotePort              int
	TotalDurationSeconds    int
	TotalNetworkRemoteWait  float64
	DbInsertDurationSeconds int
	LocalClientName         string
	RemoteClientName        string
	SyncHistory             string
}

func startMetricsContainer(
	resp api.ApiResponse,
	remoteAddr api.Address,
	n persistence.DbNode) *CurrentOutboundSyncMetrics {
	var c CurrentOutboundSyncMetrics
	c.RemoteIp = string(remoteAddr.Location)
	c.RemotePort = int(remoteAddr.Port)
	c.LocalIp = string(globals.BackendConfig.GetExternalIp())
	c.LocalPort = int(globals.BackendConfig.GetExternalPort())
	c.LocalClientName = globals.BackendTransientConfig.AppIdentifier
	c.RemoteClientName = resp.Address.Client.ClientName
	c.SyncHistory = "First sync"
	if n.BoardsLastCheckin > 0 ||
		n.ThreadsLastCheckin > 0 ||
		n.PostsLastCheckin > 0 ||
		n.VotesLastCheckin > 0 ||
		n.KeysLastCheckin > 0 ||
		n.TruststatesLastCheckin > 0 ||
		n.AddressesLastCheckin > 0 {
		c.SyncHistory = "Resync"
	}
	return &c
}

func generateStartMessage(c *CurrentOutboundSyncMetrics, clr *color.Color) string {
	openMessage := clr.Sprintf("\nOPEN: %s:%d (%s) >>> %s:%d (%s) - %s ",
		c.LocalIp, c.LocalPort, c.LocalClientName, c.RemoteIp, c.RemotePort,
		c.RemoteClientName, c.SyncHistory,
	)
	return openMessage
}

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

	logging.Log(2, fmt.Sprintf("SYNC STARTED with node: %s:%d", a.Location, a.Port))
	start := time.Now()
	addr, NODE_STATIC, apiResp, err := Check(a)
	if err != nil {
		return err
	}

	// FULLY TRUSTED ADDRESS ENTRY
	// Anything here will be committed in and will write over existing data, since all of this data is either coming from a first-party remote, or from the client.
	addrs := []api.Address{addr}
	errs := persistence.InsertOrUpdateAddresses(&addrs)
	if len(errs) > 0 {
		err := errors.New(fmt.Sprintf("Some errors were encountered when the Sync attempted InsertOrUpdateAddresses. Sync aborted. Errors: %s", errs))
		logging.Log(1, err)
		abortClr := color.New(color.FgWhite, color.BgRed)
		logging.Log(1, abortClr.Sprintf("SYNC ABORTED. Err: %s", err))
		return err
	}
	// - Check if there is a record of this node in the nodes table. If not so, create and commit.
	var n persistence.DbNode
	var err4 error
	n, err4 = persistence.ReadNode(api.Fingerprint(apiResp.NodeId))
	if err4 != nil && strings.Contains(err4.Error(), "The node you have asked for could not be found") {
		// Node does not exist in the DB. Create it and commit it to DB.
		n.Fingerprint = api.Fingerprint(apiResp.NodeId)
		err5 := persistence.InsertNode(n)
		if err5 != nil {
			// DB commit error, or node was using the same id as ours.
			return err5
		}
	} else if err4 != nil {
		// We have an error in node query and it's not 'node not found'
		return err4
	}
	c := startMetricsContainer(apiResp, a, n)
	openClr := color.New(color.FgWhite, color.BgYellow)
	logging.Log(2, generateStartMessage(c, openClr))
	// For every endpoint, hit the caches. If the node is not static, hit the POSTs too.
	endpoints := map[string]api.Timestamp{
		"boards":      n.BoardsLastCheckin,
		"threads":     n.ThreadsLastCheckin,
		"posts":       n.PostsLastCheckin,
		"votes":       n.VotesLastCheckin,
		"addresses":   n.AddressesLastCheckin,
		"keys":        n.KeysLastCheckin,
		"truststates": n.TruststatesLastCheckin}
	logging.Log(2, fmt.Sprintf("SYNC:PULL STARTED with data from node: %s:%d", a.Location, a.Port))
	logging.Log(2, fmt.Sprintf("Endpoints: %#v", endpoints))
	ims := []persistence.InsertMetrics{}
	// callOrder := []string{"boards", "threads", "posts", "votes", "addresses", "keys", "truststates"}
	callOrder := []string{"addresses", "votes", "truststates", "posts", "threads", "boards", "keys"}
	// callOrder := []string{"boards"}
	for _, endpointName := range callOrder {
		fmt.Println(endpointName)
		start := time.Now()
		// // GET
		// Do an endpoint GET with the timestamp. (Mind that the timestamp is being provided into the GetEndpoint, it will only fetch stuff after that timestamp.)
		logging.Log(2, fmt.Sprintf("Asking for entity type: %s", endpointName))
		resp, err6 := api.GetEndpoint(string(a.Location), string(a.Sublocation), a.Port, endpointName, endpoints[endpointName])
		if err6 != nil {
			logging.Log(2, fmt.Sprintf("Getting GET Endpoint for the entity type '%s' failed. Error: %s, Address: %#v", endpointName, err6, a))
		}
		logging.Log(2, fmt.Sprintf("Response to be moved to the interface pack: %#v", resp))
		elapsed := time.Since(start) // We end this counter before DB insert starts, because this is the network-time counter.
		// Move the objects into an interface to prepare them to be committed.
		iface := prepareForBatchInsert(&resp)
		// Save the response to the database.
		im, err := persistence.BatchInsert(*iface)
		if err != nil {
			logging.LogCrash(err)
		}
		ims = append(ims, im)
		// Set the last checkin timestamp for each entity type to the beginning of this process. (We will update this later before committing the node checkin set based on the POST response receipts, if any)
		// Check if the apiResp.Timestamp is newer or older than the timestamp we have. It might actually be older,because we might have received a POST response from this node, and that might have been a later Timestamp than that of the last cache's.

		if endpoints[endpointName] < resp.MostRecentSourceTimestamp {
			endpoints[endpointName] = resp.MostRecentSourceTimestamp
		}
		// Insert the metrics into the container.
		if endpointName == "boards" {
			c.BoardsGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "threads" {
			c.ThreadsGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "posts" {
			c.PostsGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "votes" {
			c.VotesGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "keys" {
			c.KeysGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "truststates" {
			c.TruststatesGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "addresses" {
			c.AddressesGETNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
		}
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

			start := time.Now()
			var elapsed time.Duration
			/*
				HEADS UP:
				Post "Pingpong"
				We flip this to true if the POST response was 1 page. Why does it matter? Because our logs provide only candidate inputs, and if the result is only one page, it is not paginated, thus doesn't have indexes, thus it is not subject to 'exists in db' checks. By the virtue of getting the result, you've already downloaded it (i.e. you can't save bandwidth by not hitting the pages, it is already one page and you already have that page), so there is no point in doing this checks as DB will check it already anyway. But it creates a confusing side effect where it looks like the node has sent unnecessary data exclusively for the endpoints that were single-page because unneeded data wasn't filtered out by the local.

				In actuality, the unnecessary data is being sent by all endpoints in POST, but when there is a multipage post response, the local client can go through the index and not hit the pages it does not need. So they look like zeroes. But that check doesn't happen for single-page responses.

				Why is unnecessary data being sent?

				Let's look at an example.

				t0 1 > 0 : FirstS : N1 gets 31 new keys (total 62) ts:t0

				t1 1 > 0 : Resync : nothing  ts:t1

				t2 0 > 1 : FirstS : N0 gets 62 keys (the keys unique to 1, plus all 0's keys) 31 NEW KEYS @ N0 at ts:t2

				t3 1 > 0 : Resync : N1 gets 31 keys (because the ts N1 has for N0 is ts:t1, and all the keys N0 just got from N1 are at ts:t2)

				t4 0 > 1 : nothing ts:t4

				t5 1 > 0 : nothing ts:t5

				The crux is at t3, the data in N1 got transmitted to N0, and N0 transmitted it back to N1. This is because Mim connections are one-way. They are not two-way syncs. The nodes cannot retain information of which node they got the results from, because that node might have been reset, and might actually need that information. The local nodes cannot communicate what they have in their database, because that would be a privacy violation.

				Why is this not an issue?
				- It does not cost bandwidth except in the case of a single page. Single page means very little bandwidth is used. Any response of decent size will be multipart, and in multipart responses, the local node will check the index, discover that the things in the index are those that it already has, and it won't download the pages.
				- This only happens to POST responses, and not GET. POST responses are the 'tip of the spear', the usual way of process is that you exhaust the caches with GET first, and then get the delta from the end of the last cache to now via the POST response. Caches are generated frequently, the time range that the POST response will cover will be in the order of hours, not days / weeks. Most of the network traffic is GET.

			*/
			var singlePage bool
			apiReq := responsegenerator.GeneratePrefilledApiResponse()
			// Here, we need to insert the last sync timestamp into the post request, so that it will be gated appropriately.
			f := api.Filter{}
			f.Type = "timestamp"
			f.Values = []string{strconv.Itoa(int(endpoints[endpointName])), strconv.Itoa(0)}
			apiReq.Filters = []api.Filter{f}
			signingErr := apiReq.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
			if signingErr != nil {
				return signingErr
			}
			reqAsJson, jsonErr := responsegenerator.ConvertApiResponseToJson(apiReq)
			if jsonErr != nil {
				return jsonErr
			}
			postResp, respDuration, err7 := api.GetPage(string(a.Location), string(a.Sublocation), a.Port, endpointsMap[endpointName], "POST", reqAsJson)
			if err7 != nil {
				return errors.New(fmt.Sprintf("Getting POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", endpointName, err7))
			}
			// // Now, check if this is an one-page response, or links to another location for a cache hit.
			if len(postResp.CacheLinks) > 0 { // This response needed more than one page, so the remote split it into multiple pages, and saved it to a cache.
				postResultResp, err8 := api.CollectMultipartPOSTResponse(string(a.Location), string(a.Sublocation), a.Port, fmt.Sprintf("responses/%s", postResp.CacheLinks[0].ResponseUrl))
				// We're adding /responses/ because that's where the singular responses will be.
				if err8 != nil {
					return errors.New(fmt.Sprintf("Getting Multi page POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", endpointName, err8))
				}
				elapsed = time.Since(start) // Ends here, since we don't want to capture DB time.
				postresultIface := prepareForBatchInsert(&postResultResp)
				im, err := persistence.BatchInsert(*postresultIface)
				if err != nil {
					logging.LogCrash(err)
				}
				ims = append(ims, im)
			} else {
				// This response is one page, so the result is embedded into the POST response itself. Simple.
				singlePage = true
				elapsed = time.Since(start) // Ends here, since we don't want to capture DB time.
				postIface := prepareForBatchInsert(&postResp)
				im, err := persistence.BatchInsert(*postIface)
				if err != nil {
					logging.LogCrash(err)
				}
				ims = append(ims, im)
			}
			endpoints[endpointName] = postResp.MostRecentSourceTimestamp
			// Insert the metrics into the container.
			if endpointName == "boards" {
				c.BoardsPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.BoardsPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.BoardsPOSTTimeToFirstResponse
				c.BoardsSinglePage = singlePage
			} else if endpointName == "threads" {
				c.ThreadsPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.ThreadsPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.ThreadsPOSTTimeToFirstResponse
				c.ThreadsSinglePage = singlePage
			} else if endpointName == "posts" {
				c.PostsPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.PostsPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.PostsPOSTTimeToFirstResponse
				c.PostsSinglePage = singlePage
			} else if endpointName == "votes" {
				c.VotesPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.VotesPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.VotesPOSTTimeToFirstResponse
				c.VotesSinglePage = singlePage
			} else if endpointName == "keys" {
				c.KeysPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.KeysPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.KeysPOSTTimeToFirstResponse
				c.KeysSinglePage = singlePage
			} else if endpointName == "truststates" {
				c.TruststatesPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.TruststatesPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.TruststatesPOSTTimeToFirstResponse
				c.TruststatesSinglePage = singlePage
			} else if endpointName == "addresses" {
				c.AddressesPOSTNetworkTime = globals.Round(elapsed.Seconds(), 0.1)
				c.AddressesPOSTTimeToFirstResponse = globals.Round(respDuration.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.AddressesPOSTTimeToFirstResponse
				c.AddressesSinglePage = singlePage
			}
		}
	}
	logging.Log(2, fmt.Sprintf("SYNC:PULL COMPLETE with data from node: %s:%d", a.Location, a.Port))
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
	logging.Log(2, fmt.Sprintf("SYNC COMPLETE with node: %s:%d. It took %d seconds", a.Location, a.Port, int(time.Since(start).Seconds())))
	closeClr := color.New(color.FgBlack, color.BgWhite)
	logging.Log(1, generateCloseMessage(c, closeClr, &ims, int(time.Since(start).Seconds()), true))
	return nil
}

func generateCloseMessage(c *CurrentOutboundSyncMetrics, clr *color.Color, ims *[]persistence.InsertMetrics, dur int, extended bool) string {
	im := persistence.InsertMetrics{}
	for _, val := range *ims {
		im.Add(val)
	}
	// spew.Dump(im)
	c.BoardsReceived = im.BoardsReceived
	c.ThreadsReceived = im.ThreadsReceived
	c.PostsReceived = im.PostsReceived
	c.PostsReceived = im.PostsReceived
	c.VotesReceived = im.VotesReceived
	c.KeysReceived = im.KeysReceived
	c.TruststatesReceived = im.TruststatesReceived
	c.AddressesReceived = im.AddressesReceived
	c.TotalDurationSeconds = dur
	c.DbInsertDurationSeconds = im.TimeElapsedSeconds
	c.BoardsDBCommitTime = im.BoardsDBCommitTime
	c.ThreadsDBCommitTime = im.ThreadsDBCommitTime
	c.PostsDBCommitTime = im.PostsDBCommitTime
	c.VotesDBCommitTime = im.VotesDBCommitTime
	c.KeysDBCommitTime = im.KeysDBCommitTime
	c.TruststatesDBCommitTime = im.TruststatesDBCommitTime
	c.AddressesDBCommitTime = im.AddressesDBCommitTime

	totalEntitiesReceived := c.BoardsReceived + c.ThreadsReceived + c.PostsReceived + c.VotesReceived + c.KeysReceived + c.TruststatesReceived + c.AddressesReceived
	insertDbDetailString := ""
	if totalEntitiesReceived > 0 {
		insertDbDetailString = fmt.Sprintf("  %d Boards, %d Threads, %d Posts, %d Votes, %d Keys, %d Truststates, %d Addresses (All before dedupe)\n  %s %s %s %s %s %s %s",
			c.BoardsReceived,
			c.ThreadsReceived,
			c.PostsReceived,
			c.VotesReceived,
			c.KeysReceived,
			c.TruststatesReceived,
			c.AddressesReceived,
			singlePageSprinter(c, "Boards"),
			singlePageSprinter(c, "Threads"),
			singlePageSprinter(c, "Posts"),
			singlePageSprinter(c, "Votes"),
			singlePageSprinter(c, "Keys"),
			singlePageSprinter(c, "Truststates"),
			singlePageSprinter(c, "Addresses"))
	}
	longTimeDetailString := ""
	shortTimeDetailString := ""
	timeDetailString := ""
	if c.TotalDurationSeconds > 0 {
		// This is where we collect db time metrics.
		dbTimeDetailString := fmt.Sprintf("\n    Boards:      %.1fs, \n    Threads:     %.1fs, \n    Posts:       %.1fs, \n    Votes:       %.1fs, \n    Keys:        %.1fs, \n    Truststates: %.1fs, \n    Addresses:   %.1fs.", c.BoardsDBCommitTime, c.ThreadsDBCommitTime, c.PostsDBCommitTime, c.VotesDBCommitTime, c.KeysDBCommitTime, c.TruststatesDBCommitTime, c.AddressesDBCommitTime)
		// network time metrics for GET and POST.
		networkTimeDetailString := fmt.Sprintf("\n    Boards:      G: %.1fs P: %.1fs (PWait: %.1f), \n    Threads:     G: %.1fs P: %.1fs (PWait: %.1f), \n    Posts:       G: %.1fs P: %.1fs (PWait: %.1f), \n    Votes:       G: %.1fs P: %.1fs (PWait: %.1f), \n    Keys:        G: %.1fs P: %.1fs (PWait: %.1f), \n    Truststates: G: %.1fs P: %.1fs (PWait: %.1f), \n    Addresses:   G: %.1fs P: %.1fs (PWait: %.1f).",
			c.BoardsGETNetworkTime, c.BoardsPOSTNetworkTime, c.BoardsPOSTTimeToFirstResponse,
			c.ThreadsGETNetworkTime, c.ThreadsPOSTNetworkTime, c.ThreadsPOSTTimeToFirstResponse,
			c.PostsGETNetworkTime, c.PostsPOSTNetworkTime, c.PostsPOSTTimeToFirstResponse,
			c.VotesGETNetworkTime, c.VotesPOSTNetworkTime, c.VotesPOSTTimeToFirstResponse,
			c.KeysGETNetworkTime, c.KeysPOSTNetworkTime, c.KeysPOSTTimeToFirstResponse,
			c.TruststatesGETNetworkTime, c.TruststatesPOSTNetworkTime, c.TruststatesPOSTTimeToFirstResponse,
			c.AddressesGETNetworkTime, c.AddressesPOSTNetworkTime, c.AddressesPOSTTimeToFirstResponse,
		)
		longTimeDetailString = fmt.Sprintf("\n  DB: %ds (%s) %s \n  Network: %ds %s", c.DbInsertDurationSeconds, globals.BackendConfig.GetDbEngine(), dbTimeDetailString, c.TotalDurationSeconds-c.DbInsertDurationSeconds, networkTimeDetailString)
		shortTimeDetailString = fmt.Sprintf("\n    DB: %ds (%s)  Network: %ds (Wait for remote: %.1fs)", c.DbInsertDurationSeconds, globals.BackendConfig.GetDbEngine(), c.TotalDurationSeconds-c.DbInsertDurationSeconds, c.TotalNetworkRemoteWait)
		if extended {
			timeDetailString = longTimeDetailString
		} else {
			timeDetailString = shortTimeDetailString
		}
	}
	closeMessage := clr.Sprintf("\nCLOSE: %s >>> %s (%s) \n(%s:%d >>> %s:%d) \nReceived: Total: %d. \n%s \nTime: Total: %ds. %s",
		c.LocalClientName, c.RemoteClientName, c.SyncHistory,
		c.LocalIp, c.LocalPort, c.RemoteIp, c.RemotePort,
		totalEntitiesReceived, insertDbDetailString,
		c.TotalDurationSeconds, timeDetailString)
	return closeMessage
}

func singlePageSprinter(c *CurrentOutboundSyncMetrics, entityType string) string {
	resp := fmt.Sprintf("single page POST")
	if entityType == "Boards" && c.BoardsSinglePage {
		return fmt.Sprintf("(%s %s - but has explicit gate)", entityType, resp)
	} else if (entityType == "Threads" && c.ThreadsSinglePage) ||
		(entityType == "Posts" && c.PostsSinglePage) ||
		(entityType == "Votes" && c.VotesSinglePage) ||
		(entityType == "Keys" && c.KeysSinglePage) ||
		(entityType == "Truststates" && c.TruststatesSinglePage) ||
		(entityType == "Addresses" && c.AddressesSinglePage) {
		return fmt.Sprintf("(%s %s)", entityType, resp)
	} else {
		return ""
	}
}

// prepareForBatchInsert verifies the items in this response container, and converts it to the correct form BatchInsert accepts.
func prepareForBatchInsert(r *api.Response) *[]interface{} {
	// // cleanedResp := verify.VerifyEntitiesInResponse(r)
	// cleanedResp := verify.BatchVerify(r)
	// resp := *cleanedResp
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
		(This is a legitimate user of GetPageRaw because the other entities that use Check sometimes need NodeId and other fields within it.)
	*/
	apiResp, err2 := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "GET", []byte{})
	if err2 != nil {
		return api.Address{}, NODE_STATIC, apiResp, err2
	}
	if apiResp.NodeId == api.Fingerprint(globals.BackendConfig.GetNodeId()) {
		/*
			This node is using the same NodeId as we do. This is, in most cases, a node connecting to itself over a loopback interface. Most router will not allow their own address to be pinged from within network, but in testing and in other rare occasions this can happen.
		*/
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, errors.New(fmt.Sprintf("This node appears to have found itself through a loopback interface, or via calling its own IP. IP: %s:%d", a.Location, a.Port))
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
		signingErr := apiReq.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
		if signingErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, signingErr
		}
		reqAsJson, jsonErr := responsegenerator.ConvertApiResponseToJson(apiReq)
		if jsonErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, jsonErr
		}
		var err3 error
		postApiResp, err3 = api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "POST", reqAsJson) // Raw call instead of regular one because we need access to the inbound remote timestamp.
		if err3 != nil {
			// Mind that this can fail for verification failure also (if 4 entities in a page fails verification, the page fails verification. This is a page that actually has entities.)
			return api.Address{}, NODE_STATIC, apiResp, errors.New(fmt.Sprintf("Getting POST Endpoint in Check() routine for this entity type failed. Endpoint type: %s, Error: %s", "node", err3))
		}
	}
	/*
		- Collect the newly built address data.
	*/
	var addr api.Address
	var lastOnline api.Timestamp
	if NODE_STATIC {
		addr = apiResp.Address
		lastOnline = apiResp.Timestamp
	} else {
		addr = postApiResp.Address // addr is what comes from remote, a is local.
		lastOnline = api.Timestamp(time.Now().Unix())
	}
	addr = *insertFirstPartyAddressData(&addr, &a, lastOnline)
	return addr, NODE_STATIC, apiResp, nil
}

// insertFirstPartyAddressData inserts the first-party data we know about this address.
func insertFirstPartyAddressData(inboundAddrPtr *api.Address, localAddrPtr *api.Address, lastOnline api.Timestamp) *api.Address {
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
	addr.LastOnline = lastOnline
	addr.EntityVersion = globals.BackendTransientConfig.EntityVersions.Address
	return &addr
}
