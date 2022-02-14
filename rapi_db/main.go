package rapi_db

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_reflect"
	"github.com/maldan/go-rapi/rapi_core"
)

var DbPath string

func Load(id string, v interface{}) {
	name := strings.ToLower(reflect.TypeOf(v).Elem().Name())
	if !cmhp_file.Exists(DbPath + "/" + name + "/" + id + ".json") {
		rapi_core.Fatal(rapi_core.Error{
			Code:        404,
			Description: fmt.Sprintf("%v with id %v not found", strings.Title(name), id),
		})
	}
	err := cmhp_file.ReadJSON(DbPath+"/"+name+"/"+id+".json", v)
	if err != nil {
		rapi_core.Fatal(rapi_core.Error{
			Description: err.Error(),
		})
	}
}

func Save(v interface{}) {
	name := strings.ToLower(reflect.TypeOf(v).Elem().Name())
	id := cmhp_crypto.UID(10)

	cmhp_reflect.SetField(v, "Id", id)
	cmhp_reflect.SetField(v, "Created", time.Now())

	// Save to file
	err := cmhp_file.Write(DbPath+"/"+name+"/"+id+".json", v)
	if err != nil {
		rapi_core.Fatal(rapi_core.Error{
			Description: err.Error(),
		})
	}
}

func Update(id string, v interface{}) {
	name := strings.ToLower(reflect.TypeOf(v).Elem().Name())

	// Save to file
	err := cmhp_file.Write(DbPath+"/"+name+"/"+id+".json", v)
	if err != nil {
		rapi_core.Fatal(rapi_core.Error{
			Description: err.Error(),
		})
	}
}
