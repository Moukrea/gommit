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
