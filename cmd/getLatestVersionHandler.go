package main

import (
	"fmt"
	"github.com/willena/super-go-mod-proxy/errors"
	"net/http"
)

func LatestVersionHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(errors.GenerateError(fmt.Errorf("Endpoint not implemented")))
}
