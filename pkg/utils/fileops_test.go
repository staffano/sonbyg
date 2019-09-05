package utils

import (
	"testing"
)

func TestDownloads(t *testing.T) {
	CreateDir(".", "test", 744)
	uri := "http://www.lefeton.se/index.htm"
	DownloadFile("test/index.html", uri)
}
func TestDownloadLarge(t *testing.T) {
	CreateDir(".", "test", 744)
	uri := "http://repo.msys2.org/distrib/x86_64/msys2-base-x86_64-20190524.tar.xz"
	DownloadFile("test/msys2-base-x86_64-20190524.tar.xz", uri)
}

func TestUnpackLarge(t *testing.T) {
	UnpackFile("test/msys2-base-x86_64-20190524.tar.xz", "test")
}

func TestUrlToFilename(t *testing.T) {
	if FilenameFromURL("http://api.plos.org/search?q=...") != "search" {
		t.Errorf("Have %s want 'search'", FilenameFromURL("http://api.plos.org/search?q=..."))
	}
}
