package cmd

func init() {
	RootCmd.AddCommand(upCmd)
	initDepCmd()
}
