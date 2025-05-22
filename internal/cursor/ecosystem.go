package cursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MCPEcosystemConfig represents comprehensive MCP server configuration
type MCPEcosystemConfig struct {
	Version     string                    `json:"version"`
	Domain      string                    `json:"domain"`
	MCPServers  map[string]MCPServerSpec  `json:"mcpServers"`
	Integrations map[string]Integration   `json:"integrations"`
	GuardRails  GuardRailsConfig         `json:"guardRails"`
	Cline       ClineConfig              `json:"cline"`
}

// MCPServerSpec defines an MCP server configuration
type MCPServerSpec struct {
	Name         string            `json:"name"`
	Command      []string          `json:"command"`
	Args         []string          `json:"args,omitempty"`
	Env          map[string]string `json:"env,omitempty"`
	Description  string            `json:"description"`
	Capabilities []string          `json:"capabilities"`
	Enabled      bool              `json:"enabled"`
	Priority     int               `json:"priority"`
}

// Integration defines third-party service integration
type Integration struct {
	Type         string            `json:"type"`
	Enabled      bool              `json:"enabled"`
	Config       map[string]string `json:"config"`
	Guards       []string          `json:"guards"`
}

// GuardRailsConfig defines domain-specific guard rails
type GuardRailsConfig struct {
	Enforcement string                `json:"enforcement"`
	Rules       map[string]GuardRule  `json:"rules"`
	Domains     map[string]DomainSpec `json:"domains"`
}

// GuardRule defines a specific guard rail rule
type GuardRule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern"`
	Severity    string   `json:"severity"`
	FileTypes   []string `json:"fileTypes"`
	Actions     []string `json:"actions"`
}

// DomainSpec defines domain-specific configurations
type DomainSpec struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Requirements []string           `json:"requirements"`
	Patterns     map[string]string  `json:"patterns"`
	Integrations []string           `json:"integrations"`
}

// ClineConfig defines Cline AI assistant configuration
type ClineConfig struct {
	Enabled      bool              `json:"enabled"`
	Model        string            `json:"model"`
	Provider     string            `json:"provider"`
	Settings     ClineSettings     `json:"settings"`
	MCPServers   []string          `json:"mcpServers"`
	GuardRails   ClineGuardRails   `json:"guardRails"`
}

// ClineSettings defines Cline-specific settings
type ClineSettings struct {
	MaxTokens          int                    `json:"maxTokens"`
	Temperature        float64                `json:"temperature"`
	SystemPrompt       string                 `json:"systemPrompt"`
	UseProjectContext  bool                   `json:"useProjectContext"`
	AutoApproveChanges bool                   `json:"autoApproveChanges"`
	DomainExpertise    map[string]interface{} `json:"domainExpertise"`
}

// ClineGuardRails defines guard rails for Cline
type ClineGuardRails struct {
	Enabled          bool     `json:"enabled"`
	StrictMode       bool     `json:"strictMode"`
	RequiredChecks   []string `json:"requiredChecks"`
	BlockedPatterns  []string `json:"blockedPatterns"`
	RequiredPatterns []string `json:"requiredPatterns"`
}

// GenerateComprehensiveMCPConfig creates a complete MCP ecosystem configuration
func GenerateComprehensiveMCPConfig(projectPath, domain string, services []string) error {
	config := &MCPEcosystemConfig{
		Version:      "2.0.0",
		Domain:       domain,
		MCPServers:   make(map[string]MCPServerSpec),
		Integrations: make(map[string]Integration),
		GuardRails:   generateGuardRailsConfig(domain),
		Cline:        generateClineConfig(domain, services),
	}

	// Add core AI Development Framework servers
	config.MCPServers["codekeeper-guardrails"] = MCPServerSpec{
		Name:        "AI Development Guard Rails",
		Command:     []string{"codekeeper", "mcp-server"},
		Description: "Enforces domain-specific guard rails and best practices",
		Capabilities: []string{"validation", "analysis", "suggestions"},
		Enabled:     true,
		Priority:    1,
		Env: map[string]string{
			"CODEKEEPER_DOMAIN":       domain,
			"CODEKEEPER_PROJECT_PATH": projectPath,
		},
	}

	config.MCPServers["codekeeper-domain"] = MCPServerSpec{
		Name:        "Domain Expert",
		Command:     []string{"codekeeper", "domain-server", "--domain", domain},
		Description: fmt.Sprintf("Provides %s domain expertise and templates", domain),
		Capabilities: []string{"templates", "recommendations", "compliance"},
		Enabled:     true,
		Priority:    2,
		Env: map[string]string{
			"CODEKEEPER_DOMAIN": domain,
		},
	}

	// Add requested ecosystem services
	for _, service := range services {
		serverSpec, integration := generateServiceConfig(service, domain)
		if serverSpec != nil {
			config.MCPServers[service] = *serverSpec
		}
		if integration != nil {
			config.Integrations[service] = *integration
		}
	}

	// Save configuration
	return saveEcosystemConfig(projectPath, config)
}

// generateServiceConfig creates MCP server and integration config for a service
func generateServiceConfig(service, domain string) (*MCPServerSpec, *Integration) {
	switch service {
	case "git":
		return &MCPServerSpec{
			Name:        "Git Integration",
			Command:     []string{"mcp-git"},
			Description: "Git version control with guard rails",
			Capabilities: []string{"commit", "branch", "merge", "status"},
			Enabled:     true,
			Priority:    3,
			Env: map[string]string{
				"GIT_GUARD_RAILS": "true",
				"DOMAIN":          domain,
			},
		}, &Integration{
			Type:    "vcs",
			Enabled: true,
			Config: map[string]string{
				"enforceCommitMessages": "true",
				"requireBranchPrefix":   "true",
				"guardRailsOnCommit":    "true",
			},
			Guards: []string{"commit-message-format", "no-secrets", "code-quality"},
		}

	case "github":
		return &MCPServerSpec{
			Name:        "GitHub Integration",
			Command:     []string{"mcp-github"},
			Description: "GitHub repository management with domain compliance",
			Capabilities: []string{"pr", "issues", "reviews", "actions"},
			Enabled:     true,
			Priority:    4,
			Env: map[string]string{
				"GITHUB_TOKEN": "${GITHUB_TOKEN}",
				"DOMAIN":       domain,
			},
		}, &Integration{
			Type:    "repository",
			Enabled: true,
			Config: map[string]string{
				"autoLinkIssues":     "true",
				"enforceReviews":     "true",
				"domainCompliance":   "true",
			},
			Guards: []string{"pr-template", "security-review", "compliance-check"},
		}

	case "jira":
		return &MCPServerSpec{
			Name:        "JIRA Integration",
			Command:     []string{"mcp-jira"},
			Description: "JIRA project management with domain workflows",
			Capabilities: []string{"tickets", "sprints", "reporting", "automation"},
			Enabled:     true,
			Priority:    5,
			Env: map[string]string{
				"JIRA_URL":    "${JIRA_URL}",
				"JIRA_TOKEN":  "${JIRA_TOKEN}",
				"DOMAIN":      domain,
			},
		}, &Integration{
			Type:    "project-management",
			Enabled: true,
			Config: map[string]string{
				"autoCreateBranches": "true",
				"linkCommits":        "true",
				"domainWorkflows":    "true",
			},
			Guards: []string{"ticket-validation", "workflow-compliance"},
		}

	case "figma":
		return &MCPServerSpec{
			Name:        "Figma Integration",
			Command:     []string{"mcp-figma"},
			Description: "Figma design collaboration with component sync",
			Capabilities: []string{"designs", "components", "tokens", "handoff"},
			Enabled:     true,
			Priority:    6,
			Env: map[string]string{
				"FIGMA_TOKEN": "${FIGMA_TOKEN}",
				"DOMAIN":      domain,
			},
		}, &Integration{
			Type:    "design",
			Enabled: true,
			Config: map[string]string{
				"syncComponents":   "true",
				"extractTokens":    "true",
				"designCompliance": "true",
			},
			Guards: []string{"design-system", "accessibility", "brand-compliance"},
		}

	case "aws":
		return &MCPServerSpec{
			Name:        "AWS Integration",
			Command:     []string{"mcp-aws"},
			Description: "AWS cloud services with security best practices",
			Capabilities: []string{"resources", "deployment", "monitoring", "security"},
			Enabled:     true,
			Priority:    7,
			Env: map[string]string{
				"AWS_REGION":         "${AWS_REGION}",
				"AWS_ACCESS_KEY_ID":  "${AWS_ACCESS_KEY_ID}",
				"AWS_SECRET_ACCESS_KEY": "${AWS_SECRET_ACCESS_KEY}",
				"DOMAIN":             domain,
			},
		}, &Integration{
			Type:    "cloud",
			Enabled: true,
			Config: map[string]string{
				"securityFirst":      "true",
				"costOptimization":   "true",
				"complianceChecks":   "true",
			},
			Guards: []string{"security-groups", "iam-policies", "encryption", "cost-limits"},
		}

	case "terraform":
		return &MCPServerSpec{
			Name:        "Terraform Integration",
			Command:     []string{"mcp-terraform"},
			Description: "Infrastructure as Code with domain compliance",
			Capabilities: []string{"plan", "apply", "state", "modules"},
			Enabled:     true,
			Priority:    8,
			Env: map[string]string{
				"TF_VAR_domain": domain,
			},
		}, &Integration{
			Type:    "infrastructure",
			Enabled: true,
			Config: map[string]string{
				"validatePlans":    "true",
				"securityScanning": "true",
				"costEstimation":   "true",
			},
			Guards: []string{"resource-naming", "security-rules", "backup-policies"},
		}

	case "docker":
		return &MCPServerSpec{
			Name:        "Docker Integration",
			Command:     []string{"mcp-docker"},
			Description: "Container management with security scanning",
			Capabilities: []string{"build", "scan", "run", "compose"},
			Enabled:     true,
			Priority:    9,
		}, &Integration{
			Type:    "containerization",
			Enabled: true,
			Config: map[string]string{
				"securityScanning": "true",
				"vulnAssessment":   "true",
				"bestPractices":    "true",
			},
			Guards: []string{"dockerfile-security", "image-scanning", "resource-limits"},
		}

	default:
		return nil, nil
	}
}

// generateGuardRailsConfig creates domain-specific guard rails configuration
func generateGuardRailsConfig(domain string) GuardRailsConfig {
	config := GuardRailsConfig{
		Enforcement: "advisory",
		Rules:       make(map[string]GuardRule),
		Domains:     make(map[string]DomainSpec),
	}

	// Add domain-specific rules
	switch domain {
	case "fintech":
		config.Rules["decimal-arithmetic"] = GuardRule{
			Name:        "Decimal Arithmetic",
			Description: "Use Decimal types for monetary calculations",
			Pattern:     `(amount|price|balance|fee).*(\+|\-|\*|\/).*[0-9]*\.[0-9]`,
			Severity:    "error",
			FileTypes:   []string{"*.ts", "*.js", "*.py"},
			Actions:     []string{"suggest-decimal", "block-commit"},
		}

		config.Rules["audit-logging"] = GuardRule{
			Name:        "Audit Logging",
			Description: "Financial operations must include audit trails",
			Pattern:     `(transaction|payment|transfer).*(?!.*audit)`,
			Severity:    "warning",
			FileTypes:   []string{"*.ts", "*.js", "*.py"},
			Actions:     []string{"suggest-audit", "require-review"},
		}

		config.Domains["fintech"] = DomainSpec{
			Name:        "Financial Technology",
			Description: "Banking, payments, and financial services",
			Requirements: []string{
				"PCI DSS compliance",
				"SOX compliance",
				"GDPR compliance",
				"Decimal arithmetic for money",
				"Audit trails for all transactions",
			},
			Patterns: map[string]string{
				"monetary":    "Use Decimal type for money calculations",
				"audit":       "Include audit logging for financial operations",
				"validation":  "Validate all financial inputs",
				"encryption":  "Encrypt sensitive financial data",
			},
			Integrations: []string{"aws", "terraform", "github", "jira"},
		}
	}

	return config
}

// generateClineConfig creates Cline AI assistant configuration
func generateClineConfig(domain string, services []string) ClineConfig {
	systemPrompt := fmt.Sprintf(`You are an AI development assistant specialized in %s domain development.

CRITICAL GUARD RAILS - ALWAYS FOLLOW:
1. Follow domain-specific best practices for %s
2. Validate all code against guard rails before suggesting
3. Use MCP servers for all external integrations
4. Ensure security and compliance requirements are met
5. Generate code that passes all domain validation rules

Available MCP servers: %v

When making changes:
- Always check guard rails before proceeding
- Use domain expertise for recommendations
- Validate security and compliance requirements
- Generate comprehensive tests
- Include proper documentation`, domain, domain, services)

	return ClineConfig{
		Enabled:  true,
		Model:    "claude-3-5-sonnet-20241022",
		Provider: "anthropic",
		Settings: ClineSettings{
			MaxTokens:          8000,
			Temperature:        0.1,
			SystemPrompt:       systemPrompt,
			UseProjectContext:  true,
			AutoApproveChanges: false,
			DomainExpertise: map[string]interface{}{
				"domain":       domain,
				"guardRails":   true,
				"compliance":   true,
				"security":     true,
			},
		},
		MCPServers: append([]string{"codekeeper-guardrails", "codekeeper-domain"}, services...),
		GuardRails: ClineGuardRails{
			Enabled:    true,
			StrictMode: domain == "fintech", // Strict mode for financial domain
			RequiredChecks: []string{
				"domain-compliance",
				"security-validation",
				"guard-rails-check",
			},
			BlockedPatterns:  getDomainBlockedPatterns(domain),
			RequiredPatterns: getDomainRequiredPatterns(domain),
		},
	}
}

// GenerateClineIntegration creates Cline-specific configuration files
func GenerateClineIntegration(projectPath, domain string, services []string) error {
	clineConfig := generateClineConfig(domain, services)

	// Create .cline directory
	clineDir := filepath.Join(projectPath, ".cline")
	if err := os.MkdirAll(clineDir, 0755); err != nil {
		return err
	}

	// Save Cline configuration
	clineConfigPath := filepath.Join(clineDir, "config.json")
	if err := saveJSON(clineConfigPath, clineConfig); err != nil {
		return err
	}

	// Generate Cline system prompts
	systemPromptsPath := filepath.Join(clineDir, "system-prompts.md")
	if err := generateClineSystemPrompts(systemPromptsPath, domain, services); err != nil {
		return err
	}

	// Generate Cline guard rails
	guardRailsPath := filepath.Join(clineDir, "guard-rails.json")
	guardRails := map[string]interface{}{
		"domain":      domain,
		"strictMode":  domain == "fintech",
		"rules":       getDomainRules(domain),
		"mcpServers":  services,
		"validation":  true,
	}
	if err := saveJSON(guardRailsPath, guardRails); err != nil {
		return err
	}

	return nil
}

// generateClineSystemPrompts creates system prompts for Cline
func generateClineSystemPrompts(filePath, domain string, services []string) error {
	content := fmt.Sprintf(`# Cline System Prompts for %s Domain

## Core Instructions

You are an AI development assistant specialized in %s development with access to comprehensive MCP ecosystem.

### CRITICAL GUARD RAILS - ALWAYS FOLLOW:

1. **Domain Compliance**: All code must follow %s domain best practices
2. **Guard Rails Validation**: Validate all suggestions against domain guard rails
3. **MCP Integration**: Use MCP servers for all external service interactions
4. **Security First**: Prioritize security and compliance in all recommendations
5. **Quality Assurance**: Ensure code passes all validation rules before suggesting

### Available MCP Servers:
%s

### Domain-Specific Requirements:

%s

### Workflow Instructions:

1. **Before Making Changes**:
   - Check current project structure using MCP
   - Validate against domain guard rails
   - Verify compliance requirements
   - Check for existing patterns

2. **When Generating Code**:
   - Use domain-specific templates from MCP
   - Follow established patterns
   - Include proper validation
   - Add comprehensive tests
   - Document security considerations

3. **After Changes**:
   - Run guard rails validation
   - Execute relevant tests
   - Update documentation
   - Commit with proper messages

### Error Handling:

If guard rails are violated:
1. Explain the violation clearly
2. Suggest compliant alternatives
3. Provide domain-specific guidance
4. Use MCP servers for validation

### Integration Guidelines:

- **Git**: Use semantic commit messages with domain prefix
- **GitHub/GitLab**: Link to relevant issues, enforce reviews
- **JIRA**: Update tickets, follow workflows
- **AWS**: Follow security best practices, cost optimization
- **Figma**: Sync design tokens, maintain component library

Remember: Quality and compliance over speed. Always validate before suggesting.
`, domain, domain, domain, formatMCPServices(services), getDomainSpecificRequirements(domain))

	return os.WriteFile(filePath, []byte(content), 0644)
}

// Helper functions
func getDomainBlockedPatterns(domain string) []string {
	switch domain {
	case "fintech":
		return []string{
			"float.*amount",
			"double.*price",
			"console.log.*password",
			"hardcoded.*key",
		}
	default:
		return []string{
			"console.log.*password",
			"hardcoded.*secret",
		}
	}
}

func getDomainRequiredPatterns(domain string) []string {
	switch domain {
	case "fintech":
		return []string{
			"Decimal.*amount",
			"audit.*log",
			"validate.*input",
			"encrypt.*sensitive",
		}
	default:
		return []string{
			"validate.*input",
			"error.*handling",
		}
	}
}

func getDomainSpecificRequirements(domain string) string {
	switch domain {
	case "fintech":
		return `**Financial Technology Requirements:**
- Use Decimal types for all monetary calculations
- Implement audit trails for financial transactions
- Follow PCI DSS compliance for payment data
- Encrypt sensitive financial information
- Implement proper authentication and authorization
- Include fraud detection considerations
- Ensure regulatory compliance (SOX, GDPR, etc.)`
	case "healthcare":
		return `**Healthcare Requirements:**
- HIPAA compliance for all patient data
- Encrypt PHI (Protected Health Information)
- Implement proper access controls
- Audit trails for data access
- Data retention and deletion policies
- Consent management systems`
	default:
		return `**General Requirements:**
- Follow security best practices
- Implement proper error handling
- Use TypeScript for type safety
- Include comprehensive tests
- Document APIs and components`
	}
}

func formatMCPServices(services []string) string {
	result := ""
	for _, service := range services {
		result += fmt.Sprintf("- **%s**: %s\n", service, getMCPServiceDescription(service))
	}
	return result
}

func getMCPServiceDescription(service string) string {
	descriptions := map[string]string{
		"git":       "Version control with guard rails enforcement",
		"github":    "Repository management and code reviews",
		"gitlab":    "DevOps platform integration",
		"jira":      "Project management and issue tracking",
		"figma":     "Design collaboration and component sync",
		"aws":       "Cloud services with security best practices",
		"terraform": "Infrastructure as Code with compliance",
		"docker":    "Container management with security scanning",
		"datadog":   "Monitoring and observability",
		"sentry":    "Error tracking and performance monitoring",
	}
	if desc, exists := descriptions[service]; exists {
		return desc
	}
	return "Integration service"
}

// saveEcosystemConfig saves the complete ecosystem configuration
func saveEcosystemConfig(projectPath string, config *MCPEcosystemConfig) error {
	// Save main ecosystem config
	mcpDir := filepath.Join(projectPath, ".mcp")
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return err
	}

	ecosystemPath := filepath.Join(mcpDir, "ecosystem.json")
	if err := saveJSON(ecosystemPath, config); err != nil {
		return err
	}

	// Update Cursor configuration with ecosystem
	cursorConfigPath := filepath.Join(projectPath, ".cursor", "config.json")
	if err := updateCursorWithEcosystem(cursorConfigPath, config); err != nil {
		return err
	}

	return nil
}

// updateCursorWithEcosystem updates Cursor config with ecosystem MCP servers
func updateCursorWithEcosystem(configPath string, ecosystemConfig *MCPEcosystemConfig) error {
	// Read existing cursor config
	var cursorConfig map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &cursorConfig)
	} else {
		cursorConfig = make(map[string]interface{})
	}

	// Add ecosystem MCP servers
	if mcpServers, ok := cursorConfig["mcp_servers"].(map[string]interface{}); ok {
		for name, server := range ecosystemConfig.MCPServers {
			mcpServers[name] = server
		}
	} else {
		cursorConfig["mcp_servers"] = ecosystemConfig.MCPServers
	}

	// Save updated config
	return saveJSON(configPath, cursorConfig)
}