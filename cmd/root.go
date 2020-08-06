package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "awslogs",
	Short: "awslogs is a tool for accessing AWS CloudWatch Logs",
}

// Execute executes the rootCmd
func Execute() {
	// do not sort the commands, keep the order as they are added
	cobra.EnableCommandSorting = false

	rootCmd.AddCommand(initGroupCmd())
	rootCmd.AddCommand(initStreamCmd())
	rootCmd.AddCommand(initGetCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
