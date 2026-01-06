package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bsalunke/dockplugin/pkg/client"
	"github.com/spf13/cobra"
)

// InspectOptions contains options for the inspect command
type InspectOptions struct {
	Format string
}

// NewInspectCommand creates the inspect command
func NewInspectCommand() *cobra.Command {
	opts := &InspectOptions{}

	cmd := &cobra.Command{
		Use:   "inspect <IMAGE> [OPTIONS]",
		Short: "Display detailed information about a specific image",
		Long:  "Display detailed information including Config SHA, Manifest Digest, architecture, and labels for a specific image.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInspect(args[0], opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Format, "format", "f", "json", "Output format (json)")

	return cmd
}

func runInspect(imageID string, opts *InspectOptions) error {
	// Create Docker client
	dc, err := client.NewDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot connect to Docker daemon. Is the daemon running?\n")
		return err
	}
	defer dc.Close()

	// Inspect image
	ctx := context.Background()
	img, err := dc.InspectImage(ctx, imageID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Image '%s' not found locally.\n", imageID)
		return err
	}

	// Format output
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(img)
}
