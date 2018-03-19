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
		establishConfigs(cmd)
		showIntro() // This isn't first because it needs configs to show app version.
		persistence.CreateDatabase()
		persistence.CheckDatabaseReady()
		startSchedules()
		// Allocate the dispatcher exclusions list.
		globals.BackendTransientConfig.DispatcherExclusions = make(map[*interface{}]time.Time)
		server.Serve()
		shutdown()
	},
}

func startSchedules() {
	logging.Log(1, "Setting up cyclical tasks is starting.")
	defer logging.Log(1, "Setting up cyclical tasks is complete.")
	/*
		Ordered by initial delay:
		Verify external port: T+0 	(immediately)
		Live dispatcher 			T+0: 	(immediately)
		UPNP Port mapper: 		T+0 	(immediately)
		Address Scanner: 			T+10m
		Static dispatcher: 		T+20m
		Cache generator: 			T+30m
	*/
	// Before doing anything, you need to validate the external port. This function takes a second or so, and it needs to block the runtime execution because if two routines call it separately, it causes a race condition. After the first initialisation, however, this function becomes safe for concurrent use.
	ports.VerifyExternalPort()
	// The dispatcher that seeks live nodes runs every minute.
	globals.BackendTransientConfig.StopLiveDispatcherCycle = scheduling.ScheduleRepeat(func() { dispatch.Dispatcher(2) }, 20*time.Second, time.Duration(0))
	// UPNP tries to port map every 10 minutes. TODO reenable
	// globals.BackendTransientConfig.StopUPNPCycle = scheduling.ScheduleRepeat(func() { upnp.MapPort() }, 10*time.Minute, time.Duration(0))
	// Address scanner goes through all prior unconnected addresses and attempts to connect to them to establish a relationship. It starts 30 minutes after a node is started, so that the node will actually have a chance to collect some addresses to check.
	globals.BackendTransientConfig.StopAddressScannerCycle = scheduling.ScheduleRepeat(func() { dispatch.AddressScanner() }, 2*time.Hour, time.Duration(10)*time.Minute)
	// The dispatcher that seeks static nodes runs every hour.
	globals.BackendTransientConfig.StopStaticDispatcherCycle = scheduling.ScheduleRepeat(func() { dispatch.Dispatcher(255) }, 10*time.Minute, time.Duration(20)*time.Minute)
	// Attempt cache generation every hour, but it will be pre-empted if the last cache generation is less than 23 hours old, and if the node is not tracking the head. So that this will run effectively every day, only.
	globals.BackendTransientConfig.StopCacheGenerationCycle = scheduling.ScheduleRepeat(func() { responsegenerator.GenerateCaches() }, 1*time.Hour, time.Duration(30)*time.Minute)
}

func shutdown() {
	logging.Log(1, "Shutdown initiated. Stopping all scheduled tasks and routines...")
	globals.BackendTransientConfig.ShutdownInitiated = true
	globals.BackendTransientConfig.StopLiveDispatcherCycle <- true // Send true through the channel to stop the dispatch.
	globals.BackendTransientConfig.StopStaticDispatcherCycle <- true
	globals.BackendTransientConfig.StopAddressScannerCycle <- true
	// globals.BackendTransientConfig.StopUPNPCycle <- true // upnp is disabled, reenable both when it's back
	globals.BackendTransientConfig.StopCacheGenerationCycle <- true
	globals.DbInstance.Close()
	logging.Log(1, "Shutdown is complete. Bye.")
	os.Exit(0)
}
