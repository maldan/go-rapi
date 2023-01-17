package rapi_valid

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_error"
	"reflect"
)

func FailIfNotIncludes(val string, values []string) {
	for _, v := range values {
		if v == val {
			return
		}
	}

	rapi_error.Fatal(rapi_error.Error{
		Description: fmt.Sprintf("Value %v not included in %v", val, values),
	})
}

func Required[T comparable](args T, fields []string) {
	for _, v := range fields {
		f := reflect.ValueOf(args)
		if f.FieldByName(v).IsZero() {
			rapi_error.Fatal(rapi_error.Error{
				Description: fmt.Sprintf("Field '%v' is required", v),
			})
		}
	}
}
