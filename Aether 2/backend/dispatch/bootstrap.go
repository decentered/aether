// Backend > Routines > Bootstrap
// This file defines the bootstrap routine which runs at warmup conditions. This routine connects to multiple remotes and downloads from them to reduce the network load on any single one.

package dispatch

import (
	"aether-core/io/api"
	"aether-core/services/configstore"
	"aether-core/services/globals"
	"aether-core/services/logging"
	// "github.com/davecgh/go-spew/spew"
	"time"
)

const (
	bsLoc    = configstore.DefaultBootstrapperLocation
	bsSubloc = configstore.DefaultBootstrapperSublocation
	bsPort   = configstore.DefaultBootstrapperPort
)

func getBootstrappers() []api.Address {
	// defaultBootstrapper := constructBootstrapper(bsLoc, bsSubloc, bsPort)
	resp, err := api.GetPageRaw(bsLoc, bsSubloc, bsPort, "bootstrappers", "GET", []byte{})
	if err != nil {
		logging.Logf(1, "Getting bootstrappers failed from this address. Error: %v, Address: %v/%v:%v", err, bsLoc, bsSubloc, bsPort)
	}
	// spew.Dump(resp)
	bsers := []api.Address{}
	if resp.Address.Type == 254 || resp.Address.Type == 3 {
		resp.Address.Location = bsLoc
		resp.Address.Sublocation = bsSubloc
		resp.Address.Port = bsPort
		bsers = append(bsers, resp.Address) // The first bootstrapper is the address we connected to if it's a bootstrapper itself (type=3)
	}
	bsers = append(bsers, resp.ResponseBody.Addresses...)
	slen := len(bsers)
	if slen > 99 {
		slen = 99
	}
	return bsers[0:slen] // limit it to protect from DDoSs
}

type execPlan struct {
	addr      api.Address
	endpoints []string
}

// constructExecPlan creates an execution plan based on the addresses given. It splits all the endpoins we have to call across.
func constructExecPlans(bootstrappers []api.Address) []execPlan {
	if len(bootstrappers) == 0 {
		return []execPlan{}
	}
	execplans := []execPlan{}
	for key, _ := range bootstrappers {
		execplans = append(execplans, execPlan{addr: bootstrappers[key]})
	}
	servingSubprots := globals.BackendConfig.GetServingSubprotocols()
	entities := []string{}
	for _, subprot := range servingSubprots {
		entities = append(entities, subprot.SupportedEntities...)
	}
	for key, _ := range entities {
		mod := key % len(execplans) // 0 % 7 = 0, 3 % 7 = 3, 7 % 7 = 0 (loops over)
		execplans[mod].endpoints = append(execplans[mod].endpoints, entities[key])
	}
	// spew.Dump(execplans)
	return execplans
}

// Bootstrap is the 'catch-up' logic that runs whenever a node falls too far behind the network head for any reason. One of the main uses is the start from the first boot, but it can also be that the node has been offline for more than bootstrap hit interval.
func doBootstrap() {
	/*
			  Logic:
			  1) Hit /bootstrappers on GET. Collect all bootstrappers.
			  2) Do a check over them to find out which ones are online.
			  3) Create a mapper that splits the endpoints apart and creates an execution plan
			  4) Run the execution plan as a series of syncs.
		    5) Run all bootstrap nodes as normal syncs to collect manifests and find missing data.
	*/
	bootstrappers := getBootstrappers()
	onlineBootstrappers := Pinger(bootstrappers)
	if len(onlineBootstrappers) == 0 {
		logging.Logf(1, "No online bootstrappers were found. Exiting bootstrap.")
		return
	}
	execPlans := constructExecPlans(onlineBootstrappers)
	// Do a partial sync for everything in the exec plan.
	errs := []error{}
	// Go through each remote in the exec plans and call them based on the types we want to pull from it.
	for key, _ := range execPlans {
		err := Sync(execPlans[key].addr, execPlans[key].endpoints)
		if err != nil {
			errs = append(errs, err)
		}
	}
	// Go through each remote in the exec plan and call all endpoints in them. This should cause a manifest scan and not much download, and a timestamp setting. This is insurance to make sure that the data we have is the union of all bootstrappers we connected to.
	for key, _ := range execPlans {
		err := Sync(execPlans[key].addr, []string{})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		logging.Logf(1, "These were the errors encountered in the bootstrap syncs: Errors: %#v", errs)
	}
	globals.BackendConfig.SetLastBootstrapAddressConnectionTimestamp(time.Now().Unix())
}

func Bootstrap() {
	bsOfflineMinutes := globals.BackendConfig.GetBootstrapAfterOfflineMinutes()
	lastBs := globals.BackendConfig.GetLastBootstrapAddressConnectionTimestamp()
	cutoff := int64(time.Now().Add(-(time.Duration(bsOfflineMinutes) * time.Minute)).Unix())
	if cutoff > lastBs {
		logging.Logf(1, "Bootstrap decided it needs to run because it's been longer than allowed cutoff since the last time it was run.")
		doBootstrap()
	}
}
