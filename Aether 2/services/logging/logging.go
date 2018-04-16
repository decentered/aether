// Services > Logging
// Logging is the universal logger. This library is responsible for checking whether logging to a file (or to stderr) is enabled, and if so, will process logs as such.

package logging

import (
	"aether-core/services/globals"
	"fmt"
	"log"
	"runtime"
)

func trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	result := fmt.Sprintf("%s,:%d %s", frame.File, frame.Line, frame.Function)
	return result
}

// Log prints to the standard logger.
func Log(level int, input interface{}) {
	// TODO: Check whether debug is enabled ONCE at application launch. If so, print to the log file. If not, be a noop.
	if globals.BackendConfig.GetLoggingLevel() >= level {
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
func LogCrash(input interface{}) {
	log.Println(globals.DumpStack())
	log.Fatal(input)
}
