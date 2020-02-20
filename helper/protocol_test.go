package helper

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIsAvailableSSH(t *testing.T) {
	f, err := ioutil.TempFile("", "id_rsa")
	require.NoError(t, err)

	found, err := IsAvailableSSH(f.Name())
	require.NoError(t, err)
	require.True(t, found)

	notFound, err := IsAvailableSSH(fmt.Sprintf("/tmp/IsAvailableSSH_%d", time.Now().UnixNano()))
	require.NoError(t, err)
	require.False(t, notFound)
}
