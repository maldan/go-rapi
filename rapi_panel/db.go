package rapi_panel

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"github.com/maldan/go-rapi/rapi_error"
	"sort"
)

const GetSettings = "getSettings"
const GetById = "getById"
const DeleteById = "deleteById"
const UpdateById = "updateById"
const Search = "search"
const Create = "create"
const Export = "export"

const TypeInt = "int"
const TypeFloat = "float"
const TypeString = "string"
const TypeDate = "date"
const TypeDateTime = "datetime"
const TypeBool = "bool"
const TypeBitmask = "bitmask"
const TypeDataUrl = "dataUrl"

type FieldInfo struct {
	Name       string `json:"name"`
	IsEdit     bool   `json:"isEdit"`
	IsHide     bool   `json:"isHide"`
	HasFilter  bool   `json:"hasFilter"`
	Type       string `json:"type"`
	Label      string `json:"label"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	IsTextarea bool   `json:"isTextarea"`
	OptionList []any  `json:"optionList"`
}

type DataSettings struct {
	IsCreatable  bool        `json:"isCreatable"`
	IsEditable   bool        `json:"isEditable"`
	IsDeletable  bool        `json:"isDeletable"`
	IsExportable bool        `json:"isExportable"`
	FieldList    []FieldInfo `json:"fieldList"`
}

type DataArgs struct {
	Id   string `json:"id"`
	Data string `json:"data"`

	Filter map[string]string
	Offset int
	Limit  int
}

type SearchResult[T any] struct {
	Count  int `json:"count"`
	Total  int `json:"total"`
	Page   int `json:"page"`
	Result []T `json:"result"`
}

type ArgsSearch struct {
	Table  string `json:"table"`
	Filter string `json:"filter"`
	Id     string `json:"id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type ArgsUpdate struct {
	Table string `json:"table"`
	Id    string `json:"id"`
	Data  string `json:"data"`
}

type DataApi struct {
}

func (u DataApi) GetSettings(args ArgsSearch) any {
	settings := Config.DataAccess[args.Table][GetSettings](DataArgs{}).(DataSettings)

	// Set deletable
	_, ok := Config.DataAccess[args.Table][DeleteById]
	settings.IsDeletable = ok

	// Set editable
	_, ok = Config.DataAccess[args.Table][UpdateById]
	settings.IsEditable = ok

	// Set creatable
	_, ok = Config.DataAccess[args.Table][Create]
	settings.IsCreatable = ok

	// Set exportable
	_, ok = Config.DataAccess[args.Table][Export]
	settings.IsExportable = ok

	return settings
}

func (u DataApi) GetTableList() []string {
	l := make([]string, 0)
	for k, _ := range Config.DataAccess {
		l = append(l, k)
	}
	sort.SliceStable(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	return l
}

func (u DataApi) GetSearch(args ArgsSearch) SearchResult[any] {
	filter := map[string]string{}
	err := json.Unmarshal(cmhp_convert.FromUrlBase64(args.Filter), &filter)
	rapi_error.FatalIfError(err)

	return Config.DataAccess[args.Table][Search](DataArgs{
		Filter: filter,
		Offset: args.Offset,
		Limit:  args.Limit,
	}).(SearchResult[any])
}

func (u DataApi) GetById(args ArgsSearch) any {
	fmt.Printf("%+v\n", args)
	return Config.DataAccess[args.Table][GetById](DataArgs{
		Id: args.Id,
	})
}

func (u DataApi) GetExport(args ArgsSearch) any {
	filter := map[string]string{}
	err := json.Unmarshal(cmhp_convert.FromUrlBase64(args.Filter), &filter)
	rapi_error.FatalIfError(err)

	return Config.DataAccess[args.Table][Export](DataArgs{
		Filter: filter,
		Offset: args.Offset,
		Limit:  args.Limit,
	})
}

func (u DataApi) DeleteById(args ArgsSearch) {
	Config.DataAccess[args.Table][DeleteById](DataArgs{
		Id: args.Id,
	})
}

func (u DataApi) PostById(args ArgsUpdate) {
	Config.DataAccess[args.Table][UpdateById](DataArgs{
		Id:   args.Id,
		Data: args.Data,
	})
}

func (u DataApi) PostCreate(args ArgsUpdate) {
	Config.DataAccess[args.Table][Create](DataArgs{
		Data: args.Data,
	})
}
