package plugins

import (
	"encoding/json"
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/fetchMethods"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

func NewVcsPlugin() types.PluginInstance {
	return &VcsPlugin{}
}

type VcsPlugin struct {
	phase           types.Phase
	vcsReplacements map[*regexp.Regexp]replacementConfig
}

type replacementConfig struct {
	Url               string                  `json:"url"`
	AuthConfiguration types.AuthConfiguration `json:"auth,omitempty"`
}

func mapToReplacementConfig(config map[string]interface{}) replacementConfig {

	url := config["url"].(string)
	auth := config["auth"].(map[string]interface{})

	return replacementConfig{
		Url: url,
		AuthConfiguration: types.AuthConfiguration{
			Type:       auth["type"].(string),
			Username:   auth["username"].(string),
			Password:   auth["password"].(string),
			PrivateKey: auth["privateKey"].(string),
		},
	}
}

func (receiver *VcsPlugin) Configure(phase types.Phase, config config.PluginConfiguration) types.PluginInstance {
	receiver.vcsReplacements = make(map[*regexp.Regexp]replacementConfig)

	receiver.phase = phase
	for i, val := range config["modules"].(map[string]interface{}) {
		var replacementConfigVar replacementConfig

		d, _ := json.Marshal(val)
		json.Unmarshal(d, &replacementConfigVar)

		re := regexp.MustCompile(i)
		logger.With(zap.Any("replace", re), zap.String("with", replacementConfigVar.Url)).Info("Configured replacement ")
		receiver.vcsReplacements[re] = replacementConfigVar
	}
	return receiver
}
func (receiver *VcsPlugin) RunPhase(context *types.RunnerContext, w http.ResponseWriter) bool {

	switch receiver.phase {
	case types.PhasePreFetch:

		for r, s := range receiver.vcsReplacements {
			if r.MatchString(context.GoModule.Path) {
				replacement := r.ReplaceAllString(context.GoModule.Path, s.Url)
				logger.
					With(zap.String("module", context.GoModule.Path), zap.Any("regex", r), zap.String("replacement", replacement)).
					Info("VCS URL has been updated")

				method, err := fetchMethods.FindForUrl(replacement, s.AuthConfiguration)
				if err != nil {
					logger.
						With(zap.String("module", context.GoModule.Path), zap.String("url", replacement)).
						Error("Could not find suitable fetch method ! ", zap.Error(err))
				}

				context.FetchMethod = method
			}
		}
		return false
	}

	return false
}
