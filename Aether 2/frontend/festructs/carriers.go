// Frontend > FrontendStructs > Carriers

// This package carries information about

package festructs

import (
	"aether-core/frontend/beapiconsumer"
	pbstructs "aether-core/protos/mimapi"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/rollingbloom"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"time"
)

type Carrier interface {
	// Posts
	GetPosts() CPostBatch
	GetPostsCATDs() CATDBatch
	GetPostsCFGs() CFGBatch
	GetPostsCMAs() CMABatch
	// Threads
	GetThreads() CThreadBatch
	GetThreadsCATDs() CATDBatch
	GetThreadsCFGs() CFGBatch
	GetThreadsCMAs() CMABatch
	// Boards
	GetBoards() CBoardBatch
	GetBoardsCATDs() CATDBatch
	GetBoardsCFGs() CFGBatch
	GetBoardsCMAs() CMABatch
	// Misc
	GetLastRefreshed() int64
	GetLastReferenced() int64
	GetWellFormed() bool
	// Util
	Save()
}

type EntityCarrier struct {
	Fingerprint       string `storm:"id"`
	ParentFingerprint string
	// Board entity data
	Boards      CBoardBatch
	BoardsCATDs CATDBatch
	BoardsCFGs  CFGBatch
	BoardsCMAs  CMABatch
	Statistics  StatisticsCarrier
	// Thread entity data
	Threads      CThreadBatch
	ThreadsCATDs CATDBatch
	ThreadsCFGs  CFGBatch
	ThreadsCMAs  CMABatch
	// Posts data
	Posts      CPostBatch
	PostsCATDs CATDBatch
	PostsCFGs  CFGBatch
	PostsCMAs  CMABatch
	// Misc
	LastRefreshed  int64
	LastReferenced int64
	WellFormed     bool
	// ^ For example, a threadcarrier without a thread entity matching is not well formed. We should still save it (i.e. posts might have arrived before thread entity) so that when it has fully arrived, we can flip it and start showing it.
	// Private fields only used in processing, and not saved:
	now int64
	// ^ All our refresh actions are atomic, and they refer to a single moment in time, which is the time we set at the beginning of the full refresh cycle, when we pull the whole delta needed.
}

func (c *EntityCarrier) GetPosts() CPostBatch      { return c.Posts }
func (c *EntityCarrier) GetPostsCATDs() *CATDBatch { return &c.PostsCATDs }
func (c *EntityCarrier) GetPostsCFGs() *CFGBatch   { return &c.PostsCFGs }
func (c *EntityCarrier) GetPostsCMAs() *CMABatch   { return &c.PostsCMAs }

func (c *EntityCarrier) GetThreads() CThreadBatch    { return c.Threads }
func (c *EntityCarrier) GetThreadsCATDs() *CATDBatch { return &c.ThreadsCATDs }
func (c *EntityCarrier) GetThreadsCFGs() *CFGBatch   { return &c.ThreadsCFGs }
func (c *EntityCarrier) GetThreadsCMAs() *CMABatch   { return &c.ThreadsCMAs }

func (c *EntityCarrier) GetBoards() CBoardBatch     { return c.Boards }
func (c *EntityCarrier) GetBoardsCATDs() *CATDBatch { return &c.BoardsCATDs }
func (c *EntityCarrier) GetBoardsCFGs() *CFGBatch   { return &c.BoardsCFGs }
func (c *EntityCarrier) GetBoardsCMAs() *CMABatch   { return &c.BoardsCMAs }

func (c *EntityCarrier) GetLastRefreshed() int64  { return c.LastRefreshed }
func (c *EntityCarrier) GetLastReferenced() int64 { return c.LastReferenced }

///////////////////////////////////////////
// Specific Carriers based on EntityCarrier
///////////////////////////////////////////

type BoardCarrier struct {
	EntityCarrier      `storm:"inline"`
	LSUHPublicTrusts   CPTBatch
	LSUHCanonicalNames CCNBatch
	LSUHF451s          CF451Batch
	LSUHPublicElects   CPEBatch
	// ^ These refer to the mass collections of signals we receive that are local to this specific board. In the local level processing, we pull the signals first, and only after we pull the users those signals point to, from the already-refreshed global scope.
}

func NewBoardCarrier(fp string, nowts int64) BoardCarrier {
	return BoardCarrier{
		EntityCarrier: EntityCarrier{
			Fingerprint:   fp,
			LastRefreshed: nowts,
		},
	}
}

func (c *BoardCarrier) Save() {
	// logging.Logf(1, "This is the board at the end of its refresh. %s", spew.Sdump(c))
	c.now = 0 // not useful anywhere else, no point in saving this.
	globals.KvInstance.Save(c)
}

func (c *BoardCarrier) Refresh(nowts int64) {
	c.now = nowts
	start := time.Now()
	// Make sure the statistics carrier is initialised.
	if !c.Statistics.Initialised {
		c.Statistics = NewStatisticsCarrier(2000)
	}
	// Refresh all local scope user header carriers
	c.refreshLocalUserHeadersAndStatistics()
	// Refresh content signals tables for board entity
	c.generateSignalsTablesForBoardEntity()
	// using those, refresh the board entity
	c.refreshBoardEntity()
	// refresh signals tables for the threads within.
	c.generateSignalsTablesForThreadEntities()
	// pull in new changes to the thread entities and using the local scope user headers, refresh those.
	c.refreshThreadEntities(c.Boards.GetBoardSpecificUserHeaders())
	c.applyMetas()
	c.LastReferenced = c.now
	// Save the number of posts. (Make sure that this is not based on incremental but on the total number of posts.)
	for k, _ := range c.Boards {
		if c.Boards[k].Fingerprint == c.Fingerprint {
			// Move compiled data that is needed on the client side to the compiled object so it can be transmitted over.
			c.Boards[k].ThreadsCount = len(c.Threads)
			c.Boards[k].UserCount = c.Statistics.UserCount
		}
	}
	c.Save()
	elapsed := time.Since(start)
	logging.Logf(2, "Compiling the board %s took %s", c.Fingerprint, elapsed)
}

func (c *BoardCarrier) applyMetas() {
	// todo
}

// refreshThreadEntities refreshes the threads that we have in a board container. Compiled versions of these entities will then be passed to each thread that we want to compile.
// refreshThreadEntities refreshes the main thread entity of this carrier. Note that the name 'threads' is still plural, but that's mostly for convenience and reuse - there should be only one thread at the end of this process.
func (c *BoardCarrier) refreshThreadEntities(boardSpecificUserHeaders CUserBatch) {
	logging.Logf(2, "Refresh thread entities in board hits.")
	newThreadEntitiesInThread := beapiconsumer.GetThreads(c.GetLastReferenced(), c.now, []string{}, c.Fingerprint, false, false)
	// for _, val := range newThreadEntitiesInThread {
	// 	fmt.Println(val)
	// }
	logging.Logf(2, "number of New Thread entities in thread: %#v, board name: %v", len(newThreadEntitiesInThread), c.Boards[0].Name)
	// ^ This can actually be plural, if we receive two updates to the same thread, etc.
	c.Threads.InsertFromProtobuf(newThreadEntitiesInThread)
	c.Threads.Refresh(c.GetThreadsCATDs(), c.GetThreadsCFGs(), c.GetThreadsCMAs(), boardSpecificUserHeaders, c.now)
	// If there is a parent, then we've actually managed to find a thread entity, which means the thread entity actually exists, which means this is a valid container.
	c.WellFormed = true
	for k, _ := range c.Threads {
		if len(c.Threads[k].Fingerprint) == 0 {
			c.WellFormed = false
			break
		}
	}
}

func (c *BoardCarrier) generateSignalsTablesForThreadEntities() {
	genSigTables("board", c.Fingerprint, "", c.GetLastReferenced(), c.now, c.GetThreadsCATDs(), c.GetThreadsCFGs(), c.GetThreadsCMAs(), true)
}

func (c *BoardCarrier) generateSignalsTablesForBoardEntity() {
	genSigTables("", "", c.Fingerprint, c.GetLastReferenced(), c.now, c.GetBoardsCATDs(), c.GetBoardsCFGs(), c.GetBoardsCMAs(), false)
}

// This refreshes the board entities in this board carrier.
func (c *BoardCarrier) refreshBoardEntity() {
	newBoardEntitiesInBoard := beapiconsumer.GetBoards(c.GetLastReferenced(), c.now, []string{c.Fingerprint}, false, false)
	c.Boards.InsertFromProtobuf(newBoardEntitiesInBoard)
	c.Boards.Refresh(c.GetBoardsCATDs(), c.GetBoardsCFGs(), c.GetBoardsCMAs(), c.Boards.GetBoardSpecificUserHeaders(), c.now)
	c.WellFormed = true
}

// refreshLocalUserHeadersAndStatistics handles the statistics updates and local user header updates. It does both because statistics count increase relies on user headers coming in, so it's more efficient than making that call two times.
func (c *BoardCarrier) refreshLocalUserHeadersAndStatistics() {
	// Get & save all the delta signals we have for this local scope.
	c.generateSignalsTablesForLocalUserHeaderScope()
	// From those signals, generate the user headers we're going to need.
	neededUserHeaderFps := c.generateNeededUserHeaderFingerprints()
	// Count all of those needed user headers into the population
	c.Statistics.Refresh(neededUserHeaderFps)
	// Refresh our local user headers bucket for this board before we start refreshing, so the refresher routines will have all the UHs they will need already in the bucket.
	c.refreshLSUserHeadersBucket(neededUserHeaderFps, c.Fingerprint)
	// For all the headers we have, apply the local signals, so that the local signals will be updated on this entity.
	c.refreshLocalScopeUserHeadersWithLocalSignals()
	// Look at all the user headers we have, and refresh them using the *global* entity updates and signal updates. This is so that we pull in any changes that happened on the global level.
	c.refreshLocalScopeUserHeadersWithGlobalUpdatesAndSignals()
}

func (c *BoardCarrier) getLocalDefaultMods() []string {
	return c.Boards.GetDefaultMods()
}

// refreshLocalScopeUserHeadersWithLocalSignals refreshes the local signals of the user header entities brought forward in this specific delta.
func (c *BoardCarrier) refreshLocalScopeUserHeadersWithLocalSignals() {
	// for every board in this board carrier
	for k, _ := range c.Boards {
		// Make it so that every user header is refreshed via the signals we have. We don't want to do a full insert - our headers' contents are already updated.
		for j, _ := range c.Boards[k].LocalScopeUserHeaders {
			c.Boards[k].LocalScopeUserHeaders[j].RefreshUserSignals(&c.LSUHPublicTrusts, &c.LSUHCanonicalNames, &c.LSUHF451s, &c.LSUHPublicElects, c.getLocalDefaultMods(), c.Fingerprint, c.Statistics.UserCount)
		}
	}
}

// refreshLocalScopeUserHeadersWithGlobalUpdatesAndSignals refreshes the user headers in the local scope with the content and signals from the global scope. Since global scope user headers update ran before, by doing this, we don't have to pull updates from the backend again, the content will automatically update. By the virtue of that, the global signals within those entities will also be updated, making us ready for a delta insert of the local signals.
func (c *BoardCarrier) refreshLocalScopeUserHeadersWithGlobalUpdatesAndSignals() {
	// for every board in this board carrier
	for k, _ := range c.Boards {
		// Make it so that every user header is refreshed via its global counterpart.
		for j, _ := range c.Boards[k].LocalScopeUserHeaders {
			globalUserHeader := CompiledUser{}
			err := globals.KvInstance.One("Fingerprint", c.Boards[k].LocalScopeUserHeaders[j].Fingerprint, &globalUserHeader)
			if err != nil {
				logging.Logf(1, "We've failed to get the global counterpart of this local user header. Error: %v", err)
				continue
			}
			c.Boards[k].LocalScopeUserHeaders[j].InsertWithSignalMerge(globalUserHeader)
		}
	}
}

func (c *BoardCarrier) generateSignalsTablesForLocalUserHeaderScope() {
	c.LSUHPublicTrusts.Insert(GetPTs("", c.Fingerprint, c.LastReferenced, c.now), c.now)
	c.LSUHCanonicalNames.Insert(GetCNs("", c.Fingerprint, c.LastReferenced, c.now), c.now)
	c.LSUHF451s.Insert(GetF451s("", c.Fingerprint, c.LastReferenced, c.now), c.now)
	c.LSUHPublicElects.Insert(GetPEs("", c.Fingerprint, c.LastReferenced, c.now), c.now)
}

func (c *BoardCarrier) generateNeededUserHeaderFingerprints() []string {
	uhfps := make(map[string]bool)
	for k, _ := range c.LSUHPublicTrusts {
		uhfps[c.LSUHPublicTrusts[k].TargetFingerprint] = true
	}
	for k, _ := range c.LSUHCanonicalNames {
		uhfps[c.LSUHCanonicalNames[k].TargetFingerprint] = true
	}
	for k, _ := range c.LSUHF451s {
		uhfps[c.LSUHF451s[k].TargetFingerprint] = true
	}
	for k, _ := range c.LSUHPublicElects {
		uhfps[c.LSUHPublicElects[k].TargetFingerprint] = true
	}
	uhfpsSlice := []string{}
	for key, _ := range uhfps {
		uhfpsSlice = append(uhfpsSlice, key)
	}
	return uhfpsSlice
}

// refreshLSUserHeadersBucket refreshes the slice in the user headers of this specific board. This is to make sure that when we refresh the signals and entities, this list will cover all of the user headers needed to display any content within.
func (c *BoardCarrier) refreshLSUserHeadersBucket(targetsfp []string, boardfp string) {
	// Attempt to scan the board's local user headers for this specific board.
	bi := c.Boards.Find(boardfp)
	b := CompiledBoard{}
	if bi != -1 {
		b = c.Boards[bi]
	}
	// For each target fp, make sure they either exist in our local set, or find it from the global set and add it to the local set.
	for _, targetfp := range targetsfp {
		for k, _ := range b.LocalScopeUserHeaders {
			if b.LocalScopeUserHeaders[k].Fingerprint == targetfp {
				return
				// We have the uh in the board. No need to do anything.
			}
		}
		// Not found in those. Pull it from globals, and return the userheader and its signals batches.
		uhc := UserHeaderCarrier{}
		err := globals.KvInstance.One("Fingerprint", targetfp, &uhc)
		if err != nil {
			logging.Logf(1, "We could not get the requested user from the global user headers. Error: %v", err)
		}
		// useri := uhc.Users.Find(targetfp)
		// uh := CompiledUser{}
		// if useri != -1 {
		// 	uh = uhc.Users[useri]
		// }
		b.LocalScopeUserHeaders.InsertWithSignalMerge(uhc.Users)
		// ^ Add it to the user headers of the current board.
	}
}

func (c *BoardCarrier) ConstructAmbientBoards() []AmbientBoard {
	abs := []AmbientBoard{}
	for key, _ := range c.Boards {
		abs = append(abs, c.Boards[key].ConvertToAmbientBoard())
	}
	return abs
}

type BCBatch []BoardCarrier

func (c *BCBatch) Find(fp string) int {
	for k, _ := range *c {
		if (*c)[k].Fingerprint == fp {
			return k
		}
	}
	return -1
}

// Distribute the oncoming board entities to existing boards.
func (c *BCBatch) Insert(newBoardEntities []*pbstructs.Board) {
	for k, _ := range newBoardEntities {
		if i := c.Find(newBoardEntities[k].GetProvable().GetFingerprint()); i != -1 {
			(*c)[i].Boards.InsertFromProtobuf([]*pbstructs.Board{newBoardEntities[k]})
		} else {
			bc := NewBoardCarrier(newBoardEntities[k].GetProvable().GetFingerprint(), globals.FrontendTransientConfig.RefresherCacheNowTimestamp)
			bc.Boards.InsertFromProtobuf([]*pbstructs.Board{newBoardEntities[k]})
			(*c) = append((*c), bc)
		}
	}
}

type ThreadCarrier struct {
	EntityCarrier `storm:"inline"`
	LinkedPosts   []CompiledPost
}

func NewThreadCarrier(fp, parentfp string, nowts int64) ThreadCarrier {
	return ThreadCarrier{
		EntityCarrier: EntityCarrier{
			Fingerprint:       fp,
			ParentFingerprint: parentfp,
			LastRefreshed:     nowts,
		},
	}
}

func (c *ThreadCarrier) Save() {
	c.now = 0
	globals.KvInstance.Save(c)
}

func (c *ThreadCarrier) getRefreshedUser(ownerfp string) *UserHeaderCarrier {
	// Get the user from global, add this specific board's user signals on it, and serve it in.
	// todo
	return &UserHeaderCarrier{}
}

// refreshThreadEntities refreshes the main thread entity of this carrier. Note that the name 'threads' is still plural, but that's mostly for convenience and reuse - there should be only one thread at the end of this process.
func (c *ThreadCarrier) refreshThreadEntities(boardSpecificUserHeaders CUserBatch) {
	// This is our logic:
	// - Update the thread entity itself (if the text has changed, etc.)
	// - Using the signal buckets we've generated prior, update the signals entity in the thread (in refresh())
	// - Using the userheaders generated globally before, and the local userheaders generated within the board scope before, update the user entity. i.e if a user is promoted to mod, this should reflect this change, etc.
	newThreadEntitiesInThread := beapiconsumer.GetThreads(c.GetLastReferenced(), c.now, []string{c.Fingerprint}, "", false, false)
	// ^ This can actually be plural, if we receive two updates to the same thread, etc.
	c.Threads.InsertFromProtobuf(newThreadEntitiesInThread)
	c.Threads.Refresh(c.GetThreadsCATDs(), c.GetThreadsCFGs(), c.GetThreadsCMAs(), boardSpecificUserHeaders, c.now)
	// If there is a parent, then we've actually managed to find a thread entity, which means the thread entity actually exists, which means this is a valid container.
	allWellFormed := true
	for k, _ := range c.Threads {
		if len(c.Threads[k].Fingerprint) == 0 {
			allWellFormed = false
			break
		}
	}
	if allWellFormed {
		c.ParentFingerprint = c.Threads[0].Board
	}
}

func (c *ThreadCarrier) generateSignalsTablesForThreadEntity() {
	genSigTables("", "", c.Fingerprint, c.GetLastReferenced(), c.now, c.GetThreadsCATDs(), c.GetThreadsCFGs(), c.GetThreadsCMAs(), false)
}

func (c *ThreadCarrier) refreshPosts(boardSpecificUserHeaders CUserBatch) {
	// Get new delta
	newPostsInThread := beapiconsumer.GetPosts(c.GetLastReferenced(), c.now, []string{}, c.Fingerprint, "thread", false, false)
	// Add delta to extant pool
	c.Posts.InsertFromProtobuf(newPostsInThread)
	// Refresh each post with signals that we have already compiled in this pass.
	c.Posts.Refresh(c.GetPostsCATDs(), c.GetPostsCFGs(), c.GetPostsCMAs(), boardSpecificUserHeaders, c.now)
}

func (c *ThreadCarrier) generateSignalsTablesForPostsInThread() {
	genSigTables("thread", c.Fingerprint, "", c.GetLastReferenced(), c.now, c.GetPostsCATDs(), c.GetPostsCFGs(), c.GetPostsCMAs(), false)
}

type TCBatch []ThreadCarrier

func (c *TCBatch) Find(fp string) int {
	for k, _ := range *c {
		if (*c)[k].Fingerprint == fp {
			return k
		}
	}
	return -1
}

// Distribute the oncoming board entities to existing boards.
func (c *TCBatch) Insert(newThreadEntities []*pbstructs.Thread) {
	for k, _ := range newThreadEntities {
		if i := c.Find(newThreadEntities[k].GetProvable().GetFingerprint()); i != -1 {
			(*c)[i].Threads.InsertFromProtobuf([]*pbstructs.Thread{newThreadEntities[k]})
		} else {
			threadCarrier := NewThreadCarrier(newThreadEntities[k].GetProvable().GetFingerprint(), newThreadEntities[k].GetBoard(), globals.FrontendTransientConfig.RefresherCacheNowTimestamp)
			threadCarrier.Threads.InsertFromProtobuf([]*pbstructs.Thread{newThreadEntities[k]})
			(*c) = append((*c), threadCarrier)
		}
	}
}

// This function is here mostly because it's used here. Basically, you can give this function any batches, and it will add the result of the update from the backend to the batches you've provided. This is useful for most types of carriers.
func genSigTables(voteParentType, parentFp, targetFp string, lastRef int64, nowts int64, extantCATDs *CATDBatch, extantCFGs *CFGBatch, extantCMAs *CMABatch, noDescendants bool) {
	// logging.Logf(1, "GenSigTables was called for: voteparenttype: %v, parentfp: %v, targetfp: %v, lastref: %v, now: %v", voteParentType, parentFp, targetFp, lastRef, nowts)
	newATDs := GetATDs(parentFp, voteParentType, targetFp, lastRef, nowts, noDescendants)
	extantCATDs.Insert(newATDs, nowts)
	newFGs := GetFGs(parentFp, voteParentType, targetFp, lastRef, nowts, noDescendants)
	extantCFGs.Insert(newFGs, nowts)
	newMAs := GetMAs(parentFp, voteParentType, targetFp, lastRef, nowts, noDescendants)
	extantCMAs.Insert(newMAs, nowts)
}

// func (c *ThreadCarrier) Reset() {
// 	cNew := NewThreadCarrier(c.Fingerprint, c.ParentFingerprint, c.now)
// 	c = &cNew
// }

// Refresh refreshes all the data in the carrier from backend, compiles the signals and does it on an incremental basis based on the lastreferenced, so only the new information is requested, processed, and compiled.
func (c *ThreadCarrier) Refresh(boardSpecificUserHeaders CUserBatch, bc *BoardCarrier, nowts int64) {
	c.now = nowts
	if bc == nil {
		// Generate signals tables for the thread entity itself.
		c.generateSignalsTablesForThreadEntity()
		// Using the signals tables and refreshed users list, refresh the thread entity itself.
		c.refreshThreadEntities(CUserBatch{})
	} else { // If it's precomputed
		ti := bc.Threads.Find(c.Fingerprint)
		if ti != -1 {
			c.Threads = CThreadBatch{bc.Threads[ti]}
		}
	}
	// Get raw signals for any content within this thread
	c.generateSignalsTablesForPostsInThread()
	// Using the signals tables and refreshed users list, refresh the posts in this thread
	c.refreshPosts(boardSpecificUserHeaders)
	// Apply flags
	c.applyMetas()
	// Apple metas
	// Set the last referenced to now, so next refresh will use it as a base.
	c.LastReferenced = c.now
	// Save the number of posts. (Make sure that this is not based on incremental but on the total number of posts.)
	for k, _ := range c.Threads {
		if c.Threads[k].Fingerprint == c.Fingerprint {
			c.Threads[k].PostsCount = len(c.Posts)
		}
	}
	// Save it to the kvstore.
	c.Save()
}

// // ForceRefresh refreshes the carrier with the assumption that it was never refreshed before, capturing the whole backend database. It will keep the data in the carrier, and only add if there's anything missing. Useful for debug, should not be needed in production.
// func (c *ThreadCarrier) ForceRefresh(boardSpecificUserHeaders CUserBatch, bc *BoardCarrier, nowts int64) {
// 	c.LastReferenced = 0
// 	c.Refresh(boardSpecificUserHeaders, bc, nowts)
// }

// // ResetRefresh is like force refresh, but before refreshing, it actually deletes information in the carrier. Useful for debugging, should not be needed in production.
// func (c *ThreadCarrier) ResetRefresh(boardSpecificUserHeaders CUserBatch, bc *BoardCarrier, nowts int64) {
// 	c.Reset()
// 	c.Refresh(boardSpecificUserHeaders, bc, nowts)
// }

func (c *ThreadCarrier) applyMetas() {
	// todo
}

///////////////////////////////////////////
// Specific Carriers based on Users
///////////////////////////////////////////

type UserHeaderCarrier struct {
	Fingerprint    string     `storm:"id"`
	Domain         string     // If it's in the local scope of a board, and which
	Users          CUserBatch // Plural out of convention, but a user header carries one user.
	PublicTrusts   CPTBatch
	CanonicalNames CCNBatch
	F451s          CF451Batch
	PublicElects   CPEBatch
	LastRefreshed  int64
	LastReferenced int64
	Self           bool
	now            int64
}

func NewUserHeaderCarrier(fp, domain string, nowts int64) UserHeaderCarrier {
	return UserHeaderCarrier{
		Fingerprint:   fp,
		Domain:        domain,
		LastRefreshed: nowts,
	}
}

// needs to run before user entity refresh
func (c *UserHeaderCarrier) refreshSignalsTables() {
	c.PublicTrusts.Insert(GetPTs(c.Fingerprint, c.Domain, c.LastReferenced, c.now), c.now)
	c.CanonicalNames.Insert(GetCNs(c.Fingerprint, c.Domain, c.LastReferenced, c.now), c.now)
	c.F451s.Insert(GetF451s(c.Fingerprint, c.Domain, c.LastReferenced, c.now), c.now)
	c.PublicElects.Insert(GetPEs(c.Fingerprint, c.Domain, c.LastReferenced, c.now), c.now)
}

func (c *UserHeaderCarrier) refreshUserEntity(localDefaultMods []string, totalPop int) {
	newUserEntities := beapiconsumer.GetKeys(c.LastReferenced, c.now, []string{c.Fingerprint}, false, false)
	// Attempt to insert every user update we have - the newest will prevail. There will most likely be only one in reality, but it's good to be defensive.
	c.Users.InsertFromProtobuf(newUserEntities, c.now)
	c.Users.Refresh(&c.PublicTrusts, &c.CanonicalNames, &c.F451s, &c.PublicElects, localDefaultMods, c.Domain, totalPop)
}

// Merge merges two different user header carriers from different privilege levels into one single whole. Merging with itself or the same UHC should not have an effect. todo
// func (c *UserHeaderCarrier) Merge(mc UserHeaderCarrier) {
// 	// todo - this is needed when calculating board user headers, which will be used at boards level and threads level.

// 	// this could also be something like 'updateWithGlobalUserData' excl for local user carriers.
// }

// This is basically only useful when we're saving these locally, and have a batch of these already at hand. In all other places, the kvstore returns them as singular.
type UHCBatch []UserHeaderCarrier

func (c *UHCBatch) Find(fp string) int {
	for k, _ := range *c {
		if (*c)[k].Fingerprint == fp {
			return k
		}
	}
	return -1
}

// We could also do something like total karma, but then that would be a little tricky.. because what happens is that it would both be probabilistic, and might hit the ceiling. Also might be fairly drastically different for different people because we're deleting the votes after 2 weeks.

// func (c *UserHeaderCarrier) Reset() {
// 	cNew := NewUserHeaderCarrier(c.Fingerprint, c.Domain, c.now)
// 	c = &cNew
// }

func (c *UserHeaderCarrier) Refresh(localDefaultMods []string, totalPop int, nowts int64) {
	c.now = nowts
	// Generate signals tables we need to use
	c.refreshSignalsTables()
	// Using those tables, refresh the user entity
	c.refreshUserEntity(localDefaultMods, totalPop)
	c.LastReferenced = c.now
	c.Save()
}

// func (c *UserHeaderCarrier) ResetRefresh(localDefaultMods []string, totalPop int, nowts int64) {
// 	c.Reset()
// 	c.Refresh(localDefaultMods, totalPop, nowts)
// }

// func (c *UserHeaderCarrier) ForceRefresh(localDefaultMods []string, totalPop int, nowts int64) {
// 	c.LastReferenced = 0
// 	c.Refresh(localDefaultMods, totalPop, nowts)
// }

func (c *UserHeaderCarrier) Save() {
	c.now = 0
	globals.KvInstance.Save(c)
}

func (c *UserHeaderCarrier) applyMetas() {
	// todo
}

type StatisticsCarrier struct {
	UserCount      int
	UserCountBloom rollingbloom.RollingBloom
	Initialised    bool
}

func NewStatisticsCarrier(bloomsize uint) StatisticsCarrier {
	return StatisticsCarrier{
		UserCountBloom: rollingbloom.NewRollingBloom(uint(globals.FrontendConfig.GetNetworkMemoryDays()), uint(globals.FrontendConfig.GetNetworkHeadDays()), bloomsize),
		Initialised:    true,
	}
}

func (g *StatisticsCarrier) Refresh(newUserFingerprints []string) {
	// Put their fingerprints into the bloom
	for k, _ := range newUserFingerprints {
		if len(newUserFingerprints[k]) > 0 {
			g.UserCountBloom.AddString(newUserFingerprints[k])
		}
	}
	g.UserCount = g.UserCountBloom.Count()
}

type GlobalStatisticsCarrier struct {
	Id                    int `storm:"id"` // always 1, this is a singleton.
	LastReferenced        int64
	RefreshStartTimestamp int64
	now                   int64
	StatisticsCarrier
}

func NewGlobalStatisticsCarrier() GlobalStatisticsCarrier {
	return GlobalStatisticsCarrier{
		Id:                1,
		StatisticsCarrier: NewStatisticsCarrier(10000),
	}
}

func (g *GlobalStatisticsCarrier) Save() {
	globals.KvInstance.Save(g)
}

func (g *GlobalStatisticsCarrier) Refresh(nowts int64) []*pbstructs.Key {
	g.now = nowts
	// Get all new user entities / updates since lastref
	newUserEntities := beapiconsumer.GetKeys(g.LastReferenced, g.now, []string{}, false, false)
	fps := []string{}
	// Put their fingerprints into the bloom
	for k, _ := range newUserEntities {
		if fp := newUserEntities[k].GetProvable().GetFingerprint(); len(fp) > 0 {
			fps = append(fps, fp)
		}
	}
	g.StatisticsCarrier.Refresh(fps)
	g.LastReferenced = g.now
	g.Save()
	return newUserEntities
}
