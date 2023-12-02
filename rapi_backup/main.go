package rapi_backup

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Id        string           `json:"id"`
	Src       string           `json:"src"`
	Dst       string           `json:"dst"`
	Status    string           `json:"status"`
	Period    string           `json:"period"`
	NextRun   time.Time        `json:"nextRun"`
	LastRun   time.Time        `json:"lastRun"`
	BeforeRun func(task *Task) `json:"-"`
	AfterRun  func(task *Task) `json:"-"`
}

func (t *Task) Exec(args ...string) {
	so, se, err := Exec(args...)
	if so != "" {
		fmt.Printf("[TASK EXEC STDOUT] - %v\n", so)
	}
	if se != "" {
		fmt.Printf("[TASK EXEC STDERR] - %v\n", se)
	}
	if err != nil {
		fmt.Printf("[TASK EXEC ERR] - %v\n", err)
	}
}

func (t *Task) DoRsync() string {
	t.Status = "progress"

	src := strings.ReplaceAll(t.Src, "%date%", time.Now().Format("2006-01-02"))
	dst := strings.ReplaceAll(t.Dst, "%date%", time.Now().Format("2006-01-02"))

	so, se, err := Exec("rsync", "-ra", src, dst)
	fmt.Printf("[BACKUP STDOUT] - %v\n", so)
	fmt.Printf("[BACKUP STDERR] - %v\n", se)
	fmt.Printf("[BACKUP ERR] - %v\n", err)

	return dst
}

func (t *Task) Start() {
	t.Status = "start"
}

func (t *Task) Done() {
	t.Status = "done"
}

func (t *Task) GetDestination() string {
	return strings.ReplaceAll(t.Dst, "%date%", time.Now().Format("2006-01-02"))
}

func (t *Task) IsReady() bool {
	return time.Now().Unix() >= t.NextRun.Unix()
}
