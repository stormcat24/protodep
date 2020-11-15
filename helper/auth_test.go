package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRepositoryURLWithSSH(t *testing.T) {
	target := &AuthProviderWithSSH{}
	actual := target.GetRepositoryURL("github.com/stormcat24/protodep")

	require.Equal(t, "ssh://github.com/stormcat24/protodep.git", actual)
}

func TestGetRepositoryURLWithSSHAgent(t *testing.T) {
	target := &AuthProviderWithSSHAgent{}
	actual := target.GetRepositoryURL("github.com/stormcat24/protodep")

	require.Equal(t, "ssh://github.com/stormcat24/protodep.git", actual)
}

func TestGetRepositoryURLHTTPS(t *testing.T) {
	target := &AuthProviderHTTPS{}
	actual := target.GetRepositoryURL("github.com/stormcat24/protodep")

	require.Equal(t, "https://github.com/stormcat24/protodep.git", actual)
}
