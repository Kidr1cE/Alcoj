package factory

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

	req := &pb.HealthCheckRequest{}
	health, err := c.HealthCheck(context.Background(), req)
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}
	log.Println(health)

	// Set env
	setenvResp, err := c.SetEnv(context.Background(), &pb.SetEnvRequest{
		Raw:        true,
		ImageName:  "python:3.6",
		Entryshell: []byte("python /sandbox/main.py"),
	})
	if err != nil {
		log.Printf("could not SetEnv: %v", err)
		return
	}

	// Send Requirements
	uploadFile(c, "main.py", "/home/alco/go-project/Alcoj/cmd/factory/main.py")
	dockerInfo, err := c.GetDockerStatus(context.Background(), &pb.GetStatusRequest{})
	if err != nil {
		log.Printf("could not GetDockerStatus: %v", err)
		return
	}
	log.Println(dockerInfo)

	// Run
	res, err := c.SimpleRun(context.Background(), &pb.SimpleRunRequest{})
	if err != nil {
		log.Printf("could not SimpleRun: %v", err)
		return
	}
	log.Println(res)

	log.Println(setenvResp)
}
