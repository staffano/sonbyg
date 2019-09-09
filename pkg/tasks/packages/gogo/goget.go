package gogo

import (
	"path"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/base"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Get calls go get and installs to ${PREFIX}/bin
func Get(b *builder.Builder, vars *utils.Variables, mod string) *builder.Task {
	t := builder.NewTask(b, "InstallProtocGenGo")
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR")
	v["GOBIN"] = path.Join("${PREFIX}", "bin")
	v["EXPORTED_VARS"] = []string{"GOBIN"}
	v.ResolveAll()
	t.Variables = v
	t.AssignDefaultSignature()
	return base.Execute(b, &t.Variables, "go", "get", "-u", mod)
}
