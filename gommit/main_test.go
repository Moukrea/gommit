package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	TEST_DIR                  = "gommit-test"
	TEMP_DIR_CREATION_FAILURE = "Failed to create temp dir: %v"
	COMMIT_MSG_EXAMPLE        = "feat: add new feature"
)

// Add this mock struct for testing
type MockConfigPathGetter struct {
	ConfigPath string
}

func (m MockConfigPathGetter) GetConfigPath() (string, error) {
	return m.ConfigPath, nil
}

func TestValidateCommitMsg(t *testing.T) {
	config := Config{
		HeaderMaxLength:   50,
		BodyLineMaxLength: 72,
		AllowedTypes:      []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert"},
	}

	tests := []struct {
		name                string
		msg                 string
		expectedErrors      []string
		needsBreakingChange bool
	}{
		{
			name:                "Valid commit message",
			msg:                 COMMIT_MSG_EXAMPLE,
			expectedErrors:      nil,
			needsBreakingChange: false,
		},
		{
			name:                "Invalid type",
			msg:                 "invalid: this is not a valid type",
			expectedErrors:      []string{"Header must be in format: <type>[optional scope][!]: <description>", "Type 'invalid' is not allowed. Allowed types are: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert"},
			needsBreakingChange: false,
		},
		{
			name:                "Header too long",
			msg:                 "feat: this header is way too long and exceeds the maximum length",
			expectedErrors:      []string{"Header must not exceed 50 characters"},
			needsBreakingChange: false,
		},
		{
			name:                "Breaking change without footer",
			msg:                 "feat!: add breaking change",
			expectedErrors:      nil,
			needsBreakingChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, needsBreakingChange := validateCommitMsg(tt.msg, config)
			if !reflect.DeepEqual(errors, tt.expectedErrors) {
				t.Errorf("validateCommitMsg() errors = %v, want %v", errors, tt.expectedErrors)
			}
			if needsBreakingChange != tt.needsBreakingChange {
				t.Errorf("validateCommitMsg() needsBreakingChange = %v, want %v", needsBreakingChange, tt.needsBreakingChange)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", TEST_DIR)
	if err != nil {
		t.Fatalf(TEMP_DIR_CREATION_FAILURE, err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock config file
	configContent := []byte(`
header_max_length: 60
body_line_max_length: 80
allowed_types:
  - feat
  - fix
  - docs
`)
	configPath := filepath.Join(tempDir, "gommit.conf.yaml")
	err = os.WriteFile(configPath, configContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock config file: %v", err)
	}

	// Use the mock ConfigPathGetter
	mockPathGetter := MockConfigPathGetter{ConfigPath: configPath}

	// Test loading the config
	config, err := loadConfig(mockPathGetter)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	// Check if the loaded config matches the expected values
	expectedConfig := Config{
		HeaderMaxLength:   60,
		BodyLineMaxLength: 80,
		AllowedTypes:      []string{"feat", "fix", "docs"},
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("loadConfig() = %v, want %v", config, expectedConfig)
	}
}

func TestContainsBreakingChange(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected bool
	}{
		{
			name:     "Contains breaking change",
			lines:    []string{"Some text", "BREAKING CHANGE: This is a breaking change"},
			expected: true,
		},
		{
			name:     "No breaking change",
			lines:    []string{"Some text", "This is not a breaking change"},
			expected: false,
		},
		{
			name:     "Empty lines",
			lines:    []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsBreakingChange(tt.lines)
			if result != tt.expected {
				t.Errorf("containsBreakingChange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAppendBreakingChange(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		description string
		expected    string
	}{
		{
			name:        "Append to message without footer",
			msg:         "feat!: add new feature",
			description: "This breaks the API",
			expected:    "feat!: add new feature\n\nBREAKING CHANGE: This breaks the API",
		},
		{
			name:        "Append to message with existing footer",
			msg:         "feat!: add new feature\n\nReviewed-by: John Doe",
			description: "This breaks the API",
			expected:    "feat!: add new feature\n\nReviewed-by: John Doe\n\nBREAKING CHANGE: This breaks the API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendBreakingChange(tt.msg, tt.description)
			if result != tt.expected {
				t.Errorf("appendBreakingChange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadFromStdin(t *testing.T) {
	// Save the original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single line input",
			input:    "Test input\n",
			expected: "Test input\n",
		},
		{
			name:     "Multi-line input",
			input:    "Line 1\nLine 2\nLine 3\n",
			expected: "Line 1\nLine 2\nLine 3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new pipe and replace stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write the test input
			go func() {
				w.Write([]byte(tt.input))
				w.Close()
			}()

			// Call the function
			result, err := readFromStdin()
			if err != nil {
				t.Fatalf("readFromStdin() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("readFromStdin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadCommitMsg(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", TEST_DIR)
	if err != nil {
		t.Fatalf(TEMP_DIR_CREATION_FAILURE, err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Single line commit message",
			content:  COMMIT_MSG_EXAMPLE,
			expected: COMMIT_MSG_EXAMPLE,
		},
		{
			name: "Multi-line commit message",
			content: `feat: add new feature

This is a detailed description of the new feature.

Closes #123`,
			expected: `feat: add new feature

This is a detailed description of the new feature.

Closes #123`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test content
			tmpfile, err := os.CreateTemp(tempDir, "commit-msg")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Call the function
			result, err := readCommitMsg(tmpfile.Name())
			if err != nil {
				t.Fatalf("readCommitMsg() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("readCommitMsg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWriteCommitMsg(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", TEST_DIR)
	if err != nil {
		t.Fatalf(TEMP_DIR_CREATION_FAILURE, err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Single line commit message",
			content: COMMIT_MSG_EXAMPLE,
		},
		{
			name: "Multi-line commit message",
			content: `feat: add new feature

This is a detailed description of the new feature.

Closes #123`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file path
			tmpfile := filepath.Join(tempDir, "commit-msg")

			// Call the function
			err := writeCommitMsg(tmpfile, tt.content)
			if err != nil {
				t.Fatalf("writeCommitMsg() error = %v", err)
			}

			// Read the content back
			content, err := os.ReadFile(tmpfile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			if string(content) != tt.content {
				t.Errorf("writeCommitMsg() wrote %v, want %v", string(content), tt.content)
			}
		})
	}
}

func TestIsRuleEnabled(t *testing.T) {
	config := Config{
		DisabledRules: []string{"rule1", "rule2"},
	}

	tests := []struct {
		name     string
		ruleName string
		expected bool
	}{
		{
			name:     "Disabled rule",
			ruleName: "rule1",
			expected: false,
		},
		{
			name:     "Enabled rule",
			ruleName: "rule3",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRuleEnabled(config, tt.ruleName)
			if result != tt.expected {
				t.Errorf("isRuleEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "Item exists in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "Item does not exist in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "date",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			item:     "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("contains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Note: Testing promptForBreakingChange is challenging due to user input.
// We could consider refactoring it to accept an io.Reader for easier testing,
// or mock the user input for testing purposes.

// Note: The following functions are not covered by unit tests due to their complexity
// and reliance on external systems or user input:
//
// - runGommit: This function orchestrates the entire commit message validation process
//   and involves user interaction, making it challenging to test in isolation.
//
// - checkAndUpdate (in autoupdate.go): This function involves network calls and user input,
//   which are difficult to mock in a unit test environment.
//
// - performUpdate (in autoupdate.go): This function performs system-level operations
//   like replacing the executable, which is not suitable for unit testing.
//
// - downloadAndReplaceBinary (in autoupdate.go): This function involves network operations
//   and file system changes, making it more appropriate for integration testing.
//
// To test these functions, consider creating integration tests or refactoring them
// into smaller, more testable units.
