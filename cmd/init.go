package cmd

func init() {
	RootCmd.AddCommand(upCmd, versionCmd)
	initDepCmd()
}
