package rapi_core

import (
	"net/http"
	"time"
)

type HandlerArgs struct {
	Route   string
	RW      http.ResponseWriter
	R       *http.Request
	RawBody []byte
}

type Context struct {
	IsSkipProcessing bool
	RW               http.ResponseWriter
	R                *http.Request
}

type Handler interface {
	Handle(args HandlerArgs)
}

type Error struct {
	Status      bool      `json:"status"`
	Code        int       `json:"code"`
	Type        string    `json:"type"`
	Field       string    `json:"field"`
	Description string    `json:"description"`
	File        string    `json:"file"`
	Line        int       `json:"line"`
	Created     time.Time `json:"created"`
}
