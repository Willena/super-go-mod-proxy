package gomodule

import (
	"fmt"
	"github.com/Masterminds/semver"
	"regexp"
	"strings"
	"time"
)

var regexIsHash = regexp.MustCompile("^[0-9a-f]{5,40}$")

type CommitRef struct {
	Ref  string
	time time.Time
}

type FullVersion struct {
	Raw       string
	Parsed    *semver.Version
	CommitRef *CommitRef
}

type GoModule struct {
	Path    string
	Version FullVersion
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
