// +build windows

package base

import (
	"path"
	"path/filepath"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Git will install git and expose the git command
func Git(b *builder.Builder, vars *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "git")
	v := vars.Copy("WORKSPACE", "DOWNLOAD_DIR", "VERBOSE", "*PATH")

	srcs := Sources{}
	pkgName := "${PN}-${PV}"
	srcs.Add("https://github.com/git-for-windows/git/releases/download/v${PV}.windows.1/MinGit-${PV}-busybox-64-bit.zip", "26076d95a90fdf5145d72d41ae396f24")
	v["SOURCES"] = srcs
	v["PV"] = "2.23.0"
	v["PN"] = "git"
	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", pkgName)

	v.ResolveAll()
	t.Variables = v
	t.AssignDefaultSignature()

	b.EstablishPath(744, v["PACKAGE_DIR"].(string))

	download := DownloadFile(b, &t.Variables)
	unpack := Unpack(b, &t.Variables)

	unpack.DependsOn(download)
	t.DependsOn(unpack)
	gitExec, _ := filepath.Abs(path.Join(v["PACKAGE_DIR"].(string), "cmd", "git.exe"))
	(*vars)["GIT"] = gitExec
	v.Printf("Created Git Task")
	return b.Add(t, func(v utils.Variables) *builder.Artifact {
		v.Printf("Git set up Done")
		return nil
	})
}
