package analysis

import "alcoj/pkg/docker"

type LinterMessage struct {
	Row       int    `json:"row,omitempty"`        // 6
	Column    int    `json:"column,omitempty"`     // 18
	ErrorCode string `json:"error_code,omitempty"` // C0103
	Message   string `json:"message,omitempty"`    // Argument name "columnPosition" doesn't conform to snake_case naming style
}

type AnalysisInterface interface {
	Analyze(cli *docker.DockerClient, path string) ([]LinterMessage, error)
}
