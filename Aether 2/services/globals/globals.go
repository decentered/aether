// Services > Globals
// This file collects all globally accessible entities.

package globals

import (
	"aether-core/services/configstore"
	// "fmt"
	"github.com/asdine/storm"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
)

var FrontendConfig *configstore.FrontendConfig
var FrontendTransientConfig *configstore.FrontendTransientConfig

var BackendConfig *configstore.BackendConfig
var BackendTransientConfig *configstore.BackendTransientConfig

var DbInstance *sqlx.DB
var KvInstance *storm.DB

// GetDbSize gets the size of the database. This is here and not in toolbox because we need to access GetUserDirectory().
func GetDbSize() int {
	dbLoc := filepath.Join(BackendConfig.GetUserDirectory(), "backend", "AetherDB.db")
	fi, _ := os.Stat(dbLoc)
	// get the size
	size := fi.Size() / 1000000
	return int(size)
}
