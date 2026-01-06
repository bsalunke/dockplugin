package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bsalunke/dockplugin/pkg/commands"
	"github.com/spf13/cobra"
)

// PluginMetadata represents Docker CLI plugin metadata
type PluginMetadata struct {
	SchemaVersion    string `json:"SchemaVersion"`
	Vendor           string `json:"Vendor"`
	Version          string `json:"Version"`
	ShortDescription string `json:"ShortDescription"`
	URL              string `json:"URL,omitempty"`
}

func main() {
	// Handle docker-cli-plugin-metadata command
	if len(os.Args) > 1 && os.Args[1] == "docker-cli-plugin-metadata" {
		metadata := PluginMetadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "dockplugin",
			Version:          "1.0.0",
			ShortDescription: "Display Docker images with Config SHA and Manifest Digest",
			URL:              "https://github.com/bsalunke/dockplugin",
		}
		json.NewEncoder(os.Stdout).Encode(metadata)
		return
	}

	// Debug: Print args (remove this after testing)
	//fmt.Fprintf(os.Stderr, "DEBUG: os.Args = %v\n", os.Args)

	// When invoked as a Docker CLI plugin, the first argument after program name is the plugin name
	// Strip it if present (e.g., "docker-imgsha imgsha list" becomes "docker-imgsha list")
	if len(os.Args) > 1 {
		// Check if the first arg is a plugin invocation name
		if os.Args[1] == "imgsha" || os.Args[1] == "img-sha" || os.Args[1] == "docker-imgsha" || os.Args[1] == "docker-img-sha" {
			os.Args = append(os.Args[:1], os.Args[2:]...)
		}
	}

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "docker-imgsha",
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
