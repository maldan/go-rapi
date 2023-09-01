package rapi_core

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-rapi/rapi_debug"
	"github.com/maldan/go-rapi/rapi_error"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func HandleError(args *HandlerArgs) {
	args.RW.Header().Add("Content-Type", "application/json")

	if err := recover(); err != nil {
		switch e := err.(type) {
		case rapi_error.Error:
			args.RW.WriteHeader(e.Code)
			message, _ := json.Marshal(e)
			args.RW.Write(message)
			if args.DebugMode {
				rapi_debug.GetRequestLog(args.Id).SetError(e)
				// rapi_debug.GetRequestLog(args.Id).SetArgs(args.MethodArgs)
			}
		default:
			_, file, line, _ := runtime.Caller(3)

			args.RW.WriteHeader(500)
			fmt.Println(string(debug.Stack()))
			ee := rapi_error.Error{
				Code:        500,
				Type:        "unknown",
				Description: fmt.Sprintf("%v", e),
				Line:        line,
				File:        file,
				Created:     time.Now(),
			}
			message, _ := json.Marshal(ee)
			args.RW.Write(message)
			if args.DebugMode {
				rapi_debug.GetRequestLog(args.Id).SetError(ee)
				// rapi_debug.GetRequestLog(args.Id).SetArgs(args.MethodArgs)
			}
		}
	}
}

func GetHandler(url string, routers map[string]Handler) (string, Handler) {
	var most string
	var handler Handler

	for k, v := range routers {
		if strings.HasPrefix(url, k) {
			if len(most) < len(k) {
				most = k
				handler = v
			}
		}
	}

	return most, handler
}

func DisableCors(rw http.ResponseWriter) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Methods", "*")
	rw.Header().Add("Access-Control-allow-Headers", "*")
}
