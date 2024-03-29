package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

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
		//	Unzip(filePath, targetDir)
		src, _ := zip.OpenReader(filePath)
		defer src.Close()
		for _, file := range src.Reader.File {
			if file.FileInfo().IsDir() {
				CreateDir(targetDir, file.Name, file.Mode())
			} else {
				CreateFile2(file, targetDir, file.Name, file.Mode())
			}
		}
		log.Printf("Unzip done")
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

// Unzip ...
// cf. https://stackoverflow.com/a/24792688
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)
		log.Printf("%s", path)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
