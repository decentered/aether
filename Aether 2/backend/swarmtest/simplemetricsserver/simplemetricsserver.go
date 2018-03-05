// This is a naive implementation of the metrics server for swarm testing purposes.

package simplemetricsserver

import (
	pb "aether-core/backend/metrics/proto"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func (s *server) SendIntro(ctx context.Context, in *pb.Intro) (*pb.Ack, error) {
	fmt.Println("A message was received from one of the swarm nodes.")
	fmt.Println(in.GetNodeId())
	fmt.Println(in.GetLocalIp())
	fmt.Println(in.GetLocalPort())
	return &pb.Ack{Message: "Ack."}, nil
}

func StartListening() {
	fmt.Println("Metrics server started listening.")
	lis, err := net.Listen("tcp", "127.0.0.1:19999")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMetricsServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	fmt.Println("Metrics server stopped listening.")

}
