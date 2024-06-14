package main

import (
    "testing"
    "os"
    "path/filepath"
)

func TestGetUserHomeDir(t *testing.T) {
    homeDir := getUserHomeDir()
    if homeDir == "" {
        t.Errorf("Expected a non-empty string")
    }
}

func TestLoadConfigFile(t *testing.T) {
    _, err := loadConfigFile()
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}

func TestWalkMatch(t *testing.T) {
    // Create a temporary file for testing
    downloadFolder := getDownloadFolder()
    tempFile := filepath.Join(downloadFolder, "test.txt")
    os.Create(tempFile)
    defer os.Remove(tempFile)

    matches, err := WalkMatch(downloadFolder, "*.txt")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if len(matches) == 0 {
        t.Errorf("Expected at least one match, got none")
    }
}