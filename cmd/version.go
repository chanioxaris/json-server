package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "unknown"
	goVersion = "unknown"
	gitCommit = "unknown"
	builtDate = "unknown"
	builtOS   = "unknown"
	builtArch = "unknown"
)

func newVersionCmd() *cobra.Command {
	// versionCmd represents the version command.
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Show version information",
		RunE:  runVersion,
	}

	return versionCmd
}
func runVersion(_ *cobra.Command, _ []string) error {
	fmt.Printf("Version:\t %s\n", version)
	fmt.Printf("Go version:\t %s\n", goVersion)
	fmt.Printf("Git commit:\t %s\n", gitCommit)
	fmt.Printf("Built:\t\t %s\n", builtDate)
	fmt.Printf("OS/Arch:\t %s/%s\n", builtOS, builtArch)

	return nil
}
