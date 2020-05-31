package main

import (
	"fmt"
	"github.com/ahmdrz/goinsta/v2"
	"github.com/manifoldco/promptui"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var insta *goinsta.Instagram
	user := os.Args[1]

	exportFile := fmt.Sprintf("goinsta.%s.json", user)

	var initErr error
	if s, _ := os.Stat(exportFile); s == nil {
		pass := getPass()
		insta = goinsta.New(user, pass)
		initErr = insta.Login()
		_ = insta.Export(exportFile)
	} else {
		insta, initErr = goinsta.Import(exportFile)
	}

	if initErr != nil {
		panic("Failed to initialize " + initErr.Error())
	}

	imgDir := "images/" + insta.Account.Username
	err := os.MkdirAll(imgDir, 0755)
	if err != nil {
		panic("Failed to create image dir " + err.Error())
	}

	feed := insta.Account.Feed()
	batchSize := 8
	total := 0
	for feed.Next() {
		for i := 0; i < len(feed.Items); i += batchSize {
			wg := sync.WaitGroup{}
			for j := 0; j < batchSize && i+j < len(feed.Items); j++ {
				wg.Add(1)
				go downloadItem(&wg, imgDir, &feed.Items[i+j], total)
				total++
			}
			wg.Wait()
		}
	}
}

func downloadItem(w *sync.WaitGroup, dir string, item *goinsta.Item, count int) {
	defer w.Done()

	uri := item.Images.GetBest()
	if uri == "" {
		return
	}

	res, err := http.Get(uri)
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("%05d_%s.jpg", count, item.ID)
	fullPath := filepath.Join(dir, fileName)
	f, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloaded %d %s \"%s\"\n", count, fullPath, item.Caption.Text)
}

func getPass() string {
	prompt := promptui.Prompt{Label: "Password", Mask: '*'}
	password, _ := prompt.Run()
	return password
}
