package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type RunRequest struct {
	Language        string `json:"language,omitempty"`
	Code            string `json:"code,omitempty"`
	StaticAnalysis  bool   `json:"static_analysis,omitempty"`
	RuntimeAnalysis bool   `json:"runtime_analysis,omitempty"`
}

type RunResponse struct {
	Output          string           `json:"output,omitempty"`
	StaticAnalysis  []StaticAnalysis `json:"static_analysis,omitempty"`
	RuntimeAnalysis RuntimeAnalysis  `json:"runtime_analysis,omitempty"`
}

type StaticAnalysis struct {
	Row     int    `json:"row,omitempty"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message,omitempty"`
}

type RuntimeAnalysis struct {
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

func responseGenerator() *RunResponse {
	return &RunResponse{
		Output: "output",
		StaticAnalysis: []StaticAnalysis{
			{
				Row:     1,
				Column:  1,
				Message: "message",
			},
		},
		RuntimeAnalysis: RuntimeAnalysis{
			Command:             "command",
			SystemTimeSeconds:   "system_time_seconds",
			UserTimeSeconds:     "user_time_seconds",
			PercentCPU:          "percent_cpu",
			ElapsedWallClock:    "elapsed_wall_clock",
			AvgSharedTextSize:   "avg_shared_text_size",
			AvgUnsharedDataSize: "avg_unshared_data_size",
			AvgStackSize:        "avg_stack_size",
			AvgTotalSize:        "avg_total_size",
			MaxResidentSetSize:  "max_resident_set_size",
			AvgResidentSetSize:  "avg_resident_set_size",
			MajorPageFaults:     "major_page_faults",
			MinorPageFaults:     "minor_page_faults",
			VoluntarySwitches:   "voluntary_switches",
			InvoluntarySwitches: "involuntary_switches",
			Swaps:               "swaps",
			FileSystemInputs:    "file_system_inputs",
			FileSystemOutputs:   "file_system_outputs",
			SocketMessagesSent:  "socket_messages_sent",
			SocketMessagesRecv:  "socket_messages_recv",
			PageSize:            "page_size",
			ExitStatus:          "exit_status",
		},
	}
}

func (req *RunRequest) String() string {
	builder := &strings.Builder{}
	builder.WriteString("Language: ")
	builder.WriteString(req.Language)
	builder.WriteString("\n")
	builder.WriteString("Code: ")
	builder.WriteString(req.Code)
	builder.WriteString("\n")
	builder.WriteString("StaticAnalysis: ")
	builder.WriteString(strconv.FormatBool(req.StaticAnalysis))
	builder.WriteString("\n")
	builder.WriteString("RuntimeAnalysis: ")
	builder.WriteString(strconv.FormatBool(req.RuntimeAnalysis))
	return builder.String()
}

func handleCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func main() {
	http.HandleFunc("/alcoj/api/v1", func(w http.ResponseWriter, r *http.Request) {
		handleCORS(w, r)

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("Request body: ", r.Body)

		var req RunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		fmt.Println("Request: ", req.String())

		resp := responseGenerator()
		json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe(":8080", nil)
}
