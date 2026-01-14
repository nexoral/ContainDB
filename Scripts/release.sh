#!/bin/bash

set -e

# === CONFIG ===
APP_NAME="containdb"
ARCH="amd64"
VERSION_FILE="./VERSION"
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')

# === Build steps ===
./Scripts/BinBuilder.sh
echo "🔨 Binary Building completed of $APP_NAME version $VERSION for $ARCH"

./Scripts/PackageBuilder.sh
echo "📦 Package Building completed of $APP_NAME version $VERSION for $ARCH"

# collect all packages from dist folder for this version
DIST_FOLDER="./dist"
DEB_FILES=("${DIST_FOLDER}"/${APP_NAME}_${VERSION}_*) # to collect all files

TAG="v$VERSION"
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)

# === Environment Variables ===
REPO="${GIT_REPOSITORY}" # GitHub Actions sets this automatically
TOKEN="${GIT_TOKEN}"     # GitHub Actions provides this

# === Checks ===
if [ ! -f "$VERSION_FILE" ]; then
  echo "❌ VERSION file not found"
  exit 1
fi

# ensure we have at least one package file
if [ ${#DEB_FILES[@]} -eq 0 ]; then
  echo "❌ No package files found for version $VERSION in $DIST_FOLDER"
  exit 1
fi

if ! command -v gh &>/dev/null; then
  echo "❌ GitHub CLI (gh) not installed"

  # Update package list and install dependencies
  apt update
  apt install -y curl gnupg software-properties-common

  # Add GitHub CLI's official package repository
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg |
    dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg

  sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" |
    tee /etc/apt/sources.list.d/github-cli.list >/dev/null

  # Install gh
  apt update
  apt install -y gh

fi

# === Authenticate GitHub CLI ===
echo "🔑 Authenticating GitHub CLI..."
echo "${TOKEN}" | gh auth login --with-token

# === Create Release ===
echo "📦 Creating GitHub release for tag $TAG..."

gh release create "$TAG" "${DEB_FILES[@]}" \
  --title "$TAG" \
  --notes "🔨 Commit: $COMMIT_HASH

📝 Message:
$COMMIT_MSG"

echo "✅ GitHub release published with .deb assets"

# --- Cleanup old releases, keep only latest two ---
echo "🗑️ Cleaning up old releases, retaining only latest two"
release_tags=($(gh release list --limit 1000 --json tagName,createdAt \
  --jq 'sort_by(.createdAt) | reverse | .[].tagName'))
  echo "Found ${#release_tags[@]} releases: ${release_tags[*]}"
if [ ${#release_tags[@]} -gt 2 ]; then
  for ((i = 2; i < ${#release_tags[@]}; i++)); do
    old_tag=${release_tags[i]}
    echo "🗑️ Deleting release $old_tag"
    gh release delete "$old_tag" --yes
    echo "🗑️ Deleting git tag $old_tag"
    gh api -X DELETE repos/${GITHUB_REPOSITORY}/git/refs/tags/${old_tag}
  done
fi
