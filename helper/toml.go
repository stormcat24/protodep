package helper

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

func WriteToml(dest string, input interface{}) error {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(input); err != nil {
		return errors.Wrapf(err, "encode config to toml format is failed. [%+v]", input)
	}

	if err := ioutil.WriteFile(dest, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "write to %s is failed", dest)
	}

	return nil
}
