package rapi_panel

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-rapi/rapi_backup"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

type BackupConfig struct {
	IsRun       bool
	HistoryFile string
	TaskList    []rapi_backup.Task `json:"taskList"`
}

type LogConfig struct {
	Name     string                                                         `json:"name"`
	Download func(from time.Time, to time.Time, writer http.ResponseWriter) `json:"-"`
}

func (b *BackupConfig) Run() {
	if runtime.GOOS == "windows" {
		fmt.Printf("[BACKUP SCHEDULE] FAILS. CURRENTLY ONLY LINUX OS SUPPORTED\n")
		return
	}

	fmt.Printf("[BACKUP SCHEDULE] START\n")

	b.ReadHistory()

	// Infinity loop
	for {
		for i := 0; i < len(b.TaskList); i++ {
			task := &(b.TaskList[i])
			if !task.IsReady() {
				continue
			}
			t := time.Now()
			if task.BeforeRun != nil {
				fmt.Printf("[BACKUP TASK BEFORE RUN] Id: %v\n", task.Id)
				err := task.BeforeRun(task)
				if err != nil {
					task.Status = "error"
					task.Error = err.Error()
					fmt.Printf("[BACKUP TASK BEFORE RUN ERR] Id: %v\n", err)
					continue
				}
			}
			fmt.Printf("[BACKUP TASK START] Id: %v\n", task.Id)
			task.Start()
			task.DoRsync()

			if task.Status == "error" {
				continue
			}

			task.Done()
			fmt.Printf("[BACKUP TASK DONE] Id: %v | Time: %v\n", task.Id, time.Since(t))
			if task.AfterRun != nil {
				err := task.AfterRun(task)
				if err != nil {
					task.Status = "error"
					task.Error = err.Error()
					fmt.Printf("[BACKUP TASK AFTER RUN ERR] Id: %v\n", err)
					continue
				}
				fmt.Printf("[BACKUP TASK AFTER DONE] Id: %v | Time: %v\n", task.Id, time.Since(t))
			}

			// Clean resource after work
			task.Clean()

			// Calculate next run
			periods := strings.Split(task.Period, " ")
			nextRun := time.Now()
			for _, period := range periods {
				if strings.Contains(period, "m") {
					periodI, err := strconv.Atoi(period[:len(period)-1])
					if err != nil {
						fmt.Printf("[TASK PARSE PERIOD ERR] %v\n", err)
					}
					nextRun = nextRun.Add(time.Minute * time.Duration(periodI))
				} else if strings.Contains(period, "h") {
					periodI, err := strconv.Atoi(period[:len(period)-1])
					if err != nil {
						fmt.Printf("[TASK PARSE PERIOD ERR] %v\n", err)
					}
					nextRun = nextRun.Add(time.Hour * time.Duration(periodI))
				} else {
					fmt.Printf("[TASK PARSE PERIOD ERR] %v\n", "unknown period")
				}
			}
			task.NextRun = nextRun
			task.LastRun = time.Now()

			b.WriteHistory()
		}

		// Each minute check task
		time.Sleep(time.Minute)
	}
}

func (b *BackupConfig) ReadHistory() {
	// Read config
	data, err := os.ReadFile(b.HistoryFile)
	if err != nil {
		fmt.Printf("[BACKUP HISTORY LOAD ERR] %v\n", err)
		return
	}
	v := map[string]rapi_backup.Task{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("[BACKUP HISTORY LOAD ERR] %v\n", err)
		return
	}

	for i := 0; i < len(b.TaskList); i++ {
		// Read last run
		vv, ok := v[b.TaskList[i].Id]
		if ok {
			b.TaskList[i].NextRun = vv.NextRun
			b.TaskList[i].LastRun = vv.LastRun
		}

		fmt.Printf(
			"Task: %v | LastRun: %v | NextRun: %v\n",
			b.TaskList[i].Id, b.TaskList[i].LastRun, b.TaskList[i].NextRun,
		)
	}
}

func (b *BackupConfig) WriteHistory() {
	v := map[string]any{}

	for i := 0; i < len(b.TaskList); i++ {
		v[b.TaskList[i].Id] = map[string]any{
			"nextRun": b.TaskList[i].NextRun,
			"lastRun": b.TaskList[i].LastRun,
		}
	}

	// Write back
	data, _ := json.Marshal(v)
	os.MkdirAll(path.Dir(b.HistoryFile), 0777)
	err := os.WriteFile(b.HistoryFile, data, 0777)
	if err != nil {
		fmt.Printf("[BACKUP HISTORY WRITE ERR] %v\n", err)
	}
}

type PanelConfig struct {
	CommandList  []PanelCommand
	ChartList    []ChartCommand
	DataAccess   map[string]map[string]func(DataArgs) any
	Password     string
	BackupConfig BackupConfig
	LogsConfig   []LogConfig
}

var Config PanelConfig
