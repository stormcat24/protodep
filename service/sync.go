package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/stormcat24/protodep/dependency"
	"github.com/stormcat24/protodep/helper"
	"github.com/stormcat24/protodep/logger"
	"github.com/stormcat24/protodep/repository"
)

type protoResource struct {
	source       string
	relativeDest string
}

type Sync interface {
	Resolve(forceUpdate bool, cleanupCache bool) error
}

type SyncImpl struct {
	authProvider  helper.AuthProvider
	userHomeDir   string
	targetDir     string
	outputRootDir string
}

func NewSync(authProvider helper.AuthProvider, userHomeDir string, targetDir string, outputRootDir string) Sync {
	return &SyncImpl{
		authProvider:  authProvider,
		userHomeDir:   userHomeDir,
		targetDir:     targetDir,
		outputRootDir: outputRootDir,
	}
}

func (s *SyncImpl) Resolve(forceUpdate bool, cleanupCache bool) error {

	dep := dependency.NewDependency(s.targetDir, forceUpdate)
	protodep, err := dep.Load()
	if err != nil {
		return err
	}

	newdeps := make([]dependency.ProtoDepDependency, 0, len(protodep.Dependencies))
	protodepDir := filepath.Join(s.userHomeDir, ".protodep")

	if cleanupCache {
		files, _ := ioutil.ReadDir(protodepDir)
		if err == nil {
			for _, file := range files {
				if file.IsDir() {
					dirpath := filepath.Join(protodepDir, file.Name())
					if err := os.RemoveAll(dirpath); err != nil {
						return err
					}
				}
			}
		}
	}

	outdir := filepath.Join(s.outputRootDir, protodep.ProtoOutdir)
	if err := os.RemoveAll(outdir); err != nil {
		return err
	}

	for _, dep := range protodep.Dependencies {
		gitrepo := repository.NewGitRepository(protodepDir, dep, s.authProvider)

		repo, err := gitrepo.Open()
		if err != nil {
			return err
		}

		sources := make([]protoResource, 0)

		protoRootDir := gitrepo.ProtoRootDir()
		filepath.Walk(protoRootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".proto") {
				if s.isIgnorePath(protoRootDir, path, dep.Ignores) {
					logger.Info("skipped %s due to ignore setting", path)
				} else {
					sources = append(sources, protoResource{
						source:       path,
						relativeDest: strings.Replace(path, protoRootDir, "", -1),
					})
				}
			}
			return nil
		})

		for _, s := range sources {
			outpath := filepath.Join(outdir, dep.Path, s.relativeDest)

			content, err := ioutil.ReadFile(s.source)
			if err != nil {
				return err
			}

			if err := helper.WriteFileWithDirectory(outpath, content, 0644); err != nil {
				return err
			}
		}

		newdeps = append(newdeps, dependency.ProtoDepDependency{
			Target:   repo.Dep.Target,
			Branch:   repo.Dep.Branch,
			Revision: repo.Hash,
			Path:     repo.Dep.Path,
			Ignores:  repo.Dep.Ignores,
		})
	}

	newProtodep := dependency.ProtoDep{
		ProtoOutdir:  protodep.ProtoOutdir,
		Dependencies: newdeps,
	}

	if dep.IsNeedWriteLockFile() {
		if err := helper.WriteToml("protodep.lock", newProtodep); err != nil {
			return err
		}
	}

	return nil
}

func (s *SyncImpl) isIgnorePath(protoRootDir string, target string, ignores []string) bool {

	for _, ignore := range ignores {
		pathPrefix := filepath.Join(protoRootDir, ignore)
		if strings.HasPrefix(target, pathPrefix) {
			return true
		}
	}

	return false
}
