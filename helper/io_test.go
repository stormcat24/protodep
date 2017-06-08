package helper

import (
	"testing"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/require"
	"io/ioutil"
)

func TestWriteFileWithDirectory(t *testing.T) {

	destDir := os.TempDir()
	testDir := filepath.Join(destDir, "hoge")
	testFile := filepath.Join(testDir, "fuga.txt")

	err := WriteFileWithDirectory(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	stat, err := os.Stat(testFile)
	require.NoError(t, err)
	require.True(t, !stat.IsDir())

	data, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, string(data), "test")

}