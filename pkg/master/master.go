package master

import (
	"context"
	"log"
	"os"

	"alcoj/pkg/analysis"
	pb "alcoj/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Master struct {
	suffix string
	conn   *grpc.ClientConn
	pbCli  pb.SandboxClient
}

var (
	master *Master
)

func StartServer() {
	var err error
	master, err = NewMaster("127.0.0.1:50051")
	if err != nil {
		log.Printf("could not create master: %v", err)
		return
	}
	defer master.Close()

	// HealthCheck
	req := &pb.HealthCheckRequest{}
	health, err := master.pbCli.HealthCheck(context.Background(), req)
	if err != nil {
		log.Printf("could not greet: %v", err)
		return
	}
	log.Println(health)

	// Set env
	setenvResp, err := master.pbCli.SetEnv(context.Background(), &pb.SetEnvRequest{
		ImageName:  "worker-python:v0.0.1",
		Entryshell: "python",
		Language:   "python",
	})
	if err != nil {
		log.Printf("could not SetEnv: %v", err)
		return
	}
	log.Println(setenvResp)

	// Set env
	// setenvResp, err := master.pbCli.SetEnv(context.Background(), &pb.SetEnvRequest{
	// 	ImageName:  "sandbox-golang:v0.0.1",
	// 	Entryshell: "go run",
	// 	Language:   "golang",
	// })
	// if err != nil {
	// 	log.Printf("could not SetEnv: %v", err)
	// 	return
	// }

	// Run
	// res, err := c.SimpleRun(context.Background(), &pb.SimpleRunRequest{
	// 	Filename: "main.py",
	// 	Input:    "12\n23\n",
	// })
	// if err != nil {
	// 	log.Printf("could not SimpleRun: %v", err)
	// 	return
	// }
	// log.Println(res)
	startHttpServer()
}

func NewMaster(address string) (*Master, error) {
	// Connect to pb server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}

	c := pb.NewSandboxClient(conn)
	master = &Master{
		conn:  conn,
		pbCli: c,
	}
	return master, nil
}

func (m *Master) Close() {
	m.conn.Close()
}

func (m *Master) Run(language string, code string, input string) (Response, error) {
	// Write code to file
	filename := uuid.New().String() + m.suffix
	err := os.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		return Response{}, err
	}

	// Run
	res, err := m.pbCli.SimpleRun(context.Background(), &pb.SimpleRunRequest{
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
