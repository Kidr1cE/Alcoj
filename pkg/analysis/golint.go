package analysis

import (
	"alcoj/pkg/docker"
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type GolangAnalysis struct{}

func parseGolintOutput(output string, scriptPath string) []LinterMessage {
	messages := []LinterMessage{}
	filename := scriptPath[strings.LastIndex(scriptPath, "/")+1:]
	log.Println("Filename:", filename)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, filename) {
			message := parseGolintLine(line)
			messages = append(messages, message)
		}
	}

	return messages
}

// parse line like "27:0: W0311: Bad indentation. Found 1 spaces, expected 4 (bad-indentation)"
func parseGolintLine(line string) LinterMessage {
	line = removeANSIEscapeCodes(line)
	piece := strings.Split(line, ":")
	row, err := strconv.Atoi(piece[0])
	if err != nil {
		log.Println("Failed to parse row:", err)
	}
	col, err := strconv.Atoi(piece[1])
	if err != nil {
		log.Println("Failed to parse col:", err)
	}
	var message string
	for i := 2; i < len(piece); i++ {
		message += piece[i]
	}

	return LinterMessage{
		Row:     row,
		Column:  col,
		Message: message,
	} // TODO: implement this function
}

func removeANSIEscapeCodes(s string) string {
	ansiEscapeRegex := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return ansiEscapeRegex.ReplaceAllString(s, "")
}

func (*GolangAnalysis) Analyze(cli *docker.DockerClient, scriptPath string) ([]LinterMessage, error) {
	ctx := context.Background()
	output, err := cli.Cmd(ctx, []string{"golangci-lint", "run", scriptPath}, "")
	if err != nil {
		log.Printf("Cmd() failed: %v", err)
		return nil, err
	}
	golintOutputs := parseGolintOutput(output, scriptPath)
	return golintOutputs, nil
}
