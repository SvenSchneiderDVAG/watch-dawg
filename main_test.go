package main

import (
    "testing"
    "os"
    "path/filepath"
)

func TestGetUserHomeDir(t *testing.T) {
    homeDir := getUserHomeDir()
    if homeDir == "" {
        t.Errorf("Expected a home directory, got an empty string")
    }
}

func TestLoadConfigFile(t *testing.T) {
    _, err := loadConfigFile()
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}

func TestWalkMatch(t *testing.T) {
    tempDir, err := os.MkdirTemp("", "test")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tempDir)

    tempFile := filepath.Join(tempDir, "test.txt")
    os.Create(tempFile)

    matches, err := WalkMatch(tempDir, "*.txt")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if len(matches) == 0 {
        t.Errorf("Expected at least one match, got none")
    }
}