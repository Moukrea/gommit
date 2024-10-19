# Manual Gommit Integration Guide

This guide provides step-by-step instructions for manually integrating Gommit into your repository. Follow these steps if you prefer not to use the `gommit-integration` tool.

## 1. Create Required Directories

First, create the necessary directories in your repository:

```bash
mkdir -p .gommit
mkdir -p scripts/git-hooks
```

## 2. Download Gommit Downloader Binaries

Download the following binaries from the [Gommit releases page](https://github.com/Moukrea/gommit/releases) and place them in the `.gommit` directory:

- gommit-downloader-linux-amd64
- gommit-downloader-linux-arm64
- gommit-downloader-darwin-amd64
- gommit-downloader-darwin-arm64
- gommit-downloader-windows-amd64.exe
- gommit-downloader-windows-arm64.exe

Ensure all downloaded binaries have executable permissions:

```bash
chmod +x .gommit/gommit-downloader-*
```

## 3. Create Gommit Script

Create a file named `scripts/gommit.sh` with the following content:

```bash
#!/bin/sh
set -e

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux*)
        case "$ARCH" in
            x86_64) DOWNLOADER=".gommit/gommit-downloader-linux-amd64" ;;
            aarch64) DOWNLOADER=".gommit/gommit-downloader-linux-arm64" ;;
            *) echo "Unsupported Linux architecture: $ARCH" >&2; exit 1 ;;
        esac
        ;;
    darwin*)
        case "$ARCH" in
            x86_64) DOWNLOADER=".gommit/gommit-downloader-darwin-amd64" ;;
            arm64) DOWNLOADER=".gommit/gommit-downloader-darwin-arm64" ;;
            *) echo "Unsupported macOS architecture: $ARCH" >&2; exit 1 ;;
        esac
        ;;
    msys*|mingw*|cygwin*)
        case "$ARCH" in
            x86_64) DOWNLOADER=".gommit/gommit-downloader-windows-amd64.exe" ;;
            aarch64) DOWNLOADER=".gommit/gommit-downloader-windows-arm64.exe" ;;
            *) echo "Unsupported Windows architecture: $ARCH" >&2; exit 1 ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS" >&2
        exit 1
        ;;
esac

if [ ! -f "$DOWNLOADER" ]; then
    echo "Gommit downloader not found. Please check the integration setup." >&2
    exit 1
fi

"$DOWNLOADER"

if [ ! -f ".gommit/gommit" ]; then
    echo "Gommit binary not found. Please check the download process." >&2
    exit 1
fi

".gommit/gommit" "$1"
```

Make this script executable:

```bash
chmod +x scripts/gommit.sh
```

## 4. Create Git Hooks

### Commit Message Hook

Create `scripts/git-hooks/commit-msg` with the following content:

```bash
#!/bin/sh
./scripts/gommit.sh "$1"
```

### Post-Checkout Hook

Create `scripts/git-hooks/post-checkout` with the following content:

```bash
#!/bin/sh

# Check if this is the initial clone
if [ "$3" = "1" ]; then
    # This is a clone operation
    echo "Initial clone detected. Setting up Gommit..."
    ./scripts/setup-hooks.sh
else
    # This is a normal checkout
    ./scripts/setup-hooks.sh
fi
```

### Post-Merge Hook

Create `scripts/git-hooks/post-merge` with the following content:

```bash
#!/bin/sh
./scripts/setup-hooks.sh
```

Make all hooks executable:

```bash
chmod +x scripts/git-hooks/*
```

## 5. Create Setup Hooks Script

Create `scripts/setup-hooks.sh` with the following content:

```bash
#!/bin/sh
set -e

REPO_ROOT=$(git rev-parse --show-toplevel)

# Check if hooks are already set up
if [ "$(git config --get core.hooksPath)" = "$REPO_ROOT/scripts/git-hooks" ]; then
    echo "Gommit hooks are already set up."
else
    # Set up the hooks path
    git config core.hooksPath "$REPO_ROOT/scripts/git-hooks"
    echo "Git hooks have been set up successfully."
fi

# Ensure the Gommit binary is downloaded
"$REPO_ROOT/scripts/gommit.sh"

echo "Gommit is ready to use."
```

Make this script executable:

```bash
chmod +x scripts/setup-hooks.sh
```

## 6. Update .gitignore

Add the following line to your `.gitignore` file to exclude the downloaded Gommit binary:

```
.gommit/gommit
```

## 7. Commit and Push Changes

After creating all the necessary files and directories, commit and push the changes:

```bash
git add .
git commit -m "chore: manually integrate Gommit for commit message validation"
git push
```

## Conclusion

You have now manually integrated Gommit into your repository. All developers who clone or pull from this repository will automatically have Gommit set up for them.

For any issues or questions, please refer to the [Gommit GitHub repository](https://github.com/Moukrea/gommit).