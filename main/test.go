package main

import (
	"context"
	"path"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/binutils"
	//"github.com/staffano/sonbyg/pkg/tasks/packages/gurka"
)

func main() {
	b := builder.NewBuilder(context.Background(), "../test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	b.Variables["SYSROOT_DIR"] = path.Join(b.Variables.Get("WORKSPACE"), "sysroot")
	b.Variables["PREFIX"] = path.Join(b.Variables["WORKSPACE"].(string), "install_dir")
	b.EstablishPath(755, b.Variables["SYSROOT_DIR"].(string))
	b.EstablishPath(755, b.Variables["PREFIX"].(string))
	task := binutils.Binutils(b, &b.Variables, "2.31.1")
	//task := gurka.Gurka(b, &b.Variables)

	b.DumpTasks()
	b.Build(task)
}
