package internal

import (
	"fmt"
	"github.com/ahmdrz/goinsta/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func DownloadImages(insta *goinsta.Instagram, imgDir string) {
	feed := insta.Account.Feed()
	batchSize := 5
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

	fileName := fmt.Sprintf("%04d_%s.jpg", count, item.ID)
	fullPath := filepath.Join(dir, fileName)

	if FileExists(fullPath) {
		fmt.Printf("Skiping %s\n", fullPath)
		return
	}

	uri := item.Images.GetBest()
	if uri == "" {
		return
	}

	res, err := http.Get(uri)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloaded %s\n", fullPath)
}
