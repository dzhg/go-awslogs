package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	version    = "dev"
	builtAt    = "unknown"
	commitSHA1 = "unknown"
)

func initVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of go-awslogs",
		Long:  "Show the version of go-awslogs, including the version number, build time and commit hash",
		Run:   printVersion,
	}
	return versionCmd
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("go-awslogs (%s %s)\nVersion: %s\nBuilt At: %s, Git Commit ID (SHA1): %s\n", runtime.GOOS, runtime.GOARCH, version, builtAt, commitSHA1)
}
