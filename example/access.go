package main

import (
	"github.com/Knetic/govaluate"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"github.com/maldan/go-rapi/rapi_error"
	"github.com/maldan/go-rapi/rapi_panel"
)

var TestAccess = map[string]func(rapi_panel.DataArgs) any{
	rapi_panel.GetSettings: func(args rapi_panel.DataArgs) any {
		return rapi_panel.DataSettings{
			// IsDeletable: true,
			// IsEditable:  true,
			FieldList: []rapi_panel.FieldInfo{
				{Name: "id", Type: "int"},
				{Name: "email", IsEdit: true, Type: "string"},
				{Name: "password", IsHide: true, Type: "string"},
				{Name: "balance", IsEdit: true, Type: "int"},
				{Name: "gay", IsEdit: true, Type: "bool"},
				{Name: "lox", Type: "bool"},

				{
					Name: "havePermission", Type: rapi_panel.TypeBitmask,
					Label: "Can Add Photo,Can Add Metadata,Can Add Documents,Can Suck,Can Fuck,Can Go Crazy,Can Kill Himself,Can Shit Under Himself",
				},

				{Name: "overridePermission", IsEdit: true, Type: rapi_panel.TypeBitmask,
					Label: "Can Add Photo,Can Add Metadata,Can Add Documents,Can Suck,Can Fuck,Can Go Crazy,Can Kill Himself,Can Shit Under Himself",
				},
			},
		}
	},
	rapi_panel.GetById: func(args rapi_panel.DataArgs) any {
		return list[args.Id-1]
	},
	rapi_panel.Search: func(args rapi_panel.DataArgs) any {
		newList := list
		if args.Filter != "" {
			// Expr
			expression, err := govaluate.NewEvaluableExpression(args.Filter)
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

		l := cmhp_slice.Paginate(newList, args.Offset, args.Limit)
		return rapi_panel.SearchResult[any]{
			Count:  len(l),
			Total:  len(newList),
			Page:   args.Offset / args.Limit,
			Result: cmhp_slice.ToAny(l),
		}
	},
	rapi_panel.UpdateById: func(args rapi_panel.DataArgs) any {
		var u = list[args.Id-1]
		data := cmhp_convert.JsonToStruct[User](args.Data)

		u.Email = data.Email
		u.Password = data.Password
		u.Balance = data.Balance
		u.Gay = data.Gay
		u.OverridePermission = data.OverridePermission
		list[args.Id-1] = u

		return data
	},
}
