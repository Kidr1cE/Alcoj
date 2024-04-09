package master

import (
	"alcoj/pkg/analysis"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var port = os.Getenv("PORT")

type JudgeRequest struct {
	Code            string `json:"code"`
	Language        string `json:"language"`
	Input           string `json:"input"`
	RuntimeAnalysis bool   `json:"runtime_analysis"`
	StaticAnalysis  bool   `json:"static_analysis"`
}

type JudgeResponse struct {
	Output          string                   `json:"output"`
	StaticAnalysis  []analysis.LinterMessage `json:"static_analysis"`
	RuntimeAnalysis analysis.TimeMessage     `json:"runtime_analysis"`
}

type RegisterRequest struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Suffix  string `json:"suffix"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

func startHttpServer(stopCh chan struct{}) {
	defer func() {
		stopCh <- struct{}{}
	}()

	http.HandleFunc("/alcoj/api/v1", alcojHandler)
	http.HandleFunc("/alcoj/api/v1/register", registerWorkerHandler)

	log.Println("Server started at port: ", port)
	http.ListenAndServe(":"+port, nil)
}

func alcojHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received: ", r.Method, r.URL.Path)
	handleCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JudgeRequest
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

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := master.AddSandboxServer(req.ID, req.Address); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(RegisterResponse{Success: true}); err != nil {
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
