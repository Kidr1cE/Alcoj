package analysis

import "alcoj/pkg/docker"

type GoAnalysis struct{}

func (*GoAnalysis) Analyze(cli *docker.DockerClient, path string) ([]LinterMessage, error) {
	return nil, nil
}
