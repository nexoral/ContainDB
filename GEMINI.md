# GEMINI.md

This file provides guidance to Gemini Code Assist when working with code in this repository.

## Project Overview

**ContainDB** - Database Container Management CLI Tool

- **Language**: Go ≥1.18
- **Type**: CLI for Docker database management
- **Platform**: Linux, macOS, Windows
- **Distribution**: npm + native binaries

## Commands

```bash
# Build & Run
go build -o containdb src/Core/main.go
go run src/Core/main.go
./Scripts/BinBuilder.sh              # Multi-platform build

# Test & Lint
go test ./...
go fmt ./...
go vet ./...
```

## Core Rules (NON-NEGOTIABLE)

1. **Cross-platform**: Must work on Linux, macOS, Windows
2. **ALWAYS test**: Test on all platforms before release
3. **ALWAYS build**: Run `go build` after changes
4. **Docker API only**: Use SDK, not direct commands
5. **Auto-rollback**: Cleanup on failures
6. **Update docs**: README, DESIGN, INSTALLATION

## Architecture

### Layers
```
CLI (main.go)
  ↓
Base Operations (src/base/)
  ↓
Docker Abstraction (src/Docker/)
  ↓
Docker Engine
```

### Structure
```
src/
├── Core/           # main.go
├── base/           # Business logic
│   ├── Create.go
│   ├── Remove.go
│   ├── Export.go
│   └── Import.go
├── Docker/         # Docker SDK wrapper
└── tools/          # Management tools
```

## Go Standards

### Error Handling
```go
// ✅ GOOD
if err != nil {
    return fmt.Errorf("failed to create container '%s': %w", name, err)
}

// ❌ BAD
if err != nil {
    return err
}
```

### Type Safety
```go
// ✅ Use structs
type ContainerOptions struct {
    Name  string
    Image string
    Port  int
}

// ❌ Too many params
func Create(name, image string, port int) error
```

## Key Patterns

### Auto-Rollback
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

### Docker SDK
```go
// ✅ Use Docker SDK
import "github.com/docker/docker/client"

cli, err := client.NewClientWithOpts(client.FromEnv)

// ❌ Don't execute commands
exec.Command("docker", "run", ...)
```

## Documentation

Update when features change:
- `README.md` - Usage, installation
- `DESIGN.md` - Architecture
- `INSTALLATION.md` - Platform guides
- Code comments for exported functions

## Security

- Never log passwords
- Validate all inputs
- Warn about root requirements
- Use Docker networks for isolation

## Testing

- Unit tests for core logic
- Integration tests for CLI
- Cross-platform testing
- Error case scenarios

## Supported Systems

- MySQL, PostgreSQL, MariaDB
- MongoDB, Redis
- phpMyAdmin, pgAdmin, RedisInsight

## Definition of "Done"

- ✅ `go build` passes
- ✅ `go test ./...` passes
- ✅ Works on Linux, macOS, Windows
- ✅ Documentation updated
- ✅ Backward compatible
- ✅ Auto-rollback works
