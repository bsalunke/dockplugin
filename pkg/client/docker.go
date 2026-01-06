package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bsalunke/dockplugin/pkg/image"
	"github.com/docker/docker/api/types"
	dockerimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// DockerClient wraps the Docker API client
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient creates a new Docker client with API negotiation
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerClient{cli: cli}, nil
}

// Close closes the Docker client connection
func (dc *DockerClient) Close() error {
	return dc.cli.Close()
}

// ListImages retrieves all local images with their Config SHA and Manifest Digest
func (dc *DockerClient) ListImages(ctx context.Context, all bool) ([]*image.ImageInfo, error) {
	images, err := dc.cli.ImageList(ctx, dockerimage.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	var result []*image.ImageInfo
	for _, img := range images {
		infos := dc.convertImageSummary(&img)
		result = append(result, infos...)
	}

	return result, nil
}

// InspectImage retrieves detailed information about a specific image
func (dc *DockerClient) InspectImage(ctx context.Context, imageID string) (*image.ImageInfo, error) {
	inspect, _, err := dc.cli.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image: %w", err)
	}

	return dc.convertImageInspect(&inspect), nil
}

// GetAPIVersion returns the Docker API version
func (dc *DockerClient) GetAPIVersion() string {
	return dc.cli.ClientVersion()
}

// convertImageSummary converts Docker API ImageSummary to our ImageInfo
func (dc *DockerClient) convertImageSummary(img *dockerimage.Summary) []*image.ImageInfo {
	var result []*image.ImageInfo

	// Extract Config SHA (Image ID)
	configSHA := img.ID

	// Extract Manifest Digest from RepoDigests
	manifestDigest := "N/A"
	if len(img.RepoDigests) > 0 {
		// RepoDigests format: "repository@sha256:digest"
		parts := strings.Split(img.RepoDigests[0], "@")
		if len(parts) == 2 {
			manifestDigest = parts[1]
		}
	}

	// Handle RepoTags
	created := time.Unix(img.Created, 0)
	if len(img.RepoTags) > 0 {
		for _, repoTag := range img.RepoTags {
			repo, tag := parseRepoTag(repoTag)
			info := &image.ImageInfo{
				Repository:     repo,
				Tag:            tag,
				ConfigSHA:      configSHA,
				ManifestDigest: manifestDigest,
				Size:           img.Size,
				Created:        created,
			}
			result = append(result, info)
		}
	} else {
		// Untagged image
		info := &image.ImageInfo{
			Repository:     "<none>",
			Tag:            "<none>",
			ConfigSHA:      configSHA,
			ManifestDigest: manifestDigest,
			Size:           img.Size,
			Created:        created,
		}
		result = append(result, info)
	}

	return result
}

// convertImageInspect converts Docker API ImageInspect to our ImageInfo
func (dc *DockerClient) convertImageInspect(inspect *types.ImageInspect) *image.ImageInfo {
	repo, tag := "<none>", "<none>"
	if len(inspect.RepoTags) > 0 {
		repo, tag = parseRepoTag(inspect.RepoTags[0])
	}

	manifestDigest := "N/A"
	if len(inspect.RepoDigests) > 0 {
		parts := strings.Split(inspect.RepoDigests[0], "@")
		if len(parts) == 2 {
			manifestDigest = parts[1]
		}
	}

	created, _ := time.Parse(time.RFC3339Nano, inspect.Created)

	return &image.ImageInfo{
		Repository:     repo,
		Tag:            tag,
		ConfigSHA:      inspect.ID,
		ManifestDigest: manifestDigest,
		Size:           inspect.Size,
		Created:        created,
		Architecture:   inspect.Architecture,
		OS:             inspect.Os,
		Labels:         inspect.Config.Labels,
	}
}

// parseRepoTag splits a repo:tag string into repository and tag
func parseRepoTag(repoTag string) (string, string) {
	parts := strings.Split(repoTag, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return repoTag, "latest"
}
