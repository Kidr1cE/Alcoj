package analysis

import (
	"alcoj/pkg/docker"
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type PythonAnalysis struct{}

func parsePylintOutput(output string) []LinterMessage {
	messages := []LinterMessage{}
	reg := regexp.MustCompile(`(\d+):(\d+):\s+(\w+)\:\s+(.+)`)
	matches := reg.FindAllString(output, -1)

	for _, match := range matches {
		message := parsePylintLine(match)
		messages = append(messages, message)
	}
	return messages
}

// parse line like "27:0: W0311: Bad indentation. Found 1 spaces, expected 4 (bad-indentation)"
func parsePylintLine(line string) LinterMessage {
	lineParts := strings.Split(line, ":")
	lineNumber, err := strconv.Atoi(lineParts[0])
	if err != nil {
		log.Fatalf("Failed to parse line number: %v", err)
	}
	columnNumber, err := strconv.Atoi(lineParts[1])
	if err != nil {
		log.Fatalf("Failed to parse column number: %v", err)
	}
	return LinterMessage{
		Row:       lineNumber,
		Column:    columnNumber,
		ErrorCode: lineParts[2],
		Message:   lineParts[3],
	}
}

func (*PythonAnalysis) Analyze(cli *docker.DockerClient, scriptPath string) ([]LinterMessage, error) {
	ctx := context.Background()
	output, err := cli.Cmd(ctx, []string{"pylint", scriptPath}, "")
	if err != nil {
		log.Printf("Cmd() failed: %v", err)
		return nil, err
	}
	pylintOutputs := parsePylintOutput(output)
	return pylintOutputs, nil
}
