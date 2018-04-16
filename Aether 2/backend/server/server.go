// Backend > Server
// This file provides the backend the server that the external world accesses to get data from the backend.

package server

import (
	"aether-core/backend/responsegenerator"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NYTimes/gziphandler"
	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Server responds to GETs with the caches and to POSTS with the live data from the database.
func Serve() {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			dir := fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, dir)
		} else { // If not GET we bail.
			w.WriteHeader(http.StatusNotFound)
		}
	})
	gzippedHandler := gziphandler.GzipHandler(handlerFunc)
	http.Handle("/v0/responses/", gzippedHandler)
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Force the content type to application/json, so even in the case of malicious file serving, it won't be executed by default.
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			switch r.URL.Path {

			case "/v0/status", "/v0/status/":
				// Status GET endpoint returns HTTP 200 only if the node is up, and 429 Too Many Requests if the node is being overloaded.
				if globals.BackendTransientConfig.TooManyConnections {
					w.WriteHeader(http.StatusTooManyRequests)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.Write([]byte{})

			case "/v0/node", "/v0/node/":
				// Node GET endpoint returns the node info.
				var resp api.ApiResponse
				r := responsegenerator.GeneratePrefilledApiResponse()
				resp = *r
				resp.Endpoint = "node"
				resp.Entity = "node"
				resp.Timestamp = api.Timestamp(time.Now().Unix())
				signingErr := resp.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
				if signingErr != nil {
					logging.Log(1, fmt.Sprintf("This cache page failed to be page-signed. Error: %#v Page: %#v\n", signingErr, resp))
				}
				jsonResp, err := responsegenerator.ConvertApiResponseToJson(&resp)
				if err != nil {
					logging.Log(1, errors.New(fmt.Sprintf("The response that was prepared to respond to this query failed to convert to JSON. Error: %#v\n", err)))
				}
				if len(jsonResp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(jsonResp)
				}

			default:
				// TODO: Convert this into a whitelist. This should not respond to the random requests, only the endpoints. It also should not list directories.
				http.ServeFile(w, r, fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), r.URL.Path))
			}

		} else if r.Method == "POST" {
			switch r.URL.Path {
			case "/v0/node", "/v0/node/":
				resp, err := NodePOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/c0/boards", "/v0/c0/boards/":
				resp, err := BoardsPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/c0/threads", "/v0/c0/threads/":
				resp, err := ThreadsPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/c0/posts", "/v0/c0/posts/":
				resp, err := PostsPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/c0/votes", "/v0/c0/votes/":
				resp, err := VotesPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/c0/keys", "/v0/c0/keys/":
				resp, err := KeysPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/addresses", "/v0/addresses/":
				resp, err := AddressesPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			case "/v0/c0/truststates", "/v0/c0/truststates/":
				resp, err := TruststatesPOST(r)
				if err != nil {
					logging.Log(1, err)
				}
				if len(resp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(resp)
				}

			default:
				logging.Log(1, fmt.Sprintf("A remote reached out to this node with a request that this node does not have a route for. The requested route: %s, The node requesting: %v", r.URL.Path, r.Body))
				w.WriteHeader(http.StatusNotFound)
			}
		} else { // If not GET or POST, we bail.
			w.WriteHeader(http.StatusNotFound)
		}
	})
	gzippedMainHandler := gziphandler.GzipHandler(mainHandler)
	http.Handle("/", gzippedMainHandler)
	port := globals.BackendConfig.GetExternalPort()
	logging.Log(1, fmt.Sprintf("Serving setup complete. Starting to serve publicly on port %d", port))

	err := http.ListenAndServe(fmt.Sprint(":", port), nil)
	if err != nil {
		logging.LogCrash(fmt.Sprintf("Server encountered a fatal error. Error: %s", err))
	}
}

// SaveRemote checks if the database has data about the remote that is reaching out. If not, save a new address. We don't insert the node, only the address. Because the remote data is untrustable.
func SaveRemote(req api.ApiResponse) error {
	// spew.Dump(req.Address.Client)
	addrs := []api.Address{req.Address}
	errs := persistence.InsertOrUpdateAddresses(&addrs)
	if len(errs) > 0 {
		err := errors.New(fmt.Sprintf("Some errors were encountered when the SaveRemote attempted InsertOrUpdateAddresses. Process aborted. Errors: %s", errs))
		logging.Log(1, err)
		return err
	}
	return nil
}

// insertLocallySourcedRemoteAddressDetails Inserts the locally sourced data about the remote into the address entity that is coming with the POST request.
func insertLocallySourcedRemoteAddressDetails(r *http.Request, req *api.ApiResponse) error {
	// This runs when a node connects to you.
	// LITTLE-TRUSTED ADDRESS ENTRY
	// Data to keep: Location, Sublocation, Port, LastOnline (sublocation is guaranteed to be empty since the connection is coming from an IP, not a static IP)
	// Delete everything else, they're untrustable.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("The address from which the remote is connecting could not be parsed. Remote Address: %s, Error: %s", r.RemoteAddr, err))
	}
	if len(host) == 0 {
		return errors.New(fmt.Sprintf("The address from which the remote is connecting seems to be empty. Remote Address: %#v. %#v", r.RemoteAddr, err))
	}
	// TODO: Decide whether making a DNS request (ParseIP makes a DNS request) below is a risk (probably not).
	ipAddrAsIP := net.ParseIP(host)
	ipV4Test := ipAddrAsIP.To4()
	if ipV4Test == nil {
		// This is an IpV6 address
		req.Address.LocationType = 6
	} else {
		req.Address.LocationType = 4
	}
	req.Address.Sublocation = "" // It's coming from an IP address, not a URL.
	req.Address.Location = api.Location(host)
	req.Address.LastOnline = api.Timestamp(time.Now().Unix())
	req.Address.Type = 2 // If it is making a request to you, it cannot be a static node, by definition.
	return nil
}

// ParsePOSTRequest receives and parses the post request given by the remote.
func ParsePOSTRequest(r *http.Request) (api.ApiResponse, error) {
	var req api.ApiResponse
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return req, errors.New(fmt.Sprintf("This HTTP body could not be read. Error: %#v\n", err))
	}
	err2 := json.Unmarshal(b, &req)
	if err2 != nil {
		return req, errors.New(fmt.Sprintf("The HTTP body could not be parsed into a valid request. Raw Body: %#v\n, Error: %#v\n", string(b), err2.Error()))
	}
	// Rules for the request: (TODO TESTS)
	// - http.Request content-type == application/json
	// - Node Id always 64 chars long
	// - Port has to exist, and > 0
	// - Type cannot be 0
	// - Protocol subprotocols have to include "c0" (aether subprotocol of mim)
	if r.Header["Content-Type"][0] == "application/json" &&
		req.Address.Port > 0 &&
		req.Address.Type != 0 {
		for _, ext := range req.Address.Protocol.Subprotocols {
			if ext.Name == "c0" {
				// We insert to the POST request the locally sourced details. (Location, Sublocation, LocationType [ipv4 or 6], LastOnline)
				err := insertLocallySourcedRemoteAddressDetails(r, &req)
				if err != nil {
					return req, err
				}
				return req, nil
			}
		}
	}
	return req, errors.New(fmt.Sprintf("The request is syntactically valid JSON, but it does not include certain vital information"))
}

func NodePOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("node", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func BoardsPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("boards", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func ThreadsPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("threads", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func PostsPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("posts", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func VotesPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("votes", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func AddressesPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("addresses", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func KeysPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("keys", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}

func TruststatesPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	err2 := SaveRemote(req)
	if err2 != nil {
		return []byte{}, err2
	}
	respAsByte, err3 := responsegenerator.GeneratePOSTResponse("truststates", req)
	if err3 != nil {
		return respAsByte, err3
	}
	return respAsByte, nil
}
