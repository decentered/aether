// Services > Globals
// This file collects all globally accessible entities.

package globals

import (
	"aether-core/services/configstore"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

var FrontendConfig *configstore.FrontendConfig
var FrontendTransientConfig *configstore.FrontendTransientConfig

var BackendConfig *configstore.BackendConfig
var BackendTransientConfig *configstore.BackendTransientConfig

var DbInstance *sqlx.DB

func GetDbSize() int {
	dbLoc := fmt.Sprintf("%s/AetherDB.db", BackendConfig.GetUserDirectory())
	fi, _ := os.Stat(dbLoc)
	// get the size
	size := fi.Size() / 1000000
	return int(size)
}
