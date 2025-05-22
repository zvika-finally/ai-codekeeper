package cursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MCPServer represents a Model Context Protocol server configuration
type MCPServer struct {
	Name        string            `json:"name"`
	Command     []string          `json:"command"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Description string            `json:"description"`
	Capabilities []string         `json:"capabilities"`
}

// CursorConfig represents Cursor IDE configuration with AI guard rails
type CursorConfig struct {
	AIRules struct {
		Enabled     bool     `json:"enabled"`
		Domain      string   `json:"domain"`
		GuardRails  []string `json:"guard_rails"`
		Enforcement string   `json:"enforcement"`
	} `json:"ai_rules"`

	MCPServers map[string]MCPServer `json:"mcp_servers"`

	Settings struct {
		AIAssistant struct {
			Model                string `json:"model"`
			Temperature          float64 `json:"temperature"`
			MaxTokens           int     `json:"max_tokens"`
			UseProjectContext   bool    `json:"use_project_context"`
			UseGuardRails       bool    `json:"use_guard_rails"`
		} `json:"ai_assistant"`

		CodeGeneration struct {
			FollowGuardRails    bool     `json:"follow_guard_rails"`
			RequiredPatterns    []string `json:"required_patterns"`
			ForbiddenPatterns   []string `json:"forbidden_patterns"`
			DomainTemplates     bool     `json:"domain_templates"`
		} `json:"code_generation"`

		AutoCompletion struct {
			EnableGuardRails    bool `json:"enable_guard_rails"`
			ValidateOnAccept    bool `json:"validate_on_accept"`
		} `json:"auto_completion"`
	} `json:"settings"`

	Rules []CursorRule `json:"rules"`
}

// CursorRule represents a specific coding rule for Cursor AI
type CursorRule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern,omitempty"`
	Replacement string   `json:"replacement,omitempty"`
	Severity    string   `json:"severity"` // "error", "warning", "info"
	Applies     []string `json:"applies"`  // File patterns
	Domain      string   `json:"domain,omitempty"`
}

// GenerateCursorConfig creates Cursor configuration with MCP integration
func GenerateCursorConfig(projectPath, domain string, guardRails []string) (*CursorConfig, error) {
	config := &CursorConfig{
		MCPServers: make(map[string]MCPServer),
		Rules:      []CursorRule{},
	}

	// Configure AI rules
	config.AIRules.Enabled = true
	config.AIRules.Domain = domain
	config.AIRules.GuardRails = guardRails
	config.AIRules.Enforcement = "advisory"

	// Configure AI assistant settings
	config.Settings.AIAssistant.Model = "claude-3-sonnet"
	config.Settings.AIAssistant.Temperature = 0.1
	config.Settings.AIAssistant.MaxTokens = 4000
	config.Settings.AIAssistant.UseProjectContext = true
	config.Settings.AIAssistant.UseGuardRails = true

	// Configure code generation
	config.Settings.CodeGeneration.FollowGuardRails = true
	config.Settings.CodeGeneration.DomainTemplates = true
	config.Settings.CodeGeneration.RequiredPatterns = getRequiredPatterns(domain)
	config.Settings.CodeGeneration.ForbiddenPatterns = getForbiddenPatterns(domain)

	// Configure auto-completion
	config.Settings.AutoCompletion.EnableGuardRails = true
	config.Settings.AutoCompletion.ValidateOnAccept = true

	// Add MCP servers
	config.MCPServers["codekeeper-guardrails"] = MCPServer{
		Name:        "Finally AI CodeKeeper Guard Rails",
		Command:     []string{"codekeeper", "mcp-server"},
		Description: "Provides guard rails validation and domain expertise",
		Capabilities: []string{"validation", "completion", "analysis"},
		Env: map[string]string{
			"CODEKEEPER_PROJECT_PATH": projectPath,
			"CODEKEEPER_DOMAIN":       domain,
		},
	}

	config.MCPServers["codekeeper-domain"] = MCPServer{
		Name:        "Finally AI CodeKeeper Domain Expert",
		Command:     []string{"codekeeper", "domain-server", "--domain", domain},
		Description: fmt.Sprintf("Provides %s domain expertise and patterns", domain),
		Capabilities: []string{"templates", "recommendations", "compliance"},
		Env: map[string]string{
			"CODEKEEPER_DOMAIN": domain,
		},
	}

	// Add domain-specific rules
	domainRules := getDomainRules(domain)
	config.Rules = append(config.Rules, domainRules...)

	return config, nil
}

// SaveCursorConfig writes Cursor configuration to .cursor directory
func SaveCursorConfig(projectPath string, config *CursorConfig) error {
	cursorDir := filepath.Join(projectPath, ".cursor")
	if err := os.MkdirAll(cursorDir, 0755); err != nil {
		return err
	}

	// Save main configuration
	configPath := filepath.Join(cursorDir, "config.json")
	if err := saveJSON(configPath, config); err != nil {
		return err
	}

	// Save rules file
	rulesPath := filepath.Join(cursorDir, "rules.json")
	rules := map[string]interface{}{
		"version": "1.0.0",
		"domain":  config.AIRules.Domain,
		"rules":   config.Rules,
	}
	if err := saveJSON(rulesPath, rules); err != nil {
		return err
	}

	// Generate .cursorrules file (plain text format)
	cursorRulesPath := filepath.Join(projectPath, ".cursorrules")
	if err := generateCursorRulesFile(cursorRulesPath, config); err != nil {
		return err
	}

	return nil
}

// generateCursorRulesFile creates the .cursorrules file that Cursor reads
func generateCursorRulesFile(filePath string, config *CursorConfig) error {
	content := fmt.Sprintf(`# Finally AI CodeKeeper Rules for %s Domain

## 🎯 Mission
Finally AI CodeKeeper ensures AI assistants follow domain-specific guard rails and best practices.

## Core Principles
- Follow %s domain best practices
- Implement security-first development
- Ensure compliance with industry standards
- Use appropriate data types for domain-specific operations

## Guard Rails
%s

## Code Generation Rules
%s

## Forbidden Patterns
%s

## Domain-Specific Instructions
%s

Generated by Finally AI CodeKeeper v1.0.0
Finally, AI development you can trust.
`, 
		config.AIRules.Domain,
		config.AIRules.Domain,
		formatGuardRails(config.AIRules.GuardRails),
		formatRequiredPatterns(config.Settings.CodeGeneration.RequiredPatterns),
		formatForbiddenPatterns(config.Settings.CodeGeneration.ForbiddenPatterns),
		getDomainInstructions(config.AIRules.Domain),
	)

	return os.WriteFile(filePath, []byte(content), 0644)
}

// Domain-specific helper functions
func getRequiredPatterns(domain string) []string {
	switch domain {
	case "fintech":
		return []string{
			"Use Decimal type for monetary amounts",
			"Include audit logging for financial operations",
			"Implement input validation for all APIs",
			"Use proper error handling with structured responses",
			"Include idempotency keys for financial transactions",
		}
	case "healthcare":
		return []string{
			"Encrypt all PII and PHI data",
			"Implement proper access controls",
			"Include audit logging for data access",
			"Use proper data retention policies",
		}
	default:
		return []string{
			"Use TypeScript for type safety",
			"Include proper error handling",
			"Implement input validation",
			"Follow security best practices",
		}
	}
}

func getForbiddenPatterns(domain string) []string {
	switch domain {
	case "fintech":
		return []string{
			"Never use floating point for monetary calculations",
			"Never log sensitive financial data",
			"Never hardcode financial limits or rates",
			"Never skip input validation on financial APIs",
		}
	case "healthcare":
		return []string{
			"Never log PHI data",
			"Never store passwords in plain text",
			"Never skip encryption for sensitive data",
			"Never expose patient data in URLs",
		}
	default:
		return []string{
			"Never hardcode secrets or API keys",
			"Never skip input validation",
			"Never ignore error handling",
			"Never log sensitive user data",
		}
	}
}

func getDomainRules(domain string) []CursorRule {
	switch domain {
	case "fintech":
		return []CursorRule{
			{
				Name:        "Decimal Arithmetic",
				Description: "Use Decimal types for all monetary calculations",
				Pattern:     `(amount|price|balance|fee).*(\+|\-|\*|\/).*[0-9]*\.[0-9]`,
				Severity:    "error",
				Applies:     []string{"*.ts", "*.js"},
				Domain:      "fintech",
			},
			{
				Name:        "Audit Logging",
				Description: "Include audit logging for financial transactions",
				Pattern:     `(transaction|payment|transfer).*(?!.*audit)`,
				Severity:    "warning",
				Applies:     []string{"*.ts", "*.js"},
				Domain:      "fintech",
			},
			{
				Name:        "Idempotency",
				Description: "Financial operations must be idempotent",
				Pattern:     `(POST|PUT).*(/payment|/transfer|/transaction)`,
				Severity:    "warning",
				Applies:     []string{"*.ts", "*.js"},
				Domain:      "fintech",
			},
		}
	default:
		return []CursorRule{
			{
				Name:        "Input Validation",
				Description: "Validate all user inputs",
				Pattern:     `req\.(body|params|query).*(?!.*validate)`,
				Severity:    "warning",
				Applies:     []string{"*.ts", "*.js"},
			},
		}
	}
}

func getDomainInstructions(domain string) string {
	switch domain {
	case "fintech":
		return `
When generating financial code:
1. Always use Decimal or BigNumber for monetary amounts
2. Implement proper transaction logging with audit trails
3. Include compliance checks for financial regulations
4. Use secure communication (HTTPS) for all financial APIs
5. Implement proper authentication and authorization
6. Include fraud detection considerations
7. Ensure PCI compliance for payment card data
8. Implement proper data encryption at rest and in transit`

	case "healthcare":
		return `
When generating healthcare code:
1. Ensure HIPAA compliance for all patient data
2. Implement proper access controls and authentication
3. Use encryption for all PHI (Protected Health Information)
4. Include audit logging for data access and modifications
5. Implement proper data retention and deletion policies
6. Ensure secure communication for all health data transfers
7. Include consent management for data processing
8. Implement proper anonymization techniques`

	default:
		return `
When generating code:
1. Follow security best practices
2. Use TypeScript for type safety
3. Implement proper error handling
4. Include comprehensive input validation
5. Follow the project's coding standards
6. Write maintainable and testable code
7. Include appropriate logging and monitoring
8. Ensure performance considerations`
	}
}

// Formatting helper functions
func formatGuardRails(guardRails []string) string {
	result := ""
	for _, rule := range guardRails {
		result += fmt.Sprintf("- %s\n", rule)
	}
	return result
}

func formatRequiredPatterns(patterns []string) string {
	result := ""
	for _, pattern := range patterns {
		result += fmt.Sprintf("- %s\n", pattern)
	}
	return result
}

func formatForbiddenPatterns(patterns []string) string {
	result := ""
	for _, pattern := range patterns {
		result += fmt.Sprintf("- %s\n", pattern)
	}
	return result
}

// StartGuardRailsServer starts the MCP guard rails server
func StartGuardRailsServer(domain, projectPath string) error {
	// For now, this is a placeholder that outputs JSON-RPC messages
	// In a real implementation, this would start an actual MCP server
	
	// Initialize the server
	fmt.Println(`{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {"listChanged": true}, "resources": {"listChanged": true}}, "clientInfo": {"name": "codekeeper-guardrails", "version": "1.0.0"}}}`)
	
	// Keep server running and respond to requests
	for {
		// This would normally handle actual MCP protocol messages
		// For demo purposes, we'll just sleep
		select {}
	}
}

// StartDomainServer starts the MCP domain expertise server
func StartDomainServer(domain string) error {
	// For now, this is a placeholder that outputs JSON-RPC messages
	// In a real implementation, this would start an actual MCP server
	
	// Initialize the server
	fmt.Println(`{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {"listChanged": true}, "resources": {"listChanged": true}}, "clientInfo": {"name": "codekeeper-domain-expert", "version": "1.0.0"}}}`)
	
	// Keep server running and respond to requests
	for {
		// This would normally handle actual MCP protocol messages
		// For demo purposes, we'll just sleep
		select {}
	}
}

// Utility function
func saveJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}