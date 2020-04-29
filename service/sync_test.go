package service

import (
	"fmt"
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
	authProviderMock.EXPECT().AuthMethod().Return(nil, nil).AnyTimes()
	authProviderMock.EXPECT().GetRepositoryURL("github.com/protocolbuffers/protobuf").Return("https://github.com/protocolbuffers/protobuf.git")
	authProviderMock.EXPECT().GetRepositoryURL("github.com/opensaasstudio/plasma").Return("https://github.com/opensaasstudio/plasma.git")

	pwd, err := os.Getwd()
	fmt.Println(pwd)

	require.NoError(t, err)

	outputRootDir := os.TempDir()

	target := NewSync(authProviderMock, dotProtoDir, pwd, outputRootDir)
	// clone
	err = target.Resolve(false, false)
	require.NoError(t, err)

	if !isFileExist(filepath.Join(outputRootDir, "proto/stream.proto")) {
		t.Error("not found file [proto/stream.proto]")
	}
	if !isFileExist(filepath.Join(outputRootDir, "proto/google/protobuf/empty.proto")) {
		t.Error("not found file [proto/google/protobuf/empty.proto]")
	}

	// check ignore worked
	// hasPrefix test - backward compatibility
	if isFileExist(filepath.Join(outputRootDir, "proto/google/protobuf/test_messages_proto3.proto")) {
		t.Error("found file [proto/google/protobuf/test_messages_proto3.proto]")
	}

	// glob test 1
	if isFileExist(filepath.Join(outputRootDir, "proto/google/protobuf/test_messages_proto2.proto")) {
		t.Error("found file [proto/google/protobuf/test_messages_proto2.proto]")
	}

	// glob test 2
	if isFileExist(filepath.Join(outputRootDir, "proto/google/protobuf/test_messages_proto2.proto")) {
		t.Error("found file [proto/google/protobuf/test_messages_proto2.proto]")
	}

	// glob test 3
	if isFileExist(filepath.Join(outputRootDir, "proto/google/protobuf/util/internal/testdata/")) {
		t.Error("found file [proto/google/protobuf/util/internal/testdata/]")
	}

	// fetch
	err = target.Resolve(false, false)
	require.NoError(t, err)
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
