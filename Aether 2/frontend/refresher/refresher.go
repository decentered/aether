// Frontend > Refresher
// This package contains the refresher loop that keeps frontend in sync with the backend it is connected to.

package refresher

import (
	"aether-core/frontend/beapiconsumer"
	"aether-core/frontend/clapiconsumer"
	"aether-core/frontend/festructs"
	"aether-core/io/api"
	pbstructs "aether-core/protos/mimapi"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/toolbox"
	// "github.com/davecgh/go-spew/spew"
	// "fmt"
	"encoding/json"
	"github.com/asdine/storm/q"
	"strings"
	"sync"
	"time"
)

var (
	GlobalStatistics festructs.GlobalStatisticsCarrier
	// This will be updated by global statistics carrier when it finishes, and we'll be using this as a basis for election calculations on the global scope. Local scopes (boards' elections) have their own population counts.
	RefreshRanBeforeOnThisRun bool
	LastRefreshDuration       time.Duration
)

func Refresh() {
	globals.FrontendTransientConfig.RefresherMutex.Lock()
	defer globals.FrontendTransientConfig.RefresherMutex.Unlock()

	/*----------  Set status visible in the client  ----------*/
	clapiconsumer.FrontendAmbientStatus.RefresherStatus = "Compiling..."
	clapiconsumer.SendFrontendAmbientStatus()

	preRefresh()

	timeStart := time.Now()
	if !RefreshRanBeforeOnThisRun {
		// logging.Logf(1, "This is the first refresh of this run. Initialising KvStore buckets.")
		festructs.InitialiseKvStore()
		RefreshRanBeforeOnThisRun = true
	}
	// Create new global statistics container at every refresh cycle.
	PrepNewGlobalStatistics()
	// GlobalStatistics.LastReferenced = 0 // todo debug
	// Prefill cache for this refresh and set its end to the global end.
	nowts := beapiconsumer.PrefillCache(GlobalStatistics.LastReferenced) // Old refresh end (lastref) is given as start
	defer beapiconsumer.ReleaseCache()
	// Save it to the frontend transient config so everyone can access it, not just refresher
	globals.FrontendTransientConfig.RefresherCacheNowTimestamp = nowts
	// RefreshGlobalStatistics refreshes basic things like total number of users in the last 6 months (total population), which is something we need when we're calculating global user headers, because signals in those global user headers deal with elections, and an election needs to know the total population to be able to determine whether it is valid (i.e. enough % of people voted) or not.
	newUserEntities := GlobalStatistics.Refresh(nowts)
	// Get the local user entity if present, and add it to new user entities, so that it will always be refreshed.

	alu := globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()
	if len(alu) != 0 {
		var key api.Key
		json.Unmarshal([]byte(alu), &key)
		kp := key.Protobuf()
		newUserEntities = append(newUserEntities, &kp)
	}
	// Refresh all users
	RefreshGlobalUserHeaders(newUserEntities, nowts)
	// Get extant ambient boards
	ambientBoards := festructs.GetCurrentAmbients()
	RefreshBoards(nowts, ambientBoards)
	ambientBoards.Save() // Save the updated ambients (update happens inside refresh boards)
	GenerateHomeView()
	GeneratePopularView()
	// at the end, delete too old lastrefresheds from the whole kvstore
	DeleteStaleData(nowts)
	// Finally, run the routines that we want after the refresh, mainly, letting the client know a refresh has happened, updating the ambients it has, and so on.
	LastRefreshDuration = time.Since(timeStart)
	postRefresh()
}

func PrepNewGlobalStatistics() {
	GlobalStatistics = festructs.GlobalStatisticsCarrier{}
	err := globals.KvInstance.One("Id", 1, &GlobalStatistics)
	if err != nil && strings.Contains(err.Error(), "not found") {
		GlobalStatistics = festructs.NewGlobalStatisticsCarrier()
	} else if err != nil {
		logging.LogCrashf("Prepare new global statistics in frontend refresh cycle has failed with the error: %v", err)
	}
}

// DeleteStaleData deletes the data that we've ceased updating. This does not mean the data is deleted from the backend store, it just means that the cache copy we keep on the frontend is. So if the user wants to see the same thing again, the click will cause a cache miss, it will be pulled and compiled from the backend again (if it's still extant there) and served to the user.
func DeleteStaleData(nowts int64) {
	logging.Logf(1, "Starting deletion of stale data.")
	cutoff := toolbox.CnvToCutoffDays(globals.FrontendConfig.GetKvStoreRetentionDays())
	// Delete stale boards
	query := globals.KvInstance.Select(q.Lte("LastRefreshed", cutoff))
	err := query.Delete(new(festructs.BoardCarrier))
	if err != nil {
		logging.Logf(1, "Deletion of stale boards errored out. Err: %v", err)
	}
	err2 := query.Delete(new(festructs.ThreadCarrier))
	if err2 != nil {
		logging.Logf(1, "Deletion of stale threads errored out. Err: %v", err2)
	}
	err3 := query.Delete(new(festructs.UserHeaderCarrier))
	if err3 != nil {
		logging.Logf(1, "Deletion of stale user headers errored out. Err: %v", err3)
	}
	logging.Logf(1, "Stale data deletion is complete.")
}

func RefreshGlobalUserHeaders(newUserEntities []*pbstructs.Key, nowts int64) {
	var uhcs []festructs.UserHeaderCarrier
	err := globals.KvInstance.All(&uhcs)
	if err != nil {
		logging.Logf(1, "Fetching all global user headers before the refresh has failed. Error: %v", err)
	}
	uhcBatch := festructs.UHCBatch(uhcs)
	for k, _ := range newUserEntities {
		if i := uhcBatch.Find(newUserEntities[k].GetProvable().GetFingerprint()); i != -1 {
			uhcBatch[i].Users.InsertFromProtobuf([]*pbstructs.Key{newUserEntities[k]}, nowts)
		} else {
			uhc := festructs.NewUserHeaderCarrier(newUserEntities[k].GetProvable().GetFingerprint(), "", nowts)
			uhc.Users.InsertFromProtobuf([]*pbstructs.Key{newUserEntities[k]}, nowts)
			uhcBatch = append(uhcBatch, uhc)
		}
	}
	for k, _ := range uhcBatch {
		uhcBatch[k].Refresh([]string{}, GlobalStatistics.UserCount, nowts)
		// ^ We have no default mods in global, and totalPop comes from global statistics.
		/*
			TODO FUTURE
			This is where you calculate and insert the global mods assigned by the CA.
		*/
	}
	// We need to add items coming in from this delta.

	// logging.Logf(1, "This is the refreshed global user headers. %s", spew.Sdump(uhcBatch))
}

func RefreshBoards(nowts int64, extantABs *festructs.AmbientBoardBatch) {
	newBoardEntities := beapiconsumer.GetBoards(GlobalStatistics.LastReferenced, nowts, []string{}, false, false)
	GlobalStatistics.LastReferenced = nowts
	GlobalStatistics.Save()
	var bcs []festructs.BoardCarrier
	err := globals.KvInstance.All(&bcs)
	if err != nil {
		logging.Logf(1, "Fetching all boards in the refresh has failed. Error: %v", err)
	}
	bcBatch := festructs.BCBatch(bcs)
	bcBatch.Insert(newBoardEntities)
	wg := sync.WaitGroup{}
	for k, _ := range bcBatch {
		wg.Add(1)
		go RefreshBoard(bcBatch[k], &wg, extantABs)
	}
	wg.Wait()
}

// RefreshBoard does a few things. First of all, it updates the board statistics, then it updates the board's own user headers, then it updates the board's own entity, then it updates the board's thread entities, then it starts the process to refresh tracked threads and gives the newly updated thread entities to those threads, so that they don't have to compile those twice.
func RefreshBoard(
	bc festructs.BoardCarrier,
	wg *sync.WaitGroup,
	extantABs *festructs.AmbientBoardBatch,
) {
	bc.RefreshWithoutSave(globals.FrontendTransientConfig.RefresherCacheNowTimestamp)
	refreshedAmbients := bc.ConstructAmbientBoards()
	extantABs.UpdateBatch(refreshedAmbients)
	RefreshThreads(&bc)
	// UpdateBoardThreadsCount(&bc)
	wg.Done()
	bc.Threads.SortByScore()
	bc.Save()
}

func RefreshThreads(bc *festructs.BoardCarrier) {
	// Determine what stuff we need to refresh
	newThreadEntities := beapiconsumer.GetThreads(bc.LastReferenced, globals.FrontendTransientConfig.RefresherCacheNowTimestamp, []string{}, bc.Fingerprint, false, false)
	bc.Threads.InsertFromProtobuf(newThreadEntities)
	wg := sync.WaitGroup{}
	for k, _ := range bc.Threads {
		wg.Add(1)
		go RefreshThread(bc.Threads[k], bc, &wg)
	}
	wg.Wait()
	// Thread refresh is done. Sort them based on score. This is the order we save them in.
}

// RefreshThread refreshes a thread. The way it does is that it first looks at whether we have an extant thread carrier for that thread. If we do, it triggers a refresh on it. If not, it creates one, fills it with the required data, and then it triggers a refresh on it.
func RefreshThread(cthread festructs.CompiledThread, bc *festructs.BoardCarrier, wg *sync.WaitGroup) {
	// Get thread carrier, create one if not present.
	tc := festructs.ThreadCarrier{}
	err := globals.KvInstance.One("Fingerprint", cthread.Fingerprint, &tc)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tc = festructs.NewThreadCarrier(cthread.Fingerprint, cthread.Board, globals.FrontendTransientConfig.RefresherCacheNowTimestamp)
		}
	}
	tc.Refresh(bc.Boards.GetBoardSpecificUserHeaders(), bc, globals.FrontendTransientConfig.RefresherCacheNowTimestamp)
	wg.Done()
}

func preRefresh() {
	// If time since the last SFW list pull is longer than an hour, refresh it first.
	lastSFWListUpdate := globals.FrontendConfig.ContentRelations.SFWList.LastUpdate
	if time.Since(time.Unix(lastSFWListUpdate, 0)).Minutes() > 60 {
		// logging.Logf(1, "It's been longer than an hour after the last SFW list refresh. Refreshing the list first.")
		globals.FrontendConfig.ContentRelations.SFWList.Refresh()
		return
	}
	// logging.Logf(1, "The SFW list was refreshed recently. Skipping the refresh.")
}

func postRefresh() {
	// ^^ As in after refresh, not refreshing posts.
	clapiconsumer.DeliverAmbients()
	clapiconsumer.PushLocalUserAmbient()
	clapiconsumer.FrontendAmbientStatus.RefresherStatus = "Idle"
	clapiconsumer.FrontendAmbientStatus.LastRefreshTimestamp = time.Now().Unix()
	clapiconsumer.FrontendAmbientStatus.LastRefreshDurationSeconds = int32(LastRefreshDuration.Seconds())
	clapiconsumer.SendFrontendAmbientStatus()
	clapiconsumer.SendHomeView()
	clapiconsumer.SendPopularView()
	clapiconsumer.SendNotifications()
	festructs.NotificationsSingleton.Save()
}
