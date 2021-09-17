package rapi_doc

import (
	_ "embed"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/maldan/go-cmhp/cmhp_string"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_rest"
)

type DebugApi struct {
}

/*type MethodStruct struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Kind  string      `json:"kind"`
	Value interface{} `json:"value"`
}*/

type MethodInput struct {
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Kind      string         `json:"kind"`
	FieldList []*MethodInput `json:"fieldList"`
}

type Method struct {
	FullPath string `json:"fullPath"`

	Controller string `json:"controller"`
	HttpMethod string `json:"httpMethod"`
	Name       string `json:"name"`

	Input *MethodInput `json:"input"`
}

var Router map[string]rapi_core.Handler

func GetInput(name string, arg interface{}) *MethodInput {
	argValue := reflect.ValueOf(arg).Elem()
	argType := reflect.TypeOf(arg).Elem()

	// If arg is struct
	if argValue.Type().Kind() == reflect.Struct {
		m := MethodInput{
			Name:      name,
			Type:      fmt.Sprintf("%v", argValue.Type()),
			Kind:      argValue.Type().Kind().String(),
			FieldList: make([]*MethodInput, 0),
		}

		// Go over fields
		amount := argValue.NumField()
		for i := 0; i < amount; i++ {
			/*m.FieldList = append(m.FieldList, &MethodInput{
				Name: cmhp_string.LowerFirst(argType.Field(i).Name),
				Type: fmt.Sprintf("%v", argValue.Field(i).Type()),
				Kind: argValue.Field(i).Kind().String(),
			})*/

			m.FieldList = append(
				m.FieldList,
				GetInput(
					cmhp_string.LowerFirst(argType.Field(i).Name),
					reflect.New(argValue.Field(i).Type()).Interface(),
				),
			)
		}

		return &m
	}

	return &MethodInput{
		Name:      name,
		Type:      fmt.Sprintf("%v", argValue.Type()),
		Kind:      argValue.Type().Kind().String(),
		FieldList: make([]*MethodInput, 0),
	}
}

func (r DebugApi) GetMethodList() []Method {
	out := make([]Method, 0)

	var re = regexp.MustCompile(`^(Get|Post|Patch|Delete|Put)(.*?)$`)

	for k, v := range Router {
		if k == "/api" {
			for route, controller := range v.(rapi_rest.ApiHandler).Controller {
				controllerType := reflect.TypeOf(controller)

				// Go over methods
				for i := 0; i < controllerType.NumMethod(); i++ {
					// Get method info
					method := controllerType.Method(i)
					methodName := cmhp_string.LowerFirst(re.ReplaceAllString(method.Name, "$2"))
					httpMethod := strings.ToUpper(re.ReplaceAllString(method.Name, `$1`))
					methodType := reflect.TypeOf(method.Func.Interface())

					// Go over args
					/* for j := 0; j < methodType.NumIn(); j++ {
						// Skip first argument
						if j == 0 {
							continue
						}
					} */

					// Get argument
					argument := reflect.New(methodType.In(1)).Interface()

					// 2 Args
					if methodType.NumIn() == 3 {
						argument = reflect.New(methodType.In(2)).Interface()
					}

					// Add method
					out = append(out, Method{
						FullPath: k + "/" + route + "/" + methodName,

						Controller: route,
						HttpMethod: httpMethod,
						Name:       methodName,

						Input: GetInput("", argument),
					})

					/*methodStruct := make([]MethodStruct, 0)

					for i := 0; i < functionType.NumIn(); i++ {
						// Skip first argument
						if i == 0 {
							continue
						}

						argument := functionType.In(i)
						argsx := reflect.New(argument).Interface()

						s := reflect.ValueOf(argsx).Elem()
						ss := reflect.TypeOf(argsx).Elem()

						if s.Type().Kind() == reflect.Struct {
							Sas(s, ss)
						}
					}*/

				}
			}
		}
	}

	return out
}
