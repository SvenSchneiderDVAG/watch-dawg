package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const DEBUG = false
const CONFIG_FILE = "config.json"

// leave empty to use default OS download folder
const DOWNLOAD_FOLDER = "" 

// Define the struct to hold the JSON data
type FileTypes struct {
	Filetypes []FileType `json:"filetypes"`
}
type FileType struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Category  string `json:"category"`
}

func main() {
	var df string
	if DOWNLOAD_FOLDER == "" {
		df = filepath.Join(getUserHomeDir(), "Downloads")
	} else {
		df = DOWNLOAD_FOLDER
	}

	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		log.Fatalf("ERROR: can't open config.json: %s\n", err)
	}
	defer file.Close()

	fileContents, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatalf("ERROR: can't read config.json: %s\n", err)
	}

	var fileData FileTypes
	if err := json.Unmarshal(fileContents, &fileData); err != nil {
		log.Fatalf("ERROR: can't decode config.json: %s\n", err)
	}

	fileTypes := fileData.Filetypes

	if DEBUG {
		for _, fileType := range fileTypes {
			fmt.Printf("Name: %s\nExtension: %s\nCategory: %s\n\n", fileType.Name, fileType.Extension, fileType.Category)
		}
	}

	fmt.Println("Watch Dawg started...")
	fmt.Println("\nChecking category folders...")

	// create category folders if they don't exist
	for _, fileType := range fileTypes {
		checkFolder(df, fileType.Category)
	}

	fmt.Printf("\nDone!\n\n")
	fmt.Printf("Observing download folder: %s\n", df)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("Event:", event)
				for _, fileType := range fileTypes {
					files, err := WalkMatch(df, "*"+fileType.Extension)
					if err != nil {
						fmt.Printf("ERROR: Walk directory: %s\n", err)
					} else {
						for _, file := range files {
							DebugPrint("Moving file", file, "to folder", fileType.Category)
							fileName := filepath.Base(file)
							destinationPath := filepath.Join(df, fileType.Category, fileName)
							if err := os.Rename(file, destinationPath); err != nil {
								fmt.Printf("ERROR: can't move file %s: %s\n", file, err)
								return
							}
						}
					}
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR:", err)
			}
		}
	}()

	if err = watcher.Add(df); err != nil {
		fmt.Println("ERROR:", err)
	}

	<-done
}

func checkFolder(df, folderName string) {
	folder := filepath.Join(df, folderName)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0777)
		DebugPrint("INFO: Creating new category folder:", folder)
	} else {
		DebugPrint("INFO: Folder already exists:", folder)
		return
	}
}

func WalkMatch(df, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(df, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func getUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("ERROR: can't get user home directory: %s\n", err)
		os.Exit(1)
	}
	return homeDir
}

func DebugPrint(str ...interface{}) {
	if DEBUG {
		fmt.Println(str...)
		return
	}
}
