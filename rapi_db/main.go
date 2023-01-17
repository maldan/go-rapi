package rapi_db

type IDataBase interface {
	GetStruct() any
	Search(int, int) []any
	GetById(int) any
	UpdateById(int, any)
	DeleteById(int)
}

type DataApi struct {
}

var DataAccess map[string]IDataBase

type ArgsSearch struct {
	Table  string `json:"table"`
	Id     int    `json:"id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (u DataApi) GetStruct(args ArgsSearch) any {
	return DataAccess[args.Table].GetStruct()
}

func (u DataApi) GetTableList() []string {
	l := make([]string, 0)
	for k, _ := range DataAccess {
		l = append(l, k)
	}
	return l
}

func (u DataApi) GetSearch(args ArgsSearch) []any {
	return DataAccess[args.Table].Search(args.Offset, args.Limit)
}

func (u DataApi) GetById(args ArgsSearch) any {
	return DataAccess[args.Table].GetById(args.Id)
}

func (u DataApi) DeleteById(args ArgsSearch) {
	DataAccess[args.Table].DeleteById(args.Id)
}

/*var DbPath string

func Load(id string, v interface{}) {
	name := strings.ToLower(reflect.TypeOf(v).Elem().Name())
	if !cmhp_file.Exists(DbPath + "/" + name + "/" + id + ".json") {
		rapi_core.Fatal(rapi_core.Error{
			Code:        404,
			Description: fmt.Sprintf("%v with id %v not found", strings.Title(name), id),
		})
	}
	err := cmhp_file.ReadJSON(DbPath+"/"+name+"/"+id+".json", v)
	if err != nil {
		rapi_core.Fatal(rapi_core.Error{
			Description: err.Error(),
		})
	}
}

func Save(v interface{}) {
	name := strings.ToLower(reflect.TypeOf(v).Elem().Name())
	id := cmhp_crypto.UID(10)

	cmhp_reflect.SetField(v, "Id", id)
	cmhp_reflect.SetField(v, "Created", time.Now())

	// Save to file
	err := cmhp_file.Write(DbPath+"/"+name+"/"+id+".json", v)
	if err != nil {
		rapi_core.Fatal(rapi_core.Error{
			Description: err.Error(),
		})
	}
}

func Update(id string, v interface{}) {
	name := strings.ToLower(reflect.TypeOf(v).Elem().Name())

	// Save to file
	err := cmhp_file.Write(DbPath+"/"+name+"/"+id+".json", v)
	if err != nil {
		rapi_core.Fatal(rapi_core.Error{
			Description: err.Error(),
		})
	}
}*/
