// Backend > Server
// This file provides the backend the server that the external world accesses to get data from the backend.

package server

import (
	// "aether-core/backend/dispatch"
	"aether-core/backend/responsegenerator"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/toolbox"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NYTimes/gziphandler"
	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net"
	"net/http"
	// "strconv"
	// "bufio"
	"crypto/tls"
	// "github.com/libp2p/go-reuseport"
	// "reflect"
	"path/filepath"
	"strings"
	"time"
)

func isReverseConn(host string, port uint16) bool {
	return host == globals.BackendTransientConfig.ReverseConnData.C1LocalLocalAddr && port == globals.BackendTransientConfig.ReverseConnData.C1LocalLocalPort
}

// Bouncer gate
func isAllowedByBouncer(r *http.Request) bool {
	remoteHost, remotePort := toolbox.SplitHostPort(r.RemoteAddr)
	reverse := isReverseConn(remoteHost, remotePort)
	return globals.BackendTransientConfig.Bouncer.RequestInboundLease(remoteHost, "", remotePort, reverse)
}

// Node type gate
func isAllowedByNodeType(method string) bool {
	nt := globals.BackendConfig.GetNodeType()
	switch method {
	case "GET":
		switch nt {
		default:
			return true
		}
	case "POST":
		switch nt {
		case 2:
			return true
		case 3:
			return true
		default:
			return false
		}
	default:
		return false
	}
	return false
}

// Server responds to GETs with the caches and to POSTS with the live data from the database.
func StartMimServer() {
	protv := globals.BackendConfig.GetProtURLVersion()
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// // START SIMULATE NAT
		// // simulate nat. this works because both apps tend to get ports in +1 -1 of the range of themselves. only accept from internal call.
		// host, port := toolbox.SplitHostPort(r.RemoteAddr)
		// if !isReverseConn(host, port) {
		// 	w.WriteHeader(http.StatusForbidden)
		// 	return
		// }
		// // END SIMULATE NAT
		if !isAllowedByNodeType(r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// We do not gate POST response directory (this handler). Because the only wy somebody would find this would be that it would have hit a POST endpoint. Only because that initial request was allowed the remote could find this link, so the remote already has had a lease and made a request.

		// if !isAllowedByBouncer(r) {
		// 	w.WriteHeader(http.StatusTooManyRequests)
		// 	return
		// }
		if r.Method == "GET" { // this is the part that serves multipage post responses.
			// Check with bouncer if this request is allowed. If not, return too busy.
			w.Header().Set("Content-Type", "application/json")
			// Some safeguards. Some of those are replicated in Go's own http library code, but it's still good to have these here just in case.
			// This disallows serving of .dotfiles and directory indexes.
			// Heads up! This will actually serve anything in the directory - if the user actually ends up putting a random file here, it will also get served, too. There's no good way to check whether the file is created by us without opening and attempting to parse the file, unfortunately.
			if strings.Contains(r.URL.Path, "..") ||
				strings.Contains(r.URL.Path, "/.") ||
				strings.Contains(r.URL.Path, "\\.") ||
				strings.HasSuffix(r.URL.Path, "/") {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			dir := filepath.Join(globals.BackendConfig.GetCachesDirectory(), r.URL.Path)
			// logging.Logf(1, "POST response directory reader was called for: %s", dir)
			w2 := CustomRespWriter{ResponseWriter: w}
			http.ServeFile(&w2, r, dir)
		} else { // If not GET we bail.
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// // START SIMULATE NAT
		// // simulate nat. this works because both apps tend to get ports in +1 -1 of the range of themselves. only accept from internal call.
		// host, port := toolbox.SplitHostPort(r.RemoteAddr)
		// if !isReverseConn(host, port) {
		// 	w.WriteHeader(http.StatusForbidden)
		// 	return
		// }
		// // END SIMULATE NAT
		if !isAllowedByNodeType(r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Check with bouncer if this request is allowed. If not, return too busy.
		if !isAllowedByBouncer(r) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		// Force the content type to application/json, so even in the case of malicious file serving, it won't be executed by default.
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/" + protv + "/status", "/" + protv + "/status/":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte{})

			case "/" + protv + "/node", "/" + protv + "/node/":
				// Node GET endpoint returns the node info.
				var resp api.ApiResponse
				resp.Prefill()
				// r := responsegenerator.GeneratePrefilledApiResponse()
				// resp = *r
				resp.Endpoint = "node"
				resp.Entity = "node"
				resp.Timestamp = api.Timestamp(time.Now().Unix())
				signingErr := resp.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
				if signingErr != nil {
					logging.Log(1, fmt.Sprintf("This cache page failed to be page-signed. Error: %#v Page: %#v\n", signingErr, resp))
				}
				jsonResp, err := resp.ToJSON()
				if err != nil {
					logging.Log(1, errors.New(fmt.Sprintf("The response that was prepared to respond to this query failed to convert to JSON. Error: %#v\n", err)))
				}
				if len(jsonResp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(jsonResp)
				}
			// FUTURE: /bootstrappers - we should probably cache this.
			case "/" + protv + "/bootstrappers", "/" + protv + "/bootstrappers/":
				// Shortcut endpoints that returns bootstrap nodes that this particular node knows.
				var resp api.ApiResponse
				resp.Prefill()
				// r := responsegenerator.GeneratePrefilledApiResponse()
				// resp = *r
				resp.Endpoint = "bootstrappers"
				resp.Entity = "addresses"
				// Get 20 addresses of type 3 (live bootstrap type) sorted by most recent localarrival and type 254 (static bootstrap type)
				addrsLiveBootstrappers, err := persistence.ReadAddresses("", "", 0, 0, 0, 20, 0, 3, "limit")
				if err != nil {
					logging.Logf(1, "There was an error when we tried to read live bootstrapper addresses for the /bootstrappers endpoint. Error: %#v", err)
				}
				addrsStaticBootstrappers, err2 := persistence.ReadAddresses("", "", 0, 0, 0, 10, 0, 254, "limit")
				if err2 != nil {
					logging.Logf(1, "There was an error when we tried to read static bootstrapper addresses for the /bootstrappers endpoint. Error: %#v", err)
				}
				resp.ResponseBody.Addresses = append(addrsLiveBootstrappers, addrsStaticBootstrappers...)
				resp.Timestamp = api.Timestamp(time.Now().Unix())
				signingErr := resp.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
				if signingErr != nil {
					logging.Log(1, fmt.Sprintf("This cache page failed to be page-signed. Error: %#v Page: %#v\n", signingErr, resp))
				}
				jsonResp, err := resp.ToJSON()
				if err != nil {
					logging.Log(1, errors.New(fmt.Sprintf("The response that was prepared to respond to this query failed to convert to JSON. Error: %#v\n", err)))
				}
				if len(jsonResp) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte{})
				} else {
					w.Write(jsonResp)
				}
			default: // this is the part that serves caches
				// Some safeguards. Some of those are replicated in Go's own http library code, but it's still good to have these here just in case.
				// This disallows serving of .dotfiles and directory indexes.
				// Heads up! This will actually serve anything in the directory - if the user actually ends up putting a random file here, it will also get served, too. There's no good way to check whether the file is created by us without opening and attempting to parse the file, unfortunately.
				if strings.Contains(r.URL.Path, "..") ||
					strings.Contains(r.URL.Path, "/.") ||
					strings.Contains(r.URL.Path, "\\.") ||
					strings.HasSuffix(r.URL.Path, "/") {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				dir := filepath.Join(globals.BackendConfig.GetCachesDirectory(), r.URL.Path)
				// logging.Logf(1, "GET directory reader was called for: %s", dir)
				w2 := CustomRespWriter{ResponseWriter: w}
				http.ServeFile(&w2, r, dir)
			}
		} else if r.Method == "POST" {
			switch r.URL.Path {
			case "/" + protv + "/node", "/" + protv + "/node/":
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

			case "/" + protv + "/c0/boards", "/" + protv + "/c0/boards/":
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

			case "/" + protv + "/c0/threads", "/" + protv + "/c0/threads/":
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

			case "/" + protv + "/c0/posts", "/" + protv + "/c0/posts/":
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

			case "/" + protv + "/c0/votes", "/" + protv + "/c0/votes/":
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

			case "/" + protv + "/c0/keys", "/" + protv + "/c0/keys/":
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

			case "/" + protv + "/c0/truststates", "/" + protv + "/c0/truststates/":
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

			case "/" + protv + "/addresses", "/" + protv + "/addresses/":
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

			default:
				logging.Log(1, fmt.Sprintf("A remote reached out to this node with a request that this node does not have a route for. The requested route: %s, The node requesting: %v", r.URL.Path, r.Body))
				// w.WriteHeader(http.StatusNotFound)
				w.WriteHeader(http.StatusNoContent)
			}
		} else { // If not GET or POST, we bail.
			// w.WriteHeader(http.StatusNotFound)
			w.WriteHeader(http.StatusNoContent)
		}
	})

	gzippedMainHandler := gziphandler.GzipHandler(mainHandler)
	http.Handle("/", gzippedMainHandler)

	gzippedHandler := gziphandler.GzipHandler(handlerFunc)
	http.Handle("/"+protv+"/responses/", gzippedHandler)

	port := globals.BackendConfig.GetExternalPort()
	extIp := globals.BackendConfig.GetExternalIp()
	logging.Log(1, fmt.Sprintf("Serving setup complete. Starting to serve Mim publicly on port %d", port))
	srv := &http.Server{
		Addr:         fmt.Sprint(extIp, ":", port),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)), // Disables HTTP2 because HTTP2 doesn't support Hijack, which we need to use to access the underlying TCP connection to perform a reverse open that we need to access remote nodes behind uncooperating NATs.
		// ConnState:    ConnStateListener,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	// srv.SetKeepAlivesEnabled(true)
	if globals.BackendTransientConfig.TLSEnabled {
		// certLoc := fmt.Sprintf("%s/backend/tls/cert.pem", )
		certLoc := filepath.Join(globals.BackendConfig.GetUserDirectory(), "backend", "tls", "cert.pem")
		// keyLoc := fmt.Sprintf("%s/backend/tls/key.pub", globals.BackendConfig.GetUserDirectory())
		keyLoc := filepath.Join(globals.BackendConfig.GetUserDirectory(), "backend", "tls", "key.pub")
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				// Secure or die
			},
		}
		srv.TLSConfig = tlsConfig
		// HSTS header is not set because node IP addresses are dynamic, and us setting HSTS for an address might mean the next user of that IP address might end up having trouble getting people to connect to it through non-TLS.

		l, err := net.Listen("tcp4", fmt.Sprint(":", port))
		// l, err := reuseport.Listen("tcp4", fmt.Sprint(extIp, ":", port))
		if err != nil {
			logging.LogCrash(err)
		}
		il := &InspectingListener{l}
		srvErr := srv.ServeTLS(il, certLoc, keyLoc)
		if srvErr != nil {
			logging.LogCrash(fmt.Sprintf("Server encountered a fatal error. (Heads up, server also exits with error even when it quits normally) Error: %s", srvErr))
		}
	} else {
		l, err := net.Listen("tcp4", fmt.Sprint(":", port))
		// l, err := reuseport.Listen("tcp4", fmt.Sprint(extIp, ":", port))
		if err != nil {
			logging.LogCrash(err)
		}
		il := &InspectingListener{l}
		srvErr := srv.Serve(il)
		if srvErr != nil {
			logging.LogCrash(fmt.Sprintf("Server encountered a fatal error. (Heads up, server also exits with error even when it quits normally) Error: %s", err))
		}
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

// --->N INBOUND.
// insertLocallySourcedRemoteAddressDetails Inserts the locally sourced data about the remote into the address entity that is coming with the POST request.
func insertLocallySourcedRemoteAddressDetails(r *http.Request, req *api.ApiResponse) error {
	// This runs when a node connects to you.
	// LITTLE-TRUSTED ADDRESS ENTRY
	// Data to keep: Location, Sublocation, Port, LastSuccessfulPing (sublocation is guaranteed to be empty since the connection is coming from an IP, not a static IP)
	// Delete everything else, they're untrustable.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("The address from which the remote is connecting could not be parsed. Remote Address: %s, Error: %s", r.RemoteAddr, err))
	}
	if len(host) == 0 {
		return errors.New(fmt.Sprintf("The address from which the remote is connecting seems to be empty. Remote Address: %#v. %#v", r.RemoteAddr, err))
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
	req.Address.LastSuccessfulPing = api.Timestamp(time.Now().Unix())
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
	// Rules for the request:
	// - http.Request content-type == application/json
	// - Node Id always 64 chars long
	// - Port has to exist, and > 0
	// - Type cannot be 0
	// - Protocol subprotocols have to include "c0" (aether subprotocol of mim)
	// - Has a valid nonce (by proxy, the timestamp is within our allowed clock skew bracket)
	// - PoW is verified.
	if r.Header["Content-Type"][0] == "application/json" &&
		req.Address.Port > 0 &&
		req.Address.Type != 0 &&
		req.VerifyNonce() {
		// Verify remote software type and version and make sure we can negotiate with it.
		if !verifyRemoteClient(req.Address.Client) {
			logging.Logf(1, "This ApiResponse is created by a remote client we do not support. Client: %#v", req.Address.Client)
			return req, errors.New(fmt.Sprintf("This ApiResponse is created by a remote client we do not support. Client: %#v", req.Address.Client))
		}

		// Check PoW, since this is a POST request, it is required to have a PoW.
		valid, err := req.VerifyPoW()
		if !valid || err != nil {
			logging.Logf(1, "This ApiResponse failed PoW verification. Possible error: %v", err)
			return req, errors.New(fmt.Sprintf("This ApiResponse failed PoW verification. Possible error: %v", err))
		}
		for _, ext := range req.Address.Protocol.Subprotocols {
			if ext.Name == "c0" {
				// We insert to the POST request the locally sourced details. (Location, Sublocation, LocationType [ipv4 or 6], LastSuccessfulPing)
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

func verifyRemoteClient(cl api.Client) bool {
	return true // DEBUG: remove
	// List known supported clients here.
	if cl.ClientName == "Aether" {
		return true
	}
	return false
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
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
	if r != nil {
		r.Body.Close()
	}
	return respAsByte, nil
}
