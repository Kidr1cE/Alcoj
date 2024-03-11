package analysis

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	PylintResult []PylintMessage
	OutputResult string
	TimeResult   TimeMessage
}

type PylintMessage struct {
	LineNumber int    // 6
	Column     int    // 18
	ErrorCode  string // C0103
	Message    string // Argument name "columnPosition" doesn't conform to snake_case naming style
}

func parsePylintOutput(output string) []PylintMessage {
	messages := []PylintMessage{}
	reg := regexp.MustCompile(`(\d+):(\d+):\s+(\w+)\:\s+(.+)`)
	matches := reg.FindAllString(output, -1)

	for _, match := range matches {
		message := parsePylintLine(match)
		messages = append(messages, message)
	}
	return messages
}

// parse line like "27:0: W0311: Bad indentation. Found 1 spaces, expected 4 (bad-indentation)"
func parsePylintLine(line string) PylintMessage {
	lineParts := strings.Split(line, ":")
	lineNumber, err := strconv.Atoi(lineParts[0])
	if err != nil {
		log.Fatalf("Failed to parse line number: %v", err)
	}
	columnNumber, err := strconv.Atoi(lineParts[1])
	if err != nil {
		log.Fatalf("Failed to parse column number: %v", err)
	}
	return PylintMessage{
		LineNumber: lineNumber,
		Column:     columnNumber,
		ErrorCode:  lineParts[2],
		Message:    lineParts[3],
	}
}

func runPylint(scriptPath string) ([]PylintMessage, error) {
	var out bytes.Buffer
	var errBuf bytes.Buffer

	pylintCmd := exec.Command("pylint", scriptPath)
	pylintCmd.Stdout = &out
	pylintCmd.Stderr = &errBuf

	// If pylint returns a non-zero exit code, it means there are errors in the script
	err := pylintCmd.Run()
	if err != nil && pylintCmd.ProcessState.ExitCode() == 1 {
		return []PylintMessage{}, err
	}

	return parsePylintOutput(out.String()), nil
}

func AnalysisPython(scriptPath string) (Result, error) {
	pylintResult, err := runPylint(scriptPath)
	if err != nil {
		return Result{}, err
	}

	outputs, timeResult, err := runTime(scriptPath)
	if err != nil {
		return Result{}, err
	}

	return Result{
		PylintResult: pylintResult,
		OutputResult: outputs,
		TimeResult:   timeResult,
	}, nil
}
