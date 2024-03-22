package master

import (
	"context"
	"log"

	pb "alcoj/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartServer() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewSandboxClient(conn)

	// HealthCheck
	req := &pb.HealthCheckRequest{}
	health, err := c.HealthCheck(context.Background(), req)
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}
	log.Println(health)

	// Set env
	setenvResp, err := c.SetEnv(context.Background(), &pb.SetEnvRequest{
		ImageName:  "worker-python:v0.0.1",
		Entryshell: "python",
		Language:   "python",
	})
	if err != nil {
		log.Printf("could not SetEnv: %v", err)
		return
	}

	// Run
	res, err := c.SimpleRun(context.Background(), &pb.SimpleRunRequest{
		Filename: "main.py",
		Input:    "12\n23\n",
	})
	if err != nil {
		log.Printf("could not SimpleRun: %v", err)
		return
	}
	log.Println(res)

	log.Println(setenvResp)
}
