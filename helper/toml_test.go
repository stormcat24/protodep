package helper

import (
	"testing"
	"github.com/stormcat24/protodep/dependency"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/require"
)

func TestWriteToml(t *testing.T) {

	config := dependency.ProtoDep{
		ProtoOutdir: "./proto",
		Dependencies: []dependency.ProtoDepDependency{
			dependency.ProtoDepDependency{
				Target:   "github.com/openfresh/plasma/protobuf",
				Branch:   "master",
				Revision: "d7ee1d95b6700756b293b722a1cfd4b905a351ba",
			},
			dependency.ProtoDepDependency{
				Target:   "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb",
				Branch:   "master",
				Revision: "c6f7a5ac629444a556bb665e389e41b897ebad39",
			},
		},
	}

	destDir := os.TempDir()
	destFile := filepath.Join(destDir, "protodep.lock")

	require.NoError(t, os.MkdirAll(os.TempDir(), 0777))
	require.NoError(t, WriteToml(destFile, config))

	stat, err := os.Stat(destFile)
	require.NoError(t, err)

	require.True(t, !stat.IsDir())
}
