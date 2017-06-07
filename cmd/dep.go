package cmd

import (
	"os"

	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/stormcat24/protodep/dependency"
	"github.com/stormcat24/protodep/helper"
	"github.com/stormcat24/protodep/repository"
	"github.com/stormcat24/protodep/logger"
)

var unitTest = false

type protoResource struct {
	source       string
	relativeDest string
}

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "get proto dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {

		isUpdate, err := cmd.Flags().GetBool("update")
		if err != nil {
			return err
		}
		logger.Info("force update = %t", isUpdate)

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		homeDir, err := homedir.Dir()
		if err != nil {
			return err
		}

		authProvider := helper.NewAuthProvider(filepath.Join(homeDir, ".ssh", "id_rsa"))

		dep := dependency.NewDependency(pwd, isUpdate)
		protodep, err := dep.Load()
		if err != nil {
			return err
		}

		newdeps := make([]dependency.ProtoDepDependency, 0, len(protodep.Dependencies))
		protodepDir := filepath.Join(homeDir, ".protodep")

		for _, dep := range protodep.Dependencies {
			gitrepo := repository.NewGitRepository(protodepDir, dep, authProvider)

			repo, err := gitrepo.Open()
			if err != nil {
				return err
			}

			var outdirRoot string
			if unitTest {
				outdirRoot = os.TempDir()
			} else {
				outdirRoot = pwd
			}

			outdir := filepath.Join(outdirRoot, protodep.ProtoOutdir, dep.Repository())

			sources := make([]protoResource, 0)

			protoRootDir := gitrepo.ProtoRootDir()
			filepath.Walk(protoRootDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(path, ".proto") {
					sources = append(sources, protoResource{
						source:       path,
						relativeDest: strings.Replace(path, protoRootDir, "", -1),
					})
				}
				return nil
			})

			for _, s := range sources {
				outpath := filepath.Join(outdir, s.relativeDest)

				content, err := ioutil.ReadFile(s.source)
				if err != nil {
					return err
				}

				if err := os.MkdirAll(outdir, 0777); err != nil {
					return err
				}

				if err := ioutil.WriteFile(outpath, content, 0644); err != nil {
					return err
				}
			}

			newdeps = append(newdeps, dependency.ProtoDepDependency{
				Target:   repo.Dep.Target,
				Branch:   repo.Dep.Branch,
				Revision: repo.Hash,
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
	},
}

func initDepCmd() {
	depCmd.PersistentFlags().BoolP("update", "u", false, "update locked file and vendors")
}
