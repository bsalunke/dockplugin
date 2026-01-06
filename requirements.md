# Requirements Document: Docker CLI Plugin for Image SHA and Digest Display

## 1. Executive Summary

### 1.1 Purpose
This document defines the requirements for a Docker CLI plugin that addresses the breaking changes introduced in Docker Engine v29, specifically the shift from local Config SHA-based Image IDs to deterministic Manifest Digest-based Image IDs. The plugin will provide backward-compatible image listing that displays both the legacy Image SHA (Config SHA) and the new Image Digest.

### 1.2 Problem Statement
Docker Engine v29's transition to the containerd image store has resulted in Image IDs now displaying the Manifest Digest instead of the Config SHA. This breaks existing automation scripts, monitoring tools, and workflows that depend on the distinct Config SHA for image identification and tracking.

### 1.3 Solution Overview
A Docker CLI plugin that wraps the existing Docker Engine API to retrieve and display both:
- **Image SHA (Config SHA)**: The content-addressable hash of the image configuration
- **Image Digest**: The manifest digest (which is now shown as the Image ID in v29)

---

## 2. Scope

### 2.1 In Scope
- CLI plugin compatible with Docker Engine v25.0+ (API v1.44+)
- Retrieval of both Config SHA and Manifest Digest for local images
- Enhanced `docker images` equivalent command with both identifiers
- Support for multi-architecture images (indexes)
- JSON and table output formats
- Filtering and sorting capabilities
- Compatibility with both legacy graph drivers and containerd image store

### 2.2 Out of Scope
- Modifications to Docker Engine core
- Remote registry operations beyond what's available via Docker API
- Image building or modification capabilities
- Support for Docker Engine versions prior to v25.0
- GUI or web interface

---

## 3. Stakeholders

| Role | Responsibility | Interest |
|------|---------------|----------|
| DevOps Engineers | Primary users of the plugin | Need backward compatibility for automation scripts |
| System Administrators | Image inventory management | Require consistent image identification across environments |
| CI/CD Pipeline Maintainers | Automated image tracking | Need stable identifiers for build artifacts |
| Security Teams | Image vulnerability tracking | Require precise image identification for scanning |

---

## 4. Functional Requirements

### 4.1 Core Functionality

#### FR-1: Image Listing with Dual Identifiers
**Priority**: MUST HAVE  
**Description**: The plugin shall retrieve and display a list of local Docker images showing both Config SHA and Manifest Digest.

**Acceptance Criteria**:
- Display repository, tag, Config SHA, Manifest Digest, size, and creation date
- Config SHA shall be extracted from image configuration blob
- Manifest Digest shall be extracted from image manifest
- Support for images stored in both legacy and containerd image stores

#### FR-2: API Version Compatibility
**Priority**: MUST HAVE  
**Description**: The plugin shall support Docker Engine API v1.44 and above.

**Acceptance Criteria**:
- Gracefully handle API version negotiation
- Display clear error messages for unsupported API versions
- Auto-detect available API version

#### FR-3: Multi-Architecture Image Support
**Priority**: MUST HAVE  
**Description**: The plugin shall correctly handle multi-platform images (OCI indexes).

**Acceptance Criteria**:
- Display index digest for multi-arch images
- Show Config SHA for each platform-specific manifest
- Indicate when an image is a multi-platform index
- Allow filtering by architecture

#### FR-4: Output Formatting
**Priority**: MUST HAVE  
**Description**: The plugin shall support multiple output formats.

**Acceptance Criteria**:
- Table format (human-readable, default)
- JSON format (machine-readable)
- Custom format using Go templates
- Wide format showing full SHA values (not truncated)

**Example Table Output**:
```
REPOSITORY          TAG       CONFIG SHA        MANIFEST DIGEST   SIZE      CREATED
nginx               latest    sha256:a1b2c3...  sha256:d4e5f6... 142MB     2 days ago
postgres            16        sha256:g7h8i9...  sha256:j0k1l2... 379MB     1 week ago
```

### 4.2 Filtering and Search

#### FR-5: Image Filtering
**Priority**: SHOULD HAVE  
**Description**: The plugin shall support filtering images based on various criteria.

**Acceptance Criteria**:
- Filter by repository name (with wildcard support)
- Filter by tag (including dangling images)
- Filter by creation date range
- Filter by size range
- Filter by architecture
- Support for combining multiple filters

#### FR-6: Search by Identifier
**Priority**: SHOULD HAVE  
**Description**: The plugin shall allow searching for images by Config SHA or Manifest Digest.

**Acceptance Criteria**:
- Accept full or truncated SHA values
- Support both SHA256 formats (with and without prefix)
- Display detailed information for matching images

### 4.3 Integration Features

#### FR-7: Docker CLI Plugin Standard
**Priority**: MUST HAVE  
**Description**: The plugin shall comply with Docker CLI plugin specifications.

**Acceptance Criteria**:
- Installable via `docker plugin install` or manual installation
- Invocable as `docker <plugin-name> <command>`
- Follow Docker CLI plugin naming conventions
- Support `--help` flag for command documentation

#### FR-8: Configuration Management
**Priority**: SHOULD HAVE  
**Description**: The plugin shall support configuration options.

**Acceptance Criteria**:
- Configuration file support (JSON/YAML)
- Environment variable overrides
- Command-line flag overrides
- Default output format configuration

---

## 5. Non-Functional Requirements

### 5.1 Performance

#### NFR-1: Response Time
**Priority**: MUST HAVE  
**Requirement**: The plugin shall list all local images within 2 seconds for systems with up to 500 images.

#### NFR-2: Resource Usage
**Priority**: SHOULD HAVE  
**Requirement**: The plugin shall consume no more than 100MB of memory during operation.

### 5.2 Compatibility

#### NFR-3: Platform Support
**Priority**: MUST HAVE  
**Requirement**: The plugin shall support the following platforms:
- Linux (x86_64, ARM64)
- macOS (x86_64, ARM64)
- Windows (x86_64)

#### NFR-4: Docker Version Support
**Priority**: MUST HAVE  
**Requirement**: The plugin shall support Docker Engine v25.0 through v29.x and beyond.

### 5.3 Usability

#### NFR-5: Error Messaging
**Priority**: MUST HAVE  
**Requirement**: The plugin shall provide clear, actionable error messages for common failure scenarios:
- Docker daemon not running
- API version incompatibility
- Permission denied
- Network connectivity issues

#### NFR-6: Documentation
**Priority**: MUST HAVE  
**Requirement**: The plugin shall include:
- README with installation instructions
- Usage examples for common scenarios
- Troubleshooting guide
- API compatibility matrix

### 5.4 Security

#### NFR-7: Credential Handling
**Priority**: MUST HAVE  
**Requirement**: The plugin shall use the Docker daemon's existing authentication without requiring separate credentials.

#### NFR-8: Data Privacy
**Priority**: MUST HAVE  
**Requirement**: The plugin shall not transmit image data or metadata to external services without explicit user consent.

---

## 6. Technical Requirements

### 6.1 Implementation

#### TR-1: Programming Language
**Requirement**: The plugin shall be implemented in Go (Golang) for consistency with Docker CLI.

#### TR-2: Docker API Client
**Requirement**: The plugin shall use the official Docker Engine API Go client library (`github.com/docker/docker/client`).

#### TR-3: Data Extraction Method
**Requirement**: The plugin shall extract identifiers using:
- **Config SHA**: From `ImageInspect.ID` or by parsing the image configuration blob
- **Manifest Digest**: From `ImageSummary.RepoDigests` or by querying the distribution manifest

### 6.2 API Interactions

#### TR-4: Required API Endpoints
The plugin shall interact with the following Docker Engine API endpoints:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/images/json` | GET | List images |
| `/images/{id}/json` | GET | Inspect image details |
| `/distribution/{name}/json` | GET | Get manifest digest |

#### TR-5: Handling Missing Data
**Requirement**: The plugin shall gracefully handle scenarios where:
- Config SHA is unavailable (display as "N/A")
- Manifest Digest is unavailable (display as "N/A")
- Image is corrupted or incomplete (display warning)

---

## 7. Command-Line Interface Specification

### 7.1 Plugin Name
**Proposed Name**: `docker-image-info` or `docker-img-sha`

### 7.2 Command Structure

```bash
docker <plugin-name> list [OPTIONS]
docker <plugin-name> inspect <IMAGE> [OPTIONS]
docker <plugin-name> version
```

### 7.3 Command: list

**Purpose**: List all local images with Config SHA and Manifest Digest

**Syntax**:
```bash
docker <plugin-name> list [OPTIONS]
```

**Options**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | `-f` | string | table | Output format (table, json, custom) |
| `--filter` | | string | | Filter images (e.g., "repository=nginx") |
| `--all` | `-a` | bool | false | Show all images (including intermediates) |
| `--digests` | | bool | true | Show digests (always true, for compatibility) |
| `--no-trunc` | | bool | false | Don't truncate output |
| `--quiet` | `-q` | bool | false | Only show image IDs |
| `--arch` | | string | | Filter by architecture |

**Examples**:
```bash
# Default table output
docker img-sha list

# JSON output
docker img-sha list --format json

# Filter by repository
docker img-sha list --filter "repository=nginx"

# Show full SHA values
docker img-sha list --no-trunc

# Custom format
docker img-sha list --format "{{.Repository}}:{{.Tag}} -> {{.ConfigSHA}}"
```

### 7.4 Command: inspect

**Purpose**: Display detailed information about a specific image

**Syntax**:
```bash
docker <plugin-name> inspect <IMAGE> [OPTIONS]
```

**Options**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | `-f` | string | json | Output format (json, yaml) |

**Output Fields**:
- Repository and tags
- Config SHA (full)
- Manifest Digest (full)
- Image size
- Architecture
- OS
- Layer count and sizes
- Creation timestamp
- Labels
- Environment variables

### 7.5 Command: version

**Purpose**: Display plugin version and compatibility information

**Output**:
```
Plugin Version: 1.0.0
Docker API Version: 1.44
Supported Engine Versions: 25.0+
```

---

## 8. Data Model

### 8.1 Image Information Structure

```go
type ImageInfo struct {
    Repository      string
    Tag             string
    ConfigSHA       string    // sha256:abc123... (Config digest)
    ManifestDigest  string    // sha256:def456... (Manifest digest)
    Size            int64
    Created         time.Time
    Architecture    string
    OS              string
    IsMultiArch     bool
    Labels          map[string]string
}
```

---

## 9. Error Handling

### 9.1 Error Categories

| Error Type | Exit Code | Example Message |
|------------|-----------|-----------------|
| Connection Error | 1 | "Cannot connect to Docker daemon. Is the daemon running?" |
| API Version Error | 2 | "Docker API version 1.43 is not supported. Requires 1.44+" |
| Permission Error | 3 | "Permission denied. Try running with sudo." |
| Image Not Found | 4 | "Image 'nginx:invalid' not found locally." |
| Invalid Arguments | 5 | "Invalid filter format. Use 'key=value'." |
| Internal Error | 99 | "Unexpected error: {details}" |

---

## 10. Testing Requirements

### 10.1 Unit Tests
- Config SHA extraction from image inspect
- Manifest Digest parsing from RepoDigests
- Filter logic validation
- Output formatting (table, JSON, custom)

### 10.2 Integration Tests
- List images with Docker v25, v27, v29
- Multi-architecture image handling
- Containerd image store compatibility
- Legacy graph driver compatibility

### 10.3 Test Coverage
**Requirement**: Minimum 80% code coverage

### 10.4 Test Scenarios

| Scenario | Expected Behavior |
|----------|-------------------|
| No local images | Display empty table with headers |
| Single-arch image | Show Config SHA and Manifest Digest |
| Multi-arch index | Show index digest, list architectures |
| Dangling images | Display `<none>` for repository/tag |
| Missing manifest digest | Display "N/A" in digest column |
| Large image set (1000+) | Complete within 5 seconds |

---

## 11. Deployment and Distribution

### 11.1 Installation Methods

#### Method 1: Manual Installation
```bash
# Download binary
curl -L https://github.com/org/docker-img-sha/releases/download/v1.0.0/docker-img-sha-linux-amd64 -o docker-img-sha

# Install
chmod +x docker-img-sha
mkdir -p ~/.docker/cli-plugins
mv docker-img-sha ~/.docker/cli-plugins/

# Verify
docker img-sha version
```

#### Method 2: Package Managers
- **Homebrew** (macOS): `brew install docker-img-sha`
- **APT** (Debian/Ubuntu): `apt install docker-img-sha`
- **YUM/DNF** (RHEL/Fedora): `yum install docker-img-sha`

### 11.2 Release Artifacts
- Binary executables for Linux (amd64, arm64)
- Binary executables for macOS (amd64, arm64)
- Binary executable for Windows (amd64)
- SHA256 checksums file
- Docker Hub image (optional containerized version)

---

## 12. Documentation Requirements

### 12.1 User Documentation
- **Installation Guide**: Step-by-step for all platforms
- **Quick Start**: Common use cases with examples
- **Command Reference**: Complete option documentation
- **Migration Guide**: For users upgrading from Docker v25-v28
- **Troubleshooting**: Common issues and solutions
- **FAQ**: Frequently asked questions

### 12.2 Developer Documentation
- **Architecture Overview**: Plugin design and components
- **Build Instructions**: How to build from source
- **Contributing Guide**: How to submit changes
- **API Reference**: Internal function documentation

---

## 13. Success Metrics

### 13.1 Adoption Metrics
- 1,000 downloads within first month
- 10 GitHub stars within first week
- Integration in at least 3 popular DevOps tools

### 13.2 Quality Metrics
- Zero critical bugs in first release
- Average issue resolution time < 48 hours
- User satisfaction rating > 4/5

### 13.3 Performance Metrics
- Image listing < 2 seconds for 500 images
- Memory usage < 100MB during operation
- CPU usage < 10% during execution

---

## 14. Risks and Mitigation

| Risk | Probability | Impact | Mitigation Strategy |
|------|-------------|--------|---------------------|
| Docker API changes in future versions | Medium | High | Use versioned API, implement graceful degradation |
| Config SHA unavailable in containerd store | Low | High | Document limitations, provide alternative identifiers |
| Performance degradation with large image sets | Medium | Medium | Implement pagination, optimize API calls |
| Security vulnerability in dependencies | Low | High | Regular dependency audits, automated security scanning |
| Incompatibility with third-party Docker distributions | Medium | Medium | Test with Docker Desktop, Rancher Desktop, Podman |

---

## 15. Timeline and Milestones

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| **Phase 1: Design & Prototyping** | 2 weeks | Architecture document, API proof-of-concept |
| **Phase 2: Core Development** | 4 weeks | Basic list and inspect commands, Config SHA extraction |
| **Phase 3: Feature Enhancement** | 3 weeks | Filtering, formatting, multi-arch support |
| **Phase 4: Testing & QA** | 2 weeks | Unit tests, integration tests, bug fixes |
| **Phase 5: Documentation** | 1 week | User docs, developer docs, examples |
| **Phase 6: Release** | 1 week | Binaries, distribution packages, announcement |

**Total Estimated Duration**: 13 weeks

---

## 16. Future Enhancements (Post-V1.0)

### 16.1 Planned Features
- **Image History with SHAs**: Show layer-by-layer Config SHAs
- **Diff Command**: Compare Config SHAs across environments
- **Export/Import**: Bulk export image metadata to CSV/JSON
- **Watch Mode**: Real-time monitoring of image changes
- **Remote Registry Support**: Query remote registries for digests
- **Integration with CI/CD**: Jenkins, GitLab CI, GitHub Actions plugins

### 16.2 Research Items
- Support for OCI image layout specification
- Integration with container security scanning tools
- Automated migration scripts for legacy workflows

---

## 17. Appendices

### Appendix A: Glossary

| Term | Definition |
|------|------------|
| **Config SHA** | SHA256 hash of the image configuration blob (legacy Image ID) |
| **Manifest Digest** | SHA256 hash of the image manifest (new default Image ID in v29) |
| **Containerd Image Store** | New default storage backend in Docker v29 |
| **Graph Driver** | Legacy storage backend (overlay2, btrfs, etc.) |
| **OCI Index** | Multi-platform image manifest |

### Appendix B: References

1. GitHub Issue: https://github.com/moby/moby/issues/51779
2. Docker Engine API Documentation: https://docs.docker.com/engine/api/
3. Docker CLI Plugin Specification: https://docs.docker.com/engine/extend/cli_plugins/
4. OCI Image Specification: https://github.com/opencontainers/image-spec

### Appendix C: Example API Responses

**ImageInspect Response (Partial)**:
```json
{
  "Id": "sha256:abc123...",
  "RepoTags": ["nginx:latest"],
  "RepoDigests": ["nginx@sha256:def456..."],
  "Config": {
    "Image": "sha256:abc123..."
  }
}
```

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-06 | [Author Name] | Initial requirements document |

**Approval**:
- [ ] Product Owner
- [ ] Technical Lead
- [ ] Security Team
- [ ] DevOps Team Lead

---

**Document Status**: Draft  
**Last Updated**: January 6, 2026  
**Next Review Date**: February 6, 2026
