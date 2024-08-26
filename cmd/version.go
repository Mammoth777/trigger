package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version of the application")
}

func PrintVersion() {
	fmt.Println("Trigger: Server Side Trigger: 1.0.0")
}



var versionCmd = &cobra.Command{
	Use:  "version",
	Aliases: []string{"-v"},
	Short: "Print the version of the application",
	Long: "Print the version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}