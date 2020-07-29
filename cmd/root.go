// Package cmd contains the functionality for the set of commands
// currently supported by the CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	// rootCmd represents the base command when called without any sub commands.
	rootCmd := &cobra.Command{
		Use:   "json-server",
		Short: "Create a dummy REST API from a json file with zero coding within seconds",
		Long:  "json-server is a cross-platform CLI tool to create within seconds a dummy REST API from a json file",
	}

	// Add sub commands to base command.
	rootCmd.AddCommand(newStartCmd())

	return rootCmd
}
