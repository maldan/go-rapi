package rapi_file

import (
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
	defer rapi_core.HandleError(args.RW, args.R)

	cwd, _ := os.Getwd()

	// Pure path without route // example /data/test -> /test
	routePath := strings.Replace(args.R.URL.Path, args.Route, "", 1)
	
	path := strings.ReplaceAll(r.Root, "@", cwd) + routePath
	path = strings.ReplaceAll(path, "\\", "/")

	rapi_core.DisableCors(args.RW)
	http.ServeFile(args.RW, args.R, path)
}
