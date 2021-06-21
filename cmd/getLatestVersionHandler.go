package main

import (
	"github.com/willena/super-go-mod-proxy/errors"
	"github.com/willena/super-go-mod-proxy/fetchMethods"
	"github.com/willena/super-go-mod-proxy/runner"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"net/http"
)

func LatestVersionHandler(writer http.ResponseWriter, request *http.Request) {
	module, err := moduleName(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(errors.GenerateError(err))
	}

	logger.Debug("Getting moduleVersion info for module ", zap.String("moduleName", module), zap.String("version", "latest"))

	err = runner.NewRunner(&types.RunnerContext{
		GoModule:    module,
		Version:     "latest",
		FetchMethod: &fetchMethods.GoProxy{Url: mainConfig.General.DefaultRelayProxy},
		Action:      types.ActionGetLatestVersion,
	}, pluginsInstances).Run(writer)

	if err != nil {
		logger.Error("Error while collecting versions info for module", zap.String("module", module), zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
