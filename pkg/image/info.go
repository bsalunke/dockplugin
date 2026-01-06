package image

import "time"

// ImageInfo represents comprehensive information about a Docker image
type ImageInfo struct {
	Repository     string            `json:"repository"`
	Tag            string            `json:"tag"`
	ConfigSHA      string            `json:"configSHA"`
	ManifestDigest string            `json:"manifestDigest"`
	Size           int64             `json:"size"`
	Created        time.Time         `json:"created"`
	Architecture   string            `json:"architecture"`
	OS             string            `json:"os"`
	IsMultiArch    bool              `json:"isMultiArch"`
	Labels         map[string]string `json:"labels,omitempty"`
}
