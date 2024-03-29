// Code generated by vfsgen; DO NOT EDIT.

// +build !dev

package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Assets statically implements the virtual filesystem provided to vfsgen.
var Assets = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2019, 9, 5, 19, 57, 45, 180787800, time.UTC),
		},
		"/001_install.patch": &vfsgen۰CompressedFileInfo{
			name:             "001_install.patch",
			modTime:          time.Date(2019, 9, 5, 19, 57, 45, 180787800, time.UTC),
			uncompressedSize: 2889,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x96\x5f\x6f\xdb\x36\x17\xc6\xaf\x5f\x7d\x8a\xf3\xc2\xbb\x68\x21\x53\x24\x25\x59\x52\x04\x04\x48\x83\xb5\x43\x81\x36\x2d\xec\x00\x43\x31\x0c\x06\x23\x1f\xc9\x5c\x64\x49\x23\xa9\xb8\xc2\xba\x7d\xf6\x81\xfe\x93\xda\x8d\xe5\xda\x5d\xd1\x5c\x50\xf1\xe1\xf3\x3c\x87\x04\x7f\x17\x67\x26\xf3\x1c\xc8\x8d\x68\x15\x2c\x74\xa7\xa3\x90\xa2\xc9\x68\x23\xb2\x85\xa8\xbc\xac\xae\xf2\x4d\xf9\xfa\xcb\xba\x43\x08\xe9\xb1\xfc\xcf\x67\xfc\x82\xb0\x11\xf1\x03\xe0\x41\x1a\xf2\x94\x73\x8f\x6d\xff\xc0\x65\x3e\x63\x8e\xeb\xba\x7d\xd1\x1b\xbf\x8d\x00\xee\xa7\xe1\x28\xe5\x81\x17\xc4\x3c\x88\x92\xe4\xd1\x7f\x75\x05\x24\x88\x87\x09\xb8\xab\xf5\xea\xca\x01\x07\x06\x70\xdd\xc1\x0c\x73\xd1\x96\x66\x08\xeb\x48\x10\x59\x86\x8d\xd1\xf6\xe7\xbd\x28\x50\x83\x96\x45\x85\x33\xb8\xeb\xe0\x1e\x3b\x0d\x66\x2e\x0c\x48\xa3\xa1\xac\x33\x51\xda\x9a\x92\x55\x61\xc3\x8c\x6a\xb5\xd1\xf0\x4c\x23\x6e\xc2\xc8\x3d\x76\x20\xaa\xd9\x4a\x6f\xc3\x1b\x51\xe0\xf3\x21\x08\x0d\x4b\x2c\x4b\xfb\x6d\xab\x4d\xfe\xb6\x9f\xe7\x90\xc1\x44\x16\x6f\xf0\x01\x4b\xb8\x84\x1b\x7c\x40\xe5\x90\xc7\x0a\x00\x5c\xc2\x18\xff\x6c\xa5\xc2\x19\xfc\x2c\x8c\xb8\x13\x1a\xdf\x35\x46\xd6\x95\x28\x1d\xf7\x89\xd5\x1d\x9c\xea\x85\x37\xf6\x46\xaf\x64\x89\x3b\x21\x9f\x77\x07\x63\x5c\xd4\x06\xbf\xd8\xdf\xc6\x39\xe0\x1c\x85\x63\x46\x17\x52\xa9\x5a\x95\x52\x1b\x6f\x21\xab\x62\x19\xf8\x87\x5e\xf4\x90\xae\x07\x9e\x43\x52\x0b\x43\x42\x58\x44\x58\x0c\x9c\xa5\x7e\x94\x8e\x92\x53\x61\xea\xcb\xdb\x85\x2b\x4c\x83\xd0\x63\x81\xef\x8f\x92\x68\x0f\xae\x70\x18\x83\x6b\x97\x0d\x5a\x03\x78\xaf\xe4\x42\xa8\x6e\xf5\xbf\xed\xe6\x7b\xb5\x2a\x1c\x77\x82\xea\x01\x15\x5c\xc2\xdc\x98\x46\xa7\x74\xd3\xd4\xeb\x44\x35\xc3\x8f\x9e\x6a\x37\x05\x4d\x57\x26\xba\x3a\x08\x95\x51\x12\x51\x07\xf6\xcc\x29\xa5\x0a\x9b\xda\x7b\x0c\x3f\xa2\xb5\x8d\x74\xdd\xaa\x0c\xf3\x5a\x15\xe8\x55\x68\x68\xa3\xea\x3f\x30\x33\xdb\x46\xb9\x2c\x51\xd3\xf1\xcb\xf7\xef\x26\xf4\xed\xeb\x9b\x5f\x7e\xed\x6b\xba\x5c\x2e\x7d\x2f\x6f\x4d\xab\x70\x29\x14\x7a\xc2\xd0\x7f\x2a\x99\xdd\xd7\xb8\x4e\x22\xeb\x0b\xec\x9d\x86\xfc\x97\x6b\x9f\x09\x56\x14\x9e\xf6\xba\x51\x78\x32\x58\x51\xf8\x7d\xc1\x5a\xe7\xed\x83\xc5\x2f\xbc\x38\x09\xa2\x28\x0e\x7f\x24\x58\x1f\x93\x68\x1a\x85\xa7\xa2\x75\x58\xfd\x4d\x70\xf5\x36\x3e\x0f\xaf\x6d\xcc\x37\x02\xb6\xb5\x9f\x81\x98\xee\xf4\xd7\x1f\x59\x77\xfa\x24\xb8\x74\xa7\xbf\x1f\x59\x9b\xb0\x7d\xac\x58\xe2\x45\x61\x14\x87\x11\xff\x41\x58\xe9\x4e\xd3\x9f\x84\xca\xe6\x27\x40\xd5\xab\x3d\x17\xa9\xc9\x87\x89\xdf\xdb\xf5\x64\xa2\x76\x8e\x73\x3e\x4f\x3b\xe6\x1e\x9a\x6a\x6d\x88\xac\xb4\x11\x65\x49\x59\x4c\x3e\xcf\x08\x9e\xdd\xda\x7f\xdf\xe3\xda\x27\x6c\x1d\x97\xaf\x11\xe3\x3e\xe1\x23\xe0\x17\x69\x18\xa5\x41\xbc\x87\x18\x3f\x84\xd8\x57\x32\xf9\x45\xcc\x08\xe3\x84\x71\x60\x3c\x65\x2c\x65\xec\x49\xa6\x25\x8d\x0f\x79\x00\x2e\x1b\x32\x8b\x1a\x59\x88\xee\x0e\xa7\xb2\x92\x66\xba\x19\x9d\xe0\xd9\x73\x87\xfc\xe5\x10\x00\x99\xc3\x6f\xf0\x7f\x20\x33\xd8\x67\xbc\xa8\xda\xa6\x80\xdf\xad\xc4\xcc\xb1\xb2\x5f\x00\xda\x6a\x45\xef\x64\x45\x77\x46\x2d\x42\x6c\xf0\xb1\xfd\xa6\x6e\xda\x52\x18\x5c\x13\x0e\x9f\x3e\xd9\xa9\x0d\x8f\x39\x14\xe6\x0a\xf5\x9c\xac\x86\xbf\x3d\xfd\x7a\x7d\xfb\xe2\xc3\xf5\xcb\xe9\xab\xd7\xe3\xc9\xed\x74\x72\xfb\x62\x7c\x7b\xb9\x55\xe4\xd2\x21\x7f\x3b\xe4\xd0\x9d\x9d\x7f\x03\x00\x00\xff\xff\xed\xec\x84\xb0\x49\x0b\x00\x00"),
		},
		"/001_mirrors.patch": &vfsgen۰CompressedFileInfo{
			name:             "001_mirrors.patch",
			modTime:          time.Date(2019, 9, 4, 18, 52, 6, 460036200, time.UTC),
			uncompressedSize: 669,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x90\xbf\x6a\xf3\x30\x14\xc5\xe7\xcf\x4f\x71\x21\xa3\x3e\xe9\xca\xc6\x91\x1d\x43\xc1\x14\x4a\xe9\xd0\x34\x34\x43\xc7\x22\x1c\x39\x55\x1b\x5b\xe6\x4a\xae\x93\xa5\xcf\x5e\xf2\x17\x9a\x0e\x4d\xa1\x1a\x84\x90\x8e\x7e\xe7\x9c\xbb\xb0\x75\x0d\x7c\xaa\x7b\x82\xc6\x6f\xbc\x4a\xd1\x84\x0a\x3b\x5d\x35\xba\x15\x0b\x6c\x2c\x91\xa3\x95\xf5\x41\x34\xb6\x5d\x0e\x2a\x3d\xc8\xae\x7f\xd2\x45\x9c\xf3\x0b\x91\xff\x12\x19\xe7\x5c\x2a\x2e\x33\x88\x65\x91\xa8\x62\x9c\x0b\x79\x5c\xc0\x64\x22\x65\xc4\x18\xbb\xd4\x7a\xcb\x9b\x70\x39\xe1\x32\x85\x44\x16\xe3\xb8\x48\x33\xa1\xc6\x52\xe6\x71\x76\xe2\x95\x25\xf0\xf4\x7f\x06\x6c\xbb\x95\x65\x04\x11\x8c\x46\x30\x23\xdb\x68\xda\xec\xce\x5b\xb7\x44\x38\x5a\x46\x7c\x6e\xe8\xdd\x10\x5c\xc1\x4b\x08\x5d\x81\x48\xa6\x73\xe2\xf4\x8e\x3b\x5f\x5c\xe7\xea\x59\xa5\x78\xa6\xf6\x05\xa2\x77\x3d\x55\xa6\x76\xb4\x34\xa2\x35\x01\x3b\x72\xaf\xa6\x0a\x1e\x77\x08\xac\xed\xca\x78\x7c\xbc\x99\x3d\xcc\xf1\xfe\x6e\x7a\xfb\x74\x42\xc1\xb9\xf1\x30\x0c\x89\xa8\xfb\xd0\x93\x19\x34\x19\xa1\x03\x7e\xb4\xb6\x7a\x73\x66\xcf\xe2\xfb\x69\x9c\x25\x82\x6f\x89\xf6\x32\xb1\xd1\xed\xc2\xac\x05\xf5\x87\x8b\x63\xa2\xaf\xdf\xd9\xaf\xea\xb3\xbf\xab\xff\x19\x00\x00\xff\xff\xe6\xfb\x17\xf3\x9d\x02\x00\x00"),
		},
		"/README.md": &vfsgen۰CompressedFileInfo{
			name:             "README.md",
			modTime:          time.Date(2019, 9, 4, 18, 33, 48, 61313100, time.UTC),
			uncompressedSize: 137,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\xcb\xe1\x09\xc2\x40\x0c\x05\xe0\xff\x82\x3b\xbc\x09\x3a\x4d\x17\x88\x97\x77\xe6\x20\x26\xa5\x49\x71\x7d\x51\x70\x80\x6f\x37\x62\x2e\x67\x61\xa6\x2b\x4f\x8c\x8c\x96\x15\x05\xa9\x62\x17\xda\xa4\xf1\x5e\xee\x78\x10\x23\x5f\xc7\x72\x2a\x56\x74\xa2\x7f\x36\xc4\xf1\x4c\x8c\x54\x6e\xf7\xdb\x6e\x2c\xfe\xed\x90\xf8\xaa\xab\xa8\xe8\x84\xa5\x2b\x0e\xe9\x61\x2c\x48\x28\xb2\x8d\x27\xaa\xaf\x39\xb7\x4f\x00\x00\x00\xff\xff\xc3\x54\x18\x88\x89\x00\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/001_install.patch"].(os.FileInfo),
		fs["/001_mirrors.patch"].(os.FileInfo),
		fs["/README.md"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
