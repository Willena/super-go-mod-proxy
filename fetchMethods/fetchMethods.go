package fetchMethods

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"sort"
	"strings"
)

var logger, _ = zap.NewDevelopment()

func FindForUrl(url string, authConfig types.AuthConfiguration) (types.FetchMethod, error) {
	const gitPrefix = "git+"

	if strings.HasPrefix(url, gitPrefix) {
		//remove git+
		newUrl := url[len(gitPrefix):]
		return &Git{Url: newUrl, Auth: authConfig}, nil
	}

	if strings.HasPrefix(url, "http") {
		return &GoProxy{url}, nil
	}

	return nil, fmt.Errorf("Method not found !")
}

func sortTags(tags []string) []string {
	vs := make([]*semver.Version, 0)
	for _, r := range tags {
		v, err := semver.NewVersion(r)
		if err != nil {
			logger.Warn("Error parsing version", zap.String("error", err.Error()), zap.String("version", r))
			continue
		}
		vs = append(vs, v)
	}
	sort.Sort(semver.Collection(vs))
	sortedTags := make([]string, len(vs))
	for i, v := range vs {
		sortedTags[i] = v.Original()
	}

	return sortedTags
}
