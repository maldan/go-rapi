package rapi_compiler

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"strings"
)

type _rewrite struct {
	rexExp     string
	httpMethod string
	url        string
}

func parseFile(t string) []_rewrite {
	lines := strings.Split(t, "\n")
	out := make([]_rewrite, 0)

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if !strings.Contains(line, "// @rewrite") {
			continue
		}

		r := _rewrite{
			rexExp: strings.ReplaceAll(line, "// @rewrite ", ""),
		}

		// Find method
		for j := 0; j < 5; j++ {
			if strings.Contains(lines[i+j], "func (") {
				methodInfo := parseFunctionInfo(lines[i+j])
				r.httpMethod = methodInfo.HttpMethod
				r.url = methodInfo.Url
				break
			}
		}

		out = append(out, r)
	}

	return out
}

func CompileRewrite(inDir string, outPath string) {
	list, _ := cmhp_file.ListAll(inDir)
	finalOut := `package core

		import "github.com/maldan/go-rapi"

		var RewriteConfig = []rapi.RewriteUrl {
	`
	for _, l := range list {
		t, _ := cmhp_file.ReadText(l.FullPath)
		ll := parseFile(t)
		for _, x := range ll {
			finalOut += fmt.Sprintf("{\"%v\", \"%v\", \"%v\"},\n", x.httpMethod, x.rexExp, x.url)
		}
	}
	finalOut += "}\n"

	cmhp_file.Write(outPath, finalOut)
}
