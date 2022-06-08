package rapi_log

import (
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

var mu sync.Mutex
var logList = make([]LogData, 0)

func Error(message string) {
	mu.Lock()
	defer mu.Unlock()
	logList = append(logList, LogData{
		Kind:    "error",
		Message: message,
		Time:    time.Now(),
	})
}

func Info(message string) {
	mu.Lock()
	defer mu.Unlock()
	logList = append(logList, LogData{
		Kind:    "info",
		Message: message,
		Time:    time.Now(),
	})
}

func (r LogApi) GetIndex() []LogData {
	return logList
}
