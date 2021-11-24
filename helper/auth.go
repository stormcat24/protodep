package helper

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type authMethod string

const (
	SSHAgent authMethod = "SSHAgent"
	SSH                 = "SSH"
	HTTPS               = "HTTPS"
)

type authOptions struct {
	method   authMethod
	pemFile  string
	username string
	password string
}

type funcAuthOption struct {
	f func(options *authOptions)
}

func (fao *funcAuthOption) apply(do *authOptions) {
	fao.f(do)
}

type AuthOption interface {
	apply(*authOptions)
}

type AuthProvider interface {
	GetRepositoryURL(reponame string) string
	AuthMethod() (transport.AuthMethod, error)
}

type AuthProviderWithSSH struct {
	pemFile  string
	password string
}

type AuthProviderWithSSHAgent struct {
}

type AuthProviderHTTPS struct {
	username string
	password string
}

func WithHTTPS(username, password string) AuthOption {
	return &funcAuthOption{
		f: func(options *authOptions) {
			options.method = HTTPS
			options.username = username
			options.password = password
		},
	}
}

func WithPemFile(pemFile, password string) AuthOption {
	return &funcAuthOption{
		f: func(options *authOptions) {
			options.method = SSH
			options.password = password
			options.pemFile = pemFile
		},
	}
}

func NewAuthProvider(opt ...AuthOption) AuthProvider {
	opts := authOptions{
		method: SSHAgent,
	}
	for _, o := range opt {
		o.apply(&opts)
	}

	var authProvider AuthProvider
	if opts.method == SSHAgent {
		authProvider = &AuthProviderWithSSHAgent{}
	} else if opts.method == SSH {
		authProvider = &AuthProviderWithSSH{
			pemFile:  opts.pemFile,
			password: opts.password,
		}
	} else {
		authProvider = &AuthProviderHTTPS{
			username: opts.username,
			password: opts.password,
		}
	}

	return authProvider
}

func (p *AuthProviderWithSSH) GetRepositoryURL(reponame string) string {
	ep, err := transport.NewEndpoint("ssh://" + reponame + ".git")
	if err != nil {
		panic(err)
	}
	return ep.String()
}

func (p *AuthProviderWithSSH) AuthMethod() (transport.AuthMethod, error) {
	am, err := ssh.NewPublicKeysFromFile("git", p.pemFile, p.password)
	if err != nil {
		return nil, err
	}
	return am, nil
}

func (p *AuthProviderWithSSHAgent) GetRepositoryURL(reponame string) string {
	ep, err := transport.NewEndpoint("ssh://" + reponame + ".git")
	if err != nil {
		panic(err)
	}
	return ep.String()
}

func (p *AuthProviderWithSSHAgent) AuthMethod() (transport.AuthMethod, error) {
	aa, err := ssh.NewSSHAgentAuth(ssh.DefaultUsername)
	if err != nil {
		panic(err)
	}
	return aa, nil
}

func (p *AuthProviderHTTPS) GetRepositoryURL(reponame string) string {
	return fmt.Sprintf("https://%s.git", reponame)
}

func (p *AuthProviderHTTPS) AuthMethod() (transport.AuthMethod, error) {
	if p.username == "" && p.password == "" {
		return nil, nil
	}
	return &http.BasicAuth{
		Username: p.username,
		Password: p.password,
	}, nil
}
