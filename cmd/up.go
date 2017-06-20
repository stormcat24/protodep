package cmd

import (
	"os"

	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/stormcat24/protodep/helper"
	"github.com/stormcat24/protodep/logger"
	"github.com/stormcat24/protodep/service"
)

var (
	authProvider helper.AuthProvider
)

type protoResource struct {
	source       string
	relativeDest string
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Populate .proto vendors existing protodep.toml and lock",
	RunE: func(cmd *cobra.Command, args []string) error {

		isForceUpdate, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}
		logger.Info("force update = %t", isForceUpdate)

		identityFile, err := cmd.Flags().GetString("identity-file")
		if err != nil {
			return err
		}
		logger.Info("identity file = %s", identityFile)

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		homeDir, err := homedir.Dir()
		if err != nil {
			return err
		}

		authProvider = helper.NewAuthProvider(filepath.Join(homeDir, ".ssh", identityFile))
		updateService := service.NewSync(authProvider, homeDir, pwd, pwd)
		return updateService.Resolve(isForceUpdate)
	},
}

func initDepCmd() {
	upCmd.PersistentFlags().BoolP("force", "f", false, "update locked file and .proto vendors")
	upCmd.PersistentFlags().StringP("identity-file", "i", "id_rsa", "set the identity file for SSH")
}
