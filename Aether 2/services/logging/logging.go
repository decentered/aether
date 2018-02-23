// Services > Logging
// Logging is the universal logger. This library is responsible for checking whether logging to a file (or to stderr) is enabled, and if so, will process logs as such.

package logging

import (
	"aether-core/services/globals"
	"log"
)

// AetherLog prints to the standard logger.
func Log(level int, input interface{}) {
	// TODO: Check whether debug is enabled ONCE at application launch. If so, print to the log file. If not, be a noop.
	if globals.BackendConfig.GetLoggingLevel() >= level {
		log.Println(input)
	}
}
func LogCrash(input interface{}) {
	log.Fatal(input)
}
