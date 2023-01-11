package rapi_test

import (
	"github.com/maldan/go-cmhp/cmhp_string"
	"reflect"
	"runtime"
	"strings"
)

type TestBlock struct {
	Id          string         `json:"id"`
	Tag         string         `json:"tag"`
	HttpMethod  string         `json:"httpMethod"`
	Url         string         `json:"url"`
	Input       []string       `json:"input"`
	Output      []string       `json:"output"`
	Constraints []string       `json:"constraints"`
	Args        map[string]any `json:"args"`
}

type TestBlockConnection struct {
	FromId    string `json:"fromId"`
	ToId      string `json:"toId"`
	FromField string `json:"fromField"`
	ToField   string `json:"toField"`
}

type TestCase struct {
	Name           string                 `json:"name"`
	Tag            string                 `json:"tag"`
	BlockList      []*TestBlock           `json:"blockList"`
	ConnectionList []*TestBlockConnection `json:"connectionList"`
}

func (t *TestBlock) SetArgs(args map[string]any) *TestBlock {
	t.Args = args
	return t
}

func (t *TestBlock) AddInput(input string) *TestBlock {
	t.Input = append(t.Input, input)
	return t
}

func (t *TestBlock) AddOutput(output string) *TestBlock {
	t.Output = append(t.Output, output)
	return t
}

func (t *TestBlock) AddConstraint(constraint string) *TestBlock {
	t.Constraints = append(t.Constraints, constraint)
	return t
}

func (t *TestCase) Add(id string, fn any) *TestBlock {
	x := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	y := strings.Split(x, "/")
	y = strings.Split(y[len(y)-1], "-")
	y = strings.Split(y[0], ".")

	out := make([]string, 0)
	httpMethod := ""

	for i, zz := range y {
		a := zz

		if i == 1 {
			a = strings.ReplaceAll(a, "Api", "")
		}
		if i == 2 {
			if strings.Contains(a, "Post") {
				httpMethod = "POST"
			}
			if strings.Contains(a, "Get") {
				httpMethod = "GET"
			}
			if strings.Contains(a, "Delete") {
				httpMethod = "DELETE"
			}
			if strings.Contains(a, "Put") {
				httpMethod = "PUT"
			}
			if strings.Contains(a, "Patch") {
				httpMethod = "PATCH"
			}
			a = strings.Replace(a, "Post", "", 1)
			a = strings.Replace(a, "Get", "", 1)
			a = strings.Replace(a, "Patch", "", 1)
			a = strings.Replace(a, "Put", "", 1)
			a = strings.Replace(a, "Delete", "", 1)
		}

		out = append(out, cmhp_string.LowerFirst(a))
	}

	block := TestBlock{
		Id:         t.Tag + "_" + t.Name + "_" + id,
		Url:        "/" + strings.Join(out, "/"),
		HttpMethod: httpMethod,
		Input:      make([]string, 0),
		Output:     make([]string, 0),
		Args:       map[string]any{},
	}

	t.BlockList = append(t.BlockList, &block)

	return &block
}

func (t *TestCase) Connect(fromId string, toId string, fromField string, toField string) {
	t.ConnectionList = append(t.ConnectionList, &TestBlockConnection{
		FromId:    t.Tag + "_" + t.Name + "_" + fromId,
		ToId:      t.Tag + "_" + t.Name + "_" + toId,
		FromField: fromField,
		ToField:   toField,
	})
}

/*type RapiTest struct {
	Host        string
	AccessToken string
	Storage     map[string]any
}

type RapiTestCase struct {
	Params           map[string]any
	HttpMethod       string
	MustStatus       int
	TestFunctionName string
	Signature        string
	Url              string
	FullName         string
	ControllerName   string
	MethodName       string
	Priority         int
	AccessTokenFrom  string
	ResponseFields   []string
	ResponseValues   map[string]any
	SaveFields       map[string]string
}

func (r *RapiTest) BigTest(args RapiTestCase) {
	out := map[string]any{}

	res := cmhp_net.Request(cmhp_net.HttpArgs{
		Url:        "http://" + r.Host + args.Url,
		Headers:    map[string]string{"Authorization": r.AccessToken},
		Method:     args.HttpMethod,
		InputJSON:  args.Params,
		OutputJSON: &out,
	})

	fmt.Printf("Testing %v %v\n", args.HttpMethod, res.Url)
	fmt.Printf("     TestFunctionName: %v\n", args.TestFunctionName)
	fmt.Printf("     Signature: %v\n", args.Signature)
	fmt.Printf("     Status: get %v / need %v\n", res.StatusCode, args.MustStatus)
	fmt.Printf("     Params: %v\n", args.Params)

	for _, v := range args.ResponseFields {
		_, ok := out[v]
		if !ok {
			fmt.Printf("%+v\n", out)
			fmt.Printf("%+v\n", r.Storage)
			fmt.Printf("     Response must have field: %v\n", v)
			os.Exit(1)
		}
	}

	for k, v := range args.ResponseValues {
		if fmt.Sprintf("%v", out[k]) != fmt.Sprintf("%v", v) {
			fmt.Printf("%+v\n", out)
			fmt.Printf("%+v\n", r.Storage)
			fmt.Printf("     Response field expected: %v but received : %v\n", v, out[k])
			os.Exit(1)
		}
	}

	for k, v := range args.SaveFields {
		r.Storage[v] = out[k]
	}

	// Check status
	if res.StatusCode != args.MustStatus {
		fmt.Printf("%+v\n", out)
		fmt.Printf("%+v\n", r.Storage)
		fmt.Printf("     ERR\n\n")
		os.Exit(1)
	}

	// Set access token
	if args.AccessTokenFrom != "" {
		r.AccessToken = out[args.AccessTokenFrom].(string)
	}

	fmt.Printf("     OK\n\n")
}*/

/*func (r *RapiTest) MustOkJson(method string, url string, args map[string]any, ok func(r map[string]any) bool) {
	out := map[string]any{"token": ""}
	res := cmhp_net.Request(cmhp_net.HttpArgs{
		Url:        "http://" + r.Host + url,
		Headers:    map[string]string{"Authorization": r.AccessToken},
		Method:     method,
		InputJSON:  args,
		OutputJSON: &out,
	})
	if res.StatusCode != 200 {
		panic(fmt.Sprintf("%v %v - FAIL\n", method, url))
	}
	if ok != nil && !ok(out) {
		panic(fmt.Sprintf("%v %v - FAIL\n", method, url))
	}

	fmt.Printf("%v %v - OK\n", method, url)
}*/
