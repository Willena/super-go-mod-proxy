package plugins

import (
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/types"
	"net/http"
)

func NewDiskCachingPlugin() types.PluginInstance {
	return &DiscCachingPlugin{}
}

type DiscCachingPlugin struct {
}

func (receiver *DiscCachingPlugin) Configure(phase types.Phase, config config.PluginConfiguration) types.PluginInstance {
	return receiver
}
func (receiver *DiscCachingPlugin) RunPhase(context *types.RunnerContext, w http.ResponseWriter) bool {
	return false
}
