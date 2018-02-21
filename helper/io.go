package helper

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func WriteFileWithDirectory(path string, data []byte, perm os.FileMode) error {

	s := strings.Split(path, "/")

	var dir string
	if len(s) > 1 {
		dir = strings.Join(s[0:len(s)-1], "/")
	} else {
		dir = path
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return errors.Wrapf(err, "create directory is failed. [%s]", dir)
	}

	if err := ioutil.WriteFile(path, data, perm); err != nil {
		return errors.Wrapf(err, "write data to file is failed. [%s]", path)
	}

	return nil
}
