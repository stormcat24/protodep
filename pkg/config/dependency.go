package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Dependency interface {
	Load() (*ProtoDep, error)
	IsNeedWriteLockFile() bool
}

type DependencyImpl struct {
	targetDir string
	tomlPath    string
	lockPath    string
	forceUpdate bool
}

func NewDependency(targetDir string, forceUpdate bool) Dependency {
	return &DependencyImpl{
		targetDir:   targetDir,
		tomlPath:    filepath.Join(targetDir, "protodep.toml"),
		lockPath:    filepath.Join(targetDir, "protodep.lock"),
		forceUpdate: forceUpdate,
	}
}

func (d *DependencyImpl) Load() (*ProtoDep, error) {

	var targetConfig string
	if d.IsNeedWriteLockFile() {
		targetConfig = d.tomlPath
	} else {
		targetConfig = d.lockPath
	}

	content, err := os.ReadFile(targetConfig)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", targetConfig, err)
	}

	var conf ProtoDep
	if _, err := toml.Decode(string(content), &conf); err != nil {
		return nil, fmt.Errorf( "decode toml: %w", err)
	}

	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf( "found invalid configuration: %w", err)
	}

	return &conf, nil
}

func (d *DependencyImpl) hasLockFile() bool {
	_, err := os.Stat(d.lockPath)
	return err == nil
}

func (d *DependencyImpl) IsNeedWriteLockFile() bool {
	return d.forceUpdate || !d.hasLockFile()
}
