package main

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-rapi"
	"github.com/maldan/go-rapi/core/handler"
	"github.com/maldan/go-rapi/rapi_config"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_file"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_panel"
	"github.com/maldan/go-rapi/rapi_rest"
	"os"
	"os/signal"
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

	m2()

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	go func() {
		for {
			s := <-sigchnl
			handlerx(s)
		}
	}()

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
		},
	})
}
