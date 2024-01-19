package worker

import (
	"alcoj/pkg/docker"
	"alcoj/pkg/util"
	pb "alcoj/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var dockerClient *docker.DockerClient

var sandBoxWorkDir = "/app/source"

type WorkerServer struct {
	cli *docker.DockerClient
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
	dockerClient, err = docker.NewWorker()
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
	log.Println("set env")

	// if dockerfile is not raw, write it to dockerfile
	if !in.Raw {
		if err := util.Write(in.Dockerfile, docker.DockerfilePath); err != nil {
			return &pb.SetEnvResponse{
				Status:  false,
				Message: err.Error(),
			}, nil
		}
		s.cli.Image = in.ImageName
		if err := s.cli.Build(ctx); err != nil {
			log.Println("build error: ", err)
			return &pb.SetEnvResponse{
				Status:  false,
				Message: err.Error(),
			}, nil
		}
	} else {
		s.cli.Image = in.ImageName
		if err := util.Write(in.Entryshell, "/app/source/run.sh"); err != nil {
			return &pb.SetEnvResponse{
				Status:  false,
				Message: err.Error(),
			}, nil
		}
	}

	// Create container
	err := s.cli.Create(ctx)
	if err != nil {
		log.Println("create error: ", err)
		return &pb.SetEnvResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}
	err = s.cli.Start(ctx, "")
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

// send requirements to /sandbox
func (s *WorkerServer) SendRequirements(stream pb.Sandbox_SendRequirementsServer) error {
	log.Println("send requirements")
	var filename string
	var content = []byte{}
	var err error
	for {
		// get chunk
		chunk, err := stream.Recv()
		if err != nil {
			log.Println("recv error: ", err)
			stream.Send(&pb.UploadStatus{
				Success: false,
				Message: err.Error(),
			})
			break
		}

		// check if it is a new file
		if chunk.Filename != filename {
			content = []byte{}
			filename = chunk.Filename
			log.Println("new file: ", filename)
		} else {
			log.Println("same file: ", filename)
		}

		content = append(content, chunk.Content...)

		if err := stream.Send(&pb.UploadStatus{Success: true}); err != nil {
			log.Println("send error: ", err)
			stream.Send(&pb.UploadStatus{
				Success: false,
				Message: err.Error(),
			})
		}

		if chunk.IsLastChunk {
			log.Println("last chunk")
			if err = util.WriteToAppFolder(filename, content); err != nil {
				stream.Send(&pb.UploadStatus{
					Success: false,
					Message: err.Error(),
				})
				log.Println("write to app folder error: ", err)
				break
			}
			stream.Send(&pb.UploadStatus{Success: true})
			return nil
		}
	}
	return err
}

func (s *WorkerServer) SimpleRun(ctx context.Context, in *pb.SimpleRunRequest) (*pb.SimpleRunResponse, error) {
	log.Println("simple run")
	s.cli.Start(ctx, "")
	result, err := s.cli.Run(ctx)
	if err != nil {
		log.Println("run error: ", err)
		return nil, err
	}
	return &pb.SimpleRunResponse{
		TestResults: []*pb.TestResult{
			{
				Output: result,
				Error:  "",
			},
		},
	}, nil
}
