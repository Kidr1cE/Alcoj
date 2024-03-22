package analysis

import "alcoj/pkg/docker"

type LinterMessage struct {
	Row       int    // 6
	Column    int    // 18
	ErrorCode string // C0103
	Message   string // Argument name "columnPosition" doesn't conform to snake_case naming style
}

type AnalysisInterface interface {
	Analyze(cli *docker.DockerClient, path string) ([]LinterMessage, error)
}
