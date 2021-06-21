package plugins

import (
	"fmt"
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/fetchMethods"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

func NewPrivatePlugin() types.PluginInstance {
	return &PrivatePlugin{}
}

type PrivatePlugin struct {
	DefaultMethod   *types.FetchMethod
	Phase           types.Phase
	privateMatchers []*regexp.Regexp
}

func (receiver *PrivatePlugin) Configure(phase types.Phase, config config.PluginConfiguration) types.PluginInstance {
	receiver.Phase = phase
	for i, val := range config["modules"].(map[string]interface{}) {
		if val.(bool) {
			logger.Debug("Loaded single matcher for PrivtaePlugin", zap.String("Regexp", i))
			receiver.privateMatchers = append(receiver.privateMatchers, regexp.MustCompile(i))
		}
	}
	return receiver
}

func (receiver *PrivatePlugin) RunPhase(context *types.RunnerContext, w http.ResponseWriter) bool {
	switch receiver.Phase {
	case types.PhasePreFetch:
		for _, v := range receiver.privateMatchers {
			if v.MatchString(context.GoModule.Path) {
				logger.Info("Found matching private module", zap.String("module", context.GoModule.Path), zap.Any("Regexp", v))
				context.FetchMethod = &fetchMethods.Git{Url: fmt.Sprintf("git+https://%s", context.GoModule.Path)}
				return false
			}
		}

		return false
	}

	return false
}
