package analysis

import (
	"log"
	"testing"
)

func TestAnalyze(t *testing.T) {
	filepath := "../../cmd/factory/main.py"
	analysisShell := "python.sh"
	// Call the function under test
	lines, _, err := Analyze(filepath, analysisShell)

	// Check the result
	if err != nil {
		t.Errorf("Analyze() failed: %v", err)
	}
	for _, line := range lines {
		log.Println(line)
	}
}
