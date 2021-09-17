package rapi_core

import (
	"net/http"
	"strings"
)

func GetHandler(url string, routers map[string]Handler) (string, Handler) {
	var most string
	var handler Handler

	for k, v := range routers {
		if strings.HasPrefix(url, k) {
			if len(most) < len(k) {
				most = k
				handler = v
			}
		}
	}

	return most, handler
}

func DisableCors(rw http.ResponseWriter) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Methods", "*")
	rw.Header().Add("Access-Control-allow-Headers", "*")
}
