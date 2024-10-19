package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	gommitRepo = "Moukrea/gommit"
	gommitDir  = ".gommit"
)

var downloaderBinaries = []string{
	"gommit-downloader-linux-amd64",
	"gommit-downloader-linux-arm64",
	"gommit-downloader-darwin-amd64",
	"gommit-downloader-darwin-arm64",
	"gommit-downloader-windows-amd64.exe",
	"gommit-downloader-windows-arm64.exe",
}

func main() {
	fmt.Println("Gommit Integration Setup")
	fmt.Println("This program will integrate Gommit into your repository.")

	if err := setupGommitIntegration(); err != nil {
		fmt.Printf("Error in Gommit integration setup: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Gommit integration setup completed successfully!")
	fmt.Println("Please add, commit, and push the changes to your repository.")
	fmt.Println("Gommit will now be automatically enabled for all developers who clone this repository.")
}

func setupGommitIntegration() error {
	if err := createDirectoryIfNotExists(gommitDir); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", gommitDir, err)
	}
	if err := createDirectoryIfNotExists(filepath.Join("scripts", "git-hooks")); err != nil {
		return fmt.Errorf("failed to create git-hooks directory: %w", err)
	}

	if err := downloadAllDownloaderBinaries(); err != nil {
		return fmt.Errorf("failed to download Gommit downloader binaries: %w", err)
	}

	if err := createGommitScript(); err != nil {
		return fmt.Errorf("failed to create gommit script: %w", err)
	}
	if err := createCommitMsgHook(); err != nil {
		return fmt.Errorf("failed to create commit-msg hook: %w", err)
	}
	if err := createSetupHooksScript(); err != nil {
		return fmt.Errorf("failed to create setup-hooks script: %w", err)
	}
	if err := createPostCheckoutHook(); err != nil {
		return fmt.Errorf("failed to create post-checkout hook: %w", err)
	}
	if err := createPostMergeHook(); err != nil {
		return fmt.Errorf("failed to create post-merge hook: %w", err)
	}

	return nil
}

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
		fmt.Printf("Created directory: %s\n", path)
	}
	return nil
}

func downloadAllDownloaderBinaries() error {
	for _, binary := range downloaderBinaries {
		if err := downloadLatestBinary(gommitRepo, binary, filepath.Join(gommitDir, binary)); err != nil {
			return fmt.Errorf("failed to download %s: %w", binary, err)
		}
	}
	return nil
}

func downloadLatestBinary(repo, binaryName, outputPath string) error {
	url := fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", repo, binaryName)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", binaryName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: unexpected status code %d", binaryName, resp.StatusCode)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file for %s: %w", binaryName, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", binaryName, err)
	}

	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions for %s: %w", binaryName, err)
	}

	fmt.Printf("Downloaded: %s\n", binaryName)
	return nil
}

func createGommitScript() error {
	content := `#!/bin/sh
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
`
	return writeFile(filepath.Join("scripts", "gommit.sh"), content)
}

func createCommitMsgHook() error {
	content := `#!/bin/sh
./scripts/gommit.sh "$1"
`
	return writeFile(filepath.Join("scripts", "git-hooks", "commit-msg"), content)
}

func createSetupHooksScript() error {
	content := `#!/bin/sh
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
`
	return writeFile(filepath.Join("scripts", "setup-hooks.sh"), content)
}

func createPostCheckoutHook() error {
	content := `#!/bin/sh

# Check if this is the initial clone
if [ "$3" = "1" ]; then
    # This is a clone operation
    echo "Initial clone detected. Setting up Gommit..."
    ./scripts/setup-hooks.sh
else
    # This is a normal checkout
    ./scripts/setup-hooks.sh
fi
`
	return writeFile(filepath.Join("scripts", "git-hooks", "post-checkout"), content)
}

func createPostMergeHook() error {
	content := `#!/bin/sh
./scripts/setup-hooks.sh
`
	return writeFile(filepath.Join("scripts", "git-hooks", "post-merge"), content)
}

func writeFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	fmt.Printf("Created: %s\n", path)
	return nil
}