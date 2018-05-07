// Backend > ResponseGenerator > MainEntityGenerate
// This file provides a set of functions that help with generation and structuralisation of main entities after they are fetched from the database in preparation to sending over.

package responsegenerator

import (
	// "fmt"
	"aether-core/io/api"
	// "aether-core/io/persistence"
	// "aether-core/services/configstore"
	"aether-core/services/globals"
	"aether-core/services/logging"
	// "aether-core/services/randomhashgen"
	// "aether-core/services/syncconfirmations"
	"aether-core/services/toolbox"
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "os"
	// "strconv"
	// "strings"
	"time"
)

func bakeEntityPages(resultPages *[]api.ApiResponse, entityCounts *[]api.EntityCount, filters *[]api.Filter, foldername string, isPOST bool, respType string, entityType string) {
	var responsedir string
	if isPOST {
		responsedir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/responses/", foldername)
	} else {
		if respType == "addresses" {
			responsedir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/", respType, "/", foldername)
		} else {
			responsedir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/c0/", respType, "/", foldername)
		}
	}
	// responsedir := fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/responses/", foldername)
	toolbox.CreatePath(responsedir)
	for i, _ := range *resultPages {
		// entityType := findEntityInApiResponse((*resultPages)[i], entityType)
		// Set timestamp, number of items in it, total page count, and which page, filters.
		if filters != nil {
			(*resultPages)[i].Filters = *filters
		}
		(*resultPages)[i].Caching.EntityCounts = *entityCounts
		(*resultPages)[i].Pagination.Pages = uint64(len(*resultPages))
		(*resultPages)[i].Pagination.CurrentPage = uint64(i)
		(*resultPages)[i].Timestamp = api.Timestamp(time.Now().Unix())
		(*resultPages)[i].Entity = entityType
		(*resultPages)[i].Endpoint = entityType
		if isPOST {
			(*resultPages)[i].Endpoint = fmt.Sprint((*resultPages)[i].Endpoint, "_post")
		}
		// Sign
		signingErr := (*resultPages)[i].CreateSignature(globals.BackendConfig.GetBackendKeyPair())
		if signingErr != nil {
			logging.Log(1, fmt.Sprintf("This result page of a multiple-page post response failed to be page-signed. Error: %#v Page: %#v\n", signingErr, (*resultPages)[i]))
		}
		// Convert to JSON
		jsonResp, err := (*resultPages)[i].ToJSON()
		if err != nil {
			logging.Log(1, fmt.Sprintf("This page of a multiple-page post response failed to convert to JSON. Error: %#v\n, Request Body: %#v\n", err, (*resultPages)[i]))
		}
		// Save to disk
		name := fmt.Sprint(i, ".json")
		saveFileToDisk(jsonResp, responsedir, name)
	}
	if isPOST {
		insertIntoPOSTResponseReuseTracker(&(*resultPages)[0], foldername)
	}
}
