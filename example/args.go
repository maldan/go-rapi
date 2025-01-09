package main

import (
	"github.com/maldan/go-rapi/rapi_const"
)

type ArgsId struct {
	Id string
}

type ArgsPhoto struct {
	File rapi_const.File `json:"file"`
}
