package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	// Version is the plugin version
	Version = "1.0.0"
	// SupportedAPIVersion is the minimum Docker API version required
	SupportedAPIVersion = "1.44"
	// SupportedEngineVersion is the minimum Docker Engine version required
	SupportedEngineVersion = "25.0+"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display plugin version and compatibility information",
		Long:  "Display the plugin version, Docker API version, and supported Docker Engine versions.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Plugin Version: %s\n", Version)
			fmt.Printf("Docker API Version: %s+\n", SupportedAPIVersion)
			fmt.Printf("Supported Engine Versions: %s\n", SupportedEngineVersion)
		},
	}

	return cmd
}
