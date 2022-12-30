package main

import "github.com/maldan/go-rapi/rapi_core"

type ArgsId struct {
	Id string
}

type ArgsPhoto struct {
	File rapi_core.File `json:"file"`
}
