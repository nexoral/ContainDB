# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ContainDB** - Database Container Management Made Simple

- **Language**: Go (Golang) ≥1.18
- **Type**: CLI Tool for Docker-based database management
- **Platform**: Linux, macOS, Windows
- **Distribution**: npm package + native binaries (.deb, direct binary)
- **Purpose**: Simplify containerized database setup and management

## Commands

```bash
# Build
go build -o containdb src/Core/main.go      # Build CLI binary
./Scripts/BinBuilder.sh                      # Build for all platforms
./Scripts/PackageBuilder.sh                  # Build Debian package

# Development
go run src/Core/main.go                      # Run directly
go test ./...                                # Run all tests
go fmt ./...                                 # Format code
go vet ./...                                 # Lint code

# Distribution
npm run build                                # Build npm package
npm publish                                  # Publish to npm
```

## Core Rules (NON-NEGOTIABLE)

1. **ALWAYS test**: Test on Linux, macOS, Windows before release
2. **ALWAYS build**: Run `go build` after code changes
3. **NEVER break compatibility**: Support Go 1.18+, maintain CLI interface
4. **Respect Docker**: All operations use Docker API, no direct container manipulation
5. **Error handling**: Always provide clear, actionable error messages
6. **Update docs**: README.md, DESIGN.md, INSTALLATION.md when features change

## Critical Constraints

### Cross-Platform Compatibility
- **Must work on**: Linux, macOS, Windows
- **Go version**: ≥1.18 (for generics and performance)
- **Docker**: Docker Engine 20.10+
- **No platform-specific code** without fallbacks

### User Experience
- **Interactive prompts** for all sensitive inputs (credentials, ports)
- **Smart defaults** for common scenarios
- **Clear error messages** with actionable solutions
- **Auto-rollback** on failures (cleanup partial state)

## Architecture Overview

See `DESIGN.md` for complete details.

### Layered Architecture
```
CLI Entry Point (main.go)
    ↓
Base Operations Package (src/base/)
    ├── Database creation logic
    ├── Container lifecycle management
    └── Volume & network handling
    ↓
Docker Abstraction (src/Docker/)
    ├── Docker SDK wrapper
    ├── Container operations
    ├── Volume operations
    └── Network operations
    ↓
Docker Engine
```

### Module Structure
```
src/
├── Core/           # main.go - CLI entry point
├── base/           # Core business logic
│   ├── Create.go           # Database creation
│   ├── Remove.go           # Resource cleanup
│   ├── Export.go           # docker-compose export
│   └── Import.go           # docker-compose import
├── Docker/         # Docker API abstraction
│   ├── Container.go        # Container operations
│   ├── Volume.go           # Volume management
│   └── Network.go          # Network operations
└── tools/          # Management tools (phpMyAdmin, pgAdmin, etc.)

Scripts/
├── BinBuilder.sh          # Multi-platform binary builder
├── PackageBuilder.sh      # Debian package builder
└── installer.sh           # Linux installer script
```

## Go Standards

### Error Handling
```go
// ✅ GOOD - Descriptive errors with context
if err != nil {
    return fmt.Errorf("failed to create container '%s': %w", name, err)
}

// ❌ BAD - Generic errors
if err != nil {
    return err
}
```

### Package Organization
```go
// ✅ GOOD - Clear package structure
package docker

import (
    "github.com/docker/docker/client"
)

func CreateContainer(opts ContainerOptions) error { }

// ❌ BAD - Mixed concerns
package main

func CreateContainer() error { }
func ParseConfig() error { }
```

### Type Safety
```go
// ✅ GOOD - Use structs for options
type ContainerOptions struct {
    Name          string
    Image         string
    Port          int
    RestartPolicy string
}

func CreateContainer(opts ContainerOptions) error { }

// ❌ BAD - Too many parameters
func CreateContainer(name, image string, port int, restart string) error { }
```

## Key Patterns

### Auto-Rollback on Failure
```go
func CreateDatabase(opts Options) error {
    // Track created resources
    created := []string{}

    // Cleanup on error
    defer func() {
        if err != nil {
            for _, resource := range created {
                cleanup(resource)
            }
        }
    }()

    // Create resources
    container, err := CreateContainer(opts)
    if err != nil {
        return err
    }
    created = append(created, container.ID)

    // Continue...
}
```

### Interactive Prompts
```go
// ✅ Always prompt for sensitive data
password := promptSecure("Enter database password:")

// ❌ Never hardcode or use default passwords
password := "password123"
```

### Docker API Usage
```go
// ✅ Use official Docker SDK
import "github.com/docker/docker/client"

cli, err := client.NewClientWithOpts(client.FromEnv)
if err != nil {
    return fmt.Errorf("failed to create Docker client: %w", err)
}

// ❌ Don't execute docker commands directly
exec.Command("docker", "run", ...)
```

## Documentation Requirements

**Update when features change**:
1. **README.md** - Usage examples, installation, features
2. **DESIGN.md** - Architecture diagrams, technical decisions
3. **INSTALLATION.md** - Platform-specific installation instructions
4. **CONTRIBUTING.md** - Development setup, guidelines
5. **Code comments** - All exported functions and complex logic

## Security

1. **Credential Handling**: Never log or display passwords
2. **Privilege Management**: Warn users about root requirements
3. **Input Validation**: Validate all user inputs (ports, names, paths)
4. **Container Isolation**: Use Docker networks for container communication
5. **No Hardcoded Secrets**: All credentials via prompts or env vars

## Testing

- **Unit tests**: Core logic in `src/base/`, `src/Docker/`
- **Integration tests**: Full CLI workflow tests
- **Cross-platform**: Test on Linux, macOS, Windows
- **Docker scenarios**: Test with/without existing containers, volumes, networks
- **Error cases**: Network failures, permission errors, conflicts

## Supported Databases

- **SQL**: MySQL, PostgreSQL, MariaDB
- **NoSQL**: MongoDB, Redis
- **Management Tools**: phpMyAdmin, pgAdmin, RedisInsight, MongoDB Compass

## Build & Distribution

### Binary Distribution
```bash
# Build for all platforms
./Scripts/BinBuilder.sh

# Output:
# - bin/containdb-linux-amd64
# - bin/containdb-darwin-amd64
# - bin/containdb-windows-amd64.exe
```

### Debian Package
```bash
# Build .deb package
./Scripts/PackageBuilder.sh

# Output: Packages/containDB_<version>.deb
```

### npm Package
```bash
# Build and publish
cd npm/
npm run build
npm publish
```

## Anti-Patterns (FORBIDDEN)

❌ Platform-specific code without fallbacks
❌ Hardcoded credentials or defaults
❌ Direct docker command execution (use SDK)
❌ Breaking changes to CLI interface
❌ Missing error handling
❌ Incomplete rollback on errors
❌ Skipping cross-platform testing
❌ Missing documentation for new features

## Workflow Guidelines

### When Adding Database Support
1. Read `DESIGN.md` to understand architecture
2. Add database configuration to `src/base/Create.go`
3. Add management tool support (if applicable)
4. Update documentation (README, DESIGN)
5. Test on all platforms
6. Update version in `VERSION` file

### When Fixing Bugs
1. Write failing test that reproduces bug
2. Fix the bug
3. Verify test passes
4. Check for similar issues elsewhere
5. Document the fix in CHANGELOG.md

### When Refactoring
1. Ensure backward compatibility
2. Run all tests before and after
3. Test on all platforms
4. Update architecture diagrams if structure changes

## Success Criteria

Every task must meet ALL:
- ✅ Builds successfully (`go build`)
- ✅ Tests pass (`go test ./...`)
- ✅ Lints pass (`go vet ./...`, `go fmt ./...`)
- ✅ Works on Linux, macOS, Windows
- ✅ Documentation updated
- ✅ Backward compatible
- ✅ Auto-rollback works on errors
- ✅ Clear error messages
