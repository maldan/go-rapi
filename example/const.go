package main

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	BlockList []TemplateBlock `json:"blockList"`
}
