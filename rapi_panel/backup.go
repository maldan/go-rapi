package rapi_panel

import "github.com/maldan/go-cmhp/cmhp_slice"

type BackupApi struct {
}

type ArgsBackup struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
	Data   string `json:"data"`
}

func (u BackupApi) GetTaskList(args ArgsBackup) any {
	return Config.BackupConfig.TaskList
}

func (u BackupApi) PostExecute(args ArgsBackup) any {
	c, ok := cmhp_slice.Find(Config.ChartList, func(c ChartCommand) bool {
		return c.Name == args.Name && c.Folder == args.Folder
	})
	if !ok {

	}
	return c.Func(args.Data)
}
