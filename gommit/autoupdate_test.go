package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestShouldCheckForUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gommit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalUpdateCheckFile := getUpdateCheckFile()
	testUpdateCheckFile := filepath.Join(tempDir, ".gommit_last_update")
	setUpdateCheckFile(testUpdateCheckFile)
	t.Cleanup(func() { setUpdateCheckFile(originalUpdateCheckFile) })

	tests := []struct {
		name     string
		setup    func()
		expected bool
	}{
		{
			name: "No update file exists",
			setup: func() {
				os.Remove(testUpdateCheckFile)
			},
			expected: true,
		},
		{
			name: "Update file exists but is old",
			setup: func() {
				oldTime := time.Now().Add(-25 * time.Hour)
				os.WriteFile(testUpdateCheckFile, []byte{}, 0644)
				os.Chtimes(testUpdateCheckFile, oldTime, oldTime)
			},
			expected: true,
		},
		{
			name: "Update file exists and is recent",
			setup: func() {
				os.WriteFile(testUpdateCheckFile, []byte{}, 0644)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := shouldCheckForUpdate(); got != tt.expected {
				t.Errorf("shouldCheckForUpdate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetLatestRelease(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectedTag    string
		expectError    bool
	}{
		{
			name: "Valid release",
			serverResponse: `{
				"tag_name": "v1.2.3",
				"assets": [
					{
						"name": "gommit-linux-amd64",
						"browser_download_url": "https://example.com/gommit-linux-amd64"
					}
				]
			}`,
			statusCode:  http.StatusOK,
			expectedTag: "v1.2.3",
			expectError: false,
		},
		{
			name:           "Not Found",
			serverResponse: `{"message": "Not Found"}`,
			statusCode:     http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "Invalid JSON",
			serverResponse: `{"invalid": "json"`,
			statusCode:     http.StatusOK,
			expectError:    true,
		},
		{
			name:           "Empty response",
			serverResponse: ``,
			statusCode:     http.StatusOK,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/releases/latest" {
					t.Errorf("Expected path /releases/latest, got %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			originalGommitRepo := getGommitRepo()
			setGommitRepo(server.URL)
			t.Cleanup(func() { setGommitRepo(originalGommitRepo) })

			release, err := getLatestRelease()
			if (err != nil) != tt.expectError {
				t.Errorf("getLatestRelease() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && release != nil && release.TagName != tt.expectedTag {
				t.Errorf("getLatestRelease() tag = %v, want %v", release.TagName, tt.expectedTag)
			}
		})
	}
}

func TestGetOSAndArch(t *testing.T) {
	result := getOSAndArch()
	if result == "" {
		t.Errorf("getOSAndArch() returned empty string")
	}
}

func TestUpdateCheckFileWrappers(t *testing.T) {
	originalFile := getUpdateCheckFile()
	newFile := "/tmp/test_update_check"

	setUpdateCheckFile(newFile)
	if got := getUpdateCheckFile(); got != newFile {
		t.Errorf("getUpdateCheckFile() = %v, want %v", got, newFile)
	}

	setUpdateCheckFile(originalFile)
	if got := getUpdateCheckFile(); got != originalFile {
		t.Errorf("getUpdateCheckFile() = %v, want %v", got, originalFile)
	}
}

func TestGommitRepoWrappers(t *testing.T) {
	originalRepo := getGommitRepo()
	newRepo := "test/repo"

	setGommitRepo(newRepo)
	if got := getGommitRepo(); got != newRepo {
		t.Errorf("getGommitRepo() = %v, want %v", got, newRepo)
	}

	setGommitRepo(originalRepo)
	if got := getGommitRepo(); got != originalRepo {
		t.Errorf("getGommitRepo() = %v, want %v", got, originalRepo)
	}
}
