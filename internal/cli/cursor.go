package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/cursor"
)

func NewCursorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cursor",
		Short: "Cursor IDE integration with AI guard rails",
		Long: `Integrates AI Development Framework with Cursor IDE:

- Configures MCP (Model Context Protocol) servers
- Sets up domain-specific AI rules  
- Enables guard rails in Cursor AI assistant
- Provides real-time code validation
- Integrates domain expertise into code completion

This ensures Cursor's AI follows your project's guard rails and domain expertise.`,
	}

	cmd.AddCommand(newCursorSetupCmd())
	cmd.AddCommand(newCursorValidateCmd())
	cmd.AddCommand(newCursorRulesCmd())

	return cmd
}

func newCursorSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Setup Cursor IDE integration",
		Long: `Sets up Cursor IDE with AI Development Framework integration:

- Creates .cursor/ configuration directory
- Configures MCP servers for guard rails
- Sets up domain-specific AI rules
- Generates .cursorrules file
- Configures AI assistant settings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Blue("🎯 Setting up Cursor IDE integration...")

			// Check if .codekeeper exists
			config, err := loadProjectConfig()
			if err != nil {
				return fmt.Errorf("no AI dev project found. Run 'codekeeper init' first")
			}

			// Generate Cursor configuration
			cursorConfig, err := cursor.GenerateCursorConfig(".", config.Domain, config.GuardRails)
			if err != nil {
				return fmt.Errorf("failed to generate Cursor config: %w", err)
			}

			// Save Cursor configuration
			if err := cursor.SaveCursorConfig(".", cursorConfig); err != nil {
				return fmt.Errorf("failed to save Cursor config: %w", err)
			}

			color.Green("✅ Cursor IDE integration configured!")
			fmt.Printf("\n📋 Next steps:\n")
			fmt.Printf("1. Restart Cursor IDE\n")
			fmt.Printf("2. Verify MCP servers are loaded in Cursor settings\n")
			fmt.Printf("3. Test AI assistant with guard rails: Cmd+K\n")
			fmt.Printf("4. Check .cursorrules is being applied\n")

			return nil
		},
	}
}

func newCursorValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate Cursor IDE configuration",
		Long:  "Validates that Cursor IDE is properly configured with guard rails",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validateCursorConfiguration()
		},
	}
}

func newCursorRulesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: "Manage Cursor AI rules",
		Long:  "View and manage AI rules for Cursor IDE",
		RunE: func(cmd *cobra.Command, args []string) error {
			return manageCursorRules()
		},
	}
}

type ProjectConfig struct {
	Domain     string   `json:"domain"`
	GuardRails []string `json:"guard_rails"`
}

func loadProjectConfig() (*ProjectConfig, error) {
	// Check if .codekeeper directory exists
	if _, err := os.Stat(".codekeeper"); os.IsNotExist(err) {
		return nil, fmt.Errorf("no AI dev project found")
	}

	// Try to load from env.local first
	if data, err := os.ReadFile(".codekeeper/env.local"); err == nil {
		content := string(data)
		config := &ProjectConfig{
			Domain: "general",
			GuardRails: []string{"input_validation", "security_headers", "error_handling"},
		}

		// Extract domain
		if strings.Contains(content, "CODEKEEPER_DOMAIN=fintech") {
			config.Domain = "fintech"
			config.GuardRails = []string{"decimal_arithmetic", "audit_trails", "encryption_at_rest", "input_validation"}
		} else if strings.Contains(content, "CODEKEEPER_DOMAIN=healthcare") {
			config.Domain = "healthcare"
			config.GuardRails = []string{"hipaa_compliance", "data_encryption", "audit_logs", "access_controls"}
		} else if strings.Contains(content, "CODEKEEPER_DOMAIN=ecommerce") {
			config.Domain = "ecommerce"
			config.GuardRails = []string{"payment_security", "pci_compliance", "inventory_validation", "user_privacy"}
		}

		return config, nil
	}

	// Default config if no env file
	return &ProjectConfig{
		Domain:     "general",
		GuardRails: []string{"input_validation", "security_headers", "error_handling"},
	}, nil
}

func validateCursorConfiguration() error {
	color.Blue("🔍 Validating Cursor IDE configuration...")
	
	checks := []struct {
		name string
		fn   func() (bool, string)
	}{
		{"Cursor configuration directory", checkCursorDirectory},
		{"Cursor settings file", checkCursorSettings},
		{"Cursor rules (.mdc files)", checkCursorRules},
		{"MCP server configuration", checkMCPConfig},
		{"Guard rails integration", checkGuardRailsIntegration},
	}

	allPassed := true
	for _, check := range checks {
		passed, message := check.fn()
		if passed {
			color.Green("✓ %s", check.name)
		} else {
			color.Red("✗ %s: %s", check.name, message)
			allPassed = false
		}
	}

	fmt.Println()
	if allPassed {
		color.Green("🎉 Cursor IDE is properly configured!")
		fmt.Printf("\n📋 Integration status:\n")
		fmt.Printf("  • Guard rails: Active\n")
		fmt.Printf("  • MCP servers: Connected\n")
		fmt.Printf("  • Domain rules: Applied\n")
		fmt.Printf("  • AI assistant: Enhanced\n")
	} else {
		color.Yellow("⚠️  Some issues found. Run 'codekeeper cursor setup' to fix.")
		return fmt.Errorf("cursor configuration issues detected")
	}

	return nil
}

func manageCursorRules() error {
	color.Blue("📋 Cursor AI Rules:")
	
	config, err := loadProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to load project configuration: %w", err)
	}

	fmt.Println()
	color.Cyan("🎯 Domain: %s", config.Domain)
	
	fmt.Println()
	color.Cyan("🛡️ Active Guard Rails:")
	for _, rule := range config.GuardRails {
		fmt.Printf("  ✓ %s\n", formatRuleName(rule))
	}

	// Load custom rules from .cursorrules file
	if customRules, err := loadCustomRules(); err == nil && len(customRules) > 0 {
		fmt.Println()
		color.Cyan("📝 Custom Rules:")
		for _, rule := range customRules {
			fmt.Printf("  • %s\n", rule)
		}
	}

	// Show enforcement status
	fmt.Println()
	color.Cyan("⚡ Enforcement Status:")
	if cursorConfigExists() {
		color.Green("  ✓ Cursor IDE integration: Active")
		color.Green("  ✓ Real-time validation: Enabled")
		color.Green("  ✓ AI assistant enhancement: Active")
	} else {
		color.Yellow("  ⚠ Cursor IDE integration: Not configured")
		fmt.Printf("    Run 'codekeeper cursor setup' to enable\n")
	}

	fmt.Println()
	color.Blue("💡 Commands:")
	fmt.Printf("  codekeeper cursor setup     - Setup/update Cursor integration\n")
	fmt.Printf("  codekeeper cursor validate  - Validate configuration\n")
	fmt.Printf("  codekeeper check           - Run guard rails validation\n")

	return nil
}

func checkCursorDirectory() (bool, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, "cannot access home directory"
	}
	
	cursorDir := filepath.Join(homeDir, ".cursor")
	if _, err := os.Stat(cursorDir); os.IsNotExist(err) {
		return false, "Cursor IDE not installed"
	}
	return true, ""
}

func checkCursorSettings() (bool, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, "cannot access home directory"
	}
	
	settingsPath := filepath.Join(homeDir, ".cursor", "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return false, "Cursor settings.json not found"
	}
	
	// Check if MCP is configured in settings
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return false, "cannot read settings file"
	}
	
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return false, "invalid settings JSON"
	}
	
	// Check for MCP configuration
	if mcpServers, exists := settings["mcp.servers"]; exists {
		if servers, ok := mcpServers.(map[string]interface{}); ok && len(servers) > 0 {
			return true, ""
		}
	}
	
	return false, "MCP servers not configured"
}

func checkCursorRules() (bool, string) {
	// Check for new .cursor/rules directory structure
	if _, err := os.Stat(".cursor/rules"); os.IsNotExist(err) {
		return false, "run 'codekeeper cursor setup' to create .cursor/rules"
	}
	
	// Check for key rule files
	requiredRules := []string{
		".cursor/rules/project-standards.mdc",
		".cursor/rules/security-standards.mdc",
	}
	
	for _, ruleFile := range requiredRules {
		if _, err := os.Stat(ruleFile); os.IsNotExist(err) {
			return false, fmt.Sprintf("missing rule file: %s", ruleFile)
		}
		
		// Verify it contains proper MDC format
		data, err := os.ReadFile(ruleFile)
		if err != nil {
			return false, fmt.Sprintf("cannot read %s", ruleFile)
		}
		
		content := string(data)
		if !strings.Contains(content, "---") || !strings.Contains(content, "description:") {
			return false, fmt.Sprintf("invalid MDC format in %s", ruleFile)
		}
	}
	
	return true, ""
}

func checkMCPConfig() (bool, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, "cannot access home directory"
	}
	
	settingsPath := filepath.Join(homeDir, ".cursor", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return false, "cannot read Cursor settings"
	}
	
	// Check for CodeKeeper MCP server
	if strings.Contains(string(data), "codekeeper-guard-rails") {
		return true, ""
	}
	
	return false, "CodeKeeper MCP server not configured"
}

func checkGuardRailsIntegration() (bool, string) {
	config, err := loadProjectConfig()
	if err != nil {
		return false, "no project configuration found"
	}
	
	if len(config.GuardRails) == 0 {
		return false, "no guard rails configured"
	}
	
	// Check if .cursor/rules reflect current configuration
	if data, err := os.ReadFile(".cursor/rules/project-standards.mdc"); err == nil {
		content := string(data)
		// Check if domain is reflected in rules
		if strings.Contains(content, config.Domain) {
			return true, ""
		}
		return false, "domain not reflected in cursor rules"
	}
	
	return false, ".cursor/rules files missing"
}

func loadCustomRules() ([]string, error) {
	var customRules []string
	
	// Read from .cursor/rules directory
	rulesDir := ".cursor/rules"
	entries, err := os.ReadDir(rulesDir)
	if err != nil {
		return nil, err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mdc") {
			data, err := os.ReadFile(filepath.Join(rulesDir, entry.Name()))
			if err != nil {
				continue
			}
			
			content := string(data)
			// Extract description from MDC metadata
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "description:") {
					desc := strings.TrimPrefix(line, "description:")
					desc = strings.Trim(desc, " \"")
					customRules = append(customRules, desc)
					break
				}
			}
		}
	}
	
	return customRules, nil
}

func cursorConfigExists() bool {
	homeDir, _ := os.UserHomeDir()
	settingsPath := filepath.Join(homeDir, ".cursor", "settings.json")
	
	if data, err := os.ReadFile(settingsPath); err == nil {
		return strings.Contains(string(data), "codekeeper")
	}
	
	return false
}

func formatRuleName(rule string) string {
	// Convert snake_case to human readable
	formatted := strings.ReplaceAll(rule, "_", " ")
	words := strings.Fields(formatted)
	
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	
	return strings.Join(words, " ")
}