package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/otiai10/copy"
	"github.com/ulikunitz/xz"
)

const (
	// CompressorNone - uncompressed file
	CompressorNone = iota

	// CompressorGzip indicates that the file i compressed using gzip
	CompressorGzip

	// CompressorBzip2 indicates a bzip2 compressor
	CompressorBzip2

	// CompressorZip indicates a zlib compressor
	CompressorZip

	// CompressorXZ indicates a xz compressor
	CompressorXZ
)

// CreateFile creates a file from an io.Reader with a given mode
func CreateFile(src io.Reader, targetdir, filename string, mode os.FileMode) {
	filePath := path.Join(targetdir, filename)
	// Make sure the directory exists...
	dstPath := filepath.Dir(filePath)
	os.MkdirAll(dstPath, mode)
	out, err := os.OpenFile(filePath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		mode,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateDir creates a directory at the given path.
func CreateDir(targetDir, dirPath string, mode os.FileMode) {
	os.MkdirAll(path.Join(targetDir, dirPath), mode)
}

// UnpackFile ...
func UnpackFile(filePath string, targetDir string) {

	compressor := CompressorNone
	tarfile := false
	suffix := ""

	if strings.HasSuffix(filePath, ".tgz") {
		suffix = ".tgz"
		compressor = CompressorGzip
		tarfile = true
	} else if strings.HasSuffix(filePath, ".tar.gz") {
		suffix = ".tar.gz"
		compressor = CompressorGzip
		tarfile = true
	} else if strings.HasSuffix(filePath, ".gz") {
		suffix = ".gz"
		compressor = CompressorGzip
	} else if strings.HasSuffix(filePath, ".tar.bz2") {
		suffix = ".tar.bz2"
		compressor = CompressorBzip2
		tarfile = true
	} else if strings.HasSuffix(filePath, ".bz") {
		suffix = ".bz"
		compressor = CompressorBzip2
	} else if strings.HasSuffix(filePath, ".zip") {
		suffix = ".zip"
		compressor = CompressorZip
	} else if strings.HasSuffix(filePath, ".tar.xz") {
		suffix = ".tar.xz"
		compressor = CompressorXZ
		tarfile = true
	} else if strings.HasSuffix(filePath, ".xz") {
		suffix = ".xz"
		compressor = CompressorXZ
	}

	var src io.Reader

	if compressor == CompressorZip {
		src, _ := zip.OpenReader(filePath)
		for _, file := range src.Reader.File {
			zippedFile, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer zippedFile.Close()
			if file.FileInfo().IsDir() {
				CreateDir(targetDir, file.Name, file.Mode())
			} else {
				CreateFile(zippedFile, targetDir, file.Name, file.Mode())
			}
		}

		return
	}
	// Get mode of archive
	fi, err := os.Lstat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	archMode := fi.Mode().Perm()

	// Open archive file
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	switch compressor {
	case CompressorGzip:
		src, err = gzip.NewReader(f)
		if err != nil {
			log.Fatalf("Could not create gzip reader for file %s", filePath)
		}
	case CompressorBzip2:
		src = bzip2.NewReader(f)

	case CompressorXZ:
		src, err = xz.NewReader(f)
		if err != nil {
			log.Fatalf("Could not create xz reader for file %s", filePath)
		}

	default:
		log.Fatalf("Unknown file compressor: %d", compressor)
	}

	if !tarfile {
		// Just write the result to the file without the extension
		name := strings.TrimSuffix(filePath, suffix)
		CreateFile(src, targetDir, name, archMode)
		return
	}

	// It is a tar file beneath the archive
	tr := tar.NewReader(src)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Header %s:\n", hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeGNUSparse:
			{
				filePath := path.Join(targetDir, hdr.Name)
				// Make sure the directory exists...
				dstPath := filepath.Dir(filePath)
				os.MkdirAll(dstPath, os.FileMode(hdr.Mode))

				outputFile, err := os.OpenFile(
					filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
					os.FileMode(hdr.Mode),
				)
				if err != nil {
					log.Fatal(err)
				}
				defer outputFile.Close()
				if _, err := io.Copy(outputFile, tr); err != nil {
					log.Fatal(err)
				}
			}
		case tar.TypeDir:
			os.MkdirAll(path.Join(targetDir, hdr.Name), os.FileMode(hdr.Mode))
		default:
			log.Fatalf("Unknown tar tag type [%d] for tag [%s]", hdr.Typeflag, hdr.Name)
		}
	}
}

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

// FileExists returns true if there is a file at the path
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// PathExists returns true if the path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetUniqueFilename returns a path in the form
// <outDir>/<id>.<nr>, where the nr is generated in a way
// to make the path unique.
func GetUniqueFilename(id, outDir string) string {
	var r string
	oldfiles, err := filepath.Glob(fmt.Sprintf("%s/%s.*", outDir, id))
	if err != nil {
		log.Fatal(err)
	}
	if len(oldfiles) == 0 {
		r, err = filepath.Abs(path.Join(outDir, id))
		if err != nil {
			log.Panic(err)
		}
		return r
	}
	for _, i := range oldfiles {
		li := strings.LastIndexByte(i, '.')
		if li == -1 {
			r = path.Join(outDir, id+".1")
			break
		}
	}
	rs, err := filepath.Abs(r)
	if err != nil {
		log.Fatal(err)
	}
	return rs
}

// FilenameFromURL extracts the filename from an URL
func FilenameFromURL(URL string) string {
	url, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	return path.Base(url.Path)
}

// CopyDir copies directory at src to dst, recursively
func CopyDir(src, dst string) {
	if err := copy.Copy(src, dst); err != nil {
		log.Panicf("Could not copy from %s to %s : %v", src, dst, err)
	}
}
