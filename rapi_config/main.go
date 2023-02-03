package rapi_config

import "net/http"

type DebugConfig struct {
	IsEnabled bool
}

type PanelConfig struct {
	Login    string
	Password string
}

type Context struct {
	AccessToken      string
	IsSkipProcessing bool
	IsServeFile      bool
	RW               http.ResponseWriter
	R                *http.Request
}

type Config struct {
	Host string

	Router []RouteHandler

	Debug DebugConfig
	Panel PanelConfig

	EnableJsonWrapper bool
}
