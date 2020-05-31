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
		fmt.Println("No username provided")
		os.Exit(1)
	}

	username := os.Args[1]
	exportFileName := fmt.Sprintf(".goinsta.%s.json", username)
	exportFile := filepath.Join(os.TempDir(), exportFileName)

	var insta *goinsta.Instagram
	var initErr error

	if internal.FileExists(exportFile) {
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
	fmt.Println("Stored session info to", exportFile)

	imgDir := "images/" + insta.Account.Username
	err = os.MkdirAll(imgDir, 0755)
	if err != nil {
		panic("Failed to create image dir " + err.Error())
	}

	internal.DownloadImages(insta, imgDir)
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
