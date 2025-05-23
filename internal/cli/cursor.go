package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/generator"
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

			// Create project spec for generator with AI-detected stack
			spec := &generator.ProjectSpec{
				Name:       filepath.Base(getCurrentWorkingDir()),
				Domain:     config.Domain,
				CoreEntity: detectCoreEntity(),
				Backend:    detectBackendStack(),
				APIStyle:   detectAPIStyle(),
			}

			// Generate Cursor configuration using new generator
			cursorIntegration := generator.NewCursorIntegration(spec)
			files, err := cursorIntegration.Generate()
			if err != nil {
				return fmt.Errorf("failed to generate Cursor config: %w", err)
			}

			// Write all generated files
			for filePath, content := range files {
				// Create directory if needed
				dir := filepath.Dir(filePath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", dir, err)
				}

				// Write file
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					return fmt.Errorf("failed to write file %s: %w", filePath, err)
				}
			}

			color.Green("✅ Cursor IDE integration configured!")
			fmt.Printf("\n📋 Next steps:\n")
			fmt.Printf("1. Restart Cursor IDE\n")
			fmt.Printf("2. Verify MCP servers are loaded in Cursor settings\n")
			fmt.Printf("3. Test AI assistant with guard rails: Cmd+K\n")
			fmt.Printf("4. Check .cursor/rules/*.mdc files are being applied\n")

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

	// Try to load from config.json first
	if data, err := os.ReadFile(".codekeeper/config.json"); err == nil {
		var projectConfig struct {
			Domain string `json:"domain"`
		}
		if json.Unmarshal(data, &projectConfig) == nil {
			config := &ProjectConfig{
				Domain: projectConfig.Domain,
			}
			
			// Set domain-specific guard rails
			switch projectConfig.Domain {
			case "fintech":
				config.GuardRails = []string{"decimal_arithmetic", "audit_trails", "encryption_at_rest", "input_validation"}
			case "healthcare":
				config.GuardRails = []string{"hipaa_compliance", "data_encryption", "audit_logs", "access_controls"}
			case "ecommerce":
				config.GuardRails = []string{"payment_security", "pci_compliance", "inventory_validation", "user_privacy"}
			default:
				config.GuardRails = []string{"input_validation", "security_headers", "error_handling"}
			}
			
			return config, nil
		}
	}

	// Try to load from env.local as fallback
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

	// Default config if no configuration files found
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
		{"MCP configuration (.cursor/mcp.json)", checkMCPConfig},
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
	// Check for project-level MCP configuration
	mcpConfigPath := ".cursor/mcp.json"
	if _, err := os.Stat(mcpConfigPath); os.IsNotExist(err) {
		return false, "run 'codekeeper cursor setup' to create .cursor/mcp.json"
	}
	
	// Verify it contains proper MCP server configuration
	data, err := os.ReadFile(mcpConfigPath)
	if err != nil {
		return false, "cannot read .cursor/mcp.json"
	}
	
	var mcpConfig map[string]interface{}
	if err := json.Unmarshal(data, &mcpConfig); err != nil {
		return false, "invalid JSON in .cursor/mcp.json"
	}
	
	// Check for mcpServers section
	if servers, exists := mcpConfig["mcpServers"]; exists {
		if serverMap, ok := servers.(map[string]interface{}); ok && len(serverMap) > 0 {
			// Check for CodeKeeper MCP servers
			for serverName := range serverMap {
				if strings.Contains(serverName, "codekeeper") {
					return true, ""
				}
			}
			return false, "no CodeKeeper MCP servers found"
		}
	}
	
	return false, "no MCP servers configured"
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
	// Check for project-level MCP configuration
	mcpConfigPath := ".cursor/mcp.json"
	if data, err := os.ReadFile(mcpConfigPath); err == nil {
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

func getCurrentWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "ai-project"
	}
	return wd
}

// detectBackendStack analyzes the project to determine the backend technology
func detectBackendStack() string {
	// Check for Go
	if _, err := os.Stat("go.mod"); err == nil {
		return "go"
	}
	if _, err := os.Stat("main.go"); err == nil {
		return "go"
	}
	
	// Check for Python
	if _, err := os.Stat("requirements.txt"); err == nil {
		return "python"
	}
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return "python"
	}
	if _, err := os.Stat("Pipfile"); err == nil {
		return "python"
	}
	
	// Check for Node.js/JavaScript
	if _, err := os.Stat("package.json"); err == nil {
		return "javascript"
	}
	
	// Check for Rust
	if _, err := os.Stat("Cargo.toml"); err == nil {
		return "rust"
	}
	
	// Check for Java
	if _, err := os.Stat("pom.xml"); err == nil {
		return "java"
	}
	if _, err := os.Stat("build.gradle"); err == nil {
		return "java"
	}
	
	// Check for C#/.NET
	if _, err := os.Stat("*.csproj"); err == nil {
		return "csharp"
	}
	
	// Check for PHP
	if _, err := os.Stat("composer.json"); err == nil {
		return "php"
	}
	
	// Default to JavaScript for web projects
	return "javascript"
}

// detectCoreEntity analyzes the project to determine the main business entity
func detectCoreEntity() string {
	// Check common domain-specific entities by examining files and directories
	commonEntities := map[string][]string{
		"User":        {"user", "account", "profile", "auth"},
		"Product":     {"product", "item", "catalog", "inventory"},
		"Order":       {"order", "purchase", "transaction", "sale"},
		"Payment":     {"payment", "billing", "invoice", "charge"},
		"Patient":     {"patient", "medical", "health", "clinical"},
		"Customer":    {"customer", "client", "member"},
		"Project":     {"project", "task", "workflow"},
		"Document":    {"document", "file", "content"},
		"Event":       {"event", "booking", "reservation"},
		"Asset":       {"asset", "resource", "property"},
	}
	
	// Scan directory names and file names for entity hints
	for entity, keywords := range commonEntities {
		for _, keyword := range keywords {
			// Check for directories
			if _, err := os.Stat(keyword); err == nil {
				return entity
			}
			// Check for files with keyword
			matches, _ := filepath.Glob("*" + keyword + "*")
			if len(matches) > 0 {
				return entity
			}
		}
	}
	
	// Default to generic entity
	return "Entity"
}

// detectAPIStyle analyzes the project to determine the API style
func detectAPIStyle() string {
	// Check for GraphQL
	if _, err := os.Stat("schema.graphql"); err == nil {
		return "GraphQL"
	}
	if _, err := os.Stat("*.graphql"); err == nil {
		return "GraphQL"
	}
	
	// Check for gRPC
	if _, err := os.Stat("*.proto"); err == nil {
		return "gRPC"
	}
	
	// Check for OpenAPI/Swagger
	if _, err := os.Stat("openapi.yaml"); err == nil {
		return "OpenAPI"
	}
	if _, err := os.Stat("swagger.yaml"); err == nil {
		return "OpenAPI"
	}
	
	// Default to REST for most web APIs
	return "REST"
}