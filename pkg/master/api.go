package master

import (
	"alcoj/pkg/analysis"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	Code            string `json:"code"`
	Language        string `json:"language"`
	Input           string `json:"input"`
	RuntimeAnalysis bool   `json:"runtime_analysis"`
	StaticAnalysis  bool   `json:"static_analysis"`
}

type Response struct {
	Output          string                   `json:"output"`
	StaticAnalysis  []analysis.LinterMessage `json:"static_analysis"`
	RuntimeAnalysis analysis.TimeMessage     `json:"runtime_analysis"`
}

func startHttpServer() {
	http.HandleFunc("/alcoj/api/v1", alcojHandler)
	http.HandleFunc("/alcoj/api/v1/register", registerWorkerHandler)
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func alcojHandler(w http.ResponseWriter, r *http.Request) {
	handleCORS(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Request received: ", req)

	// Run
	res, err := master.Run(req.Language, req.Code, req.Input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("response: ", res.Output)

	// Remove input from output's head
	prefix := strings.Replace(req.Input, "\n", "\r\n", -1)
	res.Output = strings.TrimPrefix(res.Output, prefix)

	// Write response
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
		Suffix  string `json:"suffix"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := master.AddSandboxServer(req.Address, req.Suffix); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
