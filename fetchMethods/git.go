package fetchMethods

import (
	"encoding/json"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/willena/super-go-mod-proxy/modulezip"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	ssh2 "golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var regexpGitHttps = regexp.MustCompile("git\\+https?://.*")
var regexpGitSSH = regexp.MustCompile("git\\+ssh://.*")

type Git struct {
	Url  string
	Auth types.AuthConfiguration
}

func (g *Git) GetVersions(module string) ([]string, error) {
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		logger.Error("GIT ERROR", zap.Error(err))
		return nil, err
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{Name: "main", URLs: []string{g.Url}})
	if err != nil {
		logger.Error("GIT ERROR", zap.Error(err))
		return nil, err
	}

	data, err := remote.List(&git.ListOptions{})
	if err != nil {
		logger.Error("GIT ERROR", zap.Error(err))
		return nil, err
	}

	tags := make([]string, 0)
	for _, d := range data {
		if d.Name().IsTag() {
			if ref := d.Strings()[0][10:]; strings.HasPrefix(ref, "v") {
				tags = append(tags, ref)
			}
		}
	}
	logger.With(zap.Any("References", tags)).Info("Found References when calling version ! ")

	return tags, nil

}

func (g *Git) GetLatestVersion(module string) (string, error) {
	panic("implement me")
}

func (g *Git) getAuth() (transport.AuthMethod, error) {
	switch g.Auth.Type {
	case "privateKey":
		_, err := os.Stat(g.Auth.PrivateKey)
		if err != nil {
			logger.Error("Could not read provided privateKey file", zap.String("privateKey", g.Auth.PrivateKey), zap.Error(err))
			return nil, err
		}

		publicKeys, err := ssh.NewPublicKeysFromFile(g.Auth.Username, g.Auth.PrivateKey, g.Auth.Password)
		if err != nil {
			logger.Error("Could not derivate publicKey from private key", zap.Error(err))
			return nil, err
		}

		//Do a copy to avoid having to re-code all the logic behind NewPublicKeysFromFile
		//Warning accepting all hosts keys !
		newPublicKey := &ssh.PublicKeys{User: publicKeys.User, Signer: publicKeys.Signer, HostKeyCallbackHelper: ssh.HostKeyCallbackHelper{
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh2.PublicKey) error { return nil },
		}}

		return newPublicKey, nil
	case "basic":
		return &http.BasicAuth{
			Username: g.Auth.Username, // yes, this can be anything except an empty string
			Password: g.Auth.Password}, nil
	}

	if g.Auth.Type != "" {
		logger.Warn("Unknown provided auth method for git ", zap.String("authMethod", g.Auth.Type))
	}
	return nil, nil
}

func (g *Git) GetModule(module string, version string) (string, error) {
	modulePath := path.Join(os.TempDir(), module, version)
	os.MkdirAll(modulePath, 0777)

	//Todo: avoid doing that !
	os.RemoveAll(modulePath)

	auth, err := g.getAuth()
	if err != nil {
		logger.Error("Error while preparing authentication", zap.Error(err))
		return "", err
	}

	r, err := git.PlainClone(modulePath, false, &git.CloneOptions{
		URL:           g.Url,
		Auth:          auth,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName("refs/tags/" + version),
	})

	if err != nil {
		logger.Error("Ref not found", zap.Error(err))
		return "", err
	}

	w, err := r.Worktree()
	if err != nil {
		logger.Error("GITERR", zap.Error(err))
		return "", err
	}

	f, err := w.Filesystem.Open("go.mod")
	if err != nil {
		logger.Error("Coudl not open go.mod file")
		return "", err
	}
	d, err := io.ReadAll(f)

	if err != nil {
		logger.Error("Error while reading go.mod file !")
	}

	s := string(d)

	return s, err

}

func (g *Git) GetVersionInfo(module string, version string) (string, error) {
	type Info struct {
		Version string    // version string
		Time    time.Time // commit time
	}

	_, err := g.GetModule(module, version)
	if err != nil {
		logger.Error("Could get module", zap.String("module", module), zap.String("version", version))
		return "", err
	}

	pathGit := path.Join(os.TempDir(), module, version)
	r, err := git.PlainOpen(pathGit)
	if err != nil {
		logger.Error("Could not open repo for module", zap.String("module", module), zap.String("version", version))
		return "", err
	}

	ref, err := r.Head()
	if err != nil {
		logger.Error("Could not get head for module", zap.String("module", module), zap.String("version", version))
		return "", err
	}
	obj, err := r.CommitObject(ref.Hash())
	if err != nil {
		logger.Error("Could not get commit object for module", zap.String("module", module), zap.String("version", version))
		return "", err
	}
	info := Info{Version: version, Time: obj.Author.When}
	data, err := json.Marshal(info)

	return string(data), err
}

func (g *Git) GetZipFile(module string, version string) (io.Reader, error) {

	_, err := g.GetModule(module, version)

	buff, err := modulezip.ZipModule(module, version)

	return buff, err
}

func (g *Git) Match(url string) bool {
	return regexpGitHttps.MatchString(url) || regexpGitSSH.MatchString(url)
}
