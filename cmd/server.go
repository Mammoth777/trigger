package cmd

import (
	"main/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start the server",
	Long: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}