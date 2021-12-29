package resolver

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/stormcat24/protodep/pkg/auth"
	"github.com/stormcat24/protodep/pkg/config"
	"github.com/stormcat24/protodep/pkg/logger"
	"github.com/stormcat24/protodep/pkg/repository"
)

type protoResource struct {
	source       string
	relativeDest string
}

type Resolver interface {
	Resolve(forceUpdate bool, cleanupCache bool) error

	SetHttpsAuthProvider(provider auth.AuthProvider)
	SetSshAuthProvider(provider auth.AuthProvider)
}

type resolver struct {
	conf *Config

	httpsProvider auth.AuthProvider
	sshProvider   auth.AuthProvider
}

func New(conf *Config) (Resolver, error) {
	s := &resolver{
		conf: conf,
	}

	err := s.initAuthProviders()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *resolver) Resolve(forceUpdate bool, cleanupCache bool) error {

	dep := config.NewDependency(s.conf.TargetDir, forceUpdate)
	protodep, err := dep.Load()
	if err != nil {
		return err
	}

	newdeps := make([]config.ProtoDepDependency, 0, len(protodep.Dependencies))
	protodepDir := filepath.Join(s.conf.HomeDir, ".protodep")

	_, err = os.Stat(protodepDir)
	if cleanupCache && err == nil {
		files, err := os.ReadDir(protodepDir)
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
		var authProvider auth.AuthProvider

		if s.conf.UseHttps {
			authProvider = s.httpsProvider
		} else {
			switch dep.Protocol {
			case "https":
				authProvider = s.httpsProvider
			case "ssh", "":
				authProvider = s.sshProvider
			default:
				return fmt.Errorf("%s protocol is not accepted (ssh or https only)", dep.Protocol)
			}
		}

		gitrepo := repository.NewGit(protodepDir, dep, authProvider)

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

			content, err := os.ReadFile(s.source)
			if err != nil {
				return err
			}

			if err := writeFileWithDirectory(outpath, content, 0644); err != nil {
				return err
			}
		}

		newdeps = append(newdeps, config.ProtoDepDependency{
			Target:   repo.Dep.Target,
			Branch:   repo.Dep.Branch,
			Revision: repo.Hash,
			Path:     repo.Dep.Path,
			Ignores:  repo.Dep.Ignores,
			Protocol: repo.Dep.Protocol,
			Subgroup: repo.Dep.Subgroup,
		})
	}

	newProtodep := config.ProtoDep{
		ProtoOutdir:  protodep.ProtoOutdir,
		Dependencies: newdeps,
	}

	if dep.IsNeedWriteLockFile() {
		if err := writeToml("protodep.lock", newProtodep); err != nil {
			return err
		}
	}

	return nil
}

func (s *resolver) SetHttpsAuthProvider(provider auth.AuthProvider) {
	s.httpsProvider = provider
}

func (s *resolver) SetSshAuthProvider(provider auth.AuthProvider) {
	s.sshProvider = provider
}

func (s *resolver) initAuthProviders() error {
	s.httpsProvider = auth.NewAuthProvider(auth.WithHTTPS(s.conf.BasicAuthUsername, s.conf.BasicAuthPassword))

	if s.conf.IdentityFile == "" && s.conf.IdentityPassword == "" {
		s.sshProvider = auth.NewAuthProvider()

		return nil
	}

	identifyPath := filepath.Join(s.conf.HomeDir, ".ssh", s.conf.IdentityFile)
	isSSH, err := isAvailableSSH(identifyPath)
	if err != nil {
		return err
	}

	if isSSH {
		s.sshProvider = auth.NewAuthProvider(auth.WithPemFile(identifyPath, s.conf.IdentityPassword))
	} else {
		logger.Warn("The identity file path has been passed but is not available. Falling back to ssh-agent, the default authentication method.")
		s.sshProvider = auth.NewAuthProvider()
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

func (s *resolver) isIgnorePath(protoRootDir string, target string, ignores []string, globIgnores []glob.Glob) bool {
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

func writeToml(dest string, input interface{}) error {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(input); err != nil {
		return errors.Wrapf(err, "encode config to toml format is failed. [%+v]", input)
	}

	if err := os.WriteFile(dest, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "write to %s is failed", dest)
	}

	return nil
}

func writeFileWithDirectory(path string, data []byte, perm os.FileMode) error {

	path = filepath.ToSlash(path)
	s := strings.Split(path, "/")

	var dir string
	if len(s) > 1 {
		dir = strings.Join(s[0:len(s)-1], "/")
	} else {
		dir = path
	}

	dir = filepath.FromSlash(dir)
	path = filepath.FromSlash(path)

	if err := os.MkdirAll(dir, 0777); err != nil {
		return errors.Wrapf(err, "create directory is failed. [%s]", dir)
	}

	if err := os.WriteFile(path, data, perm); err != nil {
		return errors.Wrapf(err, "write data to file is failed. [%s]", path)
	}

	return nil
}

// isAvailableSSH is Check whether this machine can use git protocol
func isAvailableSSH(identifyPath string) (bool, error) {
	if _, err := os.Stat(identifyPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// TODO: validate ssh key
	return true, nil
}

