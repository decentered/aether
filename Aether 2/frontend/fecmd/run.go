package fecmd

import (
	// "aether-core/frontend/beapiconsumer"
	"aether-core/frontend/besupervisor"
	// "aether-core/frontend/clapiconsumer"
	"aether-core/frontend/feapiserver"
	// "aether-core/protos/clapi"
	// "aether-core/frontend/festructs"
	// "aether-core/frontend/inflights"
	"aether-core/frontend/kvstore"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/ports"
	"aether-core/services/scheduling"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"aether-core/frontend/refresher"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	var loggingLevel int
	var clientIp string
	var clientPort int
	cmdRun.Flags().IntVarP(&loggingLevel, "logginglevel", "", 0, "Sets the frontend logging level.")
	cmdRun.Flags().StringVarP(&clientIp, "clientip", "", "127.0.0.1", "This is the IP of the client that is starting the frontend instance. THis is almost always 127.0.0.1 since clients and frontends almost always live in the same computer.")
	cmdRun.Flags().IntVarP(&clientPort, "clientport", "", 0, "The port of the client instance starting the frontend. Frontend will call back at this endpoint via GRPC and confirm it's ready.")
	cmdRoot.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Start an Aether frontend instance that maintains a compiled data source for the frontend to use..",
	Long: `This starts a Aether frontend process. This is the main process that communicates and responds to client requests. This is where all the different pieces of data from the network actually comes together and ends up as boards, threads, posts, users and so on.
`,
	Run: func(cmd *cobra.Command, args []string) {
		EstablishConfigs(cmd)
		// start frontend server
		gotValidPort := make(chan bool)
		go feapiserver.StartFrontendServer(gotValidPort)
		<-gotValidPort // Only proceed after this is true.
		go besupervisor.StartLocalBackend()
		for globals.FrontendTransientConfig.BackendReady != true {
			// Block until the backend tells the frontend via gRPC that it is ready.
		}
		// debug
		// go testBackend()
		// end debug
		kvstore.OpenKVStore()
		defer kvstore.CloseKVStore()
		kvstore.CheckKVStoreReady()
		startSchedules()
		// feapiserver.SendAmbients(false)

		select {} // todo handle this. having the local backend on its own goroutine helps with the excessive cpu use we're seeing, but this select {} should become a more idiomatic channel block. and we should determine what we can do in the case the backend crashes.
		// err := besupervisor.StartLocalBackend()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println("does it ever come here?")
		// select {}
		/*
			^ Why is that needed? Because we want to start the frontend serve before the backend runtime, so that the BE can reach out to frontend server to send an ACK. But if we put the FE server startup in a goroutine and start the BE server, then we have no place to actually put in the fe>cl ack, which needs to be sent AFTER the BE is ready. Basiaclly, hmm.
		*/
	},
}

func startSchedules() {
	logging.Log(1, "Setting up cyclical frontend tasks is starting.")
	defer logging.Log(1, "Setting up cyclical frontend tasks is complete.")
	ports.VerifyFrontendPorts()
	globals.FrontendTransientConfig.StopRefresherCycle = scheduling.ScheduleRepeat(func() {
		start := time.Now()
		refresher.Refresh()
		elapsed := time.Since(start)
		logging.Logcf(1, "We've refreshed the frontend. It took: %s", elapsed)

	}, 2*time.Minute, time.Duration(0), nil)
}

// func testBackend() {
// 	logging.Logf(1, "This is the backend test. Will sleep for 2 seconds and attempt to request data from the BE.")
// 	time.Sleep(3 * time.Second)
// 	logging.Logf(1, "Attempting to request data from the backend.")
// 	// start := time.Now()
// 	// keys := beapiconsumer.GetKeys(0, 0, []string{})
// 	// elapsed := time.Since(start)
// 	// logging.Logcf(1, "We've got %v keys from the backend. It took: %s", len(keys), elapsed)

// 	// start2 := time.Now()
// 	// votes := beapiconsumer.GetVotes(0, 0, []string{})
// 	// elapsed2 := time.Since(start2)
// 	// logging.Logcf(1, "We've got %v votes from the backend. It took: %s", len(votes), elapsed2)

// 	// start3 := time.Now()
// 	// boards := beapiconsumer.GetBoards(0, 0, []string{})
// 	// elapsed3 := time.Since(start3)
// 	// logging.Logcf(1, "We've got %v boards from the backend. It took: %s", len(boards), elapsed3)

// 	// start4 := time.Now()
// 	// namemaps := festructs.GetNameMaps(0)
// 	// elapsed4 := time.Since(start4)
// 	// logging.Logcf(1, "We've got %v namemaps from the backend. It took: %s", len(namemaps), elapsed4)
// 	// refresher.Refresh()
// 	start5 := time.Now()
// 	refresher.Refresh()
// 	elapsed5 := time.Since(start5)
// 	logging.Logcf(1, "We've refreshed the frontend. It took: %s", elapsed5)
// 	logging.Logf(1, "Finished attempting to request data from the backend.")
// }
