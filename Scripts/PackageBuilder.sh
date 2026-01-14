#!/bin/bash

set -e

# === CONFIG ===
APP_NAME="containdb"
AVAILABLE_OPTIONS=("amd64" "arm64" "i386")
BINARY_PATH="./bin/ContainDB" # Path to your Go binary
VERSION_FILE="./VERSION"
DIST_FOLDER="./dist"

# === Get version from VERSION file ===
if [ -f "$VERSION_FILE" ]; then
  VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
else
  echo "❌ VERSION file not found in project Core"
  exit 1
fi

# === Check binary ===
if [ ! -f "$BINARY_PATH" ]; then
  echo "❌ Binary not found at $BINARY_PATH"
  exit 1
fi

# === Install xpack ===
echo "🔧 Installing xpack..."
if ! command -v xpack &>/dev/null; then
  curl -fsSL https://raw.githubusercontent.com/nexoral/xpack/main/Scripts/installer.sh | sudo bash -
  echo "✅ xpack installed successfully"
else
  echo "✅ xpack already installed"
fi

# === Clean up old dist folder ===
rm -rf "$DIST_FOLDER"

# === Build packages for each architecture using xpack ===
for ARCH in "${AVAILABLE_OPTIONS[@]}"; do
  echo "📦 Building package for architecture: $ARCH using xpack"

  xpack -app "$APP_NAME" -arch "$ARCH" -v "$VERSION" -i "$BINARY_PATH"

  echo "✅ Package created: ${APP_NAME}_${VERSION}_${ARCH}"
done

echo "✅ All packages built successfully and available in $DIST_FOLDER"
