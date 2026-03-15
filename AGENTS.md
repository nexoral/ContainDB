# AGENTS.md

OpenAI Codex CLI Instructions for ContainDB

## Project Overview

**ContainDB** - Database Container Management CLI

- **Language**: Go ≥1.18
- **Type**: Docker database management tool
- **Platform**: Cross-platform (Linux, macOS, Windows)
- **Distribution**: npm package + native binaries

## Build Commands

```bash
# Development
go build -o containdb src/Core/main.go
go run src/Core/main.go
go test ./...

# Production
./Scripts/BinBuilder.sh              # Multi-platform binaries
./Scripts/PackageBuilder.sh          # Debian package
npm run build && npm publish         # npm distribution
```

## Core Principles

### 1. Cross-Platform Compatibility
- Must work on Linux, macOS, Windows
- Test on all platforms
- No platform-specific code without fallbacks

### 2. User Experience
- Interactive prompts for sensitive data
- Smart defaults
- Clear error messages
- Auto-rollback on failures

### 3. Docker Integration
- Use official Docker SDK
- No direct docker commands
- Handle all Docker scenarios gracefully

## Architecture

### Layers
1. **CLI Entry Point** (`src/Core/main.go`)
2. **Base Operations** (`src/base/`) - Business logic
3. **Docker Abstraction** (`src/Docker/`) - SDK wrapper
4. **Docker Engine** - Container runtime

### Structure
```
src/
├── Core/main.go        # CLI entry
├── base/               # Logic
│   ├── Create.go       # DB creation
│   ├── Remove.go       # Cleanup
│   ├── Export.go       # docker-compose export
│   └── Import.go       # docker-compose import
├── Docker/             # SDK wrapper
│   ├── Container.go
│   ├── Volume.go
│   └── Network.go
└── tools/              # Management tools
```

## Go Standards

### Error Handling
```go
// ✅ GOOD - Descriptive with context
if err != nil {
    return fmt.Errorf("failed to create container '%s': %w", name, err)
}

// ❌ BAD - Generic
if err != nil {
    return err
}
```

### Type Safety
```go
// ✅ Struct for options
type ContainerOptions struct {
    Name          string
    Image         string
    Port          int
    RestartPolicy string
}

func CreateContainer(opts ContainerOptions) error

// ❌ Too many parameters
func CreateContainer(name, image string, port int, restart string) error
```

## Key Patterns

### Auto-Rollback
```go
func CreateDatabase(opts Options) error {
    created := []string{}

    defer func() {
        if err != nil {
            for _, resource := range created {
                cleanup(resource)
            }
        }
    }()

    container, err := CreateContainer(opts)
    if err != nil {
        return err
    }
    created = append(created, container.ID)
}
```

### Docker SDK Usage
```go
// ✅ Use official SDK
import "github.com/docker/docker/client"

cli, err := client.NewClientWithOpts(client.FromEnv)

// ❌ Don't execute commands
exec.Command("docker", "run", ...)
```

## Documentation

Update when features change:
- `README.md` - Usage, features, installation
- `DESIGN.md` - Architecture, technical decisions
- `INSTALLATION.md` - Platform-specific guides
- `CONTRIBUTING.md` - Development guidelines

## Security

- Never log or display passwords
- Validate all user inputs
- Warn users about root requirements
- Use Docker networks for isolation
- No hardcoded secrets

## Testing

- Unit tests: Core logic
- Integration tests: Full CLI workflows
- Cross-platform: Linux, macOS, Windows
- Error scenarios: Network failures, conflicts

## Supported Databases

- **SQL**: MySQL, PostgreSQL, MariaDB
- **NoSQL**: MongoDB, Redis
- **Tools**: phpMyAdmin, pgAdmin, RedisInsight, MongoDB Compass

## Anti-Patterns

❌ Platform-specific code without fallbacks
❌ Hardcoded credentials
❌ Direct docker commands (use SDK)
❌ Breaking CLI interface
❌ Missing error handling
❌ Incomplete rollback
❌ Skipping cross-platform tests

## Success Criteria

- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ Works on Linux, macOS, Windows
- ✅ Documentation updated
- ✅ Backward compatible
- ✅ Auto-rollback on errors
- ✅ Clear error messages
