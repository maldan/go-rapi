package rapi_rest

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_core"
)

type UndefinedHandler struct {
}

func (r UndefinedHandler) Handle(args rapi_core.HandlerArgs) {
	// Handle panic
	defer rapi_core.HandleError(args.RW, args.R)

	// Disable cors
	rapi_core.DisableCors(args.RW)

	rapi_core.Fatal(rapi_core.Error{
		Code: 404,
		Description: fmt.Sprintf(
			"Resource for '%v' route not found",
			args.Route,
		),
	})
}
