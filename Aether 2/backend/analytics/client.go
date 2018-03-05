// Backend > Metrics

// This package is the metrics service defined and used by the metrics package. Metrics is used in non-release (alpha, beta, etc.) versions to observe network behaviour. It does not collect any information regarding the user, only about the backend that the user is using and how it is behaving in the network.

package metrics

import (
	pb "aether-core/backend/metrics/proto"
	"aether-core/services/globals"
	// "aether-core/services/logging"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// "log"
	"time"
)

const (
	addr = "127.0.0.1:19999"
)

func Prep() (pb.MetricsClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect: %v", err)
	}
	// defer conn.Close()
	c := pb.NewMetricsClient(conn)
	return c, conn
}

func IntroYourself(client pb.MetricsClient) *pb.Ack {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := client.SendIntro(
		ctx,
		&pb.Intro{
			NodeId:    globals.BackendConfig.GetNodeId(),
			LocalIp:   globals.BackendConfig.GetExternalIp(),
			LocalPort: fmt.Sprint(globals.BackendConfig.GetExternalPort())})
	if err != nil {
		fmt.Println("Could not Intro: %v", err)
	}
	return r
}
