// Unlike others, this test package is not named syncconfirmations_test because we need access to internals of the package to be able to check proper behaviour.

package syncconfirmations

import (
	"aether-core/io/api"
	// "aether-core/services/syncconfirmations"
	"aether-core/services/globals"
	// "fmt"
	"os"
	"testing"
	"time"
)

// Infrastructure, setup and teardown

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func createConfirmation(minAgo int) api.Timestamp {
	return api.Timestamp(time.Now().Add(-time.Duration(minAgo) * time.Minute).Unix())
}

func generateConfirmations() []api.Timestamp {
	confs := []api.Timestamp{}
	TwoMinutesAgoConf := createConfirmation(2)
	ThreeMinutesAgoConf := createConfirmation(3)
	FiveMinutesAgoConf := createConfirmation(5)
	SevenMinutesAgoConf := createConfirmation(7)
	EightMinutesAgoConf := createConfirmation(8)
	TenMinutesAgoConf := createConfirmation(10)
	TwelveMinutesAgoConf := createConfirmation(12)
	FourteenMinutesAgoConf := createConfirmation(14)
	SixteenMinutesAgoConf := createConfirmation(16)
	ThirtyMinutesAgoConf := createConfirmation(30)
	confs = append(confs, TwoMinutesAgoConf, ThreeMinutesAgoConf, FiveMinutesAgoConf, SevenMinutesAgoConf, EightMinutesAgoConf, TenMinutesAgoConf, TwelveMinutesAgoConf, FourteenMinutesAgoConf, SixteenMinutesAgoConf, ThirtyMinutesAgoConf)
	return confs
}

var testConfs []api.Timestamp
var thresholdMinutes int

func setup() {
	globals.SetGlobals()
	testConfs = generateConfirmations()
	// Fix these values as part of the test harness.
	globals.BackendConfig.SetTimeBlockSizeMinutes(5)
	globals.BackendConfig.SetPastBlocksToCheck(3)
	thresholdMinutes = globals.BackendConfig.GetTimeBlockSizeMinutes() * globals.BackendConfig.GetPastBlocksToCheck()
}

func teardown() {
}

// Check maintainer behaviour.
func TestInsertViaMaintainer_Success(t *testing.T) {
	for _, val := range testConfs {
		insertConfirmationInThePast(val)
	}

	for _, testConf := range testConfs {
		var confirmed bool
		for _, savedConf := range confirmations {
			if testConf == savedConf {
				confirmed = true
				break
			}
		}
		// If the confirmation is not saved, check whether if it was older than threshold minutes. If not, it's a genuine failure. If it is, the maintain did its job and removed it, so it's a success.
		if !confirmed && testConf > createConfirmation(thresholdMinutes) {
			t.Errorf("Not all confirmations are inserted")
		}
	}
}

func TestMaintainerCutoff_Success(t *testing.T) {
	confirmations = []api.Timestamp{}
	// Add a TS that is older than the allowed threshold
	tsThatShouldBeRemoved := createConfirmation(thresholdMinutes + 1)
	tsThatShouldNotBeRemoved := createConfirmation(thresholdMinutes - 1)
	insertConfirmationInThePast(tsThatShouldBeRemoved)
	insertConfirmationInThePast(tsThatShouldNotBeRemoved)
	// Check whether that confirmation is within the bucket.
	for _, val := range confirmations {
		if val == tsThatShouldBeRemoved {
			t.Errorf("This confirmation should not have been here because it's older than threshold minutes, but it is. Confirmation: %#v, All confirmations: %#v", val, confirmations)
		}
	}
}

// Check the external Insert API.
func TestInsert_Success(t *testing.T) {
	confirmations = []api.Timestamp{}
	TwoMinutesAgoConf := createConfirmation(2)
	ThreeMinutesAgoConf := createConfirmation(3)
	insertConfirmationInThePast(TwoMinutesAgoConf)
	insertConfirmationInThePast(ThreeMinutesAgoConf)
	Insert() // Insert a NOW confirmation.
	if len(confirmations) != 3 {
		t.Errorf("NOW confirmation insert did not succeed. All confirmations: %#v", confirmations)
	}
}

// Check overall tracking head endpoint.
func TestNodeTrackingHead_Success(t *testing.T) {
	confirmations = []api.Timestamp{}
	testConfs = generateConfirmations()
	for _, val := range testConfs {
		insertConfirmationInThePast(val)
	}
	tracking, err := NodeIsTrackingHead()
	if err != nil {
		t.Errorf("Node should be tracking the head, but it returned an error. All confirmations: %#v, Error: %#v, Tracking: %#v", confirmations, err, tracking)
	}
}

func TestNodeTrackingHead_Failure(t *testing.T) {
	confirmations = []api.Timestamp{}
	// First time block will fail (0 to -5)
	SixMinutesAgoConf := createConfirmation(6)
	// Second wil succeed (-5 to -10)
	TwelveMinutesAgoConf := createConfirmation(9)
	// Third will succeed (-10 to -15)
	FifteenMinutesAgoConf := createConfirmation(14)
	// Third will succeed (-10 to -15)
	// Result will fail
	testConfs = []api.Timestamp{SixMinutesAgoConf, TwelveMinutesAgoConf, FifteenMinutesAgoConf}
	for _, val := range testConfs {
		insertConfirmationInThePast(val)
	}
	_, err := NodeIsTrackingHead()
	// fmt.Println(tracking, err)
	if err == nil {
		t.Errorf("Node should not be tracking the head, returned but returned true. All confirmations: %#v", confirmations)
	}
}
