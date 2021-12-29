package config

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
	require.Equal(t, 2, len(actual.Dependencies))

	withBranch := actual.Dependencies[0]
	withRevision := actual.Dependencies[1]

	require.Equal(t, "github.com/protocolbuffers/protobuf/src", withBranch.Target)
	require.Equal(t, "master", withBranch.Branch)
	require.Equal(t, "", withBranch.Revision)
	require.Equal(t, "", withBranch.Protocol)

	require.Equal(t, "github.com/grpc-ecosystem/grpc-gateway/examples/internal/helloworld", withRevision.Target)
	require.Equal(t, "", withRevision.Branch)
	require.Equal(t, "v2.7.2", withRevision.Revision)
	require.Equal(t, "grpc-gateway/examples/internal/helloworld", withRevision.Path)
	require.Equal(t, "ssh", withRevision.Protocol)
}
