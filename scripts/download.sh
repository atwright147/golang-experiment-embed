#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Create necessary directories
mkdir -p downloads
mkdir -p assets/nix
mkdir -p assets/windows

# Get the current ExifTool version
echo "Fetching current ExifTool version..."
VERSION=$(curl -s https://exiftool.org/ver.txt | tr -d '\r\n')
if [ -z "$VERSION" ]; then
    echo "Failed to retrieve version information"
    exit 1
fi
echo "Current ExifTool version: $VERSION"

# Download the tar.gz version (NIX)
NIX_URL="https://exiftool.org/Image-ExifTool-$VERSION.tar.gz"
NIX_FILE="downloads/Image-ExifTool-$VERSION.tar.gz"
echo "Downloading NIX (Unix/Linux) version from $NIX_URL..."
if curl -L -f -o "$NIX_FILE" "$NIX_URL"; then
    echo "Successfully downloaded NIX version to $NIX_FILE"
else
    echo "Failed to download NIX version from $NIX_URL"
    exit 1
fi

# Download the 64-bit zip version (WINDOWS)
WINDOWS_URL="https://exiftool.org/exiftool-${VERSION}_64.zip"
WINDOWS_FILE="downloads/exiftool-${VERSION}_64.zip"
echo "Downloading WINDOWS version from $WINDOWS_URL..."
if curl -L -f -o "$WINDOWS_FILE" "$WINDOWS_URL"; then
    echo "Successfully downloaded WINDOWS version to $WINDOWS_FILE"
else
    echo "Failed to download WINDOWS version from $WINDOWS_URL"
    exit 1
fi

# Verify the files exist and have non-zero size before extraction
if [ ! -s "$NIX_FILE" ]; then
    echo "NIX file is empty or does not exist"
    exit 1
fi

if [ ! -s "$WINDOWS_FILE" ]; then
    echo "WINDOWS file is empty or does not exist"
    exit 1
fi

# Extract the NIX (tar.gz) version to assets/nix with suppressed output
echo "Extracting NIX version to assets/nix..."
if tar -xzf "$NIX_FILE" -C assets/nix > /dev/null 2>&1; then
    echo "Successfully extracted NIX version to assets/nix"
else
    echo "Failed to extract NIX version"
    exit 1
fi

# Extract the WINDOWS (zip) version to assets/windows with suppressed output
echo "Extracting WINDOWS version to assets/windows..."
if unzip -q -o "$WINDOWS_FILE" -d assets/windows > /dev/null 2>&1; then
    echo "Successfully extracted WINDOWS version to assets/windows"
else
    echo "Failed to extract WINDOWS version"
    exit 1
fi

echo "All downloads and extractions completed successfully!"
