package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func moduleName(r *http.Request) (string, error) {
	if module, ok := mux.Vars(r)["module"]; ok {
		return module, nil
	}

	return "", fmt.Errorf("Privided module name is not valid")
}

func moduleVersion(r *http.Request) (string, error) {
	if module, ok := mux.Vars(r)["moduleVersion"]; ok {
		return module, nil
	}

	return "", fmt.Errorf("Privided moduleVersion is not valid")
}
