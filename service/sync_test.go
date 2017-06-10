package service

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stormcat24/protodep/helper"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
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
	authProviderMock.EXPECT().GetRepositoryURL(gomock.Any()).Return("https://github.com/openfresh/plasma.git").AnyTimes()

	pwd, err := os.Getwd()
	require.NoError(t, err)

	outputRootDir := os.TempDir()

	target := NewSync(authProviderMock, dotProtoDir, pwd, outputRootDir)
	// clone
	err = target.Resolve(false)
	require.NoError(t, err)


	// fetch
	err = target.Resolve(false)
	require.NoError(t, err)
}