package rapi_panel

type PanelCommand struct {
	Folder string           `json:"folder"`
	Name   string           `json:"name"`
	Func   func(string) any `json:"-"`
}

type PanelConfig struct {
	CommandList []PanelCommand
	DataAccess  map[string]map[string]func(DataArgs) any
}

var Config PanelConfig