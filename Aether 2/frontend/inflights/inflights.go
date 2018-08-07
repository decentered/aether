// This package receives the content and signals emitted by the client and converts it into a queue that is fed into the minter.

package inflights

import (
	"aether-core/frontend/clapiconsumer"
	"aether-core/protos/clapi"
	"aether-core/protos/feapi"
	beObj "aether-core/protos/mimapi"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"sync"
	"time"
)

var Inflights *inflights

// ^ This is our instance store - we want to use a single instance

/*----------  Base carrier  ----------*/

type inflights struct {
	lock                sync.Mutex
	ingestLock          sync.Mutex
	ingestRanOnce       bool
	ingestRunning       bool
	ID                  int
	InflightBoards      []InflightBoard
	InflightThreads     []InflightThread
	InflightPosts       []InflightPost
	InflightVotes       []InflightVote
	InflightKeys        []InflightKey
	InflightTruststates []InflightTruststate
	/*----------  Complete entries  ----------*/
	FulfilledBoards      []InflightBoard
	FulfilledThreads     []InflightThread
	FulfilledPosts       []InflightPost
	FulfilledVotes       []InflightVote
	FulfilledKeys        []InflightKey
	FulfilledTruststates []InflightTruststate
}

/*----------  Protobuf conversions  ----------*/

func (o *inflights) Protobuf() *clapi.Inflights {
	o.lock.Lock()
	defer o.lock.Unlock()
	opb := clapi.Inflights{}
	for k, _ := range o.InflightBoards {
		opb.Boards = append(opb.Boards, o.InflightBoards[k].Protobuf())
	}
	for k, _ := range o.InflightThreads {
		opb.Threads = append(opb.Threads, o.InflightThreads[k].Protobuf())
	}
	for k, _ := range o.InflightPosts {
		opb.Posts = append(opb.Posts, o.InflightPosts[k].Protobuf())
	}
	for k, _ := range o.InflightVotes {
		opb.Votes = append(opb.Votes, o.InflightVotes[k].Protobuf())
	}
	for k, _ := range o.InflightKeys {
		opb.Keys = append(opb.Keys, o.InflightKeys[k].Protobuf())
	}
	for k, _ := range o.InflightTruststates {
		opb.Truststates = append(opb.Truststates, o.InflightTruststates[k].Protobuf())
	}
	return &opb
}

func (o *InflightStatus) Protobuf() *clapi.InflightStatus {
	return &clapi.InflightStatus{
		CompletionPercent:   int32(o.CompletionPercent),
		StatusText:          o.StatusText,
		RequestedTimestamp:  o.RequestedTimestamp,
		LastActionTimestamp: o.LastActionTimestamp,
		EventType:           o.EventType,
	}
}

func (o *InflightBoard) Protobuf() *clapi.InflightBoard {
	return &clapi.InflightBoard{
		Status: o.Status.Protobuf(),
		Entity: &o.Entity,
	}
}

func (o *InflightThread) Protobuf() *clapi.InflightThread {
	return &clapi.InflightThread{
		Status: o.Status.Protobuf(),
		Entity: &o.Entity,
	}
}

func (o *InflightPost) Protobuf() *clapi.InflightPost {
	return &clapi.InflightPost{
		Status: o.Status.Protobuf(),
		Entity: &o.Entity,
	}
}

func (o *InflightVote) Protobuf() *clapi.InflightVote {
	return &clapi.InflightVote{
		Status: o.Status.Protobuf(),
		Entity: &o.Entity,
	}
}

func (o *InflightKey) Protobuf() *clapi.InflightKey {
	return &clapi.InflightKey{
		Status: o.Status.Protobuf(),
		Entity: &o.Entity,
	}
}

func (o *InflightTruststate) Protobuf() *clapi.InflightTruststate {
	return &clapi.InflightTruststate{
		Status: o.Status.Protobuf(),
		Entity: &o.Entity,
	}
}

/*----------  Status header  ----------*/

type InflightStatus struct {
	CompletionPercent   int // 0-100
	StatusText          string
	RequestedTimestamp  int64 // We grab the oldest requested to start the process
	LastActionTimestamp int64
	EventType           string
}

func (s *InflightStatus) Fulfilled() bool {
	if s.CompletionPercent == 100 || s.CompletionPercent == -1 {
		return true
	}
	return false
}

const (
	STATUS_WAITING                   = "Waiting for processing."
	STATUS_MINTING                   = "Minting proof-of-work for the entity..."
	STATUS_ADDING_TO_BACKEND         = "Adding to the local backend"
	STATUS_WAITING_TO_SERVE          = "Waiting for a remote inbound to serve entity"
	STATUS_WAITING_TO_FIND_IN_REMOTE = "Waiting to see a remote node serving the entity to confirm distribution"
	STATUS_RECOMPILING_FRONTEND      = "Inserting into graph..."
	STATUS_COMPLETE                  = "Successfully posted."
	STATUS_REMOTE_COMPLETE           = "The entity is communicated to the network and its availability is verified."
	STATUS_FAILED                    = "The insertion for this entity has failed. If you see this and you didn't expect to, please reach out, it's worth a look."
)

var statusesOrdered = []string{
	STATUS_WAITING,
	STATUS_MINTING,
	// STATUS_ADDING_TO_BACKEND,
	STATUS_RECOMPILING_FRONTEND,
	// STATUS_WAITING_TO_SERVE,
	// STATUS_WAITING_TO_FIND_IN_REMOTE,
	STATUS_COMPLETE,
	// STATUS_REMOTE_COMPLETE,
}

func NewInflightStatus(status string, etype string) InflightStatus {
	st := InflightStatus{
		StatusText:          status,
		RequestedTimestamp:  time.Now().Unix(),
		LastActionTimestamp: time.Now().Unix(),
		EventType:           etype,
	}
	st.setCompletionPercent()
	return st
}

func (o *InflightStatus) Update(status string) {
	o.StatusText = status
	o.LastActionTimestamp = time.Now().Unix()
	o.setCompletionPercent()
}

func (o *InflightStatus) setCompletionPercent() {
	if o.StatusText == STATUS_FAILED {
		o.CompletionPercent = -1
		return
	}
	clp := 0
	for k, _ := range statusesOrdered {
		if statusesOrdered[k] == o.StatusText {
			clp = (k * 100) / (len(statusesOrdered) - 1)
			// ^ If you have 6 items, k will be 0-5 range, if you don't do -1 it will never be 100%.
			if clp == 0 {
				clp = 1
				// So that the progress bar won't look broken, with nothing in it.
			}
			break
		}
	}
	o.CompletionPercent = clp
}

/*----------  Inflight types  ----------*/

type InflightBoard struct {
	Status InflightStatus
	Entity beObj.Board
}

type InflightThread struct {
	Status InflightStatus
	Entity beObj.Thread
}

type InflightPost struct {
	Status InflightStatus
	Entity beObj.Post
}

type InflightVote struct {
	Status InflightStatus
	Entity beObj.Vote
}

type InflightKey struct {
	Status InflightStatus
	Entity beObj.Key
}

type InflightTruststate struct {
	Status InflightStatus
	Entity beObj.Truststate
}

/*----------  Read from and write to KvStore  ----------*/

func GetInflights() *inflights {
	if Inflights != nil {
		return Inflights
	}
	o := inflights{}
	err := globals.KvInstance.One("ID", 1, &o)
	if err != nil && err.Error() != "not found" {
		logging.Logf(1, "An error occurred while getting the inflights from KvInstance. Error: %v", err)
	}
	if !o.ingestRanOnce {
		go o.Ingest()
	}
	Inflights = &o
	return Inflights
}

func (o *inflights) commit() {
	o.ID = 1 // Always singleton
	globals.KvInstance.Save(o)
}

func (o *inflights) ManualSaveToKvStore() {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.commit()
}

/*----------  Insert & get next from the stack  ----------*/
// TODO: This is where we add preempts for changing a vote

func (o *inflights) Insert(input interface{}) {
	o.lock.Lock()
	defer o.lock.Unlock()
	switch i := input.(type) {
	case feapi.ContentEventPayload:
		// CREATE or UPDATE for Boards, Threads, Posts, Keys
		if i.GetBoardData() != nil {
			ifObj := createInflightBoard(&i)
			o.InflightBoards = append(o.InflightBoards, ifObj)
			o.commit()
			go o.Ingest()
			return
		}
		if i.GetThreadData() != nil {
			ifObj := createInflightThread(&i)
			o.InflightThreads = append(o.InflightThreads, ifObj)
			o.commit()
			go o.Ingest()
			return
		}
		if i.GetPostData() != nil {
			ifObj := createInflightPost(&i)
			o.InflightPosts = append(o.InflightPosts, ifObj)
			o.commit()
			go o.Ingest()
			return
		}
		if i.GetKeyData() != nil {
			ifObj := createInflightKey(&i)
			o.InflightKeys = append(o.InflightKeys, ifObj)
			o.commit()
			go o.Ingest()
			return
		}
	case feapi.SignalEventPayload:
		// CREATE or UPDATE for Votes, Truststates
		if targetType := i.GetSignalTargetType(); targetType == feapi.SignalTargetType_CONTENT {
			ifObj := createInflightVote(&i)
			o.InflightVotes = append(o.InflightVotes, ifObj)
			o.cleanRepeatVotes()
			o.commit()
			go o.Ingest()
			return
		}
		if targetType := i.GetSignalTargetType(); targetType == feapi.SignalTargetType_USER {
			ifObj := createInflightTruststate(&i)
			o.InflightTruststates = append(o.InflightTruststates, ifObj)
			o.cleanRepeatTruststates()
			o.commit()
			go o.Ingest()
			return
		}
	default:
		logging.Logf(1, "The type of event payload requested to be inserted into the inflights queue could not be determined. Event: %#v", i)
	}
}

/*----------  Repeat signal cleaners  ----------*/

// This function goes through all the votes in the queue and if there are more than one pointing at the same item, it only leaves the most recent. The idea is that you can change your mind before it hits the mint.
func (o *inflights) cleanRepeatVotes() {
	// First, create a list of unprocessed and in progress votes
	unprocessedVotes := []InflightVote{} // votes that hasn't started processing
	inProgressVotes := []InflightVote{}  // votes that have
	for k, _ := range o.InflightVotes {
		if o.InflightVotes[k].Status.StatusText == STATUS_WAITING {
			unprocessedVotes = append(unprocessedVotes, o.InflightVotes[k])
			continue
		}
		inProgressVotes = append(inProgressVotes, o.InflightVotes[k])
	}
	// Create a map of the latest versions of those unprocessed votes
	signalsTargetsInWaiting := make(map[string]int64) //target:timestamp
	for k, _ := range unprocessedVotes {
		if signalsTargetsInWaiting[unprocessedVotes[k].Entity.GetTarget()] <= unprocessedVotes[k].Status.RequestedTimestamp {
			// "<=" bc. if it's coming later, likely it happened later.
			signalsTargetsInWaiting[unprocessedVotes[k].Entity.GetTarget()] = unprocessedVotes[k].Status.RequestedTimestamp
		}
	}
	// If they're the latest versions, grab them from the unproc votes list.
	dedupedVotes := []InflightVote{}
	for k, _ := range unprocessedVotes {
		/*
		   unprocessedVotes[k] = scan through 0 1 2 3 4 5 ..
		   unprocessedVotes[len(unprocessedVotes)-1-k] = scan thru .. 5 4 3 2 1 0
		   reverse is important, because we want to prefer the latter addition if both happened on the same second, which by the virtue of it being latter in the queue, a newer one.
		*/
		if unprocessedVotes[len(unprocessedVotes)-1-k].Status.RequestedTimestamp == signalsTargetsInWaiting[unprocessedVotes[len(unprocessedVotes)-1-k].Entity.GetTarget()] {
			dedupedVotes = append(dedupedVotes, unprocessedVotes[len(unprocessedVotes)-1-k])
			// And clean out the map, so no other vote can enter through the same fp:ts pair.
			delete(signalsTargetsInWaiting, unprocessedVotes[len(unprocessedVotes)-1-k].Entity.GetTarget())

		}
	}
	// Add deduped votes into in progress votes and attach it back to the base inflights
	o.InflightVotes = append(inProgressVotes, dedupedVotes...)
}

func (o *inflights) cleanRepeatTruststates() {
	// First, create a list of unprocessed and in progress truststates
	unprocessedTses := []InflightTruststate{} // tses that hasn't started processing
	inProgressTses := []InflightTruststate{}  // tses that have
	for k, _ := range o.InflightTruststates {
		if o.InflightTruststates[k].Status.StatusText == STATUS_WAITING {
			unprocessedTses = append(unprocessedTses, o.InflightTruststates[k])
			continue
		}
		inProgressTses = append(inProgressTses, o.InflightTruststates[k])
	}
	// Create a map of the latest versions of those unprocessed tses
	signalsTargetsInWaiting := make(map[string]int64) //target:timestamp
	for k, _ := range unprocessedTses {
		if signalsTargetsInWaiting[unprocessedTses[k].Entity.GetTarget()] < unprocessedTses[k].Status.RequestedTimestamp {
			// "<=" bc. if it's coming later, likely it happened later.
			signalsTargetsInWaiting[unprocessedTses[k].Entity.GetTarget()] = unprocessedTses[k].Status.RequestedTimestamp
		}
	}
	// If they're the latest versions, grab them from the unproc votes list.
	dedupedTs := []InflightTruststate{}
	for k, _ := range unprocessedTses {
		if unprocessedTses[len(unprocessedTses)-1-k].Status.RequestedTimestamp == signalsTargetsInWaiting[unprocessedTses[len(unprocessedTses)-1-k].Entity.GetTarget()] {
			dedupedTs = append(dedupedTs, unprocessedTses[len(unprocessedTses)-1-k])
			// And clean out the map, so no other vote can enter through the same fp:ts pair.
			delete(signalsTargetsInWaiting, unprocessedTses[len(unprocessedTses)-1-k].Entity.GetTarget())
		}
	}
	// Add deduped votes into in progress votes and attach it back to the base inflights
	o.InflightTruststates = append(inProgressTses, dedupedTs...)
}

/*----------  Type class and type parsers for signals  ----------*/

func parseTypeClass(tc feapi.SignalTypeClass) int32 {
	// bs, ints := tc.EnumDescriptor()
	val := feapi.SignalTypeClass_value[tc.String()]
	switch val {
	/*----------  Vote type classes  ----------*/
	case 1: // atd
		return 1
	case 2: // fg
		return 2
	case 3: // ma
		return 3

	/*----------  Truststate type classes  ----------*/
	case 4: // pt
		return 1
	case 5: // naming
		return 2
	case 6: // f451
		return 3
	case 7: // pe
		return 4
	default:
		return 0
	}
}

func parseType(t feapi.SignalType) int32 {
	val := feapi.SignalType_value[t.String()]
	switch val {
	/*----------  Vote types  ----------*/
	/**----------  ATD vote types  ----------*/
	case 1: // upvote
		return 1
	case 2: // downvote
		return 2
	/**----------  FG vote types  ----------*/
	case 3: // report to mod
		return 1
	/**----------  MA vote types  ----------*/
	case 4: // modblock
		return 1
	case 5: // modapprove
		return 2
	/*----------  Truststate types  ----------*/
	/**----------  PT truststate types  ----------*/
	case 6: // follow
		return 1
	case 7: // block
		return 2
	/**----------  Naming truststate types  ----------*/
	case 8: // name assign
		return 1
	/**----------  F451 truststate types  ----------*/
	case 9: // censor assign
		return 1
	/*----------  PE truststate types  ----------*/
	case 10: // elect
		return 1
	case 11: // disqualify
		return 2
	default:
		return 0
	}
}

/*----------  Core functions for creating inflight objects  ----------*/
// The create / update state is passed through because the mint will make the decision on what fields to use or not, not this conversion.

/*----------  Content events  ----------*/

func createInflightBoard(i *feapi.ContentEventPayload) InflightBoard {
	return InflightBoard{
		Status: NewInflightStatus(STATUS_WAITING, i.GetEvent().GetEventType().String()),
		Entity: beObj.Board{
			/*----------  Identity fields  ----------*/
			Provable: &beObj.Provable{
				Fingerprint: i.GetEvent().GetPriorFingerprint(),
				Creation:    time.Now().Unix(),
			},
			Owner: i.GetEvent().GetOwnerFingerprint(),
			/*----------  Data fields  ----------*/
			Name:        i.GetBoardData().GetName(),
			Description: i.GetBoardData().GetDescription(),
			Meta:        i.GetBoardData().GetMeta(),
			// TODO FUTURE: Add board mods here after adding the UI for it.
		},
	}
}

func createInflightThread(i *feapi.ContentEventPayload) InflightThread {
	return InflightThread{
		Status: NewInflightStatus(STATUS_WAITING, i.GetEvent().GetEventType().String()),
		Entity: beObj.Thread{
			/*----------  Identity fields  ----------*/
			Provable: &beObj.Provable{
				Fingerprint: i.GetEvent().GetPriorFingerprint(),
				Creation:    time.Now().Unix(),
			},
			Owner: i.GetEvent().GetOwnerFingerprint(),
			/*----------  Data fields  ----------*/
			Board: i.GetThreadData().GetBoard(),
			Name:  i.GetThreadData().GetName(),
			Body:  i.GetThreadData().GetBody(),
			Link:  i.GetThreadData().GetLink(),
			Meta:  i.GetThreadData().GetMeta(),
		},
	}
}

func createInflightPost(i *feapi.ContentEventPayload) InflightPost {
	return InflightPost{
		Status: NewInflightStatus(STATUS_WAITING, i.GetEvent().GetEventType().String()),
		Entity: beObj.Post{
			/*----------  Identity fields  ----------*/
			Provable: &beObj.Provable{
				Fingerprint: i.GetEvent().GetPriorFingerprint(),
				Creation:    time.Now().Unix(),
			},
			Owner: i.GetEvent().GetOwnerFingerprint(),
			/*----------  Data fields  ----------*/
			Board:  i.GetPostData().GetBoard(),
			Thread: i.GetPostData().GetThread(),
			Parent: i.GetPostData().GetParent(),
			Body:   i.GetPostData().GetBody(),
			Meta:   i.GetPostData().GetMeta(),
		},
	}
}

func createInflightKey(i *feapi.ContentEventPayload) InflightKey {
	return InflightKey{
		Status: NewInflightStatus(STATUS_WAITING, i.GetEvent().GetEventType().String()),
		Entity: beObj.Key{
			/*----------  Identity fields  ----------*/
			Provable: &beObj.Provable{
				Fingerprint: i.GetEvent().GetPriorFingerprint(),
				Creation:    time.Now().Unix(),
			},
			/*----------  Data fields  ----------*/
			Type:   i.GetKeyData().GetType(),
			Key:    i.GetKeyData().GetKey(),
			Expiry: i.GetKeyData().GetExpiry(),
			Name:   i.GetKeyData().GetName(),
			Info:   i.GetKeyData().GetInfo(),
			Meta:   i.GetKeyData().GetMeta(),
			// Being able to pass through here doesn't mean minter is going to let changing key or key type though - both of those are immutable.
		},
	}
}

/*----------  Signal events  ----------*/

func createInflightVote(i *feapi.SignalEventPayload) InflightVote {
	return InflightVote{
		Status: NewInflightStatus(STATUS_WAITING, i.GetEvent().GetEventType().String()),
		Entity: beObj.Vote{
			/*----------  Identity fields  ----------*/
			Provable: &beObj.Provable{
				Fingerprint: i.GetEvent().GetPriorFingerprint(),
				Creation:    time.Now().Unix(),
			},
			Owner: i.GetEvent().GetOwnerFingerprint(),
			/*----------  Data fields  ----------*/
			Board:     i.GetTargetBoard(),
			Thread:    i.GetTargetThread(),
			Target:    i.GetTargetFingerprint(),
			TypeClass: parseTypeClass(i.GetSignalTypeClass()),
			Type:      parseType(i.GetSignalType()),
			Meta:      i.GetMeta(),
		},
	}
}

func createInflightTruststate(i *feapi.SignalEventPayload) InflightTruststate {
	return InflightTruststate{
		Status: NewInflightStatus(STATUS_WAITING, i.GetEvent().GetEventType().String()),
		Entity: beObj.Truststate{
			/*----------  Identity fields  ----------*/
			Provable: &beObj.Provable{
				Fingerprint: i.GetEvent().GetPriorFingerprint(),
				Creation:    time.Now().Unix(),
			},
			Owner: i.GetEvent().GetOwnerFingerprint(),
			/*----------  Data fields  ----------*/
			Target:    i.GetTargetFingerprint(),
			Domain:    i.GetTargetDomain(),
			Expiry:    i.GetTargetExpiry(),
			TypeClass: parseTypeClass(i.GetSignalTypeClass()),
			Type:      parseType(i.GetSignalType()),
			Meta:      i.GetMeta(),
		},
	}
}

/*----------  Push inflight changes to client  ----------*/

func (o *inflights) PushChangesToClient() {
	as := clapi.AmbientStatusPayload{Inflights: o.Protobuf()}
	clapiconsumer.SendAmbientStatus(&as)
}

/*----------  Prune the completed and failed items  ----------*/

func (o *inflights) Prune() {
	if o.ingestRunning {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()

	newInflightBoards := []InflightBoard{}
	for k, _ := range o.InflightBoards {
		if o.InflightBoards[k].Status.Fulfilled() {
			o.FulfilledBoards = append(o.FulfilledBoards, o.InflightBoards[k])
			continue
		}
		newInflightBoards = append(newInflightBoards, o.InflightBoards[k])
	}
	o.InflightBoards = newInflightBoards

	newInflightThreads := []InflightThread{}
	for k, _ := range o.InflightThreads {
		if o.InflightThreads[k].Status.Fulfilled() {
			o.FulfilledThreads = append(o.FulfilledThreads, o.InflightThreads[k])
			continue
		}
		newInflightThreads = append(newInflightThreads, o.InflightThreads[k])
	}
	o.InflightThreads = newInflightThreads

	newInflightPosts := []InflightPost{}
	for k, _ := range o.InflightPosts {
		if o.InflightPosts[k].Status.Fulfilled() {
			o.FulfilledPosts = append(o.FulfilledPosts, o.InflightPosts[k])
			continue
		}
		newInflightPosts = append(newInflightPosts, o.InflightPosts[k])
	}
	o.InflightPosts = newInflightPosts

	newInflightVotes := []InflightVote{}
	for k, _ := range o.InflightVotes {
		if o.InflightVotes[k].Status.Fulfilled() {
			o.FulfilledVotes = append(o.FulfilledVotes, o.InflightVotes[k])
			continue
		}
		newInflightVotes = append(newInflightVotes, o.InflightVotes[k])
	}
	o.InflightVotes = newInflightVotes

	newInflightKeys := []InflightKey{}
	for k, _ := range o.InflightKeys {
		if o.InflightKeys[k].Status.Fulfilled() {
			o.FulfilledKeys = append(o.FulfilledKeys, o.InflightKeys[k])
			continue
		}
		newInflightKeys = append(newInflightKeys, o.InflightKeys[k])
	}
	o.InflightKeys = newInflightKeys

	newInflightTruststates := []InflightTruststate{}
	for k, _ := range o.InflightTruststates {
		if o.InflightTruststates[k].Status.Fulfilled() {
			o.FulfilledTruststates = append(o.FulfilledTruststates, o.InflightTruststates[k])
			continue
		}
		newInflightTruststates = append(newInflightTruststates, o.InflightTruststates[k])
	}
	o.InflightTruststates = newInflightTruststates

	o.commit()
}
