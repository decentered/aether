// This is a naive implementation of the metrics server for swarm testing purposes.

package simplemetricsserver

import (
	pb "aether-core/backend/metrics/proto"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

var Buf map[int64][]pb.Metrics // There can be multiple metrics pages arriving in the same UNIX timestamp, hence the []slice.

type server struct{}

func (s *server) RequestMetricsToken(ctx context.Context, remoteMachine *pb.Machine) (*pb.Machine_MetricsToken, error) {
	fmt.Printf("A message was received from the node w/ port: %d. It's requesting a metrics token.\n", remoteMachine.GetPort())
	// saveNode(remoteMachine.Client.GetName(), int(remoteMachine.GetPort()))
	metricsToken := pb.Machine_MetricsToken{Token: "testtoken"}
	return &metricsToken, nil
}

func (s *server) UploadMetrics(ctx context.Context, metrics *pb.Metrics) (*pb.MetricsDeliveryResponse, error) {
	fmt.Printf("A message was received from the node w/ port: %d. It's sending metrics.\n", metrics.GetMachine().GetPort())
	// This saves inbound metrics into a file, so that we will have a record of what swarm nodes are doing in the network.
	now := time.Now().Unix()
	Buf[now] = append(Buf[now], *metrics)
	// spew.Dump(metrics)
	return &pb.MetricsDeliveryResponse{}, nil
}

func StartListening() {
	Buf = make(map[int64][]pb.Metrics)
	fmt.Println("Metrics server started listening.")
	lis, err := net.Listen("tcp", "127.0.0.1:19999")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMetricsServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	fmt.Println("Metrics server stopped listening.")

}
