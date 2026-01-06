package image

import "testing"

func TestFilterMatches(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		image    ImageInfo
		expected bool
	}{
		{
			name:   "No filters - should match",
			filter: Filter{},
			image: ImageInfo{
				Repository:   "nginx",
				Tag:          "latest",
				Architecture: "amd64",
			},
			expected: true,
		},
		{
			name: "Repository exact match",
			filter: Filter{
				Repository: "nginx",
			},
			image: ImageInfo{
				Repository: "nginx",
				Tag:        "latest",
			},
			expected: true,
		},
		{
			name: "Repository no match",
			filter: Filter{
				Repository: "postgres",
			},
			image: ImageInfo{
				Repository: "nginx",
				Tag:        "latest",
			},
			expected: false,
		},
		{
			name: "Repository wildcard match",
			filter: Filter{
				Repository: "ng*",
			},
			image: ImageInfo{
				Repository: "nginx",
				Tag:        "latest",
			},
			expected: true,
		},
		{
			name: "Tag exact match",
			filter: Filter{
				Tag: "latest",
			},
			image: ImageInfo{
				Repository: "nginx",
				Tag:        "latest",
			},
			expected: true,
		},
		{
			name: "Tag no match",
			filter: Filter{
				Tag: "alpine",
			},
			image: ImageInfo{
				Repository: "nginx",
				Tag:        "latest",
			},
			expected: false,
		},
		{
			name: "Architecture match",
			filter: Filter{
				Architecture: "amd64",
			},
			image: ImageInfo{
				Repository:   "nginx",
				Architecture: "amd64",
			},
			expected: true,
		},
		{
			name: "Architecture no match",
			filter: Filter{
				Architecture: "arm64",
			},
			image: ImageInfo{
				Repository:   "nginx",
				Architecture: "amd64",
			},
			expected: false,
		},
		{
			name: "Multiple filters all match",
			filter: Filter{
				Repository:   "nginx",
				Tag:          "latest",
				Architecture: "amd64",
			},
			image: ImageInfo{
				Repository:   "nginx",
				Tag:          "latest",
				Architecture: "amd64",
			},
			expected: true,
		},
		{
			name: "Multiple filters one no match",
			filter: Filter{
				Repository:   "nginx",
				Tag:          "alpine",
				Architecture: "amd64",
			},
			image: ImageInfo{
				Repository:   "nginx",
				Tag:          "latest",
				Architecture: "amd64",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.Matches(&tt.image)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		value    string
		pattern  string
		expected bool
	}{
		{"nginx", "nginx", true},
		{"nginx", "postgres", false},
		{"nginx", "*", true},
		{"nginx", "ng*", true},
		{"nginx", "*nx", true},
		{"nginx-alpine", "nginx*", true},
		{"my-nginx", "*nginx", true},
		{"my-nginx-app", "*nginx*", true},
		{"postgres", "ng*", false},
	}

	for _, tt := range tests {
		result := matchesPattern(tt.value, tt.pattern)
		if result != tt.expected {
			t.Errorf("matchesPattern(%q, %q) = %v, want %v",
				tt.value, tt.pattern, result, tt.expected)
		}
	}
}
