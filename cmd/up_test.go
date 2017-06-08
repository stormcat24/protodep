package cmd

import (
	"testing"
	"os"
	"github.com/stretchr/testify/require"
)

func TestUp(t *testing.T) {

	as := []string{"up"}
	osargs := []string{"cmd"}
	os.Args = append(osargs, as...)

	unitTest = true
	err := upCmd.Execute()
	require.NoError(t, err)
}