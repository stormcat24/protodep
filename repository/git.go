package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"

	"github.com/stormcat24/protodep/dependency"
	"github.com/stormcat24/protodep/helper"
	"github.com/stormcat24/protodep/logger"
)

type GitRepository interface {
	Open() (*OpenedRepository, error)
	ProtoRootDir() string
}

type GitHubRepository struct {
	protodepDir  string
	dep          dependency.ProtoDepDependency
	authProvider helper.AuthProvider
}

func NewGitRepository(protodepDir string, dep dependency.ProtoDepDependency, authProvider helper.AuthProvider) GitRepository {
	return &GitHubRepository{
		protodepDir:  protodepDir,
		dep:          dep,
		authProvider: authProvider,
	}
}

type OpenedRepository struct {
	Repository *git.Repository
	Dep        dependency.ProtoDepDependency
	Hash       string
}

func (r *GitHubRepository) Open() (*OpenedRepository, error) {

	branch := "master"
	if r.dep.Branch != "" {
		branch = r.dep.Branch
	}

	revision := r.dep.Revision

	reponame := r.dep.Repository()
	repopath := filepath.Join(r.protodepDir, reponame)

	auth, err := r.authProvider.AuthMethod()
	if err != nil {
		return nil, err
	}

	var rep *git.Repository

	if stat, err := os.Stat(repopath); err == nil && stat.IsDir() {
		spinner := logger.InfoWithSpinner("Getting %s ", reponame)

		rep, err = git.PlainOpen(repopath)
		if err != nil {
			return nil, errors.Wrap(err, "open repository is failed")
		}
		spinner.Stop()

		fetchOpts := &git.FetchOptions{
			Auth: auth,
		}

		// TODO: Validate remote setting.
		// TODO: If .protodep cache remains with SSH, change remote target to HTTPS.

		if err := rep.Fetch(fetchOpts); err != nil {
			if err != git.NoErrAlreadyUpToDate {
				return nil, errors.Wrap(err, "fetch repository is failed")
			}
		}
		spinner.Finish()

	} else {
		spinner := logger.InfoWithSpinner("Getting %s ", reponame)
		// IDEA: Is it better to register both ssh and HTTP?
		rep, err = git.PlainClone(repopath, false, &git.CloneOptions{
			Auth: auth,
			URL:  r.authProvider.GetRepositoryURL(reponame),
		})
		if err != nil {
			return nil, errors.Wrap(err, "clone repository is failed")
		}
		spinner.Finish()
	}

	wt, err := rep.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "get worktree is failed")
	}

	if revision == "" {
		target, err := r.resolveReference(rep, branch)
		if err != nil {
			return nil, errors.Wrapf(err, "change branch to %s is failed", branch)
		}

		if err := wt.Checkout(&git.CheckoutOptions{Hash: target.Hash()}); err != nil {
			return nil, errors.Wrapf(err, "checkout to %s is failed", revision)
		}

		head := plumbing.NewHashReference(plumbing.HEAD, target.Hash())
		if err := rep.Storer.SetReference(head); err != nil {
			return nil, errors.Wrapf(err, "set head to %s is failed", branch)
		}
	} else {
		var opts git.CheckoutOptions

		tag := plumbing.NewTagReferenceName(revision)
		_, err := rep.Reference(tag, false)
		if err != nil && err != plumbing.ErrReferenceNotFound {
			return nil, errors.Wrapf(err, "tag = %s", tag)
		} else {
			if err != nil {
				// Tag not found, revision must be a hash
				logger.Info("%s is not a tag, checking out by hash", revision)
				hash := plumbing.NewHash(revision)
				opts = git.CheckoutOptions{Hash: hash}
			} else {
				logger.Info("%s is a tag, checking out by tag", revision)
				opts = git.CheckoutOptions{Branch: tag}
			}
		}

		if err := wt.Checkout(&opts); err != nil {
			return nil, errors.Wrapf(err, "checkout to %s is failed", revision)
		}
	}

	commiter, err := rep.Log(&git.LogOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get commit is failed")
	}

	current, err := commiter.Next()
	if err != nil {
		return nil, errors.Wrap(err, "get commit current is failed")
	}

	return &OpenedRepository{
		Repository: rep,
		Dep:        r.dep,
		Hash:       current.Hash.String(),
	}, nil
}

func (r *GitHubRepository) ProtoRootDir() string {
	return filepath.Join(r.protodepDir, r.dep.Target)
}

func (r *GitHubRepository) resolveReference(rep *git.Repository, branch string) (*plumbing.Reference, error) {
	if branch != "master" {
		return r.getReference(rep, branch)
	}
	// If master branch is failed, try main branch.
	target, err := r.getReference(rep, branch)
	if err == plumbing.ErrReferenceNotFound {
		return r.getReference(rep, "main")
	}
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (r *GitHubRepository) getReference(rep *git.Repository, branch string) (*plumbing.Reference, error) {
	return rep.Storer.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)))
}
