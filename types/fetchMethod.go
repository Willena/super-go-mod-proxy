package types

import "io"

type FetchMethod interface {
	GetVersions(module string) ([]string, error)
	GetLatestVersion(module string) (string, error)
	GetModule(module string, version string) (string, error)
	GetVersionInfo(module string, version string) (string, error)
	GetZipFile(module string, version string) (io.Reader, error)
	Match(url string) bool
}
