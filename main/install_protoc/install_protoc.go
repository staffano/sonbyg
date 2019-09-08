package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/staffano/sonbyg/pkg/utils"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/protobuf"
	//"github.com/staffano/sonbyg/pkg/tasks/packages/gurka"
)

func main() {
	b := builder.NewBuilder(context.Background(), "", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	var err error
	b.Variables["PREFIX"], err = filepath.Abs("../../dependencies")
	if err != nil {
		log.Fatalf("Could not determine PREFIX path")
	}
	utils.CreateDir(b.Variables["PREFIX"].(string), "", 755)

	ga := protobuf.ImportGoogleAPIs(b, &b.Variables)
	task := protobuf.Protoc(b, &b.Variables, ".", "--go_out=output", "test.proto")
	task.DependsOn(ga)
	b.DumpTasks()
	b.Build(task)
}
