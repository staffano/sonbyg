package utils

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
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

// IsDir checks if the path is a dir
func IsDir(p string) bool {
	info, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CopyFile copies a file from src to dst and setting the mode of dst
func CopyFile(src, dst string, mode os.FileMode) {
	dstPath := filepath.Dir(dst)
	os.MkdirAll(dstPath, mode)
	out, err := os.OpenFile(dst,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		mode,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	source, err := os.Open(src)
	if err != nil {
		log.Fatalf("Could not open src: [%s], err: %v", src, err)
	}
	defer source.Close()
	_, err = io.Copy(out, source)
	if err != nil {
		log.Fatal(err)
	}
}

// Install copies stuff from src to dst
func Install(src, dst string, mode os.FileMode) {

	var err error

	if len(src) == 0 || len(dst) == 0 {
		log.Fatalf("Empty paths not allowed. src: [%s], dst: [%s]", src, dst)
	}

	if !PathExists(src) {
		log.Fatalf("The src path does not exist: %s", src)
	}

	// Check if overwrite or insert
	overwrite := true
	if dst[len(dst)-1] == '/' {
		overwrite = false
	}

	// remove any trailing path separator from src, they don't matter
	src = strings.TrimRight(src, "\\/")

	if src, err = filepath.Abs(src); err != nil {
		log.Fatalf("src path erroneous: %s, %v", src, err)
	}
	if dst, err = filepath.Abs(dst); err != nil {
		log.Fatalf("dst path erroneous: %s, %v", dst, err)
	}

	src = filepath.ToSlash(src)
	dst = filepath.ToSlash(dst)
	// Handle the case when src is a single file
	if FileExists(src) {
		if !IsDir(dst) {
			log.Fatalf("When src is a directory, the dst must also be a directory. src: %s, dst: %s", src, dst)
		}
		dstPath := path.Join(dst, filepath.Base(src))
		CopyFile(src, dstPath, mode)
		return
	}
	if !overwrite {
		var sb strings.Builder
		sb.WriteString(dst)
		sb.WriteRune(os.PathSeparator)
		sb.WriteString(filepath.Base(src))
		dst = sb.String()
	}

	walkFn := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", p, err)
			return err
		}
		p = filepath.ToSlash(p)
		relDir := strings.TrimPrefix(p, src)
		d := path.Join(dst, relDir)
		if !info.IsDir() {
			CopyFile(p, d, mode)
		} else {
			CreateDir(d, "", mode)
		}
		log.Printf("%s", d)
		return nil
	}

	log.Printf("Copying from dst=%s", dst)
	err = filepath.Walk(src, walkFn)
	if err != nil {
		log.Fatalf("Could not Install src: %s, dst %s  %v", src, dst, err)
	}
	return
}
