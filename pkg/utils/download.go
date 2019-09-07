package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/url"
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

// DownloadGit downloads a url in the form
// src = https://github.com/staffano/sonbyg
// git clone https://github.com/staffano/sonbyg
// dst is the directory that maps to sonbyg.
func DownloadGit(src, dst string) {
	v := Variables{}
	ref := "master"
	su, _ := url.Parse(src)
	q := su.Query()
	if len(q) > 0 {
		ref = q.Get("ref")
	}

	u2 := url.URL{Scheme: su.Scheme, Host: su.Host, Path: su.Path}
	if IsDir(dst) {
		v["CWD"] = dst
		Exec("GitCheckoutRef", os.TempDir(), v, "git", "checkout", "master")
		Exec("GitPull", os.TempDir(), v, "git", "pull", "--depth", "1", "--ff-onley")
	} else {
		v["CWD"] = filepath.Dir(dst)
		Exec("GitClone", os.TempDir(), v, "git", "clone", "--depth", "1", u2.String(), dst)
	}
	v["CWD"] = dst
	Exec("GitCheckoutRef", os.TempDir(), v, "git", "checkout", ref)
	Exec("GitUpdateSubmodules", os.TempDir(), v, "git", "submodule", "update", "--init", "--recursive", "--depth", "1")
}
