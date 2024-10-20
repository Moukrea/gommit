package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	gommitRepo = "Moukrea/gommit"
	gommitDir  = ".gommit"
)

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	if err := downloadGommit(); err != nil {
		fmt.Printf("Error downloading Gommit: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Gommit downloaded successfully")
}

func downloadGommit() error {
	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	assetURL, err := getAssetURL(release)
	if err != nil {
		return fmt.Errorf("failed to get asset URL: %w", err)
	}

	if err := downloadAndSaveBinary(assetURL); err != nil {
		return fmt.Errorf("failed to download and save binary: %w", err)
	}

	return nil
}

func getLatestRelease() (*GithubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", gommitRepo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

func getAssetURL(release *GithubRelease) (string, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	assetName := fmt.Sprintf("gommit-integration-%s-%s", os, arch)
	if os == "windows" {
		assetName += ".exe"
	}

	url := fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", gommitRepo, assetName)
	return url, nil
}

func downloadAndSaveBinary(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := os.MkdirAll(gommitDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", gommitDir, err)
	}

	outputPath := filepath.Join(gommitDir, "gommit")
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write binary: %w", err)
	}

	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	return nil
}
