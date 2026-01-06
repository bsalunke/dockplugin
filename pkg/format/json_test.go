package format

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/bsalunke/dockplugin/pkg/image"
)

func TestJSONFormatter(t *testing.T) {
	images := []*image.ImageInfo{
		{
			Repository:     "nginx",
			Tag:            "latest",
			ConfigSHA:      "sha256:abcdef1234567890",
			ManifestDigest: "sha256:1234567890abcdef",
			Size:           142000000,
			Created:        time.Now().Add(-48 * time.Hour),
			Architecture:   "amd64",
			OS:             "linux",
		},
	}

	var buf bytes.Buffer
	formatter := &JSONFormatter{Writer: &buf}

	err := formatter.Format(images)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify it's valid JSON
	var result []*image.ImageInfo
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 image, got %d", len(result))
	}

	if result[0].Repository != "nginx" {
		t.Errorf("Expected repository 'nginx', got '%s'", result[0].Repository)
	}
}
