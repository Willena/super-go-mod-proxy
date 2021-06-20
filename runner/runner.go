package runner

import (
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"net/http"
)

var logger, _ = zap.NewDevelopment()

type Runner struct {
	runContext *types.RunnerContext
	plugins    *types.PhasesPluginsInstance
	logger     *zap.Logger
}

func NewRunner(ctx *types.RunnerContext, plugins *types.PhasesPluginsInstance) *Runner {

	return &Runner{runContext: ctx, plugins: plugins, logger: logger.With(zap.Any("context", ctx))}
}

func (r *Runner) runPhase(phase types.Phase, steps types.PluginInstances, w http.ResponseWriter) bool {
	r.logger.Info("Running phase", zap.Any("phase", phase))

	for _, step := range steps {
		if step.RunPhase(r.runContext, w) {
			logger.Info("Last phase terminated the process")
			return true
		}
	}

	return false
}

func (r *Runner) Run(w http.ResponseWriter) error {
	r.logger.Info("Runner started !")

	logger.With(zap.Any("context", r.runContext)).Info("Phase Receive...")
	if r.runPhase(types.PhaseReceive, r.plugins.Receive, w) {
		return nil
	}

	logger.With(zap.Any("context", r.runContext)).Info("Phase Prefetch...")
	if r.runPhase(types.PhasePreFetch, r.plugins.PreFetch, w) {
		return nil
	}

	logger.With(zap.Any("context", r.runContext)).Info("Phase Fetch...")
	if r.runPhase(types.PhaseFetch, r.plugins.Fetch, w) {
		return nil
	}

	return nil
}
