// Services > Globals
// This file collects all constants and user settings, and handles the persistence of the aforementioned.

// This is a temporary file. This should be handled by userconfig.

package globals

import (
	"aether-core/services/signaturing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"time"
)

var KeyPair *ecdsa.PrivateKey
var MarshaledPubKey string
var LastCacheGenerationTimestamp int64
var VerificationEnabled bool

func SetVerificationEnabled(enabled bool) {
	if enabled {
		VerificationEnabled = true
	}
}

func GenerateUserKeyPair() {
	privKey, _ := signaturing.CreateKeyPair()
	KeyPair = privKey
	MarshaledPubKey = hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
}

type EntityPageSizes struct {
	Boards            int
	BoardIndexes      int
	Threads           int
	ThreadIndexes     int
	Posts             int
	PostIndexes       int
	Votes             int
	VoteIndexes       int
	Addresses         int
	AddressIndexes    int
	Keys              int
	KeyIndexes        int
	Truststates       int
	TruststateIndexes int
}

var EntityPageSizesObj EntityPageSizes

// The default base size is 1x (The thread size). At the base size, a page gets 100 entries.
func setEntityPageAndIndexSizes() {
	EntityPageSizesObj.Boards = 500              // 0.2x
	EntityPageSizesObj.BoardIndexes = 4000       // 0.025x
	EntityPageSizesObj.Threads = 100             // 1x
	EntityPageSizesObj.ThreadIndexes = 4000      // 0.025x
	EntityPageSizesObj.Posts = 100               // 1x
	EntityPageSizesObj.PostIndexes = 3000        // 0.033x
	EntityPageSizesObj.Votes = 500               // 0.2x
	EntityPageSizesObj.VoteIndexes = 1000        // 0.1x
	EntityPageSizesObj.Addresses = 4000          // 0.025x
	EntityPageSizesObj.AddressIndexes = 4000     // 0.025x - Address is its own index
	EntityPageSizesObj.Keys = 500                // 0.2x
	EntityPageSizesObj.KeyIndexes = 5000         // 0.02x
	EntityPageSizesObj.Truststates = 4000        // 0.025x
	EntityPageSizesObj.TruststateIndexes = 10000 // 0.01x
	// Every regular page is about 500kb that way.
	// Every index page is about 1mb.
}

type MinPoWStrengthsStruct struct {
	Board            int64
	BoardUpdate      int64
	Thread           int64
	Post             int64
	Vote             int64
	VoteUpdate       int64
	Key              int64
	KeyUpdate        int64
	Truststate       int64
	TruststateUpdate int64
}

var MinPoWStrengths MinPoWStrengthsStruct

func SetMinPoWStrengths(minstr int64) {
	if minstr == 0 {
		minstr = 20
	}
	MinPoWStrengths.Board = minstr
	MinPoWStrengths.BoardUpdate = minstr
	MinPoWStrengths.Thread = minstr
	MinPoWStrengths.Post = minstr
	MinPoWStrengths.Vote = minstr
	MinPoWStrengths.VoteUpdate = minstr
	MinPoWStrengths.Key = minstr
	MinPoWStrengths.KeyUpdate = minstr
	MinPoWStrengths.Truststate = minstr
	MinPoWStrengths.TruststateUpdate = minstr
}

type PoWBailoutTimeStruct struct {
	BailoutTimeSeconds int
}

var PoWBailoutTime PoWBailoutTimeStruct

func SetBailoutTime() {
	PoWBailoutTime.BailoutTimeSeconds = 30
}

var NodeId string
var AddressPort uint16
var AddressType int
var ProtocolVersionMajor int
var ProtocolVersionMinor int
var ProtocolExtensions []string
var ClientVersionMajor int
var ClientVersionMinor int
var ClientVersionPatch int
var ClientName string
var UserDirectory string
var PostResponseExpiryMinutes int
var CachesLocation string
var ConnectionTimeout time.Duration
var TCPConnectTimeout time.Duration
var TLSHandshakeTimeout time.Duration
var PingerPageSize int
var OnlineAddressFinderPageSize int
var DispatcherExclusionsExpiryLiveAddress time.Duration
var DispatcherExclusionsExpiryStaticAddress time.Duration
var LoggingLevel int
var ExternalIp string

/*
Application state: These are set while running. At every start, they will start from their default state given here. Do not change these until you want to test the application already being in that state. (i.e. These are not 'settings' but just the runtime variables, other parts of the code will use these to set variables that won't persist between restarts.)
*/
var TooManyConnections bool // If the system is overloaded, set this bit to true and it'll start to return HTTP 429 Too Many Requests to status endpoint.

/*
Why is this an interface instead of api.Address? Because I can't import address here, it creates a circular reference.
*/
var DispatcherExclusions map[*interface{}]time.Time
var StopLiveDispatcherCycle chan bool
var StopStaticDispatcherCycle chan bool
var StopMatureCacheGenerationCycle chan bool
var StopImmatureCacheGenerationCycle chan bool
var StopAddressScannerCycle chan bool
var StopUPNPCycle chan bool
var AddressesScannerActive bool

func SetApplicationState() {
	TooManyConnections = false
	DispatcherExclusions = make(map[*interface{}]time.Time)
	AddressesScannerActive = false
}

func SetGlobals() {
	// This function is useful until we get the configstore running.
	GenerateUserKeyPair()
	SetMinPoWStrengths(4)
	SetVerificationEnabled(false)
	SetBailoutTime()
	NodeId = "my node id"
	AddressPort = 23420
	AddressType = 2
	ProtocolVersionMajor = 0
	ProtocolVersionMinor = 1
	ProtocolExtensions = []string{"aether"}
	ClientVersionMajor = 2
	ClientVersionMinor = 0
	ClientVersionPatch = 0
	ClientName = "Aether"
	LastCacheGenerationTimestamp = 0
	setEntityPageAndIndexSizes()
	UserDirectory = "/Users/Helios/Dropbox/Aether_Catchall/Aether_Main_Repo/Aether_2/aether-core/userdir"
	PostResponseExpiryMinutes = 30
	CachesLocation = fmt.Sprint(UserDirectory, "/statics/caches/v0")
	ConnectionTimeout = 2 * time.Second
	TCPConnectTimeout = 1 * time.Second
	TLSHandshakeTimeout = 1 * time.Second
	PingerPageSize = 100
	OnlineAddressFinderPageSize = 99
	DispatcherExclusionsExpiryLiveAddress = 5 * time.Minute
	DispatcherExclusionsExpiryStaticAddress = 72 * time.Hour
	LoggingLevel = 0
	SetApplicationState()

}
