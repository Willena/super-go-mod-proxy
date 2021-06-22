package fetchMethods

import (
	"github.com/stretchr/testify/assert"
	"github.com/willena/super-go-mod-proxy/types"
	"testing"
)

func TestFindMethodGit(t *testing.T) {
	method, err := FindForUrl("git+https://github.com/www/toto", types.AuthConfiguration{})
	assert.Nil(t, err)
	assert.IsType(t, &Git{}, method, "Types is Git HTTP")
	method, err = FindForUrl("git+https://github.com/www/toto", types.AuthConfiguration{})
	assert.Nil(t, err)
	assert.IsType(t, &Git{}, method, "Types is Git HTTP")
	method, err = FindForUrl("git+ssh://github.com/www/toto", types.AuthConfiguration{})
	assert.Nil(t, err)
	assert.IsType(t, &Git{}, method, "Types is Git SSH")
}

func TestFindMethodHTTP(t *testing.T) {
	method, err := FindForUrl("https://proxy.golang.org/", types.AuthConfiguration{})
	assert.Nil(t, err)
	assert.IsType(t, &GoProxy{}, method, "Types is Go Proxy")
}

func TestFindMethodNoMethod(t *testing.T) {
	_, err := FindForUrl("sjhfodf", types.AuthConfiguration{})
	assert.NotNil(t, err)
}

func TestSortTagsTest(t *testing.T) {
	unsortedTags := []string{
		"v1.2.5",
		"v0.0.1",
		"v10.5.3",
		"v1.1.1",
		"x.z.e",
	}

	sortedtags := []string{
		"v0.0.1",
		"v1.1.1",
		"v1.2.5",
		"v10.5.3",
	}

	resultSort := sortTags(unsortedTags)
	assert.Equal(t, sortedtags, resultSort, "Tags are sorted.")
}
