package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	updateCheckInterval = 24 * time.Hour
	updateCheckFile     = ".gommit/last_update"
	gommitRepo          = "Moukrea/gommit"
)

var version = "dev" // This will be overwritten by the build process

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Add these wrapper functions
var (
	currentUpdateCheckFile = updateCheckFile
	currentGommitRepo      = gommitRepo
)

func ensureGommitDir() error {
	dir := filepath.Dir(currentUpdateCheckFile)
	return os.MkdirAll(dir, 0755)
}

func getUpdateCheckFile() string {
	return currentUpdateCheckFile
}

func setUpdateCheckFile(path string) {
	currentUpdateCheckFile = path
}

func getGommitRepo() string {
	return currentGommitRepo
}

func setGommitRepo(repo string) {
	currentGommitRepo = repo
}

func shouldCheckForUpdate() bool {
	info, err := os.Stat(getUpdateCheckFile())
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Error checking update file: %v\n", err)
		return false
	}
	return time.Since(info.ModTime()) > updateCheckInterval
}

func checkAndUpdate(isHook bool) error {
	if err := ensureGommitDir(); err != nil {
		return fmt.Errorf("failed to ensure .gommit directory exists: %w", err)
	}

	if !shouldCheckForUpdate() {
		return nil
	}

	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	if release.TagName > version {
		fmt.Printf("A new version of Gommit is available: %s\n", release.TagName)
		if !isHook {
			fmt.Print("Do you want to update? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				return nil
			}
		}
		if err := performUpdate(release); err != nil {
			return fmt.Errorf("failed to perform update: %w", err)
		}
		if isHook {
			fmt.Println("Gommit has been updated. Re-running the check...")
			if err := rerunCheck(); err != nil {
				return fmt.Errorf("failed to re-run check after update: %w", err)
			}
		} else {
			fmt.Println("Update successful. Please restart Gommit.")
			os.Exit(0)
		}
	}

	// Update the last check time
	if err := os.WriteFile(getUpdateCheckFile(), []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to update check file: %w", err)
	}

	return nil
}

func rerunCheck() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getLatestRelease() (*GithubRelease, error) {
	url := getGommitRepo() + "/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var release GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

func performUpdate(release *GithubRelease) error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	if err := downloadAndReplaceBinary(release, executable); err != nil {
		return fmt.Errorf("failed to download and replace binary: %w", err)
	}

	fmt.Println("Update successful. Please restart Gommit.")
	os.Exit(0)
	return nil
}

func downloadAndReplaceBinary(release *GithubRelease, executable string) error {
	assetURL := ""
	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, "gommit-"+getOSAndArch()) {
			assetURL = asset.BrowserDownloadURL
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("no suitable binary found for this system")
	}

	resp, err := http.Get(assetURL)
	if err != nil {
		return fmt.Errorf("failed to download new binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download new binary: status code %d", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "gommit-update")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write new binary: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	if err := os.Rename(tempFile.Name(), executable); err != nil {
		return fmt.Errorf("failed to replace old binary: %w", err)
	}

	return nil
}

func getOSAndArch() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}
