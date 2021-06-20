package types

import (
	"github.com/willena/super-go-mod-proxy/config"
	"net/http"
)

type PluginInstances []PluginInstance

type PhasesPluginsInstance struct {
	Receive  PluginInstances
	PreFetch PluginInstances
	Fetch    PluginInstances
}
type PluginInitializationFunction func() PluginInstance

type PluginInstance interface {
	Configure(phase Phase, config config.PluginConfiguration) PluginInstance
	RunPhase(context *RunnerContext, w http.ResponseWriter) bool
}
