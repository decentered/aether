package cmd

import (
	"aether-core/backend/dispatch"
	// "aether-core/backend/metrics"
	"aether-core/backend/server"
	"aether-core/io/api"
	"aether-core/io/persistence"
	"aether-core/services/create"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/scheduling"
	// "aether-core/services/ports"
	"encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"io/ioutil"
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
	var metricsDebugMode bool
	var swarmPlan string
	var killTimeout int
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
	cmdOrchestrate.Flags().BoolVarP(&metricsDebugMode, "metricsdebugmode", "d", false, "Enable sending of debug-mode metrics. These metrics are designed to provide a debuggable view of what the node is doing.")
	cmdOrchestrate.Flags().StringVarP(&swarmPlan, "swarmplan", "s", "", "This flag allows you to load a swarm plan to your swarm nodes. This swarm plan does have a list of TO-FROM node connections with certain delays, so that you can schedule connections to happen in a fashion that is pre-mediated. This allows you to kickstart a few node connections and see how network behaves based on new data, for example.")
	cmdOrchestrate.Flags().IntVarP(&killTimeout, "killtimeout", "k", 120, "If given, this sets up a maximum lifetime in seconds for the node. This is useful in swarm tests in which all swarm nodes have to exit so that the test can move on to the data analysis stage. ")
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
			logging.LogCrash(fmt.Sprintf("Please provide either all of the parts of a bootstrap node (IP:Port/Type) or none. You've provided:%s:%d, Type: %d",
				flags.bootstrapIp.value.(string),
				flags.bootstrapPort.value.(int),
				flags.bootstrapType.value.(int)))
		}
		addr := constructCallAddress(api.Location(flags.bootstrapIp.value.(string)), uint16(flags.bootstrapPort.value.(int)), uint8(flags.bootstrapType.value.(int)))
		// Prep the database
		persistence.CreateDatabase()
		addrs := []api.Address{addr}
		persistence.InsertOrUpdateAddresses(&addrs)
		showIntro()
		// startup()
		globals.DispatcherExclusions = make(map[*interface{}]time.Time)
		if flags.syncAndQuit.value.(bool) {
			// First, verify external port, so that our metrics will report the right port.
			// ports.VerifyExternalPort()
			// We just want to connect, pull and quit.
			dispatch.Sync(addr)
			logging.Log(1, fmt.Sprintf("We've gotten everything in the node %#v.", addr))
			fmt.Printf("This is the external port for this node: %d\n", globals.BackendConfig.GetExternalPort())
			// if flags.swarmPlan.changed {
			// 	scheduleSwarmPlan(flags.swarmPlan.value.(string))
			// }
			// client, conn := metrics.StartConnection()
			// defer conn.Close()
			// r := metrics.DeliverBackendMetrics(client)
			// fmt.Println("deliver metrics result:")
			// fmt.Println(r)
			// freeport := ports.GetFreePort()
			// fmt.Printf("We asked OS to give us a free port. This was the result: %s\n", freeport)
		} else {
			startSchedules()
			if flags.swarmPlan.changed {
				scheduleSwarmPlan(flags.swarmPlan.value.(string))
			}
			if flags.killTimeout.changed {
				logging.Log(1, fmt.Sprintf("This node set to shut down in %d seconds.", flags.killTimeout.value.(int)))
				scheduling.ScheduleOnce(func() {
					shutdown()
				}, time.Duration(flags.killTimeout.value.(int))*time.Second)
			}
			server.Serve()
		}
	},
}

// Orchestrate endpoint will allow us to first generate random data, then pull that into the local database.

func constructCallAddress(ip api.Location, port uint16, addrtype uint8) api.Address {
	subprots := []api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}}
	addr, err := create.CreateAddress(ip, "", 4, port, addrtype, 0, 1, 0, subprots, 2, 0, 0, "Aether")
	if err != nil {
		logging.LogCrash(err)
	}
	return addr
}

func parsePlansForThisNode(planAsByte []byte) []map[string]interface{} {
	// This is just JSON parsing without a backing struct. The other alternative was copying over the struct (I don't want to have a swarmtest dependency here, I can't import from there), so this is arguably cleaner.
	var f interface{}
	err2 := json.Unmarshal(planAsByte, &f)
	if err2 != nil {
		logging.LogCrash(fmt.Sprintf("The swarm plan JSON parsing failed. Error: %s", err2))
	}
	sch := f.([]interface{})
	var plansForThisNode []map[string]interface{}
	for _, val := range sch {
		valMapped := val.(map[string]interface{})
		if valMapped["FromNodeAppName"] == globals.BackendTransientConfig.AppIdentifier {
			plansForThisNode = append(plansForThisNode, valMapped)
		}
	}
	return plansForThisNode
}

// scheduleSwarmPlan finds and reads through the swarm plan json file, and determines which schedules apply to it. For those which apply, it inserts into the scheduler logic.
func scheduleSwarmPlan(planloc string) {
	// Read the file
	planAsByte, err := ioutil.ReadFile(planloc)
	if err != nil {
		logging.LogCrash(fmt.Sprintf("The swarm plan document could not be read. Error: %s", err))
	}
	// Parse the plans relevant to this specific node based on the AppIdentifier
	plans := parsePlansForThisNode(planAsByte)
	// Insert the plans into the scheduler.
	for _, val := range plans {
		ip := val["ToIp"].(string)
		port := uint16(val["ToPort"].(float64))
		nodetype := uint8(val["ToType"].(float64))
		triggerafter := time.Duration(int64(val["TriggerAfter"].(float64)))
		logging.Log(1, fmt.Sprintf("This node is going to attempt to connect to the address: %s:%d in %v", ip, port, triggerafter))
		scheduling.ScheduleOnce(func() {
			logging.Log(1, fmt.Sprintf("Injecting the address: %s:%d into the database now.", ip, port))
			addr := constructCallAddress(api.Location(ip), port, nodetype)
			addrs := []api.Address{addr}
			persistence.InsertOrUpdateAddresses(&addrs)
		}, triggerafter)
	}
	// spew.Dump(plans)
	// spew.Dump(schMapped)
	// spew.Dump(schMapped[2]["ToPort"])

	// fmt.Println(schedules[0]["TriggerAfter"])
	// for _, val := range itemsMap {

	// }
}
