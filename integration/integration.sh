#!/bin/sh

set -e

# Function to download a file
download_file() {
    URL=$1
    OUTPUT=$2
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$OUTPUT" "$URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$OUTPUT" "$URL"
    elif command -v fetch >/dev/null 2>&1; then
        fetch -o "$OUTPUT" "$URL"
    else
        echo "Error: No supported download tool found (curl, wget, or fetch)."
        echo "Please install one of these tools and try again."
        exit 1
    fi
}

# Create .gommit directory
mkdir -p .gommit

# Download gommit-hook-setup.sh
download_file "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/gommit-hook-setup.sh" ".gommit/gommit-hook-setup.sh"
chmod +x .gommit/gommit-hook-setup.sh

# Download gommit-hook-setup.ps1
download_file "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/gommit-hook-setup.ps1" ".gommit/gommit-hook-setup.ps1"

# Download Makefile
download_file "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/Makefile" "Makefile"

echo "Gommit integration files have been set up successfully."
echo "You can now run 'make setup-gommit' to complete the setup."