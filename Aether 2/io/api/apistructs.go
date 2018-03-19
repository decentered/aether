// API > Structs
// This file provides the struct definitions for the protocol. This is what should be arriving from the network, and what should be sent over to other nodes.

package api

import (
	"database/sql/driver"
	// "fmt"
	"aether-core/services/fingerprinting"
	"aether-core/services/globals"
	"aether-core/services/proofofwork"
	"aether-core/services/signaturing"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
)

// Structs for the entity types. There are 7 types. Board, Thread, Post, Vote, Key, Address, Truststate.

// Low-level types

// type Fingerprint [64]byte // 64 char ASCII
type Fingerprint string // 64 char ASCII
type Timestamp int64    // UNIX Timestamp
// type ProofOfWork [1024]byte
type ProofOfWork string // temp
// type Signature [512]byte
type Signature string // temp
type Location string

func (f Fingerprint) Value() (driver.Value, error) {
	return string(f), nil
}

func (f *Fingerprint) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*f = Fingerprint(stringVal)
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *Timestamp) Scan(value interface{}) error {
	numVal := value.(int64)
	*t = Timestamp(numVal)
	return nil
}

func (p ProofOfWork) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *ProofOfWork) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*p = ProofOfWork(stringVal)
	return nil
}

func (s Signature) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Signature) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*s = Signature(stringVal)
	return nil
}

func (l Location) Value() (driver.Value, error) {
	return string(l), nil
}

func (l *Location) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*l = Location(stringVal)
	return nil
}

// Basic properties

type ProvableFieldSet struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Creation    Timestamp   `json:"creation"`
	ProofOfWork ProofOfWork `json:"proof_of_work"`
	Signature   Signature   `json:"signature"`
}

type UpdateableFieldSet struct { // Common set of properties for all objects that are updateable.
	LastUpdate        Timestamp   `json:"last_update"`
	UpdateProofOfWork ProofOfWork `json:"update_proof_of_work"`
	UpdateSignature   Signature   `json:"update_signature"`
}

// Subentities

type BoardOwner struct {
	KeyFingerprint Fingerprint `json:"key_fingerprint"` // Fingerprint of the key the ownership is associated to.
	Expiry         Timestamp   `json:"expiry"`          // When the ownership expires.
	Level          uint8       `json:"level"`           // Admin (2), mod (1), or abdicated admin (0)
}

type CurrencyAddress struct {
	// CurrencyCode [5]byte    `json:"currency_code"`
	CurrencyCode string `json:"currency_code"`
	Address      string `json:"address"`
}

type Subprotocol struct {
	Name              string   `json:"name"`
	VersionMajor      uint8    `json:"version_major"`
	VersionMinor      uint16   `json:"version_minor"`
	SupportedEntities []string `json:"supported_entities"`
}

type Protocol struct {
	VersionMajor uint8         `json:"version_major"`
	VersionMinor uint16        `json:"version_minor"`
	Subprotocols []Subprotocol `json:"subprotocols"`
}

type Client struct {
	VersionMajor uint8  `json:"version_major"`
	VersionMinor uint16 `json:"version_minor"`
	VersionPatch uint16 `json:"version_patch"`
	ClientName   string `json:"name"` // Max 255
}

// Entities

type Board struct {
	ProvableFieldSet
	Name        string       `json:"name"`         // Max 255 char unicode
	BoardOwners []BoardOwner `json:"board_owners"` // max 100 owners
	Description string       `json:"description"`  // Max 65535 char unicode
	Owner       Fingerprint  `json:"owner"`
	Meta        string       `json:"meta"` // This is the dynamic JSON field
	UpdateableFieldSet
}

type Thread struct {
	ProvableFieldSet
	Board Fingerprint `json:"board"`
	Name  string      `json:"name"`
	Body  string      `json:"body"`
	Link  string      `json:"link"`
	Owner Fingerprint `json:"owner"`
	Meta  string      `json:"meta"`
}

type Post struct {
	ProvableFieldSet
	Board  Fingerprint `json:"board"`
	Thread Fingerprint `json:"thread"`
	Parent Fingerprint `json:"parent"`
	Body   string      `json:"body"`
	Owner  Fingerprint `json:"owner"`
	Meta   string      `json:"meta"`
}

type Vote struct {
	ProvableFieldSet
	Board  Fingerprint `json:"board"`
	Thread Fingerprint `json:"thread"`
	Target Fingerprint `json:"target"`
	Owner  Fingerprint `json:"owner"`
	Type   uint8       `json:"type"`
	Meta   string      `json:"meta"`
	UpdateableFieldSet
}

type Address struct {
	Location     Location  `json:"location"`
	Sublocation  Location  `json:"sublocation"`
	LocationType uint8     `json:"location_type"`
	Port         uint16    `json:"port"`
	Type         uint8     `json:"type"`
	LastOnline   Timestamp `json:"last_online"`
	Protocol     Protocol  `json:"protocol"`
	Client       Client    `json:"client"`
}

type Key struct {
	ProvableFieldSet
	Type              string            `json:"type"`
	Key               string            `json:"key"`
	Name              string            `json:"name"`
	CurrencyAddresses []CurrencyAddress `json:"currency_addresses"` // max 10 addresses
	Info              string            `json:"info"`
	Meta              string            `json:"meta"`
	UpdateableFieldSet
}

type Truststate struct {
	ProvableFieldSet
	Target  Fingerprint   `json:"target"`
	Owner   Fingerprint   `json:"owner"`
	Type    uint8         `json:"type"`
	Domains []Fingerprint `json:"domain"` // max 100 domains fingerprint
	Expiry  Timestamp     `json:"expiry"`
	Meta    string        `json:"meta"`
	UpdateableFieldSet
}

type ResultCache struct { // These are caches shown in the index endpoint of a particular entity.
	ResponseUrl string    `json:"response_url"`
	StartsFrom  Timestamp `json:"starts_from"`
	EndsAt      Timestamp `json:"ends_at"`
}

// Index Form Entities: These are index forms of the entities above.

type BoardIndex struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Creation    Timestamp   `json:"creation"`
	LastUpdate  Timestamp   `json:"last_update"`
	PageNumber  int         `json:"page_number"`
}

type ThreadIndex struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Board       Fingerprint `json:"board"`
	Creation    Timestamp   `json:"creation"`
	PageNumber  int         `json:"page_number"`
}

type PostIndex struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Board       Fingerprint `json:"board"`
	Thread      Fingerprint `json:"thread"`
	Creation    Timestamp   `json:"creation"`
	PageNumber  int         `json:"page_number"`
}

type VoteIndex struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Board       Fingerprint `json:"board"`
	Thread      Fingerprint `json:"thread"`
	Target      Fingerprint `json:"target"`
	Creation    Timestamp   `json:"creation"`
	LastUpdate  Timestamp   `json:"last_update"`
	PageNumber  int         `json:"page_number"`
}

type AddressIndex Address

type KeyIndex struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Creation    Timestamp   `json:"creation"`
	LastUpdate  Timestamp   `json:"last_update"`
	PageNumber  int         `json:"page_number"`
}

type TruststateIndex struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Target      Fingerprint `json:"target"`
	Creation    Timestamp   `json:"creation"`
	LastUpdate  Timestamp   `json:"last_update"`
	PageNumber  int         `json:"page_number"`
}

// Response types

type Pagination struct {
	Pages       uint64 `json:"pages"`
	CurrentPage uint64 `json:"current_page"`
}

type Caching struct {
	ServedFromCache bool   `json:"served_from_cache"`
	CacheScope      string `json:"cache_scope"`
	PrevCacheUrl    string `json:"prev_cache_url"`
	CurrentCacheUrl string `json:"current_cache_url"`
}

type Filter struct { // Timestamp filter or embeds, or fingerprint
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type Answer struct { // Bodies of API Endpoint responses from remote. This will be filled and unused field will be omitted.
	Boards            []Board           `json:"boards,omitempty"`
	BoardIndexes      []BoardIndex      `json:"boards_index,omitempty"`
	Threads           []Thread          `json:"threads,omitempty"`
	ThreadIndexes     []ThreadIndex     `json:"threads_index,omitempty"`
	Posts             []Post            `json:"posts,omitempty"`
	PostIndexes       []PostIndex       `json:"posts_index,omitempty"`
	Votes             []Vote            `json:"votes,omitempty"`
	VoteIndexes       []VoteIndex       `json:"votes_index,omitempty"`
	Keys              []Key             `json:"keys,omitempty"`
	KeyIndexes        []KeyIndex        `json:"keys_index,omitempty"`
	Addresses         []Address         `json:"addresses,omitempty"`
	AddressIndexes    []AddressIndex    `json:"addresses_index,omitempty"`
	Truststates       []Truststate      `json:"truststates,omitempty"`
	TruststateIndexes []TruststateIndex `json:"truststates_index,omitempty"`
}

// Response styles.

// Response is the interface junction that batch processing functions take and emit. This is the 'internal' communication structure within the backend. It is the big carrier type for the end result of a pull from a remote.
type Response struct {
	AvailableTypes            []string
	Boards                    []Board
	BoardIndexes              []BoardIndex
	Threads                   []Thread
	ThreadIndexes             []ThreadIndex
	Posts                     []Post
	PostIndexes               []PostIndex
	Votes                     []Vote
	VoteIndexes               []VoteIndex
	Keys                      []Key
	KeyIndexes                []KeyIndex
	Addresses                 []Address
	AddressIndexes            []AddressIndex
	Truststates               []Truststate
	TruststateIndexes         []TruststateIndex
	CacheLinks                []ResultCache
	MostRecentSourceTimestamp Timestamp
}

// ApiResponse is the blueprint of all requests and responses. This is the 'external' communication structure backend uses to talk to other backends.
type ApiResponse struct {
	NodeId       Fingerprint   `json:"node_id,omitempty"`
	Address      Address       `json:"address,omitempty"`
	Entity       string        `json:"entity,omitempty"`
	Endpoint     string        `json:"endpoint,omitempty"`
	Filters      []Filter      `json:"filters,omitempty"`
	Timestamp    Timestamp     `json:"timestamp,omitempty"`
	StartsFrom   Timestamp     `json:"starts_from,omitempty"`
	EndsAt       Timestamp     `json:"ends_at,omitempty"`
	Pagination   Pagination    `json:"pagination,omitempty"`
	Caching      Caching       `json:"caching,omitempty"`
	Results      []ResultCache `json:"results,omitempty"`  // Pages
	ResponseBody Answer        `json:"response,omitempty"` // Entities, Full size or Index versions.
}

// // Interfaces

type Fingerprintable interface {
	GetFingerprint() Fingerprint // Field accessor
	CreateFingerprint()
	VerifyFingerprint() bool
}

type PoWAble interface {
	GetProofOfWork() ProofOfWork // Field accessor
	CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error
	VerifyPoW(pubKey string) (bool, error)
}

type Signable interface {
	GetSignature() Signature // Field accessor
	CreateSignature(keyPair *ecdsa.PrivateKey) error
	VerifySignature(pubKey string) (bool, error)
}

type Provable interface {
	Fingerprintable
	PoWAble
	Signable
	GetOwner() Fingerprint
}

type Updateable interface {
	GetUpdateProofOfWork() ProofOfWork // Field accessor
	GetUpdateSignature() Signature     // Field accessor
	CreateUpdatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error
	CreateUpdateSignature(keyPair *ecdsa.PrivateKey) error
}

// Accessor methods. These methods allow access to fields from the interfaces. The reason why we need these is that interfaces cannot take struct fields, so I have to create these accessor methods to let them be accessible over interfaces.

// Fingerprint accessors

func (entity *Board) GetFingerprint() Fingerprint      { return entity.Fingerprint }
func (entity *Thread) GetFingerprint() Fingerprint     { return entity.Fingerprint }
func (entity *Post) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *Vote) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *Key) GetFingerprint() Fingerprint        { return entity.Fingerprint }
func (entity *Truststate) GetFingerprint() Fingerprint { return entity.Fingerprint }

// Signature accessors

func (entity *Board) GetSignature() Signature      { return entity.Signature }
func (entity *Thread) GetSignature() Signature     { return entity.Signature }
func (entity *Post) GetSignature() Signature       { return entity.Signature }
func (entity *Vote) GetSignature() Signature       { return entity.Signature }
func (entity *Key) GetSignature() Signature        { return entity.Signature }
func (entity *Truststate) GetSignature() Signature { return entity.Signature }

// UpdateSignature accessors

func (entity *Board) GetUpdateSignature() Signature      { return entity.UpdateSignature }
func (entity *Vote) GetUpdateSignature() Signature       { return entity.UpdateSignature }
func (entity *Key) GetUpdateSignature() Signature        { return entity.UpdateSignature }
func (entity *Truststate) GetUpdateSignature() Signature { return entity.UpdateSignature }

// ProofOfWork accessors

func (entity *Board) GetProofOfWork() ProofOfWork      { return entity.ProofOfWork }
func (entity *Thread) GetProofOfWork() ProofOfWork     { return entity.ProofOfWork }
func (entity *Post) GetProofOfWork() ProofOfWork       { return entity.ProofOfWork }
func (entity *Vote) GetProofOfWork() ProofOfWork       { return entity.ProofOfWork }
func (entity *Key) GetProofOfWork() ProofOfWork        { return entity.ProofOfWork }
func (entity *Truststate) GetProofOfWork() ProofOfWork { return entity.ProofOfWork }

// UpdateProofOfWork accessors

func (entity *Board) GetUpdateProofOfWork() ProofOfWork      { return entity.UpdateProofOfWork }
func (entity *Vote) GetUpdateProofOfWork() ProofOfWork       { return entity.UpdateProofOfWork }
func (entity *Key) GetUpdateProofOfWork() ProofOfWork        { return entity.UpdateProofOfWork }
func (entity *Truststate) GetUpdateProofOfWork() ProofOfWork { return entity.UpdateProofOfWork }

// Signature accessors

func (entity *Board) GetOwner() Fingerprint  { return entity.Owner }
func (entity *Thread) GetOwner() Fingerprint { return entity.Owner }
func (entity *Post) GetOwner() Fingerprint   { return entity.Owner }
func (entity *Vote) GetOwner() Fingerprint   { return entity.Owner }

// Owner of the key is itself.
func (entity *Key) GetOwner() Fingerprint        { return entity.Fingerprint }
func (entity *Truststate) GetOwner() Fingerprint { return entity.Owner }

// // Create ProofOfWork

func (b *Board) CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *b
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	b.ProofOfWork = ProofOfWork(pow)
	return nil
}

func (t *Thread) CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *t
	// Non-updateable
	cpI.Fingerprint = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	t.ProofOfWork = ProofOfWork(pow)
	return nil
}

func (p *Post) CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *p
	// Non-updateable
	cpI.Fingerprint = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	p.ProofOfWork = ProofOfWork(pow)
	return nil
}

func (v *Vote) CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *v
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	v.ProofOfWork = ProofOfWork(pow)
	return nil
}

func (k *Key) CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *k
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	k.ProofOfWork = ProofOfWork(pow)
	return nil
}

func (ts *Truststate) CreatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *ts
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	ts.ProofOfWork = ProofOfWork(pow)
	return nil
}

// Create UpdateProofOfWork

func (b *Board) CreateUpdatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *b
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	b.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func (v *Vote) CreateUpdatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *v
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	v.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func (k *Key) CreateUpdatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *k
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	k.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func (ts *Truststate) CreateUpdatePoW(keyPair *ecdsa.PrivateKey, difficulty int) error {
	cpI := *ts
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	ts.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

// Verify ProofOfWork

func (b *Board) VerifyPoW(pubKey string) (bool, error) {
	cpI := *b
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func (t *Thread) VerifyPoW(pubKey string) (bool, error) {
	cpI := *t
	// Non-updateable
	cpI.Fingerprint = ""
	// Save PoW to be verified
	pow := string(cpI.ProofOfWork)
	// Delete PoW so that the PoW will match
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= globals.BackendConfig.GetMinimumPoWStrengths().Thread {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func (p *Post) VerifyPoW(pubKey string) (bool, error) {
	cpI := *p
	// Non-updateable
	cpI.Fingerprint = ""
	// Save PoW to be verified
	pow := string(cpI.ProofOfWork)
	// Delete PoW so that the PoW will match
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= globals.BackendConfig.GetMinimumPoWStrengths().Post {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func (v *Vote) VerifyPoW(pubKey string) (bool, error) {
	cpI := *v
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func (k *Key) VerifyPoW(pubKey string) (bool, error) {
	cpI := *k
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func (ts *Truststate) VerifyPoW(pubKey string) (bool, error) {
	cpI := *ts
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

// Create Fingerprint

func (b *Board) CreateFingerprint() {
	cpI := *b
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	var emptyBOList []BoardOwner
	cpI.BoardOwners = emptyBOList
	cpI.Description = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	b.Fingerprint = Fingerprint(fp)
}

func (t *Thread) CreateFingerprint() {
	cpI := *t
	// Remove ALL mutable fields
	// (Thread does not have any mutable fields)
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	t.Fingerprint = Fingerprint(fp)
}

func (p *Post) CreateFingerprint() {
	cpI := *p
	// Remove ALL mutable fields
	// (Post does not have any mutable fields)
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	p.Fingerprint = Fingerprint(fp)
}

func (v *Vote) CreateFingerprint() {
	cpI := *v
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.Type = uint8(0)
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	v.Fingerprint = Fingerprint(fp)
}

func (k *Key) CreateFingerprint() {
	cpI := *k
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	var emptyCAList []CurrencyAddress
	cpI.CurrencyAddresses = emptyCAList
	cpI.Info = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	fmt.Println("old fingerprint")
	fmt.Println(cpI.Fingerprint)
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	k.Fingerprint = Fingerprint(fp)
	fmt.Println("new fingerprint")
	fmt.Println(k.Fingerprint)
}

func (ts *Truststate) CreateFingerprint() {
	cpI := *ts
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.Type = uint8(0)
	var emptyDList []Fingerprint
	cpI.Domains = emptyDList
	cpI.Expiry = 0
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	ts.Fingerprint = Fingerprint(fp)
}

// Verify Fingerprint

func (b *Board) VerifyFingerprint() bool {
	cpI := *b
	var fp string
	fp = string(cpI.Fingerprint)
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	var emptyBOList []BoardOwner
	cpI.BoardOwners = emptyBOList
	cpI.Description = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func (t *Thread) VerifyFingerprint() bool {
	cpI := *t
	var fp string
	fp = string(cpI.Fingerprint)
	// Remove ALL mutable fields
	// (Thread does not have any mutable fields)
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func (p *Post) VerifyFingerprint() bool {
	cpI := *p
	var fp string
	fp = string(cpI.Fingerprint)
	// Remove ALL mutable fields
	// (Post does not have any mutable fields)
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func (v *Vote) VerifyFingerprint() bool {
	cpI := *v
	var fp string
	fp = string(cpI.Fingerprint)
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.Type = uint8(0)
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func (k *Key) VerifyFingerprint() bool {
	cpI := *k
	var fp string
	fp = string(cpI.Fingerprint)
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	var emptyCAList []CurrencyAddress
	cpI.CurrencyAddresses = emptyCAList
	cpI.Info = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func (ts *Truststate) VerifyFingerprint() bool {
	cpI := *ts
	var fp string
	fp = string(cpI.Fingerprint)
	// Remove ALL mutable fields
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.Type = uint8(0)
	var emptyDList []Fingerprint
	cpI.Domains = emptyDList
	cpI.Expiry = 0
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

// Signature

func (b *Board) CreateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *b
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	b.Signature = Signature(signature)
	return nil
}

func (t *Thread) CreateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *t
	// Non-updateable
	cpI.Fingerprint = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	t.Signature = Signature(signature)
	return nil
}

func (p *Post) CreateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *p
	// Non-updateable
	cpI.Fingerprint = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	p.Signature = Signature(signature)
	return nil
}

func (v *Vote) CreateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *v
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	v.Signature = Signature(signature)
	return nil
}

func (k *Key) CreateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *k
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	k.Signature = Signature(signature)
	return nil
}

func (ts *Truststate) CreateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *ts
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	ts.Signature = Signature(signature)
	return nil
}

// Create UpdateSignature

func (b *Board) CreateUpdateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *b
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	b.UpdateSignature = Signature(signature)
	return nil
}

func (v *Vote) CreateUpdateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *v
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	v.UpdateSignature = Signature(signature)
	return nil
}

func (k *Key) CreateUpdateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *k
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	k.UpdateSignature = Signature(signature)
	return nil
}

func (ts *Truststate) CreateUpdateSignature(keyPair *ecdsa.PrivateKey) error {
	cpI := *ts
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	ts.UpdateSignature = Signature(signature)
	return nil
}

// Verify Signature

func (b *Board) VerifySignature(pubKey string) (bool, error) {
	cpI := *b
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func (t *Thread) VerifySignature(pubKey string) (bool, error) {
	cpI := *t
	// Save signature to be verified
	signature := string(cpI.Signature)
	// Non-updateable
	cpI.Fingerprint = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Delete signature so that the signature will match
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func (p *Post) VerifySignature(pubKey string) (bool, error) {
	cpI := *p
	// Save signature to be verified
	signature := string(cpI.Signature)
	// Non-updateable
	cpI.Fingerprint = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Delete signature so that the signature will match
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func (v *Vote) VerifySignature(pubKey string) (bool, error) {
	cpI := *v
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func (k *Key) VerifySignature(pubKey string) (bool, error) {
	cpI := *k
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func (ts *Truststate) VerifySignature(pubKey string) (bool, error) {
	cpI := *ts
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}
