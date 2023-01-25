package rapi_core

import (
	"net/http"
)

type HandlerArgs struct {
	Route              string
	RW                 http.ResponseWriter
	R                  *http.Request
	RawBody            []byte
	Context            *Context
	DisableJsonWrapper bool
	DebugMode          bool
	Id                 string
	MethodArgs         map[string]any
}

type Context struct {
	AccessToken      string
	IsSkipProcessing bool
	IsServeFile      bool
	RW               http.ResponseWriter
	R                *http.Request
}

type File struct {
	Name string
	Mime string
	Size int
	Data []byte
}

type Handler interface {
	Handle(args HandlerArgs)
}
