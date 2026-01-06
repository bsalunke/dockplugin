package image

import "strings"

// Filter represents filtering criteria for images
type Filter struct {
	Repository   string
	Tag          string
	Architecture string
}

// Matches checks if an ImageInfo matches the filter criteria
func (f *Filter) Matches(img *ImageInfo) bool {
	if f.Repository != "" {
		if !matchesPattern(img.Repository, f.Repository) {
			return false
		}
	}

	if f.Tag != "" {
		if !matchesPattern(img.Tag, f.Tag) {
			return false
		}
	}

	if f.Architecture != "" {
		if img.Architecture != f.Architecture {
			return false
		}
	}

	return true
}

// matchesPattern performs simple wildcard matching
func matchesPattern(value, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if !strings.Contains(pattern, "*") {
		return value == pattern
	}

	// Simple wildcard matching
	parts := strings.Split(pattern, "*")
	if len(parts) == 2 {
		prefix, suffix := parts[0], parts[1]
		return strings.HasPrefix(value, prefix) && strings.HasSuffix(value, suffix)
	}

	return strings.Contains(value, strings.Trim(pattern, "*"))
}
