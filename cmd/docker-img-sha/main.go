package main

import (
	"fmt"
	"os"

	"github.com/bsalunke/dockplugin/pkg/commands"
	"github.com/spf13/cobra"
)

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "docker-img-sha",
		Short: "Docker CLI plugin for displaying Config SHA and Manifest Digest",
		Long: `Docker CLI plugin that addresses breaking changes in Docker Engine v29,
where Image IDs changed from Config SHA to Manifest Digest.

This plugin displays both identifiers to maintain backward compatibility
for automation scripts and monitoring tools.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add commands
	rootCmd.AddCommand(commands.NewListCommand())
	rootCmd.AddCommand(commands.NewInspectCommand())
	rootCmd.AddCommand(commands.NewVersionCommand())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
