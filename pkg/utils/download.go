package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
)

// DownloadFile downloads a file to a specific location
func DownloadFile(file, url string) {
	// url = url + "?archive=false"
	var err error
	var req *grab.Request
	if file, err = filepath.Abs(file); err != nil {
		log.Fatal(err)
	}

	client := grab.NewClient()
	req, err = grab.NewRequest(file, url)
	if err != nil {
		log.Panic(err)
	}
	resp := client.Do(req)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		log.Panicf("Download failed: %v\n", err)
	}

	if req.Filename != resp.Filename {
		log.Panicf("Requested file downloaded to wrong filename Has: %s Wants %s", resp.Filename, req.Filename)
	}

	log.Printf("Download saved to %s", resp.Filename)
}

// FileMD5 calculates the MD5 sum of a file
func FileMD5(path string) []byte {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}
