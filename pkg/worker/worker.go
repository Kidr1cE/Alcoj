package worker

import (
	"alcoj/pkg/analysis"
	"alcoj/pkg/docker"
	pb "alcoj/proto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var (
	dockerClient *docker.DockerClient
	timeBash     = []string{"/usr/bin/time", "-v"}
	port         = os.Getenv("PORT")
	masterAddr   = os.Getenv("MASTER_ADDR")
	id           = os.Getenv("WORKER_ID")
)

type WorkerServer struct {
	cli           *docker.DockerClient
	entryShell    []string
	filenameIndex int
	analysis      analysis.AnalysisInterface
	runner        analysis.Runner
	pb.UnimplementedSandboxServer
}

type RegisterRequest struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Suffix  string `json:"suffix"`
}

func StartServer() {
	stopCh := make(chan struct{})
	go func() {
		err := startServer(stopCh)
		if err != nil {
			log.Println("failed to start server: ", err)
			return
		}
	}()
	time.Sleep(1 * time.Second)
	err := register()
	if err != nil {
		log.Println("failed to register: ", err)
		return
	}
	stopCh <- struct{}{}
}

func register() error {
	req := RegisterRequest{
		ID:      id,
		Address: "host.docker.internal:" + port,
	}
	content, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Connect to master
	res, err := http.Post(
		"http://"+masterAddr+"/alcoj/api/v1/register",
		"application/json",
		bytes.NewBuffer(content),
	)
	if err != nil {
		return err
	}
	log.Println("register response: ", res.Status)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register: %s", res.Status)
	}
	return nil
}

func startServer(stopCh chan struct{}) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	dockerClient, err = docker.NewDocker(id + "-box")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	ws := &WorkerServer{
		cli: dockerClient,
	}
	pb.RegisterSandboxServer(s, ws)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	stopCh <- struct{}{}
	return nil
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

	switch in.Language {
	case "python":
		s.analysis = &analysis.PythonAnalysis{}
	case "golang":
		s.analysis = &analysis.GolangAnalysis{}
	default:
		return &pb.SetEnvResponse{
			Status:  false,
			Message: "unsupported language",
		}, nil
	}

	s.runner = analysis.Runner{}

	// Create container
	err := s.cli.Create(ctx, id+"-box")
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
	output, timeMessage, err := s.runner.RunTime(s.cli, in.Filename, s.entryShell, in.Input)
	if err != nil {
		return &pb.SimpleRunResponse{}, err
	}
	ts := &pb.TimeResult{
		SystemTimeSeconds:   timeMessage.SystemTimeSeconds,
		UserTimeSeconds:     timeMessage.UserTimeSeconds,
		PercentCpu:          timeMessage.PercentCPU,
		AvgSharedTextSize:   timeMessage.AvgSharedTextSize,
		AvgUnsharedDataSize: timeMessage.AvgUnsharedDataSize,
		MaxResidentSetSize:  timeMessage.MaxResidentSetSize,
		FileSystemInputs:    timeMessage.FileSystemInputs,
		FileSystemOutputs:   timeMessage.FileSystemOutputs,
		ExitStatus:          timeMessage.ExitStatus,
	}

	linterMessage, err := s.analysis.Analyze(s.cli, in.Filename)
	if err != nil {
		return &pb.SimpleRunResponse{}, err
	}
	ar := make([]*pb.AnalysisResult, len(linterMessage))
	for i, v := range linterMessage {
		ar[i] = &pb.AnalysisResult{
			Row:     int32(v.Row),
			Column:  int32(v.Column),
			Message: v.Message,
		}
	}

	return &pb.SimpleRunResponse{
		Output:          output,
		AnalysisResults: ar,
		TimeResult:      ts,
	}, nil
}
