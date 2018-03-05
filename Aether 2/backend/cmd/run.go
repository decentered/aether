package cmd

import (
	"aether-core/backend/dispatch"
	"aether-core/backend/responsegenerator"
	"aether-core/backend/server"
	// "aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/scheduling"
	"aether-core/services/upnp"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	var loggingLevel int
	cmdRun.Flags().IntVarP(&loggingLevel, "logginglevel", "l", 0, "Global logging level of the app.")
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
		startSchedules()
		// Allocate the dispatcher exclusions list.
		globals.DispatcherExclusions = make(map[*interface{}]time.Time)
		server.Serve()
		shutdown()
	},
}

// func startup() {
// 	// Configs are established at the base app level, not in /run command here.
// 	// Intro is shown on the base config level as well.
// 	// Flags are set in the init of this file, and read in the main Run: function.
// 	persistence.CreateDatabase()
// 	// TEST Insert the localhost data.

// 	// // TODO MOVE THESE TO CONFIGSTORE
// 	// var addrLocal api.Address
// 	// addrLocal.Location = "127.0.0.1"
// 	// addrLocal.Sublocation = ""
// 	// addrLocal.LocationType = 4
// 	// addrLocal.Port = 8001
// 	// addrLocal.LastOnline = 1111111
// 	// addrLocal.Protocol.VersionMajor = 1
// 	// addrLocal.Protocol.VersionMinor = 1
// 	// addrLocal.Protocol.Subprotocols = []api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}}
// 	// addrLocal.Client.VersionMajor = 1
// 	// addrLocal.Client.VersionMinor = 1
// 	// addrLocal.Client.VersionPatch = 1
// 	// addrLocal.Client.ClientName = "Aether"
// 	// persistence.BatchInsert([]interface{}{addrLocal})
// 	// dispatch.Sync(addrLocal)
// 	startSchedules()
// 	logging.Log(1, "Startup complete.")
// }

func startSchedules() {
	logging.Log(1, "Setting up cyclical tasks is starting.")
	defer logging.Log(1, "Setting up cyclical tasks is complete.")
	// The dispatcher that seeks live nodes runs every minute.
	globals.StopLiveDispatcherCycle = scheduling.Schedule(func() { dispatch.Dispatcher(2) }, 1*time.Second)
	// The dispatcher that seeks static nodes runs every hour.
	globals.StopStaticDispatcherCycle = scheduling.Schedule(func() { dispatch.Dispatcher(255) }, 1*time.Minute)
	// Address scanner goes through all prior unconnected addresses and attempts to connect to them to establish a relationship.
	globals.StopAddressScannerCycle = scheduling.Schedule(func() { dispatch.AddressScanner() }, 6*time.Hour)
	// UPNP tries to port map every 10 minutes.
	globals.StopUPNPCycle = scheduling.Schedule(func() { upnp.MapPort() }, 10*time.Minute)
	// Attempt cache generation every hour, but it will be pre-empted if the last cache generation is less than 23 hours old, so that this will run effectively every day, only.
	globals.StopCacheGenerationCycle = scheduling.Schedule(func() { responsegenerator.GenerateCaches() }, 1*time.Hour)

	// time.AfterFunc(5*time.Second, func() {
	// })

	/*
	   For cache generation, the logic is like this:
	   - Start a schedule that checks every 5 minutes if the node is mature
	   - If node is mature, start the mature cycle and stop the immature cycle.
	*/
	// maturityChecker := func() {
	//  mature, err := persistence.LocalNodeIsMature()
	//  if err != nil {
	//    logging.LogCrash(err)
	//  }
	//  if mature {
	//    // If the node is mature, stop the immature cycle and start the mature.
	//    logging.Log(1, "The local node is as of now mature. Stopping the maturity check scheduling and starting the cache generation schedule")
	//    globals.StopMatureCacheGenerationCycle = scheduling.Schedule(func() { responsegenerator.GenerateCaches() }, 6*time.Hour)
	//    globals.StopImmatureCacheGenerationCycle <- true
	//  }
	// }
	// globals.StopImmatureCacheGenerationCycle = scheduling.Schedule(maturityChecker, 5*time.Minute)

}

func shutdown() {
	logging.Log(1, "Shutdown initiated.")
	fmt.Println("Shutdown initiated.")
	globals.StopLiveDispatcherCycle <- true // Send true through the channel to stop the dispatch.
	globals.StopStaticDispatcherCycle <- true
	globals.StopAddressScannerCycle <- true
	globals.StopUPNPCycle <- true
	// mature, err := persistence.LocalNodeIsMature()
	// if err != nil {
	//  logging.LogCrash(err)
	// }
	// if mature {
	//  globals.StopMatureCacheGenerationCycle <- true
	// } else {
	//  globals.StopImmatureCacheGenerationCycle <- true
	// }
	logging.Log(1, "Shutdown is complete.")
	fmt.Println("Shutdown is complete. Bye.")
	os.Exit(0)
}
