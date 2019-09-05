package msys64

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages"
	"github.com/staffano/sonbyg/pkg/tasks/packages/assets"
	"github.com/staffano/sonbyg/pkg/tasks/packages/diff"
	msys_assets "github.com/staffano/sonbyg/pkg/tasks/packages/msys64/data"
	"github.com/staffano/sonbyg/pkg/utils"
)

func doBash(v utils.Variables) *builder.Artifact {
	log.Println("Install BASH...")
	return nil
}

// ExecBashScript will create a script file at
// <outdir>/<id>.run.<nr> and execute that script, sending
// the output to <outdir>/<id>.out.<nr>
func doExecBashScript(v utils.Variables) *builder.Artifact {
	id := v["_TASK_ID"].(string)
	script := v["SCRIPT"].(string)
	args := v["ARGS"].([]string)
	tmpDir := v["TMP_DIR"].(string)
	scriptFilePath := utils.GetUniqueFilename(id+".run", tmpDir)

	// Write script to file. Can be good to be able to check what's actually
	// been run.
	err := ioutil.WriteFile(scriptFilePath, []byte(script), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create argument list
	a := []string{}
	a = append(a, scriptFilePath)
	a = append(a, args...)
	cmd := v["BASH"].(string)

	utils.Exec("bash_exec-"+id, tmpDir, v, cmd, a...)

	return nil
}

// ExecuteBash runs a bash script
func ExecuteBash(b *builder.Builder, v *utils.Variables, script string, args ...string) *builder.Task {
	t := builder.NewTask(b, "ExecuteBash")
	t.Variables = v.Copy("WORKSPACE", "TARGET", "HOST", "BUILD", "DOWNLOAD_DIR", "PREFIX", "VERBOSE", "PATH", "TMP_DIR", "CWD")
	t.Variables["SCRIPT"] = script
	t.Variables["ARGS"] = args
	t.Variables.ResolveAll()
	t.DependsOn(Bash(b, &t.Variables))
	t.AssignDefaultSignature()
	v.Printf("Created ExecuteBash Task")
	return b.Add(t, doExecBashScript)
}

// Bash installs bash and exports the path to the bash cmd
func Bash(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "install-bash")
	t.Variables = v.Copy("WORKSPACE", "TARGET", "HOST", "BUILD", "DOWNLOAD_DIR", "PREFIX", "VERBOSE", "PATH", "TMP_DIR")
	t.DependsOn(BuildSystem(b, &t.Variables))
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	t = b.Add(t, doBash)
	(*v)["BASH"] = fmt.Sprintf("%s/bash.exe", t.Variables.Get("BIN_DIR"))
	v.Printf("Created BASH Task")
	return t
}

func doUpgradeBuildSystem(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doUpgradeBuildSystem ")
	nv := utils.Variables{}
	bashExec := path.Join(v.Get("PACKAGE_DIR"), "msys64/usr/bin/bash.exe")
	// cf. https://chocolatey.org/packages/msys2
	success := false
	for i := 1; i < 5; i++ {
		utils.Exec("UpgradeBuildSystem-"+string(i), v["TMP_DIR"].(string), nv, bashExec, "-l", "-c", "pacman --noconfirm -Syuu | tee /update.log")
		dat, err := ioutil.ReadFile(path.Join(v["SRC_DIR"].(string), "update.log"))
		if err != nil {
			log.Fatalf("Could not check installation status of Msys")
		}
		if strings.Count(string(dat), "there is nothing to do") == 2 {
			v.Printf("Msys2 base updated...")
			success = true
			break
		}
		if strings.Count(string(dat), "Inget behöver göras.") == 2 {
			v.Printf("Msys2 base updated")
			success = true
			break
		}
	}
	if !success {
		log.Fatalf("Could not update Msys2 base")
	}
	v.Printf("End doUpgradeBuildSystem ")
	return nil
}

// UpgradeBuildSystem uses pacman to upgrade the buildsystem
func UpgradeBuildSystem(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "upgrade-build-system")
	t.Variables = v.Copy("WORKSPACE", "PACKAGE_DIR", "DOWNLOAD_DIR", "PREFIX", "VERBOSE", "PATH", "TMP_DIR", "SRC_DIR")
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created UpgradeBuildSystem Task")
	return b.Add(t, doUpgradeBuildSystem)
}

func doRunCmd(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doRunCmd ")
	bashExec := path.Join(v.Get("PACKAGE_DIR"), "msys64/usr/bin/bash.exe")
	bsCmd := v["BUILD_SYSTEM_CMD"].(string)
	utils.Exec(strings.Replace(bsCmd, " ", "", -1), v["TMP_DIR"].(string), v, bashExec,"-c", "bash -c " + bsCmd)
	v.Printf("End doRunCmd")
	return nil
}

// RunCmd will execute a command using "bash -lc" inside the build system
// Example RunCmd(b, &v ,"pacman --noconfirm -S --needed base-devel mingw-w64-x86_64-toolchain")
func RunCmd(b *builder.Builder, vars *utils.Variables, cwd, cmd string) *builder.Task {
	t := builder.NewTask(b, "RunCmd")
	v := vars.Copy("WORKSPACE", "TMP_DIR", "PATH")
	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", "msys64")
	v["BUILD_SYSTEM_CMD"] = cmd
	v.ResolveAll()

	t.Variables = v
	t.AssignDefaultSignature()

	return b.Add(t, doRunCmd)
}

// RunAlwaysCmd will always execute a command using "bash -c" inside the build system
// Example RunAllwaysCmd(b, &v ,"pacman --noconfirm -S --needed base-devel mingw-w64-x86_64-toolchain")
func RunAlwaysCmd(b *builder.Builder, vars *utils.Variables, cwd, cmd string) *builder.Task {
	t := builder.NewTask(b, "RunCmd")
	v := vars.Copy("WORKSPACE", "TMP_DIR", "PATH")
	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", "msys64")
	v["BUILD_SYSTEM_CMD"] = cmd
	v.ResolveAll()

	t.Variables = v
	t.RunAlwaysSignature()
	return b.Add(t, doRunCmd)
}

// BuildSystem installs a gcc/make/autoconf system that can be used to
// build packages
func BuildSystem(b *builder.Builder, vars *utils.Variables) *builder.Task {
	pkgName := "msys64"
	t := builder.NewTask(b, "install-build-system")
	v := vars.Copy("WORKSPACE", "TARGET", "HOST", "BUILD", "DOWNLOAD_DIR", "PREFIX", "VERBOSE", "PATH")
	srcs := packages.Sources{}
	srcs.Add("http://repo.msys2.org/distrib/x86_64/msys2-base-x86_64-20190524.tar.xz", "b9fddc5a8ea27d5f0eed232795e99725")

	v["PACKAGE_DIR"] = path.Join("${WORKSPACE}", pkgName)
	v["SRC_DIR"] = path.Join("${PACKAGE_DIR}", pkgName)
	v["TMP_DIR"] = path.Join("${PACKAGE_DIR}", "tmp")
	v["ASSETS_DIR"] = path.Join("${PACKAGE_DIR}", "assets")
	v["SOURCES"] = srcs
	v.ResolveAll()
	t.Variables = v
	(*vars)["BIN_DIR"] = b.WorkspacePath(t.Variables.Get("PACKAGE_DIR"), "msys64", "usr", "bin")

	mirrorsPatch := path.Join(v["ASSETS_DIR"].(string), "001_install.patch")

	b.EstablishPath(755, t.Variables["PACKAGE_DIR"].(string))
	b.EstablishPath(755, t.Variables["TMP_DIR"].(string))

	assts := assets.Unpack(b, "msys_assets", msys_assets.Assets, t.Variables["ASSETS_DIR"].(string))
	patch := diff.Patch(b, &t.Variables, t.Variables["SRC_DIR"].(string), mirrorsPatch)
	dwnld := packages.Download(b, &t.Variables)
	unpck := packages.Unpack(b, &t.Variables)
	installBase := RunCmd(b, &t.Variables, ".", "pacman -S base --force --noconfirm")
	upgrd := UpgradeBuildSystem(b, &t.Variables)
	instl := RunCmd(b, &t.Variables, ".", "pacman --noconfirm -S --needed base-devel mingw-w64-x86_64-toolchain")

	unpck.DependsOn(dwnld)
	patch.DependsOn(assts, unpck)
	installBase.DependsOn(patch)
	upgrd.DependsOn(installBase)
	instl.DependsOn(upgrd)
	t.DependsOn(instl)
	t.AssignDefaultSignature()

	(*vars)["MINGW64_BIN_DIR"] = path.Join(v["SRC_DIR"].(string), "mingw64\\bin")
	(*vars)["USR_BIN_DIR"] = path.Join(v["SRC_DIR"].(string), "usr\\bin")
	(*vars)["MSYS_ROOT"] = v["SRC_DIR"]

	v.Printf("Created BuildSystem Task")

	return b.Add(t, builder.Message(v, "Build system installed"))
}
