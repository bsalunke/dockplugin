package format

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bsalunke/dockplugin/pkg/image"
)

func TestTableFormatter(t *testing.T) {
	images := []*image.ImageInfo{
		{
			Repository:     "nginx",
			Tag:            "latest",
			ConfigSHA:      "sha256:abcdef1234567890",
			ManifestDigest: "sha256:1234567890abcdef",
			Size:           142000000,
			Created:        time.Now().Add(-48 * time.Hour),
		},
		{
			Repository:     "postgres",
			Tag:            "16",
			ConfigSHA:      "sha256:fedcba0987654321",
			ManifestDigest: "sha256:0987654321fedcba",
			Size:           379000000,
			Created:        time.Now().Add(-7 * 24 * time.Hour),
		},
	}

	var buf bytes.Buffer
	formatter := &TableFormatter{
		NoTrunc: false,
		Writer:  &buf,
	}

	err := formatter.Format(images)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "REPOSITORY") {
		t.Error("Output should contain header 'REPOSITORY'")
	}
	if !strings.Contains(output, "nginx") {
		t.Error("Output should contain 'nginx'")
	}
	if !strings.Contains(output, "postgres") {
		t.Error("Output should contain 'postgres'")
	}
}

func TestTableFormatterNoTrunc(t *testing.T) {
	images := []*image.ImageInfo{
		{
			Repository:     "nginx",
			Tag:            "latest",
			ConfigSHA:      "sha256:abcdef1234567890abcdef1234567890",
			ManifestDigest: "sha256:1234567890abcdef1234567890abcdef",
			Size:           142000000,
			Created:        time.Now(),
		},
	}

	var buf bytes.Buffer
	formatter := &TableFormatter{
		NoTrunc: true,
		Writer:  &buf,
	}

	err := formatter.Format(images)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	// With NoTrunc, full SHA should be visible
	if !strings.Contains(output, "abcdef1234567890abcdef1234567890") {
		t.Error("Output should contain full ConfigSHA when NoTrunc is true")
	}
}

func TestFormatSize(t *testing.T) {
	formatter := &TableFormatter{}

	tests := []struct {
		bytes    int64
		expected string
	}{
		{512, "512B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1048576, "1.0MB"},
		{142000000, "135.4MB"},
		{1073741824, "1.0GB"},
	}

	for _, tt := range tests {
		result := formatter.formatSize(tt.bytes)
		if result != tt.expected {
			t.Errorf("formatSize(%d) = %s, want %s", tt.bytes, result, tt.expected)
		}
	}
}

func TestFormatTime(t *testing.T) {
	formatter := &TableFormatter{}
	now := time.Now()

	tests := []struct {
		time     time.Time
		contains string
	}{
		{now.Add(-30 * time.Second), "Less than a minute"},
		{now.Add(-2 * time.Minute), "2 minutes ago"},
		{now.Add(-1 * time.Hour), "1 hour ago"},
		{now.Add(-3 * time.Hour), "3 hours ago"},
		{now.Add(-1 * 24 * time.Hour), "1 day ago"},
		{now.Add(-3 * 24 * time.Hour), "3 days ago"},
		{now.Add(-8 * 24 * time.Hour), "1 week ago"},
		{now.Add(-15 * 24 * time.Hour), "2 weeks ago"},
	}

	for _, tt := range tests {
		result := formatter.formatTime(tt.time)
		if !strings.Contains(result, tt.contains) {
			t.Errorf("formatTime() = %s, want to contain %s", result, tt.contains)
		}
	}
}
