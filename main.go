package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type FileType struct {
	Name      string
	Extension string
	Category  string
}

func main() {
	userDir := getUserHomeDir()
	downloadFolder := userDir + "\\Downloads\\"

	fileTypes := []FileType{
		{Name: "PDF", Extension: ".pdf", Category: "Documents"},
		{Name: "Text", Extension: ".txt", Category: "Documents"},
		{Name: "Word", Extension: ".docx", Category: "Documents"},
		{Name: "Program", Extension: ".exe", Category: "Executables"},
		{Name: "Setup", Extension: ".msi", Category: "Executables"},
		{Name: "ZIP", Extension: ".zip", Category: "Compressed"},
		{Name: "RAR", Extension: ".rar", Category: "Compressed"},
		{Name: "PNG", Extension: ".png", Category: "Images"},
		{Name: "JPEG", Extension: ".jpg", Category: "Images"},
		{Name: "GIF", Extension: ".gif", Category: "Images"},
		{Name: "MP3", Extension: ".mp3", Category: "Sounds"},
		{Name: "WAVE", Extension: ".wav", Category: "Sounds"},
		{Name: "MP4", Extension: ".mp4", Category: "Sounds"},
		{Name: "Playlist", Extension: ".pls", Category: "Sounds"},
		{Name: "OGG", Extension: ".ogg", Category: "Sounds"},
		{Name: "MKV", Extension: ".mkv", Category: "Videos"},
		{Name: "AVI", Extension: ".avi", Category: "Videos"},
		{Name: "MPG", Extension: ".mpg", Category: "Videos"},
		{Name: "Font", Extension: ".ttf", Category: "Fonts"},
	}

	fmt.Println("Watch Dawg started...")
	setupAll(downloadFolder)
	fmt.Println("Observing download folder: ", downloadFolder)

	// sort.Slice(fileTypes, func(i, j int) bool {
	// 	return fileTypes[i].Category < fileTypes[j].Category
	// })

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("Event:", event)
				for _, fileType := range fileTypes {
					files, err := WalkMatch(downloadFolder, "*"+fileType.Extension)
					if err != nil {
						fmt.Printf("Error while walking %s: %s\n", downloadFolder, err)
					} else {
						for _, file := range files {
							// fmt.Printf("Moving file %s to folder %s.\n", file, fileType.Category)
							fileName := filepath.Base(file)
							// fmt.Println("Name of file: ", fileName)
							err := os.Rename(file, downloadFolder+fileType.Category+"\\"+fileName)
							if err != nil {
								fmt.Printf("Error while moving file %s: %s\n", file, err)
								return
							}
						}
					}
				}
			case err := <-watcher.Errors:
				fmt.Println("Error: ", err)
			}
		}
	}()

	err = watcher.Add(downloadFolder)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	<-done
}

func setupAll(downloadFolder string) {
	fmt.Println("Checking category folders...")

	checkFolder(downloadFolder, "Compressed")
	checkFolder(downloadFolder, "Documents")
	checkFolder(downloadFolder, "Executables")
	checkFolder(downloadFolder, "Images")
	checkFolder(downloadFolder, "Sounds")
	checkFolder(downloadFolder, "Fonts")
}

func checkFolder(downloadFolder, folderName string) {
	folder := downloadFolder + folderName
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0777)
		fmt.Printf("Creating new category folder: '%s'\n", folder)
	} else {
		// fmt.Printf("Folder '%s' does already exist.\n", folder)
		return
	}
}

func WalkMatch(downloadFolder, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(downloadFolder, func(path string, info os.FileInfo, err error) error {
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
		fmt.Printf("Error while getting User home directory: %s\n", err)
	}
	return homeDir
}
