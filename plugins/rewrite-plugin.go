package plugins

import (
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

func NewRewritePlugin() types.PluginInstance {
	return &RewritePlugin{}
}

type RewritePlugin struct {
	Phase               types.Phase
	rewriteReplacements map[*regexp.Regexp]string
}

func (receiver *RewritePlugin) Configure(phase types.Phase, config config.PluginConfiguration) types.PluginInstance {
	receiver.rewriteReplacements = make(map[*regexp.Regexp]string)
	receiver.Phase = phase
	for i, val := range config["modules"].(map[string]interface{}) {
		re := regexp.MustCompile(i)
		logger.With(zap.Any("replace", re), zap.String("with", val.(string))).Info("Configured replacement ")
		receiver.rewriteReplacements[re] = val.(string)
	}
	return receiver
}
func (receiver *RewritePlugin) RunPhase(context *types.RunnerContext, w http.ResponseWriter) bool {

	switch receiver.Phase {
	case types.PhasePreFetch:

		for regex, replacement := range receiver.rewriteReplacements {
			old := context.GoModule
			context.GoModule = regex.ReplaceAllString(context.GoModule, replacement)
			if old != context.GoModule {
				logger.
					With(zap.String("from", old), zap.String("to", context.GoModule),
						zap.Any("withFind", regex), zap.String("replacement", replacement)).
					Debug("Rewritten gomodule")
			}

		}

		return false
	}

	return false
}
