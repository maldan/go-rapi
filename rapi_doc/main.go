package rapi_doc

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"reflect"
	"regexp"
	"strings"

	"github.com/maldan/go-cmhp/cmhp_string"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_rest"
)

type DebugApi struct {
}

type DebugPanelApi struct {
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
	Url      string `json:"url"`

	Controller string `json:"controller"`
	HttpMethod string `json:"httpMethod"`
	Name       string `json:"name"`

	Input *MethodInput `json:"input"`
}

type ArgsPostmanCollection struct {
	IsHttps bool `json:"isHttps"`
}

var Router map[string]rapi_core.Handler
var PanelPage string
var PanelJs string
var PanelCss string
var Host string

// s
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

func (r DebugPanelApi) GetIndex(ctx *rapi_core.Context) {
	ctx.IsSkipProcessing = true
	ctx.RW.Header().Set("Content-Type", "text/html")
	ctx.RW.Write([]byte(PanelPage))
}

func (r DebugPanelApi) GetJs(ctx *rapi_core.Context) {
	ctx.IsSkipProcessing = true
	ctx.RW.Header().Set("Content-Type", "text/javascript")
	ctx.RW.Write([]byte(PanelJs))
}

func (r DebugPanelApi) GetCss(ctx *rapi_core.Context) {
	ctx.IsSkipProcessing = true
	ctx.RW.Header().Set("Content-Type", "text/css")
	ctx.RW.Write([]byte(PanelCss))
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
					var argument interface{}

					// Get argument
					if methodType.NumIn() == 2 {
						argument = reflect.New(methodType.In(1)).Interface()
					}

					// 2 Args
					if methodType.NumIn() == 3 {
						argument = reflect.New(methodType.In(2)).Interface()
					}

					// Add method
					m := Method{
						FullPath: k + "/" + route + "/" + methodName,
						Url:      "//" + Host + k + "/" + route + "/" + methodName,

						Controller: route,
						HttpMethod: httpMethod,
						Name:       methodName,
					}
					if argument != nil {
						m.Input = GetInput("", argument)
					}
					out = append(out, m)

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

func FillBody(input *MethodInput, out *map[string]any) {
	//(*out)[input.Name] = ""
	for _, f := range input.FieldList {
		if f.Type == "int" || f.Type == "uint" || f.Type == "int32" || f.Type == "uint32" {
			(*out)[f.Name] = 0
		} else if f.Type == "float32" || f.Type == "float64" {
			(*out)[f.Name] = 0.0
		} else if f.Kind == reflect.Slice.String() {
			(*out)[f.Name] = make([]string, 0)
		} else {
			(*out)[f.Name] = ""
		}
	}
}

func QueryStr(input *MethodInput) string {
	x := ""
	for _, f := range input.FieldList {
		x += f.Name + "=&"
	}
	return x
}

func (r DebugApi) GetPostmanCollection(ctx *rapi_core.Context) any {
	// items := make([]any, 0)
	methodList := r.GetMethodList()

	// Protocol
	protocol := "http://"
	if ctx.R.Header.Get("X-Forwarded-Proto") != "" {
		protocol = ctx.R.Header.Get("X-Forwarded-Proto") + "://"
	}

	type folder struct {
		Name string `json:"name"`
		Item []any  `json:"item"`
	}
	folders := make(map[string]folder)

	for _, m := range methodList {
		if m.HttpMethod == "GET" {
			item := map[string]any{
				"name": m.Name,
				"request": map[string]any{
					"url":    protocol + ctx.R.Host + m.FullPath + "?" + QueryStr(m.Input),
					"method": m.HttpMethod,
				},
			}
			_, ok := folders[m.Controller]
			if !ok {
				items := make([]any, 0)
				items = append(items, item)
				folders[m.Controller] = folder{
					Name: m.Controller,
					Item: items,
				}
			} else {
				items := folders[m.Controller].Item
				items = append(items, item)
				folders[m.Controller] = folder{
					Name: m.Controller,
					Item: items,
				}
			}
		} else {
			bodyFormat := map[string]any{}
			FillBody(m.Input, &bodyFormat)
			bodyStr, _ := json.Marshal(&bodyFormat)

			item := map[string]any{
				"name": m.Name,
				"request": map[string]any{
					"url":    protocol + ctx.R.Host + m.FullPath,
					"method": m.HttpMethod,
					"header": []any{
						map[string]any{
							"key":   "Content-Type",
							"value": "application/json",
						},
					},
					"body": map[string]any{
						"mode": "raw",
						"raw":  string(bodyStr),
					},
				},
			}

			_, ok := folders[m.Controller]
			if !ok {
				items := make([]any, 0)
				items = append(items, item)
				folders[m.Controller] = folder{
					Name: m.Controller,
					Item: items,
				}
			} else {
				items := folders[m.Controller].Item
				items = append(items, item)
				folders[m.Controller] = folder{
					Name: m.Controller,
					Item: items,
				}
			}
		}
	}

	// Combine all
	items := make([]any, 0)
	for _, folder := range folders {
		items = append(items, folder)
	}

	return map[string]any{
		"info": map[string]any{
			"name":        ctx.R.Host + " Api",
			"_postman_id": cmhp_crypto.Sha1(ctx.R.Host),
			"description": fmt.Sprintf("%v Api", ctx.R.Host),
			"schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"item": items,
	}
}
