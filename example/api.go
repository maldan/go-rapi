package main

import (
	"fmt"
	"time"

	"github.com/maldan/go-rapi/rapi_core"
)

type TestApi struct{}
type Test2Api struct{}

type Sas struct {
	X string
}

type XArgs struct {
	A string `validation:"required" json:"a"`
	B int
	C Sas
	D []string
	E []int
	F map[string]string
	G bool
	H time.Time
}

type ArgsId struct {
	Id string
}

type ArgsFile struct {
	Sus rapi_core.File
}

func (u TestApi) GetSasageo(args ArgsId) string {
	fmt.Println(6, args)
	return "99"
}

func (u TestApi) PostSasageo(args XArgs) string {
	fmt.Println(6, args)
	return "99"
}

func (u TestApi) PatchSasageo(args XArgs) string {
	fmt.Println(6, args)
	return "99"
}

func (u TestApi) DeleteSasageo(args ArgsId) string {
	fmt.Println(6, args)
	return "99"
}

func (u Test2Api) GetSasageo(args ArgsId) string {
	fmt.Println(6, args)
	return "99"
}

func (u Test2Api) PostSasageo(args ArgsFile) string {
	fmt.Println(args.Sus.Name)
	return "99"
}

func (u Test2Api) PatchSasageo(args XArgs) string {
	fmt.Println(6, args)
	return "99"
}

func (u Test2Api) DeleteSasageo(args ArgsId) string {
	fmt.Println(6, args)
	return "99"
}
