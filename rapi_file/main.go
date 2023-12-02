package rapi_file

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_hash"
	"github.com/maldan/go-cmhp/cmhp_process"
	"github.com/maldan/go-rapi/rapi_debug"
	"github.com/maldan/go-rapi/rapi_error"
	"net/http"
	"os"
	"strings"

	"github.com/maldan/go-rapi/rapi_core"
)

type FileHandler struct {
	Root string
}

func (r FileHandler) Handle(args rapi_core.HandlerArgs) {
	// Handle panic
	defer rapi_core.HandleError(&args)

	cwd, _ := os.Getwd()

	// Pure path without route // example /data/test -> /test
	routePath := strings.Replace(args.R.URL.Path, args.Route, "", 1)

	path := strings.ReplaceAll(r.Root, "@", cwd) + routePath
	path = strings.ReplaceAll(path, "\\", "/")

	rapi_core.DisableCors(args.RW)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		rapi_error.Fatal(rapi_error.Error{Code: 404, Description: "File not found"})
	}

	// Create dir for converting images
	if !cmhp_file.Exists(os.TempDir() + "/rapi_convert") {
		cmhp_file.Mkdir(os.TempDir() + "/rapi_convert")
	}

	// Get mime of file
	mtype, _ := mimetype.DetectFile(path)
	uniqueName := cmhp_hash.Sha1(fmt.Sprintf("%v/%v", path, args.R.URL.Query()))

	if mtype.String() == "image/png" {
		convertTo := args.R.URL.Query().Get("convertTo")
		if convertTo == "webp" {
			tmpFile := os.TempDir() + "/rapi_convert/" + uniqueName + ".webp"

			// Convert
			if !cmhp_file.Exists(tmpFile) {
				cmhp_process.Exec("magick", path, tmpFile)
				fmt.Printf("CONVERTED: %v\n", tmpFile)
			}
			path = tmpFile
		}
	}

	http.ServeFile(args.RW, args.R, path)

	if args.DebugMode {
		rapi_debug.GetRequestLog(args.Id).SetResponse("SERVE FILE: " + path)
		// rapi_debug.GetRequestLog(args.Id).SetArgs(args.MethodArgs)
	}
}
