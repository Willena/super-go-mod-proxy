package gomodule

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"strings"
	"time"
)

const timeFormat = "20060102150405"

var regexIsHash = regexp.MustCompile("^[0-9a-f]{5,40}$")

func ShortHash(ref *plumbing.Reference) string {
	return ref.Hash().String()[:12]
}

func FormatAsValidVersionVersion(version string, commit *object.Commit, ref *plumbing.Reference, pre bool) string {

	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if commit == nil || ref == nil {
		return version
	}

	if pre {
		return fmt.Sprintf("%s.%s-%s", version, commit.Committer.When.UTC().Format(timeFormat), ShortHash(ref))
	} else {
		return fmt.Sprintf("%s-%s-%s", version, commit.Committer.When.UTC().Format(timeFormat), ShortHash(ref))
	}

}

type CommitRef struct {
	Ref  string
	time time.Time
}

type FullVersion struct {
	Raw       string
	Parsed    *semver.Version
	CommitRef *CommitRef
}

func (v *FullVersion) String() string {
	return v.Raw
}

type GoModule struct {
	Path    string
	Version *FullVersion
}

func (m *GoModule) SetVersion(s string) {
	if s == "" {
		m.Version = nil
		return
	}

	m.Version = ParseFullModuleVersion(s)
}

func (m *GoModule) MinimalGoModFile() string {
	return fmt.Sprintf("module %s", m.Path)
}

func (c *CommitRef) IsHash() bool {
	return regexIsHash.MatchString(c.Ref)
}

func parsePrerelease(pre string) *CommitRef {

	if pre == "" {
		return nil
	}

	timeAndHash := pre
	lastIndex := strings.LastIndex(pre, ".")
	if lastIndex != -1 {
		timeAndHash = pre[lastIndex:]
	}

	splitedTImeAndHash := strings.Split(timeAndHash, "-")
	if len(splitedTImeAndHash) == 2 {
		t, _ := time.Parse("20060102150405", splitedTImeAndHash[0])
		return &CommitRef{
			Ref:  splitedTImeAndHash[1],
			time: t,
		}
	}

	return nil

}

func ParseFullModuleVersion(version string) *FullVersion {
	v, err := semver.NewVersion(version)

	if err != nil {
		fmt.Println("Could not parse semversion. Assuming a branch ref")
		return &FullVersion{Raw: version, Parsed: nil, CommitRef: &CommitRef{Ref: version}}
	}

	return &FullVersion{
		Raw:       version,
		Parsed:    v,
		CommitRef: parsePrerelease(v.Prerelease()),
	}
}

func NewGoModule(name, version string) *GoModule {
	return &GoModule{
		Path:    name,
		Version: ParseFullModuleVersion(version),
	}
}
