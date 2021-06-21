package types

import (
	"github.com/willena/super-go-mod-proxy/gomodule"
	"io"
)

type FetchMethod interface {
	GetVersions(module *gomodule.GoModule) ([]string, error)
	GetLatestVersion(module *gomodule.GoModule) (string, error)
	GetModule(module *gomodule.GoModule) (string, error)
	GetVersionInfo(module *gomodule.GoModule) (string, error)
	GetZipFile(module *gomodule.GoModule) (io.Reader, error)
	Match(url string) bool
}
