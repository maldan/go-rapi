package rapi_vfs

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-rapi/rapi_core"
)

type VFSHandler struct {
	Root string
	Fs   embed.FS
}

func (r VFSHandler) Handle(args rapi_core.HandlerArgs) {
	// Handle panic
	defer rapi_core.HandleError(args.RW, args.R)

	// Prepare path
	p := strings.Replace(args.R.URL.Path, r.Root, "", 1)
	if p == "" || p == "/" {
		p = "/index.html"
	}

	// Read file
	data, err := r.Fs.ReadFile(r.Root + p)
	rapi_core.FatalIfError(err)

	// Write to temp dir
	p2 := os.TempDir() + "/rapi_vfs/" + fmt.Sprintf("%v", os.Getpid()) + "/" + p
	err = os.MkdirAll(filepath.Dir(p2), 0777)
	if err != nil {
		panic(err)
	}
	err = cmhp_file.Write(p2, data)
	if err != nil {
		panic(err)
	}

	// Serve file
	rapi_core.DisableCors(args.RW)
	http.ServeFile(args.RW, args.R, p2)
}
