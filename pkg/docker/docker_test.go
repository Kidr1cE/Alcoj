package docker

import (
	"alcoj/pkg/analysis"
	"context"
	"log"
	"strings"
	"testing"
)

func TestDocker(t *testing.T) {
	// docker client init
	cli, err := NewDocker("test")
	if err != nil {
		t.Errorf("NewDocker() failed: %v", err)
	}

	cli.ContainerID = "49c0c826177439649f0f52bd0181cc0f49bf79c22ebf4d8c5b828c55f9e1ed34"

	ctx := context.Background()
	defer cli.Clean(ctx)

	// Time
	output, err := cli.Cmd(ctx, []string{"/usr/bin/time", "-v", "ps"})
	if err != nil {
		t.Errorf("Cmd() failed: %v", err)
	}
	log.Println("output: ", output)

	lines := strings.Split(output, "\n")

	timeMessage := analysis.TimeMessage{}
	commandOutputs := lines[0 : len(lines)-24]
	timeOutputs := lines[len(lines)-24:]
	for _, line := range timeOutputs {
		analysis.ParseTimeLine(line, &timeMessage)
	}

	log.Println("timeMessage: ", timeMessage)
	log.Println("commandOutputs:")
	for _, output := range commandOutputs {
		log.Println(output)
	}

	// Pylint
	output, err = cli.Cmd(ctx, []string{"pylint", "main.py"})
	if err != nil {
		t.Errorf("Cmd() failed: %v", err)
	}
	log.Println("output: ", output)
	pylintOutputs := analysis.ParsePylintOutput(output)
	for _, message := range pylintOutputs {
		log.Println(message.Column, message.ErrorCode, message.LineNumber, message.Message)
	}
}
