package cmd

import (
	"log"
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
		log.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	DefaultConfig.limitDir = dir
	// fmt.Println("Current directory: ", DefaultConfig.limitDir)
}


var rootCmd = &cobra.Command{
  Use:   "trigger",
  Short: "Trigger is a CLI tool to trigger commands on a remote server",
  Long: `Trigger is a CLI tool to trigger commands on a remote server`,
  Run: func(cmd *cobra.Command, args []string) {
		log.Println("Trigger is a CLI tool to trigger commands on a remote server")
  },
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}