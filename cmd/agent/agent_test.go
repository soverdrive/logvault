package main

import (
	"os"
	"testing"
)

func TestGetFileList(t *testing.T) {
	expect := []string{"agent.go", "agent_test.go"}

	dir, err := os.Getwd()
	if err != nil {
		t.Error("Canot get directory ", err.Error())
	}
	files, err := getFileList("agent.go,agent_test.go", dir)
	if err != nil {
		t.Error("Error when getting file list ", err.Error())
	}
	if len(expect) != len(files) {
		t.Errorf("Error, expecting %+v but got %+v", expect, files)
	}
}
