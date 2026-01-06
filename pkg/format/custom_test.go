package format

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bsalunke/dockplugin/pkg/image"
)

func TestCustomFormatter(t *testing.T) {
	images := []*image.ImageInfo{
		{
			Repository:     "nginx",
			Tag:            "latest",
			ConfigSHA:      "sha256:abcdef1234567890",
			ManifestDigest: "sha256:1234567890abcdef",
			Size:           142000000,
			Created:        time.Now(),
		},
	}

	tests := []struct {
		name     string
		template string
		contains string
	}{
		{
			name:     "Simple template",
			template: "{{.Repository}}:{{.Tag}}",
			contains: "nginx:latest",
		},
		{
			name:     "ConfigSHA template",
			template: "{{.ConfigSHA}}",
			contains: "sha256:abcdef1234567890",
		},
		{
			name:     "Complex template",
			template: "{{.Repository}}:{{.Tag}} -> {{.ConfigSHA}}",
			contains: "nginx:latest -> sha256:abcdef1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter := &CustomFormatter{
				Template: tt.template,
				Writer:   &buf,
			}

			err := formatter.Format(images)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Output should contain '%s', got: %s", tt.contains, output)
			}
		})
	}
}

func TestCustomFormatterInvalidTemplate(t *testing.T) {
	images := []*image.ImageInfo{
		{
			Repository: "nginx",
			Tag:        "latest",
		},
	}

	var buf bytes.Buffer
	formatter := &CustomFormatter{
		Template: "{{.InvalidField}}",
		Writer:   &buf,
	}

	err := formatter.Format(images)
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}
}
