package main

import (
	"github.com/maldan/go-rapi"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_file"
	"github.com/maldan/go-rapi/rapi_rest"
)

func main() {
	rapi.Start(rapi.Config{
		Host: "127.0.0.1:5000",
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
