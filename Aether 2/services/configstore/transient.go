// Services > ConfigStore
// This module handles saving and reading values from a config user file.

package configstore

import (
	// "aether-core/services/fingerprinting"
	// "aether-core/services/randomhashgen"
	// "aether-core/services/signaturing"
	// "aether-core/services/toolbox"
	// "crypto/ecdsa"
	// "crypto/elliptic"
	// "crypto/x509"
	// "encoding/hex"
	pb "aether-core/backend/metrics/proto"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/fatih/color"
	// "crypto/elliptic"
	// "encoding/hex"
	// cdir "github.com/shibukawa/configdir"
	// "golang.org/x/crypto/ed25519"
	// "log"
	// "runtime"
	"sync"
	"time"
)

// TRANSIENT CONFIG

// These are the items that are set in runtime, and do not change until the application closes. This is different from the application state in the way that they're set-once for the runtime.

// These do not have getters and setters.

var Btc BackendTransientConfig
var Ftc FrontendTransientConfig

// Backend

// Default entity versions for this version of the app. This is not user adjustable.

const (
	defaultBoardEntityVersion      = 1
	defaultThreadEntityVersion     = 1
	defaultPostEntityVersion       = 1
	defaultVoteEntityVersion       = 1
	defaultKeyEntityVersion        = 1
	defaultTruststateEntityVersion = 1
	defaultAddressEntityVersion    = 1
)

type EntityVersions struct {
	Board      int
	Thread     int
	Post       int
	Vote       int
	Key        int
	Truststate int
	Address    int
}

/*
#### NONCOMMITTED ITEMS

## PermConfigReadOnly
When enabled, this prevents anything from saved into the config. This value itself is NOT saved into the config, so when the application restarts, this value is reset to false. This is useful in the case that you provide flags to the executable, but you don't want the values in the flags to be permanently saved into the config file. Any flags being provided into the executable will set this to true, therefore any runs with flags will effectively treat the config as read-only.

## AppIdentifier
This is the name of the app as registered to the operating system. This is useful to have here, because what we can do is we can vary this number in the swarm testing (petridish) and each of these nodes will act like a network in a single local machine, each with their own databases and different config files.

## OrgIdentifier
Same as above, but it's probably best to keep it under the same org name just to keep the local machine clean.

## PrintToStdout
This is useful because the logging things the normal kind does not pass the output to the swarm test orchestrator. This flag being enabled routes the logs to stdout so that the orchestrator can show it.

## MetricsDebugMode
This being enabled temporarily makes this node send much more detailed metrics more frequently, so that network connectivity issues can be debugged. This is a transient config on purpose, so that this cannot be enabled permanently. If a frontend connects to a backend with debug mode enabled, it has to show a warning to its user that says this backend node has debugging enabled, and only connect if the user agrees. Mind that the backend doesn't have to be truthful about whether it has the debug mode on. Having this mode on does not immediately compromise the frontend's privacy / identity, but the longer the frontend stays on that backend and the more actions a user commits, the higher the likelihood.

## ExternalPortVerified
Whether the port that was in the config was actually checked to be free and clear. This is important because we'll check once before the server starts to run, and when it starts, that port will no longer be available, and will start to return 'not available'. That will make all subsequent checks fail and that will trigger the port to be moved to a port that is free - but not bound to any server, since the server is bound to the old port, and that in fact is the reason the checks return false.

## SwarmNodeId
This is the number that this specific node will route to the main swarm orchestrator when it's reporting logs. Make sure that the App identifier (Usually in the format of "Aether-N") matches this number N, or it can be confusing.

## ShutdownInitiated
This is set when the shutdown of the backend service is initiated. The processes that take a long time to return should be checking this value periodically, and if it is set, they should stop whatever they're doing and do a graceful shutdown.

## DispatcherExclusions
This is the temporary exclusions for the dispatcher. When you connect to a node, that node is placed in the exclusions list for a while, so that you don't repeatedly keep connecting back to that node again.

## StopStaticDispatcherCycle
This is the channel to send the message to when you want to stop the static dispatcher repeated task.

## StopAddressScannerCycle
This is the channel to send the message to when you want to stop the address scanner repeated task.

## StopUPNPCycle
This is the channel to send the message to when you want to stop the UPNP mapper repeated task.

## StopCacheGenerationCycle
This is the channel to send the message to when you want to stop the cache generator repeated task.

## AddressesScannerActive
This is the mutex that gets activated when the address scanner is active, so that it cannot be triggered twice at the same time.

## LiveDispatchRunning
This is the mutex that gets activated when the live dispatcher is active, so that it cannot be triggered twice at the same time.

## StaticDispatchRunning
This is the mutex that gets activated when the static dispatcher is active, so that it cannot be triggered twice at the same time.

## CurrentMetricsPage
This is the current metrics struct that we are building to send to the metrics server, if enabled.

## ConfigMutex
This is the mutex that prevents configuration from being written from multiple places.

## FingerprintCheckEnabled
Determines whether the entities coming over from the wire are fingerprint-checked for integrity.

## SignatureCheckEnabled
Determines whether the entities coming over from the wire are signature-checked for ownership.

## ProofOfWorkCheckEnabled
Determines whether the entities coming over from the wire are PoW-checked for anti-spam.

## PageSignatureCheckEnabled
Determines whether the pages (entity containers) coming over from the wire are signature-checked for integrity.

## EntityVersions
These are the versions of the entities that we can issue in this version of the app. Mind that this is for issuance, not for acceptance - we should still accept older versions gracefully.

# POSTResponseRepo
This is the repository that we keep our post responses in, so that they can be reused. This resets at every restart.

# NeighboursList
Our list of neighbours that we are checking in with at given intervals.
*/

type BackendTransientConfig struct {
	PermConfigReadOnly        bool
	AppIdentifier             string
	OrgIdentifier             string
	PrintToStdout             bool
	MetricsDebugMode          bool
	TooManyConnections        bool
	ExternalPortVerified      bool
	SwarmNodeId               int
	ShutdownInitiated         bool
	DispatcherExclusions      map[*interface{}]time.Time
	StopLiveDispatcherCycle   chan bool
	StopStaticDispatcherCycle chan bool
	StopAddressScannerCycle   chan bool
	StopUPNPCycle             chan bool
	StopCacheGenerationCycle  chan bool
	AddressesScannerActive    sync.Mutex
	ActiveOutbound            sync.Mutex
	LiveDispatchRunning       bool
	StaticDispatchRunning     bool
	CurrentMetricsPage        pb.Metrics
	ConfigMutex               *sync.Mutex
	FingerprintCheckEnabled   bool
	SignatureCheckEnabled     bool
	ProofOfWorkCheckEnabled   bool
	PageSignatureCheckEnabled bool
	EntityVersions            EntityVersions
	POSTResponseRepo          POSTResponseRepo // empty at start, empty at every app start
	NeighboursList            NeighboursList
}

// Set transient backend config defaults. Only need to set defaults that are not the type default.

// Mind that if you somehow manage to call something before SetDefaults is called, it will return its zero value without warning. This transient config does not have a Initialised gate that we can check, because adding that gate would have us convert everything in this place to getters / setters. We might do that in the future, but the point of BTC/FTC is that these are the things where the default value of the thing is the empty value of that variable.

// The problem here is that the default value of the field being empty value of that variable type and configs that need to be transient don't exactly match. So we will probably eventually move to a get/set model where it checks for init.

func (config *BackendTransientConfig) SetDefaults() {
	config.AppIdentifier = "Aether"
	config.OrgIdentifier = "Air Labs"
	config.ConfigMutex = &sync.Mutex{}
	config.FingerprintCheckEnabled = true
	config.SignatureCheckEnabled = true
	config.ProofOfWorkCheckEnabled = true
	config.PageSignatureCheckEnabled = true
	ev := EntityVersions{
		Board:      defaultBoardEntityVersion,
		Thread:     defaultThreadEntityVersion,
		Post:       defaultPostEntityVersion,
		Vote:       defaultVoteEntityVersion,
		Key:        defaultKeyEntityVersion,
		Truststate: defaultTruststateEntityVersion,
		Address:    defaultAddressEntityVersion,
	}
	config.EntityVersions = ev
}

// Frontend

type FrontendTransientConfig struct {
	PermConfigReadOnly bool
	MetricsDebugMode   bool
	ConfigMutex        *sync.Mutex
}

// Set transient frontend config defaults

func (config *FrontendTransientConfig) SetDefaults() {
	config.PermConfigReadOnly = false
	config.ConfigMutex = &sync.Mutex{}
}
