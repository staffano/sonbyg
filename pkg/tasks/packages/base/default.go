package base

import (
	"encoding/gob"
	"encoding/hex"
	"log"
	"os"
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
	Protocol  string // file or git
}

// Sources is a collection of SrcSpec
type Sources []SrcSpec

// Add adds a source to Sources
func (s *Sources) Add(uri, md5 string) {
	*s = append(*s, SrcSpec{URI: uri, MD5: md5})
}

// NewSource is a helper function for easy creation of a single source
func NewSource(url, md5 string) Sources {
	return Sources{SrcSpec{URI: url, MD5: md5, Protocol: "file"}}
}

// NewGit ...
func NewGit(url string, localpath string) Sources {
	return Sources{SrcSpec{URI: url, Protocol: "git", LocalName: localpath}}
}
func doExecute(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doExecute")
	id := v["_TASK_ID"].(string)
	cmd := v["EXEC_FILE"].(string)
	args := v["ARGS"].([]string)
	tmpDir, ok := v["TMP_DIR"].(string)
	if !ok {
		tmpDir = os.TempDir()
	}
	utils.Exec(id, tmpDir, v, cmd, args...)
	v.Printf("End doExecute")
	return nil
}

// Execute executes the file
func Execute(b *builder.Builder, v *utils.Variables, executable string, args ...string) *builder.Task {
	t := builder.NewTask(b, "Execute")
	t.Variables = v.Copy("all")
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
		if s.Protocol == "file" || s.Protocol == "" {
			dst := path.Join(v["DOWNLOAD_DIR"].(string), filepath.Base(uri))
			v.Printf("Downloading %s", uri)
			utils.DownloadFile(dst, uri)
			fileMD5 := utils.FileMD5(dst)
			hexMd5 := hex.EncodeToString(fileMD5)
			if hexMd5 != s.MD5 {
				log.Fatalf("MD5 checksum [%s] of downloaded file doesn't match expected [%s] checksum", hexMd5, s.MD5)
			}
			log.Printf("%s downloaded", uri)
		} else if s.Protocol == "git" {
			dst := s.LocalName
			v.Printf("Downloading %s", uri)
			utils.DownloadGit(uri, dst)
			log.Printf("%s downloaded", uri)
		}

	}
	// We don't report an artifact for this
	v.Printf("End doDownload")
	return nil
}

// Download downloads all uris specified in SOURCES
func Download(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "download_sources")
	t.Variables = v.Copy("DOWNLOAD_DIR", "SOURCES", "*VERBOSE", "*PATH")
	t.Variables.ResolveAll()
	t.DependsOn(Git(b, v))
	t.AssignDefaultSignature()
	v.Printf("Created Download Task")
	return b.Add(t, doDownload)
}

// DownloadFile downloads all uris specified in SOURCES
func DownloadFile(b *builder.Builder, v *utils.Variables) *builder.Task {
	t := builder.NewTask(b, "download_sources")
	t.Variables = v.Copy("DOWNLOAD_DIR", "SOURCES", "*VERBOSE", "*PATH")
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
	t.Variables = v.Copy("SOURCES", "PACKAGE_DIR", "DOWNLOAD_DIR", "*VERBOSE", "*PATH")
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created Unpack Task")
	return b.Add(t, doUnpack)
}

func doInstall(v utils.Variables) *builder.Artifact {
	v.Printf("Begin doInstall")

	src := v["SRC_DIR"].(string)
	dst := v["DST_DIR"].(string)
	mode := v["INSTALL_MODE"].(os.FileMode)

	utils.Install(src, dst, mode)

	v.Printf("End doInstall")
	return nil
}

// Install will copy what's under src to dst and set it's mode. The behaviour depends
// on if dst ends with a '/' or not. 'dst/' will copy src to 'dst/src', while 'dst'
// will replace dst with src, potentially overwriting dst. Note that the dst and src
// will be evaluated at time of the build, not at time of configuration.
func Install(b *builder.Builder, v *utils.Variables, src, dst string, mode os.FileMode) *builder.Task {
	t := builder.NewTask(b, "Install")
	t.Variables = v.Copy("VERBOSE")
	t.Variables["SRC_DIR"] = src
	t.Variables["DST_DIR"] = dst
	t.Variables["INSTALL_MODE"] = mode
	t.Variables.ResolveAll()
	t.AssignDefaultSignature()
	v.Printf("Created Install Task %s -> %s (%d)", src, dst, mode)
	return b.Add(t, doInstall)
}
