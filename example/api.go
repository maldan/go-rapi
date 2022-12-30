package main

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
)

type UserApi struct{}
type TemplateApi struct{}

/*type Sas struct {
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

type ArgsFile struct {
	Sus rapi_core.File
}*/

var user User
var templateList = make([]Template, 0)

func (u UserApi) GetIndex() User {
	return user
}

func (u UserApi) GetList() []string {
	return []string{"a", "b", "c"}
}

func (u UserApi) PostIndex(args User) {
	user.Email = args.Email
	user.Password = args.Password
}

func (u UserApi) PatchIndex(args User) {
	user.Email = args.Email
	user.Password = args.Password
}

func (u UserApi) DeleteIndex(args ArgsId) {
	user.Email = ""
	user.Password = ""
}

func (u TemplateApi) PostIndex(args Template) {
	templateList = append(templateList, args)
}

func (u TemplateApi) PostPhoto(args ArgsPhoto) {
	fmt.Printf("%v\n", args.File.Name)
	fmt.Printf("%v\n", args.File.Mime)
	cmhp_file.Write("sas.png", args.File.Data)
}

func (u TemplateApi) GetList() []Template {
	return templateList
}
