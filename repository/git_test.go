package repository

import (
	"testing"
	"os/user"
	"github.com/stormcat24/protodep/dependency"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {

	user, _ := user.Current()
	target := NewGitRepository(user.HomeDir, dependency.ProtoDepDependency{
		Name: "github.com/openfresh/plasma/protobuf",
		Branch: "master",
	})

	_, err := target.Open()
	require.NoError(t, err)
}