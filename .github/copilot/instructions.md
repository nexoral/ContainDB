# GitHub Copilot Instructions for ContainDB

## Project Overview

**ContainDB** - Database Container Management CLI

- **Language**: Go ≥1.18
- **Type**: CLI tool for Docker database management
- **Platform**: Linux, macOS, Windows
- **Distribution**: npm package + native binaries

## Core Rules (NON-NEGOTIABLE)

### 1. Cross-Platform Compatibility
**ALWAYS test on Linux, macOS, Windows**
- Use `filepath` package for paths
- No platform-specific code without fallbacks
- Test on all target platforms

### 2. Docker SDK Only
```go
// ✅ Use official Docker SDK
import "github.com/docker/docker/client"

cli, err := client.NewClientWithOpts(client.FromEnv)

// ❌ NEVER execute docker commands
exec.Command("docker", "run", ...)
```

### 3. Auto-Rollback on Failure
```go
created := []string{}
defer func() {
    if err != nil {
        for _, id := range created {
            cleanup(id)
        }
    }
}()
```

### 4. Clear Error Messages
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

### 5. Interactive Prompts
```go
// ✅ ALWAYS prompt for sensitive data
password := promptSecure("Enter database password:")

// ❌ NEVER hardcode
password := "password123"
```

## Architecture

```
CLI Entry Point (main.go)
    ↓
Base Operations (src/base/)
    ↓
Docker Abstraction (src/Docker/)
    ↓
Docker Engine
```

## Commands

```bash
# Build & Test
go build -o containdb src/Core/main.go
go test ./...
go fmt ./...
go vet ./...

# Distribution
./Scripts/BinBuilder.sh       # Multi-platform build
./Scripts/PackageBuilder.sh   # Debian package
```

## Documentation

Update when features change:
- README.md - Usage, installation
- DESIGN.md - Architecture
- INSTALLATION.md - Platform guides
- CHANGELOG.md - Version changes

## Security

- Never log or display passwords
- Validate all user inputs
- Warn about root requirements
- No hardcoded secrets

## Testing

- Unit tests for logic
- Integration tests for CLI
- Cross-platform testing
- Error scenario testing

## Success Criteria

- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ Works on Linux, macOS, Windows
- ✅ Documentation updated
- ✅ Backward compatible
- ✅ Auto-rollback works
