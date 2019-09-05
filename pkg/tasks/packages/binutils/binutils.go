package binutils

import (
	"path"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages"
	"github.com/staffano/sonbyg/pkg/tasks/packages/msys64"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Binutils demostrates the concept of a package. It creates a set
// of tasks needed in order to fulfill the package. This Package allows
// us to specify a specific version of the packate to be installerd.
func Binutils(b *builder.Builder, vars *utils.Variables, version string) *builder.Task {
	vars.Printf("Binutils Begin Create")
	t := builder.NewTask(b, "Binutils")
	srcs := packages.Sources{}
	pkgName := "${PN}-${PV}"
	srcs.Add("http://ftp.acc.umu.se/mirror/gnu.org/gnu/binutils/binutils-${PV}.tar.bz2", "84edf97113788f106d6ee027f10b046a")

	v := vars.Copy("WORKSPACE", "DOWNLOAD_DIR", "TARGET", "HOST", "BUILD", "SYSROOT_DIR", "PREFIX", "VERBOSE", "PATH")

	v["SOURCES"] = srcs
	v["PV"] = version
	v["PN"] = "binutils"
	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", pkgName)
	v["SRC_DIR"] = path.Join("${PACKAGE_DIR}", pkgName)
	v["TMP_DIR"] = path.Join("${PACKAGE_DIR}", "tmp")
	v["CONFIG_OPTS"] = "--prefix=${PREFIX} " +
		"--program-prefix=${TARGET}- " +
		"--target=${TARGET} " +
		"--host=${HOST} " +
		"--build=${BUILD} " +
		"--with-sysroot=${SYSROOT_DIR} "

	v["MAKE_SWITCHES"] = "-j8"

	v.ResolveAll()
	t.Variables = v
	t.AssignDefaultSignature()

	b.EstablishPath(744, v["PACKAGE_DIR"].(string))
	b.EstablishPath(744, v["TMP_DIR"].(string))

	download := packages.Download(b, &t.Variables)
	unpack := packages.Unpack(b, &t.Variables)
	configure := msys64.Configure(b, &t.Variables)
	compile := msys64.Make(b, &t.Variables, "all")
	install := msys64.Make(b, &t.Variables, "install")

	install.DependsOn(compile)
	compile.DependsOn(configure)
	configure.DependsOn(unpack)
	unpack.DependsOn(download)
	t.DependsOn(install)

	v.Printf("Created Binutils Task")
	return b.Add(t, builder.Message(v, "Binutils Done"))
}
