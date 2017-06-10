package helper

import (
	"fmt"

	"github.com/stormcat24/protodep/logger"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type AuthProvider interface {
	GetRepositoryURL(reponame string) string
	AuthMethod() transport.AuthMethod
}

type AuthProviderWithSSH struct {
	pemFile string
}

type AuthProviderHTTPS struct {
}

func NewAuthProvider(pemFile string) AuthProvider {
	if pemFile != "" {
		logger.Info("use SSH protocol")
		return &AuthProviderWithSSH{
			pemFile: pemFile,
		}
	} else {
		logger.Info("use HTTP/HTTPS protocol")
		return &AuthProviderHTTPS{}
	}
}

func (p *AuthProviderWithSSH) GetRepositoryURL(reponame string) string {
	ep, err := transport.NewEndpoint("ssh://" + reponame + ".git")
	if err != nil {
		panic(err)
	}
	return ep.String()
}

func (p *AuthProviderWithSSH) AuthMethod() transport.AuthMethod {
	am, err := ssh.NewPublicKeysFromFile("git", p.pemFile, "")
	if err != nil {
		panic(err)
	}
	return am
}

func (p *AuthProviderHTTPS) GetRepositoryURL(reponame string) string {
	return fmt.Sprintf("https://%s.git", reponame)
}

func (p *AuthProviderHTTPS) AuthMethod() transport.AuthMethod {
	// nil is ok.
	return nil
}
