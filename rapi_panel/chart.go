package rapi_panel

import "github.com/maldan/go-cmhp/cmhp_slice"

type ChartApi struct {
}

type ArgsChart struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
	Data   string `json:"data"`
}

func (u ChartApi) GetList(args ArgsControl) any {
	return Config.ChartList
}

func (u ChartApi) PostExecute(args ArgsChart) any {
	c, ok := cmhp_slice.Find(Config.ChartList, func(c ChartCommand) bool {
		return c.Name == args.Name && c.Folder == args.Folder
	})
	if !ok {

	}
	return c.Func(args.Data)
}
