package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"github.com/maldan/go-rapi/rapi_db"
	"github.com/maldan/go-rapi/rapi_error"
)

type TestData[T any] struct {
}

func (d TestData[T]) GetSettings() rapi_db.DataSettings {
	return rapi_db.DataSettings{
		IsDeletable: true,
		IsEditable:  true,
		FieldList: []rapi_db.FieldInfo{
			{Name: "id", Type: "int"},
			{Name: "email", IsEdit: true, Type: "string"},
			{Name: "password", IsHide: true, Type: "string"},
			{Name: "balance", IsEdit: true, Type: "int"},
			{Name: "gay", IsEdit: true, Type: "bool"},
			{Name: "lox", Type: "bool"},
		},
	}
}

func (d TestData[T]) Search(filter string, offset int, limit int) rapi_db.SearchResult[any] {
	newList := list
	if filter != "" {
		// Expr
		expression, err := govaluate.NewEvaluableExpression(filter)
		rapi_error.FatalIfError(err)

		// expr := strings.Replace(filter, "email=", "", 1)
		newList = cmhp_slice.Filter(list, func(t *User) bool {
			result, _ := expression.Evaluate(map[string]any{
				"id":       t.Id,
				"email":    t.Email,
				"password": t.Password,
			})
			return result.(bool)
		})
	}

	l := cmhp_slice.Paginate(newList, offset, limit)
	return rapi_db.SearchResult[any]{
		Count:  len(l),
		Total:  len(newList),
		Page:   offset / limit,
		Result: cmhp_slice.ToAny(l),
	}
}

func (d TestData[T]) GetById(id int) any {
	return list[id-1]
}

func (d TestData[T]) UpdateById(id int, raw string) {
	var u = list[id-1]
	data := cmhp_convert.JsonToStruct[User](raw)
	fmt.Printf("%v", raw)
	u.Email = data.Email
	u.Password = data.Password
	u.Balance = data.Balance
	u.Gay = data.Gay
	list[id-1] = u
}

func (d TestData[T]) DeleteById(id int) {
	list = cmhp_slice.Filter(list, func(x *User) bool { return x.Id != id })
}
