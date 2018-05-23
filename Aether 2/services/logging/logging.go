// Services > Logging
// Logging is the universal logger. This library is responsible for checking whether logging to a file (or to stderr) is enabled, and if so, will process logs as such.

package logging

import (
	"aether-core/services/globals"
	"aether-core/services/toolbox"
	"fmt"
	"log"
	// "runtime"
)

var loggingLevel int
var loglevelInitd bool

func captLogLevel() {
	loggingLevel = globals.BackendConfig.GetLoggingLevel()
	loglevelInitd = true
}

// Log prints to the standard logger.
func Log(level int, input interface{}) {
	if !loglevelInitd {
		captLogLevel()
	}
	if loggingLevel >= level {
		// If print to stdout is enabled, instead of logging, route to stdout. This means it's running in a swarm setup that wants the results that way for collation.
		if globals.BackendTransientConfig.PrintToStdout {
			if globals.BackendTransientConfig.SwarmNodeId != -1 {
				fmt.Printf("%d: %s\n", globals.BackendTransientConfig.SwarmNodeId, input)
			} else {
				fmt.Println(input)
			}
		} else {
			// If not routed to stdout, log normally.
			log.Println(input)
		}
	}
}

func Logf(level int, input string, v ...interface{}) {
	if !loglevelInitd {
		captLogLevel()
	}
	if loggingLevel >= level {
		// If print to stdout is enabled, instead of logging, route to stdout. This means it's running in a swarm setup that wants the results that way for collation.
		if globals.BackendTransientConfig.PrintToStdout {
			if globals.BackendTransientConfig.SwarmNodeId != -1 {
				fmt.Printf("%d: %s\n", globals.BackendTransientConfig.SwarmNodeId, fmt.Sprintf(input, v...))
			} else {
				fmt.Printf(input, v...)
			}
		} else {
			// If not routed to stdout, log normally.
			log.Printf(input, v...)
		}
	}
}

func LogCrash(input interface{}) {
	// If we are already shutting down, do not crash.
	if globals.BackendTransientConfig.ShutdownInitiated {
		return
	}
	log.Println(toolbox.DumpStack())
	log.Fatal(input)
}

func LogObj(level int, objName string, input interface{}) {
	Logf(level, "%s: %#v", objName, input)
}
