package main

import (
	"github.com/maldan/go-rapi"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_file"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_rest"
)

func main() {
	rapi_log.Info("Fuck")
	rapi_log.Info("Suck")
	rapi_log.Error("Oak")

	rapi.Start(rapi.Config{
		Host: "127.0.0.1:16000",
		Router: map[string]rapi_core.Handler{
			"/": rapi_file.FileHandler{Root: "@"},
			"/api": rapi_rest.ApiHandler{
				Controller: map[string]interface{}{
					"test":  TestApi{},
					"test2": Test2Api{},
				},
			},
		},
	})
}
