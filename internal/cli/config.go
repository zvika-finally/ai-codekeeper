package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage framework configuration",
		Long: `Manage AI development framework configuration:

- View current settings
- Set preferences and API keys
- Reset to defaults
- Export/import configuration`,
	}

	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigResetCmd())

	return cmd
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show current configuration",
		Long:  "Displays current framework configuration and preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listConfiguration()
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set configuration value",
		Long:  "Sets a configuration key to the specified value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]
			
			return setConfiguration(key, value)
		},
	}
}

func newConfigResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  "Resets all configuration to default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			return resetConfiguration()
		},
	}
}

type CodeKeeperConfig struct {
	Domain          string            `json:"domain"`
	IDE             string            `json:"ide"`
	GuardRails      GuardRailsSettings `json:"guard_rails"`
	AIModel         AIModelSettings   `json:"ai_model"`
	Environment     EnvSettings       `json:"environment"`
	LastUpdated     string            `json:"last_updated"`
}

type GuardRailsSettings struct {
	Enforcement string   `json:"enforcement"` // "strict", "advisory", "disabled"
	PreCommit   bool     `json:"pre_commit"`
	CI          bool     `json:"ci"`
	CustomRules []string `json:"custom_rules"`
}

type AIModelSettings struct {
	Provider    string `json:"provider"`    // "openai", "anthropic", "local"
	Model       string `json:"model"`       // "gpt-4", "claude-3", etc.
	Temperature float64 `json:"temperature"`
	MaxTokens   int    `json:"max_tokens"`
}

type EnvSettings struct {
	DevContainer bool              `json:"dev_container"`
	AutoSetup    bool              `json:"auto_setup"`
	Variables    map[string]string `json:"variables"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".codekeeper", "config.json")
}

func getDefaultConfig() *CodeKeeperConfig {
	return &CodeKeeperConfig{
		Domain: "general",
		IDE:    "vscode",
		GuardRails: GuardRailsSettings{
			Enforcement: "advisory",
			PreCommit:   true,
			CI:          true,
			CustomRules: []string{},
		},
		AIModel: AIModelSettings{
			Provider:    "auto",
			Model:       "auto",
			Temperature: 0.7,
			MaxTokens:   2048,
		},
		Environment: EnvSettings{
			DevContainer: true,
			AutoSetup:    true,
			Variables:    make(map[string]string),
		},
	}
}

func loadConfiguration() (*CodeKeeperConfig, error) {
	configPath := getConfigPath()
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// If config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := getDefaultConfig()
		if err := saveConfiguration(defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return defaultConfig, nil
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config CodeKeeperConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

func saveConfiguration(config *CodeKeeperConfig) error {
	configPath := getConfigPath()
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func listConfiguration() error {
	config, err := loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	color.Blue("📋 Current Configuration:")
	fmt.Println()
	
	color.Cyan("🎯 Domain Settings:")
	fmt.Printf("  Domain: %s\n", config.Domain)
	fmt.Printf("  IDE: %s\n", config.IDE)
	fmt.Println()
	
	color.Cyan("🛡️ Guard Rails:")
	fmt.Printf("  Enforcement: %s\n", config.GuardRails.Enforcement)
	fmt.Printf("  Pre-commit hooks: %t\n", config.GuardRails.PreCommit)
	fmt.Printf("  CI/CD integration: %t\n", config.GuardRails.CI)
	if len(config.GuardRails.CustomRules) > 0 {
		fmt.Printf("  Custom rules: %s\n", strings.Join(config.GuardRails.CustomRules, ", "))
	} else {
		fmt.Printf("  Custom rules: none\n")
	}
	fmt.Println()
	
	color.Cyan("🤖 AI Model:")
	fmt.Printf("  Provider: %s\n", config.AIModel.Provider)
	fmt.Printf("  Model: %s\n", config.AIModel.Model)
	fmt.Printf("  Temperature: %.1f\n", config.AIModel.Temperature)
	fmt.Printf("  Max tokens: %d\n", config.AIModel.MaxTokens)
	fmt.Println()
	
	color.Cyan("🔧 Environment:")
	fmt.Printf("  DevContainer: %t\n", config.Environment.DevContainer)
	fmt.Printf("  Auto setup: %t\n", config.Environment.AutoSetup)
	if len(config.Environment.Variables) > 0 {
		fmt.Printf("  Variables: %d defined\n", len(config.Environment.Variables))
	} else {
		fmt.Printf("  Variables: none\n")
	}
	fmt.Println()
	
	color.Yellow("📁 Config file: %s", getConfigPath())
	
	return nil
}

func setConfiguration(key, value string) error {
	config, err := loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Parse key path (e.g., "guard_rails.enforcement")
	parts := strings.Split(key, ".")
	
	switch parts[0] {
	case "domain":
		config.Domain = value
	case "ide":
		config.IDE = value
	case "guard_rails":
		if len(parts) < 2 {
			return fmt.Errorf("guard_rails setting requires a sub-key (e.g., guard_rails.enforcement)")
		}
		switch parts[1] {
		case "enforcement":
			if value != "strict" && value != "advisory" && value != "disabled" {
				return fmt.Errorf("enforcement must be 'strict', 'advisory', or 'disabled'")
			}
			config.GuardRails.Enforcement = value
		case "pre_commit":
			config.GuardRails.PreCommit = value == "true"
		case "ci":
			config.GuardRails.CI = value == "true"
		default:
			return fmt.Errorf("unknown guard_rails setting: %s", parts[1])
		}
	case "ai_model":
		if len(parts) < 2 {
			return fmt.Errorf("ai_model setting requires a sub-key (e.g., ai_model.provider)")
		}
		switch parts[1] {
		case "provider":
			config.AIModel.Provider = value
		case "model":
			config.AIModel.Model = value
		default:
			return fmt.Errorf("unknown ai_model setting: %s", parts[1])
		}
	case "environment":
		if len(parts) < 2 {
			return fmt.Errorf("environment setting requires a sub-key (e.g., environment.dev_container)")
		}
		switch parts[1] {
		case "dev_container":
			config.Environment.DevContainer = value == "true"
		case "auto_setup":
			config.Environment.AutoSetup = value == "true"
		default:
			// Treat as environment variable
			if len(parts) == 3 && parts[1] == "var" {
				config.Environment.Variables[parts[2]] = value
			} else {
				return fmt.Errorf("unknown environment setting: %s", parts[1])
			}
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", parts[0])
	}
	
	if err := saveConfiguration(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	
	color.Green("✅ Configuration updated: %s = %s", key, value)
	return nil
}

func resetConfiguration() error {
	color.Yellow("🔄 Resetting configuration to defaults...")
	
	defaultConfig := getDefaultConfig()
	if err := saveConfiguration(defaultConfig); err != nil {
		return fmt.Errorf("failed to reset configuration: %w", err)
	}
	
	color.Green("✅ Configuration reset to defaults")
	color.Blue("📋 Run 'codekeeper config list' to see default settings")
	
	return nil
}