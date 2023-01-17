package rapi_rest

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_error"
)

type UndefinedHandler struct {
}

func (r UndefinedHandler) Handle(args rapi_core.HandlerArgs) {
	// Handle panic
	defer rapi_core.HandleError(args)

	// Disable cors
	rapi_core.DisableCors(args.RW)

	rapi_error.Fatal(rapi_error.Error{
		Code: 404,
		Description: fmt.Sprintf(
			"Resource for '%v' route not found",
			args.Route,
		),
	})
}
