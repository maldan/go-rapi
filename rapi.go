package rapi

import (
	_ "embed"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_debug"
	"github.com/maldan/go-rapi/rapi_doc"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_panel"
	"github.com/maldan/go-rapi/rapi_rest"
	"github.com/maldan/go-rapi/rapi_test"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

//go:embed panel/dist/index.html
var PanelPage string

//go:embed panel/dist/assets/index.js
var PanelJs string

//go:embed panel/dist/assets/index.css
var PanelCss string

type RewriteUrl struct {
	Method string
	RegExp string
	Url    string
}

type Config struct {
	Host               string
	DebugHost          string
	Router             map[string]rapi_core.Handler
	IsHttps            bool
	KeyFile            string
	CertFile           string
	DisableJsonWrapper bool
	Rewrite            []RewriteUrl
	TestList           []rapi_test.TestCase
	// DataAccess         map[string]map[string]func(rapi_panel.DataArgs) any
	PanelConfig rapi_panel.PanelConfig
	DebugMode   bool
	Log         rapi_debug.LogConfig
}

type Rewrite struct {
	RegExp     *regexp.Regexp
	HttpMethod string
	Url        string
}

// var redirectRegExpMap = map[*regexp.Regexp]string{}
var redirectList = make([]Rewrite, 0)

func HandleRewrite(httpMethod string, url *url.URL) {
	for _, r := range redirectList {
		if !strings.Contains(r.HttpMethod, httpMethod) {
			continue
		}

		// Check regex
		match := r.RegExp.FindStringSubmatch(url.Path)
		if len(match) == 0 {
			continue
		}

		// Replace path groups to query
		for i, name := range r.RegExp.SubexpNames() {
			if i != 0 && name != "" {
				url.RawQuery += "&" + name + "=" + match[i]
			}
		}

		// Set redirect
		url.Path = r.Url
		return
	}
	/*for reg, p := range redirectRegExpMap {
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
	}*/
}

func Start(config Config) {
	// Set debug api
	config.Router["/debug"] = rapi_rest.ApiHandler{
		Controller: map[string]interface{}{
			"api":     rapi_doc.DebugApi{},
			"log":     rapi_log.LogApi{},
			"panel":   rapi_doc.DebugPanelApi{},
			"data":    rapi_panel.DataApi{},
			"control": rapi_panel.ControlApi{},
			"chart":   rapi_panel.ChartApi{},
			"backup":  rapi_panel.BackupApi{},
		},
	}

	// Prepare rewrite
	for _, v := range config.Rewrite {
		redirectList = append(redirectList, Rewrite{
			Url:        v.Url,
			RegExp:     regexp.MustCompile(v.RegExp),
			HttpMethod: v.Method,
		})
	}

	// Set router for debug
	rapi_doc.PanelPage = PanelPage
	rapi_doc.PanelJs = PanelJs
	rapi_doc.PanelCss = PanelCss
	rapi_doc.Router = config.Router
	rapi_doc.Host = config.DebugHost
	if rapi_doc.Host == "" {
		rapi_doc.Host = config.Host
	}
	rapi_doc.TestList = config.TestList
	rapi_doc.OnRequestSearch = config.Log.OnSearch
	rapi_panel.Config = config.PanelConfig

	// Run backup schedule
	go rapi_panel.Config.BackupConfig.Run()

	rapi_log.Info("Start RApi server %v", config.Host)
	rapi_log.Info("Disable json wrapper %v", config.DisableJsonWrapper)
	rapi_log.Info("Debug host %v", rapi_doc.Host)
	rapi_log.Info("Debug mode %v", config.DebugMode)

	// Unauth access
	allowList := []string{
		"/debug/panel", "/debug/panel/", "/debug/panel/js",
		"/debug/panel/css", "/debug/api/auth",
	}

	// Entry point
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		id := cmhp_crypto.UID(16)
		debugMode := config.DebugMode

		// Auth mode
		if strings.HasPrefix(r.URL.Path, "/debug/") {
			debugMode = false

			// Api access
			itsOkay := false
			for _, aUrl := range allowList {
				if r.URL.Path == aUrl {
					itsOkay = true
					break
				}
			}

			// Checking access
			if !itsOkay {
				if r.Header.Get("DebugAuthKey") != cmhp_convert.ToUrlBase64("admin_"+config.PanelConfig.Password) {
					// @TODO enable
					// rapi_error.Fatal(rapi_error.Error{Code: 401})
				}
			}
		}

		if debugMode {
			// Ip
			ipAddress := r.Header.Get("X-Real-Ip")
			if ipAddress == "" {
				ipAddress = r.Header.Get("X-Forwarded-For")
			}
			if ipAddress == "" {
				ipAddress = r.RemoteAddr
			}
			ipAddress = strings.Split(ipAddress, ":")[0]
			rapi_debug.
				GetRequestLog(id).
				SetRequest(r.Method, r.URL.Path).
				SetRemoteAddr(ipAddress)
		}

		// Redirect handler
		HandleRewrite(r.Method, r.URL)

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
			DebugMode:          debugMode,
			Id:                 id,
		})

		// On log request
		if debugMode && config.Log.OnRequest != nil {
			config.Log.OnRequest(*rapi_debug.GetRequestLog(id))
		}
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
