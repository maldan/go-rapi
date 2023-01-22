package main

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-rapi"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_db"
	"github.com/maldan/go-rapi/rapi_file"
	"github.com/maldan/go-rapi/rapi_log"
	"github.com/maldan/go-rapi/rapi_rest"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var list = make([]User, 0)

func handler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		time.Sleep(time.Second * 2)
		os.Exit(0)
	} else {
		fmt.Println("Ignoring signal: ", signal)
	}
}

func main() {
	rapi_log.Info("Fuck")
	rapi_log.Info("Suck")
	rapi_log.Error("Oak")

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	go func() {
		for {
			s := <-sigchnl
			handler(s)
		}
	}()

	// Test
	for i := 0; i < 40000; i++ {
		list = append(list, User{Id: i + 1, Email: fmt.Sprintf("lox_%v", i), Password: cmhp_crypto.UID(32)})
	}

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
		DataAccess: map[string]rapi_db.IDataBase{
			"test": TestData[User]{},
		},
	})
}
