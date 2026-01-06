package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bsalunke/dockplugin/pkg/client"
	"github.com/bsalunke/dockplugin/pkg/format"
	"github.com/bsalunke/dockplugin/pkg/image"
	"github.com/spf13/cobra"
)

// ListOptions contains options for the list command
type ListOptions struct {
	Format       string
	Filter       string
	All          bool
	NoTrunc      bool
	Quiet        bool
	Architecture string
}

// NewListCommand creates the list command
func NewListCommand() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List all local images with Config SHA and Manifest Digest",
		Long:  "List all local Docker images displaying both the Config SHA (legacy Image ID) and Manifest Digest (new Image ID).",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Format, "format", "f", "table", "Output format (table, json, custom)")
	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter images (e.g., 'repository=nginx')")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show all images (including intermediates)")
	cmd.Flags().BoolVar(&opts.NoTrunc, "no-trunc", false, "Don't truncate output")
	cmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "Only show image IDs")
	cmd.Flags().StringVar(&opts.Architecture, "arch", "", "Filter by architecture")

	return cmd
}

func runList(opts *ListOptions) error {
	// Create Docker client
	dc, err := client.NewDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot connect to Docker daemon. Is the daemon running?\n")
		return err
	}
	defer dc.Close()

	// List images
	ctx := context.Background()
	images, err := dc.ListImages(ctx, opts.All)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list images: %v\n", err)
		return err
	}

	// Apply filters
	if opts.Filter != "" || opts.Architecture != "" {
		filter := parseFilter(opts.Filter, opts.Architecture)
		images = filterImages(images, filter)
	}

	// Quiet mode - only show Config SHA
	if opts.Quiet {
		for _, img := range images {
			fmt.Println(img.ConfigSHA)
		}
		return nil
	}

	// Format output
	switch strings.ToLower(opts.Format) {
	case "json":
		formatter := &format.JSONFormatter{Writer: os.Stdout}
		return formatter.Format(images)
	case "table":
		formatter := &format.TableFormatter{
			NoTrunc: opts.NoTrunc,
			Writer:  os.Stdout,
		}
		return formatter.Format(images)
	default:
		// Custom template format
		formatter := &format.CustomFormatter{
			Template: opts.Format,
			Writer:   os.Stdout,
		}
		return formatter.Format(images)
	}
}

func parseFilter(filterStr, arch string) *image.Filter {
	filter := &image.Filter{
		Architecture: arch,
	}

	if filterStr == "" {
		return filter
	}

	// Parse filter string (e.g., "repository=nginx" or "tag=latest")
	parts := strings.SplitN(filterStr, "=", 2)
	if len(parts) == 2 {
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "repository":
			filter.Repository = value
		case "tag":
			filter.Tag = value
		}
	}

	return filter
}

func filterImages(images []*image.ImageInfo, filter *image.Filter) []*image.ImageInfo {
	var result []*image.ImageInfo
	for _, img := range images {
		if filter.Matches(img) {
			result = append(result, img)
		}
	}
	return result
}
