package rapi

import (
	_ "embed"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_doc"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_rest"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

//go:embed panel/dist/index.html
var PanelPage string

//go:embed panel/dist/assets/index.js
var PanelJs string

//go:embed panel/dist/assets/index.css
var PanelCss string

type Config struct {
	Host               string
	Router             map[string]rapi_core.Handler
	IsHttps            bool
	KeyFile            string
	CertFile           string
	DisableJsonWrapper bool
	Rewrite            map[string]string
}

var redirectRegExpMap = map[*regexp.Regexp]string{}

func HandleRewrite(url *url.URL) {
	for reg, p := range redirectRegExpMap {
		match := reg.FindStringSubmatch(url.Path)
		if len(match) == 0 {
			continue
		}

		// Replace path groups to query
		for i, name := range reg.SubexpNames() {
			if i != 0 && name != "" {
				url.RawQuery += "&" + name + "=" + match[i]
			}
		}

		// Set redirect
		url.Path = p
		return
	}
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

	// Prepare rewrite
	for k, v := range config.Rewrite {
		redirectRegExpMap[regexp.MustCompile(k)] = v
	}

	// Set router for debug
	rapi_doc.PanelPage = PanelPage
	rapi_doc.PanelJs = PanelJs
	rapi_doc.PanelCss = PanelCss
	rapi_doc.Router = config.Router
	rapi_doc.Host = config.Host

	rapi_log.Info("Start RApi server %v", config.Host)
	rapi_log.Info("Disable json wrapper %v", config.DisableJsonWrapper)

	// Entry point
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// Redirect handler
		HandleRewrite(r.URL)

		route, handler := rapi_core.GetHandler(r.URL.Path, config.Router)
		if handler == nil {
			handler = rapi_rest.UndefinedHandler{}
			route = r.URL.Path
		}

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
