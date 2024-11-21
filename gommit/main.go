package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

const (
	COMMIT_MSG_INVALID_MSG = "✘ Commit message is invalid."
	AUTO_BREAKING_CHANGE   = "auto-breaking-change"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	headerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	detailStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
)

var (
	typePattern           = `^(feat|fix|build|chore|ci|docs|style|refactor|perf|test)`
	scopePattern          = `(\([a-z0-9\-]+\))?`
	breakingChangeMarker  = `!?`
	descriptionPattern    = `: .+`
	headerPattern         = regexp.MustCompile(typePattern + scopePattern + breakingChangeMarker + descriptionPattern + `$`)
	footerPattern         = regexp.MustCompile(`^([A-Z\-]+)(\s+)?:(\s+)?(.+)$`)
	breakingChangePattern = regexp.MustCompile(`^BREAKING[\s-]CHANGE: `)
)

type Rule struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Config struct {
	DisabledRules     []string `yaml:"disabled_rules"`
	HeaderMaxLength   int      `yaml:"header_max_length"`
	BodyLineMaxLength int      `yaml:"body_line_max_length"`
	AllowedTypes      []string `yaml:"allowed_types"`
}

var defaultConfig = Config{
	HeaderMaxLength:   50,
	BodyLineMaxLength: 72,
	AllowedTypes: []string{
		"feat", "fix", "docs", "style", "refactor",
		"perf", "test", "build", "ci", "chore", "revert",
	},
}

var defaultRules = []Rule{
	{Name: "header-format", Description: "Header must be in format: <type>[optional scope][!]: <description>"},
	{Name: "header-max-length", Description: "Header must not exceed the configured max length"},
	{Name: "header-lowercase", Description: "Header (short description) must be all lowercase"},
	{Name: "description-case", Description: "Description must start with lowercase"},
	{Name: "body-line-max-length", Description: "Body lines must not exceed the configured max length"},
	{Name: "footer-format", Description: "Footer must be in format: <token>: <value>"},
	{Name: "breaking-change", Description: "Breaking changes must be indicated in footer"},
	{Name: AUTO_BREAKING_CHANGE, Description: "Automatically add BREAKING CHANGE to footer when '!' is present in header"},
	{Name: "type-enum", Description: "Type must be one of the allowed types"},
	{Name: "type-case", Description: "Type must be in lowercase"},
	{Name: "type-empty", Description: "Type must not be empty"},
	{Name: "scope-case", Description: "Scope must be in lowercase"},
	{Name: "subject-empty", Description: "Subject must not be empty"},
}

type model struct {
	textInput textinput.Model
	err       error
}

func initialModel(initialContent string) model {
	ti := textinput.New()
	ti.Placeholder = "Edit your commit message"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80
	ti.SetValue(initialContent)

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Edit your commit message:\n\n%s\n\n%s",
		m.textInput.View(),
		"(Press Enter to save, or Esc to cancel)",
	) + "\n"
}

type ConfigPathGetter interface {
	GetConfigPath() (string, error)
}

type DefaultConfigPathGetter struct{}

func (d DefaultConfigPathGetter) GetConfigPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %w", err)
	}
	execDir := filepath.Dir(ex)
	return filepath.Join(execDir, "gommit.conf.yaml"), nil
}

func loadConfig(pathGetter ConfigPathGetter) (Config, error) {
	configPath, err := pathGetter.GetConfigPath()
	if err != nil {
		return Config{}, err
	}

	config := defaultConfig

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Check for .gommit/gommit.conf.yaml in the current working directory
			cwd, err := os.Getwd()
			if err != nil {
				return config, fmt.Errorf("error getting current working directory: %w", err)
			}
			localConfigPath := filepath.Join(cwd, ".gommit", "gommit.conf.yaml")
			data, err = os.ReadFile(localConfigPath)
			if err != nil {
				if os.IsNotExist(err) {
					return config, nil // Return default config if no config file exists
				}
				return config, fmt.Errorf("error reading local config file: %w", err)
			}
		} else {
			return config, fmt.Errorf("error reading config file: %w", err)
		}
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}

	// Use default values if not specified in the config file
	if config.HeaderMaxLength == 0 {
		config.HeaderMaxLength = defaultConfig.HeaderMaxLength
	}
	if config.BodyLineMaxLength == 0 {
		config.BodyLineMaxLength = defaultConfig.BodyLineMaxLength
	}
	if len(config.AllowedTypes) == 0 {
		config.AllowedTypes = defaultConfig.AllowedTypes
	}

	return config, nil
}

func isRuleEnabled(config Config, ruleName string) bool {
	for _, disabledRule := range config.DisabledRules {
		if disabledRule == ruleName {
			return false
		}
	}
	return true
}

func validateCommitMsg(msg string, config Config) ([]string, bool) {
	var errors []string
	needsBreakingChange := false

	msg = strings.TrimSpace(msg)
	if msg == "" {
		errors = append(errors, "Commit message is empty")
		return errors, needsBreakingChange
	}

	lines := strings.Split(msg, "\n")
	header := lines[0]
	headerParts := strings.SplitN(header, ": ", 2)

	// Rule: header-format
	if isRuleEnabled(config, "header-format") && !headerPattern.MatchString(header) {
		errors = append(errors, "Header must be in format: <type>[optional scope][!]: <description>")
	}

	// Rule: header-max-length
	if isRuleEnabled(config, "header-max-length") && len(header) > config.HeaderMaxLength {
		errors = append(errors, fmt.Sprintf("Header must not exceed %d characters", config.HeaderMaxLength))
	}

	// Rule: header-lowercase
	if isRuleEnabled(config, "header-lowercase") && strings.ToLower(header) != header {
		errors = append(errors, "Header (short description) must be all lowercase")
	}

	// Rule: type-enum
	if isRuleEnabled(config, "type-enum") && len(headerParts) > 0 {
		typeScope := strings.Split(headerParts[0], "(")
		commitType := strings.TrimSuffix(typeScope[0], "!") // Remove '!' if present
		if !contains(config.AllowedTypes, commitType) {
			errors = append(errors, fmt.Sprintf("Type '%s' is not allowed. Allowed types are: %s", commitType, strings.Join(config.AllowedTypes, ", ")))
		}
	}

	// Rule: type-case
	if isRuleEnabled(config, "type-case") && len(headerParts) > 0 {
		typeScope := strings.Split(headerParts[0], "(")
		commitType := typeScope[0]
		if commitType != strings.ToLower(commitType) {
			errors = append(errors, "Type must be in lowercase")
		}
	}

	// Rule: type-empty
	if isRuleEnabled(config, "type-empty") && len(headerParts) > 0 {
		typeScope := strings.Split(headerParts[0], "(")
		commitType := typeScope[0]
		if commitType == "" {
			errors = append(errors, "Type must not be empty")
		}
	}

	// Rule: scope-case
	if isRuleEnabled(config, "scope-case") && len(headerParts) > 0 {
		typeScope := strings.Split(headerParts[0], "(")
		if len(typeScope) > 1 {
			scope := strings.TrimRight(typeScope[1], ")")
			if scope != strings.ToLower(scope) {
				errors = append(errors, "Scope must be in lowercase")
			}
		}
	}

	// Rule: subject-empty
	if isRuleEnabled(config, "subject-empty") && len(headerParts) < 2 {
		errors = append(errors, "Subject must not be empty")
	}

	// Rule: description-case
	if isRuleEnabled(config, "description-case") && len(headerParts) == 2 && len(headerParts[1]) > 0 {
		firstChar := headerParts[1][0]
		if firstChar >= 'A' && firstChar <= 'Z' {
			errors = append(errors, "Description must start with lowercase")
		}
	}

	// Rule: body-line-max-length
	if isRuleEnabled(config, "body-line-max-length") {
		for i, line := range lines[1:] {
			if len(line) > config.BodyLineMaxLength {
				errors = append(errors, fmt.Sprintf("Body line %d exceeds %d characters", i+2, config.BodyLineMaxLength))
			}
		}
	}

	// Rule: footer-format
	if isRuleEnabled(config, "footer-format") {
		for i, line := range lines[1:] {
			if footerPattern.MatchString(line) && !breakingChangePattern.MatchString(line) {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 || len(strings.TrimSpace(parts[1])) == 0 {
					errors = append(errors, fmt.Sprintf("Footer line %d must be in format: <token>: <value>", i+2))
				}
			}
		}
	}

	// Rule: breaking-change
	if isRuleEnabled(config, "breaking-change") {
		if strings.Contains(header, "!") && !containsBreakingChange(lines[1:]) {
			needsBreakingChange = true
		}
	}

	return errors, needsBreakingChange
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func containsBreakingChange(footerLines []string) bool {
	for _, line := range footerLines {
		if breakingChangePattern.MatchString(line) {
			return true
		}
	}
	return false
}

func promptForBreakingChange() string {
	fmt.Print("Describe briefly the breaking change (leave empty for no description): ")
	var description string
	fmt.Scanln(&description)
	return description
}

func appendBreakingChange(msg, description string) string {
	lines := strings.Split(msg, "\n")
	breakingChangeFooter := "BREAKING CHANGE: " + description

	// Ensure there's an empty line before adding the breaking change
	if len(lines) > 0 && lines[len(lines)-1] != "" {
		lines = append(lines, "")
	}

	// Add the breaking change footer
	lines = append(lines, breakingChangeFooter)

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func readFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var output strings.Builder

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading from stdin: %w", err)
		}
		output.WriteString(input)
	}

	return output.String(), nil
}

func readCommitMsg(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading commit message file: %w", err)
	}
	return string(content), nil
}

func writeCommitMsg(filename, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing commit message file: %w", err)
	}
	return nil
}

func runGommit(pathGetter ConfigPathGetter) error {
	fmt.Printf("Gommit version: %s\n", version)

	isHook := len(os.Args) >= 2

	if err := checkAndUpdate(isHook); err != nil {
		fmt.Printf("Warning: Failed to check for updates: %v\n", err)
	}

	config, err := loadConfig(pathGetter)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	var commitMsgFile string

	if !isHook {
		fmt.Println(headerStyle.Render("No file provided. Running in test mode."))
		tempFile, err := os.CreateTemp("", "COMMIT_EDITMSG")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tempFile.Name()) // Clean up the temp file when done
		commitMsgFile = tempFile.Name()

		fmt.Println(headerStyle.Render("Enter your commit message (Ctrl+D when finished):"))
		commitMsg, err := readFromStdin()
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		if err := writeCommitMsg(commitMsgFile, commitMsg); err != nil {
			return fmt.Errorf("failed to write test commit message: %w", err)
		}
	} else {
		commitMsgFile = os.Args[1]
	}

	originalMsg, err := readCommitMsg(commitMsgFile)
	if err != nil {
		return fmt.Errorf("failed to read commit message file: %w", err)
	}

	commitMsg := originalMsg
	errors, needsBreakingChange := validateCommitMsg(commitMsg, config)

	if needsBreakingChange && isRuleEnabled(config, AUTO_BREAKING_CHANGE) {
		description := promptForBreakingChange()
		commitMsg = appendBreakingChange(commitMsg, description)
		errors, _ = validateCommitMsg(commitMsg, config) // Revalidate after adding BREAKING CHANGE
	}

	if len(errors) > 0 {
		fmt.Println(errorStyle.Render(COMMIT_MSG_INVALID_MSG))
		fmt.Println(errorStyle.Render("Commit message does not follow the configured rules:"))
		for _, err := range errors {
			fmt.Println(detailStyle.Render(fmt.Sprintf("  • %s", err)))
		}
		fmt.Println(headerStyle.Render("Please edit your commit message to follow the rules:"))

		p := tea.NewProgram(initialModel(commitMsg))
		m, err := p.Run()
		if err != nil {
			return fmt.Errorf("error running text input program: %w", err)
		}

		commitMsg = m.(model).textInput.Value()
		if commitMsg == originalMsg {
			fmt.Println(errorStyle.Render(COMMIT_MSG_INVALID_MSG))
			fmt.Println(errorStyle.Render("Commit message was not modified."))
			return fmt.Errorf("commit message validation failed")
		}

		// Re-validate the edited commit message
		errors, needsBreakingChange = validateCommitMsg(commitMsg, config)

		if needsBreakingChange && isRuleEnabled(config, AUTO_BREAKING_CHANGE) {
			description := promptForBreakingChange()
			commitMsg = appendBreakingChange(commitMsg, description)
			errors, _ = validateCommitMsg(commitMsg, config) // Revalidate after adding BREAKING CHANGE
		}

		if len(errors) > 0 {
			fmt.Println(errorStyle.Render(COMMIT_MSG_INVALID_MSG))
			fmt.Println(errorStyle.Render("Commit message is still invalid:"))
			for _, err := range errors {
				fmt.Println(detailStyle.Render(fmt.Sprintf("  • %s", err)))
			}
			return fmt.Errorf("commit message validation failed")
		}
	}

	err = writeCommitMsg(commitMsgFile, commitMsg)
	if err != nil {
		return fmt.Errorf("failed to write commit message: %w", err)
	}

	fmt.Println(successStyle.Render("✔ Commit message is valid."))
	fmt.Println(headerStyle.Render("Final commit message:"))
	fmt.Println(detailStyle.Render(commitMsg))
	fmt.Println(successStyle.Render(successArt))

	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("gommit encountered an unexpected error: %v", r)))
			fmt.Println(errorStyle.Render(failureArt))
			os.Exit(1)
		}
	}()

	err := runGommit(DefaultConfigPathGetter{})
	if err != nil {
		fmt.Println(errorStyle.Render(COMMIT_MSG_INVALID_MSG))
		fmt.Println(errorStyle.Render(err.Error()))
		fmt.Println(errorStyle.Render(failureArt))
		os.Exit(1)
	}
}
