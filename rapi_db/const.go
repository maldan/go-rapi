package rapi_db

type FieldInfo struct {
	Name   string `json:"name"`
	IsEdit bool   `json:"isEdit"`
	IsHide bool   `json:"isHide"`
	Type   string `json:"type"`
}

type DataSettings struct {
	IsEditable  bool        `json:"isEditable"`
	IsDeletable bool        `json:"isDeletable"`
	FieldList   []FieldInfo `json:"fieldList"`
}

type IDataBase interface {
	GetSettings() DataSettings
	Search(string, int, int) SearchResult[any]
	GetById(int) any
	UpdateById(int, string)
	DeleteById(int)
}
