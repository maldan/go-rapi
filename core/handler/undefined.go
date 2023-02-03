package handler

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_config"
	"github.com/maldan/go-rapi/rapi_error"
)

type Undefined struct {
}

func (r Undefined) Handle(args rapi_config.HandlerArgs) {
	rapi_error.Fatal(rapi_error.Error{
		Code: 404,
		Description: fmt.Sprintf(
			"Resource for '%v' route not found",
			1, //args.Route,
		),
	})
}
