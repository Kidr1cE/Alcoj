package analysis

import (
	"log"
	"testing"
)

func TestAnalysisPython(t *testing.T) {
	targetScript := "../../cmd/factory/main.py"

	result, err := AnalysisPython(targetScript)
	if err != nil {
		t.Errorf("Error running AnalysisPython: %v", err)
	}
	log.Println("Result: ", result)
}
