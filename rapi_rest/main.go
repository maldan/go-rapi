package rapi_rest

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-rapi/rapi_const"
	"github.com/maldan/go-rapi/rapi_debug"
	"github.com/maldan/go-rapi/rapi_error"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maldan/go-rapi/rapi_core"
)

type ApiHandler struct {
	Controller map[string]interface{}
}

type Response struct {
	Status   bool        `json:"status"`
	Response interface{} `json:"response"`
}

func (r ApiHandler) Handle(args rapi_core.HandlerArgs) {
	// Handle panic
	defer rapi_core.HandleError(&args)

	// Disable cors
	rapi_core.DisableCors(args.RW)

	// Fuck options
	if args.R.Method == "OPTIONS" {
		args.RW.WriteHeader(200)
		fmt.Fprintf(args.RW, "")
		return
	}

	// Get authorization
	authorization := args.R.Header.Get("Authorization")
	authorization = strings.Replace(authorization, "Token ", "", 1)

	// Create context
	remoteIp := ""
	if len(args.R.Header["X-Forwarded-For"]) > 0 {
		remoteIp = args.R.Header["X-Forwarded-For"][0]
	}
	args.Context = &rapi_core.Context{
		AccessToken: authorization,
		RW:          args.RW,
		R:           args.R,
		RemoteIP:    remoteIp,
	}

	// Collect params
	params := map[string]interface{}{
		"accessToken":    authorization,
		"timeZoneOffset": args.R.Header.Get("timeZoneOffset"),
		"buildVersion":   args.R.Header.Get("buildVersion"),
	}

	for key, element := range args.R.URL.Query() {
		params[key] = element[0]
	}

	// Parse body
	if strings.Contains(args.R.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		err := args.R.ParseForm()
		rapi_error.FatalIfError(err)

		for key, value := range args.R.PostForm {
			params[key] = value[0]
		}
	} else if strings.Contains(args.R.Header.Get("Content-Type"), "multipart/form-data") {
		// Parse multipart body and collect params
		args.R.ParseMultipartForm(0)
		for key, element := range args.R.MultipartForm.Value {
			params[key] = element[0]
		}

		// Collect files
		if len(args.R.MultipartForm.File) > 0 {
			for kk, fileHeaders := range args.R.MultipartForm.File {
				for _, fileHeader := range fileHeaders {
					f, _ := fileHeader.Open()
					buffer := make([]byte, fileHeader.Size)
					f.Read(buffer)
					f.Close()
					params[kk] = rapi_const.File{
						Name: fileHeader.Filename,
						Mime: fileHeader.Header.Get("Content-Type"),
						Size: int(fileHeader.Size),
						Data: buffer,
					}
				}
			}
		}
	} else {
		defer args.R.Body.Close()

		// Read body
		bodyBytes, err := ioutil.ReadAll(args.R.Body)
		if err != nil {
			rapi_error.Fatal(rapi_error.Error{
				Description: err.Error(),
			})
		}
		args.RawBody = bodyBytes

		// Parse json body and collect params
		jsonMap := make(map[string]interface{})
		json.Unmarshal(args.RawBody, &jsonMap)
		for key, element := range jsonMap {
			params[key] = element
		}
	}

	// Set args for debug
	if args.DebugMode {
		args.MethodArgs = params
	}

	// Get controller
	path := strings.Split(strings.Replace(args.R.URL.Path, args.Route, "", 1), "/")
	controllerName := path[1]
	methodName := ""

	if len(path) > 2 {
		methodName = path[2]
	}
	if methodName == "" {
		methodName = "Index"
	}

	// Check controller
	_, ok := r.Controller[controllerName]
	if !ok {
		rapi_error.Fatal(rapi_error.Error{
			Code: 404,
			Description: fmt.Sprintf(
				"Controller %v not found",
				controllerName,
			),
		})
	}

	// Get method
	method := GetMethod(r.Controller[controllerName], methodName, args.R.Method)
	if method == nil {
		rapi_error.Fatal(rapi_error.Error{
			Code: 404,
			Description: fmt.Sprintf(
				"Method %v not found in controller %v",
				strings.Title(strings.ToLower(args.R.Method))+strings.Title(methodName),
				controllerName,
			),
		})
	}

	if args.DebugMode {
		debugParams := make(map[string]any)
		for k, v := range params {
			switch v.(type) {
			case rapi_const.File:
				debugParams[k] = rapi_const.File{
					Name: v.(rapi_const.File).Name,
					Mime: v.(rapi_const.File).Mime,
					Size: v.(rapi_const.File).Size,
				}
			default:
				debugParams[k] = v
			}
		}
		d, _ := json.Marshal(debugParams)
		rapi_debug.GetRequestLog(args.Id).SetInput(string(d))
	}

	// Call method
	value := ExecuteMethod(r.Controller[controllerName], args, *method, params)

	// Skip prepare and write
	if args.Context.IsSkipProcessing {
		return
	}

	// If return file path to server
	if args.Context.IsServeFile {
		rapi_debug.GetRequestLog(args.Id).SetResponse("FILE " + value.Interface().(string))
		http.ServeFile(args.RW, args.R, value.Interface().(string))
		return
	}

	// Prepare response
	var res any
	if args.DisableJsonWrapper {
		res = value.Interface()
	} else {
		res = Response{
			Status:   true,
			Response: value.Interface(),
		}
	}
	data, err := json.Marshal(&res)

	if err != nil {
		rapi_error.Fatal(rapi_error.Error{
			Description: err.Error(),
		})
	}

	// Write response
	args.RW.Header().Add("Content-Type", "application/json")
	args.RW.Write(data)
	if args.DebugMode {
		rapi_debug.GetRequestLog(args.Id).SetResponse(string(data))
	}
}
