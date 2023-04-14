package rapi_valid

import (
	"fmt"
	"github.com/maldan/go-rapi/rapi_error"
	"reflect"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func CheckEmail(email string) string {
	// Trim email from spaces
	email = strings.Trim(strings.ToLower(email), " ")

	// Check
	if !emailRegex.MatchString(email) {
		rapi_error.Fatal(rapi_error.Error{Code: 500, Field: "email", Description: "Incorrect email"})
	}

	// Return cleaned email
	return email
}

func CheckPassword(password1 string, password2 string) {
	// Check password
	if len(password1) < 6 {
		rapi_error.Fatal(rapi_error.Error{
			Code: 500, Field: "password",
			Description: "Password must contain at least 6 characters",
		})
	}
	if password1 != password2 {
		rapi_error.Fatal(rapi_error.Error{
			Code: 500, Field: "password",
			Description: "Passwords do not match",
		})
	}
}

func Required[T any](args T, fields []string) {
	for _, v := range fields {
		f := reflect.ValueOf(args)
		if f.FieldByName(v).Kind() == reflect.Invalid {
			rapi_error.Fatal(rapi_error.Error{
				Description: fmt.Sprintf("Field '%v' is required", v),
			})
		}
		if f.FieldByName(v).IsZero() {
			rapi_error.Fatal(rapi_error.Error{
				Description: fmt.Sprintf("Field '%v' is required", v),
			})
		}
	}
}

func TrimAll[T any](args *T) {
	typeOf := reflect.TypeOf(args)
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Type.Kind() == reflect.String {
			Trim(args, []string{typeOf.Field(i).Name})
		}
	}
}

func Trim[T any](args *T, fields []string) {
	for _, field := range fields {
		// struct dereference
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
