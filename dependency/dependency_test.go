package dependency

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {

	pwd, _ := os.Getwd()
	target := NewDependency(pwd)

	actual, err := target.Load()
	require.NoError(t, err)
	require.Equal(t, 2, len(actual.Dependencies))

	withBranch := actual.Dependencies[0]
	withRevision := actual.Dependencies[1]

	require.Equal(t, "github.com/google/protobuf/examples", withBranch.Name)
	require.Equal(t, "master", withBranch.Branch)
	require.Equal(t, "", withBranch.Revision)

	require.Equal(t, "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb", withRevision.Name)
	require.Equal(t, "", withRevision.Branch)
	require.Equal(t, "v1.2.2", withRevision.Revision)
}
