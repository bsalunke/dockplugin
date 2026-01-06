# docker-img-sha

A Docker CLI plugin that addresses breaking changes in Docker Engine v29, where Image IDs changed from Config SHA to Manifest Digest. This plugin displays both identifiers to maintain backward compatibility for automation scripts and monitoring tools.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Examples](#examples)
- [Building from Source](#building-from-source)
- [Troubleshooting](#troubleshooting)
- [API Compatibility Matrix](#api-compatibility-matrix)
- [Contributing](#contributing)
- [License](#license)

## Overview

Docker Engine v29's transition to the containerd image store changed how Image IDs are displayed. Previously, the Image ID showed the **Config SHA** (a hash of the image configuration). Now it shows the **Manifest Digest** (a hash of the image manifest). This breaking change affects:

- Automation scripts that track images by Config SHA
- Monitoring tools that identify images
- CI/CD pipelines that depend on stable image identifiers
- Image vulnerability scanners

This plugin provides both identifiers, ensuring backward compatibility while supporting the new containerd image store.

## Features

✅ **Dual Identifier Display**: Shows both Config SHA and Manifest Digest  
✅ **Multiple Output Formats**: Table (default), JSON, and custom Go templates  
✅ **Image Filtering**: Filter by repository, tag, architecture  
✅ **Multi-Architecture Support**: Correctly handles OCI indexes  
✅ **API Version Negotiation**: Works with Docker Engine v25.0+  
✅ **Docker CLI Plugin Standard**: Integrates seamlessly with Docker CLI  
✅ **Cross-Platform**: Linux, macOS, Windows (x86_64, ARM64)

## Installation

### Method 1: Manual Installation

1. **Download the binary** for your platform from the [Releases](https://github.com/bsalunke/dockplugin/releases) page:

```bash
# Linux (x86_64)
curl -L https://github.com/bsalunke/dockplugin/releases/latest/download/docker-img-sha-linux-amd64 -o docker-img-sha

# Linux (ARM64)
curl -L https://github.com/bsalunke/dockplugin/releases/latest/download/docker-img-sha-linux-arm64 -o docker-img-sha

# macOS (x86_64)
curl -L https://github.com/bsalunke/dockplugin/releases/latest/download/docker-img-sha-darwin-amd64 -o docker-img-sha

# macOS (ARM64/Apple Silicon)
curl -L https://github.com/bsalunke/dockplugin/releases/latest/download/docker-img-sha-darwin-arm64 -o docker-img-sha
```

2. **Make it executable**:

```bash
chmod +x docker-img-sha
```

3. **Install as Docker CLI plugin**:

```bash
mkdir -p ~/.docker/cli-plugins
mv docker-img-sha ~/.docker/cli-plugins/
```

4. **Verify installation**:

```bash
docker img-sha version
```

### Method 2: Build from Source

See [Building from Source](#building-from-source) section below.

## Usage

The plugin is invoked as a Docker subcommand:

```bash
docker img-sha <command> [options]
```

## Commands

### `docker img-sha list`

List all local images with Config SHA and Manifest Digest.

**Syntax:**
```bash
docker img-sha list [OPTIONS]
```

**Options:**
- `--format, -f` - Output format: `table` (default), `json`, or custom Go template
- `--filter` - Filter images (e.g., `repository=nginx`)
- `--all, -a` - Show all images including intermediates
- `--no-trunc` - Don't truncate output (show full SHAs)
- `--quiet, -q` - Only show Config SHA values
- `--arch` - Filter by architecture (e.g., `amd64`, `arm64`)

**Example Output:**
```
REPOSITORY          TAG       CONFIG SHA        MANIFEST DIGEST   SIZE      CREATED
nginx               latest    sha256:a1b2c3...  sha256:d4e5f6... 142MB     2 days ago
postgres            16        sha256:g7h8i9...  sha256:j0k1l2... 379MB     1 week ago
redis               alpine    sha256:m3n4o5...  sha256:p6q7r8... 32.3MB    3 days ago
```

### `docker img-sha inspect`

Display detailed information about a specific image.

**Syntax:**
```bash
docker img-sha inspect <IMAGE> [OPTIONS]
```

**Options:**
- `--format, -f` - Output format: `json` (default)

**Arguments:**
- `<IMAGE>` - Image name, ID, or digest

### `docker img-sha version`

Display plugin version and compatibility information.

**Syntax:**
```bash
docker img-sha version
```

## Examples

### Basic Usage

**List all images:**
```bash
docker img-sha list
```

**List with full SHA values:**
```bash
docker img-sha list --no-trunc
```

**Show only Config SHAs (for scripting):**
```bash
docker img-sha list --quiet
```

### Filtering

**Filter by repository:**
```bash
docker img-sha list --filter "repository=nginx"
```

**Filter by tag:**
```bash
docker img-sha list --filter "tag=latest"
```

**Filter by architecture:**
```bash
docker img-sha list --arch arm64
```

### Output Formats

**JSON output:**
```bash
docker img-sha list --format json
```

**Custom template:**
```bash
docker img-sha list --format "{{.Repository}}:{{.Tag}} -> {{.ConfigSHA}}"
```

Output:
```
nginx:latest -> sha256:a1b2c3d4e5f6...
postgres:16 -> sha256:g7h8i9j0k1l2...
```

### Image Inspection

**Inspect a specific image:**
```bash
docker img-sha inspect nginx:latest
```

**Inspect by Config SHA:**
```bash
docker img-sha inspect sha256:a1b2c3d4e5f6...
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Docker (for testing)
- Make (optional, for convenience)

### Build Steps

1. **Clone the repository:**
```bash
git clone https://github.com/bsalunke/dockplugin.git
cd dockplugin
```

2. **Download dependencies:**
```bash
go mod download
```

3. **Build the binary:**
```bash
make build
# Or without Make:
go build -ldflags "-s -w" -o bin/docker-img-sha ./cmd/docker-img-sha
```

4. **Install locally:**
```bash
make install
# Or manually:
mkdir -p ~/.docker/cli-plugins
cp bin/docker-img-sha ~/.docker/cli-plugins/
chmod +x ~/.docker/cli-plugins/docker-img-sha
```

### Cross-Platform Builds

Build for all supported platforms:
```bash
make build-all
```

This creates binaries in the `bin/` directory:
- `docker-img-sha-linux-amd64`
- `docker-img-sha-linux-arm64`
- `docker-img-sha-darwin-amd64`
- `docker-img-sha-darwin-arm64`
- `docker-img-sha-windows-amd64.exe`

## Troubleshooting

### Docker Daemon Not Running

**Error:**
```
Error: Cannot connect to Docker daemon. Is the daemon running?
```

**Solution:**
- Ensure Docker Desktop or Docker daemon is running
- On Linux, check: `sudo systemctl status docker`
- Try: `sudo systemctl start docker`

### Permission Denied

**Error:**
```
Error: Permission denied while connecting to Docker daemon socket
```

**Solution:**
- Add your user to the `docker` group: `sudo usermod -aG docker $USER`
- Log out and back in, or run: `newgrp docker`
- Alternatively, run with sudo: `sudo docker img-sha list`

### Plugin Not Found

**Error:**
```
docker: 'img-sha' is not a docker command.
```

**Solution:**
- Verify installation: `ls -la ~/.docker/cli-plugins/docker-img-sha`
- Check permissions: `chmod +x ~/.docker/cli-plugins/docker-img-sha`
- Verify plugin is recognized: `docker plugin ls` (for extension plugins) or `docker --help` (should show img-sha)

### No Images Found / Empty Output

**Behavior:** The `list` command shows only headers, no images.

**Solution:**
- Pull some images first: `docker pull nginx`
- Verify images exist: `docker images`
- Check if you're listing all images: `docker img-sha list --all`

### Manifest Digest Shows "N/A"

**Behavior:** Manifest Digest column shows "N/A" for some images.

**Explanation:**
- Locally built images may not have a manifest digest until pushed to a registry
- Images pulled without digest reference may not have RepoDigests

**Solution:**
- For pushed images, pull with digest: `docker pull nginx@sha256:...`
- For local images, this is expected behavior

## API Compatibility Matrix

| Docker Engine Version | API Version | Compatibility | Notes |
|-----------------------|-------------|---------------|-------|
| v25.0 - v25.x | 1.44 | ✅ Fully Supported | Legacy image store |
| v26.0 - v26.x | 1.45 | ✅ Fully Supported | Hybrid mode |
| v27.0 - v27.x | 1.46 | ✅ Fully Supported | Containerd preview |
| v28.0 - v28.x | 1.47 | ✅ Fully Supported | Containerd beta |
| v29.0 - v29.x | 1.48 | ✅ Fully Supported | Containerd default |
| v30.0+ | 1.49+ | ✅ Expected | Future versions |

**Minimum Requirements:**
- Docker Engine: v25.0+
- Docker API: v1.44+

## Configuration

The plugin uses Docker's existing configuration and authentication. No separate configuration is required.

### Environment Variables

The plugin respects standard Docker environment variables:
- `DOCKER_HOST` - Docker daemon socket to connect to
- `DOCKER_API_VERSION` - Specify Docker API version
- `DOCKER_CERT_PATH` - Path to TLS certificates
- `DOCKER_TLS_VERIFY` - Enable TLS verification

Example:
```bash
export DOCKER_HOST=tcp://remote-docker:2376
docker img-sha list
```

## Performance

- **Image Listing**: < 2 seconds for systems with up to 500 images
- **Memory Usage**: < 100MB during operation
- **CPU Usage**: Minimal (< 10% during execution)

## Security

- ✅ Uses Docker daemon's existing authentication
- ✅ No transmission of data to external services
- ✅ No storage of credentials
- ✅ Read-only operations (no image modifications)

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Connection Error (Docker daemon not accessible) |
| 2 | API Version Error (unsupported API version) |
| 3 | Permission Error (insufficient permissions) |
| 4 | Image Not Found |
| 5 | Invalid Arguments |
| 99 | Internal Error |

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Docker team for the Docker Engine API
- Community feedback on [GitHub Issue #51779](https://github.com/moby/moby/issues/51779)
- OCI image specification maintainers

## Related Resources

- [Docker CLI Plugin Specification](https://docs.docker.com/engine/extend/cli_plugins/)
- [Docker Engine API Documentation](https://docs.docker.com/engine/api/)
- [OCI Image Specification](https://github.com/opencontainers/image-spec)
- [GitHub Issue: Image ID changed from Config SHA to Manifest Digest](https://github.com/moby/moby/issues/51779)

---

**Need help?** Open an issue on [GitHub](https://github.com/bsalunke/dockplugin/issues) or check the [Troubleshooting](#troubleshooting) section.
