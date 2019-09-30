package protobuf

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/base"
	"github.com/staffano/sonbyg/pkg/tasks/packages/gogo"
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
	v["SOURCES"] = base.NewSource("https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-win64.zip", "535d2d23004f90067ee1d3cb2599d0e5")
	v.ResolveAll()
	t.Variables = v
	t.DependsOnSerial(
		base.Download(b, &v),
		base.Unpack(b, &v),
		base.Install(b, &v, v["PACKAGE_DIR"].(string), v["PREFIX"].(string), 755),
	)
	t.AssignDefaultSignature()
	return b.Add(t, builder.Message(t.Variables, "Protoc compiler installed"))
}

func doProtoc(v utils.Variables) *builder.Artifact {
	v.Printf("Start doProtoc")
	args := v["ARGS"].([]string)
	includes := v["PROTOC_OPTS"].([]string)
	prefixBin := path.Join(v["PREFIX"].(string), "bin")
	protocExe, _ := filepath.Abs(path.Join(prefixBin, "protoc.exe"))
	tmpdir := os.TempDir()
	opts := append(includes, args...)
	v["PATH_PREPEND"] = prefixBin
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
	repoDir := "${PACKAGE_DIR}/" + pkgName
	googleDir := v.Resolve(path.Join("${PACKAGE_DIR}", pkgName, "google"))
	installDir := v.Resolve("${PREFIX}/include/google")
	v["SOURCES"] = base.NewGit("https://github.com/googleapis/googleapis.git?ref=master", repoDir)
	v.ResolveAll()
	t.Variables = v
	t.DependsOnSerial(
		base.Download(b, &v),
		base.Install(b, &v, googleDir, installDir, 755),
	)
	t.AssignDefaultSignature()
	b.EstablishPath(744, v["PACKAGE_DIR"].(string))
	return b.Add(t, builder.Message(t.Variables, "Protoc compiler installed"))
}

func appendOpts(v *utils.Variables, opts ...string) {
	if o, ok := (*v)["PROTOC_OPTS"].([]string); !ok {
		(*v)["PROTOC_OPTS"] = opts
	} else {
		(*v)["PROTOC_OPTS"] = append(opts, o...)
	}
}

// ImportGoogleAPIs updates the include path for protoc to make the google apis
// available
func ImportGoogleAPIs(b *builder.Builder, vars *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "import-google-apis")
	t.Variables = vars.Copy("PREFIX")
	t.Variables.ResolveAll()
	includeDir := t.Variables.Resolve("${PREFIX}/include/google")
	t.DependsOn(GoogleAPIs(b, vars))
	t.RunAlwaysSignature()
	appendOpts(vars, "-I"+includeDir)
	return (b.Add(t, builder.Message(t.Variables, "Google APIs imported.")))
}

// Protoc runs the protoc compiler on a source.
// THe include directories will be resolved at runtime and use
// the includes that has been set by any dependencies.
func Protoc(b *builder.Builder, vars *utils.Variables, cwd string, args ...string) *builder.Task {
	t := builder.NewTask(b, "Protoc")
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR", "*PROTOC_OPTS", "*PATH")
	v["ARGS"] = args
	var err error
	v["CWD"], err = filepath.Abs(cwd)
	if err != nil {
		log.Fatalf("Could not determine CWD, %v", err)
	}
	prefixIncl := path.Join(v["PREFIX"].(string), "include")
	appendOpts(&v, "-I"+v["CWD"].(string), "-I"+prefixIncl)
	t.Variables = v
	t.DependsOn(gogo.Get(b, &t.Variables, "github.com/golang/protobuf/protoc-gen-go"))
	t.DependsOn(InstallProtoc(b, &t.Variables))
	t.AssignDefaultSignature()
	t.RunAlways = true
	v.Printf("Created Protoc Task")
	return b.Add(t, doProtoc)
}
