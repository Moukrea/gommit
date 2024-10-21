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

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
    OS="darwin"
elif [ "$OS" = "linux" ]; then
    OS="linux"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Download gommit
GOMMIT_URL="https://github.com/Moukrea/gommit/releases/latest/download/gommit-$OS-$ARCH"
download_file "$GOMMIT_URL" ".gommit/gommit"

# Make gommit executable
chmod +x .gommit/gommit

# Prepare commit-msg hook content for gommit
GOMMIT_HOOK_CONTENT="
# Gommit commit message validation
./.gommit/gommit \$1
exit \$?
"

# Handle commit-msg hook
DEST_FILE=".git/hooks/commit-msg"

if [ -f "$DEST_FILE" ]; then
    echo "Existing commit-msg hook found."
    if grep -q "/.gommit/gommit" "$DEST_FILE"; then
        echo "Gommit hook already exists in commit-msg. No changes needed."
    else
        printf "Choose action (overwrite/append/skip): "
        read -r choice
        case "$choice" in
            overwrite)
                printf "#!/bin/sh%s" "$GOMMIT_HOOK_CONTENT" > "$DEST_FILE"
                chmod +x "$DEST_FILE"
                echo "Overwrote commit-msg hook with Gommit."
                ;;
            append)
                if grep -q "exit" "$DEST_FILE"; then
                    sed -i.bak '/exit/i\'"$GOMMIT_HOOK_CONTENT" "$DEST_FILE"
                    rm "${DEST_FILE}.bak"
                else
                    echo "$GOMMIT_HOOK_CONTENT" >> "$DEST_FILE"
                fi
                echo "Appended Gommit to existing commit-msg hook."
                ;;
            skip)
                echo "Skipped modifying commit-msg hook."
                ;;
            *)
                echo "Invalid choice. Skipping commit-msg hook modification."
                ;;
        esac
    fi
else
    printf "#!/bin/sh%s" "$GOMMIT_HOOK_CONTENT" > "$DEST_FILE"
    chmod +x "$DEST_FILE"
    echo "Created new commit-msg hook with Gommit."
fi

echo "Gommit has been set up successfully."