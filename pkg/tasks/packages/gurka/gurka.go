package gurka

import (
	"path"

	"github.com/staffano/sonbyg/pkg/tasks/packages/msys64"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Gurka is a test task
func Gurka(b *builder.Builder, vars *utils.Variables) *builder.Task {
	pkgName := "gurka"
	t := builder.NewTask(b, "gurka")
	t.Variables = vars.Copy("WORKSPACE", "TARGET", "HOST", "BUILD", "DOWNLOAD_DIR", "PREFIX", "VERBOSE", "PATH")
	t.Variables["PACKAGE_DIR"] = path.Join("${WORKSPACE}", pkgName)
	t.Variables["TMP_DIR"] = path.Join("${PACKAGE_DIR}", "tmp")
	t.Variables.ResolveAll()
	t.DependsOn(msys64.BuildSystem(b, &t.Variables))
	t.Variables.PrependEnv("PATH", t.Variables["USR_BIN_DIR"].(string))
	t.Variables.PrependEnv("PATH", t.Variables["MINGW64_BIN_DIR"].(string))

	//t.DependsOn(msys64.RunAlwaysCmd(b, &t.Variables, ".", "env | grep PATH"))
	t.DependsOn(msys64.RunAlwaysCmd(b, &t.Variables, ".", "env | grep PATH"))
	t.DependsOn(msys64.RunAlwaysCmd(b, &t.Variables, ".", "gcc"))
	
	t.RunAlwaysSignature()
	b.EstablishPath(755, t.Variables["TMP_DIR"].(string))
	t.Variables.Printf("Created Gurka Task")
	return b.Add(t, func(v utils.Variables) *builder.Artifact {
		v.Printf("Patch Done")
		return nil
	})
}
