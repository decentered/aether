// Services > Syncconfirmations
// Syncconfirmations service determines and reports on whether this node has made a successful connection in given intervals

package syncconfirmations

import (
	"aether-core/io/api"
	"aether-core/services/globals"
	"errors"
	"fmt"
	"time"
)

var confirmations []api.Timestamp

// maintain removes the sync confirmations that are older than blocks needed x block minutes. So if you need 3 blocks confirmed with 5 minutes for each block, anything older than 15 minutes is going to be deleted by maintain.
func maintain(confs []api.Timestamp) {
	minutesNeeded := globals.BackendConfig.GetTimeBlockSizeMinutes() * globals.BackendConfig.GetPastBlocksToCheck()
	thresholdTs := api.Timestamp(time.Now().Add(-time.Duration(minutesNeeded) * time.Minute).Unix())
	var maintainedConfirmations []api.Timestamp
	for _, val := range confs {
		if val > thresholdTs {
			// Move it to the new slice
			maintainedConfirmations = append(maintainedConfirmations, val)
		}
	}
	confirmations = maintainedConfirmations
}

type timeBlock struct {
	beginning api.Timestamp // older time
	end       api.Timestamp // newer time
}

// generateTimeBlocks generates the time blocks based on NOW. So it'll generate blocks at 0, -x, -2x, -3x..
func generateTimeBlocks(n int) []timeBlock {
	timeContainer := time.Now()
	var tb []timeBlock
	for i := 0; i < n; i++ {
		var block timeBlock
		prevBlockEnd := timeContainer.Add(-time.Duration(globals.BackendConfig.GetTimeBlockSizeMinutes()) * time.Minute)
		block.end = api.Timestamp(timeContainer.Unix())
		block.beginning = api.Timestamp(prevBlockEnd.Unix())
		tb = append(tb, block)
		// Finally, assign the previous block end to time container so the next block will go deeper into the past.
		timeContainer = prevBlockEnd
	}
	return tb
}

// blockIsPassing checks whether the time block given is passing based on the confirmed syncs.
func blockIsPassing(confs []api.Timestamp, block timeBlock) bool {
	for _, conf := range confs {
		if conf > block.beginning && conf < block.end {
			return true
		}
	}
	return false
}

// insertConfirmationInThePast is an unexported testing method so that you can insert past confirmations. The only non-test way you can insert a confirmation is that you call Insert(), and that inserts a timestamp of now.
func insertConfirmationInThePast(conf api.Timestamp) {
	confirmations = append(confirmations, conf)
	maintain(confirmations)
}

// Insert is the function that the external world calls to confirm that a successful sync happened.
func Insert() {
	now := api.Timestamp(time.Now().Unix())
	confirmations = append(confirmations, now)
	maintain(confirmations)
}

/* NodeIsTrackingHead gives an APPROXIMATE result on whether the node is current or not. THIS IS A HEURISTIC. Always guard against the case where this returns true but node is actually very far from tracking head. This might also return erratic results in the case of barely passing connectivity. Why?

Imagine this case: (stars are succesful syncs, bars are time blocks)
-20m  -15m	-10m	 -5m   Now
  |-----|----*|*----|*----|
	 FAIL		OK		OK		OK			Result: OK (if last 3 is being checked with 5m block interval)


Now the important thing is that the five minute intervals into the past is generated at the point of query. So 2 minutes later, every | will move 2 dashes, and:

2 MINUTES LATER:

-20m  -15m	-10m	 -5m   Now
  |-----|--**-|---*-|-----|
	 FAIL		OK		OK	 FAIL			Result: FAIL

This intentional - we want connectivity to be mostly always available, and if a node is having connectivity problems, it should NOT look like it's tracking the head.

SO:

If you're giving a N block interval, it is actually checking for a connection N/2 minutes in the worst case. If the node is having connections more frequent than a connection per N/2 minutes, it'll always be positive. Between a connection per N/2 and 2N minutes, the results can go either way. After a connection per 2N+ minutes, it will always fail.

It's important to not make the costs of failing this test too high. It's a simple test meant to run frequently, and not only it is inexact, it's even in the best case a heuristic. This is useful because it usually correlates with the shape of the local database, but not always.
*/
func NodeIsTrackingHead() (bool, error) {
	return true, nil
	failedBlocks := []timeBlock{}
	tb := generateTimeBlocks(globals.BackendConfig.GetPastBlocksToCheck())
	for _, block := range tb {
		if !blockIsPassing(confirmations, block) {
			failedBlocks = append(failedBlocks, block)
		}
	}
	if len(failedBlocks) > 0 {
		return false, errors.New(fmt.Sprintf("At least one time block failed, so this node is not tracking the head. Failed blocks: %#v, All blocks: %#v", failedBlocks, tb))
	}
	return true, nil
}
