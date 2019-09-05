package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/shurcooL/httpfs/vfsutil"
)

// UnpackAssets unpacks assets into a destination directory
func UnpackAssets(dstDir string, assets http.FileSystem) error {

	walkFn := func(p string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", p, err)
			return nil
		}
		fmt.Println(p)
		if !fi.IsDir() {

			filePath := path.Join(dstDir, p)
			// Make sure the directory exists...
			dstPath := filepath.Dir(filePath)
			os.MkdirAll(dstPath, os.FileMode(fi.Mode()))

			outputFile, err := os.OpenFile(
				filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				os.FileMode(fi.Mode()),
			)
			if err != nil {
				log.Fatal(err)
			}
			defer outputFile.Close()
			if _, err := io.Copy(outputFile, r); err != nil {
				log.Fatal(err)
			}
			log.Printf("File %s copied to %s", fi.Name(), dstPath)
		}
		return nil
	}

	err := vfsutil.WalkFiles(assets, "/", walkFn)
	if err != nil {
		panic(err)
	}
	return nil
}
