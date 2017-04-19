package dependency

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {

	matched := ProtoDepDependency{
		Name: "github.com/google/protobuf",
	}
	require.Equal(t, "github.com/google/protobuf", matched.Repository())

	protruded := ProtoDepDependency{
		Name: "github.com/google/protobuf/examples",
	}
	require.Equal(t, "github.com/google/protobuf", protruded.Repository())
}

func TestDirectory(t *testing.T) {

	matched := ProtoDepDependency{
		Name: "github.com/google/protobuf",
	}
	require.Equal(t, ".", matched.Directory())

	protruded := ProtoDepDependency{
		Name: "github.com/google/protobuf/examples",
	}

	require.Equal(t, "./examples", protruded.Directory())
}