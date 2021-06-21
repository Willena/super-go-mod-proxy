package fetchMethods

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type GoProxy struct {
	Url string
}

func (g *GoProxy) GetVersions(module string) ([]string, error) {

	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/list", g.Url, module))
	if err != nil {
		logger.With(zap.String("module", module), zap.Error(err)).Error("Could not fetch module versions")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no module version found for %s; request status: %d;", module, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	versions := strings.Split(s, "\n")
	logger.With(zap.String("module", module)).Debug("Found versions ", zap.Any("version", versions))
	return versions, nil
}

func (g *GoProxy) GetLatestVersion(module string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@latest", g.Url, module))
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", "latest"), zap.Error(err)).Error("Could not fetch module version")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("no module version found for %s:%s; request status: %d;", module, "latest", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", "latest"), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	logger.With(zap.String("module", module)).Debug("Found version ", zap.Any("version", "latest"), zap.String("content", s))
	return s, nil
}

func (g *GoProxy) GetModule(module string, version string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/%s.mod", g.Url, module, version))
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", version), zap.Error(err)).Error("Could not fetch module go.mod file")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("no go mod file found for %s:%s; request status: %d;", module, version, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", version), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	logger.With(zap.String("module", module)).Debug("Found version ", zap.Any("version", version), zap.String("content", s))
	return s, nil
}

func (g *GoProxy) GetVersionInfo(module string, version string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/%s.info", g.Url, module, version))
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", version), zap.Error(err)).Error("Could not fetch module version")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("no module version found for %s:%s; request status: %d;", module, version, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", version), zap.Error(err)).Error("Could not read proxy response body")
	}

	s := string(body)
	logger.With(zap.String("module", module)).Debug("Found version ", zap.Any("version", version), zap.String("content", s))
	return s, nil
}

func (g *GoProxy) GetZipFile(module string, version string) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/@v/%s.zip", g.Url, module, version))
	if err != nil {
		logger.With(zap.String("module", module), zap.String("version", version), zap.Error(err)).Error("Could not fetch module version")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Zip not found for %s:%s; request status: %d;", module, version, resp.StatusCode)
	}

	return resp.Body, nil
}

func (g *GoProxy) Match(url string) bool {
	return true
}
