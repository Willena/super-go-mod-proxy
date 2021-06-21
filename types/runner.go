package types

import "github.com/willena/super-go-mod-proxy/gomodule"

type RunnerContext struct {
	GoModule    *gomodule.GoModule
	FetchMethod FetchMethod
	Action      Action
}

type Phase int

const (
	PhaseReceive Phase = iota
	PhasePreFetch
	PhaseFetch
)

type Action int

const (
	ActionListVersion Action = iota
	ActionGetLatestVersion
	ActionGetModuleZip
	ActionGetVersionInfo
	ActionGetModFile
)
