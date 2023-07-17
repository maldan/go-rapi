package rapi_panel

type PanelCommand struct {
	Folder string           `json:"folder"`
	Name   string           `json:"name"`
	Func   func(string) any `json:"-"`
}

type ChartCommand struct {
	Folder string             `json:"folder"`
	Name   string             `json:"name"`
	Func   func(string) []any `json:"-"`
}

type PanelConfig struct {
	CommandList []PanelCommand
	ChartList   []ChartCommand
	DataAccess  map[string]map[string]func(DataArgs) any
	Password    string
}

var Config PanelConfig
