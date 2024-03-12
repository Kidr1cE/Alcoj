package analysis

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TimeMessage struct {
	Command             string
	SystemTimeSeconds   string
	UserTimeSeconds     string
	PercentCPU          string
	ElapsedWallClock    string
	AvgSharedTextSize   string
	AvgUnsharedDataSize string
	AvgStackSize        string
	AvgTotalSize        string
	MaxResidentSetSize  string
	AvgResidentSetSize  string
	MajorPageFaults     string
	MinorPageFaults     string
	VoluntarySwitches   string
	InvoluntarySwitches string
	Swaps               string
	FileSystemInputs    string
	FileSystemOutputs   string
	SocketMessagesSent  string
	SocketMessagesRecv  string
	PageSize            string
	ExitStatus          string
}

func parseTimeOutput(output string) TimeMessage {
	timeMessage := TimeMessage{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		ParseTimeLine(line, &timeMessage)
	}
	return timeMessage
}

func ParseTimeLine(line string, timeMessage *TimeMessage) bool {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Command being timed":
		timeMessage.Command = value
	case "User time (seconds)":
		timeMessage.UserTimeSeconds = value
	case "System time (seconds)":
		timeMessage.SystemTimeSeconds = value
	case "Percent of CPU this job got":
		timeMessage.PercentCPU = value
	// case "Elapsed (wall clock) time (h:mm:ss or m:ss)":
	// 	timeMessage.ElapsedWallClock = value
	case "Average shared text size (kbytes)":
		timeMessage.AvgSharedTextSize = value
	case "Average unshared data size (kbytes)":
		timeMessage.AvgUnsharedDataSize = value
	case "Average stack size (kbytes)":
		timeMessage.AvgStackSize = value
	case "Average total size (kbytes)":
		timeMessage.AvgTotalSize = value
	case "Maximum resident set size (kbytes)":
		timeMessage.MaxResidentSetSize = value
	case "Average resident set size (kbytes)":
		timeMessage.AvgResidentSetSize = value
	case "Major (requiring I/O) page faults":
		timeMessage.MajorPageFaults = value
	case "Minor (reclaiming a frame) page faults":
		timeMessage.MinorPageFaults = value
	case "Voluntary context switches":
		timeMessage.VoluntarySwitches = value
	case "Involuntary context switches":
		timeMessage.InvoluntarySwitches = value
	case "Swaps":
		timeMessage.Swaps = value
	case "File system inputs":
		timeMessage.FileSystemInputs = value
	case "File system outputs":
		timeMessage.FileSystemOutputs = value
	case "Socket messages sent":
		timeMessage.SocketMessagesSent = value
	case "Socket messages received":
		timeMessage.SocketMessagesRecv = value
	case "Page size (bytes)":
		timeMessage.PageSize = value
	case "Exit status":
		timeMessage.ExitStatus = value
	default:
		return false
	}

	return true
}

func runTime(scriptPath string) (string, TimeMessage, error) {
	timeCmd := exec.Command("/usr/bin/time", "-v", "python3", scriptPath)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	timeCmd.Stdout = &stdoutBuf
	timeCmd.Stderr = &stderrBuf

	err := timeCmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		return "", TimeMessage{}, err
	}

	// Read the standard output of the command
	scanner := bufio.NewScanner(&stdoutBuf)
	var programOutputBuffer bytes.Buffer
	for scanner.Scan() {
		programOutputLine := scanner.Text()
		programOutputBuffer.WriteString(programOutputLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard output:", err)
	}

	timeOutput := stderrBuf.String()
	timeResult := parseTimeOutput(timeOutput)

	return programOutputBuffer.String(), timeResult, nil
}
