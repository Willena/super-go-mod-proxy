package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willena/super-go-mod-proxy/errors"
	"github.com/willena/super-go-mod-proxy/fetchMethods"
	"github.com/willena/super-go-mod-proxy/runner"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"net/http"
)

var (
	listVersionCallCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "list_module_versions",
		Help: "The number of times the list endpoint has been called",
	})
)

func ListVersionHandler(writer http.ResponseWriter, request *http.Request) {
	listVersionCallCounter.Inc()
	module, err := moduleFromRequest(request)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(errors.GenerateError(err))
	}
	logger.Debug("Listing moduleVersion for module ", zap.String("module", module.Path))

	err = runner.NewRunner(&types.RunnerContext{
		GoModule:    module,
		FetchMethod: &fetchMethods.GoProxy{Url: mainConfig.General.DefaultRelayProxy},
		Action:      types.ActionListVersion,
	}, pluginsInstances).Run(writer)

	if err != nil {
		logger.Error("Error while collecting versions list for module", zap.String("module", module.Path), zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
