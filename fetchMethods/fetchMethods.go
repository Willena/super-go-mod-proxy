package fetchMethods

import (
	"fmt"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
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
		return &GoProxy{"https://proxy.golang.org"}, nil
	}

	return nil, fmt.Errorf("Method not found !")

}
