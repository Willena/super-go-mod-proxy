package types

type RunnerContext struct {
	GoModule    string
	FetchMethod FetchMethod
	Action      Action
	Version     string
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
