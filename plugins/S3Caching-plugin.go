package plugins

import (
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/types"
	"net/http"
)

func NewS3CachingPlugin() types.PluginInstance {
	return &S3CachingPlugin{}
}

type S3CachingPlugin struct {
}

func (receiver *S3CachingPlugin) Configure(phase types.Phase, config config.PluginConfiguration) types.PluginInstance {
	return receiver
}
func (receiver *S3CachingPlugin) RunPhase(context *types.RunnerContext, w http.ResponseWriter) bool {
	return false
}
