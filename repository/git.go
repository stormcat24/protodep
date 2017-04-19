package repository

import (
	"github.com/stormcat24/protodep/dependency"
	"gopkg.in/src-d/go-git.v4"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

	branch := r.dep.Branch
	revision := r.dep.Revision

	reponame := r.dep.Repository()
	repopath := filepath.Join(r.protodepDir, ".protodep/src", reponame)

	var rep *git.Repository

	if stat, err := os.Stat(repopath); err == nil && stat.IsDir() {
		rep, err = git.PlainOpen(repopath)
		if err != nil {
			return nil, errors.Wrap(err, "open repository is failed")
		}

		fetchOpts := &git.FetchOptions{
			//Auth: &gitssh.PublicKeys{
			//	User:   "git",
			//	Signer: signer,
			//},
			Progress: os.Stdout,
		}

		if err := rep.Fetch(fetchOpts); err != nil {
			return nil, errors.Wrap(err, "fetch repository is failed")
		}

		if revision != "" {

			wt, err := rep.Worktree()
			if err != nil {
				return nil, err
			}

			wt.Checkout()
		}


		//branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))

	} else {
		rep, err = git.PlainClone(repopath, false, &git.CloneOptions{
			URL:      fmt.Sprintf("https://%s.git", reponame),
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, errors.Wrap(err, "clone repository is failed")
		}

		if revision != "" {

		}

	}

	return rep, nil
}