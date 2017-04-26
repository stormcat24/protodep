package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stormcat24/protodep/dependency"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type GitRepository interface {
	Open() (*git.Repository, error)
}

type GitHubRepository struct {
	protodepDir string
	dep         dependency.ProtoDepDependency
}

func NewGitRepository(protodepDir string, dep dependency.ProtoDepDependency) GitRepository {
	return &GitHubRepository{
		protodepDir: protodepDir,
		dep:         dep,
	}
}

func (r *GitHubRepository) Open() (*git.Repository, error) {

	branch := "master"
	if r.dep.Branch != "" {
		branch = r.dep.Branch
	}

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
			if err != git.NoErrAlreadyUpToDate {
				return nil, errors.Wrap(err, "fetch repository is failed")
			}
			fmt.Println("#NoErrAlreadyUpToDate")
		}

	} else {
		rep, err = git.PlainClone(repopath, false, &git.CloneOptions{
			URL: fmt.Sprintf("https://%s.git", reponame),
		})
		if err != nil {
			return nil, errors.Wrap(err, "clone repository is failed")
		}

	}

	wt, err := rep.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "get worktree is failed")
	}

	if branch != "master" {
		target, err := rep.Storer.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)))
		if err != nil {
			return nil, errors.Wrapf(err, "change branch to %s is failed", branch)
		}

		if err := wt.Checkout(target.Hash()); err != nil {
			return nil, errors.Wrapf(err, "checkout to %s is failed", revision)
		}

		head := plumbing.NewHashReference(plumbing.HEAD, target.Hash())
		if err := rep.Storer.SetReference(head); err != nil {
			return nil, errors.Wrapf(err, "set head to %s is failed", branch)
		}
	}

	if revision != "" {
		hash := plumbing.NewHash(revision)
		if err := wt.Checkout(hash); err != nil {
			return nil, errors.Wrapf(err, "checkout to %s is failed", revision)
		}

		head := plumbing.NewHashReference(plumbing.HEAD, hash)
		if err := rep.Storer.SetReference(head); err != nil {
			return nil, errors.Wrapf(err, "set head to %s is failed", revision)
		}
	}

	return rep, nil
}
