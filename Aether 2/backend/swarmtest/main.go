// Backend > Swarmtest
// This package / file contains the testing framework that creates an e2e swarm testing environment for the backend nodes by spawning multiple nodes and letting them interact with each other.

/*
  This package does a few things in order.
  - Requests the database generator to generate a random database with given parameters. This requires node to be set up and dependencies installed on the local machine.
  - Run a static server with the generated database.
  - Spin up the backend node with the appropriate command, so that the backend node will ingest the data from the static node into its own database. Then let it die.
  - Start the metrics receiver
  - Start the backend node again, this time with full-on status. Inject all non-first nodes with the address of the starter node, and see what happens.

*/

package main

import (
	sas "aether-core/backend/swarmtest/simplemetricsserver"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Utility functions

func createPath(path string) {
	os.MkdirAll(path, 0755)
}

func saveFileToDisk(fileContents []byte, path string, filename string) {
	ioutil.WriteFile(fmt.Sprint(path, "/", filename), fileContents, 0755)
}

// We could do proper flag parsing but it feels unnecessary here, tbh. If this test grows in size we could probably do that.
type settingsStruct struct {
	swarmsize     int // number of nodes to be created and tested against.
	staticnodeloc string
}

type node struct {
	appname           string // appname determines which folder it will be saved within the OS. Changing this is our way of making sure multiple instances of the app can work separately from each other.
	generatedDataPath string // Data path that the static generator will put the resulting static node into.
	staticServerPort  int
}

var nodes []node
var settings settingsStruct

func setDefaults() {
	settings.swarmsize = 10
	settings.staticnodeloc = "temp_generated_data"
}

func generateSwarmNames() {
	for i := 0; i < settings.swarmsize; i++ {
		n := node{}
		n.appname = fmt.Sprintf("Aether-%d", i)
		n.staticServerPort = 17000 + i
		nodes = append(nodes, n)
	}
}

func generateNodeData(n *node) {
	log.Printf("Generating the required random database for %s.", n.appname)
	nodeTempPath := fmt.Sprintf("%s/%s", settings.staticnodeloc, n.appname)
	// Set up a folder with the appropriate name for each node that is being run
	// createPath(nodeTempPath)
	abspath, err := filepath.Abs(nodeTempPath)
	if err != nil {
		panic(err)
	}
	n.generatedDataPath = abspath
	// Check if there is anything that exists in the folder. If so, skip this.
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		cmd := exec.Command("node", "main.js", "--xsmall", "--nosign", abspath)
		cmd.Dir = "../../../Documentation/database-generator/"
		cmd.Run()
		log.Printf("Random database generation for %s is complete.", n.appname)
	} else {
		log.Printf("Skipping data generation because the static node generation folder for %s already exists. If you want it regenerated, please delete the appropriate folder.", n.appname)
	}
}

func startServingStaticNodeAsDataDonor(n *node) {
	fs := http.FileServer(http.Dir(n.generatedDataPath))
	// http.Handle("/", fs)
	svmux := http.NewServeMux()
	svmux.Handle("/", fs)
	go http.ListenAndServe(fmt.Sprint(":", n.staticServerPort), svmux)
}

func insertDataIntoBackendNodeInstance(n *node) {
	// Here, we call our orchestrate command to pull in the data that we need into the db of the given node.
	/*
		The command we need is:
		wipeterm && go run main.go orchestrate --appname="A2Test" --bootstrapip 127.0.0.1 --bootstrapport 9000 --bootstraptype 255 --logginglevel 1 --syncandquit
	*/
	log.Printf("Database insert of randomly generated node data is starting for the node %s", n.appname)
	cmd := exec.Command(
		"go", "run", "main.go", "orchestrate",
		fmt.Sprintf("--appname=%s", n.appname),
		"--bootstrapip=127.0.0.1",
		fmt.Sprintf("--bootstrapport=%d", n.staticServerPort),
		"--bootstraptype=255",
		"--logginglevel=1",
		"--printtostdout",
		"--syncandquit")
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cmd.Stdout = os.Stdout
	cmd.Dir = "../../../aether-core/backend/"
	cmd.Run()
	// fmt.Println(stdout)
	log.Printf("Database insert of randomly generated node data is complete for the node %s", n.appname)
}

func main() {
	setDefaults()
	go sas.StartListening()
	// For each node that we have requested
	generateSwarmNames()
	for _, node := range nodes {
		generateNodeData(&node)
		startServingStaticNodeAsDataDonor(&node)
		insertDataIntoBackendNodeInstance(&node)
	}
	fmt.Printf("%#v\n", nodes)
	// select {}
}
