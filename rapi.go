package rapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_db"
	"github.com/maldan/go-rapi/rapi_doc"
	"github.com/maldan/go-rapi/rapi_rest"
)

type Config struct {
	Host   string
	Router map[string]rapi_core.Handler
	DbPath string
}

func Start(config Config) {
	// Set debug api
	config.Router["/debug"] = rapi_rest.ApiHandler{
		Controller: map[string]interface{}{
			"api": rapi_doc.DebugApi{},
		},
	}

	// Set router for debug
	rapi_doc.Router = config.Router
	rapi_db.DbPath = config.DbPath

	// Entry point
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		route, handler := rapi_core.GetHandler(r.URL.Path, config.Router)
		handler.Handle(rapi_core.HandlerArgs{
			Route: route,
			RW:    rw,
			R:     r,
		})
	})

	// Start server
	fmt.Println("Start server", config.Host)
	if err := http.ListenAndServe(config.Host, nil); err != nil {
		log.Fatal(err)
		return
	}
}
