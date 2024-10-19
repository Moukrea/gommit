package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
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

const (
	currentVersion      = "v1.0.0" // Update this with each release
	updateCheckInterval = 24 * time.Hour
	updateCheckFile     = ".gommit_last_update"
	gommitRepo          = "Moukrea/gommit"
)

type Rule struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Config struct {
	DisabledRules []string `yaml:"disabled_rules"`
}

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var defaultRules = []Rule{
	{Name: "header-format", Description: "Header must be in format: <type>[optional scope][!]: <description>"},
	{Name: "header-max-length", Description: "Header must not exceed 100 characters"},
	{Name: "header-lowercase", Description: "Header (short description) must be all lowercase"},
	{Name: "description-case", Description: "Description must start with lowercase"},
	{Name: "body-line-max-length", Description: "Body lines must not exceed 100 characters"},
	{Name: "footer-format", Description: "Footer must be in format: <token>: <value>"},
	{Name: "breaking-change", Description: "Breaking changes must be indicated in footer"},
	{Name: "auto-breaking-change", Description: "Automatically add BREAKING CHANGE to footer when '!' is present in header"},
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

func getExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %w", err)
	}
	return filepath.Dir(ex), nil
}

func loadConfig() (Config, error) {
	execPath, err := getExecutablePath()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(execPath, "gommit.yml")
	config := Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Return default config if file doesn't exist
		}
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
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
	// ... [rest of the validation logic remains unchanged]
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
	// ... [rest of the function remains unchanged]
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

func shouldCheckForUpdate() bool {
	info, err := os.Stat(updateCheckFile)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Error checking update file: %v\n", err)
		return false
	}
	return time.Since(info.ModTime()) > updateCheckInterval
}

func checkAndUpdate() error {
	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	if release.TagName > currentVersion {
		fmt.Printf("A new version of Gommit is available: %s\n", release.TagName)
		fmt.Print("Do you want to update? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			if err := performUpdate(release); err != nil {
				return fmt.Errorf("failed to perform update: %w", err)
			}
		}
	}

	// Update the last check time
	if err := os.WriteFile(updateCheckFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to update check file: %w", err)
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

func runGommit() error {
	// Check for updates before proceeding
	if shouldCheckForUpdate() {
		if err := checkAndUpdate(); err != nil {
			fmt.Printf("Error checking for updates: %v\n", err)
			// Continue with normal operation even if update check fails
		}
	}

	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	var commitMsgFile string
	var isTestMode bool

	if len(os.Args) < 2 {
		fmt.Println(headerStyle.Render("No file provided. Running in test mode."))
		tempFile, err := os.CreateTemp("", "COMMIT_EDITMSG")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tempFile.Name()) // Clean up the temp file when done
		commitMsgFile = tempFile.Name()
		isTestMode = true

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
	
	if needsBreakingChange && isRuleEnabled(config, "auto-breaking-change") {
		description := promptForBreakingChange()
		commitMsg = appendBreakingChange(commitMsg, description)
		errors, _ = validateCommitMsg(commitMsg, config)  // Revalidate after adding BREAKING CHANGE
	}

	if len(errors) > 0 {
		fmt.Println(errorStyle.Render("Commit message does not follow the configured rules:"))
		for _, err := range errors {
			fmt.Println(detailStyle.Render(fmt.Sprintf("  • %s", err)))
		}
		fmt.Println()
		fmt.Println(headerStyle.Render("Please edit your commit message to follow the rules:"))
		fmt.Println()

		p := tea.NewProgram(initialModel(commitMsg))
		m, err := p.Run()
		if err != nil {
			return fmt.Errorf("error running text input program: %w", err)
		}

		commitMsg = m.(model).textInput.Value()
		if commitMsg == originalMsg {
			fmt.Println(errorStyle.Render(failureArt))
			return fmt.Errorf("commit message was not modified")
		}
	}

	if commitMsg != originalMsg {
		err = writeCommitMsg(commitMsgFile, commitMsg)
		if err != nil {
			return fmt.Errorf("failed to write commit message: %w", err)
		}
	}

	fmt.Println(successStyle.Render("✔ Commit message is valid."))
	fmt.Println(headerStyle.Render("Final commit message:"))
	fmt.Println(detailStyle.Render(commitMsg))
	fmt.Println(successStyle.Render(successArt))

	return nil
}

const successArt = `
                                                                                                    
                                                                                                    
                                                            ██▓██                                   
                                                        ████▒░░░▒▓█████                             
                                              ██    ██▒░▓█▓░░░░░░░░░▒▒██                            
                                            ████  ██░░░░░░░░░░░░░░░░░░░▓█                           
                                           ██░░▓█▓░░░░░░░░░░░▓▒░░░░░░░░░▒██                         
                                           █▓░░░░░░░▓▓▒▓▓▓██▓▒▓█▓░░░░░░░░░▒██                       
                                           ██░░░░░▒▓░░░░░░░░░░░░▓█▒▓█████▒░▒██                      
                                           █▓▒▓▒▓█░░░░░░░░░░░░░░░░░░░░░░██▒░▓██                     
                                           █▒░▓▒▒▓▒░░░░░░░░░░▒█▓░░░░░░░░█▓░░░██                     
                                           ███▒░░░░░░░░░░░░░░░░░▓░░░░░▓▓▒░░░░██                     
                                            █▓░░░▒░░░░░░░░░░░░░░░░░░░░▓▒▓▓░░▒██                     
                                           ██░░░██░░░░▓▓░░░▓███▓░░░░░░▓▒█▒▓░▓██                     
              ███                          █▒░░░▒▓░░░▓▓░░░░░░░▒▓▒░░░░░▒▒▓▓▒███                      
           ██▓░░▓█                         █▒░░░░░░░█▓░░░░░░░░░░░░░░░░░▓▒▓███                       
          █▒░░░░▒█                         █░░░░░░░▓█▒░░░░░░░░░░░░░░░░▒▓░░███                       
          ▓░░░░░▓█                         █░░░░░░░░▒▓░░░░░░▒░░░░░░░░░░░░░▓██                       
         █▓░░░░░▓                          █▒░▒█░░░░░░░░░░░▓██▓░░░░░░░░░░░░██                       
          █░░░░░▓                          █▓░▓█░▒▓▓▓▓▓▓▓░░░░▓█▒░░░░░░░░░░███                       
          █▓░░░░▒█                         ██▒░▓░▓█▒░░░░▒██▒░░█░░░░░░░▓████                         
     ███████▓▒░░░░▒██                       ██▒░░░░░▒▓▓▒░░░░░░░░░░░░░███                            
    ██░░░░░░░░▒█▓░░░▓█                       ██░░░░░░▓▓░░░░░░░░░░░░▓██                              
   ██░░░░░░░░░░░░█▒░░▓████                    ███░░░░░░░░░░░░░░░░▓██▒▓███                           
    █▓░░░░░░░░░░░▓▓░░░█████████████████████████████░░░░░░░░░░░░░░░▓▒░░░▒█▓███                       
     █▓▒▒░░▒▒▒▓███▒░░░██▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓█░▒▓█▓░░░░░░░░░▒█▒░░░░▓█▓▓▓▓███                    
    ██░░░░░░░░░░░▓▓░░▒█░▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓█░░░▒███▓▒░▒▓██▒░░░░░▓█▓▓▓▓▓▓▓▓███                 
    ██▒░░░░░░░░░░░█▒▒█▓░▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█▒░░░░░▒▒▒▒▒░░░░░░░▒██▓▓▓▓▓▓▓▓▓▓▓▓██               
      █▓░░░░▒▓████▒░░▓▒░▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█▓░░░░░░░░░░░░░░▒██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██             
     ██░░░░░░░░░░█▓░▒▓░▒▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█░░░░░░░░░░▒▓██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██           
      ██▓▒▒░░░▒▒▓█▒▒▓░▒█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█▒▒▒▒▒▒▒▒▓▓░░░░░░░▓█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███         
        ███▒░░░░▒▓█▒░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██████▒▒▒▒▒▒▒▒█▒░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██        
            ███████████████████████     ██▒▒▒▒▒▒▒▒█░░░░░░░▒█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█       
                                        ██▒▒▒▒▒▒▒▓█░░░░░░░▒█▓▓▓▓▓▓▓▓▓▓▓█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██     
                                        ██▒▒▒▒▒▒▒█▓░░░░░░░▒█▓▓▓▓▓▓▓▓▓▓▓▓▓▓███▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██    
                                        ██▒▒▒▒▒▒▒█▒░░░░░░░▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█████▓▓▓▓▓▓▓▓▓▓▓▓▓██    
                                        ██▒▒▒▒▒▒▒█▒░░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██   ██▓▓▓▓▓▓▓▓▓▓▓██    
                                        ██▓▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███  ██▓▓▓▓▓▓▓▓▓▓▓██     
                                         ██▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██  ██▓▓▓▓▓▓▓▓▓▓▓▓██     
                                         ██▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██ ██▓▓▓▓▓▓▓▓▓▓▓▓██      
                                         ███▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███ █▓▓▓▓▓▓▓▓▓▓▓▓██       
                                          █▒▒██▓▒█▒░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████▓▓▓▓▓▓▓▓▓▓▓▓██        
                                          █▓░░░░▓█░░░░░░░░▓▓▓▓▓▓▓▓▓▓▓████▓████▓▓▓▓▓▓▓▓▓▓▓██         
                                          █▓░░░░░░░░░░░░░░░▓████▓▓▓▒▒░░░░░█████▓▓▓▓▓▓▓▓▓██          
                                           █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██▒░▒███▓▓▓▓▓██           
                                           █▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓█░░░░░▒▓▓████            
                                           ███▒░░░░░░░░░░░░░░░░░░░░░░░░░░▒██░░░░░░░█▓▒██            
                                           ██▒▒▓▓██▓▒░░░░░░░░░░░░▒▒▓▓███████▒░░░░░▓▓░░▒██           
                                           ██▒▒▒▒▓▓▓▓▓███████████▓▓▓▓▓▓▓▓▓▓██░░░░░░▒█░▓█            
                                           ██▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██▓▓▒▒▒▒▒▓██            
                                           ██▓▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████                    
                                            ██▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███                    
                                             █▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███                   
                                                  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                         
                                                          ▓▓▓▓▓▓▓                                   
                                                                                                    
`

const failureArt = `
                                                                                                    
                                                                                                    
                                                                                                    
                                                             ████                                   
                                                           ██▒░░▒▒███                               
                                               ██    ██▓▓██▓░░░░░░░░▒▒▓█                            
                                             ███   ██░░░░░░░░░░░░░░░░░░▒▓                           
                                            █▓░▒██▓▒░░░░░░░░░░▒░░░░░░░░░░▓██                        
                                            █▒░░░░░░░▓▓▓▓▒▒▓█▓███▒░░░░░░░░░▓█                       
                                            █▓░░░░░▒▓░░░░░░░░░░░▒▓▓▒▒▓███▓▒░▒█                      
                                            █▓▓▒▒▒▓░░░░░░░░░░░░░░░░░░░░░░█▓░░▓█                     
                                            █░░▓▒░░░░░░░░░░░░░░░░░░░░░░░░█▒░░▒██                    
                                            █▓█▒░░░░░░░▒░░░░░░░░░░░░░░░▓▓░░░░▒██                    
                                             █▒░▒▒▒▒▒▒▒█▒▒▒▒▒▒▒▒▒░░░░░░▓▒█▒░░▒██                    
                                            █▓░░░██░░░░▓▒░░░░██░░░░░░░░▒▒▓▒▒░▓█                     
                                            █░░░░▓▓░░░▓▓░░░░░█▒░░▒░░░░░░▓▒▓▒▒██                     
                                           ██░░░░░░░▒█▒░░░░░░░░░▒░░░░░░▒▓░▓▓██                      
                                           ██░░░░░░░██░░░░░░░░░░░░░░░░░█▒░░██                       
                                           ██░░░░░░░░▓▓░░░░░░░░░░░░░░░░░░░░██                       
                                           ██░░░████████████▓▒░░░░░░░░░░░░░▒██                      
                                           ██░█▓░░░▒████▓░░░░░█▒░░░░░░░░░░░██                       
          ████                             ██▓░░░░░░░░░░░░░░░░░░░░░░░░░▒████                        
      █▓░░░░░░░░▓█                           ██░░░░░░░░░░░░░░░░░░░░░░░██                            
    █▓░░░░░░░░░░░░▒█                          ██░░░░░░░░░░░░░░░░░░░░▓██                             
    █▒░░░░░░░░░░░░░▒█████                      ██▓░░░░░░░░░░░░░░░░▓██▒▒██                           
    ██▓▓▒▒░░░░░░░░░░▓█▓▒▓▓▓██████████████████▓▓▓▓▓█▓░░░░░░░░░░░░▒▒░▓░░░░▒█▓██                       
   █░░░░░░░░░░░░░░░░░▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓░▒▓▓▒░░░░░░░░░▓█░░░░░█▓▓▓▓▓███                   
  ██░░░░░░░░░░░░░░░░░█▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓░░░▒███▓▒▒▓▓█▓░░░░░░██▓▓▓▓▓▓▓▓██                 
   ███▓▓▒░░░░░░░░░░░░▓▒▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█░░░░░░░▒▒▒░░░░░░░░▒█▓▓▓▓▓▓▓▓▓▓▓▓▓██              
  █▓░░░░░░░░░░░░░░░░▒▒▒▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█▒░░░░░░░░░░░░░░▓█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██            
  █▒░░░░░░░░░░░░░░░▓▓░▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓░░░░░░░░░▒▒███▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██          
   ██▒▒░░░░▒░░░░░░█▒░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█▓▒▒▒▒▒▒▒▒▓▒░░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██         
         █░░░░░▒█▓░▒█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███▓▒▒▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██       
        █▒░░░░▓▓░▒█▓▓▓▓▓▓▓▓▓▓▓██████     █▓▒▒▒▒▒▒▒▒▓░░░░░░░▒█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█      
        █░░░░░▓     ████                 █▓▒▒▒▒▒▒▒▓▒░░░░░░░▓█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█     
        █░░░░░▓                          █▓▒▒▒▒▒▒▒▓░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██   
        █░░░░░░█                         █▓▒▒▒▒▒▒▒▓░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██   
        ██░░░░░▓█                        █▓▒▒▒▒▒▒▒▓░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██   ██▓▓▓▓▓▓▓▓▓▓▓██   
          ██████                          ▓▒▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██   ██▓▓▓▓▓▓▓▓▓▓▓██    
                                          ▓▒▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██  ██▓▓▓▓▓▓▓▓▓▓▓▓█     
                                          █▒▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███ ██▓▓▓▓▓▓▓▓▓▓▓▓██     
                                          ██▒▒▒▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██ ██▓▓▓▓▓▓▓▓▓▓▓▓██      
                                          ██▒██▒▒▒█░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████▓▓▓▓▓▓▓▓▓▓▓▓██       
                                           █░░░░▓██░░░░░░░░█▓▓▓▓▓▓▓▓▓▓▓█████▓█▓▓▓▓▓▓▓▓▓▓▓▓██        
                                           █▒░░░░░░░░░░░░░░▓███████▓▒▒▒░░░░█████▓▓▓▓▓▓▓▓▓██         
                                           █▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██▒░▓██▓▓▓▓▓███          
                                           ██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█▓░░░░▒▓█████            
                                           ███▓▒░░░░░░░░░░░░░░░░░░░░░░░░░░▒█▓░░░░░░░█▒▒█            
                                            █▓▒▒▓██▓▒░░░░░░░░░░░░░░▒▓▓███████░░░░░░█░░░▓█           
                                            ██▒▒▒▒▓▓▓▓▓███████████▓▓▓▓▓▓▓▓▓▓█▓░░░░░░██░▓█           
                                            ██▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█▓▒▒▒▒▒▒▒▓█            
                                            ██▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███  ███               
                                            ██▓▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███                   
                                              ▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                     
                                                   ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                          
                                                                                                    
                                                                                                    
                                                                                                    
                                                                                                    
                                                                                                    
`

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("gommit encountered an unexpected error: %v", r)))
			fmt.Println(errorStyle.Render(failureArt))
			os.Exit(1)
		}
	}()

	err := runGommit()
	if err != nil {
		fmt.Println(errorStyle.Render(err.Error()))
		fmt.Println(errorStyle.Render(failureArt))
		os.Exit(1)
	}
}