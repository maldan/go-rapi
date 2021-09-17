package rapi_file

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/maldan/go-rapi/rapi_core"
)

type FileHandler struct {
	Root string
}

func (r FileHandler) Handle(args rapi_core.HandlerArgs) {
	cwd, _ := os.Getwd()
	path := strings.ReplaceAll(r.Root, "@", cwd) + args.R.URL.Path
	path = strings.ReplaceAll(path, "\\", "/")
	fmt.Println(path)

	rapi_core.DisableCors(args.RW)
	http.ServeFile(args.RW, args.R, path)
}
