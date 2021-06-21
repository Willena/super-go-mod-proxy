package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/willena/super-go-mod-proxy/gomodule"
	"net/http"
)

func moduleFromRequest(r *http.Request) (*gomodule.GoModule, error) {
	version := "master"
	if v, ok := mux.Vars(r)["moduleVersion"]; ok {
		version = v
	}
	if module, ok := mux.Vars(r)["module"]; ok {
		return gomodule.NewGoModule(module, version), nil
	}

	return nil, fmt.Errorf("Privided module name is not valid")
}
