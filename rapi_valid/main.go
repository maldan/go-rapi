package rapi_valid

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_core"
)

func FailIfNotIncludes(val string, values []string) {
	for _, v := range values {
		if v == val {
			return
		}
	}

	rapi_core.Fatal(rapi_core.Error{
		Description: fmt.Sprintf("Value %v not included in %v", val, values),
	})
}
