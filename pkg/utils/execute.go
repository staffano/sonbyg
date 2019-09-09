package utils

import (
	"fmt"
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
	v.Printf("Writing output to %s", op)
	// outfile, err := os.Create(op)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer outfile.Close()
	c := exec.Command(cmd, args...)

	if cwd, ok := v["CWD"]; ok {
		c.Dir = cwd.(string)
	}

	// Collect variables we want to export to the process
	// Start with the original environment
	tv := NewVariables(os.Environ())

	// Do we want to prepend the PATH variable?
	pth, ok := v["PATH_PREPEND"].(string)
	if ok {
		if runtime.GOOS == "windows" {
			tv.PrependEnv("Path", pth)
		} else {
			tv.PrependEnv("PATH", pth)
		}
	}

	// Do we have any explicit exports?
	vcol, ok := v["EXPORTED_VARS"].([]string)
	if ok {
		for _, vc := range vcol {
			if val, ok := v[vc]; ok {
				tv[vc] = fmt.Sprintf("%v", val)
			}
		}
	}
	if len(tv) > 0 {
		c.Env = tv.ExportEnv()
	}
	// w := bufio.NewWriter(outfile)
	// defer w.Flush()
	// mw := io.MultiWriter(w, os.Stdout)
	// for _, cv := range c.Env {
	// 	log.Print(cv)
	// }
	c.Stderr = os.Stdout
	c.Stdout = os.Stdout
	v.Printf("exec:\n  CWD=%s\n  CMD=%s\n  ARGS=%v\n  OUTDIR=%s", c.Dir, c.Path, c.Args, outDir)
	if err = c.Run(); err != nil {
		log.Fatalf("Error running command %s :%v", cmd, err)
	}
}
