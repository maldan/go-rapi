package rapi_core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func HandleError(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	if err := recover(); err != nil {
		switch e := err.(type) {
		case Error:
			rw.WriteHeader(e.Code)
			message, _ := json.Marshal(e)
			rw.Write(message)
		default:
			_, file, line, _ := runtime.Caller(3)

			rw.WriteHeader(500)
			fmt.Println(string(debug.Stack()))
			ee := Error{
				Code:        500,
				Type:        "unknown",
				Description: fmt.Sprintf("%v", e),
				Line:        line,
				File:        file,
				Created:     time.Now(),
			}
			message, _ := json.Marshal(ee)
			rw.Write(message)
		}
	}
}

func Fatal(err Error) {
	_, file, line, _ := runtime.Caller(1)
	ff := strings.Split(file, "/")

	if err.Code == 0 {
		err.Code = 500
	}
	if err.Type == "" {
		err.Type = "unknown"
	}
	err.File = strings.Join(ff[len(ff)-2:], "/")
	err.Line = line
	err.Created = time.Now()

	panic(err)
}

func FatalIfError(err error) {
	if err != nil {
		Fatal(Error{
			Description: err.Error(),
		})
	}
}
