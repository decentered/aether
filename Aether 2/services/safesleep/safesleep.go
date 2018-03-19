// Services > SafeSleep
// Safesleep provides a sleep function that does not prevent application shutdown.

package safesleep

import (
	"aether-core/services/globals"
	"aether-core/services/logging"
	"errors"
	"fmt"
	"time"
)

func Sleep(dur time.Duration) error {
	// Split time.duration into 10 second intervals
	sec := int(dur.Seconds())
	var blocks int
	if sec <= 10 {
		// Sleep for the exact time if it's less than or exactly 10 seconds.
		time.Sleep(dur)
		if globals.BackendTransientConfig.ShutdownInitiated {
			logging.Log(2, fmt.Sprintf("Shutdown in progress. SafeSleep is exiting. Duration was: %s", dur))
			return errors.New("The application is shutting down. Please exit gracefully.")
		} else {
			return nil
		}
	} else {
		blocks = sec / 10
	}
	for i := 0; i < blocks; i++ {
		time.Sleep(time.Duration(10) * time.Second)
		if globals.BackendTransientConfig.ShutdownInitiated {
			logging.Log(2, fmt.Sprintf("Shutdown in progress. SafeSleep is exiting. Duration was: %s", dur))
			return errors.New("The application is shutting down. Please exit gracefully.")
		}
	}
	return nil
}
