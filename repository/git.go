package repository

import (
	"github.com/stormcat24/protodep/dependency"
	"gopkg.in/src-d/go-git.v4"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type GitRepository interface {
	Open() (*git.Repository, error)
}


type GitHubRepository struct {
	protodepDir string
	dep dependency.ProtoDepDependency
}

func NewGitRepository(protodepDir string, dep dependency.ProtoDepDependency) GitRepository {
	return &GitHubRepository{
		protodepDir: protodepDir,
		dep: dep,
	}
}

func (r *GitHubRepository) Open() (*git.Repository, error) {

	reponame := r.dep.Repository()
	repopath := filepath.Join(r.protodepDir, ".protodep/src", reponame)

	if stat, err := os.Stat(repopath); err == nil && stat.IsDir() {
		rep, err := git.PlainOpen(repopath)
		if err != nil {
			return nil, errors.Wrap(err, "open repository is failed")
		}
		return rep, nil

	} else {
		rep, err := git.PlainClone(repopath, false, &git.CloneOptions{
			URL:      fmt.Sprintf("https://%s.git", reponame),
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, errors.Wrap(err, "clone repository is failed")
		}

		return rep, nil
	}

}