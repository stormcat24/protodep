package cmd

import (
	"testing"
	"os"
	"github.com/stretchr/testify/require"
)

func TestDep(t *testing.T) {

	as := []string{"dep"}
	osargs := []string{"cmd"}
	os.Args = append(osargs, as...)

	err := depCmd.Execute()
	require.NoError(t, err)
}