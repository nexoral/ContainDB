---
name: containdb-development
description: Development rules and patterns for ContainDB CLI tool
version: 1.0.0
tags: [golang, docker, cli, cross-platform, containers]
author: ContainDB Team
---

# ContainDB Development Skill

## Project Identity

**ContainDB** - Database Container Management Made Simple

- Go (Golang) ≥1.18
- Cross-platform CLI tool (Linux, macOS, Windows)
- Docker-based database management
- Distribution: npm package + native binaries

## Mandatory Workflows

### After EVERY Code Change
```bash
go build -o containdb src/Core/main.go
go test ./...
go fmt ./...
go vet ./...
```

### For ANY Feature Change
1. Test on Linux, macOS, Windows
2. Update documentation (README, DESIGN, INSTALLATION)
3. Update CHANGELOG.md
4. Verify auto-rollback works
5. Test error scenarios

## Definition of "Done"

A task is NOT complete until ALL are true:
- ✅ Code follows Go standards
- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ `go vet ./...` passes
- ✅ Works on Linux, macOS, Windows
- ✅ Documentation updated
- ✅ Backward compatible
- ✅ Auto-rollback implemented
- ✅ Error messages clear and actionable

## Architecture

### Layered Design
```
CLI Entry Point (src/Core/main.go)
    ↓
Base Operations (src/base/)
    ├── Create.go    - Database creation
    ├── Remove.go    - Resource cleanup
    ├── Export.go    - docker-compose export
    └── Import.go    - docker-compose import
    ↓
Docker Abstraction (src/Docker/)
    ├── Container.go - Container operations
    ├── Volume.go    - Volume management
    └── Network.go   - Network operations
    ↓
Docker Engine
```

### Module Structure
```
src/
├── Core/           # main.go - CLI entry point
├── base/           # Business logic
├── Docker/         # Docker SDK wrapper
└── tools/          # Management tools (phpMyAdmin, etc.)

Scripts/
├── BinBuilder.sh       # Multi-platform builder
├── PackageBuilder.sh   # Debian package builder
└── installer.sh        # Linux installer
```

## Go Standards (STRICT)

### Error Handling - Descriptive Context
```go
// ✅ REQUIRED
if err != nil {
    return fmt.Errorf("failed to create container '%s': %w", name, err)
}

// ❌ FORBIDDEN
if err != nil {
    return err  // Too generic!
}
```

### Type Safety - Use Structs
```go
// ✅ GOOD - Clear options struct
type ContainerOptions struct {
    Name          string
    Image         string
    Port          int
    RestartPolicy string
    Volumes       map[string]string
    Environment   map[string]string
}

func CreateContainer(opts ContainerOptions) error { }

// ❌ BAD - Too many parameters
func CreateContainer(name, image string, port int, restart string, vols map[string]string, env map[string]string) error { }
```

### Package Organization
```go
// ✅ GOOD - Clear package boundaries
package docker

import (
    "github.com/docker/docker/client"
)

func CreateContainer(opts ContainerOptions) error { }

// ❌ BAD - Mixed concerns in main
package main

func CreateContainer() error { }
func RemoveContainer() error { }
func ParseConfig() error { }
```

## Key Patterns

### 1. Auto-Rollback on Failure (CRITICAL)
```go
func CreateDatabase(opts Options) error {
    // Track all created resources
    created := []string{}

    // Cleanup on any error
    defer func() {
        if err != nil {
            log.Println("Error occurred, rolling back...")
            for _, resourceID := range created {
                if err := cleanup(resourceID); err != nil {
                    log.Printf("Failed to cleanup %s: %v", resourceID, err)
                }
            }
        }
    }()

    // Create container
    container, err := CreateContainer(opts)
    if err != nil {
        return fmt.Errorf("container creation failed: %w", err)
    }
    created = append(created, container.ID)

    // Create volume
    volume, err := CreateVolume(opts.VolumeName)
    if err != nil {
        return fmt.Errorf("volume creation failed: %w", err)
    }
    created = append(created, volume.Name)

    return nil
}
```

### 2. Interactive Prompts for Sensitive Data
```go
// ✅ ALWAYS prompt for passwords
import "github.com/manifoldco/promptui"

prompt := promptui.Prompt{
    Label: "Enter database password",
    Mask:  '*',
}
password, err := prompt.Run()

// ❌ NEVER hardcode or use defaults
password := "password123"  // FORBIDDEN
```

### 3. Docker SDK Usage (NOT Commands)
```go
// ✅ REQUIRED - Use official Docker SDK
import (
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
)

cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
    return fmt.Errorf("failed to create Docker client: %w", err)
}

resp, err := cli.ContainerCreate(ctx, &container.Config{
    Image: "mysql:8.0",
}, nil, nil, nil, "my-mysql")

// ❌ FORBIDDEN - Don't execute docker commands
cmd := exec.Command("docker", "run", "-d", "--name", "mysql", "mysql:8.0")
```

### 4. Cross-Platform Path Handling
```go
// ✅ GOOD - Use filepath package
import "path/filepath"

configPath := filepath.Join(homeDir, ".containdb", "config.yaml")

// ❌ BAD - Hardcoded separators
configPath := homeDir + "/.containdb/config.yaml"  // Fails on Windows
```

## Security Standards

### 1. Credential Handling
```go
// ✅ NEVER log passwords
log.Printf("Creating database with user: %s", username)  // OK
log.Printf("Creating database with password: %s", password)  // FORBIDDEN

// ✅ Prompt for sensitive data
password := promptSecure("Enter password:")

// ❌ No hardcoded secrets
const DefaultPassword = "admin123"  // FORBIDDEN
```

### 2. Input Validation
```go
// ✅ Validate all user inputs
func ValidatePort(port int) error {
    if port < 1024 || port > 65535 {
        return fmt.Errorf("port must be between 1024 and 65535")
    }
    return nil
}

func ValidateContainerName(name string) error {
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
        return fmt.Errorf("name must contain only alphanumeric, dash, or underscore")
    }
    return nil
}
```

### 3. Privilege Management
```go
// ✅ Warn users about root requirements
if os.Geteuid() != 0 {
    log.Println("WARNING: ContainDB requires root/sudo privileges for Docker operations")
    log.Println("Please run: sudo containdb")
    os.Exit(1)
}
```

## Documentation Requirements

### Update when features change:

1. **README.md**
   - Installation instructions
   - Usage examples
   - Feature list
   - Troubleshooting

2. **DESIGN.md**
   - Architecture diagrams
   - Technical decisions
   - Module descriptions

3. **INSTALLATION.md**
   - Platform-specific guides
   - Prerequisites
   - Verification steps

4. **CHANGELOG.md**
   - Version changes
   - New features
   - Bug fixes

5. **Code Comments**
   ```go
   // CreateDatabase creates a new containerized database instance.
   //
   // It performs the following steps:
   // 1. Creates a Docker network (if not exists)
   // 2. Creates a volume for data persistence
   // 3. Starts the database container
   // 4. Auto-rollback on any failure
   //
   // Parameters:
   //   - opts: Configuration options for the database
   //
   // Returns:
   //   - error: nil on success, descriptive error on failure
   func CreateDatabase(opts DatabaseOptions) error
   ```

## Testing Requirements

### Unit Tests
```go
func TestCreateContainer(t *testing.T) {
    tests := []struct {
        name    string
        opts    ContainerOptions
        wantErr bool
    }{
        {"valid MySQL", ContainerOptions{Name: "mysql", Image: "mysql:8.0"}, false},
        {"invalid name", ContainerOptions{Name: "my sql", Image: "mysql:8.0"}, true},
        {"missing image", ContainerOptions{Name: "mysql"}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := CreateContainer(tt.opts)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateContainer() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests
- Test full CLI workflows
- Test with real Docker containers
- Test auto-rollback scenarios
- Test on Linux, macOS, Windows

## Build & Distribution

### Multi-Platform Build
```bash
# Build for all platforms
./Scripts/BinBuilder.sh

# Output:
# bin/containdb-linux-amd64
# bin/containdb-darwin-amd64
# bin/containdb-darwin-arm64
# bin/containdb-windows-amd64.exe
```

### Debian Package
```bash
# Build .deb package
./Scripts/PackageBuilder.sh

# Output: Packages/containDB_<version>.deb
```

### npm Distribution
```bash
cd npm/
npm run build
npm publish
```

## Anti-Patterns (FORBIDDEN)

❌ Platform-specific code without fallbacks
❌ Hardcoded credentials or secrets
❌ Direct docker command execution (use SDK)
❌ Breaking changes to CLI interface
❌ Missing error handling
❌ Incomplete auto-rollback
❌ Skipping cross-platform testing
❌ Generic error messages
❌ Missing documentation updates
❌ Logging sensitive data

## Workflow Guidelines

### When Adding Database Support
1. Read `DESIGN.md` architecture section
2. Add database config to `src/base/Create.go`
3. Add default ports, images, environment vars
4. Implement auto-rollback for new resources
5. Add management tool support (if available)
6. Update README with usage examples
7. Test on all platforms
8. Update CHANGELOG.md

### When Fixing Bugs
1. Write failing test that reproduces bug
2. Fix the bug
3. Verify test passes
4. Test on all platforms
5. Check for similar issues elsewhere
6. Document in CHANGELOG.md

### When Refactoring
1. Ensure backward compatibility
2. Run all tests before and after
3. Test on all platforms
4. Update architecture docs if structure changes
5. Maintain CLI interface

## Success Criteria

Every task must meet ALL:
- ✅ Builds successfully (`go build`)
- ✅ Tests pass (`go test ./...`)
- ✅ Lints pass (`go vet ./...`, `go fmt ./...`)
- ✅ Works on Linux, macOS, Windows
- ✅ Documentation updated (README, DESIGN, INSTALLATION)
- ✅ Backward compatible (CLI interface unchanged)
- ✅ Auto-rollback implemented and tested
- ✅ Clear, actionable error messages
- ✅ No hardcoded credentials
- ✅ Input validation implemented
