package plugins

import (
	"fmt"
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"sort"
)

var logger, _ = zap.NewDevelopment()

var availablePlugins = map[string]types.PluginInitializationFunction{
	"rewrite":      NewRewritePlugin,
	"private":      NewPrivatePlugin,
	"vcs":          NewVcsPlugin,
	"s3-caching":   NewS3CachingPlugin,
	"disk-caching": NewDiskCachingPlugin,
	"default":      NewDefaultPlugin,
}

func CreateFromConfig(config *config.Config) *types.PhasesPluginsInstance {
	logger.Info("Instantiate plugins for main Runner")
	var phasesPlugins types.PhasesPluginsInstance

	//For each Phase
	phasesPlugins.Receive = LoadPlugins(types.PhaseReceive, config.Phases.Receive)
	phasesPlugins.PreFetch = LoadPlugins(types.PhasePreFetch, config.Phases.PreFetch)
	phasesPlugins.Fetch = LoadPlugins(types.PhaseFetch, config.Phases.Fetch)
	phasesPlugins.Package = LoadPlugins(types.PhasePackage, config.Phases.Package)
	phasesPlugins.Cache = LoadPlugins(types.PhaseCache, config.Phases.Cache)

	return &phasesPlugins

}

func sortedKeys(mapinfo map[string]config.PluginDefinition) []string {
	keys := make([]string, 0)

	for k := range mapinfo {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func LoadPlugins(phase types.Phase, definitions config.PluginsDefinitions) types.PluginInstances {
	instances := make(types.PluginInstances, 0)
	logger.Info(fmt.Sprintf("Loading plugins for phase %d", phase))

	defaultLoaded := false

	for _, key := range sortedKeys(definitions) {
		def := definitions[key]
		if def.Kind == "default" {
			defaultLoaded = true
		}

		logger.Info("Found plugin", zap.String("name", key), zap.String("Kind", def.Kind))
		if fn, ok := availablePlugins[def.Kind]; ok {
			instances = append(instances, fn().Configure(phase, def.Config))
		} else {
			logger.Warn("Missing plugin Kind ! ", zap.String("name", key), zap.String("Kind", def.Kind))
		}
	}

	if !defaultLoaded {
		instances = append(instances, NewDefaultPlugin().Configure(phase, map[string]interface{}{}))
	}

	return instances
}
