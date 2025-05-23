package generator

import (
	"encoding/json"
	"fmt"
)

// generateCursorIntegration is a deprecated method for old generator
func (g *Generator) generateCursorIntegration() error {
	// This is for backward compatibility only
	// New modular generator uses MCP integration instead
	return nil
}

// CursorIntegration handles Cursor IDE integration for new modular generator
type CursorIntegration struct {
	spec *ProjectSpec
}

// NewCursorIntegration creates a new Cursor integration generator
func NewCursorIntegration(spec *ProjectSpec) *CursorIntegration {
	return &CursorIntegration{spec: spec}
}

// Generate creates Cursor IDE configuration files
func (ci *CursorIntegration) Generate() (map[string]string, error) {
	files := make(map[string]string)
	
	// Generate main Cursor configuration
	ci.generateCursorSettings(files)
	
	// Generate Cursor rules and prompts
	ci.generateCursorRules(files)
	
	// Generate Cursor chat configuration
	ci.generateCursorChat(files)
	
	// Generate development workflows
	ci.generateCursorWorkflows(files)
	
	return files, nil
}

// generateCursorSettings creates the main Cursor settings
func (ci *CursorIntegration) generateCursorSettings(files map[string]string) {
	settings := map[string]interface{}{
		"cursor.general.enableCodeActions": true,
		"cursor.general.enableInlineEdit":  true,
		"cursor.chat.enabled":              true,
		"cursor.prediction.enabled":        true,
		
		// Project-specific settings
		"cursor.chat.systemMessage": ci.generateSystemMessage(),
		"cursor.chat.rules":         ci.getCursorRules(),
		
		// Domain-specific configuration
		"cursor.customInstructions": ci.getDomainInstructions(),
		
		// File associations
		"files.associations": map[string]string{
			"*.env.example": "plaintext",
			"Dockerfile*":   "dockerfile",
			"*.tf":          "hcl",
			"*.tfvars":      "hcl",
		},
		
		// Language-specific settings
		"typescript.preferences.importModuleSpecifier": "relative",
		"javascript.preferences.importModuleSpecifier": "relative",
		"go.formatTool":                                "goimports",
		
		// Exclude patterns
		"files.exclude": map[string]bool{
			"**/node_modules":     true,
			"**/dist":            true,
			"**/build":           true,
			"**/.git":            true,
			"**/terraform.tfstate": true,
		},
	}
	
	settingsJSON, _ := json.MarshalIndent(settings, "", "  ")
	files[".cursor/settings.json"] = string(settingsJSON)
}

// generateCursorRules creates Cursor coding rules
func (ci *CursorIntegration) generateCursorRules(files map[string]string) {
	rules := `# Cursor Coding Rules for ` + ci.spec.Name + `

## Project Context
- **Domain**: ` + ci.spec.Domain + `
- **Core Entity**: ` + ci.spec.CoreEntity + `
- **Backend**: ` + ci.spec.Backend + `
- **API Style**: ` + ci.spec.APIStyle + `

## Coding Standards

### General Rules
1. Follow the architectural patterns defined in docs/ARCHITECTURE.md
2. Implement security-first development practices
3. Write clean, self-documenting code with minimal comments
4. Use TypeScript strict mode for frontend development
5. Follow the established project structure in docs/frontend/PROJECT_STRUCTURE.md

### Backend Rules (` + ci.spec.Backend + `)
` + ci.getBackendRules() + `

### Frontend Rules (React TypeScript)
1. Use functional components with hooks
2. Implement proper TypeScript interfaces for all props
3. Follow component naming conventions (PascalCase)
4. Use custom hooks for complex logic
5. Implement proper error boundaries and loading states

### Domain-Specific Rules (` + ci.spec.Domain + `)
` + ci.getDomainSpecificRules() + `

### Security Rules
1. Never log sensitive information (passwords, tokens, PII)
2. Validate all user inputs at boundaries
3. Use environment variables for configuration
4. Implement proper authentication and authorization
5. Follow HTTPS-only communication patterns

### Testing Rules
1. Write tests for business logic and critical paths
2. Use descriptive test names that explain the scenario
3. Mock external dependencies appropriately
4. Maintain test coverage above 80% for core functionality

### Infrastructure Rules
1. Use Infrastructure as Code (Terraform for AWS)
2. Implement proper secrets management
3. Configure monitoring and alerting
4. Use Docker for consistent environments
5. Follow 12-factor app principles

## File Naming Conventions
- Components: PascalCase (UserProfile.tsx)
- Hooks: camelCase with 'use' prefix (useUserData.ts)
- Utils: camelCase (formatDate.ts)
- Constants: UPPER_SNAKE_CASE
- Config files: kebab-case

## Git Workflow
- Use conventional commit messages (feat:, fix:, docs:, etc.)
- Create feature branches for new development
- Keep commits atomic and focused
- Write meaningful commit messages

## MCP Integration
This project includes MCP servers for enhanced AI assistance:
- Git operations and project structure
- Domain-specific compliance and patterns
- Project context and documentation access

Refer to docs/MCP_INTEGRATION.md for usage details.
`

	files[".cursor/rules.md"] = rules
}

// generateCursorChat creates chat configuration
func (ci *CursorIntegration) generateCursorChat(files map[string]string) {
	chatConfig := map[string]interface{}{
		"version": "1.0",
		"systemMessage": ci.generateSystemMessage(),
		"rules": []string{
			"Always follow the coding standards in .cursor/rules.md",
			"Consider domain-specific requirements for " + ci.spec.Domain,
			"Use the project documentation in docs/ for context",
			"Follow security best practices for " + ci.spec.Domain + " applications",
			"Implement proper error handling and validation",
		},
		"context": map[string]interface{}{
			"includeFiles": []string{
				"docs/README.md",
				"docs/ARCHITECTURE.md", 
				"docs/frontend/STANDARDS.md",
				"docs/API_DESIGN.md",
				".cursor/rules.md",
			},
			"excludePatterns": []string{
				"node_modules/**",
				"dist/**",
				"build/**",
				"*.log",
				"terraform.tfstate",
			},
		},
	}
	
	configJSON, _ := json.MarshalIndent(chatConfig, "", "  ")
	files[".cursor/chat.json"] = string(configJSON)
}

// generateCursorWorkflows creates development workflow templates
func (ci *CursorIntegration) generateCursorWorkflows(files map[string]string) {
	workflows := `# Development Workflows for Cursor

## Quick Commands

### Project Setup
- **Start Development**: ` + "`docker-compose up`" + `
- **Run Tests**: ` + "`npm test`" + ` (in apps/backend or apps/frontend)
- **Type Check**: ` + "`npm run type-check`" + ` (frontend)
- **Lint Code**: ` + "`npm run lint`" + `

### Infrastructure Commands
- **Deploy to AWS**: ` + "`cd infra/aws && terraform apply`" + `
- **Deploy to Render**: Automatic on git push to main
- **Health Check**: ` + "`./scripts/health-check.sh`" + `

### MCP Server Commands
- **Start Git MCP**: ` + "`cd scripts/mcp && npm run start:git`" + `
- **Start Project MCP**: ` + "`cd scripts/mcp && npm run start:project`" + `
` + ci.getDomainMCPCommands() + `

## Common Tasks

### Adding a New Feature
1. Create feature branch: ` + "`git checkout -b feature/new-feature`" + `
2. Implement following .cursor/rules.md
3. Add tests for new functionality
4. Update documentation if needed
5. Create pull request

### Debugging Issues
1. Check application logs: ` + "`docker-compose logs -f`" + `
2. Run health checks: ` + "`./scripts/health-check.sh`" + `
3. Verify environment variables: ` + "`cat .env.example`" + `

### Deployment Process
1. Ensure all tests pass: ` + "`npm test`" + `
2. Run security checks: ` + "`npm audit`" + `
3. Deploy to staging first
4. Run integration tests
5. Deploy to production

## Domain-Specific Workflows (` + ci.spec.Domain + `)
` + ci.getDomainSpecificWorkflows() + `

## Cursor-Specific Tips
- Use Ctrl/Cmd+K for inline code generation
- Use Ctrl/Cmd+L for chat interface
- Reference .cursor/rules.md for project guidelines
- Use MCP servers for enhanced project context
- Ask about domain-specific patterns for ` + ci.spec.Domain + `
`

	files[".cursor/workflows.md"] = workflows
}

// Helper methods

func (ci *CursorIntegration) generateSystemMessage() string {
	return fmt.Sprintf(`You are an expert %s developer working on "%s", a %s application.

**Project Context:**
- Domain: %s with specific compliance and security requirements
- Core Entity: %s
- Backend: %s with %s API design
- Architecture: Clean architecture with modular components

**Your Role:**
- Follow domain-specific best practices for %s
- Implement security-first development patterns
- Use the established project structure and conventions
- Consider compliance requirements and regulatory standards
- Write production-ready, maintainable code

**Key Guidelines:**
1. Always reference project documentation in docs/
2. Follow the coding rules in .cursor/rules.md
3. Consider domain-specific patterns and requirements
4. Implement proper error handling and validation
5. Use MCP servers for enhanced project context

**Available Resources:**
- Project documentation in docs/
- Domain-specific patterns in docs/frontend/
- Infrastructure configs in infra/
- MCP servers for real-time project context

Always ask clarifying questions about requirements and consider the broader system impact of your changes.`,
		ci.spec.Domain,
		ci.spec.Name,
		ci.spec.Domain,
		ci.spec.Domain,
		ci.spec.CoreEntity,
		ci.spec.Backend,
		ci.spec.APIStyle,
		ci.spec.Domain)
}

func (ci *CursorIntegration) getCursorRules() []string {
	rules := []string{
		"Follow project documentation in docs/",
		"Use TypeScript strict mode",
		"Implement proper error handling",
		"Follow security best practices",
		"Write tests for critical functionality",
	}
	
	// Add domain-specific rules
	switch ci.spec.Domain {
	case "fintech":
		rules = append(rules, 
			"Use Decimal arithmetic for monetary calculations",
			"Implement audit trails for financial operations",
			"Follow PCI-DSS compliance requirements",
		)
	case "healthcare":
		rules = append(rules,
			"Ensure HIPAA compliance for patient data",
			"Implement proper consent management",
			"Use encryption for sensitive medical data",
		)
	case "ecommerce":
		rules = append(rules,
			"Implement proper inventory management",
			"Use secure payment processing",
			"Consider cart persistence and session management",
		)
	}
	
	return rules
}

func (ci *CursorIntegration) getDomainInstructions() string {
	switch ci.spec.Domain {
	case "fintech":
		return "Focus on financial compliance, secure transaction processing, and audit trail implementation. Always consider PCI-DSS and SOX requirements."
	case "healthcare":
		return "Prioritize HIPAA compliance, patient data security, and proper consent management. Consider FHIR standards for interoperability."
	case "ecommerce":
		return "Focus on user experience, secure payment processing, and inventory management. Consider scalability for high traffic."
	default:
		return "Follow general web application best practices with focus on security and maintainability."
	}
}

func (ci *CursorIntegration) getBackendRules() string {
	switch ci.spec.Backend {
	case "javascript", "node", "nodejs":
		return `1. Use async/await instead of callbacks
2. Implement proper error middleware
3. Use environment variables for configuration
4. Follow Express.js best practices
5. Implement request validation with joi or zod`
	case "python":
		return `1. Follow PEP 8 style guidelines
2. Use type hints for function signatures
3. Implement proper exception handling
4. Use virtual environments for dependencies
5. Follow FastAPI or Django best practices`
	case "go":
		return `1. Follow effective Go principles
2. Use proper error handling patterns
3. Implement context cancellation
4. Use interfaces for abstraction
5. Write idiomatic Go code`
	default:
		return `1. Follow language-specific best practices
2. Implement proper error handling
3. Use appropriate design patterns
4. Write clean, maintainable code
5. Consider performance implications`
	}
}

func (ci *CursorIntegration) getDomainSpecificRules() string {
	switch ci.spec.Domain {
	case "fintech":
		return `1. Never use floating-point arithmetic for money calculations
2. Implement comprehensive audit logging for all financial operations
3. Use strong encryption for sensitive financial data
4. Implement proper authentication and authorization
5. Follow PCI-DSS requirements for payment card data
6. Implement idempotency for financial transactions
7. Use proper decimal libraries (Decimal.js, big.js, etc.)
8. Implement rate limiting and fraud detection
9. Ensure compliance with financial regulations (SOX, GDPR, etc.)
10. Use secure communication protocols (TLS 1.3+)`
	case "healthcare":
		return `1. Encrypt all PHI (Protected Health Information) at rest and in transit
2. Implement proper access controls with role-based permissions
3. Maintain audit logs for all data access and modifications
4. Use FHIR standards for healthcare data interoperability
5. Implement proper consent management workflows
6. Follow HIPAA compliance requirements
7. Use secure authentication with MFA when possible
8. Implement data retention and deletion policies
9. Ensure proper backup and disaster recovery
10. Use anonymization techniques for analytics`
	case "ecommerce":
		return `1. Implement proper inventory management and stock tracking
2. Use secure payment processing with PCI compliance
3. Implement cart persistence across sessions
4. Use proper product catalog management
5. Implement order management workflows
6. Use CDN for static assets and image optimization
7. Implement proper search and filtering capabilities
8. Use caching strategies for performance
9. Implement customer data protection (GDPR)
10. Use analytics for business intelligence`
	default:
		return `1. Follow general web application security practices
2. Implement proper input validation and sanitization
3. Use HTTPS for all communications
4. Implement proper authentication and authorization
5. Use secure coding practices
6. Implement proper logging and monitoring
7. Follow API design best practices
8. Use proper error handling and user feedback
9. Implement data backup and recovery
10. Consider scalability and performance`
	}
}

func (ci *CursorIntegration) getDomainMCPCommands() string {
	switch ci.spec.Domain {
	case "fintech":
		return `- **Start Compliance MCP**: ` + "`cd scripts/mcp && npm run start:compliance`"
	case "healthcare":
		return `- **Start HIPAA MCP**: ` + "`cd scripts/mcp && npm run start:hipaa`"
	case "ecommerce":
		return `- **Start E-commerce MCP**: ` + "`cd scripts/mcp && npm run start:ecommerce`"
	default:
		return `- **No domain-specific MCP servers**`
	}
}

func (ci *CursorIntegration) getDomainSpecificWorkflows() string {
	switch ci.spec.Domain {
	case "fintech":
		return `### Financial Transaction Workflow
1. Implement transaction validation
2. Add audit trail logging
3. Test with compliance checks
4. Verify encryption and security
5. Deploy with monitoring

### Compliance Workflow
1. Run compliance MCP server for guidance
2. Implement required audit trails
3. Test with sample financial data
4. Verify PCI-DSS requirements
5. Document compliance measures`
	case "healthcare":
		return `### Patient Data Workflow
1. Implement PHI encryption
2. Add audit logging for access
3. Test consent management
4. Verify HIPAA compliance
5. Deploy with security monitoring

### FHIR Integration Workflow
1. Design FHIR-compliant data models
2. Implement FHIR API endpoints
3. Test interoperability
4. Verify data validation
5. Document API specifications`
	case "ecommerce":
		return `### Product Management Workflow
1. Implement product catalog features
2. Add inventory management
3. Test search and filtering
4. Verify performance with large datasets
5. Deploy with CDN integration

### Payment Processing Workflow
1. Integrate payment provider APIs
2. Implement secure checkout flow
3. Test payment scenarios
4. Verify PCI compliance
5. Monitor transaction success rates`
	default:
		return `### General Development Workflow
1. Plan feature implementation
2. Write code following project standards
3. Add comprehensive tests
4. Review security implications
5. Deploy with proper monitoring`
	}
}