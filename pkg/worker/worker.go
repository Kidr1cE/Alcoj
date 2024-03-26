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
	analysis      analysis.AnalysisInterface
	runner        analysis.Runner
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
	ws := &WorkerServer{
		cli: dockerClient,
	}
	pb.RegisterSandboxServer(s, ws)

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
