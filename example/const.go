package main

type User struct {
	Id                 int    `json:"id"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Balance            int    `json:"balance"`
	Gay                bool   `json:"gay"`
	Lox                bool   `json:"lox"`
	HavePermission     uint64 `json:"havePermission"`
	OverridePermission uint64 `json:"overridePermission"`
}

type TemplateBlock struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type TemplateHeader struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Template struct {
	Id        int             `json:"id"`
	Header    TemplateHeader  `json:"header"`
	BlockList []TemplateBlock `json:"block_list"`
}
