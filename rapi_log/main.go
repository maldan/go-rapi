package rapi_log

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"github.com/maldan/go-cmhp/cmhp_time"
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

var mu sync.Mutex
var logList = make([]LogData, 0)

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

func (r LogApi) GetIndex() []LogData {
	return logList
}

func (r LogApi) GetSearch(args ArgsSearch) []LogData {
	nLogs := logList

	if args.Date.Year() > 1 {
		from := time.Date(args.Date.Year(), args.Date.Month(), args.Date.Day(), 0, 0, 0, 0, args.Date.Location())
		to := time.Date(args.Date.Year(), args.Date.Month(), args.Date.Day(), 23, 59, 59, 0, args.Date.Location())

		nLogs = cmhp_slice.Filter(logList, func(t LogData) bool {
			return t.Time.Unix() >= from.Unix() && t.Time.Unix() <= to.Unix()
		})
	}
	return nLogs
}
