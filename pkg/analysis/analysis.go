package analysis

import (
	"bytes"
	"os/exec"
	"strings"
)

func Analyze(filepath string, analysisShell string) ([]string, []string, error) {
	cmd := exec.Command("bash", analysisShell, filepath)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return []string{}, []string{}, err
	}

	// catch each line to []string
	lines := strings.Split(out.String(), "\n")

	return lines, []string{}, nil
}
