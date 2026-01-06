package image

import (
	"testing"
	"time"
)

func TestImageInfo(t *testing.T) {
	img := &ImageInfo{
		Repository:     "nginx",
		Tag:            "latest",
		ConfigSHA:      "sha256:abc123",
		ManifestDigest: "sha256:def456",
		Size:           142000000,
		Created:        time.Now(),
		Architecture:   "amd64",
		OS:             "linux",
		IsMultiArch:    false,
		Labels:         map[string]string{"version": "1.0"},
	}

	if img.Repository != "nginx" {
		t.Errorf("Expected repository to be 'nginx', got '%s'", img.Repository)
	}

	if img.ConfigSHA != "sha256:abc123" {
		t.Errorf("Expected ConfigSHA to be 'sha256:abc123', got '%s'", img.ConfigSHA)
	}
}
