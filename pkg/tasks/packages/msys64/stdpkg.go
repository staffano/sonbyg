package msys64

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/utils"
)

const (
	configureScriptTmpl = `
{{range $key, $value := .variables  }} 
export {{ $key }}="{{ $value }}"
{{end}}

function run {
	./configure {{ .variables.CONFIG_OPTS }}
}

cd {{ .variables.SRC_DIR }} && run
`
)

// Configure calls the './configure' script in the PACKAGE_DIR
func Configure(b *builder.Builder, v *utils.Variables) *builder.Task {
	vars := v.Copy("all")
	vars["CWD"] = vars["SRC_DIR"]
	vars.ResolveAll()
	// Convert back-slash in all variables to slash. Hackish...
	tv := utils.Variables{}
	for k, vv := range vars {
		if str, ok := vv.(string); ok {
			tv[k] = filepath.ToSlash(str)
		}
	}
	tv.PrependEnv("PATH", "/mingw64/bin:/usr/local/bin:/usr/bin:/bin:/opt/bin:/")
	tmpl := map[string]interface{}{
		"variables": tv,
		"script":    configureScriptTmpl,
	}
	t := template.Must(template.New("bash-script").Parse(configureScriptTmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, tmpl); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	log.Println(buf.String())
	return ExecuteBash(b, &vars, buf.String())
}

const (
	makeScriptTmpl = `
{{range $key, $value := .variables }} 
export {{ $key }}="{{ $value }}"
{{end}}

function run {
	make {{ .variables.MAKE_OPTIONS }} {{ .variables.MAKE_TARGET }}
}

cd {{ .variables.SRC_DIR }} && run
`
)

// Make creates a ExecuteBash task with a script that calls make in the PACKAGE_DIR
func Make(b *builder.Builder, v *utils.Variables, target string) *builder.Task {

	vars := v.Copy("WORKSPACE", "DOWNLOAD_DIR", "TARGET", "TMP_DIR", "HOST", "SRC_DIR", "BUILD", "SYSROOT_DIR", "PREFIX", "VERBOSE", "PATH")
	vars["MAKE_TARGET"] = target
	vars["MAKE_OPTS"] = "-j8"
	vars["CWD"] = vars["SRC_DIR"]
	vars.ResolveAll()

	// Convert back-slash in all variables to slash. Hackish...
	tv := utils.Variables{}
	for k, vv := range vars {
		if str, ok := vv.(string); ok {
			tv[k] = filepath.ToSlash(str)
		}
	}
	tv.PrependEnv("PATH", "/mingw64/bin:/usr/local/bin:/usr/bin:/bin:/opt/bin:/")
	tmpl := map[string]interface{}{
		"variables": tv,
		"script":    makeScriptTmpl,
	}
	t := template.Must(template.New("bash-script").Parse(makeScriptTmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, tmpl); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	return ExecuteBash(b, &vars, buf.String())
}
