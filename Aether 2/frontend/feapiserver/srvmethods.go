// Frontend > FrontendAPI
// This package implements the FrontendAPI GRPC server.

package feapiserver

import (
	"aether-core/frontend/beapiconsumer"
	"aether-core/frontend/clapiconsumer"
	"aether-core/frontend/festructs"
	// "aether-core/frontend/objpool"
	"aether-core/frontend/inflights"
	// "aether-core/io/api"
	"aether-core/protos/clapi"
	pb "aether-core/protos/feapi"
	pbObj "aether-core/protos/feobjects"
	"aether-core/services/globals"
	"aether-core/services/logging"
	// "encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

type server struct{}

func (s *server) BackendReady(ctx context.Context, req *pb.BEReadyRequest) (*pb.BEReadyResponse, error) {
	fmt.Printf("Backend sent us a backendready request. Req: %#v Context: %#v\n", req, ctx)
	globals.FrontendConfig.SetBackendAPIAddress(req.GetAddress())
	globals.FrontendConfig.SetBackendAPIPort(int(req.GetPort()))
	globals.FrontendTransientConfig.BackendReady = true
	// Send the ack to the client that we are all ready.
	clapiconsumer.SendFrontendReady()
	resp := pb.BEReadyResponse{}
	return &resp, nil
}

func (s *server) GetThreadAndPosts(ctx context.Context, req *pb.ThreadAndPostsRequest) (*pb.ThreadAndPostsResponse, error) {
	// Get board
	fp := req.GetBoardFingerprint()
	// Get the board carrier
	bc := festructs.BoardCarrier{}
	err := globals.KvInstance.One("Fingerprint", fp, &bc)
	if err != nil {
		logging.Logf(1, "Getting BoardCarrier for in GetThreadAndPosts encountered an error. Error: %v", err)
	}
	b := festructs.CompiledBoard{}
	for key, _ := range bc.Boards {
		if bc.Boards[key].Fingerprint == fp {
			b = bc.Boards[key]
		}
	}
	resp := pb.ThreadAndPostsResponse{}
	resp.Board = b.Protobuf()
	// Check if board is subbed / notify enabled.
	subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(resp.Board.Fingerprint)
	resp.Board.Subscribed = subbed
	resp.Board.Notify = notify
	resp.Board.LastSeen = lastseen
	// Get thread
	tfp := req.GetThreadFingerprint()
	// Get the board carrier
	tc := festructs.ThreadCarrier{}
	err2 := globals.KvInstance.One("Fingerprint", tfp, &tc)
	if err2 != nil {
		logging.Logf(1, "Getting ThreadCarrier for in GetThreadAndPosts encountered an error. Error: %v", err2)
	}
	t := festructs.CompiledThread{}
	for key, _ := range tc.Threads {
		if tc.Threads[key].Fingerprint == tfp {
			t = tc.Threads[key]
		}
	}
	resp.Thread = t.Protobuf()
	// Get posts
	postprotos := []*pbObj.CompiledPostEntity{}
	for k, _ := range tc.Posts {
		postprotos = append(postprotos, tc.Posts[k].Protobuf())
	}
	resp.Posts = postprotos
	return &resp, nil
}

func (s *server) GetBoardAndThreads(ctx context.Context, req *pb.BoardAndThreadsRequest) (*pb.BoardAndThreadsResponse, error) {
	fp := req.GetBoardFingerprint()
	// Get the board carrier
	bc := festructs.BoardCarrier{}
	err := globals.KvInstance.One("Fingerprint", fp, &bc)
	if err != nil {
		logging.Logf(1, "Getting BoardCarrier for in GetBoardAndThreads encountered an error. Error: %v", err)
	}
	b := festructs.CompiledBoard{}
	for key, _ := range bc.Boards {
		if bc.Boards[key].Fingerprint == fp {
			b = bc.Boards[key]
		}
	}
	resp := pb.BoardAndThreadsResponse{}
	resp.Board = b.Protobuf()
	// Check if board is subbed / notify enabled.
	subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(resp.Board.Fingerprint)
	resp.Board.Subscribed = subbed
	resp.Board.Notify = notify
	resp.Board.LastSeen = lastseen
	// Get the threads that board contains
	thrcs := []festructs.ThreadCarrier{}
	err2 := globals.KvInstance.Find("ParentFingerprint", fp, &thrcs)
	if err2 != nil {
		logging.Logf(1, "Getting Threads for in GetBoardAndThreads encountered an error. Error: %v", err2)
	}
	threads := []festructs.CompiledThread{}
	for k1, _ := range thrcs {
		for k2, _ := range thrcs[k1].Threads {
			if thrcs[k1].Threads[k2].Board == fp {
				threads = append(threads, thrcs[k1].Threads[k2])
			}
		}
	}
	// Convert all threads to protos
	tprotos := []*pbObj.CompiledThreadEntity{}
	for k, _ := range threads {
		tprotos = append(tprotos, threads[k].Protobuf())
	}
	resp.Threads = tprotos
	return &resp, nil
}

func (s *server) GetAllBoards(ctx context.Context, req *pb.AllBoardsRequest) (*pb.AllBoardsResponse, error) {
	fmt.Println("We received a get all boards request.")
	var boards []festructs.BoardCarrier
	start := time.Now()
	err := globals.KvInstance.All(&boards)
	if err != nil {
		logging.Logcf(1, "Getting all boards from KvInstance encountered an error. Error: %v", err)
	}

	cb := []*pbObj.CompiledBoardEntity{}
	for key, _ := range boards {
		for k2, _ := range boards[key].Boards {
			item := boards[key].Boards[k2].Protobuf()
			subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(item.Fingerprint)
			whitelisted := globals.FrontendConfig.ContentRelations.IsWhitelistedBoard(item.Fingerprint)
			item.Subscribed = subbed
			item.Notify = notify
			item.LastSeen = lastseen
			item.Whitelisted = whitelisted
			cb = append(cb, item)
		}
	}
	fmt.Printf("Number of items found in get all boards: %v\n", len(boards))
	resp := pb.AllBoardsResponse{cb}
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	return &resp, nil
}

func (s *server) SetClientAPIServerPort(ctx context.Context, req *pb.SetClientAPIServerPortRequest) (*pb.SetClientAPIServerPortResponse, error) {
	logging.Logf(1, "We received a set client api server port request. Old port was: %v and the new one is %v", globals.FrontendConfig.GetClientPort(), req.Port)
	globals.FrontendConfig.SetClientPort(int(req.Port))
	clapiconsumer.DeliverAmbients()
	inflights := inflights.GetInflights()
	as := clapi.AmbientStatusPayload{Inflights: inflights.Protobuf()}
	clapiconsumer.SendAmbientStatus(&as)
	clapiconsumer.PushLocalUserAmbient() // todo let's disable this for a minute.
	// SendAmbients(false)
	resp := pb.SetClientAPIServerPortResponse{}
	return &resp, nil
}

func (s *server) SetBoardSignal(ctx context.Context, req *pb.BoardSignalRequest) (*pb.BoardSignalResponse, error) {
	// logging.Logf(1, "We've received a set board signal request.")
	cr := globals.FrontendConfig.GetContentRelations()
	committed := cr.SetBoardSignal(req.Fingerprint, req.Subscribed, req.Notify, req.LastSeen, req.LastSeenOnly)
	globals.FrontendConfig.SetContentRelations(cr)
	resp := pb.BoardSignalResponse{Committed: committed}
	clapiconsumer.DeliverAmbients()
	return &resp, nil
}

func (s *server) GetUserAndGraph(ctx context.Context, req *pb.UserAndGraphRequest) (*pb.UserAndGraphResponse, error) {
	fp := req.GetFingerprint()
	resp := pb.UserAndGraphResponse{}
	resp.UserEntityRequested = req.GetUserEntityRequested()
	resp.UserBoardsRequested = req.GetUserBoardsRequested()
	resp.UserThreadsRequested = req.GetUserThreadsRequested()
	resp.UserPostsRequested = req.GetUserPostsRequested()
	if req.GetUserEntityRequested() {
		uh := festructs.UserHeaderCarrier{}
		err := globals.KvInstance.One("Fingerprint", fp, &uh)
		if err != nil {
			logging.Logf(1, "Getting User Header Carrier for GetUserAndGraph encountered an error. Error: %v", err)
		}
		u := festructs.CompiledUser{}
		for key, _ := range uh.Users {
			if uh.Users[key].Fingerprint == fp {
				u = uh.Users[key]
			}
		}
		resp.User = u.Protobuf()
	}
	// if req.GetUserBoardsRequested() {
	// 	// todo
	// }
	// if req.GetUserThreadsRequested() {
	// 	// todo
	// }
	// if req.GetUserPostsRequested() {
	// 	// todo
	// }
	// ^ We actually made it so that these data are actually provided directly from the backend as uncompiled payloads.
	logging.Logf(1, "resp: %v", resp)
	return &resp, nil
}

func (s *server) SendContentEvent(ctx context.Context, req *pb.ContentEventPayload) (*pb.ContentEventResponse, error) {
	logging.Logf(1, "We've received a content event. Event: %v", *req)

	inflights := inflights.GetInflights()
	inflights.Insert(*req)
	// logging.Logf(1, "Pool: %#v", pool)
	as := clapi.AmbientStatusPayload{Inflights: inflights.Protobuf()}
	clapiconsumer.SendAmbientStatus(&as)
	// SendAmbients(true)
	resp := pb.ContentEventResponse{}
	return &resp, nil
}

func (s *server) SendSignalEvent(ctx context.Context, req *pb.SignalEventPayload) (*pb.SignalEventResponse, error) {
	logging.Logf(1, "We've received a signal event. Event: %v", *req)
	inflights := inflights.GetInflights()
	inflights.Insert(*req)
	// logging.Logf(1, "Pool: %#v", pool)
	as := clapi.AmbientStatusPayload{Inflights: inflights.Protobuf()}
	clapiconsumer.SendAmbientStatus(&as)
	// SendAmbients(true)
	resp := pb.SignalEventResponse{}
	return &resp, nil
}

func (s *server) GetUncompiledEntityByKey(ctx context.Context, req *pb.UncompiledEntityByKeyRequest) (*pb.UncompiledEntityByKeyResponse, error) {
	logging.Logf(1, "We've received an uncompiled entity by key request. Event: %v", *req)
	switch req.GetEntityType() {
	case pb.UncompiledEntityType_BOARD:
		entities := beapiconsumer.GetBoardsByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_BOARD,
			Boards:     entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_THREAD:
		entities := beapiconsumer.GetThreadsByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_THREAD,
			Threads:    entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_POST:
		entities := beapiconsumer.GetPostsByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_POST,
			Posts:      entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_VOTE:
		entities := beapiconsumer.GetVotesByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_VOTE,
			Votes:      entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_KEY:
		entities := beapiconsumer.GetKeysByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_KEY,
			Keys:       entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_TRUSTSTATE:
		entities := beapiconsumer.GetTruststatesByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType:  pb.UncompiledEntityType_TRUSTSTATE,
			Truststates: entities,
		}
		return &resp, nil
	default:
		return nil, errors.New("Entity type could not be determined")
	}
}

func (s *server) SendInflightsPruneRequest(ctx context.Context, req *pb.InflightsPruneRequest) (*pb.InflightsPruneResponse, error) {
	inflights := inflights.GetInflights()
	inflights.Prune()
	inflights.PushChangesToClient()
	resp := pb.InflightsPruneResponse{}
	return &resp, nil
}
