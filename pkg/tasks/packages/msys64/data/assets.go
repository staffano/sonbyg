// +build dev

// This file will be used when -tags dev is used with go command.
// It means that data will actually be fetched from the "files" folder
// during development, but the statically linked assets will be used
// in all other scenarios.
// Don't forget to do 'go generate ./...' before releasing, so that the
// statically linked assets are up to date!

package data

//go:generate go run -tags=dev ../assets_generate.go

import "net/http"

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("assets")
