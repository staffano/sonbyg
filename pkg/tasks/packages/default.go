package packages

import (
	"encoding/gob"
	"encoding/hex"
	"log"
	"path"
	"path/filepath"

	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/utils"
)

func init() {
	gob.Register(Sources{})
}

// A SrcSpec defines the location of a source file used by tasks and the checksum
// associated with the source
type SrcSpec struct {
	URI       string
	MD5       string
	LocalName string
}

// Sources is a collection of SrcSpec
type Sources []SrcSpec

// Add adds a source to Sources
func (s *Sources) Add(uri, md5 string) {
	*s = append(*s, SrcSpec{URI: uri, MD5: md5})
}

func doExecute(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doExecute")
	id := v["_TASK_ID"].(string)
	cmd := v["EXEC_FILE"].(string)
	args := v["ARGS"].([]string)
	tmpDir := v["TMP_DIR"].(string)
	utils.Exec(id, tmpDir, v, cmd, args...)
	v.Printf("End doExecute")
	return nil
}

// Execute executes the file
func Execute(b *builder.Builder, v *utils.Variables, executable string, args ...string) *builder.Task {
	t := builder.NewTask(b, "Execute")
	t.Variables = v.Copy("VERBOSE", "PATH", "TMP_DIR", "CWD")
	t.Variables["EXEC_FILE"] = executable
	t.Variables["ARGS"] = args
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created Execute Task")
	return b.Add(t, doExecute)
}

func doDownload(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doDownload")
	for _, s := range v["SOURCES"].(Sources) {
		uri := v.Resolve(s.URI)
		filePath := path.Join(v["DOWNLOAD_DIR"].(string), filepath.Base(uri))
		v.Printf("Downloading %s", uri)
		// Check if file exists
		if !utils.FileExists(filePath) {
			utils.DownloadFile(filePath, uri)
		}
		fileMD5 := utils.FileMD5(filePath)
		hexMd5 := hex.EncodeToString(fileMD5)
		if hexMd5 != s.MD5 {
			log.Fatalf("MD5 checksum [%s] of downloaded file doesn't match expected [%s] checksum", hexMd5, s.MD5)
		}
		log.Printf("%s downloaded", uri)
	}
	// We don't report an artifact for this
	v.Printf("End doDownload")
	return nil
}

// Download downloads all files specified in SOURCES
func Download(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "download_sources")
	t.Variables = v.Copy("DOWNLOAD_DIR", "SOURCES", "VERBOSE", "PATH")
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created Download Task")
	return b.Add(t, doDownload)
}

func doCreatePackageBuildDir(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doCreatePackageBuildDir")
	utils.CreateDir(v["PACKAGE_DIR"].(string), "", 744)
	log.Printf("Package build dir [%s] created", v["PACKAGE_DIR"].(string))
	v.Printf("End doCreatePackageBuildDir")
	return nil
}

// CreatePackageBuildDir creates a location in the WorkDir where the package can unpack
// and execute its tasks
func CreatePackageBuildDir(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "CreatePackageBuildDir")
	t.Variables = v.Copy("WORKSPACE", "PACKAGE_DIR", "VERBOSE", "PATH")
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created CreatePackageBuildDir Task")
	return b.Add(t, doCreatePackageBuildDir)
}

func doUnpack(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doUnpack")
	for _, s := range v["SOURCES"].(Sources) {
		uri := utils.FilenameFromURL(v.Resolve(s.URI))
		p := path.Join(v["DOWNLOAD_DIR"].(string), filepath.Base(uri))
		d := v.Get("PACKAGE_DIR")
		utils.UnpackFile(p, d)
		log.Printf("%s unpacked", p)
	}
	v.Printf("End doUnpack")
	return nil
}

// Unpack unpacks all files in the SOURCES variable into the build dir
func Unpack(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "Unpack-Package")
	t.Variables = v.Copy("SOURCES", "PACKAGE_DIR", "DOWNLOAD_DIR", "VERBOSE", "PATH")
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created Unpack Task")
	return b.Add(t, doUnpack)
}

// Patch applies all patches provided in the $BUILD_DIR/patches directory
func Patch(b *builder.Builder, v *utils.Variables) *builder.Task { return nil }
