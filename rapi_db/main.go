package rapi_db

import (
	"github.com/Knetic/govaluate"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"github.com/maldan/go-rapi/rapi_error"
)

type SearchResult[T any] struct {
	Count  int `json:"count"`
	Total  int `json:"total"`
	Page   int `json:"page"`
	Result []T `json:"result"`
}

var DataAccess map[string]IDataBase

type ArgsSearch struct {
	Table  string `json:"table"`
	Filter string `json:"filter"`
	Id     int    `json:"id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type ArgsUpdate struct {
	Table string `json:"table"`
	Id    int    `json:"id"`
	Data  string `json:"data"`
}

type DataApi struct {
}

func FilterByExpression[T any](slice []T, expr string, filter func(t *T) map[string]any) []T {
	expression, err := govaluate.NewEvaluableExpression(expr)
	rapi_error.FatalIfError(err)

	return cmhp_slice.Filter(slice, func(v *T) bool {
		m := filter(v)
		result, _ := expression.Evaluate(m)
		return result.(bool)
	})
}

func (u DataApi) GetSettings(args ArgsSearch) any {
	return DataAccess[args.Table].GetSettings()
}

func (u DataApi) GetTableList() []string {
	l := make([]string, 0)
	for k, _ := range DataAccess {
		l = append(l, k)
	}
	return l
}

func (u DataApi) GetSearch(args ArgsSearch) SearchResult[any] {
	f := cmhp_convert.FromBase64(args.Filter)
	return DataAccess[args.Table].Search(string(f), args.Offset, args.Limit)
}

func (u DataApi) GetById(args ArgsSearch) any {
	return DataAccess[args.Table].GetById(args.Id)
}

func (u DataApi) DeleteById(args ArgsSearch) {
	DataAccess[args.Table].DeleteById(args.Id)
}

func (u DataApi) PostById(args ArgsUpdate) {
	DataAccess[args.Table].UpdateById(args.Id, args.Data)
}
