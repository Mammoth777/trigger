package cmd

import (
	"main/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().BoolP("deamon", "d", false, "Run the server as a deamon")
	serverCmd.PersistentFlags().BoolP("stop", "s", false, "Stop the server")
}

var serverCmd = &cobra.Command{
	Use: "serve",
	Short: "Start the server",
	Long: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		deamon, _ := cmd.Flags().GetBool("deamon")
		stop, _ := cmd.Flags().GetBool("stop")
		if deamon {
			server.StartServerDeamon()
			return
		}
		if stop {
			server.StopServer()
			return
		}
		server.StartServer()

	},
}