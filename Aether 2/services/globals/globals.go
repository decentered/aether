// Services > Globals
// This file collects all globally accessible entities.

package globals

import (
	"aether-core/services/configstore"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
	"runtime"
	"strconv"
)

var FrontendConfig *configstore.FrontendConfig
var FrontendTransientConfig *configstore.FrontendTransientConfig

var BackendConfig *configstore.BackendConfig
var BackendTransientConfig *configstore.BackendTransientConfig

var DbInstance *sqlx.DB

func Round(x, unit float64) float64 {
	r := math.Round(x/unit) * unit
	formatted, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", r), 64)
	return formatted
}

func DumpStack() string {
	_, file, line, _ := runtime.Caller(1)
	_, file2, line2, _ := runtime.Caller(2)
	_, file3, line3, _ := runtime.Caller(3)
	_, file4, line4, _ := runtime.Caller(4)
	_, file5, line5, _ := runtime.Caller(5)
	return fmt.Sprintf("\nSTACK TRACE\n%s:%d\n%s:%d\n%s:%d \n%s:%d \n%s:%d\n",
		file, line, file2, line2, file3, line3, file4, line4, file5, line5)
}
