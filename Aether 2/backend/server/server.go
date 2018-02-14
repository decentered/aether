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
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Server responds to GETs with the caches and to POSTS with the live data from the database.
func Serve() {
	http.HandleFunc("/responses/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			dir := fmt.Sprint(globals.UserDirectory, "/statics", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, dir)
		} else { // If not GET we bail.
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Force the content type to application/json, so even in the case of malicious file serving, it won't be executed by default.
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			switch r.URL.Path {

			case "/v0/status", "/v0/status/":
				// Status GET endpoint returns HTTP 200 only if the node is up, and 429 Too Many Requests if the node is being overloaded.
				if globals.TooManyConnections {
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
				http.ServeFile(w, r, fmt.Sprint(globals.UserDirectory, "/statics/caches", r.URL.Path))
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

			case "/v0/c0/addresses", "/v0/c0/addresses/":
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
				w.WriteHeader(http.StatusNotFound)
			}
		} else { // If not GET or POST, we bail.
			w.WriteHeader(http.StatusNotFound)
		}
	})
	logging.Log(1, "Serving setup complete. Starting to serve publicly.")
	http.ListenAndServe(fmt.Sprint("127.0.0.1", ":", 8089), nil)
}

// MaybeSaveRemote checks if the database has data about the remote that is reaching out. If not, save a new address.
func MaybeSaveRemote(req api.ApiResponse) {
	// We don't insert the node, only the address. Because the remote data is untrustable.
	persistence.InsertOrUpdateAddress(req.Address)
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
		return errors.New(fmt.Sprintf("The address from which the remote is connecting seems to be empty. Remote Address: %s", r.RemoteAddr, err))
	}
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
	req.Address.Protocol.Subprotocols = []api.Subprotocol{}
	req.Address.Protocol.VersionMajor = 0
	req.Address.Protocol.VersionMinor = 0
	req.Address.Client.ClientName = ""
	req.Address.Client.VersionMajor = 0
	req.Address.Client.VersionMinor = 0
	req.Address.Client.VersionPatch = 0
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
		len(req.NodeId) == 64 &&
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
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("node", req)
	return respAsByte, err
}

func BoardsPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("boards", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}

func ThreadsPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("threads", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}

func PostsPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("posts", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}

func VotesPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("votes", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}

func AddressesPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("addresses", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}

func KeysPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("keys", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}

func TruststatesPOST(r *http.Request) ([]byte, error) {
	req, err := ParsePOSTRequest(r)
	if err != nil {
		logging.Log(1, fmt.Sprintf("POST request parsing failed. Error: %#v\n, Request Header: %#v\n, Request Body: %#v\n", err, r.Header, req))
		return []byte{}, nil
	}
	MaybeSaveRemote(req)
	respAsByte, err := responsegenerator.GeneratePOSTResponse("truststates", req)
	if err != nil {
		return respAsByte, err
	}
	return respAsByte, nil
}
