package rapi_file

import (
	"errors"
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

	http.ServeFile(args.RW, args.R, path)

	if args.DebugMode {
		rapi_debug.GetRequestLog(args.Id).SetResponse("SERVE FILE: " + path)
		// rapi_debug.GetRequestLog(args.Id).SetArgs(args.MethodArgs)
	}
}
