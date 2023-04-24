package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const DEBUG = false

type FileType struct {
	Name      string
	Extension string
	Category  string
}

func main() {
	userDir := getUserHomeDir()
	df := userDir + "/Downloads/"

	cat := []string{"Documents", "Executables", "Images",
		"Fonts", "Sounds", "Misc", "Compressed", "Videos"}

	fileTypes := []FileType{
		{Name: "PDF", Extension: ".pdf", Category: cat[0]},
		{Name: "Text", Extension: ".txt", Category: cat[0]},
		{Name: "Word", Extension: ".docx", Category: cat[0]},
		{Name: "Exel", Extension: ".xlsx", Category: cat[0]},
		{Name: "Powerpoint", Extension: ".pptx", Category: cat[0]},
		{Name: "Program", Extension: ".exe", Category: cat[1]},
		{Name: "Setup", Extension: ".msi", Category: cat[1]},
		{Name: "ZIP", Extension: ".zip", Category: cat[6]},
		{Name: "RAR", Extension: ".rar", Category: cat[6]},
		{Name: "PNG", Extension: ".png", Category: cat[2]},
		{Name: "JPEG", Extension: ".jpg", Category: cat[2]},
		{Name: "BMP", Extension: ".bmp", Category: cat[2]},
		{Name: "TIFF", Extension: ".tiff", Category: cat[2]},
		{Name: "GIF", Extension: ".gif", Category: cat[2]},
		{Name: "MP3", Extension: ".mp3", Category: cat[4]},
		{Name: "WAVE", Extension: ".wav", Category: cat[4]},
		{Name: "MP4", Extension: ".mp4", Category: cat[4]},
		{Name: "Playlist", Extension: ".pls", Category: cat[4]},
		{Name: "OGG", Extension: ".ogg", Category: cat[4]},
		{Name: "MKV", Extension: ".mkv", Category: cat[7]},
		{Name: "AVI", Extension: ".avi", Category: cat[7]},
		{Name: "MPG", Extension: ".mpg", Category: cat[7]},
		{Name: "Font", Extension: ".ttf", Category: cat[3]},
	}

	fmt.Println("Watch Dawg started...")

	fmt.Println("\nChecking category folders...")
	for _, c := range cat {
		checkFolder(df, c)
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
							if err := os.Rename(file, df+fileType.Category+"/"+fileName); err != nil {
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
	folder := df + folderName
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
