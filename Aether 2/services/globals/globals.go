// Services > Globals
// This file collects all globally accessible entities.

package globals

import (
	"aether-core/services/configstore"
	"github.com/jmoiron/sqlx"
)

var FrontendConfig *configstore.FrontendConfig
var FrontendTransientConfig *configstore.FrontendTransientConfig

var BackendConfig *configstore.BackendConfig
var BackendTransientConfig *configstore.BackendTransientConfig

var DbInstance *sqlx.DB
