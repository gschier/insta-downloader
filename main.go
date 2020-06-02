package main

import (
	"errors"
	"fmt"
	"github.com/ahmdrz/goinsta/v2"
	"github.com/gschier/insta/internal"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify output directory as first argument")
		os.Exit(1)
	}

	dir := os.Args[1]
	username := getUsername()
	exportFileName := fmt.Sprintf(".goinsta.%s.json", username)
	exportFile := filepath.Join(os.TempDir(), exportFileName)

	var insta *goinsta.Instagram
	var initErr error

	if internal.FileExists(exportFile) {
		fmt.Println("Restored session info to:\n  ", exportFile)
		insta, initErr = goinsta.Import(exportFile)
	} else {
		pass := getPass()
		insta = goinsta.New(username, pass)
		initErr = insta.Login()
	}

	if initErr != nil {
		fmt.Println("Failed to initialize", initErr)
		os.Exit(1)
	}

	// Save export file again
	err := insta.Export(exportFile)
	if err != nil {
		fmt.Println("Failed to export config", err)
		os.Exit(1)
	}
	fmt.Println("Persisted session info to:\n  ", exportFile)

	imgDir := filepath.Join(dir, insta.Account.Username)
	numDownloaded := internal.DownloadImages(insta, imgDir)

	fmt.Printf("Finished downloading %d images to %s\n", numDownloaded, imgDir)
}

func getPass() string {
	prompt := promptui.Prompt{
		Label: "Password",
		Mask:  '*',
		Validate: func(s string) error {
			if s == "" {
				return errors.New("required")
			}
			return nil
		},
	}

	password, _ := prompt.Run()
	return password
}

func getUsername() string {
	prompt := promptui.Prompt{
		Label: "Username",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("required")
			}
			return nil
		},
	}

	username, _ := prompt.Run()
	return username
}
