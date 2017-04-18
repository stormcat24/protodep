package cmd

func init() {
	RootCmd.AddCommand(depCmd)
	initDepCmd()
}
