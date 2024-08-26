package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	limitDir string;
}

var DefaultConfig = &Config{}


func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	DefaultConfig.limitDir = dir
	fmt.Println("Current directory: ", DefaultConfig.limitDir)
}


var rootCmd = &cobra.Command{
  Use:   "trigger",
  Short: "Trigger is a CLI tool to trigger commands on a remote server",
  Long: `Trigger is a CLI tool to trigger commands on a remote server`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Println("Trigger: Server Side Trigger: 1.0.0")
			os.Exit(0)
		}
	},
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
		fmt.Println("Executing command")
  },
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// server.StartServer()
}