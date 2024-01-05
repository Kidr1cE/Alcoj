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
	// uploadFile(c, "GZ", "/home/alco/golang-project/golang-oj-worker/cmd/factory/shifu_cloud_frontend.tar.gz")

	dockerInfo, err := c.GetDockerStatus(context.Background(), &pb.GetStatusRequest{})
	if err != nil {
		log.Printf("could not GetDockerStatus: %v", err)
		return
	}
	log.Println(dockerInfo)

	setenvResp, err := c.SetEnv(context.Background(), &pb.SetEnvRequest{
		Raw:        true,
		ImageName:  "python:3.6",
		Dockerfile: []byte("d1as32d13as21d3a1sd2as"),
		Entryshell: []byte("python /app/source/main.py"),
	})
	if err != nil {
		log.Printf("could not SetEnv: %v", err)
		return
	}
	log.Println(setenvResp)
}
