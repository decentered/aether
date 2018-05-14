// API > Structs
// This file provides the struct definitions for the protocol. This is what should be arriving from the network, and what should be sent over to other nodes.

package api

import (
	"database/sql/driver"
	// "fmt"
	"aether-core/services/logging"
	"aether-core/services/toolbox"
	"fmt"
	"golang.org/x/crypto/ed25519"
	// "github.com/davecgh/go-spew/spew"
	"aether-core/services/globals"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
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

func (t Timestamp) Humanise() string {
	if t != 0 {
		return fmt.Sprintf("%s (%d)", time.Unix(int64(t), 0).Format(time.Stamp), t)
	} else {
		return fmt.Sprint("Blank")
	}
}

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
	Verified    bool        `json:"-"`
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
	Level          uint8       `json:"level"`           // mod(1)
}

type Subprotocol struct {
	Name              string   `json:"name"` //2-16 chars
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
	ClientName   string `json:"name"` // Max 128
}

// Entities

type Board struct { // Mutables: BoardOwners, Description, Meta
	ProvableFieldSet
	Name           string       `json:"name"`         // Max 255 char unicode
	BoardOwners    []BoardOwner `json:"board_owners"` // max 128 owners
	Description    string       `json:"description"`  // Max 65535 char unicode
	Owner          Fingerprint  `json:"owner"`
	OwnerPublicKey string       `json:"owner_publickey"`
	EntityVersion  int          `json:"entity_version"`
	Language       string       `json:"language"`
	Meta           string       `json:"meta"` // This is the dynamic JSON field
	RealmId        Fingerprint  `json:"realm_id"`
	EncrContent    string       `json:"encrcontent"`
	UpdateableFieldSet
}

type Thread struct { // Mutables: Body, Meta
	ProvableFieldSet
	Board          Fingerprint `json:"board"`
	Name           string      `json:"name"`
	Body           string      `json:"body"`
	Link           string      `json:"link"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Post struct { // Mutables: Body, Meta
	ProvableFieldSet
	Board          Fingerprint `json:"board"`
	Thread         Fingerprint `json:"thread"`
	Parent         Fingerprint `json:"parent"`
	Body           string      `json:"body"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Vote struct { // Mutables: Type, Meta
	ProvableFieldSet
	Board          Fingerprint `json:"board"`
	Thread         Fingerprint `json:"thread"`
	Target         Fingerprint `json:"target"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	Type           int         `json:"type"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

// TODO: blocking the json output might be a good way to prevent values leaking out or in.
type Address struct { // Mutables: None
	Location           Location    `json:"location"`
	Sublocation        Location    `json:"sublocation"`
	LocationType       uint8       `json:"location_type"`
	Port               uint16      `json:"port"`
	Type               uint8       `json:"type"`
	LastSuccessfulPing Timestamp   `json:"-"`
	LastSuccessfulSync Timestamp   `json:"-"`
	Protocol           Protocol    `json:"protocol"`
	Client             Client      `json:"client"`
	EntityVersion      int         `json:"entity_version"`
	RealmId            Fingerprint `json:"realm_id"`
	Verified           bool        `json:"-"` // This is normally part of the provable field set, but address is not provable, so provided here separately.
}

type Key struct { // Mutables: Expiry, Info, Meta
	ProvableFieldSet
	Type          string      `json:"type"`
	Key           string      `json:"key"`
	Expiry        Timestamp   `json:"expiry"`
	Name          string      `json:"name"`
	Info          string      `json:"info"`
	EntityVersion int         `json:"entity_version"`
	Meta          string      `json:"meta"`
	RealmId       Fingerprint `json:"realm_id"`
	EncrContent   string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Truststate struct { // Mutables: Type, Domains, Expiry, Meta
	ProvableFieldSet
	Target         Fingerprint   `json:"target"`
	Owner          Fingerprint   `json:"owner"`
	OwnerPublicKey string        `json:"owner_publickey"`
	Type           int           `json:"type"`
	Domains        []Fingerprint `json:"domain"` // max 100 domains fingerprint
	Expiry         Timestamp     `json:"expiry"`
	EntityVersion  int           `json:"entity_version"`
	Meta           string        `json:"meta"`
	RealmId        Fingerprint   `json:"realm_id"`
	EncrContent    string        `json:"encrcontent"`
	UpdateableFieldSet
}

// Index Form Entities: These are index forms of the entities above.

type BoardIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type ThreadIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Board         Fingerprint `json:"board"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type PostIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Board         Fingerprint `json:"board"`
	Thread        Fingerprint `json:"thread"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type VoteIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Board         Fingerprint `json:"board"`
	Thread        Fingerprint `json:"thread"`
	Target        Fingerprint `json:"target"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type AddressIndex Address

type KeyIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type TruststateIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Target        Fingerprint `json:"target"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

// Response types

type Pagination struct {
	Pages       uint64 `json:"pages"`
	CurrentPage uint64 `json:"current_page"`
}

type Caching struct {
	Pregenerated    bool          `json:"pregenerated"`
	CurrentCacheUrl string        `json:"current_cache_url"`
	EntityCounts    []EntityCount `json:"entity_counts"`
}

type EntityCount struct {
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type Filter struct { // Timestamp filter or embeds, or fingerprint
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type ResultCache struct { // These are caches shown in the index endpoint of a particular entity.
	ResponseUrl string    `json:"response_url"`
	StartsFrom  Timestamp `json:"starts_from"`
	EndsAt      Timestamp `json:"ends_at"`
}

type Answer struct { // Bodies of API Endpoint responses from remote. This will be filled and unused field will be omitted.
	Boards      []Board      `json:"boards,omitempty"`
	Threads     []Thread     `json:"threads,omitempty"`
	Posts       []Post       `json:"posts,omitempty"`
	Votes       []Vote       `json:"votes,omitempty"`
	Keys        []Key        `json:"keys,omitempty"`
	Truststates []Truststate `json:"truststates,omitempty"`
	Addresses   []Address    `json:"addresses,omitempty"`

	BoardIndexes      []BoardIndex      `json:"boards_index,omitempty"`
	ThreadIndexes     []ThreadIndex     `json:"threads_index,omitempty"`
	PostIndexes       []PostIndex       `json:"posts_index,omitempty"`
	VoteIndexes       []VoteIndex       `json:"votes_index,omitempty"`
	KeyIndexes        []KeyIndex        `json:"keys_index,omitempty"`
	TruststateIndexes []TruststateIndex `json:"truststates_index,omitempty"`
	AddressIndexes    []AddressIndex    `json:"addresses_index,omitempty"`

	BoardManifests      []PageManifest `json:"boards_manifest,omitempty"`
	ThreadManifests     []PageManifest `json:"threads_manifest,omitempty"`
	PostManifests       []PageManifest `json:"posts_manifest,omitempty"`
	VoteManifests       []PageManifest `json:"votes_manifest,omitempty"`
	KeyManifests        []PageManifest `json:"keys_manifest,omitempty"`
	TruststateManifests []PageManifest `json:"truststates_manifest,omitempty"`
	AddressManifests    []PageManifest `json:"addresses_manifest,omitempty"`
}

// Manifest type
type PageManifest struct {
	Page     uint64               `json:"page_number"`
	Entities []PageManifestEntity `json:"entities"`
}

type PageManifestEntity struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	LastUpdate  Timestamp   `json:"last_update"`
}

// ApiResponse is the blueprint of all requests and responses. This is the 'external' communication structure backend uses to talk to other backends.
type ApiResponse struct {
	NodeId        Fingerprint   `json:"-"` // Generated and used at the ApiResponse signature verification, from the NodePublicKey. It doesn't transmit in or out, only generated on the fly. This blocks both inbound and outbound.
	NodePublicKey string        `json:"node_public_key,omitempty"`
	Signature     Signature     `json:"page_signature,omitempty"`
	Address       Address       `json:"address,omitempty"`
	Entity        string        `json:"entity,omitempty"`
	Endpoint      string        `json:"endpoint,omitempty"`
	Filters       []Filter      `json:"filters,omitempty"`
	Timestamp     Timestamp     `json:"timestamp,omitempty"`
	StartsFrom    Timestamp     `json:"starts_from,omitempty"`
	EndsAt        Timestamp     `json:"ends_at,omitempty"`
	Pagination    Pagination    `json:"pagination,omitempty"`
	Caching       Caching       `json:"caching,omitempty"`
	Results       []ResultCache `json:"results,omitempty"`  // Pages
	ResponseBody  Answer        `json:"response,omitempty"` // Entities, Full size or Index versions.
}

// GetProvables gets all provables in an ApiResponse.
func (r *ApiResponse) GetProvables() *[]Provable {
	p := []Provable{}
	for key, _ := range r.ResponseBody.Boards {
		p = append(p, Provable(&r.ResponseBody.Boards[key]))
	}
	for key, _ := range r.ResponseBody.Threads {
		p = append(p, Provable(&r.ResponseBody.Threads[key]))
	}
	for key, _ := range r.ResponseBody.Posts {
		p = append(p, Provable(&r.ResponseBody.Posts[key]))
	}
	for key, _ := range r.ResponseBody.Votes {
		p = append(p, Provable(&r.ResponseBody.Votes[key]))
	}
	for key, _ := range r.ResponseBody.Keys {
		p = append(p, Provable(&r.ResponseBody.Keys[key]))
	}
	for key, _ := range r.ResponseBody.Truststates {
		p = append(p, Provable(&r.ResponseBody.Truststates[key]))
	}
	return &p
}

// Dump dumps the apiresponse in JSON format to a predetermined location on disk for inspection.
func (r *ApiResponse) Dump() error {
	fileContents, err := json.Marshal(r)
	if err != nil {
		return errors.New(fmt.Sprint(
			"This ApiResponse failed to convert to JSON. Error: %#v, ApiResponse: %#v", err, r))
	}
	path := globals.BackendConfig.GetCachesDirectory() + "/dumps/"
	toolbox.CreatePath(path)
	filename := fmt.Sprint("ApiRespDump-", time.Now().Unix(), ".json")
	// fmt.Printf("%s%s%s\n", path, "/dumps/", filename)
	err2 := ioutil.WriteFile(fmt.Sprint(path, filename), fileContents, 0755)
	if err2 != nil {
		logging.LogCrash(err2)
	}
	return nil
}

// Verify verifies all items and flags them appropriately in a response.
func (r *ApiResponse) Verify() []error {
	errs := []error{}
	// First of all, run boundary check on the apiresponse. This does NOT check for the .Address field, neither does it check the entities contained within. We'll do those afterwards.
	boundsOk, err := r.CheckBounds()
	if len(r.ResponseBody.PostIndexes) > 0 && boundsOk {
		// fmt.Println("this api response has post indexes in it and it verified.")
	} else if len(r.ResponseBody.PostIndexes) > 0 {
		fmt.Println("this api response has post indexes in it and it did not verify.")
	}
	if err != nil {
		return []error{errors.New(fmt.Sprintf("This ApiResponse failed the boundary check for its general structure (not for its contents -- it didn't come to that.) Error: %#v, ApiResponse: %#v", err, r))}
	}
	if !boundsOk {
		// logging.LogCrash("yo")
		return []error{errors.New(fmt.Sprintf("This ApiResponse failed the boundary check for its general structure (not for its contents -- it didn't come to that.) ApiResponse: %#v", r))}
	}
	remoteAddrOk, err := r.Address.CheckBounds()
	if err != nil {
		return []error{err}
	}
	if !remoteAddrOk {
		return []error{errors.New(fmt.Sprintf("This ApiResponse's remote Address failed the boundary check. ApiResponse.Address: %#v", r.Address))}
	}
	// This is all the verification we need for addresses - just a bounds check. It does not go into the more involved Verify() flow.
	for key, _ := range r.ResponseBody.Addresses { // this is a concrete type..
		err := Verify(&r.ResponseBody.Addresses[key])
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	provables := r.GetProvables()
	for _, e := range *provables { // provable is an interface, so pointer..
		err := Verify(e)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	for _, err := range errs {
		logging.Log(1, err)
	}
	return errs
}

func (r *ApiResponse) ToJSON() ([]byte, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return result, errors.New(fmt.Sprint(
			"This ApiResponse failed to convert to JSON. Error: %#v, ApiResponse: %#v", err, r))
	}
	return result, nil
}

func (r *ApiResponse) Prefill() {
	subprotsAsShims := globals.BackendConfig.GetServingSubprotocols()
	subprotsSupported := []Subprotocol{}
	for _, val := range subprotsAsShims {
		subprotsSupported = append(subprotsSupported, Subprotocol(val))
	}
	r.NodePublicKey = globals.BackendConfig.GetMarshaledBackendPublicKey()
	addr := Address{}
	addr.LocationType = globals.BackendConfig.GetExternalIpType()
	addr.Type = 2 // This is a live node.
	addr.Port = uint16(globals.BackendConfig.GetExternalPort())
	addr.Protocol.VersionMajor = globals.BackendConfig.GetProtocolVersionMajor()
	addr.Protocol.VersionMinor = globals.BackendConfig.GetProtocolVersionMinor()
	addr.Protocol.Subprotocols = subprotsSupported
	addr.Client.VersionMajor = globals.BackendConfig.GetClientVersionMajor()
	addr.Client.VersionMinor = globals.BackendConfig.GetClientVersionMinor()
	addr.Client.VersionPatch = globals.BackendConfig.GetClientVersionPatch()
	addr.Client.ClientName = globals.BackendConfig.GetClientName()
	addr.EntityVersion = globals.BackendTransientConfig.EntityVersions.Address
	r.Address = addr
}

// // Interfaces

type Fingerprintable interface {
	GetFingerprint() Fingerprint // Field accessor
	CreateFingerprint() error
	VerifyFingerprint() bool
}

type PoWAble interface {
	GetProofOfWork() ProofOfWork // Field accessor
	CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error
	VerifyPoW(pubKey string) (bool, error)
}

type Signable interface {
	GetSignature() Signature   // Field accessor
	GetOwnerPublicKey() string // Field accessor
	CreateSignature(keyPair *ed25519.PrivateKey) error
	VerifySignature(pubKey string) (bool, error)
}

type BoundsCheckable interface {
	CheckBounds() (bool, error)
}

type Verifiable interface {
	Fingerprintable
	PoWAble
	Signable
	BoundsCheckable
	SetVerified(bool)
	GetVerified() bool
}

type Encryptable interface {
	GetEncrContent() string
}

type Shardable interface {
	GetRealmId() Fingerprint
}

type Provable interface {
	Verifiable
	Shardable
	Encryptable
	GetOwner() Fingerprint
	GetLastUpdate() Timestamp
}

type Updateable interface {
	GetUpdateProofOfWork() ProofOfWork // Field accessor
	GetUpdateSignature() Signature     // Field accessor
	CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error
	CreateUpdateSignature(keyPair *ed25519.PrivateKey) error
}

type Versionable interface {
	GetVersion() int
}

// Accessor methods. These methods allow access to fields from the interfaces. The reason why we need these is that interfaces cannot take struct fields, so I have to create these accessor methods to let them be accessible over interfaces.

// Version accessors

func (entity *Board) GetVersion() int      { return entity.EntityVersion }
func (entity *Thread) GetVersion() int     { return entity.EntityVersion }
func (entity *Post) GetVersion() int       { return entity.EntityVersion }
func (entity *Vote) GetVersion() int       { return entity.EntityVersion }
func (entity *Key) GetVersion() int        { return entity.EntityVersion }
func (entity *Truststate) GetVersion() int { return entity.EntityVersion }
func (entity *Address) GetVersion() int    { return entity.EntityVersion }

// Fingerprint accessors

func (entity *Board) GetFingerprint() Fingerprint      { return entity.Fingerprint }
func (entity *Thread) GetFingerprint() Fingerprint     { return entity.Fingerprint }
func (entity *Post) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *Vote) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *Key) GetFingerprint() Fingerprint        { return entity.Fingerprint }
func (entity *Truststate) GetFingerprint() Fingerprint { return entity.Fingerprint }

// LastUpdate accessors

func (entity *Board) GetLastUpdate() Timestamp      { return entity.LastUpdate }
func (entity *Thread) GetLastUpdate() Timestamp     { return entity.LastUpdate }
func (entity *Post) GetLastUpdate() Timestamp       { return entity.LastUpdate }
func (entity *Vote) GetLastUpdate() Timestamp       { return entity.LastUpdate }
func (entity *Key) GetLastUpdate() Timestamp        { return entity.LastUpdate }
func (entity *Truststate) GetLastUpdate() Timestamp { return entity.LastUpdate }

// Signature accessors

func (entity *Board) GetSignature() Signature      { return entity.Signature }
func (entity *Thread) GetSignature() Signature     { return entity.Signature }
func (entity *Post) GetSignature() Signature       { return entity.Signature }
func (entity *Vote) GetSignature() Signature       { return entity.Signature }
func (entity *Key) GetSignature() Signature        { return entity.Signature }
func (entity *Truststate) GetSignature() Signature { return entity.Signature }

// OwnerPublicKey accessors

func (entity *Board) GetOwnerPublicKey() string  { return entity.OwnerPublicKey }
func (entity *Thread) GetOwnerPublicKey() string { return entity.OwnerPublicKey }
func (entity *Post) GetOwnerPublicKey() string   { return entity.OwnerPublicKey }
func (entity *Vote) GetOwnerPublicKey() string   { return entity.OwnerPublicKey }

// Heads up, this is slightly different in Key below.
func (entity *Key) GetOwnerPublicKey() string        { return entity.Key }
func (entity *Truststate) GetOwnerPublicKey() string { return entity.OwnerPublicKey }

// Verifiable accessors / setters
func (entity *Board) GetVerified() bool      { return entity.Verified }
func (entity *Thread) GetVerified() bool     { return entity.Verified }
func (entity *Post) GetVerified() bool       { return entity.Verified }
func (entity *Vote) GetVerified() bool       { return entity.Verified }
func (entity *Key) GetVerified() bool        { return entity.Verified }
func (entity *Truststate) GetVerified() bool { return entity.Verified }
func (entity *Address) GetVerified() bool    { return entity.Verified }

func (entity *Board) SetVerified(v bool)      { entity.Verified = v }
func (entity *Thread) SetVerified(v bool)     { entity.Verified = v }
func (entity *Post) SetVerified(v bool)       { entity.Verified = v }
func (entity *Vote) SetVerified(v bool)       { entity.Verified = v }
func (entity *Key) SetVerified(v bool)        { entity.Verified = v }
func (entity *Truststate) SetVerified(v bool) { entity.Verified = v }
func (entity *Address) SetVerified(v bool)    { entity.Verified = v }

// UpdateSignature accessors

func (entity *Board) GetUpdateSignature() Signature      { return entity.UpdateSignature }
func (entity *Thread) GetUpdateSignature() Signature     { return entity.UpdateSignature }
func (entity *Post) GetUpdateSignature() Signature       { return entity.UpdateSignature }
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
func (entity *Thread) GetUpdateProofOfWork() ProofOfWork     { return entity.UpdateProofOfWork }
func (entity *Post) GetUpdateProofOfWork() ProofOfWork       { return entity.UpdateProofOfWork }
func (entity *Vote) GetUpdateProofOfWork() ProofOfWork       { return entity.UpdateProofOfWork }
func (entity *Key) GetUpdateProofOfWork() ProofOfWork        { return entity.UpdateProofOfWork }
func (entity *Truststate) GetUpdateProofOfWork() ProofOfWork { return entity.UpdateProofOfWork }

// Signature accessors

func (entity *Board) GetOwner() Fingerprint  { return entity.Owner }
func (entity *Thread) GetOwner() Fingerprint { return entity.Owner }
func (entity *Post) GetOwner() Fingerprint   { return entity.Owner }
func (entity *Vote) GetOwner() Fingerprint   { return entity.Owner }

// (For below, owner of the entity is itself.)
func (entity *Key) GetOwner() Fingerprint        { return entity.Fingerprint }
func (entity *Truststate) GetOwner() Fingerprint { return entity.Owner }

// RealmId accessors

func (entity *Board) GetRealmId() Fingerprint      { return entity.RealmId }
func (entity *Thread) GetRealmId() Fingerprint     { return entity.RealmId }
func (entity *Post) GetRealmId() Fingerprint       { return entity.RealmId }
func (entity *Vote) GetRealmId() Fingerprint       { return entity.RealmId }
func (entity *Key) GetRealmId() Fingerprint        { return entity.RealmId }
func (entity *Truststate) GetRealmId() Fingerprint { return entity.RealmId }

// EncrContent accessors

func (entity *Board) GetEncrContent() string      { return entity.EncrContent }
func (entity *Thread) GetEncrContent() string     { return entity.EncrContent }
func (entity *Post) GetEncrContent() string       { return entity.EncrContent }
func (entity *Vote) GetEncrContent() string       { return entity.EncrContent }
func (entity *Key) GetEncrContent() string        { return entity.EncrContent }
func (entity *Truststate) GetEncrContent() string { return entity.EncrContent }

// Response styles.

// Response is the interface junction that batch processing functions take and emit. This is the 'internal' communication structure within the backend. It is the big carrier type for the end result of a pull from a remote.
type Response struct {
	Boards      []Board
	Threads     []Thread
	Posts       []Post
	Votes       []Vote
	Keys        []Key
	Addresses   []Address
	Truststates []Truststate

	BoardIndexes      []BoardIndex
	ThreadIndexes     []ThreadIndex
	PostIndexes       []PostIndex
	VoteIndexes       []VoteIndex
	KeyIndexes        []KeyIndex
	AddressIndexes    []AddressIndex
	TruststateIndexes []TruststateIndex

	BoardManifests      []PageManifest
	ThreadManifests     []PageManifest
	PostManifests       []PageManifest
	VoteManifests       []PageManifest
	KeyManifests        []PageManifest
	TruststateManifests []PageManifest
	AddressManifests    []PageManifest

	CacheLinks                []ResultCache
	MostRecentSourceTimestamp Timestamp
}

func (r *Response) Empty() bool {
	return len(r.Boards) == 0 &&
		len(r.Threads) == 0 &&
		len(r.Posts) == 0 &&
		len(r.Votes) == 0 &&
		len(r.Keys) == 0 &&
		len(r.Truststates) == 0 &&
		len(r.Addresses) == 0 &&

		len(r.BoardIndexes) == 0 &&
		len(r.ThreadIndexes) == 0 &&
		len(r.PostIndexes) == 0 &&
		len(r.VoteIndexes) == 0 &&
		len(r.KeyIndexes) == 0 &&
		len(r.TruststateIndexes) == 0 &&
		len(r.AddressIndexes) == 0 &&

		len(r.BoardManifests) == 0 &&
		len(r.ThreadManifests) == 0 &&
		len(r.PostManifests) == 0 &&
		len(r.VoteManifests) == 0 &&
		len(r.KeyManifests) == 0 &&
		len(r.TruststateManifests) == 0 &&
		len(r.AddressManifests) == 0 &&

		len(r.CacheLinks) == 0
}

func (r *Response) Insert(r2 *Response) {
	r.Boards = append(r.Boards, r2.Boards...)
	r.Threads = append(r.Threads, r2.Threads...)
	r.Posts = append(r.Posts, r2.Posts...)
	r.Votes = append(r.Votes, r2.Votes...)
	r.Keys = append(r.Keys, r2.Keys...)
	r.Truststates = append(r.Truststates, r2.Truststates...)
	r.Addresses = append(r.Addresses, r2.Addresses...)

	r.BoardIndexes = append(r.BoardIndexes, r2.BoardIndexes...)
	r.ThreadIndexes = append(r.ThreadIndexes, r2.ThreadIndexes...)
	r.PostIndexes = append(r.PostIndexes, r2.PostIndexes...)
	r.VoteIndexes = append(r.VoteIndexes, r2.VoteIndexes...)
	r.KeyIndexes = append(r.KeyIndexes, r2.KeyIndexes...)
	r.TruststateIndexes = append(r.TruststateIndexes, r2.TruststateIndexes...)
	r.AddressIndexes = append(r.AddressIndexes, r2.AddressIndexes...)

	r.BoardManifests = append(r.BoardManifests, r2.BoardManifests...)
	r.ThreadManifests = append(r.ThreadManifests, r2.ThreadManifests...)
	r.PostManifests = append(r.PostManifests, r2.PostManifests...)
	r.VoteManifests = append(r.VoteManifests, r2.VoteManifests...)
	r.KeyManifests = append(r.KeyManifests, r2.KeyManifests...)
	r.TruststateManifests = append(r.TruststateManifests, r2.TruststateManifests...)
	r.AddressManifests = append(r.AddressManifests, r2.AddressManifests...)

	r.CacheLinks = append(r.CacheLinks, r2.CacheLinks...)

	if r.MostRecentSourceTimestamp < r2.MostRecentSourceTimestamp {
		r.MostRecentSourceTimestamp = r2.MostRecentSourceTimestamp
	} else {
		r.MostRecentSourceTimestamp = r.MostRecentSourceTimestamp
	}
}
