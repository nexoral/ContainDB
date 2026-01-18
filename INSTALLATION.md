# Installation Guide

This document provides detailed instructions for installing ContainDB on various Linux distributions.

## System Requirements

- Linux-based operating system (Ubuntu, Debian, Fedora, etc.)
- Root/sudo privileges
- Internet connectivity for downloading packages
- 50MB of disk space (plus additional space for database data)

## Installation Methods

ContainDB can be installed using one of the following methods:

### Method 1: Using the Debian Package (Recommended for Debian-based systems)

This is the simplest method for Debian-based systems like Ubuntu, Debian, Linux Mint, etc.

```bash
# Download the latest .deb release
wget https://github.com/nexoral/ContainDB/releases/download/v5.16.37-stable/containDB_5.16.37-stable_amd64.deb

# Install the package
sudo dpkg -i containDB_5.16.37-stable_amd64.deb

# If you see dependency errors, run:
sudo apt-get install -f
```

### Method 2: Using the Tarball (For any Linux distribution)

```bash
# Download the latest tar.gz release
wget https://github.com/nexoral/ContainDB/releases/download/v4.12.18-stable/containDB_4.12.18-stable_amd64.tar.gz

# Extract the archive
tar -xzf containDB_4.12.18-stable_amd64.tar.gz

# Move the binary to a directory in your PATH
sudo mv containDB /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/containDB
```

### Method 3: Building from Source

This method gives you the most control and works on any Linux distribution with Go installed.

```bash
# Install Go (if not already installed)
# For Ubuntu/Debian:
sudo apt-get update
sudo apt-get install golang-go

# For Fedora/RHEL:
sudo dnf install golang

# Clone the repository
git clone https://github.com/nexoral/ContainDB.git
cd ContainDB

# Build the CLI
./Scripts/BinBuilder.sh

# Install binary to /usr/local/bin
sudo mv ./bin/ContainDB /usr/local/bin/containDB
```

## Verifying the Installation

To verify that ContainDB is correctly installed, run:

```bash
sudo containDB --version
```

You should see the version number displayed (e.g., "ContainDB CLI Version: 4.12.18-stable").

## Docker Installation

ContainDB requires Docker to run. If Docker is not installed on your system, ContainDB will offer to install it for you when you first run it. Alternatively, you can install Docker manually using the official instructions: [Install Docker Engine](https://docs.docker.com/engine/install/).

## Upgrading ContainDB

### Upgrading a Debian Package Installation

```bash
# Download the new version
wget https://github.com/nexoral/ContainDB/releases/download/v[NEW_VERSION]/containDB_[NEW_VERSION]_amd64.deb

# Install the update
sudo dpkg -i containDB_[NEW_VERSION]_amd64.deb
```

### Upgrading a Tarball Installation

```bash
# Download the new version
wget https://github.com/nexoral/ContainDB/releases/download/v[NEW_VERSION]/containDB_[NEW_VERSION]_amd64.tar.gz

# Extract the archive
tar -xzf containDB_[NEW_VERSION]_amd64.tar.gz

# Replace the old binary
sudo mv containDB /usr/local/bin/

# Make sure it's executable
sudo chmod +x /usr/local/bin/containDB
```

### Upgrading a Source Installation

```bash
cd ContainDB
git pull
./Scripts/BinBuilder.sh
sudo mv ./bin/ContainDB /usr/local/bin/containDB
```

## Uninstalling ContainDB

### Uninstalling a Debian Package Installation

```bash
sudo apt-get remove containdb
```

### Uninstalling a Tarball or Source Installation

```bash
sudo rm /usr/local/bin/containDB
```

Note: Uninstalling ContainDB does not remove any Docker containers, images, or volumes that were created using ContainDB. To remove these resources, use the ContainDB CLI's cleanup options or Docker commands.

## Troubleshooting Installation Issues

### "Command not found" Error

If you see "command not found" when trying to run ContainDB, ensure that:

1. The installation completed successfully
2. The binary is in a directory in your PATH
3. The binary has executable permissions

### Permission Denied

ContainDB requires root privileges to manage Docker containers. Always run it with `sudo`.

### Docker-related Errors

If you encounter Docker-related errors:

1. Verify Docker is installed: `docker --version`
2. Ensure Docker service is running: `sudo systemctl status docker`
3. Check that your user has permission to use Docker: `sudo usermod -aG docker $USER`

For additional installation help, please [open an issue](https://github.com/nexoral/ContainDB/issues) on our GitHub repository.
