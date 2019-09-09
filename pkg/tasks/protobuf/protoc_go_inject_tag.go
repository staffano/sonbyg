package protobuf

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/gogo"
	"github.com/staffano/sonbyg/pkg/utils"
)

func doGenTags(v utils.Variables) *builder.Artifact {
	v.Printf("Start doGenTags")
	prefixBin := path.Join(v["PREFIX"].(string), "bin")
	exe, _ := filepath.Abs(path.Join(prefixBin, "protoc-go-inject-tag.exe"))
	tmpdir := os.TempDir()
	arg := fmt.Sprintf("--input=%s", v["FILE"].(string))
	v["PATH_PREPEND"] = prefixBin
	utils.Exec("protoc-go-inject-tag", tmpdir, v, exe, arg)
	v.Printf("End doGenTags")
	return nil
}

// GenTags runs the protoc-go-inject-tag.
// cf. https://github.com/favadi/protoc-go-inject-tag
func GenTags(b *builder.Builder, vars *utils.Variables, file string) *builder.Task {
	t := builder.NewTask(b, "GenTags")
	t.RunAlways = true
	v := vars.Copy("WORKSPACE", "PREFIX", "VERBOSE", "DOWNLOAD_DIR", "*PATH")
	v["FILE"] = file
	t.Variables = v
	//t.Variables.Dump()
	t.DependsOn(gogo.Get(b, &t.Variables, "github.com/favadi/protoc-go-inject-tag"))
	t.AssignDefaultSignature()
	v.Printf("Created GenTags Task")
	return b.Add(t, doGenTags)
}
