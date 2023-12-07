package core

import (
	pb "alcoj/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type WorkerServer struct {
	pb.UnimplementedSandboxServer
}

func Run(stopChan chan bool) {

	startServer()
}

func startServer() {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSandboxServer(s, &WorkerServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *WorkerServer) GetEnv(ctx context.Context, req *pb.SetEnvRequest) (*pb.SetEnvResponse, error) {

	return nil, nil
}

func (s *WorkerServer) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {

	return nil, nil
}

func (s *WorkerServer) Run(ctx context.Context, req *pb.RunRequest) (*pb.RunResponse, error) {

	return nil, nil
}
