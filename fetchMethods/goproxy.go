package fetchMethods

import (
	"fmt"
	"github.com/willena/super-go-mod-proxy/gomodule"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type GoProxy struct {
	Url string
}

func (g *GoProxy) GetVersions(module *gomodule.GoModule) ([]string, error) {

	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/list", g.Url, module.Path))
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.Error(err)).Error("Could not fetch module versions")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no module version found for %s; request status: %d;", module.Path, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.Error(err)).Error("Could not read proxy response body")
	}

	sortedVersion := sortTags(strings.Split(string(body), "\n"))
	logger.With(zap.String("module", module.Path)).Debug("Found versions ", zap.Any("version", sortedVersion))
	return sortedVersion, nil
}

func (g *GoProxy) GetLatestVersion(module *gomodule.GoModule) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@latest", g.Url, module.Path))
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", "latest"), zap.Error(err)).Error("Could not fetch module version")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("no module version found for %s:%s; request status: %d;", module.Path, "latest", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", "latest"), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	logger.With(zap.String("module", module.Path)).Debug("Found version ", zap.Any("version", "latest"), zap.String("content", s))
	return s, nil
}

func (g *GoProxy) GetModule(module *gomodule.GoModule) (string, error) {
	url := fmt.Sprintf("%s/%s/@v/%s.mod", g.Url, module.Path, module.Version.String())
	logger.With(zap.String("ModuleUrl", url)).Debug("Calling go proxy...")
	resp, err := http.Get(url)
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", module.Version.String()), zap.Error(err)).Error("Could not fetch module go.mod file")
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("no go mod file found for %s:%s; request status: %d;", module.Path, module.Version.String(), resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", module.Version.String()), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	logger.With(zap.String("module", module.Path)).Debug("Found version ", zap.Any("version", module.Version.String()), zap.String("content", s))
	return s, nil
}

func (g *GoProxy) GetVersionInfo(module *gomodule.GoModule) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/%s.info", g.Url, module.Path, module.Version.String()))
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", module.Version.String()), zap.Error(err)).Error("Could not fetch module version")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("no module version found for %s:%s; request status: %d;", module.Path, module.Version.String(), resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", module.Version.String()), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	logger.With(zap.String("module", module.Path)).Debug("Found version ", zap.Any("version", module.Version.String()), zap.String("content", s))
	return s, nil
}

func (g *GoProxy) GetZipFile(module *gomodule.GoModule) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/%s.zip", g.Url, module.Path, module.Version.String()))
	if err != nil {
		logger.With(zap.String("module", module.Path), zap.String("version", module.Version.String()), zap.Error(err)).Error("Could not fetch module version")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Zip not found for %s:%s; request status: %d;", module.Path, module.Version.String(), resp.StatusCode)
	}

	return resp.Body, nil
}

func (g *GoProxy) Match(url string) bool {
	return true
}
