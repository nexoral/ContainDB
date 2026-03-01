# ContainDB System Design Document

**Version:** 7.17.42-stable
**Last Updated:** 2026-03-01
**Document Purpose:** Technical architecture and design documentation for developers and contributors

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [System Overview](#system-overview)
   - [Product Definition](#product-definition)
   - [Key Features](#key-features)
   - [Supported Systems](#supported-systems)
3. [High-Level Architecture](#high-level-architecture)
   - [Layered Architecture](#layered-architecture)
   - [Module Dependencies](#module-dependencies)
   - [Data Flow](#data-flow)
4. [Core Components](#core-components)
   - [CLI Entry Point](#cli-entry-point-srccoremain go)
   - [Base Operations Package](#base-operations-package-srcbase)
   - [Docker Abstraction Package](#docker-abstraction-package-srcdocker)
   - [Tools Package](#tools-package-srctools)
5. [Container Management](#container-management)
   - [Container Lifecycle](#container-lifecycle)
   - [Supported Databases](#supported-databases)
   - [Management Tools](#management-tools)
6. [Data Persistence](#data-persistence)
   - [Volume Strategy](#volume-strategy)
   - [Volume Lifecycle](#volume-lifecycle)
7. [Network Architecture](#network-architecture)
   - [Network Design](#network-design)
   - [Network Topology](#network-topology)
   - [DNS Resolution](#dns-resolution)
8. [Configuration Management](#configuration-management)
   - [Docker Compose Export](#docker-compose-export)
   - [Docker Compose Import](#docker-compose-import)
9. [Error Handling & Recovery](#error-handling--recovery)
   - [Auto-Rollback Mechanism](#auto-rollback-mechanism)
   - [Signal Handling](#signal-handling)
10. [Security Model](#security-model)
    - [Privilege Management](#privilege-management)
    - [Credential Handling](#credential-handling)
    - [Container Isolation](#container-isolation)
11. [External Integrations](#external-integrations)
    - [Database Systems](#database-systems)
    - [Management Tools Integration](#management-tools-integration)
    - [Docker Engine Integration](#docker-engine-integration)
12. [Build & Distribution](#build--distribution)
    - [Build System](#build-system)
    - [Distribution Channels](#distribution-channels)
13. [Development Workflow](#development-workflow)
    - [Project Setup](#project-setup)
    - [Code Organization](#code-organization)
    - [Testing Strategy](#testing-strategy)
14. [Design Decisions & Trade-offs](#design-decisions--trade-offs)
15. [Operational Considerations](#operational-considerations)
16. [Appendices](#appendices)
    - [Code Metrics](#code-metrics)
    - [File Reference](#file-reference)
    - [Glossary](#glossary)

---

## Executive Summary

### What is ContainDB?

ContainDB is **not a database system** itself—it is a **command-line interface (CLI) tool** built in Go that simplifies the management of containerized database systems using Docker. Think of it as a productivity multiplier for developers who frequently work with database containers.

### Problem Statement

Developers face several friction points when setting up local database environments:

1. **Complex Docker Commands**: Running `docker run` with numerous flags for ports, volumes, environment variables, and networks is error-prone and hard to remember
2. **Platform-Specific Issues**: MongoDB experiences "Core Dumped" errors on certain Debian-based systems
3. **Fragmented Setup**: Each database and its management tools require separate configuration
4. **Manual Network Management**: Connecting containers requires manual network creation and configuration
5. **Configuration Drift**: No standardized way to export/import environment setups

### Core Value Proposition

ContainDB solves these problems by providing:

- **Zero-Configuration Setup**: Interactive menus guide users through database installation with sensible defaults
- **Unified Interface**: Single CLI for managing MongoDB, MySQL, PostgreSQL, MariaDB, Redis, and AxioDB
- **Automatic Orchestration**: Network creation, data persistence, and tool integration handled automatically
- **Portable Configurations**: Export/import via Docker Compose for reproducible environments
- **Auto-Rollback**: Automatic cleanup on errors prevents partial or broken states
- **Cross-Platform**: Runs on Linux, macOS, and Windows with consistent behavior

### Key Capabilities at a Glance

```
┌─────────────────────────────────────────────────────────────┐
│                    ContainDB Features                        │
├─────────────────────────────────────────────────────────────┤
│ ✓ Interactive database installation (6 databases)           │
│ ✓ Management tool integration (4 GUI tools)                 │
│ ✓ Docker Compose export/import                              │
│ ✓ Named volume management for data persistence              │
│ ✓ Automatic network isolation (ContainDB-Network)           │
│ ✓ Signal handling with auto-rollback on Ctrl+C              │
│ ✓ Cross-platform binaries (Linux, macOS, Windows)           │
│ ✓ NPM distribution with auto-update support                 │
└─────────────────────────────────────────────────────────────┘
```

---

## System Overview

### Product Definition

**Target Users:**
- Backend developers setting up local development environments
- DevOps engineers testing database configurations
- Students learning database systems
- Teams needing reproducible database setups

**Use Cases:**
- Local development with multiple database systems
- Testing database migrations and schema changes
- Prototyping applications with different databases
- Educational environments requiring quick database provisioning
- CI/CD pipelines needing ephemeral database containers

**Non-Goals:**
- Production database hosting (no HA, clustering, or replication management)
- Database performance tuning or optimization
- Database backup/restore orchestration
- Multi-host Docker deployments
- Database query execution or data manipulation

### Key Features

#### 1. Interactive Database Installation
- Menu-driven selection of database type
- Configurable port mapping (custom or default)
- Optional data persistence with named volumes
- Automatic image pulling from Docker Hub
- Environment variable configuration (credentials, etc.)

#### 2. Management Tool Integration
- **phpMyAdmin** for MySQL/MariaDB (containerized web UI)
- **pgAdmin** for PostgreSQL (containerized web UI)
- **RedisInsight** for Redis (containerized web UI)
- **MongoDB Compass** for MongoDB (native desktop application download)

#### 3. Docker Compose Export/Import
- Export running containers to `docker-compose.yml`
- Import and deploy services from existing Compose files
- Environment variable filtering (excludes system variables)
- Network configuration preservation

#### 4. Data Persistence Management
- Named volume creation with standard naming convention
- Volume reuse detection with user prompts
- Separate volume lifecycle from container lifecycle
- Volume usage validation before removal

#### 5. Network Isolation
- Dedicated Docker network: `ContainDB-Network`
- All containers attached to shared network
- DNS-based service discovery (container-to-container communication)
- Automatic network creation on startup

#### 6. Error Handling & Rollback
- SIGINT (Ctrl+C) handler for graceful interrupts
- Automatic cleanup of failed containers
- Dangling image removal
- Temporary file cleanup (e.g., MongoDB Compass downloads)

#### 7. Cross-Platform Support
- Linux (primary platform with sudo requirement)
- macOS (Intel and Apple Silicon)
- Windows (with optional Administrator privileges)
- Platform-specific binary selection via npm wrapper

### Supported Systems

#### Databases

| Database | Default Port | Docker Image | Data Path | Volume Pattern |
|----------|--------------|--------------|-----------|----------------|
| MongoDB | 27017 | `mongo:latest` | `/data/db` | `mongodb-data` |
| MySQL | 3306 | `mysql:latest` | `/var/lib/mysql` | `mysql-data` |
| PostgreSQL | 5432 | `postgres:latest` | `/var/lib/postgresql/data` | `postgresql-data` |
| MariaDB | 3306 | `mariadb:latest` | `/var/lib/mysql` | `mariadb-data` |
| Redis | 6379 | `redis:latest` | `/data` | `redis-data` |
| AxioDB | 27018 | `theankansaha/axiodb` | `/app/AxioDB` | `axiodb-data` |

#### Management Tools

| Tool | Target Database | Docker Image | Default Port | Protocol |
|------|-----------------|--------------|--------------|----------|
| phpMyAdmin | MySQL/MariaDB | `phpmyadmin/phpmyadmin` | 8080 | HTTP |
| pgAdmin | PostgreSQL | `dpage/pgadmin4:latest` | 5050 | HTTP |
| RedisInsight | Redis | `redis/redisinsight:latest` | 8001 (host) → 5540 (container) | HTTP |
| MongoDB Compass | MongoDB | N/A (native .deb download) | N/A | Native app |

---

## High-Level Architecture

### Layered Architecture

ContainDB follows a **5-layer architecture** that separates concerns and provides clear abstraction boundaries:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                       Layer 1: User Interface                           │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐     │
│  │   CLI Flags      │  │  Interactive     │  │   Banner &       │     │
│  │   (--version,    │  │  Prompts         │  │   Help Text      │     │
│  │    --export,     │  │  (promptui)      │  │                  │     │
│  │    --import)     │  │                  │  │                  │     │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘     │
│                                                                         │
│  Files: src/Core/main.go, src/base/flagHandler.go, src/base/Banner.go │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     Layer 2: Application Logic                          │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐     │
│  │  Menu System     │  │  Database        │  │  Container       │     │
│  │  (9 operations)  │  │  Selector        │  │  Orchestration   │     │
│  │                  │  │                  │  │                  │     │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘     │
│                                                                         │
│  Files: src/base/BaseCaseHandler.go, src/base/DatabaseSelector.go,    │
│         src/base/StartContainer.go, src/base/FilePathSelector.go      │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                  Layer 3: Docker Abstraction                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ │
│  │  Container   │  │   Volume     │  │   Network    │  │   Compose  │ │
│  │  Operations  │  │   Mgmt       │  │   Mgmt       │  │   Export/  │ │
│  │              │  │              │  │              │  │   Import   │ │
│  └──────────────┘  └──────────────┘  └──────────────┘  └────────────┘ │
│                                                                         │
│  Files: src/Docker/docker.go, src/Docker/docker_container.go,         │
│         src/Docker/Docker_Network.go, src/Docker/DockerComposeMaker.go,│
│         src/Docker/ImportDockerServices.go                             │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                Layer 4: System Integration                              │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐     │
│  │  Platform        │  │  Docker          │  │  System Req      │     │
│  │  Detection       │  │  Installation    │  │  Validation      │     │
│  │  (OS/Arch)       │  │  Check           │  │  (RAM/Disk)      │     │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘     │
│                                                                         │
│  Files: src/Docker/platform.go, src/Docker/docker_installation.go,    │
│         src/Docker/SysRequirement.go                                   │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Layer 5: Docker Engine                             │
│                                                                         │
│                   Docker CLI Commands (via os/exec)                    │
│         docker run | docker ps | docker network | docker volume       │
│                                                                         │
│                    External Dependency: Docker 20.10+                  │
└─────────────────────────────────────────────────────────────────────────┘
```

### Module Dependencies

ContainDB's Go packages have a clear hierarchical dependency structure:

```
src/Core/
  └── main.go
       ├─→ src/base/
       │    ├── BaseCaseHandler.go
       │    ├── DatabaseSelector.go
       │    ├── StartContainer.go
       │    ├── flagHandler.go
       │    ├── Banner.go
       │    ├── FilePathSelector.go
       │    └── DockerStarterPack.go
       │         └─→ src/Docker/
       │              ├── docker.go
       │              ├── docker_container.go
       │              ├── Docker_Network.go
       │              ├── DockerComposeMaker.go
       │              ├── ImportDockerServices.go
       │              ├── docker_installation.go
       │              ├── SysRequirement.go
       │              └── platform.go
       │
       └─→ src/tools/
            ├── PhpMyAdmin.go
            ├── PgAdmin.go
            ├── Redis_Insight.go
            ├── MongoDB_Tools.go
            ├── rollback.go
            ├── askForInput.go
            └── AfterContainerToolInstaller.go
                 └─→ src/Docker/ (shared dependency)

External Dependencies:
  ├─→ github.com/manifoldco/promptui v0.9.0 (interactive prompts)
  ├─→ github.com/fatih/color v1.18.0 (colored output)
  ├─→ gopkg.in/yaml.v2 v2.4.0 (YAML parsing)
  └─→ github.com/chzyer/readline v1.5.1 (terminal input)
```

**Dependency Rules:**
- `main.go` orchestrates initialization but delegates to `base` and `tools`
- `base` package handles user interaction and workflows
- `Docker` package provides stateless Docker API abstraction
- `tools` package implements specific database tool installers
- No circular dependencies between packages

### Data Flow

#### User Interaction Flow

```
   User Launch
       │
       ├─→ CLI Flags? ──→ (--version, --help, --export, --import)
       │                       │
       │                       └─→ Execute & Exit
       │
       ├─→ Check OS Support (src/Core/main.go:59-63)
       │
       ├─→ Check Privileges (Linux: sudo required) (src/Core/main.go:66-76)
       │
       ├─→ Check Docker Installation (src/base/DockerStarterPack.go)
       │
       ├─→ Create ContainDB-Network (src/Docker/Docker_Network.go:8-21)
       │
       ├─→ Show Banner (src/base/Banner.go)
       │
       └─→ Main Menu (src/base/BaseCaseHandler.go:16-25)
            │
            ├─→ Install Database
            │    ├─→ Select Database (src/base/DatabaseSelector.go)
            │    ├─→ Configure (port, persistence, credentials)
            │    └─→ Start Container (src/base/StartContainer.go)
            │         └─→ Docker run command execution
            │
            ├─→ List Databases ──→ docker ps --filter network=ContainDB-Network
            │
            ├─→ Remove Database ──→ Select → Confirm → docker rm
            │
            ├─→ Remove Image ──→ Check usage → Confirm → docker rmi
            │
            ├─→ Remove Volume ──→ Check usage → Confirm → docker volume rm
            │
            ├─→ Export Services ──→ Inspect containers → Generate YAML
            │
            ├─→ Import Services ──→ Parse YAML → Validate → docker compose up
            │
            ├─→ Update ContainDB ──→ npm update -g OR installer script
            │
            └─→ Exit
```

#### Docker Compose Export Flow

```
   User Selects "Export Services"
       │
       ├─→ List Containers on ContainDB-Network
       │    (docker ps --filter network=ContainDB-Network)
       │
       ├─→ For Each Container:
       │    │
       │    ├─→ docker inspect <container> --format {{.Config.Image}}
       │    ├─→ docker inspect <container> --format {{.NetworkSettings.Ports}}
       │    ├─→ docker inspect <container> --format {{.Mounts}}
       │    ├─→ docker inspect <container> --format {{.Config.Env}}
       │    └─→ Collect metadata into ContainerInfo struct
       │
       ├─→ Filter Environment Variables
       │    (Exclude: PATH, PHP_*, MONGO_*, system vars)
       │
       ├─→ Generate YAML Structure
       │    │
       │    ├─→ version: "3"
       │    ├─→ services: { ... }
       │    ├─→ volumes: { ... }
       │    └─→ networks: { ContainDB-Network: external: true }
       │
       └─→ Write to ./docker-compose.yml
            (src/Docker/DockerComposeMaker.go)
```

#### Docker Compose Import Flow

```
   User Selects "Import Services" + File Path
       │
       ├─→ Parse YAML File (gopkg.in/yaml.v2)
       │    └─→ DockerComposeConfig struct
       │
       ├─→ Pre-flight Validation:
       │    ├─→ Check port availability
       │    ├─→ Check for existing services with same names
       │    └─→ Check volume conflicts
       │
       ├─→ Create Missing Volumes
       │    (docker volume create <volume-name>)
       │
       ├─→ Deploy Services
       │    └─→ docker compose -f <file> up -d
       │
       └─→ Verify Deployment
            (src/Docker/ImportDockerServices.go)
```

#### Auto-Rollback Flow

```
   Trigger Event:
   ├─→ SIGINT (Ctrl+C) Signal
   └─→ Container Creation Error

       │
       ├─→ Signal Handler Catches Interrupt
       │    (src/Core/main.go:50-57)
       │
       └─→ tools.Cleanup() Execution
            (src/tools/rollback.go:13-47)
            │
            ├─→ Remove Failed Containers
            │    └─→ For status in [exited, dead, created]:
            │         docker ps -a --filter status=<status> --format {{.ID}}
            │         docker rm -f <container-id>
            │
            ├─→ Remove Dangling Images
            │    └─→ docker image prune -f
            │
            ├─→ Clean Temporary Files
            │    └─→ Remove /tmp/mongodb-compass.deb
            │
            └─→ Exit Program (exit code 1)
```

---

## Core Components

### CLI Entry Point (src/Core/main.go)

**Responsibilities:**
- Parse command-line flags and execute corresponding actions
- Verify OS compatibility and privilege requirements
- Register SIGINT handler for graceful interrupts
- Initialize Docker network infrastructure
- Delegate to base package for interactive operations

**Key Functions:**

```go
// main.go:13-95
func main() {
    VERSION := "7.17.42-stable"

    // Handle flags: --version, --help, --install-docker, --uninstall-docker
    // ... (lines 17-47)

    // Register SIGINT handler for Ctrl+C
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt)
    go func() {
        <-sigCh
        fmt.Println("\n⚠️ Interrupt received, rolling back...")
        tools.Cleanup()  // Auto-rollback
        os.Exit(1)
    }()

    // Verify OS support (lines 59-63)
    Docker.CheckOSSupport()

    // Check privileges (lines 66-76)
    // Linux: require sudo, Windows: warn if not admin

    // Ensure Docker is installed
    base.DockerStarter()

    // Process additional flags (export, import)
    base.FlagHandler()

    // Create network if not exists
    Docker.CreateDockerNetworkIfNotExists()

    // Show banner and start interactive menu
    base.ShowBanner()
    base.BaseCaseHandler()
}
```

**Signal Handling Implementation:**

The SIGINT handler (lines 50-57) ensures graceful cleanup on Ctrl+C interrupts. This is critical for preventing orphaned containers or partial deployments.

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/Core/main.go:50-57`

---

### Base Operations Package (src/base/)

This package provides the interactive user interface and orchestrates high-level workflows.

#### BaseCaseHandler.go - Main Menu System

**9-Operation Menu:**

```go
// BaseCaseHandler.go:16-19
actionPrompt := promptui.Select{
    Label: "What do you want to do?",
    Items: []string{
        "Install Database",     // 1. Database installation workflow
        "List Databases",       // 2. Show running databases
        "Remove Database",      // 3. Stop and remove container
        "Remove Image",         // 4. Remove Docker image
        "Remove Volume",        // 5. Delete persistent data
        "Import Services",      // 6. Deploy from docker-compose.yml
        "Export Services",      // 7. Generate docker-compose.yml
        "Update ContainDB",     // 8. Self-update mechanism
        "Exit"                  // 9. Exit program
    },
}
```

**Operation Implementations:**

1. **Install Database** (lines 28-41):
   - Calls `SelectDatabase()` for user choice
   - Routes to tool-specific installers (phpMyAdmin, pgAdmin, etc.) or `StartContainer()`

2. **List Databases** (lines 43-70):
   - Executes `Docker.ListRunningDatabases()`
   - Filters out management tools (phpmyadmin, pgadmin, redisinsight)
   - Displays container names

3. **Remove Database** (lines 72-96):
   - Lists running databases
   - Interactive selection with promptui
   - Calls `Docker.RemoveDatabase(name)`

4. **Remove Image** (lines 98-144):
   - Lists database images
   - Checks if image is in use via `Docker.IsImageInUse()`
   - Requires confirmation before removal

5. **Remove Volume** (lines 146-192):
   - Lists ContainDB-managed volumes
   - Checks usage via `Docker.IsVolumeInUse()`
   - Warning: "This will delete ALL DATA in this volume!"

6. **Export Services** (lines 193-206):
   - Warning about configuration-only export (no data)
   - Calls `Docker.MakeDockerComposeWithAllServices()`
   - Returns file path of generated `docker-compose.yml`

7. **Import Services** (lines 207-224):
   - Interactive file path selection
   - Calls `Docker.ImportDockerServices(filePath)`
   - Deploys services via `docker compose up -d`

8. **Update ContainDB** (lines 227-262):
   - Detects installation source (npm vs installer script)
   - npm: `npm update -g containdb`
   - Linux installer: Re-runs installer script via curl

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/base/BaseCaseHandler.go:14-269`

#### DatabaseSelector.go - Database Selection

Provides interactive menu for selecting database type or management tool.

**Supported Options:**
- MongoDB, MySQL, PostgreSQL, MariaDB, Redis, AxioDB
- phpMyAdmin, MongoDB Compass, PgAdmin, Redis Insight

#### StartContainer.go - Container Orchestration

**Responsibilities:**
- Pull Docker image if not present
- Configure port mapping (default or custom)
- Handle data persistence (create volume, reuse, or none)
- Set environment variables (database credentials)
- Execute `docker run` with all configurations
- Attach container to ContainDB-Network

**Configuration Flow:**
1. Determine default port and image name
2. Prompt for custom port or use default
3. Ask about data persistence
4. If persistent: check volume exists, prompt for reuse/recreate
5. Collect database-specific credentials (MySQL password, PostgreSQL user/password, etc.)
6. Build docker run command with all flags
7. Execute command and display result

---

### Docker Abstraction Package (src/Docker/)

This package provides a **stateless abstraction layer** over Docker CLI commands using `os/exec`.

#### docker.go - Core Docker Operations

**Key Functions:**

- `ListRunningDatabases() ([]string, error)`
  Executes: `docker ps --filter network=ContainDB-Network --format {{.Names}}`

- `RemoveDatabase(name string) error`
  Executes: `docker stop <name> && docker rm <name>`

- `ListDatabaseImages() ([]string, error)`
  Lists images for: mongo, mysql, postgres, redis, mariadb, phpmyadmin, pgadmin4

- `RemoveImage(image string) error`
  Executes: `docker rmi <image>`

- `IsImageInUse(image string) (bool, string, error)`
  Checks: `docker ps -a --filter ancestor=<image> --format {{.Names}}`

#### docker_container.go - Container & Volume Management

**Key Functions:**

- `VolumeExists(name string) bool`
  Executes: `docker volume inspect <name>`

- `CreateVolume(name string) error`
  Executes: `docker volume create <name>`

- `RemoveVolume(name string) error`
  Executes: `docker volume rm <name>`

- `ListContainDBVolumes() ([]string, error)`
  Lists volumes matching patterns: `mongodb-data`, `mysql-data`, `postgresql-data`, `redis-data`, `mariadb-data`, `axiodb-data`

- `IsVolumeInUse(volume string) (bool, string, error)`
  Executes: `docker ps -a --format "{{.Names}} {{.Mounts}}"` and parses output

- `AskYesNo(question string) bool`
  Interactive yes/no prompt using promptui

#### Docker_Network.go - Network Management

**Network Creation:**

```go
// Docker_Network.go:8-21
func CreateDockerNetworkIfNotExists() error {
    // Check if network exists
    cmdCheck := exec.Command("docker", "network", "inspect", "ContainDB-Network")
    err := cmdCheck.Run()
    if err == nil {
        return nil  // Network exists
    }

    // Create network
    cmdCreate := exec.Command("docker", "network", "create", "ContainDB-Network")
    cmdCreate.Stdout = os.Stdout
    cmdCreate.Stderr = os.Stderr
    return cmdCreate.Run()
}
```

**Network Characteristics:**
- Name: `ContainDB-Network` (hardcoded)
- Driver: Bridge (Docker default)
- Scope: Local to host machine
- Idempotent: Safe to call multiple times

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/Docker/Docker_Network.go:8-21`

#### DockerComposeMaker.go - Export Logic

**Data Structures:**

```go
type ContainerInfo struct {
    Name          string
    Image         string
    Ports         []string
    Volumes       []string
    EnvVars       []string
    Networks      []string
    Dependencies  []string
    RestartPolicy string
    Command       string
}

type DockerComposeService struct {
    Ports       []string
    Volumes     []string
    Environment map[string]string
}

type DockerComposeConfig struct {
    Version  string
    Services map[string]DockerComposeService
    Volumes  map[string]interface{}
}
```

**Export Process:**
1. List containers on ContainDB-Network
2. For each container, execute `docker inspect` to extract:
   - Image name
   - Port mappings
   - Volume mounts
   - Environment variables
   - Network configuration
3. Filter environment variables (exclude system vars like PATH, PHP_*, MONGO_*)
4. Build DockerComposeConfig struct
5. Marshal to YAML using `gopkg.in/yaml.v2`
6. Write to `./docker-compose.yml`

**Environment Variable Filtering:**
Excluded patterns: `PATH`, `PHP_*`, `MONGO_*`, `HOSTNAME`, system-level vars

#### ImportDockerServices.go - Import Logic

**Import Process:**
1. Read and parse YAML file
2. Validate port availability (check for conflicts)
3. Check for existing services with same names
4. Create missing volumes via `docker volume create`
5. Execute `docker compose -f <file> up -d`
6. Verify deployment success

**Pre-flight Checks:**
- Port conflict detection
- Service name collision detection
- Volume availability verification

---

### Tools Package (src/tools/)

Provides database management tool installers and cleanup utilities.

#### PhpMyAdmin.go

Installs phpMyAdmin container for MySQL/MariaDB management.

**Configuration:**
- Image: `phpmyadmin/phpmyadmin:latest`
- Default port: 8080 (configurable)
- Environment variables: `PMA_HOST`, `PMA_PORT`, `PMA_USER`
- Network: ContainDB-Network

**Workflow:**
1. Prompt for target MySQL/MariaDB container name
2. Prompt for port
3. Execute `docker run` with environment linking

#### PgAdmin.go

Installs pgAdmin container for PostgreSQL management.

**Configuration:**
- Image: `dpage/pgadmin4:latest`
- Default port: 5050
- Environment variables: `PGADMIN_DEFAULT_EMAIL`, `PGADMIN_DEFAULT_PASSWORD`
- Network: ContainDB-Network

**Workflow:**
1. Prompt for email and password
2. Start pgAdmin container
3. Display container IP for PostgreSQL connection configuration

#### Redis_Insight.go

Installs RedisInsight container for Redis management.

**Configuration:**
- Image: `redis/redisinsight:latest`
- Port mapping: 8001 (host) → 5540 (container)
- Network: ContainDB-Network

**Features:**
- Provides DNS-based connection instructions for Redis containers
- Web UI accessible at `http://localhost:8001`

#### MongoDB_Tools.go

Downloads and installs MongoDB Compass (native desktop application).

**Platform Support:** Linux only (Debian-based)

**Workflow:**
1. Download `.deb` package from mongodb.com
2. Install using `dpkg -i`
3. Cleanup downloaded package

**Limitations:** Containerized version not available; native app only

#### rollback.go - Cleanup Mechanism

**Cleanup Function:**

```go
// rollback.go:13-47
func Cleanup() {
    fmt.Println("🧹 Cleaning up resources...")

    // Remove failed containers (exited, dead, created states)
    statuses := []string{"exited", "dead", "created"}
    for _, status := range statuses {
        cmd := exec.Command("docker", "ps", "-a", "--filter",
                           fmt.Sprintf("status=%s", status),
                           "--format", "{{.ID}}")
        output, err := cmd.Output()
        if err == nil {
            containerIDs := strings.Fields(strings.TrimSpace(string(output)))
            for _, id := range containerIDs {
                rmCmd := exec.Command("docker", "rm", "-f", id)
                rmCmd.Run()
            }
        }
    }

    // Remove dangling images
    exec.Command("docker", "image", "prune", "-f").Run()

    // Remove MongoDB Compass download
    tempDir := Docker.GetTempDir()
    debPath := filepath.Join(tempDir, "mongodb-compass.deb")
    os.Remove(debPath)

    fmt.Println("✅ Cleanup completed.")
    os.Exit(1)
}
```

**Trigger Conditions:**
- SIGINT signal (Ctrl+C)
- Explicit call on operation failure
- Interrupt during interactive prompts

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/tools/rollback.go:13-47`

---

## Container Management

### Container Lifecycle

ContainDB containers follow a standard Docker lifecycle with specific ContainDB management:

```
   ┌──────────────────────────────────────────────────┐
   │                Container Lifecycle                │
   └──────────────────────────────────────────────────┘

    [User Selection]
           │
           ▼
    ┌──────────────┐
    │   PULL       │  docker pull <image>:latest
    │              │  - mongo, mysql, postgres, etc.
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │  CONFIGURE   │  - Port mapping (default or custom)
    │              │  - Volume mount (optional persistence)
    │              │  - Environment vars (credentials)
    │              │  - Network: ContainDB-Network
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │   CREATE     │  docker run -d --network ContainDB-Network
    │              │  --name <db>-container
    │              │  -p <host>:<container>
    │              │  -v <volume>:<data-path>
    │              │  -e <ENV_VARS>
    │              │  --restart unless-stopped
    │              │  <image>
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │   RUNNING    │  Container actively serving requests
    │              │  - DNS: <container-name>.ContainDB-Network
    │              │  - Accessible from host via port mapping
    │              │  - Accessible from other containers via DNS
    └──────┬───────┘
           │
           ├─→ [User: Stop/Remove]
           │        │
           │        ▼
           │   ┌──────────────┐
           │   │   STOPPED    │  docker stop <container>
           │   └──────┬───────┘
           │          │
           │          ▼
           │   ┌──────────────┐
           │   │   REMOVED    │  docker rm <container>
           │   │              │  (Volume persists if created)
           │   └──────────────┘
           │
           └─→ [Error/Interrupt]
                    │
                    ▼
               ┌──────────────┐
               │  ROLLBACK    │  tools.Cleanup()
               │              │  - Remove exited/dead/created
               │              │  - Prune dangling images
               └──────────────┘
```

### Supported Databases

#### MongoDB

**Image:** `mongo:latest`
**Default Port:** 27017
**Data Path:** `/data/db`
**Volume:** `mongodb-data`
**Authentication:** Optional (no auth by default)

**Container Name:** `mongodb-container`

**Environment Variables:** None required (can optionally set `MONGO_INITDB_ROOT_USERNAME` and `MONGO_INITDB_ROOT_PASSWORD`)

**Management Tool:** MongoDB Compass (native .deb download for Linux)

---

#### MySQL

**Image:** `mysql:latest`
**Default Port:** 3306
**Data Path:** `/var/lib/mysql`
**Volume:** `mysql-data`
**Authentication:** Required

**Container Name:** `mysql-container`

**Environment Variables:**
- `MYSQL_ROOT_PASSWORD` (required)

**Management Tool:** phpMyAdmin (containerized)

---

#### PostgreSQL

**Image:** `postgres:latest`
**Default Port:** 5432
**Data Path:** `/var/lib/postgresql/data`
**Volume:** `postgresql-data`
**Authentication:** Required

**Container Name:** `postgresql-container`

**Environment Variables:**
- `POSTGRES_USER` (required)
- `POSTGRES_PASSWORD` (required)

**Management Tool:** pgAdmin (containerized)

---

#### MariaDB

**Image:** `mariadb:latest`
**Default Port:** 3306
**Data Path:** `/var/lib/mysql`
**Volume:** `mariadb-data`
**Authentication:** Required

**Container Name:** `mariadb-container`

**Environment Variables:**
- `MARIADB_ROOT_PASSWORD` (required)

**Management Tool:** phpMyAdmin (containerized)

---

#### Redis

**Image:** `redis:latest`
**Default Port:** 6379
**Data Path:** `/data`
**Volume:** `redis-data`
**Authentication:** None (optional CONFIG password)

**Container Name:** `redis-container`

**Environment Variables:** None

**Management Tool:** RedisInsight (containerized)

---

#### AxioDB

**Image:** `theankansaha/axiodb:latest`
**Default Port:** 27018
**Data Path:** `/app/AxioDB`
**Volume:** `axiodb-data`
**Authentication:** Depends on image configuration

**Container Name:** `axiodb-container`

**Management Tool:** None (custom database)

---

### Management Tools

#### phpMyAdmin

**Purpose:** Web-based MySQL/MariaDB administration
**Image:** `phpmyadmin/phpmyadmin:latest`
**Default Port:** 8080
**Access:** `http://localhost:8080`

**Configuration:**
- `PMA_HOST`: Target MySQL/MariaDB container name
- `PMA_PORT`: Database port (default 3306)
- `PMA_USER`: Optional default username

**Network:** ContainDB-Network (communicates with database via DNS)

---

#### pgAdmin

**Purpose:** Web-based PostgreSQL administration
**Image:** `dpage/pgadmin4:latest`
**Default Port:** 5050
**Access:** `http://localhost:5050`

**Configuration:**
- `PGADMIN_DEFAULT_EMAIL`: Login email
- `PGADMIN_DEFAULT_PASSWORD`: Login password

**Connection Setup:** User must manually add PostgreSQL server using container IP or name

**Network:** ContainDB-Network

---

#### RedisInsight

**Purpose:** Web-based Redis administration
**Image:** `redis/redisinsight:latest`
**Port Mapping:** 8001 (host) → 5540 (container)
**Access:** `http://localhost:8001`

**Configuration:** None required (auto-discovery on network)

**Connection:** Use Redis container name (e.g., `redis-container`) and port 6379

**Network:** ContainDB-Network

---

#### MongoDB Compass

**Purpose:** Native desktop MongoDB administration
**Distribution:** `.deb` package download
**Platform:** Linux only (Debian/Ubuntu)

**Installation:**
1. Download from mongodb.com
2. Install via `dpkg -i mongodb-compass.deb`
3. Launch as native application

**Connection:** Use `localhost:27017` (or custom port)

**Network:** N/A (native app connects to exposed port)

---

## Data Persistence

### Volume Strategy

ContainDB uses **Docker named volumes** for data persistence with a standardized naming convention:

**Naming Pattern:** `{database-type}-data`

**Examples:**
- `mongodb-data`
- `mysql-data`
- `postgresql-data`
- `mariadb-data`
- `redis-data`
- `axiodb-data`

**Rationale:**
- **Predictable Naming**: Easy to identify which volume belongs to which database
- **Platform-Independent**: Named volumes work consistently across Linux, macOS, Windows
- **Docker-Managed**: Docker handles storage location and permissions
- **Portable**: Can be backed up via `docker volume inspect` or native database tools

**Trade-offs:**
- Less visibility into data location compared to bind mounts
- Requires Docker volume commands for direct access
- Not easily browsable from host filesystem

### Volume Lifecycle

```
   ┌──────────────────────────────────────────────────┐
   │              Volume Lifecycle                     │
   └──────────────────────────────────────────────────┘

    [User Enables Persistence]
           │
           ▼
    ┌──────────────┐
    │   CHECK      │  VolumeExists("{database}-data")
    │   EXISTS     │
    └──────┬───────┘
           │
           ├─→ [Exists] ──→ Prompt User:
           │                 ├─→ "Use existing" (mount to new container)
           │                 ├─→ "Create fresh" (delete + recreate)
           │                 └─→ "Exit" (cancel operation)
           │
           └─→ [Not Exists]
                    │
                    ▼
             ┌──────────────┐
             │   CREATE     │  docker volume create {database}-data
             └──────┬───────┘
                    │
                    ▼
             ┌──────────────┐
             │   ATTACH     │  docker run -v {volume}:{data-path}
             │              │
             │              │  Mount points:
             │              │  - MongoDB: /data/db
             │              │  - MySQL: /var/lib/mysql
             │              │  - PostgreSQL: /var/lib/postgresql/data
             │              │  - Redis: /data
             └──────┬───────┘
                    │
                    ▼
             ┌──────────────┐
             │   PERSIST    │  Data survives container removal
             │              │  Lifecycle independent of container
             └──────┬───────┘
                    │
                    ├─→ [Container Removed]
                    │        │
                    │        └─→ Volume remains intact
                    │
                    └─→ [User: Remove Volume]
                             │
                             ▼
                      ┌──────────────┐
                      │   CHECK      │  IsVolumeInUse()
                      │   USAGE      │
                      └──────┬───────┘
                             │
                             ├─→ [In Use] ──→ Error: "Volume in use by container X"
                             │
                             └─→ [Not In Use]
                                      │
                                      ▼
                               ┌──────────────┐
                               │   CONFIRM    │  "Delete ALL DATA?"
                               └──────┬───────┘
                                      │
                                      ├─→ [Yes] ──→ docker volume rm {volume}
                                      └─→ [No]  ──→ Cancel
```

**Volume Reuse Scenario:**

When a user installs a database with the same name as a previously removed container:

1. ContainDB detects existing volume
2. Prompts user with options:
   - **Use existing**: Mount existing volume (preserves data)
   - **Create fresh**: Delete old volume, create new one (fresh start)
   - **Exit**: Cancel operation
3. User choice determines volume handling

**Volume Usage Validation:**

Before allowing volume removal:
1. Execute `docker ps -a --format "{{.Names}} {{.Mounts}}"`
2. Parse output to check if volume is mounted
3. If in use, display error with container name
4. Require container removal before volume deletion

**Data Backup Recommendations:**

ContainDB does **not** provide automatic backup. Users should:
- Use native database tools (mongodump, mysqldump, pg_dump, redis-cli SAVE)
- Export via `docker volume` commands
- Use Docker Compose export as configuration backup (not data backup)

---

## Network Architecture

### Network Design

ContainDB creates a **dedicated Docker bridge network** for container isolation and communication:

**Network Name:** `ContainDB-Network`
**Driver:** Bridge (default Docker network driver)
**Scope:** Local (single host)
**Creation:** Idempotent check in `main.go` initialization

**Purpose:**
1. **Service Discovery**: Containers can reach each other using container names as hostnames
2. **Isolation**: Separates ContainDB containers from other Docker containers
3. **Simplification**: No manual IP management required

**Creation Logic:**

```go
// src/Docker/Docker_Network.go:8-21
func CreateDockerNetworkIfNotExists() error {
    cmdCheck := exec.Command("docker", "network", "inspect", "ContainDB-Network")
    err := cmdCheck.Run()
    if err == nil {
        return nil  // Network already exists
    }

    cmdCreate := exec.Command("docker", "network", "create", "ContainDB-Network")
    cmdCreate.Stdout = os.Stdout
    cmdCreate.Stderr = os.Stderr
    return cmdCreate.Run()
}
```

**Characteristics:**
- **Idempotent**: Safe to call multiple times (checks existence first)
- **Persistent**: Survives ContainDB restarts
- **Shared**: All ContainDB containers join this network
- **External**: Marked as `external: true` in exported docker-compose.yml

### Network Topology

```
┌────────────────────────────────────────────────────────────────────┐
│                         Host Machine                                │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │           ContainDB-Network (Bridge: 172.18.0.0/16)          │  │
│  │                                                               │  │
│  │   ┌─────────────────┐      ┌─────────────────┐              │  │
│  │   │ mongodb-        │      │ mysql-          │              │  │
│  │   │ container       │      │ container       │              │  │
│  │   │                 │      │                 │              │  │
│  │   │ IP: 172.18.0.2  │      │ IP: 172.18.0.3  │              │  │
│  │   │ DNS: mongodb-   │      │ DNS: mysql-     │              │  │
│  │   │      container  │      │      container  │              │  │
│  │   │                 │      │                 │              │  │
│  │   │ Port: 27017     │      │ Port: 3306      │              │  │
│  │   └────────┬────────┘      └────────┬────────┘              │  │
│  │            │                        │                        │  │
│  │            │ ◄──────DNS──────────►  │                        │  │
│  │            │   mongodb-container    │                        │  │
│  │            │   resolves to IP       │                        │  │
│  │            │                        │                        │  │
│  │   ┌────────▼────────┐      ┌───────▼─────────┐              │  │
│  │   │ pgadmin         │      │ phpmyadmin      │              │  │
│  │   │                 │      │                 │              │  │
│  │   │ IP: 172.18.0.4  │      │ IP: 172.18.0.5  │              │  │
│  │   │                 │      │                 │              │  │
│  │   │ Connects to:    │      │ Connects to:    │              │  │
│  │   │ postgresql-     │      │ mysql-container │              │  │
│  │   │ container:5432  │      │ (via PMA_HOST)  │              │  │
│  │   └─────────────────┘      └─────────────────┘              │  │
│  │                                                               │  │
│  └────────┬────────────────────────┬─────────────────────────────┘  │
│           │                        │                                │
│           │ Port Mapping           │ Port Mapping                   │
│           │ 27017:27017           │ 8080:80                         │
│           │ 5050:80               │                                 │
│           │                        │                                │
│  ┌────────▼────────────────────────▼─────────────────────────────┐  │
│  │                    Host Network Stack                         │  │
│  │         localhost:27017, localhost:8080, etc.                 │  │
│  └────────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  User Access:                                                 │  │
│  │  - Database: mongo://localhost:27017                          │  │
│  │  - phpMyAdmin: http://localhost:8080                          │  │
│  │  - pgAdmin: http://localhost:5050                             │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────┘
```

### DNS Resolution

Docker's embedded DNS server provides **automatic service discovery** within the ContainDB-Network:

**How It Works:**
1. Container A wants to connect to Container B (e.g., phpMyAdmin → mysql-container)
2. Container A queries `mysql-container` as hostname
3. Docker's DNS server (127.0.0.11) resolves to Container B's IP
4. Connection established via bridge network

**Example: phpMyAdmin Connecting to MySQL**

```bash
# phpMyAdmin container configuration
docker run -d \
  --name phpmyadmin \
  --network ContainDB-Network \
  -e PMA_HOST=mysql-container \  # Uses container name as hostname
  -e PMA_PORT=3306 \
  -p 8080:80 \
  phpmyadmin/phpmyadmin

# Inside phpMyAdmin container:
# mysql-container resolves to MySQL container's IP (e.g., 172.18.0.3)
# Connection succeeds via network
```

**Benefits:**
- No hardcoded IP addresses
- Containers can be recreated without configuration changes
- Simpler user experience (just use container names)

**Port Management:**

**Container-to-Container Communication:**
- Uses **container ports** (e.g., 3306, 27017, 5432)
- No host port mapping required for inter-container communication
- Traffic stays within Docker network (faster, isolated)

**Host-to-Container Communication:**
- Requires **port mapping** `-p <host>:<container>`
- Example: `-p 27017:27017` (host port 27017 → container port 27017)
- Host port can differ from container port: `-p 27018:27017`

**Port Conflict Handling:**
- ContainDB does **not** automatically detect port conflicts
- If port is in use, Docker will return an error
- User must choose a different port and retry

---

## Configuration Management

### Docker Compose Export

**Purpose:** Generate a portable `docker-compose.yml` file from running ContainDB containers.

**Use Cases:**
- Share database configurations with team members
- Migrate setups between development machines
- Document current environment state
- Recreate environment after system reset

**Export Workflow:**

```
User Selects "Export Services"
    │
    ├─→ List containers on ContainDB-Network
    │    └─→ docker ps --filter network=ContainDB-Network --format {{.Names}}
    │
    ├─→ For each container:
    │    │
    │    ├─→ docker inspect <container> --format {{.Config.Image}}
    │    ├─→ docker inspect <container> --format {{json .NetworkSettings.Ports}}
    │    ├─→ docker inspect <container> --format {{range .Mounts}}...{{end}}
    │    ├─→ docker inspect <container> --format {{.Config.Env}}
    │    ├─→ docker inspect <container> --format {{.HostConfig.RestartPolicy}}
    │    └─→ Collect into ContainerInfo struct
    │
    ├─→ Filter environment variables:
    │    ├─→ Exclude: PATH, HOSTNAME, HOME
    │    ├─→ Exclude: PHP_*, MONGO_*, MYSQL_* (system-generated)
    │    └─→ Include: User-defined credentials (MYSQL_ROOT_PASSWORD, etc.)
    │
    ├─→ Build YAML structure:
    │    │
    │    ├─→ version: "3"
    │    │
    │    ├─→ services:
    │    │    ├─→ mongodb-container:
    │    │    │    ├─→ image: mongo:latest
    │    │    │    ├─→ ports: ["27017:27017"]
    │    │    │    ├─→ volumes: ["mongodb-data:/data/db"]
    │    │    │    ├─→ environment: {...}
    │    │    │    ├─→ networks: [ContainDB-Network]
    │    │    │    └─→ restart: unless-stopped
    │    │    └─→ ... (other containers)
    │    │
    │    ├─→ volumes:
    │    │    ├─→ mongodb-data: {}
    │    │    └─→ mysql-data: {}
    │    │
    │    └─→ networks:
    │         └─→ ContainDB-Network:
    │              └─→ external: true
    │
    └─→ Write to ./docker-compose.yml
```

**Generated YAML Example:**

```yaml
version: "3"

services:
  mongodb-container:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongodb-data:/data/db
    networks:
      - ContainDB-Network
    restart: unless-stopped

  mysql-container:
    image: mysql:latest
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    environment:
      MYSQL_ROOT_PASSWORD: "password123"
    networks:
      - ContainDB-Network
    restart: unless-stopped

volumes:
  mongodb-data: {}
  mysql-data: {}

networks:
  ContainDB-Network:
    external: true
```

**Important Limitations:**

**Warning Displayed to User:**
```
⚠️  IMPORTANT: The export functionality only exports container configurations, not the actual data.
   Even if you used data persistence during installation, the exported compose file only
   references local volume paths from your current machine which won't exist on other systems.
   For data backup, please use each database's native backup tools.
```

**What's Exported:**
- Container configurations
- Port mappings
- Volume mount points (not volume data)
- Environment variables (filtered)
- Network configuration

**What's NOT Exported:**
- Actual database data
- Volume contents
- Host-specific paths

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/Docker/DockerComposeMaker.go`

### Docker Compose Import

**Purpose:** Deploy database services from an existing `docker-compose.yml` file.

**Use Cases:**
- Recreate environment from exported configuration
- Deploy team-shared database setups
- Migrate between development machines
- Initialize CI/CD test environments

**Import Workflow:**

```
User Selects "Import Services" + Provides File Path
    │
    ├─→ Read docker-compose.yml
    │    └─→ YAML parse using gopkg.in/yaml.v2
    │         └─→ DockerComposeConfig struct
    │
    ├─→ Pre-flight validation:
    │    │
    │    ├─→ Check port availability:
    │    │    └─→ For each service, verify host port is not in use
    │    │
    │    ├─→ Check for existing services:
    │    │    └─→ docker ps -a --format {{.Names}}
    │    │    └─→ Warn if container name conflicts
    │    │
    │    └─→ Check volume conflicts:
    │         └─→ docker volume ls --format {{.Name}}
    │         └─→ Prompt if volumes exist
    │
    ├─→ Create missing volumes:
    │    └─→ For each volume in compose file:
    │         └─→ docker volume create <volume-name>
    │
    ├─→ Deploy services:
    │    └─→ docker compose -f <file> up -d
    │         ├─→ Pull images
    │         ├─→ Create containers
    │         ├─→ Attach to networks
    │         └─→ Start containers
    │
    └─→ Verify deployment:
         └─→ docker ps --filter network=ContainDB-Network
```

**Pre-flight Checks:**

1. **Port Availability:**
   - Parses port mappings from YAML
   - Checks if host ports are already bound
   - Errors if conflicts detected

2. **Service Name Conflicts:**
   - Lists existing containers
   - Warns if container with same name exists
   - User must remove conflicting containers first

3. **Volume Conflicts:**
   - Lists existing volumes
   - Prompts user to reuse or recreate
   - Prevents accidental data loss

**Error Handling:**

- Invalid YAML: Parse error with line number
- Missing images: Docker pull failures displayed
- Port conflicts: Error message with conflicting port
- Permission errors: Escalation hints (sudo required on Linux)

**Network Handling:**

If compose file includes:
```yaml
networks:
  ContainDB-Network:
    external: true
```

Docker will attach services to existing ContainDB-Network. If not present, ContainDB recommends manual network specification.

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/Docker/ImportDockerServices.go`

---

## Error Handling & Recovery

### Auto-Rollback Mechanism

ContainDB implements a **comprehensive cleanup system** that automatically removes failed resources to prevent partial or broken states.

**Cleanup Function Implementation:**

```go
// src/tools/rollback.go:13-47
func Cleanup() {
    fmt.Println("🧹 Cleaning up resources...")

    // Step 1: Remove failed containers
    fmt.Println("- Removing failed containers...")
    statuses := []string{"exited", "dead", "created"}

    for _, status := range statuses {
        cmd := exec.Command("docker", "ps", "-a",
                           "--filter", fmt.Sprintf("status=%s", status),
                           "--format", "{{.ID}}")
        output, err := cmd.Output()
        if err == nil {
            containerIDs := strings.Fields(strings.TrimSpace(string(output)))
            for _, id := range containerIDs {
                if id != "" {
                    rmCmd := exec.Command("docker", "rm", "-f", id)
                    rmCmd.Run()
                }
            }
        }
    }

    // Step 2: Remove dangling images
    fmt.Println("- Removing dangling images...")
    exec.Command("docker", "image", "prune", "-f").Run()

    // Step 3: Clean temporary files
    tempDir := Docker.GetTempDir()
    debPath := filepath.Join(tempDir, "mongodb-compass.deb")
    os.Remove(debPath)

    fmt.Println("✅ Cleanup completed.")
    os.Exit(1)
}
```

**Trigger Conditions:**

1. **SIGINT Signal (Ctrl+C):**
   - User interrupt during operation
   - Caught by signal handler in main.go
   - Calls `tools.Cleanup()` before exit

2. **Interactive Prompt Cancellation:**
   - User presses Ctrl+C during promptui selection
   - promptui returns error
   - Cleanup called in error handler

3. **Container Creation Failure:**
   - Docker run command fails
   - Error detected and logged
   - Cleanup ensures no orphaned resources

**Cleanup Actions:**

**Container Cleanup:**
- Target states: `exited`, `dead`, `created`
- Rationale: These states indicate failed or incomplete operations
- Command: `docker rm -f <container-id>`
- No confirmation required (automatic)

**Image Cleanup:**
- Target: Dangling images (untagged, unreferenced)
- Command: `docker image prune -f`
- Frees disk space from failed pulls

**File Cleanup:**
- MongoDB Compass .deb file in temp directory
- Prevents accumulation of downloaded packages

**What's NOT Cleaned Up:**
- Running containers (state: running, restarting)
- Named volumes (preserves user data)
- Networks (ContainDB-Network persists)
- Successfully created images

**Exit Behavior:**
- Always exits with code 1 after cleanup
- Prevents continuation after interrupt
- Forces clean program termination

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/tools/rollback.go:13-47`

### Signal Handling

**SIGINT Handler Registration:**

```go
// src/Core/main.go:50-57
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt)

go func() {
    <-sigCh
    fmt.Println("\n⚠️ Interrupt received, rolling back...")
    tools.Cleanup()
    os.Exit(1)
}()
```

**How It Works:**

1. **Channel Creation:** Buffered channel of size 1 for signal reception
2. **Signal Registration:** `signal.Notify(sigCh, os.Interrupt)` registers SIGINT listener
3. **Goroutine Launch:** Background goroutine blocks on channel receive
4. **Interrupt Handling:**
   - User presses Ctrl+C
   - OS sends SIGINT to process
   - Signal delivered to channel
   - Goroutine unblocks
   - Prints warning message
   - Calls cleanup function
   - Exits with code 1

**Why Goroutine?**
- Non-blocking: Main program continues execution
- Always listening: Handler active throughout program lifetime
- Immediate response: Interrupt handled as soon as signal arrives

**Signal Coverage:**
- **SIGINT**: Covered (Ctrl+C)
- **SIGTERM**: Not covered (program runs interactively, not as daemon)
- **SIGKILL**: Cannot be caught (OS-level force kill)

**Reference:** `/home/ankan/Documents/Projects/ContainDB/src/Core/main.go:50-57`

---

## Security Model

### Privilege Management

ContainDB enforces **platform-specific privilege requirements** for Docker access:

**Linux:**
- **Requirement:** Root privileges (sudo)
- **Check:** `os.Geteuid() == 0`
- **Enforcement:** Hard requirement (program exits if not root)
- **Rationale:** Docker daemon requires root access on most Linux systems

**Code Implementation:**

```go
// src/Core/main.go:66-70
if runtime.GOOS == "linux" {
    if !Docker.IsAdmin() {
        fmt.Println("❌ Please run this program with sudo (Linux requires root for Docker)")
        os.Exit(1)
    }
}
```

**Windows:**
- **Requirement:** Administrator privileges (optional)
- **Check:** `net session` command execution
- **Enforcement:** Warning only (program continues)
- **Rationale:** Docker Desktop may work without admin rights

**Code Implementation:**

```go
// src/Core/main.go:71-76
if runtime.GOOS == "windows" {
    if !Docker.IsAdmin() {
        fmt.Println("⚠️  Warning: Not running as Administrator. Docker may require admin privileges.")
        fmt.Println("   Continuing anyway...")
    }
}
```

**macOS:**
- **Requirement:** Optional admin privileges
- **Check:** None (assumes Docker Desktop handles permissions)
- **Enforcement:** No check
- **Rationale:** Docker Desktop on macOS handles privilege escalation internally

**Admin Detection Function:**

```go
// src/Docker/platform.go
func IsAdmin() bool {
    if runtime.GOOS == "windows" {
        cmd := exec.Command("net", "session")
        err := cmd.Run()
        return err == nil
    } else {
        return os.Geteuid() == 0
    }
}
```

### Credential Handling

ContainDB follows secure credential management practices:

**Interactive Collection:**
- All database credentials collected via interactive prompts
- Uses `promptui.Prompt` for password input
- No command-line arguments for passwords (prevents shell history exposure)

**No Persistence:**
- Credentials **never** written to files by ContainDB
- Passed directly to Docker as environment variables
- No local storage, cache, or configuration files

**Environment Variable Passing:**

```bash
# Example: MySQL container creation
docker run -d \
  --name mysql-container \
  --network ContainDB-Network \
  -e MYSQL_ROOT_PASSWORD="<user-provided-password>" \
  -p 3306:3306 \
  mysql:latest
```

**Credential Lifecycle:**
1. User inputs password via prompt
2. Password stored in memory (string variable)
3. Passed to `docker run` command
4. Docker creates container with env var
5. Variable cleared when function returns

**Export Consideration:**

**Warning:** Credentials ARE included in exported docker-compose.yml:

```yaml
environment:
  MYSQL_ROOT_PASSWORD: "password123"  # ⚠️ Plaintext in file
```

Users should:
- Treat exported files as sensitive
- Use `.gitignore` to exclude from version control
- Consider using Docker secrets for production

**Security Recommendations:**
- Use strong passwords
- Rotate credentials regularly
- Don't commit compose files with credentials to Git
- Use environment variable files (`.env`) for team sharing

### Container Isolation

**Network Isolation:**

All ContainDB containers run on a **dedicated bridge network** (`ContainDB-Network`):

**Benefits:**
- Isolated from other Docker containers on the host
- Cannot communicate with containers on default bridge
- Reduced attack surface

**Limitations:**
- Containers within ContainDB-Network can communicate freely
- No firewall rules between ContainDB containers
- User must implement application-level authentication

**Volume Permissions:**

- Named volumes inherit Docker daemon permissions
- On Linux: owned by root (or container user if specified)
- Data accessible to any container mounting the volume
- No encryption at rest (depends on host filesystem)

**Host Network Mode:**

ContainDB **never** uses `--network host`:
- All containers use ContainDB-Network
- Port mapping required for host access
- Provides network namespace isolation

**Container User:**

- Containers run as default image user (varies by image)
- MongoDB: `mongodb` user
- MySQL: `mysql` user
- PostgreSQL: `postgres` user
- No explicit user override by ContainDB

**Security Hardening Opportunities:**

Future improvements could include:
- User namespace remapping
- Read-only root filesystems
- Capability dropping
- SELinux/AppArmor profiles
- Network segmentation (separate networks per database type)

---

## External Integrations

### Database Systems

ContainDB integrates with **6 database systems** via official Docker images:

#### MongoDB Integration

**Image:** `mongo:latest` (Official MongoDB image)
**Version:** Latest stable (pulled from Docker Hub)
**Documentation:** https://hub.docker.com/_/mongo

**Configuration:**
- Default port: 27017
- Data directory: `/data/db`
- Authentication: Disabled by default (can enable via env vars)
- Storage engine: WiredTiger

**Optional Environment Variables:**
- `MONGO_INITDB_ROOT_USERNAME`
- `MONGO_INITDB_ROOT_PASSWORD`
- `MONGO_INITDB_DATABASE`

**Management Tool:** MongoDB Compass (native app download)

---

#### MySQL Integration

**Image:** `mysql:latest` (Official MySQL image)
**Version:** Latest stable (8.x series)
**Documentation:** https://hub.docker.com/_/mysql

**Configuration:**
- Default port: 3306
- Data directory: `/var/lib/mysql`
- Authentication: Required (`MYSQL_ROOT_PASSWORD`)

**Required Environment Variables:**
- `MYSQL_ROOT_PASSWORD` (mandatory)

**Optional Environment Variables:**
- `MYSQL_DATABASE` (create database on init)
- `MYSQL_USER`, `MYSQL_PASSWORD` (create user)

**Management Tool:** phpMyAdmin (containerized)

---

#### PostgreSQL Integration

**Image:** `postgres:latest` (Official PostgreSQL image)
**Version:** Latest stable (16.x series)
**Documentation:** https://hub.docker.com/_/postgres

**Configuration:**
- Default port: 5432
- Data directory: `/var/lib/postgresql/data`
- Authentication: Required

**Required Environment Variables:**
- `POSTGRES_USER` (default: postgres)
- `POSTGRES_PASSWORD` (mandatory)

**Optional Environment Variables:**
- `POSTGRES_DB` (create database on init)

**Management Tool:** pgAdmin (containerized)

---

#### MariaDB Integration

**Image:** `mariadb:latest` (Official MariaDB image)
**Version:** Latest stable (11.x series)
**Documentation:** https://hub.docker.com/_/mariadb

**Configuration:**
- Default port: 3306
- Data directory: `/var/lib/mysql`
- Authentication: Required

**Required Environment Variables:**
- `MARIADB_ROOT_PASSWORD` (mandatory)

**Optional Environment Variables:**
- `MARIADB_DATABASE`
- `MARIADB_USER`, `MARIADB_PASSWORD`

**Management Tool:** phpMyAdmin (containerized, shared with MySQL)

---

#### Redis Integration

**Image:** `redis:latest` (Official Redis image)
**Version:** Latest stable (7.x series)
**Documentation:** https://hub.docker.com/_/redis

**Configuration:**
- Default port: 6379
- Data directory: `/data`
- Authentication: None by default (can enable via config)

**Persistence:**
- RDB snapshots (default)
- AOF available via config

**Management Tool:** RedisInsight (containerized)

---

#### AxioDB Integration

**Image:** `theankansaha/axiodb:latest` (Custom image)
**Maintainer:** Ankan Saha (ContainDB author)
**Documentation:** Custom database system

**Configuration:**
- Default port: 27018
- Data directory: `/app/AxioDB`
- Authentication: Depends on image configuration

**Management Tool:** None (custom CLI or API)

---

### Management Tools Integration

#### phpMyAdmin

**Image:** `phpmyadmin/phpmyadmin:latest`
**Purpose:** Web-based MySQL/MariaDB administration
**Documentation:** https://hub.docker.com/r/phpmyadmin/phpmyadmin

**Setup Process:**
1. User selects "Install Database" → "phpmyadmin"
2. Prompts for target MySQL/MariaDB container name
3. Prompts for port (default: 8080)
4. Creates container with environment linking

**Environment Configuration:**
- `PMA_HOST`: Target database container name (DNS resolution)
- `PMA_PORT`: Database port (default: 3306)
- `PMA_USER`: Optional default username

**Network:** ContainDB-Network (same as target database)

**Access:** `http://localhost:8080` (or custom port)

**Connection:** Automatic via DNS (uses container name as hostname)

---

#### pgAdmin

**Image:** `dpage/pgadmin4:latest`
**Purpose:** Web-based PostgreSQL administration
**Documentation:** https://hub.docker.com/r/dpage/pgadmin4

**Setup Process:**
1. User selects "Install Database" → "PgAdmin"
2. Prompts for email and password (pgAdmin credentials)
3. Prompts for port (default: 5050)
4. Creates container
5. Displays connection instructions with container IP

**Environment Configuration:**
- `PGADMIN_DEFAULT_EMAIL`: Login email
- `PGADMIN_DEFAULT_PASSWORD`: Login password

**Network:** ContainDB-Network

**Access:** `http://localhost:5050` (or custom port)

**Connection Setup:**
- User must manually add PostgreSQL server in pgAdmin UI
- Use container name (e.g., `postgresql-container`) and port 5432
- Or use container IP displayed during setup

---

#### RedisInsight

**Image:** `redis/redisinsight:latest`
**Purpose:** Web-based Redis administration and monitoring
**Documentation:** https://hub.docker.com/r/redis/redisinsight

**Setup Process:**
1. User selects "Install Database" → "Redis Insight"
2. Automatically configures port: 8001 (host) → 5540 (container)
3. Creates container
4. Displays connection instructions

**Port Mapping:**
- Fixed mapping: `8001:5540`
- Rationale: RedisInsight listens on port 5540 internally

**Network:** ContainDB-Network

**Access:** `http://localhost:8001`

**Connection Setup:**
- In RedisInsight UI, add database
- Use Redis container name (e.g., `redis-container`)
- Port: 6379
- DNS resolution automatic

---

#### MongoDB Compass

**Distribution:** Native desktop application (.deb package)
**Platform:** Linux only (Debian/Ubuntu)
**Documentation:** https://www.mongodb.com/try/download/compass

**Download Process:**
1. User selects "Install Database" → "MongoDB Compass"
2. ContainDB downloads `.deb` package from mongodb.com
3. Version: 1.46.2 (amd64)
4. URL: `https://downloads.mongodb.com/compass/mongodb-compass_1.46.2_amd64.deb`

**Installation:**
```bash
# Download to temp directory
tempDir := Docker.GetTempDir()
debPath := filepath.Join(tempDir, "mongodb-compass.deb")

# Install using dpkg
sudo dpkg -i mongodb-compass.deb
```

**Why Not Containerized?**
- MongoDB Compass is a desktop GUI application (Electron-based)
- Requires X11 or Wayland display server
- Containerizing desktop apps adds complexity
- Native installation provides better user experience

**Connection:**
- Launch MongoDB Compass application
- Connect to `localhost:27017` (or custom port)
- No network dependency (connects to exposed port)

**Cleanup:**
- Downloaded .deb file removed during `tools.Cleanup()`

---

### Docker Engine Integration

**Integration Method:** Docker CLI via `os/exec` package

**Why CLI Instead of Docker SDK?**

**Decision Rationale:**
- **Simplicity**: CLI commands are straightforward to execute
- **Version Independence**: No API version compatibility issues
- **Minimal Dependencies**: No additional Go SDK dependencies
- **Debugging**: Easier to debug (see exact commands executed)
- **Flexibility**: Can use latest Docker features immediately

**Trade-offs:**
- **Performance**: Process spawn overhead (~10-50ms per command)
- **Error Parsing**: String parsing of stderr instead of structured errors
- **Type Safety**: No compile-time validation of Docker operations

**Docker Version Requirements:**

**Minimum Version:** Docker 20.10.0
**Recommended Version:** Docker 24.0+

**Feature Dependencies:**
- `docker network create` (Docker 1.9+)
- `docker volume create` (Docker 1.9+)
- `docker compose` (Compose V2, Docker 20.10+)

**Compatibility Check:**

ContainDB does **not** verify Docker version. It assumes:
- Docker CLI is in PATH
- Docker daemon is running
- User has permissions to execute docker commands

**Command Execution Pattern:**

```go
// Example: Creating a network
cmdCreate := exec.Command("docker", "network", "create", "ContainDB-Network")
cmdCreate.Stdout = os.Stdout  // Forward stdout to user
cmdCreate.Stderr = os.Stderr  // Forward stderr to user
err := cmdCreate.Run()

if err != nil {
    // Handle error (error message already printed to stderr)
    return err
}
```

**Standard Commands Used:**

**Container Operations:**
- `docker run -d --network ... -p ... -v ... -e ... --name ... <image>`
- `docker ps --filter network=ContainDB-Network`
- `docker stop <container>`
- `docker rm <container>`
- `docker inspect <container> --format <template>`

**Volume Operations:**
- `docker volume create <volume>`
- `docker volume ls --format {{.Name}}`
- `docker volume rm <volume>`
- `docker volume inspect <volume>`

**Network Operations:**
- `docker network create <network>`
- `docker network inspect <network>`

**Image Operations:**
- `docker pull <image>`
- `docker rmi <image>`
- `docker images`
- `docker image prune -f`

**Compose Operations:**
- `docker compose -f <file> up -d`
- `docker compose -f <file> down`

**Output Handling:**
- **Stdout**: Displayed to user (container IDs, list results)
- **Stderr**: Displayed to user (error messages, warnings)
- **Exit Code**: Used for error detection (`err != nil`)

---

## Build & Distribution

### Build System

ContainDB uses a **multi-platform build system** that compiles Go binaries for various operating systems and architectures.

**Build Script:** `/home/ankan/Documents/Projects/ContainDB/Scripts/BinBuilder.sh`

**Target Platforms:**

| Platform | Architecture | Binary Name | Output Path |
|----------|--------------|-------------|-------------|
| Linux | amd64 (x86-64) | `containdb_linux_amd64` | `bin/` |
| macOS | amd64 (Intel) | `containdb_darwin_amd64` | `bin/` |
| macOS | arm64 (Apple Silicon) | `containdb_darwin_arm64` | `bin/` |
| Windows | amd64 (x86-64) | `containdb_windows_amd64.exe` | `bin/` |

**Build Process:**

1. **Version Management:**
   - Read version from `VERSION` file
   - Update `main.go`, `package.json`, documentation

2. **Cross-Compilation:**
   ```bash
   # Linux
   GOOS=linux GOARCH=amd64 go build -o bin/containdb_linux_amd64 src/Core/main.go

   # macOS Intel
   GOOS=darwin GOARCH=amd64 go build -o bin/containdb_darwin_amd64 src/Core/main.go

   # macOS Apple Silicon
   GOOS=darwin GOARCH=arm64 go build -o bin/containdb_darwin_arm64 src/Core/main.go

   # Windows
   GOOS=windows GOARCH=amd64 go build -o bin/containdb_windows_amd64.exe src/Core/main.go
   ```

3. **NPM Package Preparation:**
   - Copy binaries to `npm/bin/` directory
   - Update `npm/package.json` with version
   - Prepare for npm publish

4. **Debian Package Build:**
   - Execute `Scripts/PackageBuilder.sh`
   - Create `.deb` package for Linux distribution
   - Output to `Packages/` directory

**Build Dependencies:**
- Go 1.24+ (specified in `go.mod:3-5`)
- Bash shell (for build scripts)
- `dpkg-deb` (for Debian package creation)

**Go Compilation Flags:**
- No additional flags (default compilation)
- Static linking handled by Go default for standalone binaries

**Binary Characteristics:**
- **Size:** ~15-20 MB per binary (depends on Go version and platform)
- **Dependencies:** None (statically linked, except libc on Linux)
- **Executable:** Directly runnable (no interpreter required)

---

### Distribution Channels

ContainDB provides **4 distribution methods** to support different platforms and workflows:

#### 1. NPM Package Distribution

**Package Name:** `containdb`
**Registry:** https://www.npmjs.com/package/containdb
**Installation:** `npm install -g containdb`

**Platform Support:**
- Linux (x64)
- macOS (Intel x64, Apple Silicon arm64)
- Windows (x64)

**How It Works:**

```javascript
// npm/InstallController.js
const os = require('os');
const { spawn } = require('child_process');

const platform = os.platform();  // 'linux', 'darwin', 'win32'
const arch = os.arch();          // 'x64', 'arm64'

// Select appropriate binary
let binaryPath;
if (platform === 'linux' && arch === 'x64') {
    binaryPath = './bin/containdb_linux_amd64';
} else if (platform === 'darwin' && arch === 'x64') {
    binaryPath = './bin/containdb_darwin_amd64';
} else if (platform === 'darwin' && arch === 'arm64') {
    binaryPath = './bin/containdb_darwin_arm64';
} else if (platform === 'win32' && arch === 'x64') {
    binaryPath = './bin/containdb_windows_amd64.exe';
} else {
    console.error('Unsupported platform:', platform, arch);
    process.exit(1);
}

// Spawn binary process
const child = spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit' });
```

**Environment Variable:**
- Sets `CONTAINDB_INSTALL_SOURCE=npm` during execution
- Used by update mechanism to detect npm installation

**Benefits:**
- Cross-platform (single install command for all OS)
- Automatic updates via `npm update -g containdb`
- Familiar workflow for Node.js developers

**Limitations:**
- Requires Node.js 16+ and npm 8+
- Larger package size (~60 MB with all binaries)

---

#### 2. Linux Installer Script

**URL:** https://raw.githubusercontent.com/nexoral/ContainDB/main/Scripts/installer.sh
**Installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/nexoral/ContainDB/main/Scripts/installer.sh | sudo bash -
```

**Process:**
1. Detect system architecture (amd64, arm64, i386)
2. Download appropriate binary from GitHub releases
3. Install to `/usr/local/bin/containdb`
4. Set executable permissions

**Platform Support:**
- Ubuntu 20.04+
- Debian 11+
- Other Linux distributions with `curl` and `bash`

**Benefits:**
- Single command installation
- No npm dependency
- System-wide installation

**Security Consideration:**
- Piping curl to bash is convenient but risky
- Users should review script before execution
- Alternative: download script, inspect, then execute

---

#### 3. Debian Package (.deb)

**Build Script:** `/home/ankan/Documents/Projects/ContainDB/Scripts/PackageBuilder.sh`
**Output:** `Packages/containdb_<version>_<arch>.deb`

**Supported Architectures:**
- amd64 (x86-64)
- arm64 (ARM 64-bit)
- i386 (x86 32-bit, legacy)

**Installation:**
```bash
sudo dpkg -i containdb_7.17.42_amd64.deb
```

**Package Contents:**
- Binary: `/usr/bin/containdb`
- Documentation: `/usr/share/doc/containdb/`
- Man page: `/usr/share/man/man1/containdb.1.gz` (if created)

**Dependencies:**
- Docker (not enforced by package, user must install separately)

**Benefits:**
- Native package manager integration
- Automatic PATH configuration
- Clean uninstallation via `apt remove containdb`

---

#### 4. Direct Binary Download

**Source:** GitHub Releases
**URL Pattern:** `https://github.com/nexoral/ContainDB/releases/download/v<version>/containdb_<platform>_<arch>`

**Download Examples:**
```bash
# Linux
wget https://github.com/nexoral/ContainDB/releases/download/v7.17.42/containdb_linux_amd64

# macOS Intel
curl -LO https://github.com/nexoral/ContainDB/releases/download/v7.17.42/containdb_darwin_amd64

# macOS Apple Silicon
curl -LO https://github.com/nexoral/ContainDB/releases/download/v7.17.42/containdb_darwin_arm64

# Windows
Invoke-WebRequest -Uri https://github.com/nexoral/ContainDB/releases/download/v7.17.42/containdb_windows_amd64.exe -OutFile containdb.exe
```

**Manual Installation:**
```bash
chmod +x containdb_linux_amd64
sudo mv containdb_linux_amd64 /usr/local/bin/containdb
```

**Benefits:**
- No dependencies (npm, package managers)
- Full control over installation location
- Suitable for CI/CD systems

**Limitations:**
- Manual PATH configuration
- Manual updates (no auto-update)

---

**Distribution Matrix:**

| Method | Platforms | Auto-Update | Dependencies | Installation Complexity |
|--------|-----------|-------------|--------------|------------------------|
| NPM | Linux, macOS, Windows | ✅ Yes | Node.js 16+ | Low |
| Installer Script | Linux | ⚠️ Re-run script | curl, bash | Low |
| Debian Package | Debian/Ubuntu | ❌ No | dpkg | Medium |
| Direct Binary | All | ❌ No | None | High |

---

## Development Workflow

### Project Setup

**Prerequisites:**
- Go 1.24+ ([installation guide](https://go.dev/doc/install))
- Docker 20.10+ ([installation guide](https://docs.docker.com/engine/install/))
- Git
- Bash (for build scripts)

**Clone Repository:**

```bash
git clone https://github.com/nexoral/ContainDB.git
cd ContainDB
```

**Install Dependencies:**

```bash
go mod download
# Downloads:
# - github.com/manifoldco/promptui v0.9.0
# - github.com/fatih/color v1.18.0
# - gopkg.in/yaml.v2 v2.4.0
# - github.com/chzyer/readline v1.5.1
```

**Build Binary:**

```bash
# Single platform (current OS)
go build -o bin/containdb src/Core/main.go

# Multi-platform (all targets)
./Scripts/BinBuilder.sh
```

**Run Locally:**

```bash
# Linux/macOS
sudo ./bin/containdb

# Windows (PowerShell as Administrator)
.\bin\containdb.exe
```

**Run from Source (Development):**

```bash
# Direct execution (hot reload)
sudo go run src/Core/main.go
```

---

### Code Organization

```
ContainDB/
├── src/                          # Go source code (production code)
│   ├── Core/                     # Entry point package
│   │   └── main.go              # main() function, initialization
│   │
│   ├── base/                     # User interface and workflows
│   │   ├── BaseCaseHandler.go   # Main menu system (9 operations)
│   │   ├── DatabaseSelector.go  # Database selection UI
│   │   ├── StartContainer.go    # Container orchestration logic
│   │   ├── flagHandler.go       # CLI flag processing (--export, --import)
│   │   ├── Banner.go            # Welcome banner display
│   │   ├── FilePathSelector.go  # Interactive file path input
│   │   └── DockerStarterPack.go # Docker installation check/bootstrap
│   │
│   ├── Docker/                   # Docker abstraction layer
│   │   ├── docker.go            # Container operations (list, remove)
│   │   ├── docker_container.go  # Volume operations, utility functions
│   │   ├── Docker_Network.go    # Network creation (ContainDB-Network)
│   │   ├── DockerComposeMaker.go # Export to docker-compose.yml
│   │   ├── ImportDockerServices.go # Import from docker-compose.yml
│   │   ├── docker_installation.go # Docker installation automation
│   │   ├── SysRequirement.go    # System validation (RAM, disk, OS)
│   │   └── platform.go          # OS/architecture detection
│   │
│   └── tools/                    # Management tools and utilities
│       ├── PhpMyAdmin.go        # phpMyAdmin container setup
│       ├── PgAdmin.go           # pgAdmin container setup
│       ├── Redis_Insight.go     # RedisInsight container setup
│       ├── MongoDB_Tools.go     # MongoDB Compass download/install
│       ├── rollback.go          # Cleanup function (auto-rollback)
│       ├── askForInput.go       # User input helpers
│       └── AfterContainerToolInstaller.go # Post-install tool suggestions
│
├── Scripts/                      # Build and deployment automation
│   ├── BinBuilder.sh            # Multi-platform Go compilation
│   ├── PackageBuilder.sh        # Debian package builder
│   ├── installer.sh             # Linux installer (curl | bash)
│   ├── versionController.sh     # Version management across files
│   └── release.sh               # GitHub release automation
│
├── npm/                          # NPM package structure
│   ├── package.json             # NPM manifest (version, keywords, metadata)
│   ├── InstallController.js     # Platform detection and binary launcher
│   ├── bin/                     # Platform-specific binaries (copied by build)
│   │   ├── containdb_linux_amd64
│   │   ├── containdb_darwin_amd64
│   │   ├── containdb_darwin_arm64
│   │   └── containdb_windows_amd64.exe
│   ├── LICENSE                  # MIT license
│   └── README.md                # NPM package documentation
│
├── Packages/                     # Debian package output (generated)
│   └── containdb_7.17.42_amd64.deb
│
├── .github/                      # GitHub repository configuration
│   └── workflows/
│       ├── run_build.yml        # CI: Build on push
│       └── ReviewBuddy.yml      # Automated code review
│
├── Documentation/                # Project documentation (markdown)
│   ├── README.md                # Main user documentation
│   ├── DEPLOYMENT.md            # Deployment guide
│   ├── INSTALLATION.md          # Installation instructions
│   ├── LEARN.md                 # Learning resources
│   ├── CONTRIBUTING.md          # Contributor guidelines
│   ├── ROADMAP.md               # Future feature plans
│   ├── SECURITY.md              # Security policy
│   ├── CODE_OF_CONDUCT.md       # Community guidelines
│   └── CHANGELOG.md             # Version history
│
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── VERSION                       # Current version number
├── LICENSE                       # MIT License
├── .gitignore                    # Git exclusions
└── README.md                     # Repository README (symlink to Documentation/README.md)
```

**Package Responsibilities:**

- **Core:** Entry point, initialization, signal handling
- **base:** User interaction, menu system, workflows
- **Docker:** Stateless Docker CLI abstraction
- **tools:** Tool-specific installers, cleanup utilities

**Import Rules:**
- `main.go` imports `base`, `Docker`, `tools`
- `base` imports `Docker`, `tools`
- `Docker` has no internal dependencies (only stdlib + external libs)
- `tools` imports `Docker`

---

### Testing Strategy

**Current Approach:** **Manual Testing**

ContainDB does **not** currently have automated tests. Testing is performed manually through:

**Test Environments:**
- **Ubuntu 22.04 LTS** (primary platform)
- **Debian 11** (Bullseye)
- **macOS 13** (Ventura, Intel)
- **macOS 14** (Sonoma, Apple Silicon)
- **Windows 11** (via Docker Desktop)

**Test Scenarios:**

1. **Database Installation:**
   - Install each database type (MongoDB, MySQL, PostgreSQL, MariaDB, Redis, AxioDB)
   - Verify container creation
   - Test port mapping (default and custom)
   - Validate data persistence (create data, remove container, recreate, verify data)

2. **Management Tools:**
   - Install phpMyAdmin, verify MySQL/MariaDB connection
   - Install pgAdmin, verify PostgreSQL connection
   - Install RedisInsight, verify Redis connection
   - Download MongoDB Compass (Linux only)

3. **Docker Compose:**
   - Export running containers to YAML
   - Verify YAML structure and content
   - Import YAML on clean system
   - Verify containers recreated correctly

4. **Error Handling:**
   - Test Ctrl+C interrupt during operations
   - Verify cleanup removes failed containers
   - Test port conflicts
   - Test volume conflicts

5. **Update Mechanism:**
   - Test npm update path
   - Test installer script update path

**Docker Version Matrix:**
- Docker 20.10.x
- Docker 24.0.x (latest)

**Manual Testing Checklist:**

```markdown
## Pre-Release Testing Checklist

- [ ] Build successful for all platforms (Linux, macOS, Windows)
- [ ] Install MongoDB with persistence
- [ ] Install MySQL with credentials
- [ ] Install PostgreSQL with custom port
- [ ] Install Redis without persistence
- [ ] Install phpMyAdmin and connect to MySQL
- [ ] Install pgAdmin and connect to PostgreSQL
- [ ] Export docker-compose.yml from running containers
- [ ] Clean environment (remove all containers)
- [ ] Import docker-compose.yml
- [ ] Verify all containers recreated
- [ ] Test Ctrl+C interrupt (verify cleanup)
- [ ] Test volume removal (verify in-use check)
- [ ] Test image removal (verify in-use check)
- [ ] Test update mechanism (npm and installer script)
```

**Future Testing Improvements:**

Automated testing opportunities:
1. **Unit Tests:**
   - Docker command construction
   - YAML parsing/generation
   - Environment variable filtering
   - Volume naming logic

2. **Integration Tests:**
   - End-to-end database installation (requires Docker)
   - Export/import round-trip
   - Cleanup verification

3. **Snapshot Tests:**
   - Generated docker-compose.yml output
   - Help text and banner output

**Testing Challenges:**
- Requires Docker daemon (integration test complexity)
- Interactive prompts (need to mock `promptui`)
- Platform-specific behavior (Linux vs macOS vs Windows)

---

## Design Decisions & Trade-offs

### Decision 1: Docker CLI vs Docker SDK

**Decision:** Use Docker CLI via `os/exec` instead of official Docker SDK

**Rationale:**
- **Simplicity**: Executing CLI commands is straightforward (`docker run`, `docker ps`, etc.)
- **No API Version Lock-in**: CLI is stable across Docker versions
- **Minimal Dependencies**: No need for large SDK libraries
- **Debugging**: Easy to see exact commands executed (logged to stdout/stderr)
- **Compatibility**: Works with any Docker version that has CLI

**Trade-offs:**
- **Performance**: Process spawn overhead (~10-50ms per command)
  - Impact: Negligible for interactive tool
  - Not suitable for high-frequency operations
- **Error Handling**: Parse stderr strings instead of structured errors
  - Impact: Harder to distinguish error types
  - Mitigated by displaying stderr directly to user
- **Type Safety**: No compile-time validation of Docker operations
  - Impact: Typos in command strings caught at runtime
  - Mitigated by extensive manual testing

**Alternative Considered:**
Using `github.com/docker/docker/client` (official Go SDK):
- Pros: Structured API, type safety, better performance
- Cons: Dependency bloat, API version management, complexity

**Verdict:** CLI approach is appropriate for ContainDB's interactive, user-facing nature.

---

### Decision 2: Interactive vs Declarative

**Decision:** Menu-driven interactive interface (promptui) instead of declarative configuration files

**Rationale:**
- **User Experience**: Guides users through options (self-documenting)
- **Discovery**: Users see available operations without reading docs
- **Safety**: Confirmation prompts prevent accidental deletions
- **Simplicity**: No need to learn YAML/JSON schema
- **Progressive Disclosure**: Shows options based on context

**Trade-offs:**
- **Automation**: Not scriptable (can't pass config file to automate setup)
  - Impact: Not suitable for CI/CD or automated provisioning
  - Mitigated by: Docker Compose export/import for reproducible setups
- **Speed**: Slower than command-line flags for power users
  - Impact: More clicks/selections required
  - Mitigated by: Flag-based commands for export/import
- **Testability**: Hard to automate testing of interactive prompts
  - Impact: Requires manual testing
  - No easy mitigation

**Alternative Considered:**
Config file approach (e.g., `containdb.yml`):
```yaml
databases:
  - type: mongodb
    port: 27017
    persistence: true
  - type: mysql
    port: 3306
    credentials:
      root_password: password123
```

- Pros: Scriptable, versionable, faster for power users
- Cons: Requires documentation, schema validation, less discoverable

**Verdict:** Interactive approach fits target audience (developers setting up local environments). Export/import addresses automation needs.

---

### Decision 3: Single Network vs Multi-Network

**Decision:** All containers on single shared network (`ContainDB-Network`)

**Rationale:**
- **Simplicity**: Single network easier to understand and manage
- **DNS-based Discovery**: Containers can reference each other by name
- **Management Tools**: phpMyAdmin, pgAdmin can easily connect to databases
- **Reduced Complexity**: No network routing or configuration needed

**Trade-offs:**
- **No Isolation**: All ContainDB databases can communicate with each other
  - Impact: Security concern if multiple projects/teams on same host
  - Mitigation: Application-level authentication still required
- **Namespace Collision**: Can't have two containers with same name
  - Impact: Can't run two MySQL instances for different projects
  - Mitigation: Users must use different container names
- **No Multi-Project Support**: All databases mixed together
  - Impact: Hard to separate dev/test/staging environments
  - Future: Could implement project-based network prefixes

**Alternative Considered:**
Per-database-type networks:
- `ContainDB-MongoDB-Network`
- `ContainDB-MySQL-Network`
- Pros: Better isolation between database types
- Cons: Management tools need multiple network attachments, complexity

**Verdict:** Single network is appropriate for typical single-developer local development use case. Future enhancement could add project-based namespaces.

---

### Decision 4: Named Volumes vs Bind Mounts

**Decision:** Use Docker named volumes instead of bind mounts for data persistence

**Rationale:**
- **Platform Independence**: Named volumes work identically on Linux, macOS, Windows
- **Docker-Managed**: Docker handles storage location and permissions
- **Portability**: Volume names portable in docker-compose exports
- **Performance**: Better performance on macOS/Windows (no filesystem translation)

**Trade-offs:**
- **Data Access**: Harder to browse data directly from host
  - Impact: Users can't easily `cd /data` to inspect files
  - Mitigation: Use database-specific tools or `docker volume inspect`
- **Backup Complexity**: Can't simply `cp -r` data directory
  - Impact: Requires `docker cp` or native database backup tools
  - Mitigation: Document backup procedures
- **Hidden Location**: Users don't know where data is stored
  - Impact: Less transparency about disk usage
  - Mitigation: `docker volume inspect` shows mount point

**Alternative Considered:**
Bind mounts to local directories:
```bash
-v ./data/mongodb:/data/db
```

- Pros: Easy data access, simple backup, transparent location
- Cons: Path separators differ on Windows, permission issues, less portable

**Verdict:** Named volumes align with Docker best practices and provide better cross-platform consistency.

---

### Decision 5: Latest Tags vs Version Pinning

**Decision:** Use `:latest` tags for all database images (e.g., `mongo:latest`, `mysql:latest`)

**Rationale:**
- **Simplicity**: No version management required
- **Up-to-date**: Users get latest stable versions automatically
- **No Maintenance**: No need to update ContainDB when database versions release

**Trade-offs:**
- **Reproducibility**: Different users may get different versions over time
  - Impact: Subtle behavior differences across installations
  - Mitigation: Export includes specific image digest after pull
- **Breaking Changes**: Database major version upgrades may break compatibility
  - Impact: MySQL 8 vs MySQL 9, MongoDB 6 vs MongoDB 7
  - Mitigation: Users can manually specify versions in import YAML
- **Security**: Auto-pulling latest may introduce vulnerabilities
  - Impact: No control over what version is installed
  - Mitigation: Users should audit their environments

**Alternative Considered:**
Version pinning (e.g., `mongo:7.0`, `mysql:8.0`):
- Pros: Reproducible, predictable, safer
- Cons: Requires ContainDB updates when versions change, stale versions

**Verdict:** `:latest` is appropriate for local development tool. Production deployments should use pinned versions (user can modify exports).

---

## Operational Considerations

### System Requirements

**Minimum Requirements:**

**RAM:** 2GB
- Enforced by: `src/Docker/SysRequirement.go`
- Check: `getTotalRAM()` function
- Rationale: Docker Engine + database containers require memory
- Recommendation: 4GB for running multiple databases

**Disk Space:** 10GB free
- Enforced by: `src/Docker/SysRequirement.go`
- Check: `getAvailableDiskSpace()` function
- Rationale: Docker images (~500MB-1GB each) + data volumes
- Recommendation: 20GB for comfortable usage

**Operating System:**

| OS | Versions | Requirement | Notes |
|----|----------|-------------|-------|
| Linux | Ubuntu 20.04+, Debian 11+, CentOS 8+ | Sudo/root | Primary platform |
| macOS | 11+ (Big Sur and later) | Optional admin | Docker Desktop required |
| Windows | 10/11 (with WSL2) | Optional admin | Docker Desktop required |

**Docker Version:** 20.10.0+
- Recommendation: Docker 24.0+ for best compatibility
- Docker Compose: V2 (integrated with Docker CLI)

**Platform Detection:**

```go
// src/Docker/platform.go
func CheckOSSupport() error {
    switch runtime.GOOS {
    case "windows", "darwin", "linux":
        return nil
    default:
        return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
    }
}
```

**Unsupported Platforms:**
- FreeBSD
- Solaris
- Other Unix variants

---

### Resource Management

**Container Limits:**

ContainDB does **not** enforce resource limits on containers. Default Docker behavior applies:

- **CPU**: Unlimited (shared among all containers)
- **Memory**: Unlimited (can use all available host RAM)
- **Disk I/O**: Unlimited (shared host disk)

**User Responsibility:**
- Monitor resource usage via `docker stats`
- Set limits manually if needed: `docker update --memory 512m <container>`
- Consider resource constraints for production-like testing

**Recommendations:**

For resource-constrained systems (e.g., 4GB RAM laptop):
```bash
# Limit MongoDB memory
docker update --memory 512m mongodb-container

# Limit MySQL memory
docker update --memory 512m mysql-container
```

**Cleanup Recommendations:**

**Periodic Maintenance:**
```bash
# Remove stopped containers
docker container prune -f

# Remove unused images
docker image prune -a -f

# Remove unused volumes (⚠️ WARNING: DATA LOSS)
docker volume prune -f

# Full cleanup (⚠️ WARNING: REMOVES EVERYTHING)
docker system prune -a --volumes -f
```

**ContainDB Cleanup:**
- Use "Remove Database", "Remove Image", "Remove Volume" menu options
- Auto-rollback cleans up failed containers automatically
- No automatic cleanup of old/unused resources

**Volume Growth:**

Database volumes grow over time as data accumulates:
- **MongoDB**: WiredTiger storage engine compresses data
- **MySQL**: InnoDB files grow, rarely shrink
- **PostgreSQL**: VACUUM reclaims space but doesn't shrink files

**Monitoring:**
```bash
# Check volume sizes
docker system df -v

# Inspect specific volume
docker volume inspect mongodb-data
```

**Backup Before Cleanup:**
Always backup data before removing volumes:
```bash
# MongoDB
docker exec mongodb-container mongodump --out /data/backup

# MySQL
docker exec mysql-container mysqldump -u root -p --all-databases > backup.sql

# PostgreSQL
docker exec postgresql-container pg_dumpall -U postgres > backup.sql
```

---

### Monitoring & Logging

**Container Logs:**

```bash
# View logs for specific container
docker logs mongodb-container

# Follow logs in real-time
docker logs -f mysql-container

# View last 100 lines
docker logs --tail 100 postgresql-container
```

**Status Checking:**

```bash
# List all ContainDB containers
docker ps --filter network=ContainDB-Network

# Check container status
docker inspect mongodb-container --format '{{.State.Status}}'

# View resource usage
docker stats --filter network=ContainDB-Network
```

**Health Checks:**

ContainDB does **not** implement health checks. Users can add manually:

```yaml
# docker-compose.yml with health check
services:
  mongodb-container:
    image: mongo:latest
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

**Metrics Collection:**

ContainDB does **not** collect metrics. For production-like monitoring:

**External Tools:**
- **Prometheus + cAdvisor**: Container metrics
- **Grafana**: Visualization
- **ELK Stack**: Log aggregation
- **Datadog/New Relic**: APM (paid)

**Database-Specific Monitoring:**
- MongoDB: MongoDB Cloud Manager, Ops Manager
- MySQL: MySQL Enterprise Monitor, Percona Monitoring
- PostgreSQL: pgAdmin statistics, pg_stat_statements
- Redis: RedisInsight dashboard

**No Built-in Metrics:**
- ContainDB is a setup tool, not a monitoring platform
- Users expected to use database-native tools

---

## Appendices

### Code Metrics

**Total Codebase:** 2,635 lines of Go code
**Files:** 23 Go source files
**Packages:** 4 (Core, base, Docker, tools)

**Package Distribution:**

| Package | Files | Lines | Percentage | Responsibility |
|---------|-------|-------|------------|----------------|
| Docker | 10 | ~1,200 | 45% | Docker API abstraction |
| base | 7 | ~800 | 30% | User interface, workflows |
| tools | 6 | ~500 | 19% | Tool installers, utilities |
| Core | 1 | ~95 | 4% | Entry point, initialization |

**External Dependencies:**

| Library | Version | Usage | Lines Impacted |
|---------|---------|-------|----------------|
| promptui | v0.9.0 | Interactive prompts | ~300 |
| fatih/color | v1.18.0 | Colored output | ~50 |
| yaml.v2 | v2.4.0 | Docker Compose YAML | ~200 |
| readline | v1.5.1 | Terminal input (indirect) | N/A |

**Code Characteristics:**

- **Average Function Length:** 25 lines
- **Longest File:** `DockerComposeMaker.go` (352 lines)
- **Shortest File:** `Docker_Network.go` (21 lines)
- **Comments:** Minimal (mostly inline documentation)
- **Error Handling:** 263 error checks across 20 files

---

### File Reference

**Core Configuration:**

- `/home/ankan/Documents/Projects/ContainDB/go.mod` - Go module definition and dependencies
- `/home/ankan/Documents/Projects/ContainDB/VERSION` - Current version number (7.17.42-stable)
- `/home/ankan/Documents/Projects/ContainDB/npm/package.json` - NPM package metadata
- `/home/ankan/Documents/Projects/ContainDB/LICENSE` - MIT License

**Source Code:**

- `/home/ankan/Documents/Projects/ContainDB/src/Core/main.go` - Main entry point (95 lines)
- `/home/ankan/Documents/Projects/ContainDB/src/base/BaseCaseHandler.go` - Menu system (270 lines)
- `/home/ankan/Documents/Projects/ContainDB/src/Docker/DockerComposeMaker.go` - Export logic (352 lines)
- `/home/ankan/Documents/Projects/ContainDB/src/tools/rollback.go` - Cleanup mechanism (47 lines)

**Build & Distribution:**

- `/home/ankan/Documents/Projects/ContainDB/Scripts/BinBuilder.sh` - Multi-platform build script
- `/home/ankan/Documents/Projects/ContainDB/Scripts/PackageBuilder.sh` - Debian package builder
- `/home/ankan/Documents/Projects/ContainDB/Scripts/installer.sh` - Linux installer
- `/home/ankan/Documents/Projects/ContainDB/npm/InstallController.js` - NPM binary launcher

**Documentation:**

- `/home/ankan/Documents/Projects/ContainDB/README.md` - Main user documentation
- `/home/ankan/Documents/Projects/ContainDB/CONTRIBUTING.md` - Contribution guidelines
- `/home/ankan/Documents/Projects/ContainDB/DEPLOYMENT.md` - Deployment guide
- `/home/ankan/Documents/Projects/ContainDB/SECURITY.md` - Security policy
- `/home/ankan/Documents/Projects/ContainDB/ROADMAP.md` - Future feature plans
- `/home/ankan/Documents/Projects/ContainDB/CHANGELOG.md` - Version history

**CI/CD:**

- `/home/ankan/Documents/Projects/ContainDB/.github/workflows/run_build.yml` - Build automation
- `/home/ankan/Documents/Projects/ContainDB/.github/workflows/ReviewBuddy.yml` - Code review automation

---

### Glossary

**ContainDB-Network**
Dedicated Docker bridge network created by ContainDB. All managed containers join this network for DNS-based service discovery and isolation from other Docker containers.

**Named Volume**
Docker-managed persistent storage with a specific name (e.g., `mongodb-data`). Lifecycle independent of containers, survives container removal.

**Management Tool**
GUI application for database administration. Examples: phpMyAdmin (MySQL), pgAdmin (PostgreSQL), RedisInsight (Redis), MongoDB Compass (MongoDB).

**Auto-Rollback**
Automatic cleanup mechanism triggered by errors or user interrupts (Ctrl+C). Removes failed containers, dangling images, and temporary files to prevent partial states.

**Docker Compose Export**
Process of generating a `docker-compose.yml` file from running ContainDB containers. Captures configuration (not data) for portability.

**Docker Compose Import**
Process of deploying database services from an existing `docker-compose.yml` file. Creates volumes, pulls images, and starts containers.

**Interactive Prompts**
User interface pattern using `promptui` library. Presents menus and options for user selection instead of command-line arguments.

**Port Mapping**
Docker feature that exposes container ports on the host. Format: `-p <host-port>:<container-port>`. Example: `-p 27017:27017` makes MongoDB accessible at `localhost:27017`.

**SIGINT Handler**
Signal handler registered in `main.go` that catches Ctrl+C interrupts and executes cleanup before program exit.

**Data Persistence**
Feature allowing database data to survive container removal via named volumes. Optional during database installation.

**Stateless Abstraction**
Design pattern where Docker package functions have no internal state. Each function independently executes Docker CLI commands without relying on previous calls.

**Container Lifecycle**
Series of states a Docker container transitions through: Pull → Create → Start → Running → Stopped → Removed.

**Volume Reuse**
Scenario where user installs database with same name as previously removed container. ContainDB detects existing volume and prompts user to reuse or recreate.

**Environment Variable Filtering**
Process during Docker Compose export that excludes system-generated variables (PATH, PHP_*, MONGO_*) while preserving user-defined credentials.

**Cross-Compilation**
Building executables for multiple operating systems and architectures from a single source code base using Go's `GOOS` and `GOARCH` environment variables.

**Platform Detection**
Logic in `npm/InstallController.js` that determines user's operating system and architecture to select appropriate binary.

**Pre-flight Checks**
Validation steps performed before Docker Compose import: port availability, service name conflicts, volume conflicts.

---

### References

**Docker Documentation:**
- Docker CLI Reference: https://docs.docker.com/engine/reference/commandline/cli/
- Docker Compose Specification: https://docs.docker.com/compose/compose-file/
- Docker Networks: https://docs.docker.com/network/
- Docker Volumes: https://docs.docker.com/storage/volumes/

**Go Documentation:**
- Go Modules: https://go.dev/ref/mod
- os/exec Package: https://pkg.go.dev/os/exec
- runtime Package: https://pkg.go.dev/runtime

**Library Documentation:**
- promptui: https://github.com/manifoldco/promptui
- fatih/color: https://github.com/fatih/color
- gopkg.in/yaml.v2: https://pkg.go.dev/gopkg.in/yaml.v2

**Database Documentation:**
- MongoDB: https://docs.mongodb.com/
- MySQL: https://dev.mysql.com/doc/
- PostgreSQL: https://www.postgresql.org/docs/
- MariaDB: https://mariadb.com/kb/
- Redis: https://redis.io/documentation

**ContainDB Resources:**
- GitHub Repository: https://github.com/nexoral/ContainDB
- NPM Package: https://www.npmjs.com/package/containdb
- Issue Tracker: https://github.com/nexoral/ContainDB/issues
- Security Policy: https://github.com/nexoral/ContainDB/blob/main/SECURITY.md

---

## Document Conventions

**File Path References:**
- Format: `/absolute/path/to/file.go:line-range`
- Example: `/home/ankan/Documents/Projects/ContainDB/src/Core/main.go:50-57`

**Code Examples:**
- Language specified in code blocks (```go, ```bash, ```yaml)
- Inline code uses backticks: `docker run`

**Diagrams:**
- ASCII art for architecture, topology, and flow diagrams
- Consistent box-drawing characters: `┌─┐│└┘├┬┤┴┼`

**Terminology:**
- Docker container (not "container instance")
- Named volume (not "Docker volume")
- Management tool (not "admin tool" or "GUI")
- Database system (not "database engine")

---

*End of ContainDB System Design Document*
*For questions or contributions, see [CONTRIBUTING.md](CONTRIBUTING.md)*
