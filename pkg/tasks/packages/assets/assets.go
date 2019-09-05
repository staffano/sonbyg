package assets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/utils"
)

// createAssetSignature creates a md5 hash based on the name of the
// asset and the filenames in the asset.
func createAssetSignature(name string, assets http.FileSystem) string {
	h := md5.New()
	io.WriteString(h, name)
	walkFn := func(p string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", p, err)
			return nil
		}
		fmt.Println(p)
		io.WriteString(h, fi.Name())
		return nil
	}
	err := vfsutil.WalkFiles(assets, "/", walkFn)
	if err != nil {
		panic(err)
	}
	sig := hex.EncodeToString(h.Sum(nil))
	return sig
}

// Unpack will create a task that unpacks the http.FileSystem to the dst folder
// and use the contents of this asset as base for creating a signature of the
// task. As such, it does not use the Variables method of creating
// a signature. The asset is transferred to the task using the closure
// of the task function.
// Naming the asset provides an additional source of uniqueness to the
// signature.
func Unpack(b *builder.Builder, name string, fs http.FileSystem, dst string) *builder.Task {
	t := builder.NewTask(b, fmt.Sprintf("UnpackAssets/%s", name))
	t.Variables = utils.Variables{}
	t.Signature = createAssetSignature(name, fs)
	return b.Add(t, func(v utils.Variables) *builder.Artifact {
		utils.UnpackAssets(dst, fs)
		return nil
	})
}
