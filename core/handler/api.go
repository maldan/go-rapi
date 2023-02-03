package handler

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-rapi/rapi_config"
	"github.com/maldan/go-rapi/rapi_error"
	"io/ioutil"
	"strings"
)

type API struct {
	ControllerList []any
}

func (a API) Handle(args rapi_config.HandlerArgs) {
	// Get authorization
	authorization := args.Request.Header.Get("Authorization")
	authorization = strings.Replace(authorization, "Token ", "", 1)

	// Collect params
	params := map[string]any{
		"accessToken": authorization,
	}

	// Read url params
	for key, element := range args.Request.URL.Query() {
		params[key] = element[0]
	}

	// Read body
	bodyBytes, _ := ioutil.ReadAll(args.Request.Body)

	// Parse json body and
	jsonMap := map[string]any{}
	err := json.Unmarshal(bodyBytes, &jsonMap)
	rapi_error.FatalIfError(err)

	// Collect params
	for key, element := range jsonMap {
		params[key] = element
	}

	// Compress all params again
	b, err := json.Marshal(params)
	rapi_error.FatalIfError(err)

	fmt.Printf("%v", string(b))
}
