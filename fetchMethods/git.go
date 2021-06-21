package fetchMethods

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/willena/super-go-mod-proxy/gomodule"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	ssh2 "golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var regexpGitHttps = regexp.MustCompile("git\\+https?://.*")
var regexpGitSSH = regexp.MustCompile("git\\+ssh://.*")

type Git struct {
	Url  string
	Auth types.AuthConfiguration
}

type Info struct {
	Version string    // version string
	Time    time.Time // commit time
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

	auth, err := g.getAuth()
	if err != nil {
		return nil, err
	}

	data, err := remote.List(&git.ListOptions{Auth: auth})
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

	vs := make([]*semver.Version, len(tags))
	for i, r := range tags {
		v, err := semver.NewVersion(r)
		if err != nil {
			logger.Warn("Error parsing version", zap.String("error", err.Error()), zap.String("version", r))
		}
		vs[i] = v
	}
	sort.Sort(semver.Collection(vs))

	for i, v := range vs {
		tags[i] = v.Original()
	}

	return tags, nil

}

func (g *Git) GetLatestVersion(module string) (string, error) {
	versions, err := g.GetVersions(module)
	if err != nil {
		return "", err
	}

	if len(versions) > 0 {
		return g.GetVersionInfo(module, versions[len(versions)-1])
	} else {
		r, _, err := g.downloadRepo(module, "")
		if err != nil {
			logger.Error("Could not download repository", zap.Error(err))
		}

		ref, err := r.Head()
		if err != nil {
			logger.Error("Could not get head for module", zap.String("module", module), zap.String("version", "latest"))
			return "", err
		}

		return g.GetVersionInfo(module, ref.Hash().String())
	}

}

func (g *Git) downloadRepo(module, version string) (*git.Repository, *git.Worktree, error) {
	fullVersion := gomodule.ParseFullModuleVersion(version)

	auth, err := g.getAuth()
	if err != nil {
		logger.Error("Error while preparing authentication", zap.Error(err))
		return nil, nil, err
	}

	var r *git.Repository

	r, err = git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:  g.Url,
		Auth: auth,
	})

	if err != nil {
		logger.Error("Could not clone repo", zap.Error(err))
		return nil, nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		logger.Error("GITERR", zap.Error(err))
		return nil, nil, err
	}

	if fullVersion.CommitRef != nil {
		resolved, err := r.ResolveRevision(plumbing.Revision(fullVersion.CommitRef.Ref))
		if err != nil {
			logger.Error("could not resolve revision", zap.String("Resv", fullVersion.CommitRef.Ref))
			return nil, nil, err
		}
		logger.Debug("Collecting checkout to commit", zap.String("commit", resolved.String()))
		err = w.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(resolved.String())})
		if err != nil {
			logger.Error("Checkout faild !", zap.String("Hash", fullVersion.CommitRef.Ref))
			return nil, nil, err
		}
	}

	return r, w, nil
}

func findNearestTagInLog(r *git.Repository, hash plumbing.Hash) (string, error) {
	log, err := r.Log(&git.LogOptions{
		From: hash,
	})

	if err != nil {
		logger.Error("Could not get git log...")
		return "", err
	}

	tags, err := r.Tags()
	if err != nil {
		logger.Error("Could not get git tags...")
		return "", err
	}

	tagMap := make(map[string]*plumbing.Reference)

	err = tags.ForEach(func(reference *plumbing.Reference) error {
		if strings.HasPrefix(reference.Name().Short(), "v") {
			tagMap[reference.Hash().String()] = reference
		}
		return nil
	})

	if err != nil {
		logger.Error("Error while processsing tags...")
		return "", err
	}

	if len(tagMap) == 0 {
		return "", nil
	}

	var tagName = ""
	log.ForEach(func(commit *object.Commit) error {
		if val, ok := tagMap[commit.Hash.String()]; ok {
			tagName = val.Name().Short()
			return nil
		}
		return nil
	})

	return tagName, nil
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

	_, w, err := g.downloadRepo(module, version)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	if _, err := w.Filesystem.Stat("go.mod"); os.IsNotExist(err) {
		logger.Warn("No go.mod file found, return default go.mod")
		return fmt.Sprintf("module %s", module), nil
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

	r, _, err := g.downloadRepo(module, version)
	if err != nil {
		logger.Error(err.Error())
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

	versionFull := gomodule.ParseFullModuleVersion(version)
	var info Info
	var tagFound string
	if versionFull.Parsed == nil {
		//No major version given !
		tagFound, err = findNearestTagInLog(r, ref.Hash())
		if err != nil || tagFound == "" {
			logger.Error("Could not found nearest tag. ")
			tagFound = "v0.0.0"
			info = Info{Version: fmt.Sprintf("%s-%s-%s", tagFound, obj.Committer.When.UTC().Format("20060102150405"), ref.Hash().String()[:12]), Time: obj.Committer.When.UTC()}
		} else {
			semVersion, err := semver.NewVersion(tagFound)
			if err != nil {
				logger.Error("Tag not valid !")
				return "", err
			}
			nversion := semVersion.IncPatch()
			nversion, _ = nversion.SetPrerelease("0")
			info = Info{Version: fmt.Sprintf("%s.%s-%s", "v"+nversion.String(), obj.Committer.When.UTC().Format("20060102150405"), ref.Hash().String()[:12]), Time: obj.Committer.When.UTC()}
		}
	} else {
		info = Info{Version: version, Time: obj.Committer.When.UTC()}
	}
	data, err := json.Marshal(info)

	return string(data), err
}

func (g *Git) GetZipFile(module string, version string) (io.Reader, error) {

	_, w, err := g.downloadRepo(module, version)
	if err != nil {
		return nil, err
	}

	buff, err := gomodule.ZipModule(w.Filesystem, module, version)

	return buff, err
}

func (g *Git) Match(url string) bool {
	return regexpGitHttps.MatchString(url) || regexpGitSSH.MatchString(url)
}
