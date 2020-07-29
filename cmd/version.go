package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   string
	goVersion string
	gitCommit string
	builtDate string
	builtOS   string
	builtArch string
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
	fmt.Println(fmt.Sprintf("Version:\t %s", version))
	fmt.Println(fmt.Sprintf("Go version:\t %s", goVersion))
	fmt.Println(fmt.Sprintf("Git commit:\t %s", gitCommit))
	fmt.Println(fmt.Sprintf("Built:\t\t %s", builtDate))
	fmt.Println(fmt.Sprintf("OS/Arch:\t %s/%s", builtOS, builtArch))

	return nil
}
