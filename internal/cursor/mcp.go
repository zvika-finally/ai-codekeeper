package cursor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

Generated by Finally AI CodeKeeper v2.0.0
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

// MCPRequest represents an MCP JSON-RPC request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents an MCP JSON-RPC response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// StartGuardRailsServer starts the MCP guard rails server
func StartGuardRailsServer(domain, projectPath string) error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			continue
		}
		
		response := handleGuardRailsRequest(request, domain, projectPath)
		if response != nil {
			responseJSON, _ := json.Marshal(response)
			fmt.Println(string(responseJSON))
		}
	}
	
	return scanner.Err()
}

// StartDomainServer starts the MCP domain expertise server
func StartDomainServer(domain string) error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			continue
		}
		
		response := handleDomainRequest(request, domain)
		if response != nil {
			responseJSON, _ := json.Marshal(response)
			fmt.Println(string(responseJSON))
		}
	}
	
	return scanner.Err()
}

// handleGuardRailsRequest processes MCP requests for the guard rails server
func handleGuardRailsRequest(request MCPRequest, domain, projectPath string) *MCPResponse {
	switch request.Method {
	case "initialize":
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{
						"listChanged": true,
					},
					"resources": map[string]interface{}{
						"listChanged": true,
					},
				},
				"serverInfo": map[string]interface{}{
					"name":    "Finally AI CodeKeeper Guard Rails",
					"version": "1.0.0",
				},
			},
		}
		
	case "tools/list":
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "validate_code",
						"description": "Validate code against domain-specific guard rails",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"code": map[string]interface{}{
									"type":        "string",
									"description": "Code to validate",
								},
								"language": map[string]interface{}{
									"type":        "string",
									"description": "Programming language",
								},
							},
							"required": []string{"code"},
						},
					},
					{
						"name":        "suggest_fix",
						"description": "Suggest fixes for guard rail violations",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"violation": map[string]interface{}{
									"type":        "string",
									"description": "Guard rail violation description",
								},
								"code": map[string]interface{}{
									"type":        "string",
									"description": "Problematic code",
								},
							},
							"required": []string{"violation", "code"},
						},
					},
				},
			},
		}
		
	case "tools/call":
		if params, ok := request.Params.(map[string]interface{}); ok {
			if name, ok := params["name"].(string); ok {
				switch name {
				case "validate_code":
					return handleValidateCode(request.ID, params, domain)
				case "suggest_fix":
					return handleSuggestFix(request.ID, params, domain)
				}
			}
		}
		
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
	
	return nil
}

// handleDomainRequest processes MCP requests for the domain expertise server
func handleDomainRequest(request MCPRequest, domain string) *MCPResponse {
	switch request.Method {
	case "initialize":
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{
						"listChanged": true,
					},
					"resources": map[string]interface{}{
						"listChanged": true,
					},
				},
				"serverInfo": map[string]interface{}{
					"name":    fmt.Sprintf("Finally AI CodeKeeper %s Expert", strings.Title(domain)),
					"version": "1.0.0",
				},
			},
		}
		
	case "tools/list":
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"tools": getDomainTools(domain),
			},
		}
		
	case "tools/call":
		if params, ok := request.Params.(map[string]interface{}); ok {
			if name, ok := params["name"].(string); ok {
				return handleDomainToolCall(request.ID, name, params, domain)
			}
		}
		
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
	
	return nil
}

// handleValidateCode validates code against guard rails
func handleValidateCode(id interface{}, params map[string]interface{}, domain string) *MCPResponse {
	code, _ := params["code"].(string)
	language, _ := params["language"].(string)
	
	violations := validateCodeAgainstGuardRails(code, language, domain)
	
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Validation Results:\n%s", formatValidationResults(violations)),
				},
			},
		},
	}
}

// handleSuggestFix suggests fixes for violations
func handleSuggestFix(id interface{}, params map[string]interface{}, domain string) *MCPResponse {
	violation, _ := params["violation"].(string)
	code, _ := params["code"].(string)
	
	suggestion := generateFixSuggestion(violation, code, domain)
	
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": suggestion,
				},
			},
		},
	}
}

// validateCodeAgainstGuardRails performs actual validation
func validateCodeAgainstGuardRails(code, language, domain string) []string {
	var violations []string
	
	// Domain-specific validations
	switch domain {
	case "fintech":
		if strings.Contains(code, "parseFloat") || strings.Contains(code, "parseInt") {
			violations = append(violations, "Use Decimal or BigNumber for financial calculations instead of parseFloat/parseInt")
		}
		if strings.Contains(code, "Math.random()") {
			violations = append(violations, "Use cryptographically secure random number generation for financial operations")
		}
		if strings.Contains(code, "console.log") && (strings.Contains(code, "amount") || strings.Contains(code, "balance")) {
			violations = append(violations, "Never log sensitive financial data")
		}
		
	case "healthcare":
		if strings.Contains(code, "console.log") && (strings.Contains(code, "patient") || strings.Contains(code, "phi")) {
			violations = append(violations, "Never log PHI (Protected Health Information)")
		}
		if strings.Contains(code, "http://") {
			violations = append(violations, "Use HTTPS for all healthcare data transmission")
		}
		
	default:
		if strings.Contains(code, "eval(") {
			violations = append(violations, "Avoid using eval() - security risk")
		}
		if strings.Contains(code, "innerHTML =") {
			violations = append(violations, "Use textContent or proper sanitization to prevent XSS")
		}
	}
	
	// General validations
	if strings.Contains(code, "password") && strings.Contains(code, "=") && !strings.Contains(code, "hash") {
		violations = append(violations, "Never store passwords in plain text")
	}
	
	return violations
}

// generateFixSuggestion creates fix suggestions
func generateFixSuggestion(violation, code, domain string) string {
	switch {
	case strings.Contains(violation, "parseFloat"):
		return "Replace parseFloat() with Decimal.js or a similar library:\n\n" +
			"// Instead of:\nconst amount = parseFloat(value);\n\n" +
			"// Use:\nconst amount = new Decimal(value);"
			
	case strings.Contains(violation, "Math.random"):
		return "Replace Math.random() with crypto.randomBytes():\n\n" +
			"// Instead of:\nconst randomValue = Math.random();\n\n" +
			"// Use:\nconst crypto = require('crypto');\nconst randomValue = crypto.randomBytes(4).readUInt32BE(0) / 0xFFFFFFFF;"
			
	case strings.Contains(violation, "console.log"):
		return "Remove or replace console.log with proper logging:\n\n" +
			"// Remove sensitive data logging or use structured logging:\nlogger.info('Operation completed', { operationId: id });"
			
	case strings.Contains(violation, "eval"):
		return "Replace eval() with safer alternatives:\n\n" +
			"// Instead of eval(), use JSON.parse() for data or Function() constructor for safer code evaluation"
			
	default:
		return fmt.Sprintf("Fix suggestion for: %s\n\nReview the %s domain best practices and apply appropriate security measures.", violation, domain)
	}
}

// getDomainTools returns available tools for a domain
func getDomainTools(domain string) []map[string]interface{} {
	switch domain {
	case "fintech":
		return []map[string]interface{}{
			{
				"name":        "generate_payment_api",
				"description": "Generate secure payment API endpoints",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"operation": map[string]interface{}{
							"type": "string",
							"enum": []string{"payment", "refund", "transfer"},
						},
					},
				},
			},
			{
				"name":        "compliance_check",
				"description": "Check code for financial compliance",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"code": map[string]interface{}{"type": "string"},
					},
				},
			},
		}
	default:
		return []map[string]interface{}{
			{
				"name":        "best_practices",
				"description": "Get domain-specific best practices",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"topic": map[string]interface{}{"type": "string"},
					},
				},
			},
		}
	}
}

// handleDomainToolCall handles tool calls for domain expertise
func handleDomainToolCall(id interface{}, name string, params map[string]interface{}, domain string) *MCPResponse {
	switch name {
	case "generate_payment_api":
		operation, _ := params["operation"].(string)
		code := generatePaymentAPI(operation, domain)
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": code},
				},
			},
		}
		
	case "compliance_check":
		code, _ := params["code"].(string)
		results := performComplianceCheck(code, domain)
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": results},
				},
			},
		}
		
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Tool not found",
			},
		}
	}
}

// Helper functions
func formatValidationResults(violations []string) string {
	if len(violations) == 0 {
		return "✅ No guard rail violations found!"
	}
	
	result := fmt.Sprintf("⚠️ Found %d guard rail violations:\n\n", len(violations))
	for i, violation := range violations {
		result += fmt.Sprintf("%d. %s\n", i+1, violation)
	}
	
	return result
}

func generatePaymentAPI(operation, domain string) string {
	switch operation {
	case "payment":
		return `// Secure Payment API Endpoint
const Decimal = require('decimal.js');
const { v4: uuidv4 } = require('uuid');

app.post('/api/payments', async (req, res) => {
  try {
    const { amount, currency, paymentMethodId } = req.body;
    
    // Validate input
    if (!amount || !currency || !paymentMethodId) {
      return res.status(400).json({ error: 'Missing required fields' });
    }
    
    // Use Decimal for monetary calculations
    const paymentAmount = new Decimal(amount);
    if (paymentAmount.lessThanOrEqualTo(0)) {
      return res.status(400).json({ error: 'Amount must be positive' });
    }
    
    // Generate idempotency key
    const idempotencyKey = req.headers['idempotency-key'] || uuidv4();
    
    // Process payment with audit logging
    const payment = await processPayment({
      amount: paymentAmount,
      currency,
      paymentMethodId,
      idempotencyKey,
      userId: req.user.id
    });
    
    // Log transaction for audit
    await auditLogger.log('PAYMENT_CREATED', {
      paymentId: payment.id,
      amount: paymentAmount.toString(),
      currency,
      userId: req.user.id
    });
    
    res.json({ success: true, payment });
  } catch (error) {
    logger.error('Payment processing failed', { error: error.message });
    res.status(500).json({ error: 'Payment processing failed' });
  }
});`
	default:
		return fmt.Sprintf("// %s API generation not yet implemented", operation)
	}
}

func performComplianceCheck(code, domain string) string {
	violations := validateCodeAgainstGuardRails(code, "javascript", domain)
	return formatValidationResults(violations)
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