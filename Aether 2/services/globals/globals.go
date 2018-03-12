// Services > Globals
// This file collects all constants and user settings, and handles the persistence of the aforementioned.

// This is a temporary file. This should be handled by userconfig.

package globals

import (
	pb "aether-core/backend/metrics/proto"
	"aether-core/services/configstore"
	"github.com/jmoiron/sqlx"
	"time"
)

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
var StopAddressScannerCycle chan bool
var StopUPNPCycle chan bool
var StopCacheGenerationCycle chan bool
var AddressesScannerActive bool
var LiveDispatchRunning bool
var StaticDispatchRunning bool

var FrontendConfig *configstore.FrontendConfig
var FrontendTransientConfig *configstore.FrontendTransientConfig

var BackendConfig *configstore.BackendConfig
var BackendTransientConfig *configstore.BackendTransientConfig

var DbInstance *sqlx.DB
var CurrentMetricsPage pb.Metrics
