package master

import (
	"log"
	"os"
	"sync/atomic"
)

var language = os.Getenv("LANGUAGE")

type Master struct {
	WorkerNum     atomic.Int64     `json:"worker_num"`
	QueueTasks    atomic.Int64     `json:"queue_tasks"`
	FinishedTasks atomic.Int64     `json:"finished_tasks"`
	Workers       []*SandboxServer `json:"workers"`
}

var master *Master

func StartServer() {
	master = &Master{
		QueueTasks:    atomic.Int64{},
		FinishedTasks: atomic.Int64{},
	}
	master.QueueTasks.Store(0)
	master.FinishedTasks.Store(0)
	// master.AddSandboxServer("host.docker.internal:50051", ".py")

	httpStopCh := make(chan struct{})
	websocketStopCh := make(chan struct{})
	go startHttpServer(httpStopCh)
	go startWebsocketServer(websocketStopCh)

	<-httpStopCh
	<-websocketStopCh
}

func (m *Master) AddSandboxServer(id string, address string) error {
	suffix := ""
	switch language {
	case "python":
		suffix = ".py"
	case "golang":
		suffix = ".go"
	}

	sandbox, err := newSandboxServer(id, address, suffix)
	if err != nil {
		return err
	}
	m.Workers = append(m.Workers, sandbox)
	go func() {
		err := sandbox.Start()
		if err != nil {
			log.Println("failed to start sandbox server: ", err)
			return
		}
	}()

	m.WorkerNum.Add(1)
	return nil
}

func (m *Master) Run(language, code, input string) (JudgeResponse, error) {
	defer func() {
		m.QueueTasks.Add(-1)
		m.FinishedTasks.Add(1)
	}()
	req := JudgeRequest{
		Code:     code,
		Language: language,
		Input:    input,
	}
	m.QueueTasks.Add(1)

	for _, worker := range m.Workers {
		select {
		case worker.reqCh <- req:

		default:
			continue
		}
		return <-worker.resCh, nil
	}

	return JudgeResponse{}, nil
}
