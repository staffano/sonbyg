package protobuf

import (
	"fmt"
	"os"
	"path"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Note that the grpc support will automatically be included in the host application
// once the "google.golang.org/grpc" import is used, so no extra dependency handling
// is required for gRPC

// InstallProtoc installs the protoc compiler.
// VARS: WORKSPACE, PREFIX,
func InstallProtoc(b *builder.Builder, vars *utils.Variables) *builder.Task {
	pkgName := "protoc-compiler"
	t := builder.NewTask(b, "InstallProtoc")
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR")
	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", pkgName)
	v["SOURCES"] = packages.NewSource("https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-win64.zip", "535d2d23004f90067ee1d3cb2599d0e5")
	v.ResolveAll()
	t.Variables = v
	t.DependsOn(packages.Download(b, &v))
	t.DependsOn(packages.Unpack(b, &v))
	t.DependsOn(packages.Install(b, &v, v["PACKAGE_DIR"].(string), v["PREFIX"].(string), 755))
	t.AssignDefaultSignature()
	return b.Add(t, builder.Message(t.Variables, "Protoc compiler installed"))
}

// InstallProtocGenGo installs the protoc-gen-go tool.
// Since we want the plugin protoc-gen-go to be in the PREFIX/bin directory
// we will specifically set the $GOBIN variable to ${PREFIX}/bin before
// calling "go get -u github.com/golang/protobuf/protoc-gen-go"
func InstallProtocGenGo(b *builder.Builder, vars *utils.Variables) *builder.Task {

	t := builder.NewTask(b, "InstallProtocGenGo")
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR")
	v["GOBIN"] = path.Join("${PREFIX}", "bin")
	v.ResolveAll()

	protocGenGoURL := "github.com/golang/protobuf/protoc-gen-go"
	t.Variables = v
	e := packages.Execute(b, &t.Variables, "go", "get", "-u", protocGenGoURL)

	t.DependsOn(InstallProtoc(b, &t.Variables))
	t.DependsOn(e)

	t.AssignDefaultSignature()
	return b.Add(t, builder.Message(t.Variables, fmt.Sprintf("protoc-gen-go installed to %v", v["GOBIN"])))
}

func doProtoc(v utils.Variables) *builder.Artifact {
	v.Printf("Start doProtoc")
	args := v["ARGS"].([]string)
	includes := v["PROTOC_OPTS"].([]string)
	prefixBin := path.Join(v["PREFIX"].(string), "bin")
	protocExe := path.Join(prefixBin, "protoc.exe")
	tmpdir := os.TempDir()

	// Build option list to protoc
	opts := []string{}
	for _, i := range includes {
		opts = append(opts, "-I"+i)
	}
	opts = append(opts, args...)
	v.PrependEnv("PATH", prefixBin)
	utils.Exec("protoc", tmpdir, v, protocExe, opts...)
	v.Printf("End doProtoc")
	return nil
}

// GoogleAPIs makes the google protobuf APIs available to the build system
func GoogleAPIs(b *builder.Builder, vars *utils.Variables) *builder.Task {
	pkgName := "protobuf-google-apis"
	t := builder.NewTask(b, pkgName)
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR")
	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", pkgName)
	// TODO, we need git support again... go-gitter
	v["SOURCES"] = packages.NewGit("github.com/googleapis/googleapis.git")
	v.ResolveAll()
	t.Variables = v
	srcDir := v.Resolve("${PACKAGE_DIR}/googleapis-master")
	t.DependsOn(packages.Download(b, &v))
	t.DependsOn(packages.Install(b, &v, srcDir, v["PREFIX"].(string), 755))
	t.AssignDefaultSignature()
	return b.Add(t, builder.Message(t.Variables, "Protoc compiler installed"))
}

func appendOpts(v *utils.Variables, opts ...string) {
	if o, ok := (*v)["PROTOC_OPTS"].([]string); !ok {
		(*v)["PROTOC_OPTS"] = opts
	} else {
		(*v)["PROTOC_OPTS"] = append(o, opts...)
	}
}

// ImportGoogleAPIs updates the include path for protoc to make the google apis
// available
func ImportGoogleAPIs(b *builder.Builder, vars *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "import-google-apis")
	v := vars.Copy("PREFIX")
	v.ResolveAll()
	includeDir := v.Resolve("${PREFIX}/googleapis")
	t.DependsOn(GoogleAPIs(b, vars))
	t.RunAlwaysSignature()
	appendOpts(vars, "-I"+includeDir)
	return (b.Add(t, builder.Message(t.Variables, "Google APIs imported.")))
}

// Protoc runs the protoc compiler on a source.
// THe include directories will be resolved at runtime and use
// the includes that has been set by any dependencies.
func Protoc(b *builder.Builder, vars *utils.Variables, args ...string) *builder.Task {
	t := builder.NewTask(b, "Protoc")
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR", "*PROTOC_OPTS")
	v["ARGS"] = args
	prefixIncl := path.Join(v["PREFIX"].(string), "include")
	appendOpts(&v, "-I.", "-I"+prefixIncl)
	t.Variables = v
	t.DependsOn(InstallProtocGenGo(b, &t.Variables))
	t.AssignDefaultSignature()
	v.Printf("Created ExecuteBash Task")
	return b.Add(t, builder.Message(t.Variables, fmt.Sprintf("protoc installed")))
}
