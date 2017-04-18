package dependency

import (
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Dependency interface {
	Load() (*ProtoDep, error)
}

type DependencyImpl struct {
	tomlpath string
}

func NewDependency(targetDir string) Dependency {
	tomlpath := filepath.Join(targetDir, "protodep.toml")
	return &DependencyImpl{
		tomlpath: tomlpath,
	}
}

func (d *DependencyImpl) Load() (*ProtoDep, error) {

	content, err := ioutil.ReadFile(d.tomlpath)
	if err != nil {
		return nil, errors.Wrapf(err, "load %s is failed", d.tomlpath)
	}

	var conf ProtoDep
	if _, err := toml.Decode(string(content), &conf); err != nil {
		return nil, errors.Wrap(err, "decode toml is failed")
	}
	return &conf, nil
}
