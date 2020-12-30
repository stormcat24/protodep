package cmd

import (
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/stormcat24/protodep/helper"
	"github.com/stormcat24/protodep/logger"
	"github.com/stormcat24/protodep/service"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Populate .proto vendors existing protodep.toml and lock",
	RunE: func(cmd *cobra.Command, args []string) error {

		isForceUpdate, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}
		logger.Info("force update = %t", isForceUpdate)

		isCleanupCache, err := cmd.Flags().GetBool("cleanup")
		if err != nil {
			return err
		}
		logger.Info("cleanup cache = %t", isCleanupCache)

		identityFile, err := cmd.Flags().GetString("identity-file")
		if err != nil {
			return err
		}
		logger.Info("identity file = %s", identityFile)

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		if password != "" {
			logger.Info("password = %s", strings.Repeat("x", len(password))) // Do not display the password.
		}

		useHttps, err := cmd.Flags().GetBool("use-https")
		if err != nil {
			return err
		}
		logger.Info("use https = %t", useHttps)

		basicAuthUsername, err := cmd.Flags().GetString("basic-auth-username")
		if err != nil {
			return err
		}
		if basicAuthUsername != "" {
			logger.Info("https basic auth username = %s", basicAuthUsername)
		}

		basicAuthPassword, err := cmd.Flags().GetString("basic-auth-password")
		if err != nil {
			return err
		}
		if basicAuthPassword != "" {
			logger.Info("https basic auth password = %s", strings.Repeat("x", len(basicAuthPassword))) // Do not display the password.
		}

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		homeDir, err := homedir.Dir()
		if err != nil {
			return err
		}

		conf := helper.SyncConfig{
			UseHttps:          useHttps,
			HomeDir:           homeDir,
			TargetDir:         pwd,
			OutputDir:         pwd,
			BasicAuthUsername: basicAuthUsername,
			BasicAuthPassword: basicAuthPassword,
			IdentityFile:      identityFile,
			IdentityPassword:  password,
		}

		updateService, err := service.NewSync(&conf)
		if err != nil {
			return err
		}

		return updateService.Resolve(isForceUpdate, isCleanupCache)
	},
}

func initDepCmd() {
	upCmd.PersistentFlags().BoolP("force", "f", false, "update locked file and .proto vendors")
	upCmd.PersistentFlags().StringP("identity-file", "i", "", "set the identity file for SSH")
	upCmd.PersistentFlags().StringP("password", "p", "", "set the password for SSH")
	upCmd.PersistentFlags().BoolP("cleanup", "c", false, "cleanup cache before exec.")
	upCmd.PersistentFlags().BoolP("use-https", "u", false, "use HTTPS to get dependencies.")
	upCmd.PersistentFlags().StringP("basic-auth-username", "", "", "set the username with Basic Auth via HTTPS")
	upCmd.PersistentFlags().StringP("basic-auth-password", "", "", "set the password or personal access token(when enabled 2FA) with Basic Auth via HTTPS")
}
