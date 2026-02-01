<div align="center">
  <h1>ContainDB</h1>
  <p><strong>Database Container Management Made Simple</strong></p>
  <p>
    <a href="#installation">Installation</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#features">Features</a> •
    <a href="#usage-examples">Usage</a> •
    <a href="#architecture">Architecture</a> •
    <a href="#troubleshooting">Troubleshooting</a>
  </p>
  
  ![License](https://img.shields.io/badge/license-MIT-green)
  ![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)
  ![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey)
  ![npm](https://img.shields.io/badge/available%20via-npm-red)
</div>

## The Problem ContainDB Solves

As developers, we often face these frustrating scenarios:

- Spending hours configuring database environments across different projects
- Dealing with conflicting versions of databases on our development machines
- Struggling with complex Docker commands for simple database tasks
- Managing database persistence, networking, and tools separately
- Lack of a unified interface for different database systems
- **MongoDB "Core Dumped" errors on Debian-based systems** that are nearly impossible for beginners to troubleshoot

This last point was a major motivation for creating ContainDB. As a developer working on Debian-based systems, I repeatedly encountered the dreaded "Core Dumped" error when trying to install MongoDB natively. After spending countless hours troubleshooting compatibility issues, library dependencies, and system configurations, I realized there needed to be a better way.

ContainDB was born out of these pain points. I wanted a simple CLI tool that could handle all database container operations with minimal effort, allowing me to focus on actual development rather than environment setup.

## What is ContainDB?

ContainDB is an open-source CLI tool that automates the creation, management, and monitoring of containerized databases using Docker. It provides a simple, interactive interface for running popular databases and their management tools without needing to remember complex Docker commands or container configurations.

## Features

- **🚀 Instant Setup**: Get databases running in seconds with sensible defaults
- **🔄 Seamless Integration**: All databases run on the same Docker network for easy inter-container communication
- **💾 Data Persistence**: Optional volume management for data durability
- **🔐 Security Controls**: Interactive prompts for credentials and access control
- **🧩 Extensible Design**: Support for multiple database types and management tools
- **⚙️ Customization**: Configure ports, restart policies, and environment variables
- **📊 Management Tools**: One-click setup for phpMyAdmin, pgAdmin, RedisInsight, and MongoDB Compass
- **🧹 Easy Cleanup**: Simple commands to remove containers, images, and volumes
- **🧠 Smart Detection**: Checks for existing resources to avoid conflicts
- **🔄 Auto-Rollback**: Automatic cleanup of resources if any errors occur during setup
- **📦 Docker Compose Export**: Export your database configurations as a docker-compose.yml file that you can run anytime, anywhere
- **📥 Docker Compose Import**: Import and deploy services from existing docker-compose.yml files with automatic conflict resolution

## Installation

### Option 1: Using the Debian Package (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/nexoral/ContainDB/main/Scripts/installer.sh | sudo bash -

```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/nexoral/ContainDB.git
cd ContainDB

# Build the CLI
./Scripts/BinBuilder.sh

# Build Package (Debian)
./Scripts/PackageBuilder.sh

# Install the package the generated .deb file under the Packages directory
sudo dpkg -i Packages/containDB_*.deb
```

## Quick Start

Run ContainDB with root privileges:

```bash
sudo containDB
```

You'll be greeted with an attractive banner and a simple menu system that guides you through the process.

## Supported Databases & Tools

| Databases  | Management Tools |
|------------|-----------------|
| MongoDB    | MongoDB Compass  |
| MySQL      | phpMyAdmin       |
| PostgreSQL | pgAdmin          |
| MariaDB    | (uses phpMyAdmin)|
| Redis      | RedisInsight     |

## Usage Examples

### Installing a Database

```bash
sudo containDB
# Select "Install Database"
# Choose your database (e.g., "mongodb")
# Follow the interactive prompts
```

### Connecting to Your Database

After installation, ContainDB provides you with connection details:

```
✅ PostgreSQL started! Access it at http://localhost:5432
Link it to your DB container 'postgresql-container' inside pgAdmin.
📋 Connection information:
   - Container name: postgresql-container
   - IP Address: 172.18.0.2
   - Port: 5432
🔐 pgAdmin login credentials:
   - Email: admin@local.com
   - Password: yourpassword
```

### Setting Up Management Tools

```bash
sudo containDB
# Select "Install Database"
# Choose "phpMyAdmin", "PgAdmin", "Redis Insight", or "MongoDB Compass"
# Select the container to manage
# Follow the interactive prompts
```

#### Using RedisInsight with Your Redis Instance

After setting up a Redis container and launching RedisInsight:

1. Access the RedisInsight web interface at `http://localhost:8001` (or your custom port)
2. Add a new Redis database connection using:
   - Host: Your Redis container name (e.g., `redis-container`)
   - Port: `6379`
   - Use the Docker network's built-in DNS to connect automatically

```
✅ RedisInsight started. Access it at: http://localhost:8001
👉 In the RedisInsight UI, add a Redis database with host: `redis-container`, port: `6379`
   (RedisInsight will resolve container name using Docker network DNS.)
```

### Managing Existing Resources

```bash
sudo containDB
# Select "List Databases" to see running containers
# Select "Remove Database" to stop and remove containers
# Select "Remove Image" to delete Docker images
# Select "Remove Volume" to delete persistent data volumes
```

### Exporting Docker Compose Configuration

Export your running databases and management tools as a Docker Compose file:

```bash
sudo containDB --export
```

Or from the interactive menu:

```bash
sudo containDB
# Select "Export Services"
```

This creates a `docker-compose.yml` file in your current directory that you can use to recreate your entire database environment on any system with Docker:

```bash
# Move the docker-compose.yml to your project
cp docker-compose.yml /path/to/your/project/

# Run it anywhere
cd /path/to/your/project
docker-compose up -d
```

⚠️ **Important Note about Data Persistence**: The exported Docker Compose file contains only the configuration of your containers, not the actual database data. If you set up data persistence when installing a database, the exported file will reference the volume paths from your original machine. When running the exported compose file on another machine or after resetting your system, your previous data will not be available. For data backup and migration, you should use each database's native backup and restore functionality.

#### How the Export Feature Works Internally

---------------------------------------

```
┌────────────────────────────┐                ┌─────────────────────┐
│ ContainDB CLI              │                │ Running Docker      │
│ (export command)           │                │ Containers          │
└─────────────┬──────────────┘                └──────────┬──────────┘
              │                                          │
              │ 1. Identify running containers           │
              ├──────────────────────────────────────────┘
              │
              │                                ┌─────────────────────┐
              │ 2. Inspect container details   │ Container           │
              ├─────────────────────────────────> Configuration Data │
              │                                └─────────────────────┘
              │
              │ 3. Extract settings:
              │ - Image name & tag
              │ - Container name
              │ - Port mappings
              │ - Environment variables
              │ - Volume mounts
              │ - Network configuration
              │ - Restart policies
              │ - Command overrides
              │
              ▼
┌────────────────────────────┐
│ docker-compose.yml         │                 ┌─────────────────────┐
│ Generation                 │                 │ Local File System   │
└─────────────┬──────────────┘                 └──────────┬──────────┘
              │                                           │
              │ 4. Write docker-compose.yml file          │
              └───────────────────────────────────────────┘
```

This diagram shows how the ContainDB export feature captures the configuration of your running containers without copying the actual data stored in volumes. The generated docker-compose.yml provides a template for recreating your database infrastructure but requires separate data migration for full restoration.

---------------------------------------

### Importing Docker Compose Configuration

Import and deploy services from an existing docker-compose.yml file:

```bash
sudo containDB --import /path/to/docker-compose.yml
```

Or from the interactive menu:

```bash
sudo containDB
# Select "Import Services"
# Provide the path to your docker-compose.yml file
```

This feature analyzes your docker-compose.yml file and deploys all services with proper configuration:

```bash
# Example docker-compose.yml import
sudo containDB --import /home/user/my-project/docker-compose.yml
```

⚠️ **Important Note about Importing**: Before importing, ContainDB will check for port conflicts and existing volumes. You'll receive warnings about potential conflicts, allowing you to make decisions before deployment proceeds. The import feature intelligently handles Docker network creation and connects all imported services to the ContainDB network for seamless integration with your existing containers.

#### How the Import Feature Works Internally

---------------------------------------

```
┌────────────────────────────┐                ┌─────────────────────┐
│ ContainDB CLI              │                │ docker-compose.yml  │
│ (import command)           │                │ File                │
└─────────────┬──────────────┘                └──────────┬──────────┘
              │                                          │
              │ 1. Read & parse compose file             │
              ├──────────────────────────────────────────┘
              │
              │ 2. Validate Docker installation
              │
              │ 3. Check for existing services   ┌─────────────────────┐
              ├─────────────────────────────────>│ Running Docker      │
              │                                  │ Containers          │
              │                                  └─────────────────────┘
              │
              │ 4. Check for port conflicts      ┌─────────────────────┐
              ├─────────────────────────────────>│ Host System Ports   │
              │                                  └─────────────────────┘
              │
              │ 5. Set up volumes                ┌─────────────────────┐
              ├─────────────────────────────────>│ Docker Volumes      │
              │                                  └─────────────────────┘
              │
              ▼
┌────────────────────────────┐
│ Docker Compose Command     │                 ┌─────────────────────┐
│ Execution                  │---------------->│ Deployed Services   │
└────────────────────────────┘                 └─────────────────────┘
```

This diagram illustrates how ContainDB imports services from a docker-compose.yml file, checking for conflicts and setting up the necessary resources before deploying the services to ensure a smooth integration with your existing environment.

---------------------------------------

## Architecture

ContainDB follows a layered architecture that separates concerns and promotes code organization.

```
                      User Interaction
                            ↓
      +-------------------------------------------+
      |             Main CLI Interface            |
      +-------------------------------------------+
                ↑                      ↑
                |                      |
       +------------------+   +------------------+
       |  Base Operations |   |   Tool Helpers   |
       +------------------+   +------------------+
                ↑                      ↑
                |                      |
       +------------------+   +------------------+
       | Docker Interface |   | System Utilities |
       +------------------+   +------------------+
                            ↓
                    Docker Engine & Host System
```

### How ContainDB Works Internally

---------------------------------------

1. **Network Creation**

ContainDB first ensures that a dedicated Docker network (`ContainDB-Network`) exists, which allows all containers to communicate with each other using container names as hostnames.

```
┌────────────────────────────┐     creates    ┌─────────────────────┐
│ ContainDB CLI              │ ───────────────> Docker Network      │
└────────────────────────────┘                └─────────────────────┘
```

---------------------------------------

2. **Container Orchestration**

When you select a database to install, ContainDB:
- Pulls the latest image (if needed)
- Checks for port conflicts
- Sets up necessary volumes for persistence
- Configures environment variables
- Creates and starts the container

```
┌────────────────────────────┐                ┌─────────────────────┐
│ ContainDB CLI              │                │ Docker Hub          │
└─────────────┬──────────────┘                └──────────┬──────────┘
              │                                          │
              │ 1. Pull image                            │
              ├──────────────────────────────────────────┘
              │
              │ 2. Create volume                ┌─────────────────────┐
              ├─────────────────────────────────> Data Volume         │
              │                                 └─────────────────────┘
              │
              │ 3. Run container               ┌──────────────────────┐
              └────────────────────────────────> Database Container   │
                                               └──────────────────────┘
```

---------------------------------------

3. **Management Tool Integration**

For tools like phpMyAdmin, pgAdmin, or MongoDB Compass, ContainDB handles:
- Tool installation and configuration
- Linking to the appropriate database container
- Providing connection details and credentials

```
┌────────────────────────────┐                ┌─────────────────────┐
│ ContainDB CLI              │ ─────────────┐ │ Database Container  │
└────────────────────────────┘              │ └─────────────────────┘
                                            │
                                            │ links
                                            ▼
                                  ┌─────────────────────┐
                                  │ Management Tool     │
                                  │ Container/App       |
                                  └─────────────────────┘
```

---------------------------------------

4. **Auto-Rollback Mechanism**

ContainDB implements a robust error handling system that automatically cleans up any partially created resources if something goes wrong:

```
┌────────────────────────────┐
│ Operation Started          │
└─────────────┬──────────────┘
              │
              ▼
┌────────────────────────────┐
│ Resource Creation          │
└─────────────┬──────────────┘
              │
              ▼
      ┌───────/\───────┐
      │ Error Occurs?  │
      └───────┬\───────┘
              │
          Yes │             ┌─────────────────────┐
              ├─────────────> Stop Containers     │
              │             └─────────┬───────────┘
              │                       │
              │             ┌─────────▼───────────┐
              │             │ Remove Containers   │
              │             └─────────┬───────────┘
              │                       │
              │             ┌─────────▼───────────┐
              │             │ Clean Temp Files    │
              │             └─────────────────────┘
              │
         No   ▼
┌────────────────────────────┐
│ Operation Completed        │
└────────────────────────────┘
```

This ensures your system stays clean even if an operation fails midway, preventing orphaned containers or dangling resources.

---------------------------------------

## The Origin Story

ContainDB was created after I found myself repeatedly setting up the same database environments across different projects. The challenges included:

1. Remembering specific Docker commands for each database
2. Managing network connectivity between containers
3. Ensuring data persistence with proper volume configuration
4. Setting up administration tools for each database type

The final straw came when I needed to set up a multi-database project that required MongoDB, PostgreSQL, and Redis - all with different configurations, management tools, and persistence requirements. I spent hours on environment setup instead of actual coding.

**The MongoDB Crisis:** The breaking point was when I tried to deploy MongoDB on a new Debian-based system. Instead of a working database, I got the cryptic "Illegal instruction (core dumped)" error - a common issue on certain Debian systems due to CPU instruction set incompatibilities with pre-built MongoDB binaries. After wasting a full day on troubleshooting this single issue, I realized containerization was the answer.

I realized that these repetitive tasks could be automated, and ContainDB was born. What started as a personal script evolved into a comprehensive tool that I now use daily and want to share with the developer community.

## How ContainDB Helps in Daily Development

ContainDB has become an essential part of my development workflow by:

- **Saving Setup Time**: What used to take 30+ minutes now takes seconds
- **Standardizing Environments**: Ensuring consistent database setups across projects
- **Simplifying Management**: Providing easy access to admin tools and interfaces
- **Isolating Services**: Preventing conflicts between different database versions
- **Managing Resources**: Making cleanup and maintenance straightforward
- **Bypassing System Compatibility Issues**: Avoiding the notorious MongoDB "core dumped" errors on Debian systems
- **Visual Database Management**: Quick setup of GUI tools like RedisInsight for better productivity

Real-world example: When working on a new microservice project, I can spin up a PostgreSQL instance, link it to pgAdmin, and have a fully functional development environment in less than a minute - all with proper network configuration and persistence.

## Troubleshooting

### Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| **"Permission Denied"** | Ensure you run ContainDB with `sudo` |
| **"Docker Not Found"** | Let ContainDB install Docker or run with `--install-docker` flag |
| **"Port Already in Use"** | Choose a different port when prompted |
| **"Volume Already Exists"** | Select to reuse or recreate the volume |
| **"Cannot Connect to Database"** | Check network settings and credentials |

## Contributing

ContainDB is an open-source project that welcomes contributions from everyone. Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit your changes**: `git commit -m 'Add amazing feature'`
4. **Push to the branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request`

For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- The Docker team for creating an amazing containerization platform
- The Go community for providing excellent libraries and tools
- Redis Labs for the RedisInsight tool
- All contributors who have helped improve ContainDB

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/AnkanSaha">Ankan Saha</a></p>
  <p>
    <a href="https://github.com/nexoral/ContainDB/stargazers">⭐ Star this project</a> •
    <a href="https://github.com/nexoral/ContainDB/issues">🐞 Report Bug</a> •
    <a href="https://github.com/nexoral/ContainDB/issues">✨ Request Feature</a>
  </p>
</div>
</div>
