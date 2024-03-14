package worker

import (
	"alcoj/pkg/analysis"
	"alcoj/pkg/docker"
	pb "alcoj/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	dockerClient *docker.DockerClient
	timeBash     = []string{"/usr/bin/time", "-v"}
)

type WorkerServer struct {
	cli           *docker.DockerClient
	entryShell    []string
	filenameIndex int
	pb.UnimplementedSandboxServer
}

func Run(stopChan chan struct{}) {
	startServer()
	stopChan <- struct{}{}
}

func startServer() {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	id := uuid.New().String()
	dockerClient, err = docker.NewDocker(id)
	if err != nil {
		return
	}

	s := grpc.NewServer()
	pb.RegisterSandboxServer(s, &WorkerServer{
		cli: dockerClient,
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *WorkerServer) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	log.Println("health check")
	return &pb.HealthCheckResponse{Status: true}, nil
}

func (s *WorkerServer) GetDockerStatus(ctx context.Context, in *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	msg, err := s.cli.Info(ctx)
	if err != nil {
		return &pb.GetStatusResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}
	return &pb.GetStatusResponse{
		Status:  true,
		Message: msg,
	}, nil
}

// raw:true -> entryshell raw:false -> dockerfile
func (s *WorkerServer) SetEnv(ctx context.Context, in *pb.SetEnvRequest) (*pb.SetEnvResponse, error) {
	s.cli.Image = in.ImageName
	s.entryShell = append(timeBash, strings.Split(in.Entryshell, " ")...)
	s.entryShell = append(s.entryShell, "nothing")
	s.filenameIndex = len(s.entryShell) - 1

	// Create container
	uuid := uuid.New().String()
	err := s.cli.Create(ctx, uuid)
	if err != nil {
		log.Println("create error: ", err)
		return &pb.SetEnvResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	err = s.cli.Start(ctx)
	if err != nil {
		log.Println("start error: ", err)
		return &pb.SetEnvResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SetEnvResponse{
		Status:  true,
		Message: "",
	}, nil
}

func (s *WorkerServer) SimpleRun(ctx context.Context, in *pb.SimpleRunRequest) (*pb.SimpleRunResponse, error) {
	s.entryShell[s.filenameIndex] = in.Filename
	if err := s.runAndTime(ctx, in.Filename); err != nil {
		return &pb.SimpleRunResponse{}, err
	}
	if err := s.runAndPylint(ctx, in.Filename); err != nil {
		return &pb.SimpleRunResponse{}, err
	}

	return &pb.SimpleRunResponse{
		TestResults: []*pb.TestResult{
			{
				Output: "output",
				Error:  "",
			},
		},
	}, nil
}

func (s *WorkerServer) runAndTime(ctx context.Context, filename string) error {
	s.entryShell[s.filenameIndex] = filename
	output, err := s.cli.Cmd(ctx, s.entryShell)
	if err != nil {
		log.Printf("Cmd() failed: %v", err)
		return err
	}
	log.Println("output: ", output)

	lines := strings.Split(output, "\n")

	timeMessage := analysis.TimeMessage{}
	commandOutputs := lines[0 : len(lines)-24]
	timeOutputs := lines[len(lines)-24:]
	for _, line := range timeOutputs {
		analysis.ParseTimeLine(line, &timeMessage)
	}

	log.Println("timeMessage: ", timeMessage)
	log.Println("commandOutputs:", commandOutputs)
	for _, output := range commandOutputs {
		log.Println(output)
	}

	return nil
}

func (s *WorkerServer) runAndPylint(ctx context.Context, filename string) error {
	output, err := s.cli.Cmd(ctx, []string{"pylint", filename})
	if err != nil {
		log.Printf("Cmd() failed: %v", err)
		return err
	}
	log.Println("output: ", output)
	pylintOutputs := analysis.ParsePylintOutput(output)
	for _, message := range pylintOutputs {
		log.Println(message.Column, message.ErrorCode, message.LineNumber, message.Message)
	}
	return nil
}
