package master

type Master struct {
	workerNum int
	workers   []*SandboxServer
}

var master *Master

func StartServer() {
	master = &Master{}
	master.AddSandboxServer("host.docker.internal:50051", ".py")

	startHttpServer()
}

func (m *Master) AddSandboxServer(address string, suffix string) error {
	s, err := newSandboxServer(address, suffix)
	if err != nil {
		return err
	}
	m.workers = append(m.workers, s)
	go s.Start()

	m.workerNum++
	return nil
}

func (m *Master) Run(language, code, input string) (Response, error) {
	req := Request{
		Code:     code,
		Language: language,
		Input:    input,
	}
	worker := m.getWorker()
	worker.reqCh <- req
	return <-worker.resCh, nil
}

func (m *Master) getWorker() *SandboxServer {
	return m.workers[0]
}
