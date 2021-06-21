package gomodule

import (
	"testing"
)

func TestParseVersions(t *testing.T) {
	ParseFullModuleVersion("v1.1.1")
	ParseFullModuleVersion("v1.1.1-yyyymmddhhmmss-abcdefabcdef")
	ParseFullModuleVersion("v1.1.2-patch+mmm45")
	ParseFullModuleVersion("v1.2.3-pre.0.yyyymmddhhmmss-abcdefabcdef")
	ParseFullModuleVersion("v1.2.4-0.yyyymmddhhmmss-abcdefabcdef")
	ParseFullModuleVersion("v1.2.4-0.yyyymmddhhmmss-abcdefabcdef+meta")
	ParseFullModuleVersion("master")
	ParseFullModuleVersion("branch")
}
