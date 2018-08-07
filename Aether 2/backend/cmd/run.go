package cmd

import (
	"aether-core/backend/beapiserver"
	"aether-core/backend/dispatch"
	"aether-core/backend/feapiconsumer"
	"aether-core/backend/responsegenerator"
	"aether-core/backend/server"
	// "aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/ports"
	"aether-core/services/scheduling"
	"aether-core/services/upnp"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	var loggingLevel int
	var backendAPIPort int
	var backendAPIPublic bool
	var adminFeAddr string
	var adminFePk string
	cmdRun.Flags().IntVarP(&loggingLevel, "logginglevel", "", 0, "Global logging level of the app.")
	cmdRun.Flags().IntVarP(&backendAPIPort, "backendapiport", "", 0, "Sets the port that the backend will attempt to serve the backend API output from. If this port is occupied, it will pick another, therefore it's not safe to assume that this will be the actual backend API port.")
	cmdRun.Flags().BoolVarP(&backendAPIPublic, "backendapipublic", "", false, "If you set this to true, your node will expose the backend api port to the public internet, as well. If not, it will be only served locally. Defaults to false. The reason you might want this is to put the backend on a VPS and make your frontend connect to it, so that it can stay online 24/7.")
	cmdRun.Flags().StringVarP(&adminFeAddr, "adminfeaddr", "", "127.0.0.1:45001", "Spawner FE Address is the address of the frontend that spawns this backend instance. The backend will reach out to this address to tell that it is ready at which port.")
	cmdRun.Flags().StringVarP(&adminFePk, "adminfepk", "", "", "Spawner FE Public Key is the public key of the frontend instance that is spawning the backend process. This is useful to give, because if admin needs to change (ex: when you want to monitor the status of the backend from a different machine than you've installed) you can move your FE config to the new machine, run the FE from the new machine and it will update the admin FE address because it can authenticate with the key.")
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
		gotValidPort := make(chan bool)
		go beapiserver.StartBackendServer(gotValidPort)
		<-gotValidPort // Only proceed after this is true.
		feapiconsumer.SendBackendReady()
		server.StartMimServer()
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
	ports.VerifyBackendPorts()
	// UPNP tries to port map every 10 minutes.
	globals.BackendTransientConfig.StopUPNPCycle = scheduling.ScheduleRepeat(func() { upnp.MapPort() }, 10*time.Minute, time.Duration(0), nil)
	dispatch.Bootstrap() // This will run only if needed.
	// The dispatcher that seeks live nodes runs every minute.
	globals.BackendTransientConfig.StopNeighbourhoodCycle = scheduling.ScheduleRepeat(func() { dispatch.NeighbourWatch() }, 1*time.Minute, time.Duration(0), nil)
	globals.BackendTransientConfig.StopExplorerCycle = scheduling.ScheduleRepeat(func() { dispatch.Explore() }, 10*time.Minute, time.Duration(10)*time.Minute, nil)
	globals.BackendTransientConfig.StopInboundConnectionCycle = scheduling.ScheduleRepeat(func() { dispatch.InboundConnectionWatch() }, 1*time.Minute, time.Duration(5)*time.Minute, nil)

	// Address scanner goes through all prior unconnected addresses and attempts to connect to them to establish a relationship. It starts 30 minutes after a node is started, so that the node will actually have a chance to collect some addresses to check.
	globals.BackendTransientConfig.StopAddressScannerCycle = scheduling.ScheduleRepeat(func() { dispatch.AddressScanner() }, 2*time.Hour, time.Duration(15)*time.Minute, nil)
	// Attempt cache generation every hour, but it will be pre-empted if the last cache generation is less than 23 hours old, and if the node is not tracking the head. So that this will run effectively every day, only.
	globals.BackendTransientConfig.StopCacheGenerationCycle = scheduling.ScheduleRepeat(func() { responsegenerator.GenerateCaches() }, 1*time.Hour, time.Duration(3)*time.Minute, nil)
}

func shutdown() {
	logging.Log(1, "Shutdown initiated. Entering lameduck mode and stopping all scheduled tasks and routines. This will take a minute.")
	// Initiate lameduck mode. This will start declining and inbound and outbound requests, as well as reverse connections requests. Ongoing database actions can still be processed.
	globals.BackendTransientConfig.LameduckInitiated = true
	logging.Log(1, "Waiting 60 seconds to let all network i/o close gracefully...")
	time.Sleep(time.Duration(60) * time.Second)
	// Initiate shutdown. At this point, if anything is still being written into the database, they will attempt to exit gracefully.
	globals.BackendTransientConfig.ShutdownInitiated = true
	// Stop routines
	globals.BackendTransientConfig.StopNeighbourhoodCycle <- true
	globals.BackendTransientConfig.StopInboundConnectionCycle <- true
	globals.BackendTransientConfig.StopExplorerCycle <- true
	globals.BackendTransientConfig.StopAddressScannerCycle <- true
	globals.BackendTransientConfig.StopUPNPCycle <- true
	globals.BackendTransientConfig.StopCacheGenerationCycle <- true
	// logging.Logf(1, "Inbounds: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.Inbounds))
	// logging.Logf(1, "Outbounds: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.Outbounds))
	// logging.Logf(1, "InboundHistory: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.InboundHistory))
	// logging.Logf(1, "OutboundHistory: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.OutboundHistory))
	// logging.Logf(1, "Last Inbound Sync Timestamp: %v", globals.BackendTransientConfig.Bouncer.GetLastInboundSyncTimestamp(false))
	// logging.Logf(1, "Last Successful Outbound Sync Timestamp: %v", globals.BackendTransientConfig.Bouncer.GetLastOutboundSyncTimestamp(true))
	// logging.Logf(1, "Inbounds in the last 5 minutes: %v", len(globals.BackendTransientConfig.Bouncer.GetInboundsInLastXMinutes(5)))
	// logging.Logf(1, "Successful outbounds in the last 5 minutes: %v", len(globals.BackendTransientConfig.Bouncer.GetOutboundsInLastXMinutes(5, true)))
	// logging.Logf(2, "Nonces: %v", globals.BackendTransientConfig.Nonces)

	logging.Log(1, "Waiting 5 seconds to let DB close gracefully...")
	time.Sleep(time.Duration(5) * time.Second) // Wait 5 seconds to let DB tasks complete.
	// And after that, we shut down the database.
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
