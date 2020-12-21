package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
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

	SetHttpsAuthProvider(provider helper.AuthProvider)
	SetSshAuthProvider(provider helper.AuthProvider)
}

type SyncImpl struct {
	conf *helper.SyncConfig

	httpsProvider helper.AuthProvider
	sshProvider   helper.AuthProvider
}

func NewSync(conf *helper.SyncConfig) (Sync, error) {
	s := &SyncImpl{
		conf: conf,
	}

	err := s.initAuthProviders()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SyncImpl) Resolve(forceUpdate bool, cleanupCache bool) error {

	dep := dependency.NewDependency(s.conf.TargetDir, forceUpdate)
	protodep, err := dep.Load()
	if err != nil {
		return err
	}

	newdeps := make([]dependency.ProtoDepDependency, 0, len(protodep.Dependencies))
	protodepDir := filepath.Join(s.conf.HomeDir, ".protodep")

	_, err = os.Stat(protodepDir)
	if cleanupCache && err == nil {
		files, err := ioutil.ReadDir(protodepDir)
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() {
				dirpath := filepath.Join(protodepDir, file.Name())
				if err := os.RemoveAll(dirpath); err != nil {
					return err
				}
			}
		}
	}

	outdir := filepath.Join(s.conf.OutputDir, protodep.ProtoOutdir)
	if err := os.RemoveAll(outdir); err != nil {
		return err
	}

	for _, dep := range protodep.Dependencies {
		var authProvider helper.AuthProvider

		if s.conf.UseHttps {
			authProvider = s.httpsProvider
		} else {
			switch dep.Agent {
			case "https":
				authProvider = s.httpsProvider
			case "ssh", "":
				authProvider = s.sshProvider
			default:
				return fmt.Errorf("%s agent is not accepted (ssh or https only)", dep.Agent)
			}
		}

		gitrepo := repository.NewGitRepository(protodepDir, dep, authProvider)

		repo, err := gitrepo.Open()
		if err != nil {
			return err
		}

		sources := make([]protoResource, 0)

		compiledIgnores := compileIngoresToGlob(dep.Ignores)

		protoRootDir := gitrepo.ProtoRootDir()
		filepath.Walk(protoRootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".proto") {
				if s.isIgnorePath(protoRootDir, path, dep.Ignores, compiledIgnores) {
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
			Agent:    repo.Dep.Agent,
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

func (s *SyncImpl) SetHttpsAuthProvider(provider helper.AuthProvider) {
	s.httpsProvider = provider
}

func (s *SyncImpl) SetSshAuthProvider(provider helper.AuthProvider) {
	s.sshProvider = provider
}

func (s *SyncImpl) initAuthProviders() error {
	s.httpsProvider = helper.NewAuthProvider(helper.WithHTTPS(s.conf.BasicAuthUsername, s.conf.BasicAuthPassword))

	if s.conf.IdentityFile == "" && s.conf.IdentityPassword == "" {
		s.sshProvider = helper.NewAuthProvider()

		return nil
	}

	identifyPath := filepath.Join(s.conf.HomeDir, ".ssh", s.conf.IdentityFile)
	isSSH, err := helper.IsAvailableSSH(identifyPath)
	if err != nil {
		return err
	}

	if isSSH {
		s.sshProvider = helper.NewAuthProvider(helper.WithPemFile(identifyPath, s.conf.IdentityPassword))
	} else {
		logger.Warn("The identity file path has been passed but is not available. Falling back to ssh-agent, the default authentication method.")
		s.sshProvider = helper.NewAuthProvider()
	}

	return nil
}

func compileIngoresToGlob(ignores []string) []glob.Glob {
	globIngores := make([]glob.Glob, len(ignores))

	for idx, ignore := range ignores {
		globIngores[idx] = glob.MustCompile(ignore)
	}

	return globIngores
}

func (s *SyncImpl) isIgnorePath(protoRootDir string, target string, ignores []string, globIgnores []glob.Glob) bool {
	// convert slashes otherwise doesnt work on windows same was as on linux
	target = filepath.ToSlash(target)

	// keeping old logic for backward compatibility
	for _, ignore := range ignores {
		// support windows paths correctly
		pathPrefix := filepath.ToSlash(filepath.Join(protoRootDir, ignore))
		if strings.HasPrefix(target, pathPrefix) {
			return true
		}
	}

	for _, ignore := range globIgnores {
		if ignore.Match(target) {
			return true
		}
	}

	return false
}
