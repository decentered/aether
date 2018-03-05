package cmd

import (
	"aether-core/backend/analytics"
	"aether-core/backend/dispatch"
	"aether-core/backend/server"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/create"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	var orgName string
	var appName string
	var loggingLevel int
	var port int
	var externalIp string
	var bootstrapIp string
	var bootstrapPort int
	var bootstrapType int
	var syncAndQuit bool
	var printToStdout bool
	cmdOrchestrate.Flags().StringVarP(&orgName, "orgname", "o", "Air Labs", "Global transient org name for the app.")
	cmdOrchestrate.Flags().StringVarP(&appName, "appname", "a", "Aether", "Global transient app name for the app.")
	cmdOrchestrate.Flags().IntVarP(&loggingLevel, "logginglevel", "l", 0, "Global logging level of the app.")
	cmdOrchestrate.Flags().IntVarP(&port, "port", "p", 49999, "The port the external world can use to communicate with this node.")
	cmdOrchestrate.Flags().StringVarP(&externalIp, "externalip", "e", "0.0.0.0", "The IP address the external world can use to communicate with this node.")
	cmdOrchestrate.Flags().StringVarP(&bootstrapIp, "bootstrapip", "x", "127.0.0.1", "The bootstrap node's ip address that will be inserted into the local database at the start. If you provide this, you also have to provide port and type.")
	cmdOrchestrate.Flags().IntVarP(&bootstrapPort, "bootstrapport", "y", 51000, "The bootstrap node's port that will be inserted into the local database at the start. If you provide this, you also have to provide ip address and type.")
	cmdOrchestrate.Flags().IntVarP(&bootstrapType, "bootstraptype", "z", 255, "The bootstrap node's type (static or live, 2 or 255) that will be inserted into the local database at the start. If you provide this, you also have to provide ip address and port.")
	cmdOrchestrate.Flags().BoolVarP(&syncAndQuit, "syncandquit", "q", false, "The only thing that happens is that we connect to the bootstrap node, ingest everything, and then quit. No ongoing processing.")
	cmdOrchestrate.Flags().BoolVarP(&printToStdout, "printtostdout", "c", false, "Route log output to stdout. This will make the logging library write to stdout instead of log. This is useful in the case of orchestrate if you want to see logs from individual backend nodes.")
	cmdRoot.AddCommand(cmdOrchestrate)
}

var cmdOrchestrate = &cobra.Command{
	Use:   "orchestrate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		flags := establishConfigs(cmd)
		// do things here
		// If all of them are changed, or if none of them are changed, we're good. If SOME of them are changed, we crash.
		if !((flags.bootstrapIp.changed && flags.bootstrapPort.changed && flags.bootstrapType.changed) || (!flags.bootstrapIp.changed && !flags.bootstrapPort.changed && !flags.bootstrapType.changed)) {
			logging.LogCrash("Please provide either all of the parts of a bootstrap node (IP:Port/Type) or none.")
		}
		addr := createBootstrapAddress(api.Location(flags.bootstrapIp.value.(string)), uint16(flags.bootstrapPort.value.(int)), uint8(flags.bootstrapType.value.(int)))
		// Prep the database
		persistence.CreateDatabase()
		err := persistence.InsertOrUpdateAddress(addr)
		if err != nil {
			logging.LogCrash(err)
		}
		showIntro()
		// startup()
		globals.DispatcherExclusions = make(map[*interface{}]time.Time)
		if flags.syncAndQuit.value.(bool) {
			// We just want to connect, pull and quit.
			dispatch.Sync(addr)
			logging.Log(1, fmt.Sprintf("We've gotten everything in the node %#v.", addr))
			client, conn := analytics.Prep()
			resp := analytics.IntroYourself(client)
			fmt.Println(resp)
			conn.Close()
		} else {
			startSchedules()
			server.Serve()
		}
	},
}

// Orchestrate endpoint will allow us to first generate random data, then pull that into the local database.

func createBootstrapAddress(ip api.Location, port uint16, addrtype uint8) api.Address {
	subprots := []api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}}
	addr, err := create.CreateAddress(ip, "", 4, port, addrtype, 0, 1, 0, subprots, 2, 0, 0, "Aether")
	if err != nil {
		logging.LogCrash(err)
	}
	return addr
}

/*
bsip := flags.bootstrapIp.value.(string)
bsport := flags.bootstrapPort.value.(int)
bstype := flags.bootstrapType.value.(int)
createBootstrapAddress(
	bsip.(api.Location),
	bsport.(uint16),
	bstype.(uint8))

*/
