package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mitchellh/go-homedir"
	"github.com/stormcat24/protodep/helper"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {

	homeDir, err := homedir.Dir()
	require.NoError(t, err)

	dotProtoDir := filepath.Join(homeDir, "protodep_ut")
	err = os.RemoveAll(dotProtoDir)
	require.NoError(t, err)

	c := gomock.NewController(t)
	defer c.Finish()

	authProviderMock := helper.NewMockAuthProvider(c)
	authProviderMock.EXPECT().AuthMethod().Return(nil).AnyTimes()
	authProviderMock.EXPECT().GetRepositoryURL("github.com/google/protobuf").Return("https://github.com/google/protobuf.git")
	authProviderMock.EXPECT().GetRepositoryURL("github.com/openfresh/plasma").Return("https://github.com/openfresh/plasma.git")
	authProviderMock.EXPECT().GetRepositoryURL("github.com/kubernetes/helm").Return("https://github.com/kubernetes/helm.git")

	pwd, err := os.Getwd()
	require.NoError(t, err)

	outputRootDir := os.TempDir()

	target := NewSync(authProviderMock, dotProtoDir, pwd, outputRootDir)
	// clone
	err = target.Resolve(false)
	require.NoError(t, err)

	if !isFileExist(filepath.Join(outputRootDir, "proto/stream.proto")) {
		t.Error("not found file [proto/stream.proto]")
	}
	if !isFileExist(filepath.Join(outputRootDir, "proto/google/protobuf/empty.proto")) {
		t.Error("not found file [proto/google/protobuf/empty.proto]")
	}
	if !isFileExist(filepath.Join(outputRootDir, "proto/chart/chart.proto")) {
		t.Error("not found file [proto/chart/chart.proto]")
	}

	// fetch
	err = target.Resolve(false)
	require.NoError(t, err)
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
