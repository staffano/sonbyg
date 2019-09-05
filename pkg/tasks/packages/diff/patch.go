package diff

import (
	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages"
	"github.com/staffano/sonbyg/pkg/tasks/packages/git"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Patch will patch the directory with a patch
func Patch(b *builder.Builder, vars *utils.Variables, dir, patch string) *builder.Task {
	v := vars.Copy("WORKSPACE", "DOWNLOAD_DIR", "VERBOSE", "PATH", "TMP_DIR")
	git := git.Git(b, &v)
	v["CWD"] = dir
	patch = v.Resolve(patch)
	v.ResolveAll()
	pp := packages.Execute(b, &v, v["GIT"].(string), "apply", "--unsafe-paths", patch)
	pp.DependsOn(git)
	return pp
}
