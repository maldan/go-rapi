package handler

import (
	"embed"
	"github.com/maldan/go-rapi/rapi_config"
)

type VFS struct {
	Root string
	Fs   embed.FS
}

func (V VFS) Handle(args rapi_config.HandlerArgs) {
	//TODO implement me
	panic("implement me")
}
