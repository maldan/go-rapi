package rapi

import (
	_ "embed"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_db"
	"github.com/maldan/go-rapi/rapi_doc"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_rest"
	"log"
	"net/http"
)

//go:embed panel/dist/index.html
var PanelPage string

//go:embed panel/dist/assets/index.d398e5c1.js
var PanelJs string

//go:embed panel/dist/assets/index.96d05a2d.css
var PanelCss string

type Config struct {
	Host               string
	Router             map[string]rapi_core.Handler
	DbPath             string
	IsHttps            bool
	KeyFile            string
	CertFile           string
	DisableJsonWrapper bool
}

func Start(config Config) {

	// Set debug api
	config.Router["/debug"] = rapi_rest.ApiHandler{
		Controller: map[string]interface{}{
			"api":   rapi_doc.DebugApi{},
			"log":   rapi_log.LogApi{},
			"panel": rapi_doc.DebugPanelApi{},
		},
	}

	// Set router for debug
	rapi_doc.PanelPage = PanelPage
	rapi_doc.PanelJs = PanelJs
	rapi_doc.PanelCss = PanelCss
	rapi_doc.Router = config.Router
	rapi_doc.Host = config.Host
	rapi_db.DbPath = config.DbPath

	// Entry point
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		route, handler := rapi_core.GetHandler(r.URL.Path, config.Router)
		handler.Handle(rapi_core.HandlerArgs{
			Route:              route,
			RW:                 rw,
			R:                  r,
			DisableJsonWrapper: config.DisableJsonWrapper,
		})
	})

	// Start server
	if config.IsHttps {
		err := http.ListenAndServeTLS(config.Host, config.CertFile, config.KeyFile, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		if err := http.ListenAndServe(config.Host, nil); err != nil {
			log.Fatal(err)
			return
		}
	}
}
