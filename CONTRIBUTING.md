# Contributing to docker-imgsha

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Code of Conduct

Be respectful and inclusive. We welcome contributions from everyone.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Docker version, plugin version)

### Suggesting Features

Feature suggestions are welcome! Please open an issue describing:
- The use case
- Why the feature would be useful
- How it might work

### Pull Requests

1. **Fork the repository** and create a new branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for new functionality
4. **Run tests** and ensure they pass: `make test`
5. **Run linters**: `make lint` (or `go fmt ./...`)
6. **Update documentation** if needed
7. **Commit your changes** with clear commit messages
8. **Push to your fork** and submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Docker (for testing)
- Make (optional)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/dockplugin.git
cd dockplugin

# Download dependencies
go mod download

# Build
make build

# Run tests
make test

# Install locally for testing
make install
```

## Coding Standards

- Follow standard Go conventions and idioms
- Run `go fmt` before committing
- Write clear, descriptive variable and function names
- Add comments for complex logic
- Keep functions focused and reasonably sized

## Testing

- Write unit tests for new functionality
- Ensure tests are deterministic
- Test edge cases and error conditions
- Aim for >80% code coverage

Run tests:
```bash
make test
```

Check coverage:
```bash
make coverage
```

## Commit Messages

Use clear, descriptive commit messages:
- Start with a verb in present tense (Add, Fix, Update, etc.)
- Keep the first line under 72 characters
- Add details in the body if needed

Examples:
```
Add filtering by architecture
Fix manifest digest parsing for multi-arch images
Update README with troubleshooting section
```

## Project Structure

```
dockplugin/
├── cmd/docker-imgsha/    # Main application entry point
├── pkg/
│   ├── client/            # Docker client wrapper
│   ├── commands/          # CLI commands
│   ├── format/            # Output formatters
│   └── image/             # Image data structures
├── .github/workflows/     # CI/CD workflows
├── Makefile               # Build automation
└── README.md              # Documentation
```

## Release Process

Releases are managed by maintainers. To propose a release:

1. Ensure all tests pass
2. Update version in `pkg/commands/version.go`
3. Update CHANGELOG.md (if it exists)
4. Create a pull request with the version update

## Questions?

If you have questions, feel free to:
- Open an issue
- Start a discussion on GitHub
- Reach out to maintainers

Thank you for contributing!
