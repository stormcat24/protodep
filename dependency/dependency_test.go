package dependency

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {

	pwd, _ := os.Getwd()
	target := NewDependency(pwd, false)

	actual, err := target.Load()
	require.NoError(t, err)
	require.Equal(t, 3, len(actual.Dependencies))

	withBranch := actual.Dependencies[0]
	withRevision := actual.Dependencies[1]
	withIgnore := actual.Dependencies[2]

	require.Equal(t, "github.com/protocolbuffers/protobuf/src", withBranch.Target)
	require.Equal(t, "master", withBranch.Branch)
	require.Equal(t, "", withBranch.Revision)
	require.Equal(t, "", withBranch.Protocol)

	require.Equal(t, "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb", withRevision.Target)
	require.Equal(t, "", withRevision.Branch)
	require.Equal(t, "v1.2.2", withRevision.Revision)
	require.Equal(t, "grpc-gateway/examplepb", withRevision.Path)
	require.Equal(t, "ssh", withRevision.Protocol)

	require.Equal(t, "github.com/kubernetes/helm/_proto/hapi", withIgnore.Target)
	require.Equal(t, "", withIgnore.Branch)
	require.Equal(t, "v2.8.1", withIgnore.Revision)
	require.Equal(t, []string{"./release", "./rudder", "./services", "./version"}, withIgnore.Ignores)
	require.Equal(t, "https", withIgnore.Protocol)
}
