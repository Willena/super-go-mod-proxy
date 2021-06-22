package gomodule

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseSimpleVersion(t *testing.T) {
	v := ParseFullModuleVersion("v1.1.1")

	assert.NotNil(t, v)
	assert.Equal(t, v.Raw, "v1.1.1")
	assert.Nil(t, v.CommitRef)
	assert.Equal(t, v.Parsed.Major(), int64(1))
	assert.Equal(t, v.Parsed.Minor(), int64(1))
	assert.Equal(t, v.Parsed.Patch(), int64(1))
}

func TestVersionWithDateAndCommit(t *testing.T) {
	v := ParseFullModuleVersion("v1.1.1-20210512055052-abcdefabcdef")
	assert.NotNil(t, v)
	assert.Equal(t, v.Raw, "v1.1.1-20210512055052-abcdefabcdef")
	assert.Equal(t, v.Parsed.Major(), int64(1))
	assert.Equal(t, v.Parsed.Minor(), int64(1))
	assert.Equal(t, v.Parsed.Patch(), int64(1))
	assert.Equal(t, v.CommitRef.Ref, "abcdefabcdef")
	dateEqual(t, v.CommitRef.time)
}

func TestWithoutCommitRefButPreversion(t *testing.T) {
	v := ParseFullModuleVersion("v1.1.2-patch+mmm45")
	assert.NotNil(t, v)
	assert.Equal(t, v.Raw, "v1.1.2-patch+mmm45")
	assert.Equal(t, v.Parsed.Major(), int64(1))
	assert.Equal(t, v.Parsed.Minor(), int64(1))
	assert.Equal(t, v.Parsed.Patch(), int64(2))
	assert.Nil(t, v.CommitRef)
}

func TestVersionPreWithCommitAndTime(t *testing.T) {
	v := ParseFullModuleVersion("v1.2.3-pre.0.20210512055052-abcdefabcdef")
	assert.NotNil(t, v)
	assert.Equal(t, v.Raw, "v1.2.3-pre.0.20210512055052-abcdefabcdef")
	assert.Equal(t, v.Parsed.Major(), int64(1))
	assert.Equal(t, v.Parsed.Minor(), int64(2))
	assert.Equal(t, v.Parsed.Patch(), int64(3))
	assert.Equal(t, v.CommitRef.Ref, "abcdefabcdef")
	dateEqual(t, v.CommitRef.time)
}

func dateEqual(t *testing.T, d time.Time) {
	assert.Equal(t, d.Year(), 2021)
	assert.Equal(t, d.Month(), time.Month(05))
	assert.Equal(t, d.Day(), 12)
	assert.Equal(t, d.Hour(), 05)
	assert.Equal(t, d.Minute(), 50)
	assert.Equal(t, d.Second(), 52)
}

func TestVersionCloseToTag(t *testing.T) {
	v := ParseFullModuleVersion("v1.2.4-0.20210512055052-abcdefabcdef")
	assert.NotNil(t, v)
	assert.Equal(t, v.Raw, "v1.2.4-0.20210512055052-abcdefabcdef")
	assert.Equal(t, v.Parsed.Major(), int64(1))
	assert.Equal(t, v.Parsed.Minor(), int64(2))
	assert.Equal(t, v.Parsed.Patch(), int64(4))
	assert.Equal(t, v.CommitRef.Ref, "abcdefabcdef")
	dateEqual(t, v.CommitRef.time)
}

func TestVersionOnlyRef(t *testing.T) {
	v := ParseFullModuleVersion("master")
	assert.NotNil(t, v)
	assert.Equal(t, v.Raw, "master")
	assert.Nil(t, v.Parsed)
	assert.Equal(t, v.CommitRef.Ref, "master")
}
