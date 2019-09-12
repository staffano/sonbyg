package gogo

import (
	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/base"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Run calls go run
func Run(b *builder.Builder, vars *utils.Variables, mod string, args ...string) *builder.Task {
	t := builder.NewTask(b, "InstallProtocGenGo")
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "*CWD")
	//	v["GOBIN"] = path.Join("${PREFIX}", "bin")
	//	v["EXPORTED_VARS"] = []string{"GOBIN"}
	v.ResolveAll()
	t.Variables = v
	t.AssignDefaultSignature()
	as := []string{"run", mod}
	as = append(as, args...)
	return base.Execute(b, &t.Variables, "go", as...)
}
