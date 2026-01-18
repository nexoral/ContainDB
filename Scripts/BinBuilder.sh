#!/bin/bash

# == Build ContainDB for Linux (Bin File) ==

set -e

BINARY_PATH="./bin/ContainDB"
BUILD_OUTPUT_DIR=$(dirname "$BINARY_PATH")

# Check if Go is installed
if ! command -v go &>/dev/null; then
  echo "Go not found. Installing Go via snap..."
  sudo snap install go --classic
else
  echo "Go is already installed."
fi

# Create output directory for binaries
mkdir -p "$BUILD_OUTPUT_DIR"


# Build the Go project with -o flag
echo "Building the project..."
go build -o "$BINARY_PATH" ./src/Core

echo "Build complete. Binary available at $BINARY_PATH"

echo "Building for Linux (AMD64)..."
GOOS=linux GOARCH=amd64 go build -o "${BUILD_OUTPUT_DIR}/containdb_linux_amd64" ./src/Core

echo "Building for macOS (Darwin AMD64)..."
GOOS=darwin GOARCH=amd64 go build -o "${BUILD_OUTPUT_DIR}/containdb_darwin_amd64" ./src/Core

echo "Building for macOS (Darwin ARM64)..."
GOOS=darwin GOARCH=arm64 go build -o "${BUILD_OUTPUT_DIR}/containdb_darwin_arm64" ./src/Core

echo "Building for Windows (AMD64)..."
GOOS=windows GOARCH=amd64 go build -o "${BUILD_OUTPUT_DIR}/containdb_windows_amd64.exe" ./src/Core


# Copy binaries to npm/bin for package distribution
NPM_BIN_DIR="./npm/bin"
mkdir -p "$NPM_BIN_DIR"

echo "Copying binaries to npm package..."
cp "${BUILD_OUTPUT_DIR}/containdb_linux_amd64" "$NPM_BIN_DIR/"
cp "${BUILD_OUTPUT_DIR}/containdb_darwin_amd64" "$NPM_BIN_DIR/"
cp "${BUILD_OUTPUT_DIR}/containdb_darwin_arm64" "$NPM_BIN_DIR/"
cp "${BUILD_OUTPUT_DIR}/containdb_windows_amd64.exe" "$NPM_BIN_DIR/"

echo "Copying metadata to npm package..."
cp "README.md" "./npm/"
cp "LICENSE" "./npm/"
cp "VERSION" "./npm/"

echo "Build complete. Binaries available in $BUILD_OUTPUT_DIR and $NPM_BIN_DIR"

 # Remove the binary file from the build output directory
echo "Removing binary file from build output directory..."
rm "${BUILD_OUTPUT_DIR}/containdb_linux_amd64"
rm "${BUILD_OUTPUT_DIR}/containdb_darwin_amd64"
rm "${BUILD_OUTPUT_DIR}/containdb_darwin_arm64"
rm "${BUILD_OUTPUT_DIR}/containdb_windows_amd64.exe"

