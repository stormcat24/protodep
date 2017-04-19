package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/stormcat24/protodep/dependency"
	"github.com/stormcat24/protodep/repository"
)

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "get proto dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {

		user, _ := user.Current()

		isUpdate, err := cmd.Flags().GetBool("update")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Printf("isUpdate=%v\n", isUpdate)

		pwd, _ := os.Getwd()

		dep := dependency.NewDependency(pwd)
		protodep, err := dep.Load()
		if err != nil {
			return err
		}

		for _, dep := range protodep.Dependencies {
			gitrepo := repository.NewGitRepository(user.HomeDir, dep)

			repo, err := gitrepo.Open()
			if err != nil {
				return err
			}

			fmt.Println(repo)
		}

		return nil
	},
}

func initDepCmd() {
	depCmd.PersistentFlags().BoolP("update", "u", false, "update locked file and vendors")
}
