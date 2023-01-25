package rapi_panel

import "github.com/maldan/go-cmhp/cmhp_slice"

type ControlApi struct {
}

type ArgsControl struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
}

func (u ControlApi) GetList(args ArgsControl) any {
	return Config.CommandList
}

func (u ControlApi) PostExecute(args ArgsControl) any {
	c, ok := cmhp_slice.Find(Config.CommandList, func(c PanelCommand) bool {
		return c.Name == args.Name && c.Folder == args.Folder
	})
	if !ok {

	}
	return c.Func("")
}
