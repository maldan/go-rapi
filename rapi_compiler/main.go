package rapi_compiler

import (
	"github.com/maldan/go-cmhp/cmhp_string"
	"regexp"
	"strings"
)

type RapiMethodInfo struct {
	HttpMethod     string
	ControllerName string
	MethodName     string
	Url            string
}

func parseMethodInfoFromName(methodName string) RapiMethodInfo {
	name := strings.ReplaceAll(methodName, "Post", "")
	name = strings.ReplaceAll(name, "Get", "")
	name = strings.ReplaceAll(name, "Patch", "")
	name = strings.ReplaceAll(name, "Delete", "")
	name = strings.ReplaceAll(name, "Put", "")

	httpMethod := ""
	if strings.Contains(methodName, "Post") {
		httpMethod = "POST"
	}
	if strings.Contains(methodName, "Get") {
		httpMethod = "GET"
	}
	if strings.Contains(methodName, "Patch") {
		httpMethod = "PATCH"
	}
	if strings.Contains(methodName, "Delete") {
		httpMethod = "DELETE"
	}
	if strings.Contains(methodName, "Put") {
		httpMethod = "PUT"
	}

	return RapiMethodInfo{
		HttpMethod: httpMethod,
		MethodName: name,
	}
}

func parseFunctionInfo(input string) RapiMethodInfo {
	out := RapiMethodInfo{}

	r2 := regexp.MustCompile("(?U)func \\([a-z] (.*)\\) ([a-zA-Z]+)\\(")
	match2 := r2.FindStringSubmatch(input)

	controller := strings.ReplaceAll(match2[1], "Api", "")

	methodName := match2[2]
	methodName = strings.ReplaceAll(methodName, "Post", "")
	methodName = strings.ReplaceAll(methodName, "Get", "")
	methodName = strings.ReplaceAll(methodName, "Patch", "")
	methodName = strings.ReplaceAll(methodName, "Delete", "")
	methodName = strings.ReplaceAll(methodName, "Put", "")

	if strings.Contains(match2[2], "Post") {
		out.HttpMethod = "POST"
	}
	if strings.Contains(match2[2], "Get") {
		out.HttpMethod = "GET"
	}
	if strings.Contains(match2[2], "Patch") {
		out.HttpMethod = "PATCH"
	}
	if strings.Contains(match2[2], "Delete") {
		out.HttpMethod = "DELETE"
	}
	if strings.Contains(match2[2], "Put") {
		out.HttpMethod = "PUT"
	}

	out.MethodName = methodName
	out.ControllerName = controller
	out.Url = "/api/" + cmhp_string.LowerFirst(controller) + "/" + cmhp_string.LowerFirst(methodName)

	return out
}
