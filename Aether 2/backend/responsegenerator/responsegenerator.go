// Backend > ResponseGenerator
// This file provides a set of functions that take a database response, and convert it into a set of paginated (or nonpaginated) results.

package responsegenerator

import (
	// "fmt"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/randomhashgen"
	"aether-core/services/syncconfirmations"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// GeneratePrefilledApiResponse constructs the basic ApiResponse and fills it with the data about the local machine.
func GeneratePrefilledApiResponse() *api.ApiResponse {
	subprotsAsShims := globals.BackendConfig.GetServingSubprotocols()
	subprotsSupported := []api.Subprotocol{}
	for _, val := range subprotsAsShims {
		subprotsSupported = append(subprotsSupported, api.Subprotocol(val))
	}
	var resp api.ApiResponse
	resp.NodeId = api.Fingerprint(globals.BackendConfig.GetNodeId())
	resp.Address.LocationType = globals.BackendConfig.GetExternalIpType()
	resp.Address.Type = 2 // This is a live node.
	resp.Address.Port = uint16(globals.BackendConfig.GetExternalPort())
	resp.Address.Protocol.VersionMajor = globals.BackendConfig.GetProtocolVersionMajor()
	resp.Address.Protocol.VersionMinor = globals.BackendConfig.GetProtocolVersionMinor()
	resp.Address.Protocol.Subprotocols = subprotsSupported
	resp.Address.Client.VersionMajor = globals.BackendConfig.GetClientVersionMajor()
	resp.Address.Client.VersionMinor = globals.BackendConfig.GetClientVersionMinor()
	resp.Address.Client.VersionPatch = globals.BackendConfig.GetClientVersionPatch()
	resp.Address.Client.ClientName = globals.BackendConfig.GetClientName()
	return &resp
}

func ConvertApiResponseToJson(resp *api.ApiResponse) ([]byte, error) {
	result, err := json.Marshal(resp)
	var jsonErr error
	if err != nil {
		jsonErr = errors.New(fmt.Sprint(
			"This ApiResponse failed to convert to JSON. Error: %#v, ApiResponse: %#v", err, *resp))
	}
	return result, jsonErr
}

type FilterSet struct {
	Fingerprints []api.Fingerprint
	TimeStart    api.Timestamp
	TimeEnd      api.Timestamp
	Embeds       []string
}

func processFilters(req *api.ApiResponse) FilterSet {
	var fs FilterSet
	for _, filter := range req.Filters {
		// Fingerprint
		if filter.Type == "fingerprint" {
			for _, fp := range filter.Values {
				fs.Fingerprints = append(fs.Fingerprints, api.Fingerprint(fp))
			}
		}
		// Embeds
		if filter.Type == "embed" {
			for _, embed := range filter.Values {
				fs.Embeds = append(fs.Embeds, embed)
			}
		}
		// If a time filter is given, timeStart is either the timestamp provided by the remote if it's larger than the end date of the last cache, or the end timestamp of the last cache.
		// In essence, we do not provide anything that is already cached from the live server.
		if filter.Type == "timestamp" {
			// now := int64(time.Now().Unix())
			start, _ := strconv.ParseInt(filter.Values[0], 10, 64)
			end, _ := strconv.ParseInt(filter.Values[1], 10, 64)

			// If there is a value given (not 0), that is, the timerange filter is active.
			// The sanitisation of these ranges are done in the DB level, so this is just intake.
			if start > 0 || end > 0 {
				fs.TimeStart = api.Timestamp(start)
				fs.TimeEnd = api.Timestamp(end)
			}

		}
	}
	return fs
}

func splitEntityIndexesToPages(fullData *api.Response) *[]api.Response {
	var entityTypes []string
	if len(fullData.BoardIndexes) > 0 {
		entityTypes = append(entityTypes, "boardindexes")
	}
	if len(fullData.ThreadIndexes) > 0 {
		entityTypes = append(entityTypes, "threadindexes")
	}
	if len(fullData.PostIndexes) > 0 {
		entityTypes = append(entityTypes, "postindexes")
	}
	if len(fullData.VoteIndexes) > 0 {
		entityTypes = append(entityTypes, "voteindexes")
	}
	if len(fullData.AddressIndexes) > 0 {
		entityTypes = append(entityTypes, "addressindexes")
	}
	if len(fullData.KeyIndexes) > 0 {
		entityTypes = append(entityTypes, "keyindexes")
	}
	if len(fullData.TruststateIndexes) > 0 {
		entityTypes = append(entityTypes, "truststateindexes")
	}

	var pages []api.Response
	// This is a lot of copy paste. This is because there is no automatic conversion from []api.Boards being recognised as []api.Provable. Without that, I have to convert them explicitly to be able to put them into a map[string:struct] which is a lot of extra work - more work than copy paste.
	for i, _ := range entityTypes {
		// Index entities
		if entityTypes[i] == "boardindexes" {
			dataSet := fullData.BoardIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().BoardIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.BoardIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "threadindexes" {
			dataSet := fullData.ThreadIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().ThreadIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.ThreadIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "postindexes" {
			dataSet := fullData.PostIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().PostIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.PostIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "voteindexes" {
			dataSet := fullData.VoteIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().VoteIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.VoteIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "keyindexes" {
			dataSet := fullData.KeyIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().KeyIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.KeyIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "addressindexes" {
			dataSet := fullData.AddressIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().AddressIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.AddressIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "truststateindexes" {
			dataSet := fullData.TruststateIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().TruststateIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.TruststateIndexes = pageData
				pages = append(pages, page)
			}
		}
	}
	if len(entityTypes) == 0 {
		// The result is empty
		var page api.Response
		pages = append(pages, page)
	}
	return &pages
}

func splitEntitiesToPages(fullData *api.Response) *[]api.Response {
	var entityTypes []string
	// We do this check set below so that we don't run pagination logic on entity types that does not exist in this response. This is a bit awkward because there's no good way to iterate over fields of a struct.
	if len(fullData.Boards) > 0 {
		entityTypes = append(entityTypes, "boards")
	}
	if len(fullData.BoardIndexes) > 0 {
		entityTypes = append(entityTypes, "boardindexes")
	}
	if len(fullData.Threads) > 0 {
		entityTypes = append(entityTypes, "threads")
	}
	if len(fullData.ThreadIndexes) > 0 {
		entityTypes = append(entityTypes, "threadindexes")
	}
	if len(fullData.Posts) > 0 {
		entityTypes = append(entityTypes, "posts")
	}
	if len(fullData.PostIndexes) > 0 {
		entityTypes = append(entityTypes, "postindexes")
	}
	if len(fullData.Votes) > 0 {
		entityTypes = append(entityTypes, "votes")
	}
	if len(fullData.VoteIndexes) > 0 {
		entityTypes = append(entityTypes, "voteindexes")
	}
	if len(fullData.Addresses) > 0 {
		entityTypes = append(entityTypes, "addresses")
	}
	if len(fullData.AddressIndexes) > 0 {
		entityTypes = append(entityTypes, "addressindexes")
	}
	if len(fullData.Keys) > 0 {
		entityTypes = append(entityTypes, "keys")
	}
	if len(fullData.KeyIndexes) > 0 {
		entityTypes = append(entityTypes, "keyindexes")
	}
	if len(fullData.Truststates) > 0 {
		entityTypes = append(entityTypes, "truststates")
	}
	if len(fullData.TruststateIndexes) > 0 {
		entityTypes = append(entityTypes, "truststateindexes")
	}

	var pages []api.Response
	// This is a lot of copy paste. This is because there is no automatic conversion from []api.Boards being recognised as []api.Provable. Without that, I have to convert them explicitly to be able to put them into a map[string:struct] which is a lot of extra work - more work than copy paste.
	for i, _ := range entityTypes {
		if entityTypes[i] == "boards" {
			dataSet := fullData.Boards
			pageSize := globals.BackendConfig.GetEntityPageSizes().Boards
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Boards = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "threads" {
			dataSet := fullData.Threads
			pageSize := globals.BackendConfig.GetEntityPageSizes().Threads
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Threads = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "posts" {
			dataSet := fullData.Posts
			pageSize := globals.BackendConfig.GetEntityPageSizes().Posts
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Posts = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "votes" {
			dataSet := fullData.Votes
			pageSize := globals.BackendConfig.GetEntityPageSizes().Votes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Votes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "addresses" {
			dataSet := fullData.Addresses
			pageSize := globals.BackendConfig.GetEntityPageSizes().Addresses
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Addresses = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "keys" {
			dataSet := fullData.Keys
			pageSize := globals.BackendConfig.GetEntityPageSizes().Keys
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Keys = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "truststates" {
			dataSet := fullData.Truststates
			pageSize := globals.BackendConfig.GetEntityPageSizes().Truststates
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Truststates = pageData
				pages = append(pages, page)
			}
		}
		// Index entities
		if entityTypes[i] == "boardindexes" {
			dataSet := fullData.BoardIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().BoardIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.BoardIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "threadindexes" {
			dataSet := fullData.ThreadIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().ThreadIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.ThreadIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "postindexes" {
			dataSet := fullData.PostIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().PostIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.PostIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "voteindexes" {
			dataSet := fullData.VoteIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().VoteIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.VoteIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "keyindexes" {
			dataSet := fullData.KeyIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().KeyIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.KeyIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "addressindexes" {
			dataSet := fullData.AddressIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().AddressIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.AddressIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "truststateindexes" {
			dataSet := fullData.TruststateIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().TruststateIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.TruststateIndexes = pageData
				pages = append(pages, page)
			}
		}
	}
	if len(entityTypes) == 0 {
		// The result is empty
		var page api.Response
		pages = append(pages, page)
	}
	return &pages
}

func convertResponsesToApiResponses(r *[]api.Response) *[]api.ApiResponse {
	var responses []api.ApiResponse
	for i, _ := range *r {
		resp := GeneratePrefilledApiResponse()
		resp.ResponseBody.Boards = (*r)[i].Boards
		resp.ResponseBody.Threads = (*r)[i].Threads
		resp.ResponseBody.Posts = (*r)[i].Posts
		resp.ResponseBody.Votes = (*r)[i].Votes
		resp.ResponseBody.Addresses = (*r)[i].Addresses
		resp.ResponseBody.Keys = (*r)[i].Keys
		resp.ResponseBody.Truststates = (*r)[i].Truststates
		// Indexes
		resp.ResponseBody.BoardIndexes = (*r)[i].BoardIndexes
		resp.ResponseBody.ThreadIndexes = (*r)[i].ThreadIndexes
		resp.ResponseBody.PostIndexes = (*r)[i].PostIndexes
		resp.ResponseBody.VoteIndexes = (*r)[i].VoteIndexes
		resp.ResponseBody.AddressIndexes = (*r)[i].AddressIndexes
		resp.ResponseBody.KeyIndexes = (*r)[i].KeyIndexes
		resp.ResponseBody.TruststateIndexes = (*r)[i].TruststateIndexes
		resp.Pagination.Pages = uint64(len(*r) - 1) // pagination starts from 0
		resp.Pagination.CurrentPage = uint64(i)
		responses = append(responses, *resp)
	}
	return &responses
}

// func randomhashgen.GenerateRandomHash() (string, error) {
// 	const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	saltBytes := make([]byte, 16)
// 	for i := range saltBytes {
// 		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
// 		if err != nil {
// 			return "", errors.New(fmt.Sprint(
// 				"Random number generator generated an error. err: ", err))
// 		}
// 		saltBytes[i] = LETTERS[int(randNum.Int64())]
// 	}
// 	calculator := sha256.New()
// 	calculator.Write(saltBytes)
// 	resultHex := fmt.Sprintf("%x", calculator.Sum(nil))
// 	return resultHex, nil
// }

func generateExpiryTimestamp() int64 {
	expiry := time.Duration(globals.BackendConfig.GetPOSTResponseExpiryMinutes()) * time.Minute
	expiryTs := int64(time.Now().Add(expiry).Unix())
	return expiryTs
}

func findEntityInApiResponse(resp api.ApiResponse) string {
	if len(resp.ResponseBody.Boards) > 0 {
		return "boards"
	}
	if len(resp.ResponseBody.Threads) > 0 {
		return "threads"
	}
	if len(resp.ResponseBody.Posts) > 0 {
		return "posts"
	}
	if len(resp.ResponseBody.Votes) > 0 {
		return "votes"
	}
	if len(resp.ResponseBody.Addresses) > 0 {
		return "addresses"
	}
	if len(resp.ResponseBody.Keys) > 0 {
		return "keys"
	}
	if len(resp.ResponseBody.Truststates) > 0 {
		return "truststates"
	}
	return ""
}

func createPath(path string) {
	os.MkdirAll(path, 0755)
}

func saveFileToDisk(fileContents []byte, path string, filename string) {
	ioutil.WriteFile(fmt.Sprint(path, "/", filename), fileContents, 0755)
}

// bakeFinalApiResponse looks at the resultpages. If there is one, it is directly provided as is. If there is more, the results are committed into the file system, and a cachelink page is provided instead.
func bakeFinalApiResponse(resultPages *[]api.ApiResponse) (*api.ApiResponse, error) {
	resp := GeneratePrefilledApiResponse()
	if len(*resultPages) > 1 {
		// Create a random SHA256 hash as folder name
		dirname, err := randomhashgen.GenerateRandomHash()
		if err != nil {
			return resp, err
		}
		// Generate the responses directory if doesn't exist. Add the expiry date to the folder name to be searched for.
		foldername := fmt.Sprint(generateExpiryTimestamp(), "_", dirname)
		responsedir := fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/responses/", foldername)
		os.MkdirAll(responsedir, 0755)
		var jsons [][]byte
		// For each response, number it, set timestamps etc. And save to disk.
		for i, _ := range *resultPages {
			resultPage := (*resultPages)[i]
			entityType := findEntityInApiResponse(resultPage)
			// Set timestamp, number of items in it, total page count, and which page.

			resultPage.Pagination.Pages = uint64(len(*resultPages))
			resultPage.Pagination.CurrentPage = uint64(i)
			resultPage.Timestamp = api.Timestamp(time.Now().Unix())
			resultPage.Entity = entityType
			resultPage.Endpoint = fmt.Sprint(entityType, "_post")
			jsonResp, err := ConvertApiResponseToJson(&resultPage)
			if err != nil {
				logging.Log(1, fmt.Sprintf("This page of a multiple-page post response failed to convert to JSON. Error: %#v\n, Request Body: %#v\n", err, resultPage))
			}
			jsons = append(jsons, jsonResp)
		}
		// Insert these jsons into the filesystem.
		for i, _ := range jsons {
			name := fmt.Sprint(i, ".json")
			createPath(responsedir)
			saveFileToDisk(jsons[i], responsedir, name)
			var c api.ResultCache
			c.ResponseUrl = foldername
			resp.Results = append(resp.Results, c)
		}
		resp.Endpoint = "multipart_post_response"

	} else if len(*resultPages) == 1 {
		// There is only one response page here.
		entityType := findEntityInApiResponse((*resultPages)[0])
		resp.Pagination.Pages = 0 // These start to count from 0
		resp.Pagination.CurrentPage = 0
		resp.Entity = entityType
		resp.Endpoint = "singular_post_response"
		resp.ResponseBody = (*resultPages)[0].ResponseBody
	} else {
		logging.LogCrash(fmt.Sprintf("This post request produced both no results and no resulting apiResponses. []ApiResponse: %#v", *resultPages))
	}
	return resp, nil
}

// sanitiseOutboundAddresses removes untrusted address data from the addresses destined to go out of this node. The remote node will also remove it, but there is no reason to leak information unnecessarily.
func sanitiseOutboundAddresses(addrsPtr *[]api.Address) *[]api.Address {
	addrs := *addrsPtr
	for key, _ := range addrs {
		addrs[key].LocationType = 0
		addrs[key].Type = 0
		addrs[key].LastOnline = 0
		addrs[key].Protocol = api.Protocol{}
		addrs[key].Protocol.Subprotocols = []api.Subprotocol{}
		addrs[key].Client = api.Client{}
	}
	return &addrs
}

// GeneratePOSTResponse creates a response that is directly returned to a custom request by the remote.
func GeneratePOSTResponse(respType string, req api.ApiResponse) ([]byte, error) {
	var resp api.ApiResponse
	// Look at filters to figure out what is being requested
	filters := processFilters(&req)
	switch respType {
	case "node":
		r := GeneratePrefilledApiResponse()
		resp = *r
		// resp.Endpoint = "node"
		resp.Entity = "node"
	case "boards", "threads", "posts", "votes", "keys", "truststates":
		localData, dbError := persistence.Read(respType, filters.Fingerprints, filters.Embeds, filters.TimeStart, filters.TimeEnd)
		if dbError != nil {
			return []byte{}, errors.New(fmt.Sprintf("The query coming from the remote caused an error in the local database while trying to respond to this request. Error: %#v\n, Request: %#v\n", dbError, req))
		}
		pages := splitEntitiesToPages(&localData)
		pagesAsApiResponses := convertResponsesToApiResponses(pages)
		finalResponse, err := bakeFinalApiResponse(pagesAsApiResponses)
		// fmt.Printf("%#v", finalResponse)
		if err != nil {
			return []byte{}, errors.New(fmt.Sprintf("An error was encountered while trying to finalise the API response. Error: %#v\n, Request: %#v\n", err, req))
		}
		resp = *finalResponse
		// resp.Endpoint = "entity"
	case "addresses": // Addresses can't do address search by loc/subloc/port. Only time search is available, since addresses don't have fingerprints defined.
		/*
			An addresses POST response returns results within the time boundary that has been seen online first-person by the remote. It does not communicate addresses that the remote has not connected to.
		*/
		logging.Log(2, fmt.Sprintf("We've gotten an address request with the filters: %#v", filters))
		addresses, dbError := persistence.ReadAddresses("", "", 0, filters.TimeStart, filters.TimeEnd, 0, 0, 0, "connected")
		addresses = *sanitiseOutboundAddresses(&addresses)
		var localData api.Response
		localData.Addresses = addresses
		if dbError != nil {
			return []byte{}, errors.New(fmt.Sprintf("The query coming from the remote caused an error in the local database while trying to respond to this request. Error: %#v\n, Request: %#v\n", dbError, req))
		}
		pages := splitEntitiesToPages(&localData)
		pagesAsApiResponses := convertResponsesToApiResponses(pages)
		finalResponse, err := bakeFinalApiResponse(pagesAsApiResponses)
		if err != nil {
			return []byte{}, errors.New(fmt.Sprintf("An error was encountered while trying to finalise the API response. Error: %#v\n, Request: %#v\n", err, req))
		}
		resp = *finalResponse
		resp.Endpoint = "entity"
	}
	// Build the response itself
	resp.Entity = respType
	resp.Timestamp = api.Timestamp(time.Now().Unix())
	// Construct the query, and run an index to determine how many entries we have for the filter.
	jsonResp, err := ConvertApiResponseToJson(&resp)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("The response that was prepared to respond to this query failed to convert to JSON. Error: %#v\n, Request Body: %#v\n", err, req))
	}
	return jsonResp, nil
}

func createBoardIndex(entity *api.Board, pageNum int) api.BoardIndex {
	var entityIndex api.BoardIndex
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	return entityIndex
}

func createThreadIndex(entity *api.Thread, pageNum int) api.ThreadIndex {
	var entityIndex api.ThreadIndex
	entityIndex.Board = entity.Board
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.PageNumber = pageNum
	return entityIndex
}

func createPostIndex(entity *api.Post, pageNum int) api.PostIndex {
	var entityIndex api.PostIndex
	entityIndex.Board = entity.Board
	entityIndex.Thread = entity.Thread
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.PageNumber = pageNum
	return entityIndex
}

func createVoteIndex(entity *api.Vote, pageNum int) api.VoteIndex {
	var entityIndex api.VoteIndex
	entityIndex.Board = entity.Board
	entityIndex.Thread = entity.Thread
	entityIndex.Target = entity.Target
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	return entityIndex
}

func createKeyIndex(entity *api.Key, pageNum int) api.KeyIndex {
	var entityIndex api.KeyIndex
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	return entityIndex
}

func createTruststateIndex(entity *api.Truststate, pageNum int) api.TruststateIndex {
	var entityIndex api.TruststateIndex
	entityIndex.Target = entity.Target
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	return entityIndex
}

// createIndexes creates the index variant of every entity in an api.Response, and puts it back inside one single container for all indexes.
func createIndexes(fullData *[]api.Response) *api.Response {
	fd := *fullData
	var resp api.Response
	if len(fd) > 0 {
		for i, _ := range fd {
			// For each Api.Response page
			if len(fd[i].Boards) > 0 {
				for j, _ := range fd[i].Boards {
					entityIndex := createBoardIndex(&fd[i].Boards[j], i)
					resp.BoardIndexes = append(resp.BoardIndexes, entityIndex)
				}
			}
			if len(fd[i].Threads) > 0 {
				for j, _ := range fd[i].Threads {
					entityIndex := createThreadIndex(&fd[i].Threads[j], i)
					resp.ThreadIndexes = append(resp.ThreadIndexes, entityIndex)
				}
			}
			if len(fd[i].Posts) > 0 {
				for j, _ := range fd[i].Posts {
					entityIndex := createPostIndex(&fd[i].Posts[j], i)
					resp.PostIndexes = append(resp.PostIndexes, entityIndex)
				}
			}
			if len(fd[i].Votes) > 0 {
				for j, _ := range fd[i].Votes {
					entityIndex := createVoteIndex(&fd[i].Votes[j], i)
					resp.VoteIndexes = append(resp.VoteIndexes, entityIndex)
				}
			}
			// Addresses: Address doesn't have an index form. It is its own index.
			// Addresses are skipped here.
			if len(fd[i].Keys) > 0 {
				for j, _ := range fd[i].Keys {
					entityIndex := createKeyIndex(&fd[i].Keys[j], i)
					resp.KeyIndexes = append(resp.KeyIndexes, entityIndex)
				}
			}
			if len(fd[i].Truststates) > 0 {
				for j, _ := range fd[i].Truststates {
					entityIndex := createTruststateIndex(&fd[i].Truststates[j], i)
					resp.TruststateIndexes = append(resp.TruststateIndexes, entityIndex)
				}
			}
		}
	}
	return &resp
}

func generateCacheName() (string, error) {
	hash, err := randomhashgen.GenerateRandomHash()
	if err != nil {
		return "", err
	}
	n := fmt.Sprint("cache_", hash)
	return n, nil
}

// CacheResponse is the internal procesing structure for generating caches to be saved to the disk.
type CacheResponse struct {
	cacheName   string
	start       api.Timestamp
	end         api.Timestamp
	entityPages *[]api.Response
	indexPages  *[]api.Response
}

func cleanTooOldEntities(localData *api.Response) *api.Response {
	networkHeadEndTs := api.Timestamp(time.Now().Add(-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
	var cleanedLocalData api.Response
	if len(localData.Boards) > 0 {
		for _, val := range localData.Boards {
			if val.Creation > networkHeadEndTs || val.LastUpdate > networkHeadEndTs {
				cleanedLocalData.Boards = append(cleanedLocalData.Boards, val)
			} else {
				logging.Log(1, fmt.Sprintf("This entity didn't make the cut for being included in any caches at the point of creation. Entity: %#v", val))
			}
		}
	} else if len(localData.Threads) > 0 {
		for _, val := range localData.Threads {
			if val.Creation > networkHeadEndTs {
				cleanedLocalData.Threads = append(cleanedLocalData.Threads, val)
			} else {
				logging.Log(1, fmt.Sprintf("This entity didn't make the cut for being included in any caches at the point of creation. Entity: %#v", val))
			}
		}
	} else if len(localData.Posts) > 0 {
		for _, val := range localData.Posts {
			if val.Creation > networkHeadEndTs {
				cleanedLocalData.Posts = append(cleanedLocalData.Posts, val)
			} else {
				logging.Log(1, fmt.Sprintf("This entity didn't make the cut for being included in any caches at the point of creation. Entity: %#v", val))
			}
		}
	} else if len(localData.Votes) > 0 {
		for _, val := range localData.Votes {
			if val.Creation > networkHeadEndTs || val.LastUpdate > networkHeadEndTs {
				cleanedLocalData.Votes = append(cleanedLocalData.Votes, val)
			} else {
				logging.Log(1, fmt.Sprintf("This entity didn't make the cut for being included in any caches at the point of creation. Entity: %#v", val))
			}
		}
	} else if len(localData.Keys) > 0 {
		for _, val := range localData.Keys {
			if val.Creation > networkHeadEndTs || val.LastUpdate > networkHeadEndTs {
				cleanedLocalData.Keys = append(cleanedLocalData.Keys, val)
			} else {
				logging.Log(1, fmt.Sprintf("This entity didn't make the cut for being included in any caches at the point of creation. Entity: %#v", val))
			}
		}
	} else if len(localData.Truststates) > 0 {
		for _, val := range localData.Truststates {
			if val.Creation > networkHeadEndTs || val.LastUpdate > networkHeadEndTs {
				cleanedLocalData.Truststates = append(cleanedLocalData.Truststates, val)
			} else {
				logging.Log(1, fmt.Sprintf("This entity didn't make the cut for being included in any caches at the point of creation. Entity: %#v", val))
			}
		}
	}
	return &cleanedLocalData
}

// GenerateCacheResponse responds to a cache generation request. This returns an Api.Response entity with entities, entity indexes, and the cache link that needs to be inserted into the index of the endpoint.
// This has no filters.
func GenerateCacheResponse(respType string, start api.Timestamp, end api.Timestamp) (CacheResponse, error) {
	var resp CacheResponse
	switch respType {
	case "boards", "threads", "posts", "votes", "keys", "truststates":
		localData, dbError := persistence.Read(respType, []api.Fingerprint{}, []string{}, start, end)
		if dbError != nil {
			return resp, errors.New(fmt.Sprintf("This cache generation request caused an error in the local database while trying to respond to this request. Error: %#v\n", dbError))
		}
		if len(localData.Boards) == 0 && len(localData.Threads) == 0 && len(localData.Posts) == 0 && len(localData.Votes) == 0 && len(localData.Keys) == 0 && len(localData.Truststates) == 0 {
			/*
				There's no data in this result. But the cache generation should continue. Why?

				1) This cache generation process is guarded by the 'is this node tracking head?' guard. So this part of the code does not need to care about accidentally generating blank caches.

				2) Consider the case that the most recent data in the network is actually, genuinely three days old. Had we stopped cache generation when empty, the caches for those two blank days would NEVER be generated, but ALWAYS attempted. So every hit of the cache generation cycle would turn out to be an attempt for the cache generation of those two days.

				How do I know? Because that's exactly what happened and this text is the bug fix.
			*/
			logging.Log(1, fmt.Sprintf("The result for this cache is empty. Entity type: %s", respType))
		}
		cleanedLocalData := cleanTooOldEntities(&localData)
		localData = *cleanedLocalData
		entityPages := splitEntitiesToPages(&localData)
		indexes := createIndexes(entityPages)
		indexPages := splitEntityIndexesToPages(indexes)
		cn, err := generateCacheName()
		if err != nil {
			return resp, errors.New(fmt.Sprintf("There was an error in the cache generation request serving. Error: %#v\n", err))
		}
		resp.cacheName = cn
		resp.start = start
		resp.end = end
		resp.indexPages = indexPages
		resp.entityPages = entityPages

	case "addresses":
		addresses, dbError := persistence.ReadAddresses("", "", 0, start, end, 0, 0, 0, "connected") // Cache generation only generates caches for addresses that this computer has personally connected to.
		if dbError != nil {
			return resp, errors.New(fmt.Sprintf("This cache generation request caused an error in the local database while trying to respond to this request. Error: %#v\n", dbError))
		}
		addresses = *sanitiseOutboundAddresses(&addresses)
		if len(addresses) == 0 {
			/*
				There's no data in this result. But the cache generation should continue. Why?

				1) This cache generation process is guarded by the 'is this node tracking head?' guard. So this part of the code does not need to care about accidentally generating blank caches.

				2) Consider the case that the most recent data in the network is actually, genuinely three days old. Had we stopped cache generation when empty, the caches for those two blank days would NEVER be generated, but ALWAYS attempted. So every hit of the cache generation cycle would turn out to be an attempt for the cache generation of those two days.

				How do I know? Because that's exactly what happened and this text is the bug fix.
			*/
			logging.Log(1, fmt.Sprintf("The result for this cache is empty. Entity type: %s", respType))
		}
		var localData api.Response
		localData.Addresses = addresses

		entityPages := splitEntitiesToPages(&localData)
		cn, err := generateCacheName()
		if err != nil {
			return resp, errors.New(fmt.Sprintf("There was an error in the cache generation request serving. Error: %#v\n", err))
		}
		resp.cacheName = cn
		resp.start = start
		resp.end = end
		resp.entityPages = entityPages

	default:
		return resp, errors.New(fmt.Sprintf("The requested entity type is unknown to the cache generator. Entity type: %s", respType))
	}
	return resp, nil
}

func updateCacheIndex(cacheIndex *api.ApiResponse, cacheData *CacheResponse) {
	// Save the cache link into the index.
	var c api.ResultCache
	c.ResponseUrl = cacheData.cacheName
	c.StartsFrom = cacheData.start
	c.EndsAt = cacheData.end
	cacheIndex.Results = append(cacheIndex.Results, c)
	cacheIndex.Timestamp = api.Timestamp(int64(time.Now().Unix()))
	cacheIndex.Caching.ServedFromCache = true
	cacheIndex.Caching.CacheScope = "day"
	// TODO: How many places am I setting this ".Caching" data?
}

func deleteTooOldCaches(respType string, cacheIndex *api.ApiResponse, entityCacheDir string) {
	var threshold api.Timestamp
	// First, count the number of caches available.
	if respType == "boards" {
		threshold = api.Timestamp(time.Now().Add(
			-time.Duration(globals.BackendConfig.GetNetworkMemoryDays()*24) * time.Hour).Unix())
	} else {
		threshold = api.Timestamp(time.Now().Add(
			-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
	}
	oldestCacheEnd := api.Timestamp(time.Now().Unix())
	for _, cache := range cacheIndex.Results {
		if oldestCacheEnd > cache.EndsAt {
			oldestCacheEnd = cache.EndsAt
		}
	}
	if threshold > oldestCacheEnd {
		// We have more caches than needed. We need to delete some starting from the oldest.
		logging.Log(1, fmt.Sprintf("We have caches for a longer duration of time than we need. (The oldest cache.EndsAt is %d, the threshold is %d) Caches will be purged starting from the oldest. Purge is starting.", oldestCacheEnd, threshold))
		oldCaches := []api.ResultCache{}
		stillValidCaches := []api.ResultCache{}
		for _, cache := range cacheIndex.Results {
			if cache.EndsAt < threshold {
				// This cache is too old.
				oldCaches = append(oldCaches, cache)
			} else {
				// This is not too old. We keep this.
				stillValidCaches = append(stillValidCaches, cache)
			}
		}
		cacheIndex.Results = stillValidCaches
		for _, cache := range oldCaches {
			// Figure out the location of the cache and nuke it.
			location := entityCacheDir + "/" + cache.ResponseUrl
			os.RemoveAll(location)
		}
		logging.Log(1, fmt.Sprintf("Old cache purging is complete. We've deleted these caches from both index and from the local file system: %#v", oldCaches))
	}
}

// saveCacheToDisk saves an entire cache's data (entities and indexes, inside a folder named based on the cache name) into the proper location on the disk.
func saveCacheToDisk(entityCacheDir string, cacheData *CacheResponse, respType string) error {
	// Create the index directory.
	cacheDir := fmt.Sprint(entityCacheDir, "/", cacheData.cacheName)
	createPath(cacheDir)
	var indexPages []api.ApiResponse
	var indexDir string
	if respType != "addresses" {
		indexDir = fmt.Sprint(entityCacheDir, "/", cacheData.cacheName, "/index")
		createPath(indexDir)
		indexPages = *convertResponsesToApiResponses(cacheData.indexPages)
	}
	// Convert api.Responses to api.ApiResponses for saving.
	entityPages := *convertResponsesToApiResponses(cacheData.entityPages)
	// Iterate over the data, convert api.ApiResponses to JSON, and save.
	for i, _ := range indexPages {
		indexPages[i].Endpoint = "entity_index"
		indexPages[i].Entity = respType
		indexPages[i].Timestamp = api.Timestamp(int64(time.Now().Unix()))
		indexPages[i].Caching.ServedFromCache = true
		indexPages[i].Caching.CurrentCacheUrl = cacheData.cacheName
		// indexPages[i].Caching.PrevCacheUrl // TODO Pulling this is expensive as heck here. Reconsider the need.
		indexPages[i].Caching.CacheScope = "day"
		// For each index, look at the page number and save the result as that.
		json, _ := ConvertApiResponseToJson(&indexPages[i])
		saveFileToDisk(json, indexDir, fmt.Sprint(indexPages[i].Pagination.CurrentPage, ".json"))
	}
	for i, _ := range entityPages {
		entityPages[i].Endpoint = "entity"
		entityPages[i].Entity = respType
		entityPages[i].Timestamp = api.Timestamp(int64(time.Now().Unix()))
		entityPages[i].Caching.ServedFromCache = true
		entityPages[i].Caching.CurrentCacheUrl = cacheData.cacheName
		// indexPages[i].Caching.PrevCacheUrl // TODO Pulling this is expensive as heck here. Reconsider the need.
		entityPages[i].Caching.CacheScope = "day"
		// For each index, look at the page number and save the result as that.
		json, _ := ConvertApiResponseToJson(&entityPages[i])
		saveFileToDisk(json, cacheDir, fmt.Sprint(entityPages[i].Pagination.CurrentPage, ".json"))
	}
	return nil
}

// CreateCache creates the cache for the given entity type for the given time range.
func CreateCache(respType string, start api.Timestamp, end api.Timestamp) error {
	// - Pull the data from the DB
	// - Look at the cache folder. If there is a cache folder and an index there, save the cache and add to index.
	// - If there is no cache present there, create the index and add it as the first entry.
	fmt.Printf("create cache was asked to generate a cache for the resp type %#vthat ended at the timestamp: %#v\n", respType, end)
	cacheData, err := GenerateCacheResponse(respType, start, end)
	if err != nil {
		if strings.Contains(err.Error(), "The result for this cache is empty") {
			logging.Log(1, errors.New(fmt.Sprintf("The result for this cache is empty. Entity type: %s", respType)))
			return nil
		} else {
			return errors.New(fmt.Sprintf("Cache creation process encountered an error. Error: %s", err))
		}
	}
	var entityCacheDir string
	if respType == "boards" || respType == "threads" || respType == "posts" || respType == "votes" || respType == "keys" || respType == "truststates" {
		entityCacheDir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/c0/", respType)
	} else if respType == "addresses" {
		entityCacheDir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/v0/", respType)
	} else {
		return errors.New(fmt.Sprintf("Unknown response type: %s", respType))
	}

	// Create the caches dir and the appropriate endpoint if does not exist.
	createPath(entityCacheDir)
	// Save the cache to disk.
	err2 := saveCacheToDisk(entityCacheDir, &cacheData, respType)
	// TODO: above needs to add caching tag, entity and endpoint fields, and the current timestamp.
	if err2 != nil {
		return errors.New(fmt.Sprintf("Cache creation process encountered an error. Error: %s", err2))
	}
	var apiResp api.ApiResponse
	// Look for the index.json in it. If it doesn't exist, create.
	cacheIndexAsJson, err3 := ioutil.ReadFile(fmt.Sprint(entityCacheDir, "/index.json"))
	if err3 != nil && strings.Contains(err3.Error(), "no such file or directory") {
		// The index.json of this cache likely doesn't exist. Create one.
		apiResp = *GeneratePrefilledApiResponse()
	} else if err3 != nil {
		return errors.New(fmt.Sprintf("Cache creation process encountered an error. Error: %s", err3))
	} else {
		// err3 is nil
		json.Unmarshal(cacheIndexAsJson, &apiResp)
	}
	// If the file exists, go through with regular processing.
	updateCacheIndex(&apiResp, &cacheData)
	deleteTooOldCaches(respType, &apiResp, entityCacheDir)
	json, err4 := ConvertApiResponseToJson(&apiResp)
	if err4 != nil {
		return err
	}
	saveFileToDisk(json, entityCacheDir, "index.json")
	return nil
}

/*
	Methods related to cache days table generation. Cache Days Table is a table of days with beginning and end timestamps that we feed into the cache generator to generate caches for those days.

	We then feed this cache generation table into our cache generator, and it creates the appropriate folder structure for us.
*/

// readCacheIndex reads the cache index of the requested endpoint from the local drive. This is then used for finding the end timestamp of the last cache generated.
func readCacheIndex(etype string) (api.ApiResponse, error) {
	var cacheDir string
	if etype == "boards" || etype == "threads" || etype == "posts" || etype == "votes" || etype == "keys" || etype == "truststates" {
		cacheDir = globals.BackendConfig.GetCachesDirectory() + "/v0/c0/" + etype
	} else if etype == "addresses" {
		cacheDir = globals.BackendConfig.GetCachesDirectory() + "/v0/" + etype
	}
	cacheIndex := cacheDir + "/index.json"
	dat, err := ioutil.ReadFile(cacheIndex)
	if err != nil {
		return api.ApiResponse{}, err
	}
	var apiresp api.ApiResponse
	err2 := json.Unmarshal([]byte(dat), &apiresp)
	if err2 != nil {
		logging.Log(1, fmt.Sprintf(fmt.Sprintf(
			"The JSON That was the cache index for the entity type is malformed. Entity type: %s, JSON: %s", etype, string([]byte(dat)))))
		// Delete the whole index folder and return 0 to generate new caches.
		os.RemoveAll(cacheDir)
		return api.ApiResponse{}, errors.New("no such file or directory")
	}
	return apiresp, nil
}

// determineLastCacheEnd figures out when was the last cache for this entity type was generated. For each entity, we need to look at the last cache that is generated by the entity and find its end timestamp.
func determineLastCacheEnd(etype string) api.Timestamp {
	cacheIndex, err := readCacheIndex(etype)
	if err != nil {
		// logging.LogCrash(err)
		// TODO add tampered caches gating
		if strings.Contains(err.Error(), "no such file or directory") {
			logging.Log(1, fmt.Sprintf("The cache for this entity type does not exist yet. We'll be generating this from scratch. Entity type: %#v", etype))
			// var blankTs api.Timestamp
			// return blankTs
		} else {
			logging.LogCrash(err)
		}
	}
	// Identify the most recent end timestamp
	var mostRecentExtantCacheEndTs api.Timestamp
	for _, cache := range cacheIndex.Results {
		if cache.EndsAt > mostRecentExtantCacheEndTs {
			mostRecentExtantCacheEndTs = cache.EndsAt
		}
	}
	// If the most recent extant cache end is lesser than our network head, make it network head (so that we don't have to generate caches starting from 1970s)
	networkHeadThreshold := api.Timestamp(time.Now().Add(-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
	networkMemoryThreshold := api.Timestamp(time.Now().Add(-time.Duration(globals.BackendConfig.GetNetworkMemoryDays()*24) * time.Hour).Unix())
	// If entity type is a board, the threshold is network memory.
	if etype == "boards" {
		if mostRecentExtantCacheEndTs < networkMemoryThreshold {
			mostRecentExtantCacheEndTs = networkMemoryThreshold
		}
	} else {
		// If not a board, then the threshold is network head.
		if mostRecentExtantCacheEndTs < networkHeadThreshold {
			mostRecentExtantCacheEndTs = networkHeadThreshold
		}
	}
	// fmt.Printf("Most recent extant cache end for type %#v is %#v.\n", etype, mostRecentExtantCacheEndTs)
	return mostRecentExtantCacheEndTs
}

// generateRequestedCachesTable determines how many caches we need to generate, and at which intervals they need to start and end.
func generateRequestedCachesTable(mostRecentExtantCacheEndTs api.Timestamp) []api.ResultCache {
	// Split the difference of most recent cache end and now into 24H slices.
	now := api.Timestamp(time.Now().Unix())
	var dayTable []api.ResultCache
	currentEndTs := mostRecentExtantCacheEndTs
	// So long as the current end + a day is lesser than timestamp of now, iterate
	for currentEndTs < now {
		// fmt.Println("current end ts smaller than now")
		newEnd := api.Timestamp(time.Unix(int64(currentEndTs), 0).Add(24 * time.Hour).Unix())
		cache := api.ResultCache{
			StartsFrom: currentEndTs,
			EndsAt:     newEnd,
		}
		currentEndTs = newEnd
		dayTable = append(dayTable, cache)
	}
	if len(dayTable) == 0 {
		logging.LogCrash(fmt.Sprintf("Day table length turned out to be zero. Your Time block size and past blocks to check values are invalid. If you have not changed them, please delete the backend configuration file and restart the application."))
	}
	// After this table generation is done, check the last cache bracket (start>end). If the time difference of its start and now() is less than 12 hours, delete the last bracket, and set the n-1th cache bracket's end timestamp to now.
	/*
		e.g.
		IF:
		 Day -3  Day -2  Day -1     Now
		|-------|-------|-------|----=--|
		(more than half day's worth data in last)

		DO:
		 Day -3  Day -2  Day -1     Now
		|-------|-------|-------|----=|
		(move the end to now)

		IF:
		 Day -3  Day -2  Day -1   Now
		|-------|-------|-------|--=----|
		(less than 12 hours of data in last)

		DO:
		 Day -3  Day -2  Day -1  Now
		|-------|-------|---------=|
		(remove the last one and move the end of n-1 to now)
	*/
	lastDTItemEndTs := dayTable[len(dayTable)-1].EndsAt
	halfDayIntoFuture := api.Timestamp(time.Now().Add(12 * time.Hour).Unix())
	if halfDayIntoFuture < lastDTItemEndTs && len(dayTable) > 1 {
		// The last cache covers less than 12 hours (i.e. it captures more than 12 hours of not-happened-yet) AND it's not just one item in the day table (in which case, don't do anything.)
		// Chop the last item off.
		dayTable = dayTable[:len(dayTable)-1]
	}
	// Make the last item of the day table come up to now, not to future.
	dayTable[len(dayTable)-1].EndsAt = api.Timestamp(time.Now().Unix())
	// fmt.Printf("Day table length is: %#v\n", len(dayTable))
	// fmt.Printf("Day table is: %#v", dayTable)
	return dayTable
}

// GenerateCacheSet determines how many caches we will need to create for a given entity type, and generates them.
func GenerateCacheSet(etype string) {
	// Read the end of the last cache, or if there are none, start from the beginning.
	lastCacheEndTs := determineLastCacheEnd(etype)
	// If the lastCacheEndTs is younger than 23 hours, we do nothing. The cache generator cycle will attempt to create a cache every hour, so this is where we gate how often we create caches.

	// If last cache end is more than 23 hours ago
	cachegenThreshold := api.Timestamp(
		time.Now().Add(-time.Duration(globals.BackendConfig.GetCacheGenerationIntervalHours()-1) * time.Hour).Unix())
	// fmt.Println("Cachegen threshold: ", cachegenThreshold)
	// fmt.Println("Last cache end TS: ", lastCacheEndTs)
	if cachegenThreshold > lastCacheEndTs {
		cachesTable := generateRequestedCachesTable(lastCacheEndTs)
		for _, val := range cachesTable {
			CreateCache(etype, val.StartsFrom, val.EndsAt)
		}
	} else {
		logging.Log(1, fmt.Sprintf("Last cache that was created for %s was newer than 23 hours ago. Please wait until after.", etype))
	}
}

// GenerateCaches generates all day caches for all entities and saves them to disk.
func GenerateCaches() {
	entityTypes := []string{"boards", "threads", "posts", "votes", "keys", "truststates", "addresses"}
	nodeIsUpToDate, err := syncconfirmations.NodeIsTrackingHead()
	if err != nil {
		logging.Log(1, fmt.Sprintf("The function that checks whether the local node is up to date returned an error. Because of that, this cache generation cycle is pre-empted. It'll be attempted again in the next interval. Error: %#v", err))
		return // If the node is not up to date, bail
	}
	// If the node IS up to date, generate the caches.
	if nodeIsUpToDate {
		for _, val := range entityTypes {
			GenerateCacheSet(val)
		}
		// We're setting this for the purposes of denying POST requests with a timestamp that is partially or wholly available within our cache bracket. (That is, it's not used to determine where to start generating caches from, we read the actual saved cache for that.)
		globals.BackendConfig.SetLastCacheGenerationTimestamp(time.Now().Unix())
	}
}
