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
	
	// Generate MCP configuration for Cursor
	ci.generateCursorMCPConfig(files)
	
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

// generateCursorMCPConfig creates MCP server configuration for Cursor
func (ci *CursorIntegration) generateCursorMCPConfig(files map[string]string) {
	mcpConfig := map[string]interface{}{
		"mcpServers": ci.generateMCPServers(),
	}
	
	mcpJSON, _ := json.MarshalIndent(mcpConfig, "", "  ")
	files[".cursor/mcp.json"] = string(mcpJSON)
}

// generateMCPServers creates the MCP server configurations
func (ci *CursorIntegration) generateMCPServers() map[string]interface{} {
	servers := make(map[string]interface{})
	
	// Core project MCP server
	servers["codekeeper-project"] = map[string]interface{}{
		"command": "node",
		"args":    []string{"scripts/mcp/project-server.js"},
		"env": map[string]string{
			"PROJECT_NAME":   ci.spec.Name,
			"PROJECT_DOMAIN": ci.spec.Domain,
			"CORE_ENTITY":    ci.spec.CoreEntity,
		},
	}
	
	// Git operations MCP server
	servers["codekeeper-git"] = map[string]interface{}{
		"command": "node",
		"args":    []string{"scripts/mcp/git-server.js"},
		"env": map[string]string{
			"PROJECT_ROOT": ".",
		},
	}
	
	// Domain-specific MCP server
	if domainServer := ci.getDomainMCPServer(); domainServer != nil {
		servers["codekeeper-domain"] = domainServer
	}
	
	return servers
}

// getDomainMCPServer returns domain-specific MCP server configuration
func (ci *CursorIntegration) getDomainMCPServer() map[string]interface{} {
	switch ci.spec.Domain {
	case "fintech":
		return map[string]interface{}{
			"command": "node",
			"args":    []string{"scripts/mcp/compliance-server.js"},
			"env": map[string]string{
				"COMPLIANCE_TYPE": "fintech",
				"PCI_DSS_MODE":    "enabled",
				"SOX_COMPLIANCE":  "enabled",
			},
		}
	case "healthcare":
		return map[string]interface{}{
			"command": "node",
			"args":    []string{"scripts/mcp/hipaa-server.js"},
			"env": map[string]string{
				"COMPLIANCE_TYPE": "healthcare",
				"HIPAA_MODE":      "enabled",
				"PHI_PROTECTION":  "enabled",
			},
		}
	case "ecommerce":
		return map[string]interface{}{
			"command": "node",
			"args":    []string{"scripts/mcp/ecommerce-server.js"},
			"env": map[string]string{
				"COMPLIANCE_TYPE": "ecommerce",
				"PCI_MODE":        "enabled",
				"GDPR_COMPLIANCE": "enabled",
			},
		}
	default:
		// For general domain, use the compliance server with general settings
		return map[string]interface{}{
			"command": "node",
			"args":    []string{"scripts/mcp/compliance-server.js"},
			"env": map[string]string{
				"COMPLIANCE_TYPE": "general",
			},
		}
	}
}

// generateCursorRules creates Cursor coding rules using proper .mdc format
func (ci *CursorIntegration) generateCursorRules(files map[string]string) {
	// Main project rules - always applied
	files[".cursor/rules/project-standards.mdc"] = ci.generateProjectStandardsRule()
	
	// Domain-specific rules - auto-attached based on file patterns
	files[".cursor/rules/domain-rules.mdc"] = ci.generateDomainSpecificRule()
	
	// Backend-specific rules - auto-attached for backend files
	files[".cursor/rules/backend-patterns.mdc"] = ci.generateBackendRule()
	
	// Frontend-specific rules - auto-attached for frontend files  
	files[".cursor/rules/frontend-patterns.mdc"] = ci.generateFrontendRule()
	
	// Security rules - always applied
	files[".cursor/rules/security-standards.mdc"] = ci.generateSecurityRule()
	
	// Testing rules - auto-attached for test files
	files[".cursor/rules/testing-standards.mdc"] = ci.generateTestingRule()
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

// New MDC rule generation methods

func (ci *CursorIntegration) generateProjectStandardsRule() string {
	return `---
description: Core project standards and architectural patterns
globs: ["**/*"]
alwaysApply: true
---

# ` + ci.spec.Name + ` - Project Standards

## Project Context
- **Domain**: ` + ci.spec.Domain + `
- **Core Entity**: ` + ci.spec.CoreEntity + `
- **Backend**: ` + ci.spec.Backend + `
- **API Style**: ` + ci.spec.APIStyle + `

## Architectural Principles
1. Follow clean architecture patterns defined in docs/ARCHITECTURE.md
2. Implement domain-driven design principles
3. Use proper separation of concerns
4. Follow SOLID principles
5. Implement proper error handling throughout the application

## Code Quality Standards
- Write self-documenting code with meaningful names
- Use TypeScript strict mode for all frontend code
- Implement proper interfaces and type definitions
- Follow established project structure in docs/frontend/PROJECT_STRUCTURE.md
- Maintain consistency with existing codebase patterns

## File Organization
- Components: PascalCase (UserProfile.tsx)
- Hooks: camelCase with 'use' prefix (useUserData.ts)
- Utils: camelCase (formatDate.ts)
- Constants: UPPER_SNAKE_CASE
- Config files: kebab-case

## Git Workflow
- Use conventional commit messages (feat:, fix:, docs:, etc.)
- Create feature branches for new development
- Keep commits atomic and focused
- Write meaningful commit messages that explain the "why"

Generated by AI CodeKeeper v2.1.0 for ` + ci.spec.Domain + ` domain.`
}

func (ci *CursorIntegration) generateDomainSpecificRule() string {
	var domainPatterns, domainGuidelines string
	
	switch ci.spec.Domain {
	case "fintech":
		domainPatterns = `["**/*financial*", "**/*payment*", "**/*transaction*", "**/*money*", "**/*currency*"]`
		domainGuidelines = `## Fintech Domain Rules

### Critical Requirements
- **NEVER use floating-point arithmetic for money calculations**
- Use Decimal.js, big.js, or similar libraries for precise calculations
- Implement comprehensive audit logging for all financial operations
- Follow PCI-DSS compliance requirements for payment data
- Use strong encryption for sensitive financial information

### Compliance Standards
- PCI-DSS: Payment Card Industry Data Security Standard
- SOX: Sarbanes-Oxley Act compliance
- KYC: Know Your Customer procedures
- AML: Anti-Money Laundering checks
- GDPR: General Data Protection Regulation

### Implementation Patterns
- Always validate monetary amounts with proper precision
- Implement idempotency for financial transactions
- Use proper authorization for financial operations
- Implement rate limiting and fraud detection
- Maintain complete audit trails for regulatory compliance

### Forbidden Patterns
- Never use JavaScript Number type for monetary calculations
- Never log sensitive financial information (account numbers, SSNs)
- Never skip validation for financial transactions
- Never implement financial logic without audit trails`

	case "healthcare":
		domainPatterns = `["**/*patient*", "**/*medical*", "**/*health*", "**/*clinical*", "**/*hipaa*"]`
		domainGuidelines = `## Healthcare Domain Rules

### HIPAA Compliance Requirements
- Encrypt all PHI (Protected Health Information) at rest and in transit
- Implement proper access controls with role-based permissions
- Maintain comprehensive audit logs for all data access
- Use FHIR standards for healthcare data interoperability
- Implement proper consent management workflows

### Security Standards
- Use strong encryption (AES-256) for sensitive medical data
- Implement proper authentication with MFA when possible
- Follow principle of least privilege for data access
- Use secure communication protocols (TLS 1.3+)
- Implement proper data retention and deletion policies

### Clinical Data Handling
- Use standardized medical codes (ICD-10, CPT, SNOMED)
- Implement proper data validation for clinical information
- Follow HL7 FHIR standards for interoperability
- Ensure data integrity for clinical decision support
- Implement proper versioning for medical records

### Forbidden Patterns
- Never log PHI in application logs
- Never store unencrypted medical data
- Never skip consent validation for data access
- Never implement medical logic without proper validation`

	case "ecommerce":
		domainPatterns = `["**/*product*", "**/*cart*", "**/*order*", "**/*payment*", "**/*inventory*"]`
		domainGuidelines = `## E-commerce Domain Rules

### Core Business Logic
- Implement proper inventory management and stock tracking
- Use secure payment processing with PCI compliance
- Implement cart persistence across user sessions
- Use proper product catalog management with search/filtering
- Implement comprehensive order management workflows

### Performance Requirements
- Use CDN for static assets and image optimization
- Implement proper caching strategies for product data
- Optimize database queries for large product catalogs
- Use pagination for large result sets
- Implement proper search indexing (Elasticsearch, etc.)

### Security Standards
- Follow PCI-DSS for payment card data
- Implement secure checkout flows
- Use proper session management for shopping carts
- Protect customer data according to GDPR/CCPA
- Implement fraud detection for transactions

### Forbidden Patterns
- Never store payment card data without PCI compliance
- Never implement checkout without proper validation
- Never allow overselling of inventory
- Never skip input validation for user data`

	default:
		domainPatterns = `["**/*"]`
		domainGuidelines = `## General Domain Rules

### Security Best Practices
- Validate all user inputs at application boundaries
- Use HTTPS for all communications
- Implement proper authentication and authorization
- Use environment variables for configuration
- Follow security best practices for your technology stack

### Data Handling
- Implement proper input validation and sanitization
- Use appropriate data types for business logic
- Follow data privacy principles (GDPR compliance)
- Implement proper error handling and user feedback
- Use secure coding practices to prevent common vulnerabilities

### Performance Considerations
- Implement proper caching strategies
- Optimize database queries and indexes
- Use appropriate design patterns for scalability
- Monitor application performance and resource usage
- Implement proper logging and monitoring`
	}

	return `---
description: Domain-specific rules for ` + ci.spec.Domain + ` applications
globs: ` + domainPatterns + `
alwaysApply: false
---

` + domainGuidelines
}

func (ci *CursorIntegration) generateBackendRule() string {
	backendPatterns := `["apps/backend/**/*", "**/*.js", "**/*.ts", "**/*.py", "**/*.go"]`
	
	var backendGuidelines string
	switch ci.spec.Backend {
	case "javascript", "node", "nodejs":
		backendGuidelines = `## Node.js/JavaScript Backend Rules

### Core Patterns
- Use async/await instead of callbacks for asynchronous operations
- Implement proper error middleware for Express.js applications
- Use environment variables for all configuration
- Follow Express.js best practices and security guidelines
- Implement request validation using joi, zod, or similar libraries

### Architecture
- Use proper MVC or clean architecture patterns
- Implement service layer for business logic
- Use repository pattern for data access
- Implement proper dependency injection
- Use middleware for cross-cutting concerns

### Error Handling
- Always use proper error handling with try/catch blocks
- Implement global error handlers
- Use appropriate HTTP status codes
- Log errors with sufficient context for debugging
- Never expose internal error details to clients

### Performance
- Use connection pooling for databases
- Implement proper caching strategies (Redis, in-memory)
- Use streaming for large data operations
- Implement proper pagination for large datasets
- Monitor and optimize database queries`

	case "python":
		backendGuidelines = `## Python Backend Rules

### Code Style
- Follow PEP 8 style guidelines strictly
- Use type hints for all function signatures
- Implement proper exception handling patterns
- Use virtual environments for dependency management
- Follow FastAPI or Django best practices

### Architecture
- Use proper MVC or clean architecture patterns
- Implement service layer pattern for business logic
- Use dependency injection where appropriate
- Follow SOLID principles in class design
- Use proper async/await patterns for I/O operations

### Error Handling
- Use specific exception types instead of generic Exception
- Implement proper logging with structured formats
- Use context managers for resource management
- Handle database connections properly
- Implement proper validation using Pydantic or similar

### Performance
- Use appropriate data structures for performance
- Implement proper database connection pooling
- Use caching strategies (Redis, memcached)
- Profile code for performance bottlenecks
- Use async frameworks for high-concurrency applications`

	case "go":
		backendGuidelines = `## Go Backend Rules

### Go Idioms
- Follow effective Go principles and idioms
- Use proper error handling patterns (no exceptions)
- Implement context cancellation for long-running operations
- Use interfaces for abstraction and testing
- Write idiomatic Go code following community standards

### Architecture
- Use clean architecture or hexagonal architecture
- Implement proper dependency injection
- Use proper package organization
- Follow Go project layout standards
- Use appropriate design patterns

### Error Handling
- Always check and handle errors explicitly
- Use error wrapping for context (fmt.Errorf with %w verb)
- Implement proper logging with structured formats
- Use context for cancellation and timeouts
- Handle panics appropriately with recover()

### Performance
- Use proper goroutine patterns
- Implement channels for communication
- Use sync package for proper synchronization
- Profile applications with pprof
- Optimize memory allocations and garbage collection`

	default:
		backendGuidelines = `## General Backend Rules

### Architecture
- Follow language-specific best practices
- Implement proper separation of concerns
- Use appropriate design patterns
- Follow clean code principles
- Implement proper testing strategies

### Error Handling
- Implement comprehensive error handling
- Use appropriate logging levels and formats
- Handle edge cases and failure scenarios
- Implement proper retry mechanisms
- Use circuit breaker patterns for external services

### Performance
- Optimize database queries and connections
- Implement proper caching strategies
- Monitor application performance
- Use appropriate concurrency patterns
- Implement proper resource management`
	}

	return `---
description: Backend development patterns and standards
globs: ` + backendPatterns + `
alwaysApply: false
---

` + backendGuidelines
}

func (ci *CursorIntegration) generateFrontendRule() string {
	return `---
description: Frontend React TypeScript development standards
globs: ["apps/frontend/**/*", "**/*.tsx", "**/*.jsx", "**/*.css", "**/*.scss"]
alwaysApply: false
---

# Frontend Development Rules

## React TypeScript Patterns

### Component Development
- Use functional components with hooks exclusively
- Implement proper TypeScript interfaces for all props
- Follow PascalCase naming convention for components
- Use proper component composition patterns
- Implement error boundaries for robust error handling

### State Management
- Use React hooks (useState, useEffect, useContext) appropriately
- Implement custom hooks for complex logic reuse
- Use proper dependency arrays in useEffect
- Avoid prop drilling - use Context API when needed
- Implement proper state normalization for complex data

### TypeScript Standards
- Use strict TypeScript configuration
- Define proper interfaces for all data structures
- Use union types and type guards appropriately
- Implement proper generic types where beneficial
- Avoid 'any' type - use proper typing

### Performance Optimization
- Use React.memo for expensive components
- Implement proper key props for list items
- Use useMemo and useCallback for expensive operations
- Implement code splitting with React.lazy
- Optimize bundle size with proper tree shaking

### Styling Standards
- Use CSS modules or styled-components for component styling
- Follow BEM methodology for CSS class naming
- Implement responsive design patterns
- Use CSS custom properties for theming
- Maintain consistent spacing and typography scales

### Testing Patterns
- Write unit tests for components using React Testing Library
- Test user interactions and behavior, not implementation
- Use proper test data and mocking strategies
- Implement accessibility testing
- Write integration tests for complex user flows

### Accessibility (a11y)
- Use semantic HTML elements appropriately
- Implement proper ARIA labels and roles
- Ensure keyboard navigation support
- Maintain proper color contrast ratios
- Test with screen readers and accessibility tools`
}

func (ci *CursorIntegration) generateSecurityRule() string {
	return `---
description: Security standards and best practices
globs: ["**/*"]
alwaysApply: true
---

# Security Standards

## Critical Security Rules

### Data Protection
- **NEVER log sensitive information** (passwords, tokens, PII, API keys)
- Use environment variables for all configuration and secrets
- Implement proper input validation at all boundaries
- Use parameterized queries to prevent SQL injection
- Implement proper authentication and authorization

### Communication Security
- Use HTTPS/TLS for all communications
- Implement proper CORS policies
- Use secure headers (CSP, HSTS, X-Frame-Options)
- Validate and sanitize all user inputs
- Implement proper session management

### Code Security
- Never hardcode secrets, API keys, or passwords
- Use proper cryptographic libraries for encryption
- Implement secure random number generation
- Follow principle of least privilege
- Use proper error handling without information disclosure

### Infrastructure Security
- Use secure defaults for all configurations
- Implement proper secrets management
- Use container security best practices
- Implement proper logging and monitoring
- Follow security scanning and vulnerability management

## Domain-Specific Security

` + ci.getDomainSecurityRules() + `

## Security Checklist
- [ ] Input validation implemented
- [ ] Authentication and authorization working
- [ ] Sensitive data properly encrypted
- [ ] Security headers configured
- [ ] Error handling doesn't leak information
- [ ] Dependencies scanned for vulnerabilities
- [ ] Secrets managed properly
- [ ] Logging excludes sensitive data`
}

func (ci *CursorIntegration) generateTestingRule() string {
	return `---
description: Testing standards and best practices
globs: ["**/*.test.*", "**/*.spec.*", "**/tests/**/*", "**/__tests__/**/*"]
alwaysApply: false
---

# Testing Standards

## Testing Philosophy
- Test behavior, not implementation details
- Write tests that provide confidence in the system
- Use descriptive test names that explain the scenario
- Follow the testing pyramid (unit > integration > e2e)
- Maintain high test coverage for critical business logic

## Unit Testing
- Test individual functions and components in isolation
- Use proper mocking for external dependencies
- Test edge cases and error conditions
- Use arrange-act-assert (AAA) pattern
- Keep tests fast and independent

## Integration Testing
- Test component interactions and API endpoints
- Use realistic test data and scenarios
- Test authentication and authorization flows
- Verify data persistence and retrieval
- Test error handling and recovery

## End-to-End Testing
- Test critical user journeys and workflows
- Use page object patterns for maintainability
- Test across different browsers and devices
- Use stable selectors (data-testid attributes)
- Keep E2E tests focused and reliable

## Test Data Management
- Use factories or builders for test data creation
- Implement proper test database setup and teardown
- Use realistic but anonymized data
- Follow data privacy guidelines
- Implement proper test isolation

## Domain-Specific Testing

` + ci.getDomainTestingGuidelines() + `

## Testing Best Practices
- Run tests in CI/CD pipeline
- Maintain test coverage reports
- Review and update tests with code changes
- Use property-based testing where appropriate
- Implement performance and security testing`
}

func (ci *CursorIntegration) getDomainSecurityRules() string {
	switch ci.spec.Domain {
	case "fintech":
		return `### Financial Services Security
- Follow PCI-DSS requirements for payment card data
- Implement strong encryption for financial data (AES-256)
- Use proper tokenization for sensitive payment information
- Implement fraud detection and prevention mechanisms
- Follow SOX compliance for financial reporting
- Use proper audit trails for all financial operations
- Implement rate limiting for financial APIs
- Use secure communication protocols (TLS 1.3+)`

	case "healthcare":
		return `### Healthcare Security (HIPAA)
- Encrypt all PHI at rest and in transit
- Implement proper access controls and audit logging
- Use strong authentication (preferably MFA)
- Follow HIPAA Security Rule requirements
- Implement proper data retention and disposal
- Use de-identification techniques for analytics
- Implement proper consent management
- Follow FHIR security guidelines for interoperability`

	case "ecommerce":
		return `### E-commerce Security
- Follow PCI-DSS for payment processing
- Implement secure checkout flows
- Protect customer data according to GDPR/CCPA
- Use proper session management for shopping carts
- Implement fraud detection for transactions
- Secure product catalog and inventory data
- Use proper authentication for user accounts
- Implement secure password policies`

	default:
		return `### General Security
- Follow OWASP Top 10 security guidelines
- Implement proper authentication and authorization
- Use secure coding practices
- Follow data protection regulations (GDPR)
- Implement proper logging and monitoring
- Use security scanning tools
- Follow secure development lifecycle (SDLC)
- Implement incident response procedures`
	}
}

func (ci *CursorIntegration) getDomainTestingGuidelines() string {
	switch ci.spec.Domain {
	case "fintech":
		return `### Financial Services Testing
- Test monetary calculations with precise decimal arithmetic
- Verify audit trails for all financial operations
- Test compliance workflows (KYC, AML procedures)
- Validate encryption of sensitive financial data
- Test fraud detection mechanisms
- Verify regulatory reporting accuracy
- Test transaction idempotency
- Validate PCI-DSS compliance in payment flows`

	case "healthcare":
		return `### Healthcare Testing
- Test PHI encryption and access controls
- Verify HIPAA compliance in data handling
- Test consent management workflows
- Validate clinical decision support accuracy
- Test FHIR interoperability
- Verify audit logging for data access
- Test emergency access (break-glass) procedures
- Validate data anonymization for analytics`

	case "ecommerce":
		return `### E-commerce Testing
- Test inventory management and stock tracking
- Verify payment processing security
- Test shopping cart persistence
- Validate order fulfillment workflows
- Test search and filtering performance
- Verify product catalog management
- Test recommendation algorithms
- Validate customer data protection`

	default:
		return `### General Testing
- Test security controls and authentication
- Verify data validation and sanitization
- Test error handling and recovery
- Validate performance under load
- Test integration with external services
- Verify logging and monitoring
- Test backup and disaster recovery
- Validate compliance with regulations`
	}
}