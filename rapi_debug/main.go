package rapi_debug

import (
	"github.com/maldan/go-rapi/rapi_error"
	"sync"
	"time"
)

type RapiDebugLog struct {
	Id         string           `json:"id"`
	HttpMethod string           `json:"httpMethod"`
	Url        string           `json:"url"`
	Args       map[string]any   `json:"args"`
	Response   any              `json:"response"`
	Error      rapi_error.Error `json:"error"`
	Created    time.Time        `json:"created"`
}

func (l *RapiDebugLog) SetRequest(method string, url string) {
	l.HttpMethod = method
	l.Url = url
}

func (l *RapiDebugLog) SetArgs(args map[string]any) {
	l.Args = args
}

func (l *RapiDebugLog) SetResponse(response any) {
	l.Response = response
}

func (l *RapiDebugLog) SetError(err rapi_error.Error) {
	l.Error = err
}

var mu sync.RWMutex
var LogList = make([]*RapiDebugLog, 0)
var LogMap = make(map[string]*RapiDebugLog)

func Log(id string) *RapiDebugLog {
	mu.Lock()
	defer mu.Unlock()

	log, ok := LogMap[id]
	if ok {
		return log
	}
	log = &RapiDebugLog{
		Id:      id,
		Created: time.Now(),
	}
	LogList = append(LogList, log)
	LogMap[id] = log
	return log
}
