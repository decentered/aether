// Frontend > ClientAPI Client
// This package is the client side of the Client's GRPC API. This is the API frontend uses to send the client some frontend health related information, updates, etc.

package clapiconsumer

import (
	"aether-core/frontend/festructs"
	"aether-core/io/api"
	pb "aether-core/protos/clapi"
	"aether-core/protos/feobjects"
	"aether-core/services/globals"
	"aether-core/services/logging"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func StartClientAPIConnection() (pb.ClientAPIClient, *grpc.ClientConn) {
	clAddr := fmt.Sprint(globals.FrontendConfig.GetClientAPIAddress(), ":", globals.FrontendConfig.GetClientPort())
	conn, err := grpc.Dial(clAddr, grpc.WithInsecure())
	if err != nil {
		logging.Logf(1, "Could not connect to the client API service. Error: %v", err)
	}
	c := pb.NewClientAPIClient(conn)
	return c, conn
}

func SendFrontendReady() {
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := pb.FEReadyRequest{
		Address: "127.0.0.1",
		Port:    int32(globals.FrontendConfig.GetFrontendAPIPort()),
	}
	_, err := c.FrontendReady(ctx, &payload)
	if err != nil {
		logging.Logf(1, "SendFrontendReady encountered an error. Err: %v", err)
	}
}

func DeliverAmbients() {
	logging.Logf(1, "Deliver ambients is called, FE>CL, Cl receiver port is: %v", globals.FrontendConfig.GetClientPort())
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := pb.AmbientsRequest{
		Boards: festructs.GetCurrentAmbients().Protobuf(),
	}
	_, err := c.DeliverAmbients(ctx, &payload)
	if err != nil {
		logging.Logf(1, "DeliverAmbients encountered an error. Err: %v", err)
	}
}

func SendAmbientStatus(cas *pb.AmbientStatusPayload) {
	logging.Logf(1, "SendAmbientStatus is called")
	if cas != nil {
		updateAmbientStatus(cas)
		// If it is nil we just use the extant ambient status in fe transient config
	}
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := globals.FrontendTransientConfig.CurrentAmbientStatus
	_, err := c.SendAmbientStatus(ctx, &payload)
	if err != nil {
		logging.Logf(1, "SendAmbientStatus encountered an error. Err: %v", err)
	}
}

// updateAmbientStatus partially updates the parts of the live ambient status. So effectively if you make an update to the inflights, this one makes it so that the update doesn't delete the existing but older ambient statuses from backend and frontend.
func updateAmbientStatus(currentAmbientStatus *pb.AmbientStatusPayload) {
	as := pb.AmbientStatusPayload{}
	if bas := currentAmbientStatus.GetBackendAmbientStatus(); bas != nil {
		as.BackendAmbientStatus = bas
	}
	if fas := currentAmbientStatus.GetFrontendAmbientStatus(); fas != nil {
		as.FrontendAmbientStatus = fas
	}
	if ifl := currentAmbientStatus.GetInflights(); ifl != nil {
		as.Inflights = ifl
	}
	globals.FrontendTransientConfig.CurrentAmbientStatus = as
	logging.Logf(1, "Current ambient status: %v", as)
}

/*----------  Ambient Local User Data  ----------*/

func SendAmbientLocalUserEntity(localUserExists bool, localUser *feobjects.CompiledUserEntity) {
	logging.Logf(1, "AmbientLocalUserEntity is called")

	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	alu := pb.AmbientLocalUserEntityPayload{}
	alu.LocalUserExists = localUserExists
	alu.LocalUserEntity = localUser
	logging.Logf(1, "alu: %#v", alu)
	_, err := c.SendAmbientLocalUserEntity(ctx, &alu)
	if err != nil {
		logging.Logf(1, "AmbientLocalUserEntity encountered an error. Err: %v", err)
	}
}

/*----------  Higher level methods  ----------*/
/*
These methods aren't 1-1 matches to the gRPC API.
*/

// pushLocalUserAmbient reads from the configstore, and if local user doesn't exist there, bails. If it does, it attempts to read the compiled user header with the same fingerprint. If that fails, the entity exists but not found, and no data is sent until the next attempt.
func PushLocalUserAmbient() {
	alu := globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()
	localUserExists := false
	var fp string
	if len(alu) == 0 {
		SendAmbientLocalUserEntity(false, nil)
		return
	}
	localUserExists = true
	var key api.Key
	json.Unmarshal([]byte(alu), &key)
	fp = string(key.Fingerprint)
	uh := festructs.UserHeaderCarrier{}
	err := globals.KvInstance.One("Fingerprint", fp, &uh)
	if err != nil {
		logging.Logf(1, "Getting the compiled user entity in PushLocalUserAmbient failed. Error: %v", err)
		// If it exists but not found in the compiled store, that means it hasn't been compiled yet. In this case, we wait and not push anything so that the client can keep its 'loading' state.
		return
	}
	u := festructs.CompiledUser{}
	for key, _ := range uh.Users {
		if uh.Users[key].Fingerprint == fp {
			u = uh.Users[key]
		}
	}
	uproto := u.Protobuf()
	SendAmbientLocalUserEntity(localUserExists, uproto)
	return
}
