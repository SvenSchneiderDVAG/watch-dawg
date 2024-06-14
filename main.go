package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
    "syscall"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const Debug = false
const ErrorMsg = "ERROR:"
const FormatString = "%s %s\n"
const ConfigFile = "config.json"

// leave empty to use your OS default download folder
const DownloadFolder = "" 

type FileTypes struct {
	Filetypes []FileType `json:"filetypes"`
}
type FileType struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Category  string `json:"category"`
}

func main() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigs
        fmt.Println("\nReceived an interrupt, goodbye...")
        os.Exit(0)
    }()

	df := getDownloadFolder()
	
	fileData, err := loadConfigFile()
    if err != nil {
        logErrorAndExit(err)
    }

	fileTypes := fileData.Filetypes

	printDebugInfo(fileTypes)

	fmt.Println("Watch Dawg started...")
	fmt.Println("\nChecking category folders...")

	createCategoryFolders(df, fileTypes)

	startWatching(df, fileTypes)
}

func createCategoryFolders(df string, fileTypes []FileType) {
	for _, fileType := range fileTypes {
		checkFolder(df, fileType.Category)
	}
	fmt.Printf("\nDone!\n\n")
    fmt.Printf("Observing download folder: %s\n", df)
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

func startWatching(df string, fileTypes []FileType) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        logErrorAndExit(err)
    }
    defer watcher.Close()

    done := make(chan bool)

    go processEvents(watcher, df, fileTypes)

    if err = watcher.Add(df); err != nil {
        logErrorAndExit(err)
    }

    <-done
}

func processEvents(watcher *fsnotify.Watcher, df string, fileTypes []FileType) {
    for {
        select {
        case event := <-watcher.Events:
            DebugPrint("Event:", event)
            processFiles(df, fileTypes)
        case err := <-watcher.Errors:
            fmt.Printf(FormatString, ErrorMsg, err)
        }
    }
}

func processFiles(df string, fileTypes []FileType) {
    for _, fileType := range fileTypes {
        files, err := WalkMatch(df, "*"+fileType.Extension)
        if err != nil {
            fmt.Printf("%s Walk directory: %s\n", ErrorMsg, err)
            continue
        }
        for _, file := range files {
            // Ignore temporary files created by browsers
            if strings.HasSuffix(file, ".crdownload") || strings.HasSuffix(file, ".part") || strings.HasSuffix(file, ".tmp") {
                continue
            }
            fmt.Printf("Moving file: %s to folder %s\n", file, fileType.Category)
            fileName := filepath.Base(file)
            destinationPath := filepath.Join(df, fileType.Category, fileName)
            if err := os.Rename(file, destinationPath); err != nil {
                fmt.Printf("%s can't move file %s: %s\n", ErrorMsg, file, err)
                return
            }
        }
    }
}

func WalkMatch(df, pattern string) ([]string, error) {
    var matches []string
    err := filepath.Walk(df, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Ignore errors caused by missing files
			if os.IsNotExist(err) {
				return nil
			}
            return err
        }
		// we are only scanning the root directory, not the subdirectories
        if info.IsDir() {
            if path != df {
                return filepath.SkipDir
            }
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
		fmt.Printf("%s can't get user home directory: %s\n", ErrorMsg, err)
		os.Exit(1)
	}
	return homeDir
}

func getDownloadFolder() string {
	var df string
	if DownloadFolder == "" {
		df = filepath.Join(getUserHomeDir(), "Downloads")
	} else {
		df = DownloadFolder
	}
	return df
}

func loadConfigFile() (FileTypes, error) {
    file, err := os.Open(ConfigFile)
    if err != nil {
        return FileTypes{}, fmt.Errorf("%s can't open config.json: %w", ErrorMsg, err)
    }
    defer file.Close()

    fileContents, err := os.ReadFile(ConfigFile)
    if err != nil {
        return FileTypes{}, fmt.Errorf("%s can't read config.json: %w", ErrorMsg, err)
    }

    var fileData FileTypes
    if err := json.Unmarshal(fileContents, &fileData); err != nil {
        return FileTypes{}, fmt.Errorf("%s can't decode config.json: %w", ErrorMsg, err)
    }
    return fileData, nil
}

func logErrorAndExit(err error) {
    fmt.Printf(FormatString, ErrorMsg, err)
    os.Exit(1)
}

func DebugPrint(str ...interface{}) {
	if Debug {
		fmt.Println(str...)
		return
	}
}

func printDebugInfo(fileTypes []FileType) {
	if Debug {
		for _, fileType := range fileTypes {
			fmt.Printf("Name: %s\nExtension: %s\nCategory: %s\n\n", fileType.Name, fileType.Extension, fileType.Category)
		}
	}
}
