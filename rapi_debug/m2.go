package rapi_debug

import (
	"github.com/maldan/go-rapi/rapi_error"
	"github.com/maldan/go-rapi/rapi_panel"
	"sync"
	"time"
)

type LogConfig struct {
	OnRequest func(debugLog RapiRequestLog)
	OnSearch  func(args RapiRequestLogSearchArgs) rapi_panel.SearchResult[RapiRequestLog]
}

type RapiRequestLog struct {
	Id         string    `json:"id"`
	HttpMethod string    `json:"httpMethod"`
	Url        string    `json:"url"`
	Input      string    `json:"input"`
	Response   string    `json:"response"`
	RemoteAddr string    `json:"remoteAddr"`
	StatusCode int       `json:"statusCode"`
	Error      string    `json:"error"`
	Created    time.Time `json:"created"`
}

type RapiRequestLogSearchArgs struct {
	Url     string
	Offset  int
	Limit   int
	Created time.Time
}

func (l *RapiRequestLog) SetRequest(method string, url string) *RapiRequestLog {
	l.HttpMethod = method
	l.Url = url
	return l
}

func (l *RapiRequestLog) SetRemoteAddr(addr string) *RapiRequestLog {
	l.RemoteAddr = addr
	return l
}

func (l *RapiRequestLog) SetInput(input string) {
	l.Input = input
}

func (l *RapiRequestLog) SetResponse(response string) {
	l.Response = response
}

func (l *RapiRequestLog) SetError(err rapi_error.Error) {
	l.StatusCode = err.Code
	l.Error = err.Description
}

var mu2 sync.RWMutex
var LogList2 = make([]*RapiRequestLog, 0)
var LogMap2 = make(map[string]*RapiRequestLog)

func GetRequestLog(id string) *RapiRequestLog {
	mu2.Lock()
	defer mu2.Unlock()

	log, ok := LogMap2[id]
	if ok {
		return log
	}
	log = &RapiRequestLog{
		Id:         id,
		StatusCode: 200,
		Created:    time.Now(),
	}
	LogList2 = append(LogList2, log)
	LogMap2[id] = log
	return log
}
