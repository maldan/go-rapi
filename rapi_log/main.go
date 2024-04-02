package rapi_log

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_time"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_panel"
	"sync"
	"time"
)

type LogApi struct {
}

type LogData struct {
	Kind    string    `json:"kind"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type ArgsSearch struct {
	Date time.Time `json:"date"`
}

type ArgsDownload struct {
	Category string    `json:"category"`
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

var mu sync.Mutex
var logList = make([]LogData, 0)
var Config rapi_panel.PanelConfig

func Error(format string, a ...any) {
	/*mu.Lock()
	defer mu.Unlock()
	logList = append(logList, LogData{
		Kind:    "error",
		Message: message,
		Time:    time.Now(),
	})*/
	msg := fmt.Sprintf(format, a...)
	fmt.Errorf("[ERR ] (%v) - %v\n", cmhp_time.Format(time.Now(), "YYYY-MM-DD HH:mm:ss.SSS"), msg)
}

func Info(format string, a ...any) {
	/*mu.Lock()
	defer mu.Unlock()
	logList = append(logList, LogData{
		Kind:    "info",
		Message: message,
		Time:    time.Now(),
	})*/
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("[INFO] (%v) - %v\n", cmhp_time.Format(time.Now(), "YYYY-MM-DD HH:mm:ss.SSS"), msg)
}

/* func (r LogApi) GetIndex() []LogData {
	return logList
} */

/*func (r LogApi) GetSearch(args ArgsSearch) []LogData {
	nLogs := logList

	if args.Date.Year() > 1 {
		from := time.Date(args.Date.Year(), args.Date.Month(), args.Date.Day(), 0, 0, 0, 0, args.Date.Location())
		to := time.Date(args.Date.Year(), args.Date.Month(), args.Date.Day(), 23, 59, 59, 0, args.Date.Location())

		nLogs = cmhp_slice.Filter(logList, func(t *LogData) bool {
			return t.Time.Unix() >= from.Unix() && t.Time.Unix() <= to.Unix()
		})
	}
	return nLogs
}*/

func (r LogApi) GetCategoryList(args ArgsSearch) any {
	return Config.LogsConfig
}

func (r LogApi) GetDownload(ctx *rapi_core.Context, args ArgsDownload) {
	ctx.IsSkipProcessing = true
	ctx.RW.Header().Set("Content-Type", "text/plain")
	ctx.RW.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(
			"attachment; filename=%v_%v_%v.txt",
			args.Category,
			args.FromDate.Format("2006-01-02"),
			args.ToDate.Format("2006-01-02"),
		),
	)

	for i := 0; i < len(Config.LogsConfig); i++ {
		if args.Category == Config.LogsConfig[i].Name {
			Config.LogsConfig[i].Download(args.FromDate, args.ToDate, ctx.RW)
			break
		}
	}
}
