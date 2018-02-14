// Scheduling
// This package provides a scheduler for functions that need to run in certain intervals.

package scheduling

import (
	// "fmt"
	"time"
)

// Schedule runs a function repeatedly until it's asked to stop. Mind that what this does that it calls the function counting after the execution of the prior execution has finished. So if your function takes 5 minutes to run, and you set it to run every 5 minutes, this function will in practice be running every 10 minutes, not 5. This means you don't need to check if two of these functions are running at the same time, there will only ever be one of them running.
func Schedule(inputFunction func(), interval time.Duration) chan bool {
	stopChan := make(chan bool)
	go func() {
		for {
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
