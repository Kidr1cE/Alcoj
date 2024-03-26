package analysis

import (
	"alcoj/pkg/docker"
	"context"
	"log"
	"strings"
)

type TimeMessage struct {
	Command             string `json:"command,omitempty"`
	SystemTimeSeconds   string `json:"system_time_seconds,omitempty"`
	UserTimeSeconds     string `json:"user_time_seconds,omitempty"`
	PercentCPU          string `json:"percent_cpu,omitempty"`
	ElapsedWallClock    string `json:"elapsed_wall_clock,omitempty"`
	AvgSharedTextSize   string `json:"avg_shared_text_size,omitempty"`
	AvgUnsharedDataSize string `json:"avg_unshared_data_size,omitempty"`
	AvgStackSize        string `json:"avg_stack_size,omitempty"`
	AvgTotalSize        string `json:"avg_total_size,omitempty"`
	MaxResidentSetSize  string `json:"max_resident_set_size,omitempty"`
	AvgResidentSetSize  string `json:"avg_resident_set_size,omitempty"`
	MajorPageFaults     string `json:"major_page_faults,omitempty"`
	MinorPageFaults     string `json:"minor_page_faults,omitempty"`
	VoluntarySwitches   string `json:"voluntary_switches,omitempty"`
	InvoluntarySwitches string `json:"involuntary_switches,omitempty"`
	Swaps               string `json:"swaps,omitempty"`
	FileSystemInputs    string `json:"file_system_inputs,omitempty"`
	FileSystemOutputs   string `json:"file_system_outputs,omitempty"`
	SocketMessagesSent  string `json:"socket_messages_sent,omitempty"`
	SocketMessagesRecv  string `json:"socket_messages_recv,omitempty"`
	PageSize            string `json:"page_size,omitempty"`
	ExitStatus          string `json:"exit_status,omitempty"`
}

type TimeInterface interface {
	RunTime(cli *docker.DockerClient, filename string, cmd []string) (string, TimeMessage, error)
}

type Runner struct{}

func (*Runner) RunTime(cli *docker.DockerClient, filename string, cmd []string, input string) (string, TimeMessage, error) {
	ctx := context.Background()
	output, err := cli.Cmd(ctx, cmd, input)
	if err != nil {
		log.Printf("Cmd() failed: %v", err)
		return "", TimeMessage{}, err
	}

	lines := strings.Split(output, "\n")

	timeMessage := TimeMessage{}
	commandOutputs := parseOutput(lines[0 : len(lines)-24])
	timeOutputs := lines[len(lines)-24:]
	for _, line := range timeOutputs {
		parseTimeLine(line, &timeMessage)
	}

	log.Println("timeMessage: ", timeMessage)
	log.Println("commandOutputs:", commandOutputs)

	return commandOutputs, timeMessage, nil
}

func parseOutput(lines []string) string {
	var outputBuffer strings.Builder
	for _, line := range lines {
		outputBuffer.WriteString(line)
		outputBuffer.WriteString("\n")
	}
	return outputBuffer.String()
}

func parseTimeLine(line string, timeMessage *TimeMessage) bool {
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
