package plugins

import (
	"fmt"
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/errors"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func NewDefaultPlugin() types.PluginInstance {
	return &DefaultPlugin{}
}

type Default struct {
	DefaultMethod *types.FetchMethod
}

func (p Default) FetchUrl() string {
	panic("implement me")
}

type DefaultPlugin struct {
	phase types.Phase
}

func (receiver *DefaultPlugin) Configure(phase types.Phase, config config.PluginConfiguration) types.PluginInstance {
	receiver.phase = phase
	return receiver
}

func (receiver *DefaultPlugin) RunPhase(context *types.RunnerContext, w http.ResponseWriter) bool {
	logger.Debug("Running default phase...")
	switch receiver.phase {
	case types.PhaseReceive:
		return doReceiver(context, w)
	case types.PhasePreFetch:
		return doPrefetch(context, w)
	case types.PhaseFetch:
		return doFetch(context, w)
	}
	return false
}

func doFetch(context *types.RunnerContext, w http.ResponseWriter) bool {
	logger.Debug("Running default->Fetch phase...")
	switch context.Action {
	case types.ActionListVersion:
		version, err := context.FetchMethod.GetVersions(context.GoModule)
		if err != nil {
			logger.Error("Error while fetching versions", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			w.Write(errors.GenerateError(err))
			return true
		}

		w.WriteHeader(http.StatusOK)
		for _, v := range version {
			w.Write([]byte(fmt.Sprintf("%s\n", v)))
		}
		return true
	case types.ActionGetVersionInfo:
		json, err := context.FetchMethod.GetVersionInfo(context.GoModule)
		if err != nil {
			logger.Error("Error while fetching version information", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			w.Write(errors.GenerateError(err))
			return true
		}

		//Forward json directly !
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))

		return true
	case types.ActionGetModFile:
		text, err := context.FetchMethod.GetModule(nil)
		if err != nil {
			logger.Error("Error while fetching version information", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			w.Write(errors.GenerateError(err))
			return true
		}

		//Forward json directly !
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))

		return true
	case types.ActionGetModuleZip:
		read, err := context.FetchMethod.GetZipFile(context.GoModule)
		if err != nil {
			logger.Error("Error while fetching version information", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			w.Write(errors.GenerateError(err))
			return true
		}

		//Forward json directly !
		w.Header().Add("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, read)
		return true
	case types.ActionGetLatestVersion:
		json, err := context.FetchMethod.GetLatestVersion(context.GoModule)
		if err != nil {
			logger.Error("Error while fetching version information", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			w.Write(errors.GenerateError(err))
			return true
		}

		//Forward json directly !
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
		return true
	}
	return false
}

func doPrefetch(context *types.RunnerContext, w http.ResponseWriter) bool {
	logger.Debug("Running default->prefetch phase...")
	return false
}

func doReceiver(context *types.RunnerContext, w http.ResponseWriter) bool {
	logger.Debug("Running default->receive phase...")
	return false
}
