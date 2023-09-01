package rapi

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-rapi/core/handler"
	"github.com/maldan/go-rapi/rapi_config"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_error"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func HandleError(args *rapi_core.HandlerArgs) {
	err := recover()
	if err == nil {
		return
	}

	// Set error output as json
	args.RW.Header().Add("Content-Type", "application/json")

	switch e := err.(type) {
	case rapi_error.Error:
		args.RW.WriteHeader(e.Code)
		message, _ := json.Marshal(e)
		args.RW.Write(message)
		if args.DebugMode {
			//rapi_debug.Log(args.Id).SetError(e)
			//rapi_debug.Log(args.Id).SetArgs(args.MethodArgs)
		}
	default:
		_, file, line, _ := runtime.Caller(3)

		for i := 0; i < 10; i++ {
			p, f, l, ok := runtime.Caller(i)
			if ok {
				fmt.Printf("%v %v:%v\n", p, f, l)
			}
		}

		args.RW.WriteHeader(500)
		// fmt.Println(string(debug.Stack()))
		ee := rapi_error.Error{
			Code:        500,
			Type:        "unknown",
			Description: fmt.Sprintf("%v", e),
			Line:        line,
			File:        file,
			// Stack:       string(debug.Stack()),
			Created: time.Now(),
		}
		message, _ := json.Marshal(ee)
		args.RW.Write(message)
		if args.DebugMode {
			//rapi_debug.Log(args.Id).SetError(ee)
			//rapi_debug.Log(args.Id).SetArgs(args.MethodArgs)
		}
	}
}

func getHandler(url string, routers []rapi_config.RouteHandler) (string, rapi_config.Handler) {
	for i := 0; i < len(routers); i++ {
		if strings.HasPrefix(url, routers[i].Path) {
			return routers[i].Path, routers[i].Handler
		}
	}

	return "", handler.Undefined{}
}

func Start2(config rapi_config.Config) {
	// Entry point
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		defer HandleError(&rapi_core.HandlerArgs{Route: "123", RW: response, R: request})

		// Disable cors for all queries
		rapi_core.DisableCors(response)

		// Fuck options
		if request.Method == "OPTIONS" {
			response.WriteHeader(200)
			return
		}

		route, h := getHandler(request.URL.Path, config.Router)
		fmt.Printf("%v - %v\n", route, h)
		h.Handle(rapi_config.HandlerArgs{Path: route, Request: request, Response: response})
	})

	err := http.ListenAndServe(config.Host, nil)
	rapi_error.FatalIfError(err)
}
