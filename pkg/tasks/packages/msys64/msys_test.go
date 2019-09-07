package msys64

import (
	"context"
	"path"
	"testing"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/base"
)

func Test_Download(t *testing.T) {

	b := builder.NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	srcs := base.Sources{}
	srcs.Add("https://ftp.acc.umu.se/mirror/gnu.org/gnu/binutils/binutils-2.32.tar.xz", "0d174cdaf85721c5723bf52355be41e6")
	b.Variables["PACKAGE_DIR"] = path.Join(b.Variables.Get("WORKSPACE"), "build_system")
	b.Variables["SOURCES"] = srcs
	b.Variables["PACKAGE_DIR"] = ` \
    --prefix=${INSTALL_DIR} \
    --program-prefix=${TARGET}- \
    --target=${TARGET} \
    --host=${HOST} \
    --build=${BUILD} \
    --with-sysroot=${SYSROOT_DIR} \
`
	download := base.Download(b, &b.Variables)
	unpack := base.Unpack(b, &b.Variables)
	configure := Configure(b, &b.Variables)
	unpack.DependsOn(download)

	b.Build(configure)
}
