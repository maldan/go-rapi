package rapi_compiler

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_string"
	"github.com/maldan/go-rapi/rapi_test"
	"regexp"
	"sort"
	"strings"
)

func getTestCase(input string) rapi_test.RapiTestCase {
	testCase := rapi_test.RapiTestCase{
		MustStatus: 200,
	}

	// Parse args
	r := regexp.MustCompile("(?U)ARG=({(.*)});")
	match := r.FindStringSubmatch(input)
	if len(match) > 1 {
		json.Unmarshal([]byte(match[1]), &testCase.Params)
	}

	// Parse save
	r = regexp.MustCompile("(?U)SV=(.*);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		json.Unmarshal([]byte(match[1]), &testCase.SaveFields)
	}

	// Parse response fields
	r = regexp.MustCompile("(?U)RF=(.*);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		json.Unmarshal([]byte(match[1]), &testCase.ResponseFields)
	}

	// Parse response fields
	r = regexp.MustCompile("(?U)RF=(.*);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		json.Unmarshal([]byte(match[1]), &testCase.ResponseFields)
	}

	// Parse response values
	r = regexp.MustCompile("(?U)RV=(.*);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		json.Unmarshal([]byte(match[1]), &testCase.ResponseValues)
	}

	// Parse status code
	r = regexp.MustCompile("(?U)SC=(\\d+);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		testCase.MustStatus = cmhp_convert.StrToInt(match[1])
	}

	// Parse priority
	r = regexp.MustCompile("(?U)PR=(\\d+);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		testCase.Priority = cmhp_convert.StrToInt(match[1])
	}

	// Access token from
	r = regexp.MustCompile("(?U)AT=([a-zA-Z]+);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		testCase.AccessTokenFrom = match[1]
	}

	// Access token from
	r = regexp.MustCompile("(?U)METHOD=([a-zA-Z]+);")
	match = r.FindStringSubmatch(input)
	if len(match) > 1 {
		testCase.FullName = match[1]
	}

	return testCase
}

func getTestCaseListFromFile(controller string, t string) []rapi_test.RapiTestCase {
	lines := strings.Split(t, "\n")
	out := make([]rapi_test.RapiTestCase, 0)

	// Get tests
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if !strings.Contains(line, "// @test") {
			continue
		}

		testCase := getTestCase(line)
		testCase.Signature = line

		// Find method
		for j := 0; j < 5; j++ {
			if strings.Contains(lines[i+j], "func (") {
				methodInfo := parseFunctionInfo(lines[i+j])
				testCase.MethodName = methodInfo.MethodName
				testCase.ControllerName = methodInfo.ControllerName
				testCase.HttpMethod = methodInfo.HttpMethod
				testCase.Url = methodInfo.Url

				break
			}
		}

		out = append(out, testCase)
	}

	// Get groups
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.Contains(line, "// @start test group") {
			for j := 0; j < 10; j++ {
				line2 := lines[i+j]

				if strings.Contains(line2, "// METHOD=") {
					testCase := getTestCase(line2)
					methodInfo := parseMethodInfoFromName(testCase.FullName)
					testCase.HttpMethod = methodInfo.HttpMethod
					testCase.MethodName = methodInfo.MethodName
					testCase.ControllerName = controller
					testCase.Signature = line2
					testCase.Url = "/api/" + controller + "/" + cmhp_string.LowerFirst(methodInfo.MethodName)
					out = append(out, testCase)
					continue
				}

				if strings.Contains(line2, "// @end test group") {
					i = i + j
					break
				}
			}
		}
	}

	return out
}

func CompileTests(apiDir string, outPath string) {
	finalOut := `package main
	import(
		"fmt"
		"encoding/json"
		"github.com/maldan/go-rapi/rapi_test"
	)

	var rt = rapi_test.RapiTest{
		Host: "127.0.0.1:9124",
		Storage: map[string]any{},
	}
	`

	list, _ := cmhp_file.ListAll(apiDir)
	finalTestList := make([]rapi_test.RapiTestCase, 0)

	// Parse tests
	for _, l := range list {
		fileText, _ := cmhp_file.ReadText(l.FullPath)
		controllerName := cmhp_string.ToPascalCase(strings.ReplaceAll(l.Name, ".go", ""))
		testList := getTestCaseListFromFile(controllerName, fileText)
		for _, t := range testList {
			finalTestList = append(finalTestList, t)
		}
	}

	sort.Slice(finalTestList, func(i, j int) bool {
		return finalTestList[i].Priority > finalTestList[j].Priority
	})

	// Generate code
	counter := 0
	fnList := make([]string, 0)
	for _, t := range finalTestList {
		fnName := fmt.Sprintf("t_%v_%v_%v_%v", t.HttpMethod, t.ControllerName, t.MethodName, counter)
		finalOut += fmt.Sprintf("func %v() {\n", fnName)
		fnList = append(fnList, fnName)
		counter += 1

		prms, _ := json.Marshal(t.Params)
		var re = regexp.MustCompile(`%([a-zA-Z]+)%`)
		params := re.ReplaceAllString(string(prms), "`+fmt.Sprintf(\"%v\", rt.Storage[\"$1\"])+`")

		responseFields := ""
		for _, rf := range t.ResponseFields {
			responseFields += `"` + rf + `",`
		}

		rv, _ := json.Marshal(t.ResponseValues)
		re = regexp.MustCompile(`%([a-zA-Z]+)%`)
		responseValues := re.ReplaceAllString(string(rv), "`+fmt.Sprintf(\"%v\", rt.Storage[\"$1\"])+`")

		saveFields, _ := json.Marshal(t.SaveFields)

		finalOut += fmt.Sprintf(`params := map[string]any{}
				json.Unmarshal([]byte(%v), &params)

				values := map[string]any{}
				json.Unmarshal([]byte(%v), &values)

				save := map[string]string{}
				json.Unmarshal([]byte(%v), &save)

				rt.BigTest(rapi_test.RapiTestCase{
					HttpMethod: "%v", 
					Url: "%v", 
					MustStatus: %v,
					AccessTokenFrom: "%v",
					ResponseFields: []string{%v},
					ResponseValues: values,
					Params: params,
					Signature: %v,
					TestFunctionName: "%v",
					SaveFields: save,
				})
			`,
			"`"+string(params)+"`",
			"`"+string(responseValues)+"`",
			"`"+string(saveFields)+"`",
			t.HttpMethod,
			t.Url,
			t.MustStatus,
			t.AccessTokenFrom,
			responseFields,
			"`"+t.Signature+"`",
			fnName,
		)

		finalOut += "}\n\n"
	}

	finalOut += "func main() {\n"
	for _, r := range fnList {
		finalOut += r + "()\n"
	}
	finalOut += "}\n"

	cmhp_file.Write(outPath, finalOut)
}
