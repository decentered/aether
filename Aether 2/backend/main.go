// Aether Core Backend
// Application start. Starts all the necessary components:
// Database, API, UI, UPNPC, Dispatcher, Validator.

package main

import (
	// "aether-core/backend/dispatch"
	"aether-core/backend/responsegenerator"
	"aether-core/backend/server"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	// "aether-core/services/verify"
	// "crypto/ecdsa"
	"aether-core/services/logging"
	// "aether-core/services/scheduling"
	// "aether-core/services/upnp"
	"flag"
	"fmt"
	"os"
	// "time"
)

func StartSchedules() {
	logging.Log(1, "Setting up cyclical tasks is starting.")
	defer logging.Log(1, "Setting up cyclical tasks is complete.")
	// The dispatcher that seeks live nodes runs every minute.
	// globals.StopLiveDispatcherCycle = scheduling.Schedule(func() { dispatch.Dispatcher(2) }, 1*time.Second)
	// // The dispatcher that seeks static nodes runs every hour.
	// globals.StopStaticDispatcherCycle = scheduling.Schedule(func() { dispatch.Dispatcher(255) }, 1*time.Minute)
	// // Address scanner goes through all prior unconnected addresses and attempts to connect to them to establish a relationship.
	// globals.StopAddressScannerCycle = scheduling.Schedule(func() { dispatch.AddressScanner() }, 6*time.Hour)
	// // UPNP tries to port map every 10 minutes.
	// globals.StopUPNPCycle = scheduling.Schedule(func() { upnp.MapPort() }, 10*time.Minute)
	fmt.Println("Caches are starting to be generated...")
	responsegenerator.GenerateCaches()
	fmt.Println("Caches generation is complete.")
	// time.AfterFunc(5*time.Second, func() {
	// })

	/*
		For cache generation, the logic is like this:
		- Start a schedule that checks every 5 minutes if the node is mature
		- If node is mature, start the mature cycle and stop the immature cycle.
	*/
	// maturityChecker := func() {
	// 	mature, err := persistence.LocalNodeIsMature()
	// 	if err != nil {
	// 		logging.LogCrash(err)
	// 	}
	// 	if mature {
	// 		// If the node is mature, stop the immature cycle and start the mature.
	// 		logging.Log(1, "The local node is as of now mature. Stopping the maturity check scheduling and starting the cache generation schedule")
	// 		globals.StopMatureCacheGenerationCycle = scheduling.Schedule(func() { responsegenerator.GenerateCaches() }, 6*time.Hour)
	// 		globals.StopImmatureCacheGenerationCycle <- true
	// 	}
	// }
	// globals.StopImmatureCacheGenerationCycle = scheduling.Schedule(maturityChecker, 5*time.Minute)

}

func ReadFlags() {
	logIntPtr := flag.Int("logginglevel", 0, "Determines the logging level of the application. Logging level 1 is core messages, 2 is everything. Mind that the more logging you have enabled, the more the app will slow down.")
	flag.Parse()
	globals.LoggingLevel = *logIntPtr
}

func ShowIntro() {
	fmt.Println(`
	                1ttfffLLLLLLLLLffft
	            11111ttfffLLLLLLLLLffftt111
	         111ttfLLLCCGGG000000GGGCCLLLfft111
	      1111ffLLCG00880000GGGG00008880GCCLLft111
	     11tfLLCG0880GCCLLLLLCCLLLLLCCGG0880CLLft111
	   111fLLC0880CCLLLLLLLLL08CLLLLLLLLLCG880GLLft11
	  11tLLCG88GCLCCLLLLLLLLL08CLLLLLLLLCCLLG0@0CLLf11
	 11tLLC0@0CLLL080CLLLLLLLG8LLLLLLLCG88CLLLG88GLLf11
	 1tLLC8@GLLLLLCG880CLLLLLG8LLLLLCG880CLLLLLC88GLLf1
	11LLL0@GLLLLLLLLLG080CLLLG8LLLCG88GCLLLLLLLLC88CLLt1
	1tLLG@8LLLLLLLLLLLLC080CLG8LCG80GCLLLLLLLLLLLG@0LLL1
	1LLL0@GLLLLLLLLLLLLLLCG8008080GCLLLLLLLLLLLLLC8@CLLt
	1LLL0@CLLG000000000000G0@@@@800000000000000CLL0@CLLt
	1LLL0@GLLLCCCCCCCCCCCCG0088080CCCCCCCCCCCCCLLC8@CLLt
	1tLLG@8LLLLLLLLLLLLCG08GLG8LC080CLLLLLLLLLLLLG@0LLL1
	11LLC8@GLLLLLLLLLCG80GLLLG8LLLC080CLLLLLLLLLC8@GLLt1
	 1tLLC8@GLLLLLLC080GLLLLLG8LLLLLC080GCLLLLLC8@GLLf1
	 11tLLC8@0CLLLG80GLLLLLLLG8LLLLLLLC080CLLLG88GLLf11
	  11tLLCG88GLLCCLLLLLLLLL08CLLLLLLLLCCLLC088CLLf11
	   111fLLC0880CLLLLLLLLLL08CLLLLLLLLLCG080GLLft111
	    111tfLLCG8880GCCLLLLLCCLLLLLLCCG0880CLLft1111
	      1111ffLLCG0088000GGGGGGG0008800CCLLft111
	         1111tffLLCCGG00000000GGGCLLLftt1111
	            11111ttfffLLLLLLLLLffftt11111
	                1ttfffLLLLLLLLLffft
		`)
	fmt.Println("Aether Runtime Environment. Version: dev.v0.0.1")
}

func Startup() {
	globals.SetGlobals()
	persistence.CreateDatabase()
	ShowIntro()
	ReadFlags()
	StartSchedules()
	logging.Log(1, "Startup complete.")
	// TEST Insert the localhost data.
	var addrLocal api.Address
	addrLocal.Location = "127.0.0.1"
	addrLocal.Sublocation = ""
	addrLocal.LocationType = 4
	addrLocal.Port = 8001
	addrLocal.LastOnline = 1111111
	addrLocal.Protocol.VersionMajor = 1
	addrLocal.Protocol.VersionMinor = 1
	addrLocal.Protocol.Subprotocols = []api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}}
	addrLocal.Client.VersionMajor = 1
	addrLocal.Client.VersionMinor = 1
	addrLocal.Client.VersionPatch = 1
	addrLocal.Client.ClientName = "Aether"
	persistence.BatchInsert([]interface{}{addrLocal})
	// dispatch.Sync(addrLocal)
}

func Shutdown() {
	logging.Log(1, "Shutdown initiated.")
	fmt.Println("Shutdown initiated.")
	globals.StopLiveDispatcherCycle <- true // Send true through the channel to stop the dispatch.
	globals.StopStaticDispatcherCycle <- true
	globals.StopAddressScannerCycle <- true
	globals.StopUPNPCycle <- true
	// mature, err := persistence.LocalNodeIsMature()
	// if err != nil {
	// 	logging.LogCrash(err)
	// }
	// if mature {
	// 	globals.StopMatureCacheGenerationCycle <- true
	// } else {
	// 	globals.StopImmatureCacheGenerationCycle <- true
	// }
	logging.Log(1, "Shutdown is complete.")
	fmt.Println("Shutdown is complete. Bye.")
	os.Exit(0)
}

func main() {
	Startup()
	server.Serve()
	Shutdown()
}
