package rapi_rest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/maldan/go-cmhp/cmhp_string"
	"github.com/maldan/go-rapi/rapi_core"
)

func ValidateFieldList(s *reflect.Value, ss reflect.Type) {
	amount := s.NumField()

	for i := 0; i < amount; i++ {
		field := s.Field(i)
		fieldName := ss.Field(i).Name
		fieldTag := ss.Field(i).Tag
		isRequired := fieldTag.Get("validation") == "required"

		// Can change field
		if field.IsValid() {
			if field.CanSet() {

				if isRequired && field.IsZero() {
					rapi_core.Fatal(rapi_core.Error{
						Type:  "emptyField",
						Field: cmhp_string.LowerFirst(fieldName),
						Description: fmt.Sprintf(
							"Field %v with type (%v) is required",
							cmhp_string.LowerFirst(fieldName),
							field.Type().Name(),
						),
					})
				}
			}
		}
	}
}

func FillFieldList(s *reflect.Value, ss reflect.Type, params map[string]interface{}) {
	amount := s.NumField()

	if params == nil {
		params = make(map[string]interface{})
	}

	for i := 0; i < amount; i++ {
		field := s.Field(i)
		fieldName := ss.Field(i).Name
		fieldTag := ss.Field(i).Tag
		jsonName := fieldTag.Get("json")

		// Can change field
		if field.IsValid() {
			if field.CanSet() {
				// Skip
				if jsonName == "-" {
					continue
				}

				// Get value
				var v interface{}
				if jsonName != "" {
					x, ok := params[jsonName]
					if x == nil {
						continue
					}
					if ok {
						v = x
					} else {
						continue
					}
				} else {
					x, ok := params[cmhp_string.LowerFirst(fieldName)]
					if x == nil {
						continue
					}
					if ok {
						v = x
					} else {
						continue
					}
				}

				// Get field type
				switch field.Kind() {
				case reflect.String:
					ApplyString(&field, v)
				case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint,
					reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
					ApplyInt(&field, v)
				case reflect.Float32, reflect.Float64:
					ApplyFloat(&field, v)
				case reflect.Bool:
					ApplyBool(&field, v)
				case reflect.Slice:
					ApplySlice(&field, v)
				case reflect.Struct:
					if field.Type().Name() == "File" {
						field.Set(reflect.ValueOf(v))
					} else if field.Type().Name() == "Time" {
						ApplyTime(&field, v)
					} else {
						if reflect.TypeOf(v).Kind() == reflect.Map {
							FillFieldList(&field, reflect.TypeOf(field.Interface()), v.(map[string]interface{}))
						}
					}
				case reflect.Ptr:
					ApplyPtr(&field, v)
					continue
				case reflect.Map:
					ApplyMap(&field, v)
				default:
					continue
				}
			}
		}
	}

	ValidateFieldList(s, ss)
}

func GetMethod(controller interface{}, methodName string, httpMethod string) *reflect.Method {
	finalMethodName := strings.Title(strings.ToLower(httpMethod)) + strings.Title(methodName)
	controllerType := reflect.TypeOf(controller)
	for i := 0; i < controllerType.NumMethod(); i++ {
		method := controllerType.Method(i)
		if method.Name == finalMethodName {
			return &method
		}
	}

	return nil
}

func virtualCall(fn reflect.Method, args ...interface{}) reflect.Value {
	function := reflect.ValueOf(fn.Func.Interface())

	in := make([]reflect.Value, len(args))
	for i, v := range args {
		in[i] = reflect.ValueOf(v)
	}
	r := function.Call(in)
	if len(r) > 0 {
		return r[0]
	}

	return reflect.ValueOf("")
}

func ExecuteMethod(
	controller interface{},
	args rapi_core.HandlerArgs,
	method reflect.Method,
	params map[string]interface{},
) reflect.Value {
	functionType := reflect.TypeOf(method.Func.Interface())

	// No args
	if functionType.NumIn() == 1 {
		return virtualCall(method, controller)
	}

	// Has 1 arg
	if functionType.NumIn() == 2 {
		arg := reflect.New(functionType.In(1)).Interface()
		argValue := reflect.ValueOf(arg).Elem()
		argType := reflect.TypeOf(arg).Elem()

		// Is struct
		if argType.Kind() == reflect.Struct {
			FillFieldList(&argValue, argType, params)
			return virtualCall(method, controller, argValue.Interface())
		}

		// Is string
		if argType.Kind() == reflect.String {
			return virtualCall(method, controller, string(args.RawBody))
		}

		// Bytes
		if argType.Kind() == reflect.Slice {
			return virtualCall(method, controller, args.RawBody)
		}

		// Context
		if argType.Kind() == reflect.Ptr {
			return virtualCall(method, controller, args.Context)
		}
	}

	// Has 2 arg
	if functionType.NumIn() == 3 {
		arg := reflect.New(functionType.In(2)).Interface()
		argValue := reflect.ValueOf(arg).Elem()
		argType := reflect.TypeOf(arg).Elem()

		// Is struct
		if argType.Kind() == reflect.Struct {
			FillFieldList(&argValue, argType, params)
			return virtualCall(method, controller, args.Context, argValue.Interface())
		}

		// Is string
		if argType.Kind() == reflect.String {
			return virtualCall(method, controller, args.Context, string(args.RawBody))
		}

		// Bytes
		if argType.Kind() == reflect.Slice {
			return virtualCall(method, controller, args.Context, args.RawBody)
		}
	}

	return reflect.ValueOf("")
}
