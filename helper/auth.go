package helper

import (
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
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
		return &AuthProviderWithSSH{
			pemFile: pemFile,
		}
	} else {
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
	am, err := gitssh.NewPublicKeysFromFile("git", p.pemFile, "")
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