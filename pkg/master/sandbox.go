package master

import (
	"alcoj/pkg/analysis"
	pb "alcoj/proto"
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SandboxServer struct {
	conn   *grpc.ClientConn
	pbCli  pb.SandboxClient
	suffix string
	stopCh chan struct{}
	reqCh  chan Request
	resCh  chan Response
}

func newSandboxServer(address string, suffix string) (*SandboxServer, error) {
	// Connect to pb server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}

	c := pb.NewSandboxClient(conn)
	server := &SandboxServer{
		conn:   conn,
		pbCli:  c,
		suffix: suffix,
		stopCh: make(chan struct{}),
		reqCh:  make(chan Request),
		resCh:  make(chan Response),
	}

	return server, nil
}

func (s *SandboxServer) Close() {
	s.stopCh <- struct{}{}
	s.conn.Close()
}

func (s *SandboxServer) Start() error {
	health, err := s.pbCli.HealthCheck(context.Background(), &pb.HealthCheckRequest{})
	if err != nil {
		return err
	}
	log.Println(health)

	setenvResp, err := s.pbCli.SetEnv(context.Background(), &pb.SetEnvRequest{
		ImageName:  "worker-python:v0.0.1",
		Entryshell: "python",
		Language:   "python",
	})
	if err != nil {
		return err
	}
	log.Println(setenvResp)

	for {
		select {
		case <-s.stopCh:
			return nil
		case req := <-s.reqCh:
			res, err := s.run(req.Language, req.Code, req.Input)
			if err != nil {
				return err
			}
			s.resCh <- res
		}
	}
}

func (s *SandboxServer) run(language string, code string, input string) (Response, error) {
	// Write code to file
	filename := uuid.New().String() + s.suffix
	err := os.WriteFile("/sandbox/"+filename, []byte(code), 0644)
	if err != nil {
		return Response{}, err
	}

	// Run
	res, err := s.pbCli.SimpleRun(context.Background(), &pb.SimpleRunRequest{
		Filename: filename,
		Input:    input,
	})
	if err != nil {
		log.Printf("could not SimpleRun: %v", err)
		return Response{}, err
	}
	message, err := parseResponse(res)
	if err != nil {
		return Response{}, err
	}
	return message, nil
}

func parseResponse(response *pb.SimpleRunResponse) (Response, error) {
	// Parse
	var res Response
	res.Output = response.Output
	if response.TimeResult != nil {
		timeResult := analysis.TimeMessage{
			SystemTimeSeconds:   response.TimeResult.SystemTimeSeconds,
			UserTimeSeconds:     response.TimeResult.UserTimeSeconds,
			PercentCPU:          response.TimeResult.PercentCpu,
			AvgSharedTextSize:   response.TimeResult.AvgSharedTextSize,
			AvgUnsharedDataSize: response.TimeResult.AvgUnsharedDataSize,
			MaxResidentSetSize:  response.TimeResult.MaxResidentSetSize,
			FileSystemInputs:    response.TimeResult.FileSystemInputs,
			FileSystemOutputs:   response.TimeResult.FileSystemOutputs,
			ExitStatus:          response.TimeResult.ExitStatus,
		}
		res.RuntimeAnalysis = timeResult
	}
	if len(response.AnalysisResults) > 0 {
		analysisResult := make([]analysis.LinterMessage, len(response.AnalysisResults))
		for i := range response.AnalysisResults {
			analysisResult[i] = analysis.LinterMessage{
				Row:     int(response.AnalysisResults[i].Row),
				Column:  int(response.AnalysisResults[i].Column),
				Message: response.AnalysisResults[i].Message,
			}
		}
		res.StaticAnalysis = analysisResult
	}
	return res, nil
}
