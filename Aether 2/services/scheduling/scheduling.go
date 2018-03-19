// Scheduling
// This package provides a scheduler for functions that need to run in certain intervals.

package scheduling

import (
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	// "aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/safesleep"
	"fmt"
	"time"
)

// Schedule runs a function repeatedly until it's asked to stop. Mind that what this does that it calls the function counting after the execution of the prior execution has finished. So if your function takes 5 minutes to run, and you set it to run every 5 minutes, this function will in practice be running every 10 minutes, not 5. This means you don't need to check if two of these functions are running at the same time, there will only ever be one of them running.
func ScheduleRepeat(inputFunction func(), interval time.Duration, initialDelay time.Duration) chan bool {
	// If there's a shutdown in progress, you cannot schedule new tasks.
	everRan := false
	stopChan := make(chan bool)
	go func() {
		for {
			if !everRan {
				err := safesleep.Sleep(initialDelay)
				if err != nil {
					logging.Log(1, fmt.Sprintf("Shutdown in progress. Cancelling this schedule. Initial delay was: %s. Interval was: %s.", initialDelay, interval))
					// This below is so that the channel isn't closed - if it is, the attempt to close this channel will block forever.
					select {
					case <-stopChan:
						return
					}
				}
				everRan = true
			}
			inputFunction()
			select {
			case <-time.After(interval):
			case <-stopChan:
				return
			}
		}
	}()
	return stopChan
}

func ScheduleOnce(inputFunction func(), initialDelay time.Duration) {
	go func() {
		// If there's a shutdown in progress, you cannot schedule new tasks.
		err := safesleep.Sleep(initialDelay)
		if err != nil {
			logging.Log(1, fmt.Sprintf("Shutdown in progress. Cancelling this schedule. Initial Delay was: %s.", initialDelay))
			return
		}
		inputFunction()
	}()
}
