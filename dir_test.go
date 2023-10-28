package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Setup
	tmpDir := filepath.Join(currentDir, "tmp")
	err = os.Mkdir(tmpDir, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	// Check is empty
	isEmpty, err := IsEmpty(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if !isEmpty {
		t.Errorf("expected %v to be empty but got nonempty", tmpDir)
	}
	// Check not empty
	randomFile := filepath.Join(tmpDir, "foobar.txt")
	err = createFile(randomFile, "")
	if err != nil {
		t.Fatal(err)
	}
	// Check is empty
	isEmpty, err = IsEmpty(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if isEmpty {
		t.Errorf("expected %v to be empty but got nonempty", tmpDir)
	}
	// Cleanup
	err = os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
}
