package format

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bsalunke/dockplugin/pkg/image"
)

// TableFormatter formats images as a table
type TableFormatter struct {
	NoTrunc bool
	Writer  io.Writer
}

// Format writes images in table format
func (f *TableFormatter) Format(images []*image.ImageInfo) error {
	w := tabwriter.NewWriter(f.Writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintln(w, "REPOSITORY\tTAG\tCONFIG SHA\tMANIFEST DIGEST\tSIZE\tCREATED")

	// Print each image
	for _, img := range images {
		configSHA := f.formatSHA(img.ConfigSHA)
		manifestDigest := f.formatSHA(img.ManifestDigest)
		size := f.formatSize(img.Size)
		created := f.formatTime(img.Created)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			img.Repository,
			img.Tag,
			configSHA,
			manifestDigest,
			size,
			created,
		)
	}

	return nil
}

// formatSHA truncates SHA if NoTrunc is false
func (f *TableFormatter) formatSHA(sha string) string {
	if sha == "N/A" || sha == "" {
		return "N/A"
	}
	if f.NoTrunc {
		return sha
	}
	// Show "sha256:" prefix + 12 characters
	if strings.HasPrefix(sha, "sha256:") {
		if len(sha) > 19 { // 7 (sha256:) + 12
			return sha[:19] + "..."
		}
	}
	return sha
}

// formatSize formats bytes to human-readable size
func (f *TableFormatter) formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatTime formats time to relative time string
func (f *TableFormatter) formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "Less than a minute ago"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
