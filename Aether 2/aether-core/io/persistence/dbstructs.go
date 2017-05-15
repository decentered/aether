// Persistence > Structs
// This file contains the struct definitions of the database objects.

package persistence

import (
	"aether-core/io/api"
	"aether-core/services/logging"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Basic properties

type DbUpdateable struct {
	LastUpdate        api.Timestamp   `db:"LastUpdate"`
	UpdateProofOfWork api.ProofOfWork `db:"UpdateProofOfWork"`
	UpdateSignature   api.Signature   `db:"UpdateSignature"`
}

type DbProvable struct {
	Creation    api.Timestamp   `db:"Creation"`
	ProofOfWork api.ProofOfWork `db:"ProofOfWork"`
	Signature   api.Signature   `db:"Signature"`
}

// Subentities

type DbBoardOwner struct {
	BoardFingerprint api.Fingerprint `db:"BoardFingerprint"`
	KeyFingerprint   api.Fingerprint `db:"KeyFingerprint"`
	Expiry           api.Timestamp   `db:"Expiry"`
	Level            uint8           `db:"Level"`
}

type DbCurrencyAddress struct {
	KeyFingerprint api.Fingerprint `db:"KeyFingerprint"`
	CurrencyCode   string          `db:"CurrencyCode"`
	Address        string          `db:"Address"`
}

// Entities

type DbBoard struct {
	Fingerprint  api.Fingerprint `db:"Fingerprint"`
	Name         string          `db:"Name"`
	Owner        api.Fingerprint `db:"Owner"`
	Description  string          `db:"Description"`
	LocalArrival api.Timestamp   `db:"LocalArrival"`
	DbProvable
	DbUpdateable
}

type DbThread struct {
	Fingerprint  api.Fingerprint `db:"Fingerprint"`
	Board        api.Fingerprint `db:"Board"`
	Name         string          `db:"Name"`
	Body         string          `db:"Body"`
	Link         string          `db:"Link"`
	Owner        api.Fingerprint `db:"Owner"`
	LocalArrival api.Timestamp   `db:"LocalArrival"`
	DbProvable
}

type DbPost struct {
	Fingerprint  api.Fingerprint `db:"Fingerprint"`
	Board        api.Fingerprint `db:"Board"`
	Thread       api.Fingerprint `db:"Thread"`
	Parent       api.Fingerprint `db:"Parent"`
	Body         string          `db:"Body"`
	Owner        api.Fingerprint `db:"Owner"`
	LocalArrival api.Timestamp   `db:"LocalArrival"`
	DbProvable
}

type DbVote struct {
	Fingerprint  api.Fingerprint `db:"Fingerprint"`
	Board        api.Fingerprint `db:"Board"`
	Thread       api.Fingerprint `db:"Thread"`
	Target       api.Fingerprint `db:"Target"`
	Owner        api.Fingerprint `db:"Owner"`
	Type         uint8           `db:"Type"`
	LocalArrival api.Timestamp   `db:"LocalArrival"`
	DbProvable
	DbUpdateable
}

type DbAddress struct {
	Location             api.Location  `db:"Location"`
	Sublocation          api.Location  `db:"Sublocation"`
	Port                 uint16        `db:"Port"`
	LocationType         uint8         `db:"IPType"`
	Type                 uint8         `db:"AddressType"`
	LastOnline           api.Timestamp `db:"LastOnline"`
	ProtocolVersionMajor uint8         `db:"ProtocolVersionMajor"`
	ProtocolVersionMinor uint16        `db:"ProtocolVersionMinor"`
	ProtocolExtensions   string        `db:"ProtocolExtensions"` // comma separated extension list
	ClientVersionMajor   uint8         `db:"ClientVersionMajor"`
	ClientVersionMinor   uint16        `db:"ClientVersionMinor"`
	ClientVersionPatch   uint16        `db:"ClientVersionPatch"`
	ClientName           string        `db:"ClientName"`
	LocalArrival         api.Timestamp `db:"LocalArrival"`
}

type DbKey struct {
	Fingerprint  api.Fingerprint `db:"Fingerprint"`
	Type         string          `db:"Type"`
	PublicKey    string          `db:"PublicKey"`
	Name         string          `db:"Name"`
	Info         string          `db:"Info"`
	LocalArrival api.Timestamp   `db:"LocalArrival"`
	DbProvable
	DbUpdateable
}

type DbTruststate struct {
	Fingerprint  api.Fingerprint `db:"Fingerprint"`
	Target       api.Fingerprint `db:"Target"`
	Owner        api.Fingerprint `db:"Owner"`
	Type         uint8           `db:"Type"`
	Domains      string          `db:"Domains"` // comma separated fingerprint list
	Expiry       api.Timestamp   `db:"Expiry"`
	LocalArrival api.Timestamp   `db:"LocalArrival"`
	DbProvable
	DbUpdateable
}

// Non-communicating entities
type DbNode struct {
	Fingerprint            api.Fingerprint `db:"Fingerprint"`
	BoardsLastCheckin      api.Timestamp   `db:"BoardsLastCheckin"`
	ThreadsLastCheckin     api.Timestamp   `db:"ThreadsLastCheckin"`
	PostsLastCheckin       api.Timestamp   `db:"PostsLastCheckin"`
	VotesLastCheckin       api.Timestamp   `db:"VotesLastCheckin"`
	AddressesLastCheckin   api.Timestamp   `db:"AddressesLastCheckin"`
	KeysLastCheckin        api.Timestamp   `db:"KeysLastCheckin"`
	TruststatesLastCheckin api.Timestamp   `db:"TruststatesLastCheckin"`
}

// Return types of APIToDB. This is necessary because some API objects, when converted to their DB form, return more than one DB object.

type BoardPack struct {
	Board       DbBoard
	BoardOwners []DbBoardOwner
}

type KeyPack struct {
	Key               DbKey
	CurrencyAddresses []DbCurrencyAddress
}

// APItoDB translates structs of API objects into structs of DB objects.
func APItoDB(object interface{}) (interface{}, error) {
	switch obj := object.(type) {
	// obj: typed API object.
	case api.Board:
		// Corner case: board owners
		var dbObj DbBoard
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Name = obj.Name
		dbObj.Owner = obj.Owner
		dbObj.Description = obj.Description
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		var DbBoardOwners []DbBoardOwner
		for _, val := range obj.BoardOwners {
			var bo DbBoardOwner
			bo.BoardFingerprint = obj.Fingerprint
			bo.KeyFingerprint = val.KeyFingerprint
			bo.Expiry = val.Expiry
			bo.Level = val.Level
			DbBoardOwners = append(DbBoardOwners, bo)
		}
		var result BoardPack
		result.Board = dbObj
		result.BoardOwners = DbBoardOwners
		return result, nil

	case api.Thread:
		var dbObj DbThread
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Board = obj.Board
		dbObj.Name = obj.Name
		dbObj.Body = obj.Body
		dbObj.Link = obj.Link
		dbObj.Owner = obj.Owner
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		return dbObj, nil

	case api.Post:
		var dbObj DbPost
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Board = obj.Board
		dbObj.Thread = obj.Thread
		dbObj.Parent = obj.Parent
		dbObj.Body = obj.Body
		dbObj.Owner = obj.Owner
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		return dbObj, nil

	case api.Vote:
		var dbObj DbVote
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Board = obj.Board
		dbObj.Thread = obj.Thread
		dbObj.Target = obj.Target
		dbObj.Owner = obj.Owner
		dbObj.Type = obj.Type
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		return dbObj, nil

	case api.Address:
		// Corner case, parsing protocol extensions into single field.
		var dbObj DbAddress
		dbObj.Location = obj.Location
		dbObj.Sublocation = obj.Sublocation
		dbObj.LocationType = obj.LocationType
		dbObj.Port = obj.Port
		dbObj.Type = obj.Type
		dbObj.LastOnline = obj.LastOnline
		dbObj.ProtocolVersionMajor = obj.Protocol.VersionMajor
		dbObj.ProtocolVersionMinor = obj.Protocol.VersionMinor
		for _, protExt := range obj.Protocol.Extensions {
			// Convert to comma separated list.
			if len(dbObj.ProtocolExtensions) == 0 {
				dbObj.ProtocolExtensions = string(protExt)
			} else {
				dbObj.ProtocolExtensions = fmt.Sprint(
					dbObj.ProtocolExtensions, ",", protExt)
			}
		}
		dbObj.ClientVersionMajor = obj.Client.VersionMajor
		dbObj.ClientVersionMinor = obj.Client.VersionMinor
		dbObj.ClientVersionPatch = obj.Client.VersionPatch
		dbObj.ClientName = obj.Client.ClientName
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		parsedStr, err := parseStringSliceToCommaSeparatedString(obj.Protocol.Extensions, 64, 100)
		if err != nil {
			return dbObj, err
		}
		dbObj.ProtocolExtensions = parsedStr
		return dbObj, nil

	case api.Key:
		// Corner case: currency addresses
		var dbObj DbKey
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Type = obj.Type
		dbObj.PublicKey = obj.Key
		dbObj.Name = obj.Name
		dbObj.Info = obj.Info
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		// Loop over currency addresses and insert to pack.
		var currAddrs []DbCurrencyAddress
		for _, val := range obj.CurrencyAddresses {
			var c DbCurrencyAddress
			c.KeyFingerprint = obj.Fingerprint
			c.CurrencyCode = val.CurrencyCode
			c.Address = val.Address
			currAddrs = append(currAddrs, c)
		}
		var kp KeyPack
		kp.Key = dbObj
		kp.CurrencyAddresses = currAddrs
		return kp, nil

	case api.Truststate:
		// Corner case: domain is a slice of fingerprints, convert to that.
		var dbObj DbTruststate
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Target = obj.Target
		dbObj.Owner = obj.Owner
		dbObj.Type = obj.Type
		dbObj.Expiry = obj.Expiry
		now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		parsedStr, err :=
			parseStringSliceToCommaSeparatedString(
				convertFingerprintSliceToStringSlice(obj.Domains), 64, 100)
		if err != nil {
			return dbObj, err
		}
		dbObj.Domains = parsedStr
		return dbObj, nil
	default:
		return nil, errors.New(
			fmt.Sprintf(
				"APItoDB only takes API (not DB) objects. Your object: %#v\n", obj))
	}
}

// DBtoAPI converts a DB type to an API type. This is useful when this object needs to communicated to the outside world, such as showing it to the user, or to send it over the wire. This merges certain objects (Boards will have their BoardOwners appended, and Keys, their CurrencyAddresses) and removes some information (local arrival time).
func DBtoAPI(object interface{}) (interface{}, error) {
	switch obj := object.(type) {
	// obj: typed DB object.
	case DbBoard:
		// Corner case: has to query BoardOwners, too.
		var apiObj api.Board
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Name = obj.Name
		apiObj.Owner = obj.Owner
		apiObj.Description = obj.Description
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		// Pull the board owners for this board from database.
		dbBoardOwners, err := ReadDBBoardOwners(obj.Fingerprint, "")
		if err != nil {
			logging.LogCrash(err)
		}
		for _, dbBoardOwner := range dbBoardOwners {
			var apiBoardOwner api.BoardOwner
			apiBoardOwner.KeyFingerprint = dbBoardOwner.KeyFingerprint
			apiBoardOwner.Expiry = dbBoardOwner.Expiry
			apiBoardOwner.Level = dbBoardOwner.Level
			apiObj.BoardOwners = append(apiObj.BoardOwners, apiBoardOwner)
		}
		return apiObj, nil

	case DbThread:
		var apiObj api.Thread
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Board = obj.Board
		apiObj.Name = obj.Name
		apiObj.Body = obj.Body
		apiObj.Link = obj.Link
		apiObj.Owner = obj.Owner
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		return apiObj, nil

	case DbPost:
		var apiObj api.Post
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Board = obj.Board
		apiObj.Thread = obj.Thread
		apiObj.Parent = obj.Parent
		apiObj.Body = obj.Body
		apiObj.Owner = obj.Owner
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		return apiObj, nil

	case DbVote:
		var apiObj api.Vote
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Board = obj.Board
		apiObj.Thread = obj.Thread
		apiObj.Target = obj.Target
		apiObj.Owner = obj.Owner
		apiObj.Type = obj.Type
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		return apiObj, nil

	case DbAddress:
		// Corner case, comma separated fingerprint parse
		var apiObj api.Address
		apiObj.Location = obj.Location
		apiObj.Sublocation = obj.Sublocation
		apiObj.LocationType = obj.LocationType
		apiObj.Port = obj.Port
		apiObj.Type = obj.Type
		apiObj.LastOnline = obj.LastOnline
		apiObj.Protocol.VersionMajor = obj.ProtocolVersionMajor
		apiObj.Protocol.VersionMinor = obj.ProtocolVersionMinor
		apiObj.Client.VersionMajor = obj.ClientVersionMajor
		apiObj.Client.VersionMinor = obj.ClientVersionMinor
		apiObj.Client.VersionPatch = obj.ClientVersionPatch
		apiObj.Client.ClientName = obj.ClientName
		parsedStrSlice, err := parseCommaSeparatedStringToStringSlice(obj.ProtocolExtensions, 64, 100)
		if err != nil {
			return apiObj, err
		}
		apiObj.Protocol.Extensions = parsedStrSlice
		return apiObj, nil

	case DbKey:
		// Corner case, has to query CurrencyAddresses, too.
		var apiObj api.Key
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Type = obj.Type
		apiObj.Key = obj.PublicKey
		apiObj.Name = obj.Name
		apiObj.Info = obj.Info
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		// Pull the currency addresses for this key from database.
		dbCurrAddrs, err := ReadDBCurrencyAddresses(obj.Fingerprint, "")
		if err != nil {
			logging.LogCrash(err)
		}
		for _, dbCurrAddr := range dbCurrAddrs {
			var apiCrrAddr api.CurrencyAddress
			apiCrrAddr.CurrencyCode = dbCurrAddr.CurrencyCode
			apiCrrAddr.Address = dbCurrAddr.Address
			apiObj.CurrencyAddresses = append(apiObj.CurrencyAddresses, apiCrrAddr)
		}
		return apiObj, nil

	case DbTruststate:
		// Corner case, comma separated fingerprint parse
		var apiObj api.Truststate
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Target = obj.Target
		apiObj.Owner = obj.Owner
		apiObj.Type = obj.Type
		apiObj.Expiry = obj.Expiry
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		parsedStrSlice, err := parseCommaSeparatedStringToStringSlice(obj.Domains, 64, 100)
		if err != nil {
			return apiObj, err
		}
		apiObj.Domains = convertStringSliceToFingerprintSlice(parsedStrSlice)
		return apiObj, nil

	case DbBoardOwner:
		return nil, errors.New(
			fmt.Sprintf(
				"This object cannot be queried on its own. Try querying the parent Board object. Your object: %#v\n", obj))
	case DbCurrencyAddress:
		return nil, errors.New(
			fmt.Sprintf(
				"This object cannot be queried on its own. Try querying the parent Key object. Your object: %#v\n", obj))
	default:
		return nil, errors.New(
			fmt.Sprintf(
				"DBtoAPI only takes DB (not API) objects. Your object: %#v\n", obj))
	}
}

// parseCommaSeparatedStringToStringSlice converts a single, comma separated string into a group of strings.
func parseCommaSeparatedStringToStringSlice(str string, maxLen int, maxCount int) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(str))
	result, err := reader.ReadAll()
	if err != nil {
		logging.LogCrash(err)
	}
	// Trim and check sanity.
	var cleaned []string
	var err2 error
	if result != nil { // If this isn't present, result[0] ends up undefined.
		for i, val := range result[0] {
			if i >= maxCount {
				// Too many items for the field. It's here, because otherwise somebody could send a 100 billion field string and brick the remote. We're still printing it, but not parsing it any more.
				return nil, errors.New(
					fmt.Sprintf(
						"The string provided has too many items. String: %#v\n", val))
			}
			if len(val) < maxLen && len(val) > 0 {
				cleaned = append(cleaned, strings.Trim(val, " "))
			} else if len(val) < maxLen { // The length of the string is zero.
				err2 = errors.New(
					fmt.Sprintf(
						"This string is empty. String: %#v\n", val))
			} else { // The string is longer than max length allowed for field.
				err2 = errors.New(
					fmt.Sprintf(
						"This string is too long for this field. String: %#v\n", val))
			}
		}
		err3 := checkDuplicatesInStringSlice(cleaned)
		if err3 != nil {
			return cleaned, err3
		}
	}
	return cleaned, err2
}

// parseStringSliceToCommaSeparatedString converts a slice of strings into one single comma separated string.
func parseStringSliceToCommaSeparatedString(strs []string, maxLen int, maxCount int) (string, error) {
	var err error
	var finalStr string
	if len(strs) > maxCount {
		return finalStr, errors.New(
			fmt.Sprintf(
				"The string slice provided has too many items. String slice: %#v\n", strs))
	} else {
		// We know it has an acceptable amount of elements. Let's check for duplicates.
		err := checkDuplicatesInStringSlice(strs)
		if err != nil {
			return finalStr, err
		}
	}
	for i, str := range strs {
		// Check sanity
		if len(str) < maxLen && len(str) > 0 {
			// Convert to comma separated string.
			// If this is the first item, no comma at the beginning.
			if i == 0 {
				finalStr = str
			} else {
				finalStr = fmt.Sprint(finalStr, ",", str)
			}
		} else if len(str) < maxLen { // The length of the string is zero.
			err = errors.New(
				fmt.Sprintf(
					"This string is empty. String: %#v\n", str))
		} else { // The string is longer than max length allowed for field.
			err = errors.New(
				fmt.Sprintf(
					"This string is too long for this field. String: %#v\n", str))
		}
	}
	return finalStr, err
}

func checkDuplicatesInStringSlice(strs []string) error {
	mappy := make(map[string]int)
	for _, val := range strs {
		mappy[val]++
	}
	for str, occurrenceCount := range mappy {
		if occurrenceCount > 1 {
			return errors.New(
				fmt.Sprintf(
					"This list includes items that are duplicates. Duplicate item: %#v\n", str))
		}
	}
	return nil
}

func convertFingerprintSliceToStringSlice(fps []api.Fingerprint) []string {
	var strSlice []string
	for _, val := range fps {
		strSlice = append(strSlice, string(val))
	}
	return strSlice
}

func convertStringSliceToFingerprintSlice(strs []string) []api.Fingerprint {
	var fpSlice []api.Fingerprint
	for _, val := range strs {
		fpSlice = append(fpSlice, api.Fingerprint(val))
	}
	return fpSlice
}
