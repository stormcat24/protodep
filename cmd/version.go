package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/stormcat24/protodep/version"
)

const art = `
 ________  ________  ________  _________  ________  ________  _______   ________   
|\   __  \|\   __  \|\   __  \|\___   ___\\   __  \|\   ___ \|\  ___ \ |\   __  \  
\ \  \|\  \ \  \|\  \ \  \|\  \|___ \  \_\ \  \|\  \ \  \_|\ \ \   __/|\ \  \|\  \ 
 \ \   ____\ \   _  _\ \  \\\  \   \ \  \ \ \  \\\  \ \  \ \\ \ \  \_|/_\ \   ____\
  \ \  \___|\ \  \\  \\ \  \\\  \   \ \  \ \ \  \\\  \ \  \_\\ \ \  \_|\ \ \  \___|
   \ \__\    \ \__\\ _\\ \_______\   \ \__\ \ \_______\ \_______\ \_______\ \__\   
    \|__|     \|__|\|__|\|_______|    \|__|  \|_______|\|_______|\|_______|\|__|`

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Show protodep version",
	RunE: func(cdm *cobra.Command, args []string) error {
		fmt.Println(art)
		fmt.Println("")
		fmt.Println(version.Get())
		return nil
	},
}