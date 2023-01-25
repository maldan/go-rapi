package rapi_valid

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_error"
	"reflect"
	"strings"
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

func Required[T any](args T, fields []string) {
	for _, v := range fields {
		f := reflect.ValueOf(args)
		if f.FieldByName(v).IsZero() {
			rapi_error.Fatal(rapi_error.Error{
				Description: fmt.Sprintf("Field '%v' is required", v),
			})
		}
	}
}

func Trim[T any](args *T, fields []string) {
	for _, field := range fields {
		// struct
		f := reflect.ValueOf(args).Elem()
		if f.FieldByName(field).Kind() == reflect.String {
			value := f.FieldByName(field).Interface().(string)

			if f.FieldByName(field).CanSet() {
				fmt.Printf("%v\n", value)
				f.FieldByName(field).SetString(strings.Trim(value, " "))
			}
		}
	}
}
