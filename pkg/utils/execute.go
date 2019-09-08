package utils

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

// Exec runs a command and stores the stdout/stderr to
// a file named <outDir>/<id>.<nr>, where nr is a number
// that makes the filename unique within the directory.
func Exec(id, outDir string, v Variables, cmd string, args ...string) {
	var err error
	op := GetUniqueFilename(id+".out", outDir)
	log.Printf("Writing output to %s", op)
	// outfile, err := os.Create(op)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer outfile.Close()
	c := exec.Command(cmd, args...)
	if cwd, ok := v["CWD"]; ok {
		c.Dir = cwd.(string)
	}
	pth, ok := v["PATH_PREPEND"].(string)
	if ok {
		tv := NewVariables(os.Environ())
		if runtime.GOOS == "windows" {
			tv.PrependEnv("Path", pth)
		} else {
			tv.PrependEnv("PATH", pth)
		}
		c.Env = v.ExportEnv()
	}
	// w := bufio.NewWriter(outfile)
	// defer w.Flush()
	// mw := io.MultiWriter(w, os.Stdout)
	// v.Dump()
	c.Stderr = os.Stdout
	c.Stdout = os.Stdout
	log.Printf("exec:\n  CWD=%s\n  CMD=%s\n  ARGS=%v\n  OUTDIR=%s", c.Dir, c.Path, c.Args, outDir)
	if err = c.Run(); err != nil {
		log.Fatalf("Error running command %s :%v", cmd, err)
	}
}
