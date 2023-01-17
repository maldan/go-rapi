package main

import (
	"github.com/maldan/go-cmhp/cmhp_slice"
	"github.com/maldan/go-rapi"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_db"
	"github.com/maldan/go-rapi/rapi_file"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_rest"
)

type TestData[T comparable] struct {
}

var list = make([]User, 0)

func (d TestData[T]) GetStruct() any {
	return map[string]string{
		"id":       "int",
		"email":    "string",
		"password": "string",
	}
}

func (d TestData[T]) Search(offset int, limit int) []any {
	var out = make([]any, 0)
	for _, x := range list {
		out = append(out, x)
	}
	return cmhp_slice.Paginate(out, offset, limit)
}

func (d TestData[T]) GetById(id int) any {
	return list[id-1]
}

func (d TestData[T]) UpdateById(id int, v any) {
	var u = list[id-1]
	u.Password = v.(map[string]any)["password"].(string)
	list[id-1] = u
}

func (d TestData[T]) DeleteById(id int) {
	list = cmhp_slice.Filter(list, func(x User) bool { return x.Id != id })
}

func main() {
	rapi_log.Info("Fuck")
	rapi_log.Info("Suck")
	rapi_log.Error("Oak")

	// Test
	list = append(list, User{Id: 1, Email: "lox"})
	list = append(list, User{Id: 2, Email: "a"})
	list = append(list, User{Id: 3, Email: "b"})

	xx := TestData[User]{}

	rapi.Start(rapi.Config{
		Host: "127.0.0.1:16000",
		Router: map[string]rapi_core.Handler{
			"/": rapi_file.FileHandler{Root: "@"},
			"/api": rapi_rest.ApiHandler{
				Controller: map[string]interface{}{
					"user":     UserApi{},
					"template": TemplateApi{},
				},
			},
		},
		DisableJsonWrapper: true,
		DebugMode:          true,
		DataAccess: map[string]rapi_db.IDataBase{
			"test": xx,
		},
	})
}
