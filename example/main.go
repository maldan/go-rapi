package main

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-rapi"
	"github.com/maldan/go-rapi/core/handler"
	"github.com/maldan/go-rapi/rapi_backup"
	"github.com/maldan/go-rapi/rapi_config"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_error"
	"github.com/maldan/go-rapi/rapi_file"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_panel"
	"github.com/maldan/go-rapi/rapi_rest"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"syscall"
	"time"

	_ "modernc.org/sqlite"
)

var list = make([]User, 0)

func handlerx(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		time.Sleep(time.Second * 1)
		os.Exit(0)
	} else {
		fmt.Println("Ignoring signal: ", signal)
	}
}

func m2() {
	x := rapi_config.Config{
		Router: []rapi_config.RouteHandler{
			{
				Path: "/api",
				Handler: handler.API{
					ControllerList: []any{UserApi{}, TemplateApi{}},
				},
			},
			{
				Path: "/data",
				Handler: handler.FS{
					ContentPath: "example",
				},
			},
			{
				Path:    "/",
				Handler: handler.VFS{},
			},
		},
	}
	fmt.Printf("%v\n", reflect.TypeOf(UserApi{}).Name())
	fmt.Printf("%v\n", x)
}

func main() {
	rapi_log.Info("Fuck")
	rapi_log.Info("Suck")
	rapi_log.Error("Oak")

	/*m2()

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	go func() {
		for {
			s := <-sigchnl
			handlerx(s)
		}
	}()*/

	// Test
	for i := 0; i < 40000; i++ {
		list = append(list, User{
			Id: i + 1, HavePermission: 3, Email: fmt.Sprintf("lox_%v", i), Password: cmhp_crypto.UID(32),
			Created: time.Now(),
		})
	}

	/*rapi.Start2(rapi_config.Config{
		Host: "127.0.0.1:16000",
		Router: []rapi_config.RouteHandler{
			{
				Path: "/api",
				Handler: handler.API{
					ControllerList: []any{UserApi{}, TemplateApi{}},
				},
			},
			{
				Path: "/data",
				Handler: handler.FS{
					ContentPath: "example",
				},
			},
			{
				Path:    "/",
				Handler: handler.VFS{},
			},
		},
	})*/

	rapi.Start(rapi.Config{
		Host: "127.0.0.1:16000",
		Router: map[string]rapi_core.Handler{
			"/": rapi_file.FileHandler{Root: "@"},
			"/api": rapi_rest.ApiHandler{
				Controller: map[string]interface{}{
					"user":     UserApi{},
					"template": TemplateApi{},
				},
			},
		},
		DisableJsonWrapper: true,
		DebugMode:          true,
		PanelConfig: rapi_panel.PanelConfig{
			CommandList: []rapi_panel.PanelCommand{
				{
					Folder: "backup", Name: "sas", Func: func(s string) any {
						fmt.Printf("%v", "gas")
						time.Sleep(time.Second)
						return ""
					},
				},
				{
					Folder: "test", Name: "sas", Func: func(s string) any {
						fmt.Printf("%v", "xxx")
						return ""
					},
				},
				{
					Folder: "test", Name: "Generate", Func: func(s string) any {
						fmt.Printf("%v\n", "xxxxax")
						return ""
					},
				},
			},
			DataAccess: map[string]map[string]func(rapi_panel.DataArgs) any{
				"test": TestAccess,
			},
			ChartList: []rapi_panel.ChartCommand{
				{
					Folder: "main",
					Name:   "test",
					Func: func(s string) []any {
						list2 := make([]any, 0)
						s1 := rand.NewSource(time.Now().UnixNano())
						r1 := rand.New(s1)

						for i := 0; i < 128; i++ {
							list2 = append(list2, map[string]any{
								"date":  "2012-01-01",
								"value": r1.Float32(),
							})
						}

						return list2
					},
				},
				{
					Folder: "main",
					Name:   "test2",
					Func: func(s string) []any {
						return []any{
							map[string]any{
								"value": 1,
							},
							map[string]any{
								"value": 2,
							},
							map[string]any{
								"value": 3,
							},
						}
					},
				},
				{
					Folder: "main2",
					Name:   "test2",
					Func: func(s string) []any {
						return []any{
							map[string]any{
								"value": 1,
							},
							map[string]any{
								"value": 2,
							},
							map[string]any{
								"value": 3,
							},
						}
					},
				},
			},
			Password: "petux",
			BackupConfig: rapi_panel.BackupConfig{
				HistoryFile: "./backup_history.json",
				TaskList: []rapi_backup.Task{
					{
						Id:        "main_db",
						IsVirtual: true,

						// Src:    []string{"./db/."},

						Dst:    []string{"./backup/"},
						Period: "1m",
						BeforeRun: func(task *rapi_backup.Task) error {
							name := fmt.Sprintf("./main_%v.tar.gz", time.Now().Format("2006-01-02_15_04_05"))
							task.Exec("tar", "-czf", name, "./db/.")
							task.SrcVirtual = []string{name}
							task.RemoveQueue = []string{name}

							return nil
						},
					},
				},
			},
			LogsConfig: []rapi_panel.LogConfig{
				{
					Name: "Test",
					Download: func(from time.Time, to time.Time, writer http.ResponseWriter) {
						fmt.Printf("%v\n", from)
						fmt.Printf("%v\n", to)
						rapi_error.Fatal(rapi_error.Error{Description: "gas"})
						writer.Write([]byte("gas"))
					},
				},
				{
					Name: "Rock",
					Download: func(from time.Time, to time.Time, writer http.ResponseWriter) {
						writer.Write([]byte("gas"))
					},
				},
			},
		},
	})
}
