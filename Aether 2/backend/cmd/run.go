package cmd

import (
	"aether-core/backend/dispatch"
	"aether-core/backend/responsegenerator"
	"aether-core/backend/server"
	// "aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/ports"
	"aether-core/services/scheduling"
	// "aether-core/services/upnp"
	// "fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	var loggingLevel int
	cmdRun.Flags().IntVarP(&loggingLevel, "logginglevel", "", 0, "Global logging level of the app.")
	cmdRoot.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Start a full-fledged Mim node that tracks the network and responds to requests.",
	Long: `Start a full-fledged Mim node that tracks the network and responds to requests. This is the default entry point if you want to use this Mim backend normally.

This will do three main things:

- Actively start fetching from other nodes and constructing the network head in the local machine
- Expose a local API for the frontend app to peruse, so that the content fetched over the network is available for consumption
- Expose an API to the external world that serves the data this computer has under the rules set by the Mim protocol.
`,
	Run: func(cmd *cobra.Command, args []string) {
		EstablishConfigs(cmd)
		showIntro() // This isn't first because it needs configs to show app version.
		persistence.CreateDatabase()
		persistence.CheckDatabaseReady()
		startSchedules()
		// Allocate the dispatcher exclusions list.
		server.Serve()
		shutdown()
	},
}

func startSchedules() {
	// logging.Logf(1, "UserDir: %v", globals.BackendConfig.GetUserDirectory())
	logging.Log(1, "Setting up cyclical tasks is starting.")
	defer logging.Log(1, "Setting up cyclical tasks is complete.")
	/*
		Ordered by initial delay:
		Verify external port: T+0 	(immediately)
		Neighbourhood dispatcher 			T+0: 	(immediately)

		UPNP Port mapper: 		T+0 	(immediately)
		Explorer dispatcher 			T+10:
		Address Scanner: 			T+15m
		Cache generator: 			T+30m
	*/
	// Before doing anything, you need to validate the external port. This function takes a second or so, and it needs to block the runtime execution because if two routines call it separately, it causes a race condition. After the first initialisation, however, this function becomes safe for concurrent use.
	ports.VerifyExternalPort()
	// UPNP tries to port map every 10 minutes. TODO reenable
	// globals.BackendTransientConfig.StopUPNPCycle = scheduling.ScheduleRepeat(func() { upnp.MapPort() }, 10*time.Minute, time.Duration(0))
	dispatch.Bootstrap() // This will run only if needed.
	// The dispatcher that seeks live nodes runs every minute.
	globals.BackendTransientConfig.StopNeighbourhoodCycle = scheduling.ScheduleRepeat(func() { dispatch.NeighbourWatch() }, 20*time.Second, time.Duration(0))
	globals.BackendTransientConfig.StopExplorerCycle = scheduling.ScheduleRepeat(func() { dispatch.Explore() }, 10*time.Minute, time.Duration(10)*time.Minute)

	// Address scanner goes through all prior unconnected addresses and attempts to connect to them to establish a relationship. It starts 30 minutes after a node is started, so that the node will actually have a chance to collect some addresses to check.
	globals.BackendTransientConfig.StopAddressScannerCycle = scheduling.ScheduleRepeat(func() { dispatch.AddressScanner() }, 2*time.Hour, time.Duration(15)*time.Minute)
	// Attempt cache generation every hour, but it will be pre-empted if the last cache generation is less than 23 hours old, and if the node is not tracking the head. So that this will run effectively every day, only.
	globals.BackendTransientConfig.StopCacheGenerationCycle = scheduling.ScheduleRepeat(func() { responsegenerator.GenerateCaches() }, 1*time.Hour, time.Duration(30)*time.Minute)
}

func shutdown() {
	logging.Log(1, "Shutdown initiated. Stopping all scheduled tasks and routines...")
	globals.BackendTransientConfig.ShutdownInitiated = true
	globals.BackendTransientConfig.StopNeighbourhoodCycle <- true // Send true through the channel to stop the dispatch.
	globals.BackendTransientConfig.StopExplorerCycle <- true      // Send true through the channel to stop the dispatch.
	globals.BackendTransientConfig.StopAddressScannerCycle <- true
	// globals.BackendTransientConfig.StopUPNPCycle <- true // upnp is disabled, reenable both when it's back
	globals.BackendTransientConfig.StopCacheGenerationCycle <- true
	logging.Log(1, "Waiting 5 seconds to let DB close gracefully...")
	time.Sleep(time.Duration(5) * time.Second) // Wait 5 seconds to let DB tasks complete.
	globals.DbInstance.Close()
	defer func() {
		// The functions that access DB can panic after the DB is closed. But after DB is closed, we don't care - the DB is out of harm's way and the only state that remains at this phase is the transient state, and that's going to be wiped out a few nanoseconds later. Recover from any panics.
		recResult := recover()
		if recResult != nil {
			logging.Logf(1, "Recovered from a panic at the end of the shutdown after DB close. A panic here can be caused by a process being interrupted. In most cases, it's normal behaviour and nothing to worry about. Panic'd error: %#v", recResult)
		}
	}()
	// We delete at shutdown and at boot, just in case deletion at shutdown didn't work.
	globals.BackendTransientConfig.POSTResponseRepo.DeleteAllFromDisk()
	logging.Log(1, "Shutdown is complete. Bye.")
	os.Exit(0)
}
