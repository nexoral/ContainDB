#!/bin/bash
set -e


# Check which package is needed amd64, arm64, or i386
ARCH=$(dpkg --print-architecture)

echo "Detected architecture: $ARCH"

VERSION="6.15.41-stable"

if [[ "$ARCH" == "amd64" ]]; then
  PKG="containdb_${VERSION}_amd64.deb"
elif [[ "$ARCH" == "arm64" ]]; then
  PKG="containdb_${VERSION}_arm64.deb"
elif [[ "$ARCH" == "i386" ]]; then
  PKG="containdb_${VERSION}_i386.deb"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

URL="https://github.com/nexoral/ContainDB/releases/download/v${VERSION}/${PKG}"
echo "Downloading package: $PKG from $URL"

# Download package
wget -q $URL -O /tmp/$PKG

# shellcheck disable=SC2181
if [[ $? -ne 0 ]]; then
  echo "Failed to download package from $URL"
  exit 1
fi

echo "Download completed."

# Install package
sudo dpkg -i /tmp/$PKG

# shellcheck disable=SC2181
if [[ $? -ne 0 ]]; then
  echo "Failed to install package $PKG"
  exit 1
fi

echo "Installation completed successfully."

# Clean up
rm /tmp/$PKG
