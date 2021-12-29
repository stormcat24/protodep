package resolver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	"github.com/stormcat24/protodep/pkg/auth"
	"github.com/stormcat24/protodep/pkg/config"
)

func TestSync(t *testing.T) {

	homeDir, err := homedir.Dir()
	require.NoError(t, err)

	dotProtoDir := filepath.Join(homeDir, "protodep_ut")
	err = os.RemoveAll(dotProtoDir)
	require.NoError(t, err)

	pwd, err := os.Getwd()
	fmt.Println(pwd)

	require.NoError(t, err)

	outputRootDir := os.TempDir()

	conf := Config{
		HomeDir:   dotProtoDir,
		TargetDir: pwd,
		OutputDir: outputRootDir,
	}

	target, err := New(&conf)
	require.NoError(t, err)

	c := gomock.NewController(t)
	defer c.Finish()

	httpsAuthProviderMock := auth.NewMockAuthProvider(c)
	httpsAuthProviderMock.EXPECT().AuthMethod().Return(nil, nil).AnyTimes()
	httpsAuthProviderMock.EXPECT().GetRepositoryURL("github.com/protocolbuffers/protobuf").Return("https://github.com/protocolbuffers/protobuf.git")

	sshAuthProviderMock := auth.NewMockAuthProvider(c)
	sshAuthProviderMock.EXPECT().AuthMethod().Return(nil, nil).AnyTimes()
	sshAuthProviderMock.EXPECT().GetRepositoryURL("github.com/opensaasstudio/plasma").Return("https://github.com/opensaasstudio/plasma.git")

	target.SetHttpsAuthProvider(httpsAuthProviderMock)
	target.SetSshAuthProvider(sshAuthProviderMock)

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

func TestWriteToml(t *testing.T) {

	config := config.ProtoDep{
		ProtoOutdir: "./proto",
		Dependencies: []config.ProtoDepDependency{
			config.ProtoDepDependency{
				Target:   "github.com/openfresh/plasma/protobuf",
				Branch:   "master",
				Revision: "d7ee1d95b6700756b293b722a1cfd4b905a351ba",
			},
			config.ProtoDepDependency{
				Target:   "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb",
				Branch:   "master",
				Revision: "c6f7a5ac629444a556bb665e389e41b897ebad39",
			},
		},
	}

	destDir := os.TempDir()
	destFile := filepath.Join(destDir, "protodep.lock")

	require.NoError(t, os.MkdirAll(os.TempDir(), 0777))
	require.NoError(t, writeToml(destFile, config))

	stat, err := os.Stat(destFile)
	require.NoError(t, err)

	require.True(t, !stat.IsDir())
}

func TestWriteFileWithDirectory(t *testing.T) {
	destDir := os.TempDir()
	testDir := filepath.Join(destDir, "hoge")
	testFile := filepath.Join(testDir, "fuga.txt")

	err := writeFileWithDirectory(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	stat, err := os.Stat(testFile)
	require.NoError(t, err)
	require.True(t, !stat.IsDir())

	data, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, string(data), "test")
}

func TestIsAvailableSSH(t *testing.T) {
	f, err := ioutil.TempFile("", "id_rsa")
	require.NoError(t, err)

	found, err := isAvailableSSH(f.Name())
	require.NoError(t, err)
	require.True(t, found)

	notFound, err := isAvailableSSH(fmt.Sprintf("/tmp/IsAvailableSSH_%d", time.Now().UnixNano()))
	require.NoError(t, err)
	require.False(t, notFound)
}

