package diff

import (
	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/base"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Patch will patch the directory with a patch
func Patch(b *builder.Builder, vars *utils.Variables, dir, patch string) *builder.Task {
	v := vars.Copy("WORKSPACE", "DOWNLOAD_DIR", "VERBOSE", "PATH", "TMP_DIR")
	git := base.Git(b, &v)
	v["CWD"] = dir
	patch = v.Resolve(patch)
	v.ResolveAll()
	pp := base.Execute(b, &v, v["GIT"].(string), "apply", "--unsafe-paths", patch)
	pp.DependsOn(git)
	return pp
}
