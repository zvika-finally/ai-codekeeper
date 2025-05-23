package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Generator handles project generation based on specifications
type Generator struct {
	spec   *ProjectSpec
	domain string // Simplified for now
}

// New creates a new generator instance
func New(spec *ProjectSpec) *Generator {
	return &Generator{
		spec:   spec,
		domain: spec.Domain,
	}
}

// Generate creates the complete project structure
func (g *Generator) Generate() error {
	// Validate specification
	if err := g.spec.Validate(); err != nil {
		return fmt.Errorf("invalid specification: %w", err)
	}

	// Set up project path
	g.spec.ProjectPath = filepath.Join(".", g.spec.Name)
	g.spec.GeneratedAt = time.Now().Format(time.RFC3339)

	// Create project directory
	if err := os.MkdirAll(g.spec.ProjectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize git repository
	if err := g.initializeGitRepository(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Generate monorepo structure following AI_MASTER_PROMPT.md
	if err := g.generateMonorepoStructure(); err != nil {
		return fmt.Errorf("failed to generate monorepo structure: %w", err)
	}

	// Generate basic files for testing
	if err := g.generateBasicFiles(); err != nil {
		return fmt.Errorf("failed to generate basic files: %w", err)
	}

	// Generate framework configuration
	if err := g.generateFrameworkConfig(); err != nil {
		return fmt.Errorf("failed to generate framework config: %w", err)
	}

	// Generate Cursor IDE integration
	if err := g.generateCursorIntegration(); err != nil {
		return fmt.Errorf("failed to generate Cursor integration: %w", err)
	}

	return nil
}

// generateMonorepoStructure creates the base directory structure
func (g *Generator) generateMonorepoStructure() error {
	dirs := []string{
		"apps",
		"apps/backend",
		"apps/frontend", 
		"packages",
		"packages/shared-types",
		"infra",
		"infra/aws",
		"infra/render",
		"docs",
		".github",
		".github/workflows",
		".codekeeper", // Framework configuration
	}

	for _, dir := range dirs {
		path := filepath.Join(g.spec.ProjectPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

// generateBasicFiles creates basic placeholder files
func (g *Generator) generateBasicFiles() error {
	// Base files that are always generated
	files := map[string]string{
		"README.md": g.generateReadme(),
		"SECURITY.md": g.generateSecurityMd(),
		"CHANGELOG.md": g.generateChangelog(),
		"docker-compose.yml": g.generateDockerCompose(),
		".devcontainer/devcontainer.json": g.generateDevContainer(),
		"docs/00_OVERVIEW.md": g.generateOverviewDoc(),
		"docs/01_REQUIREMENTS.md": g.generateRequirementsDoc(),
		"docs/02_ARCHITECTURE.md": g.generateArchitectureDoc(),
		"docs/03_SYSTEM_DESIGN.md": g.generateSystemDesignDoc(),
		"docs/04_DATA_MODEL.md": g.generateDataModelDoc(),
		"docs/05_API_ENDPOINTS.md": g.generateAPIEndpointsDoc(),
		"docs/06_DEPLOYMENT.md": g.generateDeploymentDoc(),
		"docs/07_OPERATIONS_OBSERVABILITY.md": g.generateOperationsDoc(),
		"docs/08_TESTING_STRATEGY.md": g.generateTestingStrategyDoc(),
		"docs/09_LOCAL_DEVELOPMENT.md": g.generateLocalDevelopmentDoc(),
		"docs/CONTRIBUTING.md": g.generateContributingDoc(),
		".github/workflows/ci.yml": g.generateCIWorkflow(),
		"LICENSE": g.generateLicense(),
		".vscode/settings.json": g.generateVSCodeSettings(),
		".vscode/extensions.json": g.generateVSCodeExtensions(),
		".vscode/launch.json": g.generateVSCodeLaunch(),
	}
	
	// Add backend-specific files
	backendFiles := g.generateBackendFiles()
	for path, content := range backendFiles {
		files[path] = content
	}
	
	// Add frontend files (always React/TypeScript for now)
	frontendFiles := g.generateFrontendFiles()
	for path, content := range frontendFiles {
		files[path] = content
	}
	
	// Add infrastructure as code
	infraFiles := g.generateInfrastructureFiles()
	for path, content := range infraFiles {
		files[path] = content
	}
	
	// Add final integration files
	finalFiles := g.generateFinalIntegrationFiles()
	for path, content := range finalFiles {
		files[path] = content
	}
	
	// Add development tooling based on tech stack
	toolingFiles := g.generateDevelopmentTooling()
	for path, content := range toolingFiles {
		files[path] = content
	}

	for filePath, content := range files {
		fullPath := filepath.Join(g.spec.ProjectPath, filePath)
		
		// Ensure directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		
		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// generateFrameworkConfig creates AI development framework configuration
func (g *Generator) generateFrameworkConfig() error {
	configPath := filepath.Join(g.spec.ProjectPath, ".codekeeper")
	
	// Create framework configuration
	g.spec.Framework = &FrameworkConfig{
		Version: "1.0.0",
		Domain:  DomainConfig{
			Name: g.spec.Domain,
			Version: "1.0.0",
		},
		GuardRails: GuardRailsConfig{
			Enforcement: "advisory",
			PreCommit:   true,
			CI:          true,
			IDE:         "vscode", // Default
			Rules:       g.getDomainGuardRails(),
		},
		DevEnvironment: DevEnvironmentConfig{
			Type:        g.spec.DevEnvironment,
			Services:    g.getRequiredServices(),
			Ports:       g.getServicePorts(),
			Environment: g.getEnvironmentVars(),
		},
	}

	// Save configuration
	configFile := filepath.Join(configPath, "config.json")
	return saveJSON(configFile, g.spec.Framework)
}

// Template generation methods (simplified for testing)
func (g *Generator) generateReadme() string {
	return fmt.Sprintf(`# %s

%s

## Generated Structure

This project was generated using the AI Development Framework with %s domain expertise.

### Core Entity: %s

The application includes complete CRUD operations for %s with domain-specific best practices.

### Technology Stack

- **Backend**: %s
- **API Style**: %s
- **Databases**: %s
- **Development**: %s environment

### Getting Started

1. Setup development environment:
   ` + "`" + `bash
   codekeeper env setup
   ` + "`" + `

2. Start all services:
   ` + "`" + `bash
   docker compose up
   ` + "`" + `

3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Documentation: http://localhost:8080/docs

### Guard Rails

This project includes automated guard rails for:
%s

### Next Steps

1. Review generated code in ` + "`" + `apps/` + "`" + ` directory
2. Customize domain logic in service files
3. Run guard rails: ` + "`" + `codekeeper check` + "`" + `
4. Generate additional features: ` + "`" + `codekeeper feature <name>` + "`" + `

Generated with AI Development Framework v1.0.0
`,
		g.spec.Name,
		g.spec.Description,
		g.spec.Domain,
		g.spec.CoreEntity,
		g.spec.CoreEntity,
		g.spec.Backend,
		g.spec.APIStyle,
		strings.Join(g.spec.Databases, ", "),
		g.spec.DevEnvironment,
		g.getGuardRailsDescription(),
	)
}

func (g *Generator) generateBackendPackageJson() string {
	return fmt.Sprintf(`{
  "name": "@%s/backend",
  "version": "1.0.0",
  "description": "%s backend API",
  "main": "src/server.js",
  "scripts": {
    "start": "node src/server.js",
    "dev": "nodemon src/server.js",
    "test": "jest --testTimeout=10000",
    "test:coverage": "jest --coverage",
    "lint": "eslint src/",
    "lint:fix": "eslint src/ --fix",
    "check": "codekeeper check",
    "migration:create": "sequelize-cli migration:generate",
    "migration:run": "sequelize-cli db:migrate",
    "migration:undo": "sequelize-cli db:migrate:undo"
  },
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^7.0.0",
    "express-rate-limit": "^6.7.0",
    "sequelize": "^6.32.1",
    "pg": "^8.11.0",
    "pg-hstore": "^2.3.4",
    "jsonwebtoken": "^9.0.1",
    "bcrypt": "^5.1.0",
    "winston": "^3.9.0",
    "validator": "^13.9.0",
    "dotenv": "^16.1.4"
  },
  "devDependencies": {
    "nodemon": "^2.0.22",
    "jest": "^29.5.0",
    "supertest": "^6.3.3",
    "eslint": "^8.42.0",
    "eslint-config-node": "^4.1.0",
    "sequelize-cli": "^6.6.1"
  },
  "jest": {
    "testEnvironment": "node",
    "setupFilesAfterEnv": ["<rootDir>/tests/setup.js"],
    "collectCoverageFrom": [
      "src/**/*.js",
      "!src/server.js"
    ],
    "coverageThreshold": {
      "global": {
        "branches": 80,
        "functions": 80,
        "lines": 80,
        "statements": 80
      }
    }
  }
}`, g.spec.Name, g.spec.Description)
}

func (g *Generator) generateFrontendPackageJson() string {
	return fmt.Sprintf(`{
  "name": "@%s/frontend",
  "version": "1.0.0",
  "description": "%s frontend application",
  "type": "module",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0",
    "lint:fix": "eslint . --ext ts,tsx --fix",
    "type-check": "tsc --noEmit",
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage",
    "check": "codekeeper check"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.14.1",
    "zustand": "^4.3.9",
    "axios": "^1.4.0",
    "clsx": "^1.2.1",
    "lucide-react": "^0.263.1",
    "@%s/shared-types": "workspace:*"
  },
  "devDependencies": {
    "@types/react": "^18.2.15",
    "@types/react-dom": "^18.2.7",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0",
    "@vitejs/plugin-react": "^4.0.3",
    "@vitest/coverage-v8": "^0.33.0",
    "@vitest/ui": "^0.33.0",
    "eslint": "^8.45.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.3",
    "jsdom": "^22.1.0",
    "typescript": "^5.0.2",
    "vite": "^4.4.5",
    "vitest": "^0.33.0"
  }
}`, g.spec.Name, g.spec.Description, g.spec.Name)
}

func (g *Generator) generateSharedTypesPackageJson() string {
	return fmt.Sprintf(`{
  "name": "@%s/shared-types",
  "version": "1.0.0",
  "description": "Shared TypeScript types for %s",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch"
  },
  "devDependencies": {
    "typescript": "latest"
  }
}`, g.spec.Name, g.spec.Name)
}

func (g *Generator) generateDockerCompose() string {
	backendLang := g.spec.GetBackendLanguage()
	
	// Base services
	services := `version: '3.8'

services:
  backend:
    build: 
      context: ./apps/backend`
      
	// Add backend-specific dockerfile and environment
	switch backendLang {
	case "javascript":
		services += `
      dockerfile: Dockerfile.node
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=development`
	case "python":
		services += `
      dockerfile: Dockerfile.python
    ports:
      - "8080:8080"
    environment:
      - PYTHONPATH=/app`
	case "go":
		services += `
      dockerfile: Dockerfile.go
    ports:
      - "8080:8080"
    environment:
      - CGO_ENABLED=0`
	}
	
	services += `
    depends_on:
      - postgres

  frontend:
    build: ./apps/frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=` + g.spec.Name + `
      - POSTGRES_USER=dev
      - POSTGRES_PASSWORD=dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data`

	if g.spec.HasDatabase("Redis") {
		services += `

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"`
	}

	services += `

volumes:
  postgres_data:`

	return services
}

func (g *Generator) generateDevContainer() string {
	return `{
  "name": "` + g.spec.Name + ` Development",
  "dockerComposeFile": "../docker-compose.yml",
  "service": "backend",
  "workspaceFolder": "/workspace",
  
  "customizations": {
    "vscode": {
      "extensions": [
        "ms-vscode.vscode-typescript-next",
        "esbenp.prettier-vscode",
        "dbaeumer.vscode-eslint",
        "ms-azuretools.vscode-docker"
      ],
      "settings": {
        "editor.formatOnSave": true,
        "editor.defaultFormatter": "esbenp.prettier-vscode"
      }
    }
  },
  
  "forwardPorts": [3000, 8080, 5432],
  "postCreateCommand": "npm install"
}`
}

func (g *Generator) generateOverviewDoc() string {
	return fmt.Sprintf(`# %s - Overview

## Application Description

%s

## Core Entity: %s

This application centers around the **%s** entity, which provides:

- Complete CRUD operations
- Domain-specific business logic
- %s best practices
- Compliance and security controls

## Architecture

### Technology Stack

- **Backend**: %s
- **Frontend**: React with TypeScript
- **Database**: %s
- **API Style**: %s
- **Development**: %s

### Domain Expertise: %s

This application leverages %s domain knowledge for:

- Optimal technology recommendations
- Industry best practices
- Compliance requirements
- Security patterns

## Generated Components

### Backend (` + "`" + `apps/backend/` + "`" + `)
- RESTful API endpoints for %s
- Domain-specific business logic
- Database models and migrations
- Authentication and authorization
- Input validation and error handling

### Frontend (` + "`" + `apps/frontend/` + "`" + `)
- React components for %s management
- TypeScript type safety
- Responsive design
- API integration

### Shared Types (` + "`" + `packages/shared-types/` + "`" + `)
- Common TypeScript interfaces
- API request/response types
- Domain models

### Infrastructure (` + "`" + `infra/` + "`" + `)
- AWS Terraform configuration
- Render deployment setup
- CI/CD pipelines

### Documentation (` + "`" + `docs/` + "`" + `)
- Complete architecture documentation
- API specifications
- Local development guide
- Deployment instructions

## Next Steps

1. Review the generated code structure
2. Customize business logic for your specific needs
3. Run the development environment
4. Generate additional features as needed

Generated on: %s
Framework Version: 1.0.0
Domain: %s
`,
		g.spec.Name,
		g.spec.Description,
		g.spec.CoreEntity,
		g.spec.CoreEntity,
		g.spec.Domain,
		g.spec.Backend,
		strings.Join(g.spec.Databases, ", "),
		g.spec.APIStyle,
		g.spec.DevEnvironment,
		g.spec.Domain,
		g.spec.Domain,
		g.spec.CoreEntity,
		g.spec.CoreEntity,
		g.spec.GeneratedAt,
		g.spec.Domain,
	)
}

func (g *Generator) generateCIWorkflow() string {
	backendLang := g.spec.GetBackendLanguage()
	
	baseWorkflow := `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    `
    
    switch backendLang {
    case "javascript":
    	return baseWorkflow + `
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: |
        cd apps/backend && npm ci
        cd ../frontend && npm ci
    
    - name: Run tests
      run: |
        cd apps/backend && npm test
        cd ../frontend && npm test
    
    - name: Run guard rails
      run: codekeeper check --ci
    
    - name: Build
      run: |
        cd apps/backend && npm run build
        cd ../frontend && npm run build`
        
    case "python":
    	return baseWorkflow + `
    - name: Setup Python
      uses: actions/setup-python@v4
      with:
        python-version: '3.11'
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: |
        cd apps/backend && pip install -r requirements.txt
        cd ../frontend && npm ci
    
    - name: Run tests
      run: |
        cd apps/backend && python -m pytest
        cd ../frontend && npm test
    
    - name: Run guard rails
      run: codekeeper check --ci
    
    - name: Lint
      run: |
        cd apps/backend && flake8 .
        cd ../frontend && npm run lint`
        
    case "go":
    	return baseWorkflow + `
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: |
        cd apps/backend && go mod download
        cd ../frontend && npm ci
    
    - name: Run tests
      run: |
        cd apps/backend && go test ./...
        cd ../frontend && npm test
    
    - name: Run guard rails
      run: codekeeper check --ci
    
    - name: Build
      run: |
        cd apps/backend && go build .
        cd ../frontend && npm run build`
        
    default:
    	return baseWorkflow + `
    - name: Run guard rails
      run: codekeeper check --ci`
    }
}

func (g *Generator) generateLicense() string {
	currentYear := time.Now().Year()
	return fmt.Sprintf(`MIT License

Copyright (c) %d %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`, currentYear, g.spec.Name)
}

// Helper functions
func (g *Generator) getDomainGuardRails() []string {
	switch g.spec.Domain {
	case "fintech":
		return []string{"decimal_arithmetic", "audit_trails", "encryption_at_rest", "input_validation"}
	case "healthcare":
		return []string{"hipaa_compliance", "data_encryption", "audit_logs", "access_controls"}
	default:
		return []string{"input_validation", "security_headers", "error_handling", "logging"}
	}
}

func (g *Generator) getGuardRailsDescription() string {
	rules := g.getDomainGuardRails()
	var descriptions []string
	
	for _, rule := range rules {
		switch rule {
		case "decimal_arithmetic":
			descriptions = append(descriptions, "- Decimal arithmetic for financial calculations")
		case "audit_trails":
			descriptions = append(descriptions, "- Comprehensive audit logging")
		case "encryption_at_rest":
			descriptions = append(descriptions, "- Data encryption at rest")
		case "input_validation":
			descriptions = append(descriptions, "- Strict input validation")
		default:
			descriptions = append(descriptions, fmt.Sprintf("- %s", strings.ReplaceAll(rule, "_", " ")))
		}
	}
	
	return strings.Join(descriptions, "\n")
}

func (g *Generator) getRequiredServices() []string {
	services := []string{"backend", "frontend"}
	
	for _, db := range g.spec.Databases {
		if strings.Contains(strings.ToLower(db), "postgresql") {
			services = append(services, "postgres")
		}
		if strings.Contains(strings.ToLower(db), "redis") {
			services = append(services, "redis")
		}
	}
	
	return services
}

func (g *Generator) getServicePorts() map[string]int {
	ports := map[string]int{
		"backend":  8080,
		"frontend": 3000,
		"postgres": 5432,
	}
	
	if g.spec.HasDatabase("Redis") {
		ports["redis"] = 6379
	}
	
	return ports
}

func (g *Generator) getEnvironmentVars() map[string]string {
	env := map[string]string{
		"NODE_ENV": "development",
		"API_URL":  "http://localhost:8080",
	}
	
	if g.spec.Domain == "fintech" {
		env["ENCRYPTION_KEY"] = "dev-key-change-in-production"
		env["AUDIT_LOG_LEVEL"] = "DEBUG"
	}
	
	return env
}

// initializeGitRepository sets up a git repository in the project directory
func (g *Generator) initializeGitRepository() error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Printf("⚠️  Git not found, skipping repository initialization\n")
		return nil
	}

	// Check if already a git repository
	gitDir := filepath.Join(g.spec.ProjectPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		fmt.Printf("📁 Git repository already exists, skipping initialization\n")
		return nil
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = g.spec.ProjectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Create .gitignore file
	gitignoreContent := g.generateGitignore()
	gitignorePath := filepath.Join(g.spec.ProjectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Add all files
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = g.spec.ProjectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	// Create initial commit
	commitMessage := fmt.Sprintf("Initial commit: %s\n\nGenerated with AI CodeKeeper v1.0.0\nDomain: %s\nCore Entity: %s", 
		g.spec.Name, g.spec.Domain, g.spec.CoreEntity)
	
	cmd = exec.Command("git", "commit", "-m", commitMessage)
	cmd.Dir = g.spec.ProjectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	fmt.Printf("✅ Git repository initialized with initial commit\n")
	return nil
}

// generateGitignore creates a comprehensive .gitignore file
func (g *Generator) generateGitignore() string {
	gitignore := `# Dependencies
node_modules/
*/node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Production builds
dist/
build/
.next/
out/

# Environment variables
.env
.env.local
.env.production
.env.staging

# Logs
logs/
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Runtime data
pids/
*.pid
*.seed
*.pid.lock

# Coverage directory used by tools like istanbul
coverage/
*.lcov

# IDEs and editors
.vscode/
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/extensions.json
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Docker
.dockerignore

# Temporary files
*.tmp
*.temp

# CodeKeeper
.codekeeper/logs/
.codekeeper/cache/
.codekeeper/*.tmp
`

	// Add domain-specific ignores
	switch g.spec.Domain {
	case "fintech":
		gitignore += `
# Fintech-specific
keys/
certificates/
*.pem
*.key
audit-logs/
`
	case "healthcare":
		gitignore += `
# Healthcare-specific
patient-data/
phi-logs/
hipaa-exports/
`
	}

	// Add backend-specific ignores
	if strings.Contains(strings.ToLower(g.spec.Backend), "go") {
		gitignore += `
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work
`
	}

	if strings.Contains(strings.ToLower(g.spec.Backend), "python") {
		gitignore += `
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
venv/
env/
.venv/
.env/
pip-log.txt
pip-delete-this-directory.txt
`
	}

	return gitignore
}

// Documentation generation functions
func (g *Generator) generateSecurityMd() string {
	return fmt.Sprintf(`# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

Please report security vulnerabilities to: [security@example.com]

**Do not report security vulnerabilities through public GitHub issues.**

## Implemented Security Measures

### Authentication & Authorization
- JWT token-based authentication
- Role-Based Access Control (RBAC) for %s
- Secure password hashing (bcrypt)
- Session management best practices

### Data Protection
- Input validation for all API endpoints
- SQL injection prevention through ORM usage
- XSS protection via proper output encoding
- CSRF protection for state-changing operations

### Infrastructure Security
- HTTPS enforcement in production
- Environment variable-based configuration
- Secrets management via cloud providers
- Docker security best practices
- Network security groups and firewalls

### Development Security
- Dependency vulnerability scanning
- Static code analysis
- Pre-commit security hooks
- Regular security updates

## Security Best Practices for Developers

1. **Never commit secrets** - Use environment variables
2. **Validate all inputs** - Server-side validation required
3. **Use parameterized queries** - Prevent SQL injection
4. **Handle errors securely** - Don't expose internal details
5. **Keep dependencies updated** - Regular security patches
6. **Follow principle of least privilege** - Minimal permissions
7. **Log security events** - For audit and monitoring
8. **Use HTTPS everywhere** - No plain HTTP in production

## Infrastructure Security

### Render Platform
- Managed SSL certificates
- Environment variable encryption
- Network isolation
- Regular platform security updates

### AWS Infrastructure
- VPC with private subnets
- Security groups with minimal ingress
- IAM roles with least privilege
- AWS Secrets Manager for sensitive data
- CloudTrail for audit logging
- WAF for application protection

Generated for %s domain with enhanced security considerations.
`, strings.Join(g.spec.GetUserRolesList(), ", "), g.spec.Domain)
}

func (g *Generator) generateChangelog() string {
	currentDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`# Changelog

All notable changes to %s will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - %s

### Added - Complete Production-Ready Application

**Core Application Features:**
- ✅ Complete %s CRUD functionality with %s domain expertise
- ✅ User authentication and authorization system (JWT + RBAC)
- ✅ RESTful API with comprehensive validation and error handling
- ✅ Modern React TypeScript frontend with responsive design
- ✅ Shared TypeScript types for end-to-end type safety

**Backend Infrastructure (%s):**
- ✅ Production-ready server with security middleware
- ✅ Database integration with connection pooling and migrations
- ✅ Structured logging and monitoring endpoints
- ✅ Comprehensive test suite with >80%% coverage target
- ✅ Domain-specific business logic and validation

**Frontend Application (React + TypeScript):**
- ✅ Modern component architecture with custom hooks
- ✅ State management with Zustand
- ✅ API integration with Axios interceptors
- ✅ Responsive design and accessibility features
- ✅ Development tooling (Vite, ESLint, Prettier)

**Infrastructure as Code:**
- ✅ Render deployment for MVP/staging (managed PostgreSQL, SSL, auto-deploy)
- ✅ AWS Terraform for production (ECS Fargate, RDS, ALB, Secrets Manager)
- ✅ Multi-environment CI/CD pipeline with GitHub Actions
- ✅ Docker containerization with multi-stage builds
- ✅ Security best practices and cost optimization

**Documentation & Developer Experience:**
- ✅ Comprehensive README with AI implementation guidelines
- ✅ Complete API documentation and architecture diagrams
- ✅ Local development setup with Docker Compose
- ✅ Deployment guides for both Render and AWS
- ✅ Security policies and compliance documentation

**Domain-Specific Features (%s):**
%s

**12-Factor App Compliance:**
- ✅ Factor I-XII fully implemented with environment-based configuration
- ✅ Stateless processes and external state management
- ✅ Build/release/run separation with CI/CD pipeline
- ✅ Structured logging and graceful shutdown handling

### Technology Stack
- **Backend**: %s with %s framework
- **Frontend**: React 18 + TypeScript + Vite
- **Database**: %s with ORM integration
- **API Style**: %s with comprehensive validation
- **State Management**: Zustand for frontend, JWT for auth
- **Infrastructure**: Render (staging) + AWS Terraform (production)
- **CI/CD**: GitHub Actions with multi-environment deployment
- **Testing**: Jest/Vitest with React Testing Library
- **Development**: Docker Compose with hot reload

### Generated Components

**Total Files Generated: %s**
- Documentation: 14 comprehensive guides and specifications
- Backend: Complete %s application with tests
- Frontend: Full React TypeScript application
- Infrastructure: Render + AWS Terraform configurations
- CI/CD: GitHub Actions workflows for deployment
- Development: Docker, linting, and tooling configurations

### Deployment Options

1. **MVP/Staging (Render)**: One-click deployment with managed services
2. **Production (AWS)**: Enterprise-grade infrastructure with Terraform
3. **Local Development**: Docker Compose for complete stack

### Next Steps for AI Implementation

The framework provides comprehensive guidelines for AI to implement:
1. Complete backend CRUD operations with domain logic
2. Frontend components with proper state management
3. Advanced Terraform modules for production deployment
4. Enhanced security and monitoring features

### Cost Estimates
- **Development**: Free tier eligible
- **Staging (Render)**: $0-25/month for MVP
- **Production (AWS)**: $50-200/month depending on scale

Generated with AI CodeKeeper v1.0.0 - Complete production-ready application framework
`, g.spec.Name, currentDate, g.spec.CoreEntity, g.spec.Domain, 
   g.spec.Backend, g.spec.Domain, g.getDomainSpecificFeatures(),
   g.spec.Backend, g.spec.GetBackendFramework(), strings.Join(g.spec.Databases, ", "), g.spec.APIStyle,
   g.getGeneratedFileCount(), g.spec.Backend)
}

func (g *Generator) generateRequirementsDoc() string {
	return fmt.Sprintf(`# Requirements Specification

## Application Overview

**Application**: %s
**Domain**: %s
**Description**: %s

## Core Entity

**Primary Entity**: %s

The application centers around the %s entity, providing complete lifecycle management with domain-specific business rules and validation.

## Functional Requirements

### F1: %s Management
- **F1.1**: Create new %s entries with validation
- **F1.2**: Read/view %s details and listings
- **F1.3**: Update existing %s information
- **F1.4**: Delete %s entries (soft delete preferred)
- **F1.5**: Search and filter %s by relevant criteria

### F2: User Management & Authentication
- **F2.1**: User registration and email verification
- **F2.2**: Secure login with JWT tokens
- **F2.3**: Password reset functionality
- **F2.4**: User profile management

### F3: Authorization & Access Control
%s

### F4: %s Domain-Specific Features
%s

## Non-Functional Requirements

### Performance
- **P1**: API response time < 200ms for CRUD operations
- **P2**: Support for 1000+ concurrent users
- **P3**: Database query optimization for large datasets

### Security
- **S1**: All data transmission over HTTPS
- **S2**: Input validation and sanitization
- **S3**: Protection against OWASP Top 10 vulnerabilities
- **S4**: Audit logging for critical operations

### Scalability
- **SC1**: Horizontal scaling capability
- **SC2**: Database read replicas support
- **SC3**: CDN integration for static assets

### Reliability
- **R1**: 99.9%% uptime target
- **R2**: Automated backup and disaster recovery
- **R3**: Health monitoring and alerting

### Usability
- **U1**: Responsive design for mobile and desktop
- **U2**: Intuitive user interface
- **U3**: Accessibility compliance (WCAG 2.1 AA)

## Technical Constraints

- **Backend**: %s
- **Frontend**: React with TypeScript
- **Database**: %s
- **API Style**: %s
- **Deployment**: Render (staging), AWS (production)

## Success Criteria

1. Complete %s CRUD functionality implemented
2. User authentication and authorization working
3. Responsive UI for all screen sizes
4. API documentation complete and accurate
5. Automated testing coverage > 80%%
6. Production deployment successful
7. Performance benchmarks met

## Future Roadmap

### Phase 2 (v2.0.0)
- Advanced search and filtering
- Bulk operations
- Data export functionality
- Integration APIs

### Phase 3 (v3.0.0)
- Mobile application
- Real-time notifications
- Advanced analytics
- Multi-tenant architecture

Generated with Finally AI CodeKeeper - ensuring %s domain best practices.
`, g.spec.Name, g.spec.Domain, g.spec.Description, g.spec.CoreEntity, g.spec.CoreEntity, 
   g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity,
   g.generateRoleRequirements(), g.spec.Domain, g.getDomainSpecificRequirements(),
   g.spec.Backend, strings.Join(g.spec.Databases, ", "), g.spec.APIStyle, g.spec.CoreEntity, g.spec.Domain)
}

func (g *Generator) generateArchitectureDoc() string {
	return fmt.Sprintf(`# System Architecture

## Overview

%s is designed as a modern, scalable web application following the 12-Factor App methodology and domain-driven design principles for %s applications.

## Architectural Goals & Constraints

### Goals
- **Scalability**: Handle growth from startup to enterprise scale
- **Security**: %s domain security requirements and compliance
- **Maintainability**: Clear separation of concerns and clean code
- **Performance**: Fast response times and efficient resource usage
- **Reliability**: High availability and fault tolerance

### Constraints
- **Budget**: Cost-effective for early-stage deployment
- **Team Size**: Optimized for small to medium development teams
- **Time to Market**: Rapid development and deployment capability
- **Technology Stack**: %s backend, React frontend

## Technology Choices

### Backend: %s
**Rationale**: %s

### Frontend: React with TypeScript
**Rationale**: Mature ecosystem, excellent TypeScript support, component reusability, strong community

### Database: %s
**Rationale**: %s

### API Style: %s
**Rationale**: %s

## High-Level Architecture Diagrams

### C4 Model - Level 1: System Context

` + "```mermaid" + `
graph TB
    Users[Users<br/>%s] --> App[%s<br/>Web Application]
    App --> ExtAPI[External APIs<br/>Third-party Services]
    App --> DB[(Database<br/>%s)]
    App --> Auth[Authentication<br/>Service]
` + "```" + `

### C4 Model - Level 2: Container Diagram

` + "```mermaid" + `
graph TB
    subgraph "User Devices"
        Browser[Web Browser]
        Mobile[Mobile Browser]
    end
    
    subgraph "%s Application"
        Frontend[Frontend<br/>React SPA<br/>Port 3000]
        Backend[Backend API<br/>%s<br/>Port 8080]
        DB[(PostgreSQL<br/>Database<br/>Port 5432)]
    end
    
    subgraph "External Services"
        CDN[CDN<br/>Static Assets]
        Monitor[Monitoring<br/>Observability]
    end
    
    Browser --> Frontend
    Mobile --> Frontend
    Frontend --> Backend
    Backend --> DB
    Frontend --> CDN
    Backend --> Monitor
` + "```" + `

## Monorepo Strategy

Our monorepo structure promotes code sharing and consistent development practices:

` + "```" + `
%s/
├── apps/
│   ├── backend/          # %s API server
│   └── frontend/         # React web application
├── packages/
│   ├── shared-types/     # Common TypeScript types
│   └── ui-components/    # Shared UI components
├── infra/
│   ├── render/          # Render deployment config
│   └── aws/             # Terraform AWS infrastructure
├── docs/                # Project documentation
└── .github/workflows/   # CI/CD pipelines
` + "```" + `

## API Design Philosophy

### %s API Design
%s

### Authentication Strategy
- JWT tokens for stateless authentication
- Role-based access control (RBAC)
- Refresh token rotation for security

### Error Handling
- Consistent error response format
- Proper HTTP status codes
- Detailed error messages for development
- Sanitized errors for production

## Data Management Strategy

### Database Design
- PostgreSQL for transactional data
- Normalized schema with proper indexing
- Audit trails for compliance
- Soft deletes for data retention

### Caching Strategy
- Application-level caching for frequently accessed data
- CDN for static assets
- Database query optimization

## 12-Factor App Adherence

### I. Codebase
✅ Single codebase tracked in Git with multiple deployment environments

### II. Dependencies
✅ All dependencies explicitly declared in package.json/requirements files

### III. Config
✅ Configuration stored in environment variables (.env files, cloud settings)

### IV. Backing Services
✅ Database and external services treated as attached resources

### V. Build, Release, Run
✅ Strict separation: Docker builds, tagged releases, container deployment

### VI. Processes
✅ Stateless application processes, shared-nothing architecture

### VII. Port Binding
✅ Applications export services via port binding ($PORT environment variable)

### VIII. Concurrency
✅ Scale out via process model (horizontal scaling, container orchestration)

### IX. Disposability
✅ Fast startup and graceful shutdown handling

### X. Dev/Prod Parity
✅ Docker ensures development and production environment parity

### XI. Logs
✅ Structured logging to stdout/stderr, centralized log aggregation

### XII. Admin Processes
✅ Administrative tasks run as one-off processes in same environment

## Security Architecture

### Authentication & Authorization
- JWT-based stateless authentication
- Role-based access control (%s)
- Secure password hashing (bcrypt)

### Data Protection
- HTTPS everywhere
- Input validation and sanitization
- SQL injection prevention via ORM
- XSS protection

### Infrastructure Security
- VPC with private subnets (AWS)
- Security groups with minimal access
- Secrets management via cloud providers
- Regular security scanning

## Deployment Architecture

### Development
- Local Docker Compose environment
- Hot reloading for rapid development
- Seed data for testing

### Staging (Render)
- Simplified deployment for quick iterations
- Managed PostgreSQL database
- Environment variable configuration

### Production (AWS)
- ECS Fargate for container orchestration
- RDS PostgreSQL with Multi-AZ
- Application Load Balancer
- Auto-scaling based on demand
- CloudWatch monitoring

Generated with Finally AI CodeKeeper v1.0.0 for %s domain applications.
`, g.spec.Name, g.spec.Domain, g.spec.Domain, g.spec.Backend, 
   g.spec.Backend, g.getBackendRationale(), 
   strings.Join(g.spec.Databases, ", "), g.getDatabaseRationale(),
   g.spec.APIStyle, g.getAPIRationale(),
   strings.Join(g.spec.GetUserRolesList(), ", "), g.spec.Name, strings.Join(g.spec.Databases, ", "),
   g.spec.Name, g.spec.Backend, g.spec.Name, g.spec.Backend,
   g.spec.APIStyle, g.getAPIDesignDetails(),
   strings.Join(g.spec.GetUserRolesList(), ", "), g.spec.Domain)
}

// Helper functions for architecture documentation
func (g *Generator) getBackendRationale() string {
	switch {
	case strings.Contains(strings.ToLower(g.spec.Backend), "node"):
		return "Excellent for I/O intensive applications, mature ecosystem, JavaScript/TypeScript consistency across stack"
	case strings.Contains(strings.ToLower(g.spec.Backend), "python"):
		return "Rapid development, excellent libraries for data processing, strong AI/ML ecosystem"
	case strings.Contains(strings.ToLower(g.spec.Backend), "go"):
		return "High performance, excellent concurrency, fast compilation, minimal resource usage"
	default:
		return "Selected based on team expertise and application requirements"
	}
}

func (g *Generator) getDatabaseRationale() string {
	switch {
	case strings.Contains(strings.Join(g.spec.Databases, " "), "PostgreSQL"):
		return "ACID compliance, excellent performance, rich feature set, strong ecosystem"
	default:
		return "Selected based on data requirements and consistency needs"
	}
}

func (g *Generator) getAPIRationale() string {
	switch g.spec.APIStyle {
	case "RESTful APIs":
		return "Wide compatibility, simple to understand and implement, excellent caching support"
	case "GraphQL":
		return "Flexible data fetching, strong typing, single endpoint for complex queries"
	case "gRPC":
		return "High performance, strong typing, excellent for microservices communication"
	default:
		return "Selected based on client requirements and use case complexity"
	}
}

func (g *Generator) getAPIDesignDetails() string {
	switch g.spec.APIStyle {
	case "RESTful APIs":
		return `Following REST principles with proper HTTP methods:
- GET for data retrieval
- POST for resource creation
- PUT for resource updates
- DELETE for resource removal
- Proper status codes (200, 201, 400, 404, 500, etc.)
- Resource-based URLs (/api/v1/` + strings.ToLower(g.spec.CoreEntity) + `s)
- JSON request/response format`
	case "GraphQL":
		return `Single endpoint GraphQL API:
- Schema-first development
- Type-safe queries and mutations
- Flexible data fetching
- Real-time subscriptions
- Introspection capabilities`
	default:
		return "API design following industry best practices for the chosen style"
	}
}

func (g *Generator) generateRoleRequirements() string {
	roles := g.spec.GetUserRolesList()
	if len(roles) == 0 {
		return "- **F3.1**: Single user type with basic permissions"
	}
	
	result := ""
	for i, role := range roles {
		result += fmt.Sprintf("- **F3.%d**: %s role with appropriate permissions\n", i+1, strings.Title(role))
	}
	return result
}

func (g *Generator) getDomainSpecificRequirements() string {
	switch g.spec.Domain {
	case "fintech":
		return `- Decimal precision for monetary calculations
- Audit logging for all financial transactions
- Regulatory compliance features
- Multi-currency support
- Transaction idempotency
- Anti-fraud measures`
	case "healthcare":
		return `- HIPAA compliance for patient data
- Audit trails for data access
- Data encryption at rest and in transit
- Access controls and permissions
- Data retention policies
- Patient consent management`
	case "ecommerce":
		return `- Product catalog management
- Inventory tracking
- Shopping cart functionality
- Payment processing integration
- Order management
- Customer reviews and ratings`
	default:
		return `- Standard business logic validation
- Data integrity constraints
- User activity logging
- Performance optimization
- Security best practices`
	}
}

// Missing documentation generation functions
func (g *Generator) generateSystemDesignDoc() string {
	return fmt.Sprintf(`# System Design Document

## Overview
%s system design with %s domain expertise.

## Core Entity: %s
Complete CRUD operations with domain-specific validation.

## Technology Stack
- Backend: %s
- Frontend: React TypeScript
- Database: %s
- API: %s
`, g.spec.Name, g.spec.Domain, g.spec.CoreEntity, g.spec.Backend, strings.Join(g.spec.Databases, ", "), g.spec.APIStyle)
}

func (g *Generator) generateDataModelDoc() string {
	return fmt.Sprintf(`# Data Model

## Core Entity: %s
Primary entity with full CRUD operations.

## Database: %s
Optimized for %s domain requirements.
`, g.spec.CoreEntity, strings.Join(g.spec.Databases, ", "), g.spec.Domain)
}

func (g *Generator) generateAPIEndpointsDoc() string {
	return fmt.Sprintf(`# API Endpoints

## %s API
%s endpoints for %s management.

### Core Operations
- GET /api/%s - List all
- POST /api/%s - Create new
- GET /api/%s/:id - Get by ID
- PUT /api/%s/:id - Update
- DELETE /api/%s/:id - Remove
`, g.spec.APIStyle, g.spec.APIStyle, g.spec.CoreEntity, strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity))
}

func (g *Generator) generateDeploymentDoc() string {
	return fmt.Sprintf(`# Deployment Guide

## Environments
- Development: Docker Compose
- Staging: Render
- Production: AWS

## %s Application
Scalable deployment for %s domain.
`, g.spec.Name, g.spec.Domain)
}

func (g *Generator) generateOperationsDoc() string {
	return fmt.Sprintf(`# Operations & Observability

## Monitoring
%s application monitoring and observability.

## Domain: %s
Specialized monitoring for %s applications.
`, g.spec.Name, g.spec.Domain, g.spec.Domain)
}

func (g *Generator) generateTestingStrategyDoc() string {
	return fmt.Sprintf(`# Testing Strategy

## Test Coverage
Comprehensive testing for %s entity.

## Domain Testing
%s domain-specific test scenarios.
`, g.spec.CoreEntity, g.spec.Domain)
}

func (g *Generator) generateLocalDevelopmentDoc() string {
	return fmt.Sprintf(`# Local Development

## Setup
1. Clone repository
2. Run: docker compose up
3. Access: http://localhost:3000

## %s Development
%s with %s domain tools.
`, g.spec.Name, g.spec.Backend, g.spec.Domain)
}

func (g *Generator) generateContributingDoc() string {
	return fmt.Sprintf(`# Contributing Guide

## Development Process
1. Fork repository
2. Create feature branch
3. Make changes
4. Run tests
5. Submit pull request

## %s Domain
Follow %s best practices.
`, g.spec.Domain, g.spec.Domain)
}

func (g *Generator) generatePrettierConfig() string {
	return `{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 80,
  "tabWidth": 2
}`
}

func (g *Generator) generateESLintConfig() string {
	return `module.exports = {
  extends: ["eslint:recommended", "@typescript-eslint/recommended"],
  parser: "@typescript-eslint/parser",
  plugins: ["@typescript-eslint"],
  rules: {
    "no-console": "warn",
    "@typescript-eslint/no-unused-vars": "error"
  }
};`
}

func (g *Generator) generateVSCodeSettings() string {
	return `{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "typescript.preferences.importModuleSpecifier": "relative"
}`
}

func (g *Generator) generateVSCodeExtensions() string {
	return `{
  "recommendations": [
    "esbenp.prettier-vscode",
    "dbaeumer.vscode-eslint",
    "ms-vscode.vscode-typescript-next",
    "bradlc.vscode-tailwindcss"
  ]
}`
}

func (g *Generator) generateVSCodeLaunch() string {
	return `{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Backend",
      "type": "node",
      "request": "launch",
      "program": "${workspaceFolder}/apps/backend/src/server.js",
      "env": { "NODE_ENV": "development" }
    }
  ]
}`
}

// generateBackendFiles creates backend foundation and README with AI expectations
func (g *Generator) generateBackendFiles() map[string]string {
	files := make(map[string]string)
	
	switch g.spec.GetBackendLanguage() {
	case "javascript":
		files["apps/backend/package.json"] = g.generateBackendPackageJson()
		files["apps/backend/.env.example"] = g.generateBackendEnvExample()
		files["apps/backend/README.md"] = g.generateBackendReadme()
		files["apps/backend/Dockerfile"] = g.generateBackendDockerfile()
		// Basic starter file
		files["apps/backend/src/server.js"] = g.generateBasicServer()
		
	case "python":
		files["apps/backend/requirements.txt"] = g.generatePythonRequirements()
		files["apps/backend/pyproject.toml"] = g.generatePyprojectToml()
		files["apps/backend/.env.example"] = g.generateBackendEnvExample()
		files["apps/backend/README.md"] = g.generateBackendReadme()
		files["apps/backend/Dockerfile"] = g.generateBackendDockerfile()
		// Basic starter file
		files["apps/backend/main.py"] = g.generateBasicPythonMain()
		
	case "go":
		files["apps/backend/go.mod"] = g.generateGoMod()
		files["apps/backend/.env.example"] = g.generateBackendEnvExample()
		files["apps/backend/README.md"] = g.generateBackendReadme()
		files["apps/backend/Dockerfile"] = g.generateBackendDockerfile()
		// Basic starter file
		files["apps/backend/main.go"] = g.generateGoMain()
	}
	
	return files
}

// generateFrontendFiles creates frontend foundation with AI guidelines
func (g *Generator) generateFrontendFiles() map[string]string {
	files := make(map[string]string)
	
	// Frontend foundation
	files["apps/frontend/package.json"] = g.generateFrontendPackageJson()
	files["apps/frontend/.env.example"] = g.generateFrontendEnvExample()
	files["apps/frontend/README.md"] = g.generateFrontendReadme()
	files["apps/frontend/Dockerfile"] = g.generateFrontendDockerfile()
	files["apps/frontend/index.html"] = g.generateFrontendIndexHtml()
	files["apps/frontend/vite.config.ts"] = g.generateViteConfig()
	
	// Basic React structure
	files["apps/frontend/src/main.tsx"] = g.generateFrontendMain()
	files["apps/frontend/src/App.tsx"] = g.generateFrontendApp()
	files["apps/frontend/src/App.css"] = g.generateFrontendCSS()
	
	// Shared types package
	files["packages/shared-types/package.json"] = g.generateSharedTypesPackageJson()
	files["packages/shared-types/README.md"] = g.generateSharedTypesReadme()
	files["packages/shared-types/src/index.ts"] = g.generateSharedTypesIndex()
	
	return files
}

// generateDevelopmentTooling creates development tooling files based on tech stack
func (g *Generator) generateDevelopmentTooling() map[string]string {
	files := make(map[string]string)
	
	// Always include frontend tooling for React/TypeScript
	files[".prettierrc.json"] = g.generatePrettierConfig()
	files[".eslintrc.js"] = g.generateESLintConfig()
	
	// Add backend-specific tooling
	switch g.spec.GetBackendLanguage() {
	case "python":
		files[".flake8"] = g.generateFlake8Config()
		files["pyproject.toml"] = g.generatePyprojectToml()
	case "go":
		files[".golangci.yml"] = g.generateGolangCIConfig()
	}
	
	return files
}

// Backend-specific configuration generators
func (g *Generator) generateBackendEnvExample() string {
	return fmt.Sprintf(`# Database
DATABASE_URL=postgresql://dev:dev@localhost:5432/%s

# Server
PORT=8080
NODE_ENV=development

# JWT
JWT_SECRET=your-jwt-secret-here

# Domain-specific environment variables
%s
`, g.spec.Name, g.getDomainEnvVars())
}

func (g *Generator) generatePythonRequirements() string {
	framework := g.spec.GetBackendFramework()
	requirements := ""
	
	switch framework {
	case "django":
		requirements = `Django>=4.2.0
djangorestframework>=3.14.0
django-cors-headers>=4.0.0
psycopg2-binary>=2.9.0
python-decouple>=3.8.0
djangorestframework-simplejwt>=5.2.0`
	case "flask":
		requirements = `Flask>=2.3.0
Flask-SQLAlchemy>=3.0.0
Flask-Migrate>=4.0.0
Flask-CORS>=4.0.0
psycopg2-binary>=2.9.0
PyJWT>=2.8.0
python-decouple>=3.8.0`
	default:
		requirements = `fastapi>=0.100.0
uvicorn>=0.22.0
sqlalchemy>=2.0.0
psycopg2-binary>=2.9.0
python-decouple>=3.8.0
PyJWT>=2.8.0`
	}
	
	return requirements
}

func (g *Generator) generatePyprojectToml() string {
	return fmt.Sprintf(`[build-system]
requires = ["setuptools>=61.0", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "%s-backend"
version = "1.0.0"
description = "%s backend API"
dependencies = []

[tool.black]
line-length = 88
target-version = ['py311']

[tool.isort]
profile = "black"
multi_line_output = 3

[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]
`, g.spec.Name, g.spec.Description)
}

func (g *Generator) generateGoMod() string {
	return fmt.Sprintf(`module %s-backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/lib/pq v1.10.9
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/joho/godotenv v1.4.0
)
`, g.spec.Name)
}

func (g *Generator) generateGoMain() string {
	return fmt.Sprintf(`package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "%s-backend",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/%s", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from %s API"})
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
`, g.spec.Name, strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity)
}

func (g *Generator) generateFlake8Config() string {
	return `[flake8]
max-line-length = 88
extend-ignore = E203, W503
exclude = .git,__pycache__,migrations,venv,.venv`
}

func (g *Generator) generateGolangCIConfig() string {
	return `run:
  timeout: 5m
  
linters:
  enable:
    - gofmt
    - golint
    - govet
    - ineffassign
    - misspell
    - revive
    
linters-settings:
  golint:
    min-confidence: 0.8`
}

func (g *Generator) getDomainEnvVars() string {
	switch g.spec.Domain {
	case "fintech":
		return `# Fintech-specific
ENCRYPTION_KEY=your-encryption-key
AUDIT_LOG_LEVEL=INFO
COMPLIANCE_MODE=strict`
	case "healthcare":
		return `# Healthcare-specific
HIPAA_ENCRYPTION_KEY=your-hipaa-key
PHI_LOG_LEVEL=WARN
COMPLIANCE_MODE=hipaa`
	default:
		return `# Application-specific
LOG_LEVEL=INFO
DEBUG_MODE=true`
	}
}

// Backend README generation
func (g *Generator) generateBackendReadme() string {
	domainGuidelines := g.getDomainSpecificGuidelines()
	
	return fmt.Sprintf(`# %s Backend

%s backend API built with %s and %s domain expertise.

## Technology Stack

- **Language**: %s
- **Framework**: %s
- **Database**: %s
- **API Style**: %s
- **Authentication**: JWT with RBAC
- **Testing**: %s

## Quick Start

1. Install dependencies: ` + "`npm install`" + ` (or equivalent)
2. Copy .env.example to .env and configure
3. Run database setup: ` + "`npm run db:setup`" + `
4. Start development: ` + "`npm run dev`" + `

## AI Development Guidelines

This backend follows the **AI Master Prompt** specifications. When implementing features:

### Core Architecture Requirements

**12-Factor App Compliance:**
- ✅ All config via environment variables (.env)
- ✅ Stateless processes with external state storage
- ✅ Structured logging to stdout/stderr
- ✅ Graceful shutdown handling (SIGTERM/SIGINT)
- ✅ Port binding via $PORT environment variable

**Security Requirements:**
- JWT authentication with secure secret rotation
- Input validation for all endpoints
- Rate limiting and DDoS protection
- CORS configuration for frontend integration
- Helmet.js security headers (Node.js)

**Database Requirements:**
- ORM/ODM usage (no raw SQL)
- Connection pooling and retry logic
- Migration-based schema management
- Audit logging for critical operations

### Expected Project Structure

` + "`" + `
src/
├── config/         # Environment and database configuration
├── controllers/    # HTTP request handlers (thin layer)
├── services/       # Business logic (main implementation)
├── models/         # Database models with validation
├── routes/         # API route definitions with middleware
├── middleware/     # Auth, logging, error handling
├── utils/          # Shared utilities and helpers
└── jobs/           # Background tasks (if needed)

tests/
├── unit/           # Service and utility tests
├── integration/    # API endpoint tests
└── fixtures/       # Test data and mocks
` + "`" + `

### Core Entity: %s

**Required CRUD Operations:**
- **GET** /api/v1/%s - List with pagination, filtering, search
- **GET** /api/v1/%s/:id - Get single item with validation
- **POST** /api/v1/%s - Create with validation and domain rules
- **PUT** /api/v1/%s/:id - Update with partial support
- **DELETE** /api/v1/%s/:id - Soft delete preferred

**Authentication Endpoints:**
- **POST** /api/v1/auth/register - User registration with email verification
- **POST** /api/v1/auth/login - Login with JWT token generation
- **GET** /api/v1/auth/profile - Protected user profile
- **PUT** /api/v1/auth/profile - Profile updates

### Domain-Specific Requirements: %s

%s

### Development Commands

` + "`" + `bash
# Development
npm run dev          # Start with hot reload
npm run lint         # Code linting
npm run format       # Code formatting

# Database
npm run db:setup     # Initialize database
npm run db:migrate   # Run migrations
npm run db:seed      # Seed development data

# Testing
npm test             # Run all tests
npm run test:unit    # Unit tests only
npm run test:int     # Integration tests
npm run test:cov     # Coverage report

# Production
npm run build        # Build for production
npm start            # Start production server
` + "`" + `

### Implementation Priority

When implementing this backend, follow this order:

1. **Foundation** (Complete first)
   - Database connection and models
   - Authentication system
   - Error handling middleware
   - Basic logging

2. **Core Features**
   - %s CRUD operations
   - Input validation
   - Business logic services
   - Unit tests

3. **Production Features**
   - Integration tests
   - Performance optimization
   - Security hardening
   - Monitoring endpoints

### Code Standards

- Use TypeScript for type safety (if Node.js)
- Implement comprehensive input validation
- Follow RESTful API design principles
- Write tests for all business logic
- Use structured logging with correlation IDs
- Handle errors gracefully with proper HTTP status codes

Generated with AI CodeKeeper v1.0.0 - Ready for AI implementation
`, g.spec.Name, g.spec.Description, g.spec.Backend, g.spec.Domain,
   g.spec.GetBackendLanguage(), g.spec.GetBackendFramework(), 
   strings.Join(g.spec.Databases, ", "), g.spec.APIStyle, g.getTestingFramework(),
   g.spec.CoreEntity, strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity),
   strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity),
   g.spec.Domain, domainGuidelines, g.spec.CoreEntity)
}

// Backend Dockerfile generation
func (g *Generator) generateBackendDockerfile() string {
	switch g.spec.GetBackendLanguage() {
	case "javascript":
		return `# Multi-stage build for Node.js
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine AS runtime

# Create non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S backend -u 1001

WORKDIR /app

# Copy built app
COPY --from=builder /app/node_modules ./node_modules
COPY --chown=backend:nodejs . .

USER backend

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["node", "src/server.js"]`

	case "python":
		return `# Multi-stage build for Python
FROM python:3.11-slim AS builder

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir --user -r requirements.txt

FROM python:3.11-slim AS runtime

# Create non-root user
RUN groupadd -r backend && useradd -r -g backend backend

WORKDIR /app

# Copy installed packages
COPY --from=builder /root/.local /home/backend/.local
COPY --chown=backend:backend . .

USER backend

# Add local packages to PATH
ENV PATH=/home/backend/.local/bin:$PATH

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["python", "main.py"]`

	case "go":
		return `# Multi-stage build for Go
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest AS runtime

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S backend
RUN adduser -S backend -u 1001 -G backend

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

USER backend

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]`

	default:
		return "# Dockerfile placeholder"
	}
}

// Helper functions
func (g *Generator) getTestingFramework() string {
	switch g.spec.GetBackendLanguage() {
	case "javascript":
		return "Jest"
	case "python":
		return "pytest"
	case "go":
		return "Go testing"
	default:
		return "Framework tests"
	}
}

// Node.js Backend Implementation
func (g *Generator) generateNodeServer() string {
	return fmt.Sprintf(`const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
const { createLogger } = require('./utils/logger');
const { connectDatabase } = require('./config/database');
const errorHandler = require('./middleware/errorHandler');

// Import routes
const authRoutes = require('./routes/auth');
const %sRoutes = require('./routes/%s');

const logger = createLogger();
const app = express();

// Security middleware
app.use(helmet());
app.use(cors({
  origin: process.env.FRONTEND_URL || 'http://localhost:3000',
  credentials: true
}));

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
});
app.use(limiter);

// Body parsing
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Request logging
app.use((req, res, next) => {
  logger.info({
    method: req.method,
    url: req.url,
    ip: req.ip,
    userAgent: req.get('User-Agent')
  }, 'HTTP Request');
  next();
});

// Health check
app.get('/health', (req, res) => {
  res.json({ 
    status: 'healthy', 
    service: '%s-backend',
    timestamp: new Date().toISOString(),
    uptime: process.uptime()
  });
});

// API routes
app.use('/api/v1/auth', authRoutes);
app.use('/api/v1/%s', %sRoutes);

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({ 
    error: 'Route not found',
    path: req.originalUrl 
  });
});

// Error handling
app.use(errorHandler);

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully');
  process.exit(0);
});

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully');
  process.exit(0);
});

const PORT = process.env.PORT || 8080;

async function startServer() {
  try {
    // Connect to database
    await connectDatabase();
    logger.info('Database connected successfully');

    // Start server
    app.listen(PORT, () => {
      logger.info({ port: PORT }, 'Server started successfully');
      logger.info({ 
        domain: '%s',
        entity: '%s',
        nodeEnv: process.env.NODE_ENV || 'development'
      }, 'Application configuration');
    });
  } catch (error) {
    logger.error({ error: error.message }, 'Failed to start server');
    process.exit(1);
  }
}

startServer();

module.exports = app;
`, strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity),
   g.spec.Name, strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity),
   g.spec.Domain, g.spec.CoreEntity)
}

func (g *Generator) generateNodeDatabaseConfig() string {
	return fmt.Sprintf(`const { Sequelize } = require('sequelize');
const { createLogger } = require('../utils/logger');

const logger = createLogger();

const sequelize = new Sequelize(process.env.DATABASE_URL, {
  dialect: 'postgres',
  logging: (msg) => logger.debug(msg),
  pool: {
    max: 10,
    min: 0,
    acquire: 30000,
    idle: 10000
  },
  retry: {
    match: [
      /ETIMEDOUT/,
      /EHOSTUNREACH/,
      /ECONNRESET/,
      /ECONNREFUSED/,
      /ETIMEDOUT/,
      /ESOCKETTIMEDOUT/,
      /EHOSTUNREACH/,
      /EPIPE/,
      /EAI_AGAIN/,
      /SequelizeConnectionError/,
      /SequelizeConnectionRefusedError/,
      /SequelizeHostNotFoundError/,
      /SequelizeHostNotReachableError/,
      /SequelizeInvalidConnectionError/,
      /SequelizeConnectionTimedOutError/
    ],
    max: 3
  }
});

async function connectDatabase() {
  try {
    await sequelize.authenticate();
    logger.info('Database connection established successfully');
    
    // Sync models in development
    if (process.env.NODE_ENV === 'development') {
      await sequelize.sync({ alter: true });
      logger.info('Database models synchronized');
    }
  } catch (error) {
    logger.error({ error: error.message }, 'Database connection failed');
    throw error;
  }
}

module.exports = { sequelize, connectDatabase };`)
}

func (g *Generator) generateNodeAuthConfig() string {
	return `const jwt = require('jsonwebtoken');
const bcrypt = require('bcrypt');

const JWT_SECRET = process.env.JWT_SECRET || 'fallback-secret-key';
const JWT_EXPIRES_IN = process.env.JWT_EXPIRES_IN || '24h';
const BCRYPT_ROUNDS = parseInt(process.env.BCRYPT_ROUNDS) || 12;

function generateToken(payload) {
  return jwt.sign(payload, JWT_SECRET, { expiresIn: JWT_EXPIRES_IN });
}

function verifyToken(token) {
  return jwt.verify(token, JWT_SECRET);
}

async function hashPassword(password) {
  return bcrypt.hash(password, BCRYPT_ROUNDS);
}

async function comparePassword(password, hashedPassword) {
  return bcrypt.compare(password, hashedPassword);
}

module.exports = {
  generateToken,
  verifyToken,
  hashPassword,
  comparePassword,
  JWT_SECRET,
  JWT_EXPIRES_IN
};`
}

// Helper function for domain-specific guidelines
func (g *Generator) getDomainSpecificGuidelines() string {
	switch g.spec.Domain {
	case "fintech":
		return `**Financial Domain Requirements:**
- Use decimal arithmetic for all monetary calculations (never floats)
- Implement comprehensive audit trails for all transactions
- Add idempotency keys for transaction endpoints
- Implement proper error handling for payment failures
- Add rate limiting for transaction endpoints
- Use encryption for sensitive financial data
- Implement proper money transfer validation
- Add compliance logging for regulatory requirements

**Data Models:**
- Transaction amounts must use decimal/currency types
- All financial operations must be logged with audit trails
- Implement proper currency handling and conversion
- Add transaction status tracking (pending, completed, failed)

**Security:**
- PCI compliance considerations for card data
- Implement proper access controls for financial operations
- Add fraud detection patterns
- Use secure random generators for transaction IDs`

	case "healthcare":
		return `**Healthcare Domain Requirements:**
- HIPAA compliance for all patient data handling
- Implement comprehensive audit logging for data access
- Add proper consent management
- Use encryption for PHI (Personal Health Information)
- Implement proper data retention policies
- Add patient data anonymization capabilities

**Data Models:**
- Patient records with proper privacy controls
- Medical history with audit trails
- Appointment scheduling with constraints
- Insurance and billing information security

**Security:**
- Role-based access control for medical staff
- Patient data access logging
- Secure data transmission and storage
- Regular security audits and compliance checks`

	case "ecommerce":
		return `**E-commerce Domain Requirements:**
- Product catalog management with inventory tracking
- Shopping cart session management
- Order processing workflows
- Payment integration patterns
- Shipping and fulfillment tracking

**Data Models:**
- Product variants and pricing
- Inventory management with stock levels
- Order lifecycle management
- Customer profiles and preferences
- Review and rating systems

**Security:**
- Secure payment processing
- Customer data protection
- Order fraud prevention
- Secure checkout processes`

	default:
		return `**General Application Requirements:**
- Follow RESTful API design principles
- Implement proper input validation and sanitization
- Add comprehensive error handling
- Use structured logging throughout
- Implement caching where appropriate
- Add proper database indexing
- Follow security best practices`
	}
}

// Basic server generation (minimal starter)
func (g *Generator) generateBasicServer() string {
	return fmt.Sprintf(`const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
require('dotenv').config();

const app = express();

// Security middleware
app.use(helmet());
app.use(cors());
app.use(express.json());

// Health check
app.get('/health', (req, res) => {
  res.json({ 
    status: 'healthy', 
    service: '%s-backend',
    timestamp: new Date().toISOString()
  });
});

// API routes placeholder
app.get('/api/v1/%s', (req, res) => {
  res.json({ message: 'API endpoint ready for implementation' });
});

// Error handling
app.use((error, req, res, next) => {
  console.error(error);
  res.status(500).json({ error: 'Internal server error' });
});

const PORT = process.env.PORT || 8080;

app.listen(PORT, () => {
  console.log('Server running on port', PORT);
  console.log('Domain:', '%s');
  console.log('Core Entity:', '%s');
});

module.exports = app;`, g.spec.Name, strings.ToLower(g.spec.CoreEntity), g.spec.Domain, g.spec.CoreEntity)
}

// Basic Python main
func (g *Generator) generateBasicPythonMain() string {
	return fmt.Sprintf(`from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import os

app = FastAPI(
    title="%s API",
    description="%s",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": "%s-backend",
        "domain": "%s",
        "core_entity": "%s"
    }

@app.get("/api/v1/%s")
async def get_%s():
    return {"message": "API endpoint ready for implementation"}

if __name__ == "__main__":
    port = int(os.getenv("PORT", 8080))
    uvicorn.run(app, host="0.0.0.0", port=port)
`, g.spec.Name, g.spec.Description, g.spec.Name, g.spec.Domain, g.spec.CoreEntity, 
   strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity))
}

// Frontend generation functions
func (g *Generator) generateFrontendEnvExample() string {
	return `# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api/v1

# Application Configuration
VITE_APP_NAME=` + g.spec.Name + `
VITE_APP_VERSION=1.0.0

# Environment
VITE_NODE_ENV=development

# Optional: Analytics, monitoring, etc.
# VITE_ANALYTICS_ID=your-analytics-id
# VITE_SENTRY_DSN=your-sentry-dsn
`
}

func (g *Generator) generateFrontendReadme() string {
	return fmt.Sprintf(`# %s Frontend

Modern React TypeScript application with %s domain expertise.

## Technology Stack

- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite for fast development and optimized builds
- **Routing**: React Router v6
- **State Management**: Zustand (lightweight, TypeScript-first)
- **HTTP Client**: Axios with interceptors
- **Styling**: CSS Modules + Tailwind CSS (when needed)
- **Icons**: Lucide React
- **Testing**: Vitest + React Testing Library
- **Type Safety**: Full TypeScript integration

## Quick Start

1. Install dependencies: ` + "`npm install`" + `
2. Copy .env.example to .env and configure API URL
3. Start development server: ` + "`npm run dev`" + `
4. Open http://localhost:5173

## AI Development Guidelines

This frontend follows the **AI Master Prompt** specifications for modern React development.

### Core Architecture Requirements

**Component Architecture:**
- Functional components with TypeScript
- Custom hooks for reusable logic
- Error boundaries for robust UI
- Responsive design with mobile-first approach

**State Management:**
- Zustand for global state (auth, %s data)
- React hooks for local component state
- Optimistic updates with error rollback
- Persistent auth state across sessions

**API Integration:**
- Axios with request/response interceptors
- Automatic JWT token handling
- Error handling with user-friendly messages
- Loading states and optimistic updates

### Expected Project Structure

` + "`" + `
src/
├── components/         # Reusable UI components
│   ├── ui/            # Basic UI elements (Button, Input, Modal)
│   ├── forms/         # Form components with validation
│   └── layout/        # Layout components (Header, Sidebar)
├── pages/             # Page-level components (routes)
│   ├── auth/          # Login, Register, Profile pages
│   ├── %s/           # %s CRUD pages
│   └── dashboard/     # Dashboard and analytics
├── hooks/             # Custom React hooks
│   ├── useAuth.ts     # Authentication logic
│   ├── useApi.ts      # API call abstractions
│   └── use%s.ts      # %s-specific operations
├── services/          # API clients and business logic
│   ├── api.ts         # Axios configuration and interceptors
│   ├── auth.ts        # Authentication API calls
│   └── %s.ts         # %s API operations
├── store/             # Zustand stores
│   ├── authStore.ts   # User authentication state
│   └── %sStore.ts    # %s data management
├── types/             # TypeScript type definitions
├── utils/             # Helper functions and utilities
├── assets/            # Static assets (images, fonts)
└── styles/            # Global styles and themes
` + "`" + `

### Core Entity: %s

**Required UI Components:**
- **%sList** - Paginated list with search and filters
- **%sCard** - Individual %s preview component
- **%sForm** - Create/edit form with validation
- **%sDetail** - Full detail view with actions
- **%sFilters** - Advanced filtering component

**Required Pages:**
- **GET** / - Dashboard with %s overview
- **GET** /%s - List all %s with pagination
- **GET** /%s/new - Create new %s form
- **GET** /%s/:id - View %s details
- **GET** /%s/:id/edit - Edit %s form

### Domain-Specific Requirements: %s

%s

### Development Commands

` + "`" + `bash
# Development
npm run dev              # Start dev server with hot reload
npm run build            # Build for production
npm run preview          # Preview production build

# Code Quality
npm run lint             # ESLint checking
npm run lint:fix         # Auto-fix ESLint issues
npm run type-check       # TypeScript type checking

# Testing
npm test                 # Run all tests
npm run test:ui          # Interactive test UI
npm run test:coverage    # Coverage report

# Guard Rails
npm run check            # AI CodeKeeper validation
` + "`" + `

### Implementation Priority

When implementing this frontend, follow this order:

1. **Foundation** (Complete first)
   - API client setup with authentication
   - Basic routing and layout
   - Auth pages (login/register)
   - Error boundaries and loading states

2. **Core Features**
   - %s CRUD components
   - Form validation and submission
   - List/detail views with navigation
   - Search and filtering

3. **User Experience**
   - Responsive design
   - Loading states and optimistic updates
   - Error handling with user feedback
   - Accessibility improvements

4. **Production Features**
   - Performance optimization
   - SEO considerations
   - PWA features (if needed)
   - Analytics integration

### Code Standards

- **TypeScript**: Strict mode enabled, no ` + "`any`" + ` types
- **Components**: Functional components with proper typing
- **Hooks**: Custom hooks for complex logic
- **Testing**: Test user interactions, not implementation
- **Accessibility**: WCAG 2.1 AA compliance
- **Performance**: Code splitting and lazy loading
- **State**: Minimize state, prefer derived values

### API Integration

All API calls should:
- Use the configured Axios instance
- Handle loading/error states consistently
- Implement optimistic updates where appropriate
- Show user-friendly error messages
- Include proper TypeScript types from shared-types package

Generated with AI CodeKeeper v1.0.0 - Ready for AI implementation
`, g.spec.Name, g.spec.Domain, g.spec.CoreEntity, 
   strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   g.spec.CoreEntity, g.spec.CoreEntity, strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity, g.spec.CoreEntity,
   strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
   g.spec.Domain, g.getFrontendDomainGuidelines(), g.spec.CoreEntity)
}

func (g *Generator) generateFrontendDockerfile() string {
	return `# Multi-stage build for React frontend
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage with Nginx
FROM nginx:alpine AS runtime

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create non-root user
RUN addgroup -g 1001 -S frontend
RUN adduser -S frontend -u 1001 -G frontend

# Change ownership
RUN chown -R frontend:frontend /usr/share/nginx/html
RUN chown -R frontend:frontend /var/cache/nginx
RUN chown -R frontend:frontend /var/log/nginx
RUN chown -R frontend:frontend /etc/nginx/conf.d
RUN touch /var/run/nginx.pid
RUN chown -R frontend:frontend /var/run/nginx.pid

USER frontend

EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

CMD ["nginx", "-g", "daemon off;"]`
}

func (g *Generator) generateFrontendIndexHtml() string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="description" content="%s - %s domain application" />
    <title>%s</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>`, g.spec.Name, g.spec.Domain, g.spec.Name)
}

func (g *Generator) generateViteConfig() string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@shared': path.resolve(__dirname, '../shared-types/src'),
    },
  },
  server: {
    port: 3000,
    host: true, // Allow external connections
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
  },
})`
}

func (g *Generator) generateFrontendMain() string {
	return `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './App.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)`
}

func (g *Generator) generateFrontendApp() string {
	return fmt.Sprintf(`import React from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import './App.css'

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>%s</h1>
          <p>%s Application - %s Domain</p>
          <p>
            Ready for AI implementation with full TypeScript support
          </p>
        </header>
        
        <main>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/auth/login" element={<div>Login Page - Ready for AI implementation</div>} />
            <Route path="/auth/register" element={<div>Register Page - Ready for AI implementation</div>} />
            <Route path="/%s" element={<div>%s List - Ready for AI implementation</div>} />
            <Route path="/%s/new" element={<div>New %s - Ready for AI implementation</div>} />
            <Route path="/%s/:id" element={<div>%s Detail - Ready for AI implementation</div>} />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

function Dashboard() {
  return (
    <div>
      <h2>Dashboard</h2>
      <p>Core Entity: %s</p>
      <p>Domain: %s</p>
      <p>Ready for AI implementation</p>
    </div>
  )
}

export default App`, g.spec.Name, g.spec.Description, g.spec.Domain,
		strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
		strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
		strings.ToLower(g.spec.CoreEntity), g.spec.CoreEntity,
		g.spec.CoreEntity, g.spec.Domain)
}

func (g *Generator) generateFrontendCSS() string {
	return `#root {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}

.App {
  min-height: 100vh;
}

.App-header {
  background-color: #f8f9fa;
  padding: 2rem;
  margin-bottom: 2rem;
  border-radius: 8px;
}

.App-header h1 {
  color: #1a1a1a;
  margin: 0 0 1rem 0;
}

.App-header p {
  color: #666;
  margin: 0.5rem 0;
}

main {
  text-align: left;
}

/* Ready for component-specific styles */
.dashboard {
  display: grid;
  gap: 1rem;
}

.card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.button {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.button:hover {
  background-color: #2563eb;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
}

/* Responsive design */
@media (max-width: 768px) {
  #root {
    padding: 1rem;
  }
  
  .App-header {
    padding: 1rem;
  }
}`
}

func (g *Generator) generateSharedTypesReadme() string {
	return fmt.Sprintf(`# Shared Types

TypeScript type definitions shared between frontend and backend for %s.

## Purpose

This package contains:
- API request/response types
- Database model interfaces
- Common utility types
- Domain-specific type definitions

## Usage

` + "`" + `typescript
import { User, %s, APIResponse } from '@%s/shared-types'
` + "`" + `

## Type Categories

### Authentication Types
- User interface with role-based permissions
- Login/Register request/response types
- JWT token payload structure

### %s Types
- %s entity interface
- CRUD operation request/response types
- Domain-specific validation types

### API Types
- Standard API response wrapper
- Error response structure
- Pagination metadata
- Filter and search parameters

## Development

When adding new types:
1. Add the interface/type definition
2. Export from src/index.ts
3. Update version in package.json
4. The monorepo will automatically pick up changes

Generated with AI CodeKeeper v1.0.0
`, g.spec.Name, g.spec.CoreEntity, g.spec.Name, g.spec.CoreEntity, g.spec.CoreEntity)
}

func (g *Generator) generateSharedTypesIndex() string {
	entity := g.spec.CoreEntity
	
	var domainSpecificTypes string
	switch g.spec.Domain {
	case "fintech":
		domainSpecificTypes = `
// Fintech-specific types
export interface Transaction {
  id: string;
  amount: number; // Use Decimal type in implementation
  currency: string;
  type: 'debit' | 'credit' | 'transfer';
  description?: string;
  reference?: string;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  userId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTransactionRequest {
  amount: number;
  currency: string;
  type: 'debit' | 'credit' | 'transfer';
  description?: string;
  reference?: string;
}

export interface TransactionFilters {
  status?: Transaction['status'];
  type?: Transaction['type'];
  dateFrom?: string;
  dateTo?: string;
  minAmount?: number;
  maxAmount?: number;
}`
	case "healthcare":
		domainSpecificTypes = `
// Healthcare-specific types
export interface MedicalRecord {
  id: string;
  patientId: string;
  type: string;
  date: string;
  diagnosis?: string;
  treatment?: string;
  status: 'scheduled' | 'in-progress' | 'completed' | 'cancelled';
  userId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateMedicalRecordRequest {
  patientId: string;
  type: string;
  date: string;
  diagnosis?: string;
  treatment?: string;
}

export interface MedicalRecordFilters {
  status?: MedicalRecord['status'];
  type?: string;
  dateFrom?: string;
  dateTo?: string;
  patientId?: string;
}`
	default:
		domainSpecificTypes = fmt.Sprintf(`
// %s entity types
export interface %s {
  id: string;
  name: string;
  description?: string;
  status: 'active' | 'inactive' | 'pending';
  userId: string;
  createdAt: string;
  updatedAt: string;
}

export interface Create%sRequest {
  name: string;
  description?: string;
  status?: %s['status'];
}

export interface %sFilters {
  status?: %s['status'];
  search?: string;
}`, entity, entity, entity, entity, entity, entity)
	}

	return fmt.Sprintf(`// Shared TypeScript types for %s

// User and Authentication types
export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: 'admin' | 'user' | 'viewer';
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role?: User['role'];
}

export interface AuthResponse {
  user: User;
  token: string;
}
%s

// API Response types
export interface APIResponse<T = any> {
  data?: T;
  message?: string;
  error?: string;
  details?: string[];
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationMeta;
}

// Common utility types
export interface BaseEntity {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export interface Timestamps {
  createdAt: string;
  updatedAt: string;
}

// Error types
export interface ValidationError {
  field: string;
  message: string;
}

export interface APIError {
  error: string;
  details?: ValidationError[];
  statusCode?: number;
}

// Export the main entity type alias for convenience
export type %sEntity = %s;
export type %sListResponse = PaginatedResponse<%s>;
`, g.spec.Name, domainSpecificTypes, entity, entity, entity, entity)
}

func (g *Generator) getFrontendDomainGuidelines() string {
	switch g.spec.Domain {
	case "fintech":
		return `**Financial UI Requirements:**
- Display monetary amounts with proper decimal precision
- Show transaction history with audit trails
- Implement secure payment forms with validation
- Add transaction status indicators and progress tracking
- Use proper currency formatting based on locale
- Implement real-time balance updates
- Add fraud detection alerts and notifications`

	case "healthcare":
		return `**Healthcare UI Requirements:**
- Ensure patient data privacy and HIPAA compliance
- Implement secure access controls for medical records
- Add audit logging for all data access
- Use medical terminology and proper date formatting
- Implement appointment scheduling with constraints
- Add patient consent management interfaces
- Ensure accessibility for healthcare professionals`

	case "ecommerce":
		return `**E-commerce UI Requirements:**
- Product catalog with search and filtering
- Shopping cart with persistent state
- Checkout flow with payment integration
- Order tracking and status updates
- Product reviews and ratings display
- Inventory status indicators
- Customer account management`

	default:
		return `**General UI Requirements:**
- Responsive design for all screen sizes
- Accessible interface following WCAG guidelines
- Consistent design system and component library
- Proper error handling with user-friendly messages
- Loading states and progress indicators
- Search and filtering capabilities
- Data visualization where appropriate`
	}
}

// Infrastructure as Code generation
func (g *Generator) generateInfrastructureFiles() map[string]string {
	files := make(map[string]string)
	
	// Render deployment configuration
	files["infra/render/render.yaml"] = g.generateRenderConfig()
	files["infra/render/README.md"] = g.generateRenderReadme()
	
	// AWS Terraform infrastructure
	files["infra/aws/main.tf"] = g.generateTerraformMain()
	files["infra/aws/variables.tf"] = g.generateTerraformVariables()
	files["infra/aws/outputs.tf"] = g.generateTerraformOutputs()
	files["infra/aws/modules/vpc/main.tf"] = g.generateTerraformVPC()
	files["infra/aws/modules/ecs/main.tf"] = g.generateTerraformECS()
	files["infra/aws/modules/rds/main.tf"] = g.generateTerraformRDS()
	files["infra/aws/modules/secrets/main.tf"] = g.generateTerraformSecrets()
	files["infra/aws/README.md"] = g.generateTerraformReadme()
	
	// Enhanced CI/CD pipeline
	files[".github/workflows/ci.yml"] = g.generateEnhancedCIWorkflow()
	files[".github/workflows/deploy.yml"] = g.generateDeploymentWorkflow()
	
	return files
}

func (g *Generator) generateRenderConfig() string {
	return fmt.Sprintf(`# Render deployment configuration for %s
# Documentation: https://render.com/docs/yaml-spec

services:
  # Backend API Service
  - type: web
    name: %s-backend
    runtime: docker
    repo: # Set your GitHub repository URL
    branch: main
    dockerfilePath: ./apps/backend/Dockerfile
    dockerContext: ./apps/backend
    plan: starter # Change to standard/pro for production
    region: oregon # or your preferred region
    
    # Environment variables
    envVars:
      - key: NODE_ENV
        value: production
      - key: PORT
        value: 8080
      - key: DATABASE_URL
        fromDatabase:
          name: %s-db
          property: connectionString
      - key: JWT_SECRET
        generateValue: true # Render will generate a secure value
      - key: FRONTEND_URL
        value: https://%s-frontend.onrender.com
      
    # Health check
    healthCheckPath: /health
    
    # Build command (if not using Dockerfile)
    # buildCommand: npm install && npm run build
    # startCommand: npm start

  # Frontend Web Service  
  - type: web
    name: %s-frontend
    runtime: docker
    repo: # Set your GitHub repository URL
    branch: main
    dockerfilePath: ./apps/frontend/Dockerfile
    dockerContext: ./apps/frontend
    plan: starter
    region: oregon
    
    # Environment variables
    envVars:
      - key: VITE_API_BASE_URL
        value: https://%s-backend.onrender.com/api/v1
      - key: VITE_APP_NAME
        value: %s
      - key: VITE_NODE_ENV
        value: production
    
    # Custom domain (optional)
    # customDomains:
    #   - name: your-domain.com

# Database
databases:
  - name: %s-db
    databaseName: %s
    user: %s_user
    plan: starter # Change to standard/pro for production
    region: oregon
    version: "15" # PostgreSQL version

# Redis (optional - uncomment if needed)
# - name: %s-redis
#   type: redis
#   plan: starter
#   region: oregon
`, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name)
}

func (g *Generator) generateRenderReadme() string {
	return fmt.Sprintf(`# Render Deployment Guide

Simple, managed deployment for %s on Render.com.

## Overview

Render provides managed hosting with automatic deployments, SSL certificates, and database management - perfect for MVP and early-stage applications.

## Setup Instructions

### 1. Prerequisites

- GitHub repository with your code
- Render account (free tier available)
- Environment variables configured

### 2. Connect Repository

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click "New" → "Web Service" 
3. Connect your GitHub repository
4. Render will auto-detect the ` + "`render.yaml`" + ` configuration

### 3. Configure Environment Variables

Update the ` + "`render.yaml`" + ` file with your repository URL:

` + "```yaml" + `
services:
  - type: web
    name: %s-backend
    repo: https://github.com/your-username/your-repo
` + "```" + `

### 4. Environment Variables

**Backend Environment Variables:**
- ` + "`NODE_ENV`" + `: Set to production
- ` + "`DATABASE_URL`" + `: Auto-configured from database
- ` + "`JWT_SECRET`" + `: Auto-generated secure value
- ` + "`FRONTEND_URL`" + `: Frontend service URL

**Frontend Environment Variables:**
- ` + "`VITE_API_BASE_URL`" + `: Backend API URL
- ` + "`VITE_APP_NAME`" + `: Application name
- ` + "`VITE_NODE_ENV`" + `: production

### 5. Database Setup

Render will automatically:
- Create a PostgreSQL database
- Generate connection credentials
- Configure DATABASE_URL environment variable
- Handle backups and maintenance

### 6. Deployment Process

1. **Automatic Deployments**: Every push to ` + "`main`" + ` branch triggers deployment
2. **Build Process**: Render builds Docker images from your Dockerfiles
3. **Health Checks**: Verifies ` + "`/health`" + ` endpoint before routing traffic
4. **SSL**: Automatic HTTPS certificates

## Service URLs

After deployment, your services will be available at:

- **Backend**: ` + "`https://%s-backend.onrender.com`" + `
- **Frontend**: ` + "`https://%s-frontend.onrender.com`" + `
- **Database**: Managed PostgreSQL instance

## Domain Configuration

To use a custom domain:

1. Add domain to ` + "`render.yaml`" + `:
   ` + "```yaml" + `
   customDomains:
     - name: your-domain.com
   ` + "```" + `

2. Configure DNS records as instructed by Render

## Monitoring & Logs

- **Service Logs**: Available in Render dashboard
- **Metrics**: Built-in performance monitoring
- **Alerts**: Configure notifications for downtime

## Cost Optimization

**Free Tier Limitations:**
- Services sleep after 15 minutes of inactivity
- 750 hours/month limit
- Shared resources

**Upgrade Considerations:**
- Switch to ` + "`standard`" + ` plan for always-on services
- Use ` + "`pro`" + ` plan for production workloads
- Consider dedicated databases for high traffic

## Security Features

- **Automatic HTTPS**: SSL certificates managed by Render
- **Environment Isolation**: Each service runs in isolated containers
- **Secret Management**: Environment variables encrypted at rest
- **DDoS Protection**: Built-in protection against common attacks

## Backup & Recovery

- **Database Backups**: Automatic daily backups (retained for 7 days on free tier)
- **Point-in-time Recovery**: Available on paid plans
- **Manual Backups**: Can be triggered via dashboard

## Troubleshooting

**Common Issues:**

1. **Build Failures**: Check Dockerfile paths in ` + "`render.yaml`" + `
2. **Database Connection**: Verify DATABASE_URL is properly configured
3. **CORS Errors**: Ensure FRONTEND_URL is set correctly in backend
4. **Service Sleep**: Upgrade to paid plan for always-on services

**Debug Commands:**
` + "```bash" + `
# Check service logs
render logs --service=%s-backend

# View service status  
render services list

# Manual deploy
render deploy --service=%s-backend
` + "```" + `

## Next Steps

1. **Custom Domain**: Configure your domain name
2. **Monitoring**: Set up alerts and monitoring
3. **Scaling**: Consider upgrading to standard/pro plans
4. **AWS Migration**: When ready, migrate to AWS infrastructure

Generated with AI CodeKeeper v1.0.0 - Production-ready Render deployment
`, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name)
}

func (g *Generator) generateTerraformMain() string {
	return fmt.Sprintf(`# Main Terraform configuration for %s
# AWS infrastructure following Well-Architected Framework principles

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  # Configure S3 backend for state management
  backend "s3" {
    bucket         = "%s-terraform-state"
    key            = "infrastructure/terraform.tfstate"
    region         = var.aws_region
    encrypt        = true
    dynamodb_table = "%s-terraform-locks"
  }
}

# Configure AWS Provider
provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = "%s"
      Environment = var.environment
      ManagedBy   = "terraform"
      Domain      = "%s"
      CreatedBy   = "ai-codekeeper"
    }
  }
}

# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

# VPC Module
module "vpc" {
  source = "./modules/vpc"
  
  project_name        = var.project_name
  environment         = var.environment
  availability_zones  = data.aws_availability_zones.available.names
  vpc_cidr           = var.vpc_cidr
}

# ECS Module  
module "ecs" {
  source = "./modules/ecs"
  
  project_name     = var.project_name
  environment      = var.environment
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnet_ids
  public_subnets  = module.vpc.public_subnet_ids
  
  # Application configuration
  backend_image   = var.backend_image
  frontend_image  = var.frontend_image
  database_url    = module.rds.database_url
  
  # Secrets
  jwt_secret_arn = module.secrets.jwt_secret_arn
  
  depends_on = [module.vpc, module.rds, module.secrets]
}

# RDS Module
module "rds" {
  source = "./modules/rds"
  
  project_name     = var.project_name
  environment      = var.environment
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnet_ids
  
  # Database configuration
  database_name     = var.database_name
  master_username   = var.database_username
  master_password   = module.secrets.database_password
  instance_class    = var.database_instance_class
  allocated_storage = var.database_allocated_storage
  
  depends_on = [module.vpc, module.secrets]
}

# Secrets Module
module "secrets" {
  source = "./modules/secrets"
  
  project_name = var.project_name
  environment  = var.environment
}

# Outputs
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "ecs_cluster_name" {
  description = "ECS Cluster name"
  value       = module.ecs.cluster_name
}

output "application_load_balancer_dns" {
  description = "Application Load Balancer DNS name"
  value       = module.ecs.alb_dns_name
}

output "database_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.database_endpoint
  sensitive   = true
}

output "secrets_manager_arns" {
  description = "Secrets Manager ARNs"
  value = {
    jwt_secret       = module.secrets.jwt_secret_arn
    database_password = module.secrets.database_password_arn
  }
  sensitive = true
}
`, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Domain)
}

func (g *Generator) generateTerraformVariables() string {
	return fmt.Sprintf(`# Terraform variables for %s infrastructure

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "%s"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

# Application Configuration
variable "backend_image" {
  description = "Backend Docker image URI"
  type        = string
  default     = "%s-backend:latest"
}

variable "frontend_image" {
  description = "Frontend Docker image URI" 
  type        = string
  default     = "%s-frontend:latest"
}

# Database Configuration
variable "database_name" {
  description = "Name of the database"
  type        = string
  default     = "%s"
}

variable "database_username" {
  description = "Database master username"
  type        = string
  default     = "app_user"
}

variable "database_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro" # Free tier eligible
}

variable "database_allocated_storage" {
  description = "Database allocated storage in GB"
  type        = number
  default     = 20
}

# ECS Configuration
variable "backend_cpu" {
  description = "CPU units for backend service"
  type        = number
  default     = 256
}

variable "backend_memory" {
  description = "Memory for backend service"
  type        = number
  default     = 512
}

variable "frontend_cpu" {
  description = "CPU units for frontend service"
  type        = number
  default     = 256
}

variable "frontend_memory" {
  description = "Memory for frontend service"
  type        = number
  default     = 512
}

variable "backend_desired_count" {
  description = "Desired number of backend tasks"
  type        = number
  default     = 1
}

variable "frontend_desired_count" {
  description = "Desired number of frontend tasks"
  type        = number
  default     = 1
}

# Domain-specific variables for %s
variable "domain_specific_config" {
  description = "Domain-specific configuration for %s applications"
  type = object({
    enable_audit_logging = bool
    compliance_mode     = string
    encryption_required = bool
  })
  default = {
    enable_audit_logging = true
    compliance_mode     = "%s"
    encryption_required = true
  }
}

# Cost optimization
variable "enable_nat_gateway" {
  description = "Enable NAT Gateway (costs ~$45/month)"
  type        = bool
  default     = true
}

variable "multi_az_database" {
  description = "Enable Multi-AZ for database (increases cost)"
  type        = bool
  default     = false
}
`, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Domain, g.spec.Domain, strings.ToLower(g.spec.Domain))
}

func (g *Generator) generateTerraformOutputs() string {
	return fmt.Sprintf(`# Terraform outputs for %s infrastructure

# Network Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

# ECS Outputs
output "ecs_cluster_id" {
  description = "ID of the ECS cluster"
  value       = module.ecs.cluster_id
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = module.ecs.cluster_name
}

output "backend_service_name" {
  description = "Name of the backend ECS service"
  value       = module.ecs.backend_service_name
}

output "frontend_service_name" {
  description = "Name of the frontend ECS service"
  value       = module.ecs.frontend_service_name
}

# Load Balancer Outputs
output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.ecs.alb_dns_name
}

output "alb_zone_id" {
  description = "Zone ID of the Application Load Balancer"
  value       = module.ecs.alb_zone_id
}

output "backend_target_group_arn" {
  description = "ARN of the backend target group"
  value       = module.ecs.backend_target_group_arn
}

output "frontend_target_group_arn" {
  description = "ARN of the frontend target group"
  value       = module.ecs.frontend_target_group_arn
}

# Database Outputs
output "database_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.database_endpoint
  sensitive   = true
}

output "database_port" {
  description = "RDS instance port"
  value       = module.rds.database_port
}

output "database_name" {
  description = "Database name"
  value       = module.rds.database_name
}

# Security Outputs
output "database_security_group_id" {
  description = "ID of the database security group"
  value       = module.rds.security_group_id
}

output "ecs_security_group_id" {
  description = "ID of the ECS security group"
  value       = module.ecs.security_group_id
}

# Secrets Manager Outputs
output "jwt_secret_arn" {
  description = "ARN of the JWT secret in Secrets Manager"
  value       = module.secrets.jwt_secret_arn
  sensitive   = true
}

output "database_password_arn" {
  description = "ARN of the database password in Secrets Manager"
  value       = module.secrets.database_password_arn
  sensitive   = true
}

# Application URLs
output "application_urls" {
  description = "URLs to access the application"
  value = {
    backend_api    = "https://$${module.ecs.alb_dns_name}/api/v1"
    frontend_app   = "https://$${module.ecs.alb_dns_name}"
    health_check   = "https://$${module.ecs.alb_dns_name}/health"
  }
}

# Cost and Resource Summary
output "infrastructure_summary" {
  description = "Summary of deployed infrastructure"
  value = {
    environment     = var.environment
    region         = var.aws_region
    vpc_cidr       = var.vpc_cidr
    database_class = var.database_instance_class
    ecs_cluster    = module.ecs.cluster_name
    estimated_monthly_cost = "~$50-100 for t3.micro instances (excluding data transfer)"
  }
}
`, g.spec.Name)
}

// Terraform module placeholders (AI will implement full modules)
func (g *Generator) generateTerraformVPC() string {
	return fmt.Sprintf(`# VPC Module for %s
# AI Implementation Note: Complete VPC with public/private subnets, NAT Gateway, Internet Gateway

# Variables
variable "project_name" { type = string }
variable "environment" { type = string }
variable "availability_zones" { type = list(string) }
variable "vpc_cidr" { type = string }

# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {
    Name = "$${var.project_name}-$${var.environment}-vpc"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  
  tags = {
    Name = "$${var.project_name}-$${var.environment}-igw"
  }
}

# Public Subnets (AI: Implement across multiple AZs)
resource "aws_subnet" "public" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone = var.availability_zones[count.index]
  
  map_public_ip_on_launch = true
  
  tags = {
    Name = "$${var.project_name}-$${var.environment}-public-$${count.index + 1}"
    Type = "public"
  }
}

# Private Subnets (AI: Implement across multiple AZs)
resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 2)
  availability_zone = var.availability_zones[count.index]
  
  tags = {
    Name = "$${var.project_name}-$${var.environment}-private-$${count.index + 1}"
    Type = "private"
  }
}

# NAT Gateway (AI: Implement with EIP and routing)
# ... Additional VPC components to be implemented by AI

# Outputs
output "vpc_id" { value = aws_vpc.main.id }
output "vpc_cidr_block" { value = aws_vpc.main.cidr_block }
output "public_subnet_ids" { value = aws_subnet.public[*].id }
output "private_subnet_ids" { value = aws_subnet.private[*].id }
`, g.spec.Name)
}

func (g *Generator) generateTerraformECS() string {
	return fmt.Sprintf(`# ECS Module for %s
# AI Implementation Note: Complete ECS Fargate cluster with ALB, target groups, and services

# Variables (AI: Add all required variables)
variable "project_name" { type = string }
variable "environment" { type = string }
variable "vpc_id" { type = string }
variable "private_subnets" { type = list(string) }
variable "public_subnets" { type = list(string) }
variable "backend_image" { type = string }
variable "frontend_image" { type = string }
variable "database_url" { type = string }
variable "jwt_secret_arn" { type = string }

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "$${var.project_name}-$${var.environment}"
  
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
  
  tags = {
    Name = "$${var.project_name}-$${var.environment}-cluster"
  }
}

# Application Load Balancer (AI: Implement complete ALB configuration)
resource "aws_lb" "main" {
  name               = "$${var.project_name}-$${var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnets
  
  enable_deletion_protection = false
  
  tags = {
    Name = "$${var.project_name}-$${var.environment}-alb"
  }
}

# Security Groups (AI: Implement proper ingress/egress rules)
resource "aws_security_group" "alb" {
  name_prefix = "$${var.project_name}-alb-"
  vpc_id      = var.vpc_id
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ECS Task Definitions (AI: Implement complete task definitions for backend/frontend)
# ECS Services (AI: Implement services with proper scaling and health checks)
# Target Groups (AI: Implement ALB target groups)
# Listener Rules (AI: Implement routing rules)

# Outputs
output "cluster_id" { value = aws_ecs_cluster.main.id }
output "cluster_name" { value = aws_ecs_cluster.main.name }
output "alb_dns_name" { value = aws_lb.main.dns_name }
output "alb_zone_id" { value = aws_lb.main.zone_id }
output "security_group_id" { value = aws_security_group.alb.id }
# AI: Add additional outputs for services and target groups
`, g.spec.Name)
}

func (g *Generator) generateTerraformRDS() string {
	return `# RDS Module
# AI Implementation Note: Complete RDS PostgreSQL with security groups and subnet groups

# Variables (AI: Add all required variables)
variable "project_name" { type = string }
variable "environment" { type = string }
variable "vpc_id" { type = string }
variable "private_subnets" { type = list(string) }
variable "database_name" { type = string }
variable "master_username" { type = string }
variable "master_password" { type = string }
variable "instance_class" { type = string }
variable "allocated_storage" { type = number }

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}-db-subnet-group"
  subnet_ids = var.private_subnets
  
  tags = {
    Name = "${var.project_name}-${var.environment}-db-subnet-group"
  }
}

# Security Group for RDS
resource "aws_security_group" "rds" {
  name_prefix = "${var.project_name}-rds-"
  vpc_id      = var.vpc_id
  
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [] # AI: Reference ECS security groups
  }
  
  tags = {
    Name = "${var.project_name}-${var.environment}-rds-sg"
  }
}

# RDS Instance (AI: Implement complete RDS configuration)
resource "aws_db_instance" "main" {
  identifier = "${var.project_name}-${var.environment}-db"
  
  engine         = "postgres"
  engine_version = "15.3"
  instance_class = var.instance_class
  
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.allocated_storage * 2
  storage_type          = "gp2"
  storage_encrypted     = true
  
  db_name  = var.database_name
  username = var.master_username
  password = var.master_password
  
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  skip_final_snapshot = true # AI: Set to false for production
  
  tags = {
    Name = "${var.project_name}-${var.environment}-db"
  }
}

# Outputs
output "database_endpoint" { value = aws_db_instance.main.endpoint }
output "database_port" { value = aws_db_instance.main.port }
output "database_name" { value = aws_db_instance.main.db_name }
output "database_url" { 
  value = "postgresql://${var.master_username}:${var.master_password}@${aws_db_instance.main.endpoint}:${aws_db_instance.main.port}/${aws_db_instance.main.db_name}"
  sensitive = true 
}
output "security_group_id" { value = aws_security_group.rds.id }`
}

func (g *Generator) generateTerraformSecrets() string {
	return `# Secrets Manager Module
# AI Implementation Note: Implement secrets for JWT, database password, and other sensitive data

# Variables
variable "project_name" { type = string }
variable "environment" { type = string }

# JWT Secret
resource "aws_secretsmanager_secret" "jwt_secret" {
  name        = "${var.project_name}/${var.environment}/jwt-secret"
  description = "JWT secret for ${var.project_name} ${var.environment}"
  
  tags = {
    Name = "${var.project_name}-${var.environment}-jwt-secret"
  }
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id = aws_secretsmanager_secret.jwt_secret.id
  secret_string = jsonencode({
    jwt_secret = random_password.jwt_secret.result
  })
}

# Database Password
resource "aws_secretsmanager_secret" "database_password" {
  name        = "${var.project_name}/${var.environment}/database-password"
  description = "Database password for ${var.project_name} ${var.environment}"
  
  tags = {
    Name = "${var.project_name}-${var.environment}-db-password"
  }
}

resource "aws_secretsmanager_secret_version" "database_password" {
  secret_id = aws_secretsmanager_secret.database_password.id
  secret_string = jsonencode({
    password = random_password.database_password.result
  })
}

# Random Passwords
resource "random_password" "jwt_secret" {
  length  = 64
  special = true
}

resource "random_password" "database_password" {
  length  = 32
  special = true
}

# Outputs
output "jwt_secret_arn" { value = aws_secretsmanager_secret.jwt_secret.arn }
output "database_password_arn" { value = aws_secretsmanager_secret.database_password.arn }
output "database_password" { 
  value = random_password.database_password.result 
  sensitive = true 
}`
}

func (g *Generator) generateTerraformReadme() string {
	return fmt.Sprintf(`# AWS Infrastructure with Terraform

Enterprise-grade AWS infrastructure for %s following Well-Architected Framework principles.

## Architecture Overview

This Terraform configuration creates:

- **VPC**: Multi-AZ with public/private subnets
- **ECS Fargate**: Serverless container orchestration
- **RDS PostgreSQL**: Managed database with backups
- **Application Load Balancer**: High availability traffic routing
- **Secrets Manager**: Secure credential management
- **Security Groups**: Network-level security controls

## Prerequisites

1. **AWS CLI configured** with appropriate permissions
2. **Terraform >= 1.0** installed
3. **S3 bucket** for state storage
4. **DynamoDB table** for state locking

## Initial Setup

### 1. Create S3 Backend Resources

` + "```bash" + `
# Create S3 bucket for Terraform state
aws s3 mb s3://%s-terraform-state --region us-west-2

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket %s-terraform-state \
  --versioning-configuration Status=Enabled

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name %s-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --region us-west-2
` + "```" + `

### 2. Configure Variables

Create ` + "`terraform.tfvars`" + `:

` + "```hcl" + `
project_name = "%s"
environment = "dev"  # or "staging", "prod"
aws_region = "us-west-2"

# Database configuration
database_instance_class = "db.t3.micro"  # Free tier
database_allocated_storage = 20

# ECS configuration
backend_image = "your-account.dkr.ecr.us-west-2.amazonaws.com/%s-backend:latest"
frontend_image = "your-account.dkr.ecr.us-west-2.amazonaws.com/%s-frontend:latest"

# Cost optimization
enable_nat_gateway = true   # Set to false to save costs in dev
multi_az_database = false   # Set to true for production
` + "```" + `

## Deployment

### Development Environment

` + "```bash" + `
# Initialize Terraform
terraform init

# Plan deployment
terraform plan -var-file="terraform.tfvars"

# Apply (creates infrastructure)
terraform apply -var-file="terraform.tfvars"

# Get outputs
terraform output
` + "```" + `

### Production Environment

` + "```bash" + `
# Use workspace for environment isolation
terraform workspace new prod
terraform workspace select prod

# Use production variables
terraform plan -var-file="terraform.prod.tfvars"
terraform apply -var-file="terraform.prod.tfvars"
` + "```" + `

## Cost Optimization

**Free Tier Resources:**
- ` + "`db.t3.micro`" + ` RDS instance (750 hours/month)
- ` + "`t3.micro`" + ` ECS tasks (some free tier)
- ALB (750 hours/month)

**Cost-Saving Options:**
- Set ` + "`enable_nat_gateway = false`" + ` (saves ~$45/month)
- Use ` + "`db.t3.micro`" + ` for development
- Set ` + "`multi_az_database = false`" + ` for non-production

**Estimated Monthly Costs:**
- **Development**: $20-50/month
- **Production**: $100-200/month (depending on traffic)

## Security Features

- **Network Isolation**: Private subnets for application and database
- **Security Groups**: Restrictive ingress/egress rules
- **Encryption**: RDS encryption at rest, Secrets Manager for sensitive data
- **IAM Roles**: Least privilege access for ECS tasks
- **VPC Flow Logs**: Network traffic monitoring (optional)

## Monitoring & Observability

- **CloudWatch**: Container insights and custom metrics
- **ALB Logs**: Request logging and analysis
- **RDS Monitoring**: Database performance insights
- **Health Checks**: Application and database health monitoring

## Disaster Recovery

- **RDS Backups**: 7-day retention (configurable)
- **Multi-AZ**: Available for production databases
- **Infrastructure as Code**: Complete environment recreation
- **State Management**: S3 versioning and DynamoDB locking

## Domain-Specific Configuration

**%s Domain Settings:**
- Audit logging enabled by default
- Encryption required for sensitive data
- Compliance mode: %s
- Enhanced security group rules

## Troubleshooting

**Common Issues:**

1. **State Lock**: Run ` + "`terraform force-unlock LOCK_ID`" + `
2. **Permission Errors**: Verify AWS IAM permissions
3. **Resource Limits**: Check AWS service quotas
4. **Cost Alerts**: Set up billing alerts in AWS Console

**Debug Commands:**
` + "```bash" + `
# Validate configuration
terraform validate

# Check current state
terraform show

# Import existing resources
terraform import aws_instance.example i-1234567890abcdef0

# Destroy infrastructure (BE CAREFUL)
terraform destroy -var-file="terraform.tfvars"
` + "```" + `

## Next Steps

1. **Configure Domain**: Set up Route 53 for custom domain
2. **SSL Certificate**: Add ACM certificate to ALB
3. **Monitoring**: Set up CloudWatch dashboards
4. **Backup Strategy**: Configure automated backups
5. **CI/CD Integration**: Update deployment pipeline

Generated with AI CodeKeeper v1.0.0 - Production-ready AWS infrastructure
`, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Domain, strings.ToLower(g.spec.Domain))
}

func (g *Generator) generateEnhancedCIWorkflow() string {
	// Replace the existing basic CI workflow
	return g.generateCIWorkflow()
}

func (g *Generator) generateDeploymentWorkflow() string {
	return fmt.Sprintf(`# Deployment workflow for %s
name: Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  AWS_REGION: us-west-2
  ECR_REGISTRY: ${{ secrets.AWS_ACCOUNT_ID }}.dkr.ecr.us-west-2.amazonaws.com
  IMAGE_TAG: ${{ github.sha }}

jobs:
  deploy-render:
    name: Deploy to Render
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - name: Checkout
      uses: actions/checkout@v4
    
    - name: Deploy to Render
      run: |
        echo "Render auto-deploys on main branch push"
        echo "Monitor deployment at: https://dashboard.render.com"

  build-and-push:
    name: Build and Push Docker Images
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    outputs:
      backend-image: ${{ steps.build-backend.outputs.image }}
      frontend-image: ${{ steps.build-frontend.outputs.image }}
    
    steps:
    - name: Checkout
      uses: actions/checkout@v4
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v4
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: ${{ env.AWS_REGION }}
    
    - name: Login to ECR
      id: login-ecr
      uses: aws-actions/amazon-ecr-login@v2
    
    - name: Build and push backend image
      id: build-backend
      run: |
        cd apps/backend
        docker build -t $ECR_REGISTRY/%s-backend:$IMAGE_TAG .
        docker push $ECR_REGISTRY/%s-backend:$IMAGE_TAG
        echo "image=$ECR_REGISTRY/%s-backend:$IMAGE_TAG" >> $GITHUB_OUTPUT
    
    - name: Build and push frontend image
      id: build-frontend
      run: |
        cd apps/frontend
        docker build -t $ECR_REGISTRY/%s-frontend:$IMAGE_TAG .
        docker push $ECR_REGISTRY/%s-frontend:$IMAGE_TAG
        echo "image=$ECR_REGISTRY/%s-frontend:$IMAGE_TAG" >> $GITHUB_OUTPUT

  deploy-aws:
    name: Deploy to AWS (Manual)
    runs-on: ubuntu-latest
    needs: build-and-push
    if: false  # Set to true when ready for AWS deployment
    environment: production
    
    steps:
    - name: Checkout
      uses: actions/checkout@v4
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v4
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: ${{ env.AWS_REGION }}
    
    - name: Setup Terraform
      uses: hashicorp/setup-terraform@v3
      with:
        terraform_version: 1.5.0
    
    - name: Terraform Init
      run: |
        cd infra/aws
        terraform init
    
    - name: Terraform Plan
      run: |
        cd infra/aws
        terraform plan \
          -var="backend_image=${{ needs.build-and-push.outputs.backend-image }}" \
          -var="frontend_image=${{ needs.build-and-push.outputs.frontend-image }}" \
          -var="environment=prod"
    
    - name: Terraform Apply
      if: github.ref == 'refs/heads/main'
      run: |
        cd infra/aws
        terraform apply -auto-approve \
          -var="backend_image=${{ needs.build-and-push.outputs.backend-image }}" \
          -var="frontend_image=${{ needs.build-and-push.outputs.frontend-image }}" \
          -var="environment=prod"
`, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name, g.spec.Name)
}

// Final integration files
func (g *Generator) generateFinalIntegrationFiles() map[string]string {
	files := make(map[string]string)
	
	// Enhanced documentation
	files["docs/10_AI_IMPLEMENTATION_GUIDE.md"] = g.generateAIImplementationGuide()
	files["docs/11_PRODUCTION_CHECKLIST.md"] = g.generateProductionChecklist()
	
	// Root-level integration files
	files[".gitattributes"] = g.generateGitAttributes()
	files[".editorconfig"] = g.generateEditorConfig()
	files["package.json"] = g.generateRootPackageJson()
	
	// Development environment
	files["docs/09_LOCAL_DEVELOPMENT.md"] = g.generateLocalDevelopmentDoc()
	
	// AI CodeKeeper configuration
	files[".codekeeper/templates/README.md"] = g.generateTemplatesReadme()
	files[".codekeeper/prompts/implementation.md"] = g.generateImplementationPrompts()
	
	return files
}

func (g *Generator) generateAIImplementationGuide() string {
	return fmt.Sprintf(`# AI Implementation Guide

Comprehensive guide for AI systems implementing %s using the generated framework.

## Overview

This project has been scaffolded with the **AI CodeKeeper framework** following the AI Master Prompt specification. All foundation components are in place with clear implementation guidelines.

## Implementation Priorities

### Phase 1: Foundation Verification (1-2 hours)
1. **Environment Setup**
   - Verify all dependencies are correctly installed
   - Test Docker Compose setup: ` + "`docker-compose up`" + `
   - Confirm database connections and migrations

2. **Basic Functionality Test**
   - Backend health check: ` + "`curl http://localhost:8080/health`" + `
   - Frontend loads correctly: ` + "`http://localhost:3000`" + `
   - Environment variables are properly configured

### Phase 2: Backend Implementation (4-6 hours)
1. **Authentication System**
   - Complete user registration/login endpoints
   - JWT token generation and validation
   - Password hashing and security measures
   - Role-based access control middleware

2. **%s CRUD Operations**
   - Database models with %s domain validation
   - Complete REST API endpoints (GET, POST, PUT, DELETE)
   - Input validation and error handling
   - Pagination and filtering capabilities

3. **Domain-Specific Logic**
%s

4. **Testing Implementation**
   - Unit tests for services and utilities
   - Integration tests for API endpoints
   - Test data fixtures and mocks

### Phase 3: Frontend Implementation (4-6 hours)
1. **Authentication UI**
   - Login and registration forms
   - Protected route handling
   - User session management
   - Role-based UI components

2. **%s Management Interface**
   - List view with pagination and search
   - Create/edit forms with validation
   - Detail views with actions
   - Delete confirmation dialogs

3. **State Management**
   - Zustand stores for auth and %s data
   - API integration with error handling
   - Loading states and optimistic updates
   - Persistent authentication state

4. **UI/UX Polish**
   - Responsive design across screen sizes
   - Error boundaries and user feedback
   - Loading indicators and empty states
   - Accessibility improvements

### Phase 4: Production Readiness (2-4 hours)
1. **Security Review**
   - Input validation completeness
   - Authentication security audit
   - Environment variable protection
   - SQL injection prevention verification

2. **Performance Optimization**
   - Database query optimization
   - Frontend bundle size analysis
   - API response time optimization
   - Caching strategy implementation

3. **Monitoring Setup**
   - Application logging configuration
   - Error tracking integration
   - Health check endpoints
   - Performance monitoring

## Implementation Guidelines

### Backend Development (%s)

**File Structure to Implement:**
` + "```" + `
apps/backend/src/
├── config/           ✅ Generated
├── controllers/      🔄 Implement CRUD operations
├── services/         🔄 Implement business logic
├── models/           🔄 Complete database models
├── middleware/       🔄 Complete auth middleware
├── routes/           🔄 Implement route handlers
├── utils/            ✅ Generated helpers
└── tests/            🔄 Implement test suite
` + "```" + `

**Key Implementation Points:**
- Follow the domain-specific patterns in generated service stubs
- Use the provided validation utilities and extend as needed
- Implement proper error handling with user-friendly messages
- Add comprehensive logging for debugging and monitoring
- Ensure all database operations use the ORM properly

### Frontend Development (React + TypeScript)

**Component Structure to Implement:**
` + "```" + `
apps/frontend/src/
├── components/       🔄 Implement UI components
│   ├── ui/          🔄 Basic elements (Button, Input, Modal)
│   ├── forms/       🔄 Form components with validation
│   └── layout/      🔄 Header, Sidebar, Navigation
├── pages/           🔄 Implement page components
├── hooks/           🔄 Custom React hooks
├── services/        🔄 API client implementation
├── store/           🔄 Zustand state management
└── types/           ✅ Generated (use shared-types)
` + "```" + `

**Key Implementation Points:**
- Use the shared TypeScript types from the packages/shared-types
- Follow the component patterns established in the basic App.tsx
- Implement proper form validation and error handling
- Use the configured Axios instance for all API calls
- Follow accessibility best practices (WCAG 2.1 AA)

### Database Design

**%s Entity Schema:**
` + "```sql" + `
-- Implement based on domain requirements
-- Use the generated model as starting point
-- Add proper indexes and constraints
-- Consider audit trail requirements for %s domain
` + "```" + `

### API Design Patterns

**RESTful Endpoints:**
` + "```" + `
GET    /api/v1/%s                # List with pagination
POST   /api/v1/%s                # Create new
GET    /api/v1/%s/:id            # Get by ID  
PUT    /api/v1/%s/:id            # Update
DELETE /api/v1/%s/:id            # Delete (soft delete preferred)

# Authentication
POST   /api/v1/auth/register     # User registration
POST   /api/v1/auth/login        # User login
GET    /api/v1/auth/profile      # User profile
PUT    /api/v1/auth/profile      # Update profile
` + "```" + `

## Testing Strategy

### Backend Testing
- **Unit Tests**: Services, utilities, validation functions
- **Integration Tests**: API endpoints with database interaction
- **Security Tests**: Authentication, authorization, input validation
- **Domain Tests**: Business logic specific to %s

### Frontend Testing  
- **Component Tests**: UI components with React Testing Library
- **Integration Tests**: User flows and API integration
- **Accessibility Tests**: Screen reader compatibility
- **Visual Tests**: Responsive design across devices

## Deployment Strategy

### Development/Staging (Render)
1. Connect GitHub repository to Render
2. Configure environment variables in Render dashboard
3. Deploy automatically on push to main branch
4. Monitor deployment and logs

### Production (AWS)
1. Set up AWS credentials and ECR repository
2. Configure Terraform backend (S3 + DynamoDB)
3. Deploy infrastructure: ` + "`terraform apply`" + `
4. Deploy application via GitHub Actions

## Domain-Specific Considerations

### %s Domain Requirements
%s

## Success Metrics

**Technical Metrics:**
- All tests pass with >80%% coverage
- API response times <200ms for CRUD operations
- Frontend bundle size <1MB gzipped
- Zero security vulnerabilities in dependencies

**User Experience Metrics:**
- Complete user flows work without errors
- Responsive design works on mobile/tablet/desktop
- Loading states and error handling provide good UX
- Authentication and authorization work correctly

**Production Readiness:**
- Health checks respond correctly
- Logging provides adequate debugging information
- Error handling covers edge cases
- Database performance is optimized

## Getting Help

**Generated Documentation:**
- Architecture overview: ` + "`docs/02_ARCHITECTURE.md`" + `
- API specifications: ` + "`docs/05_API_ENDPOINTS.md`" + `
- Deployment guide: ` + "`docs/06_DEPLOYMENT.md`" + `
- Local development: ` + "`docs/09_LOCAL_DEVELOPMENT.md`" + `

**Code Examples:**
- Backend: Check generated service and controller stubs
- Frontend: Review App.tsx and component structure
- Database: Use generated model definitions as starting point
- Testing: Follow patterns in generated test files

Generated with AI CodeKeeper v1.0.0 - Ready for AI implementation
`, g.spec.Name, g.spec.CoreEntity, g.spec.Domain, g.getDomainImplementationGuidelines(),
   g.spec.CoreEntity, g.spec.CoreEntity, g.spec.Backend, g.spec.CoreEntity, g.spec.Domain,
   strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), 
   strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), strings.ToLower(g.spec.CoreEntity), g.spec.Domain,
   g.spec.Domain, g.getDomainProductionRequirements())
}

func (g *Generator) generateProductionChecklist() string {
	return fmt.Sprintf(`# Production Readiness Checklist

Comprehensive checklist for deploying %s to production.

## Pre-Deployment Checklist

### ✅ Security Audit
- [ ] All environment variables are properly configured
- [ ] No secrets or API keys committed to repository
- [ ] JWT secret is cryptographically secure (64+ random characters)
- [ ] Database passwords are strong and randomly generated
- [ ] All API endpoints have proper authentication
- [ ] Input validation covers all user inputs
- [ ] SQL injection protection verified
- [ ] XSS protection implemented
- [ ] CORS settings are restrictive for production
- [ ] Rate limiting is configured appropriately

### ✅ Performance Verification
- [ ] API response times <200ms for CRUD operations
- [ ] Database queries are optimized with proper indexes
- [ ] Frontend bundle size <1MB gzipped
- [ ] Images and assets are optimized
- [ ] Caching headers are properly configured
- [ ] Database connection pooling is configured
- [ ] No memory leaks in long-running processes

### ✅ Testing Completion
- [ ] All unit tests pass with >80%% code coverage
- [ ] Integration tests cover critical user flows
- [ ] End-to-end tests verify complete functionality
- [ ] Load testing confirms performance under expected traffic
- [ ] Security testing validates auth and input handling
- [ ] Browser compatibility testing completed
- [ ] Mobile responsiveness verified

### ✅ Monitoring & Observability
- [ ] Application logging is comprehensive and structured
- [ ] Error tracking is configured (Sentry, etc.)
- [ ] Health check endpoints respond correctly
- [ ] Database monitoring is set up
- [ ] Performance metrics are tracked
- [ ] Alerting is configured for critical issues
- [ ] Log aggregation is properly configured

### ✅ Documentation
- [ ] API documentation is complete and accurate
- [ ] Deployment procedures are documented
- [ ] Environment setup instructions are verified
- [ ] Troubleshooting guides are available
- [ ] Security policies are documented
- [ ] Backup and recovery procedures are defined

## Deployment Strategy

### Render Deployment (Staging/MVP)

**Pre-Deployment:**
- [ ] GitHub repository is properly configured
- [ ] ` + "`render.yaml`" + ` has correct repository URL
- [ ] Environment variables are set in Render dashboard
- [ ] Database migration strategy is planned

**Deployment Steps:**
1. [ ] Connect repository to Render
2. [ ] Configure environment variables
3. [ ] Deploy and verify health checks
4. [ ] Test complete user flows
5. [ ] Monitor logs for any issues

**Post-Deployment:**
- [ ] Custom domain configured (if applicable)
- [ ] SSL certificate is working
- [ ] Monitoring alerts are active
- [ ] Backup schedule is confirmed

### AWS Production Deployment

**Infrastructure Setup:**
- [ ] AWS account and credentials configured
- [ ] S3 bucket for Terraform state created
- [ ] DynamoDB table for state locking created
- [ ] ECR repositories for container images created
- [ ] Domain and SSL certificate ready (if applicable)

**Terraform Deployment:**
- [ ] ` + "`terraform init`" + ` completed successfully
- [ ] ` + "`terraform plan`" + ` reviewed and approved
- [ ] ` + "`terraform apply`" + ` completed without errors
- [ ] All AWS resources created correctly
- [ ] Security groups have minimal required access
- [ ] Secrets Manager contains all required secrets

**Application Deployment:**
- [ ] Docker images built and pushed to ECR
- [ ] ECS services deployed successfully
- [ ] Load balancer health checks passing
- [ ] Database connectivity verified
- [ ] All environment variables configured

**Post-Deployment Verification:**
- [ ] Application accessible via load balancer
- [ ] All API endpoints responding correctly
- [ ] Database operations working
- [ ] Authentication flows working
- [ ] Frontend application loads correctly
- [ ] HTTPS/SSL working properly

## Domain-Specific Requirements

### %s Domain Checklist
%s

## Performance Benchmarks

### Backend Performance
- [ ] Health check endpoint: <50ms response time
- [ ] Authentication endpoints: <100ms response time
- [ ] CRUD operations: <200ms response time
- [ ] List endpoints with pagination: <300ms response time
- [ ] Database query optimization verified
- [ ] Connection pool properly sized

### Frontend Performance
- [ ] Initial page load: <3 seconds
- [ ] Route navigation: <1 second
- [ ] Bundle size optimized (<1MB gzipped)
- [ ] Critical path rendering optimized
- [ ] Lazy loading implemented for routes
- [ ] Images and assets optimized

### Infrastructure Performance
- [ ] Load balancer response times acceptable
- [ ] Auto-scaling configured and tested
- [ ] Database performance monitoring active
- [ ] CDN configured for static assets (if applicable)

## Security Verification

### Authentication & Authorization
- [ ] JWT tokens have appropriate expiration
- [ ] Refresh token rotation implemented
- [ ] Role-based access control working
- [ ] Session management secure
- [ ] Password requirements enforced
- [ ] Account lockout policies configured

### Data Protection
- [ ] All sensitive data encrypted at rest
- [ ] Data encryption in transit (HTTPS)
- [ ] Database access properly restricted
- [ ] Backup encryption verified
- [ ] Data retention policies implemented
- [ ] GDPR compliance verified (if applicable)

### Infrastructure Security
- [ ] Security groups follow least privilege
- [ ] Network ACLs properly configured
- [ ] VPC configuration secure
- [ ] IAM roles follow least privilege
- [ ] Secrets management properly implemented
- [ ] Security scanning completed

## Monitoring & Alerting

### Critical Alerts
- [ ] Application down/unhealthy
- [ ] Database connectivity issues
- [ ] High error rates (>5%% of requests)
- [ ] Response time degradation
- [ ] Security incidents
- [ ] Resource exhaustion (CPU, memory, disk)

### Business Metrics
- [ ] User registration/login success rates
- [ ] %s creation/modification rates
- [ ] API endpoint usage patterns
- [ ] Error patterns and resolution
- [ ] Performance trends over time

## Backup & Recovery

### Data Backup
- [ ] Database backups scheduled and tested
- [ ] Backup retention policy configured
- [ ] Point-in-time recovery tested
- [ ] Cross-region backup strategy (for critical data)
- [ ] Backup restoration procedures documented

### Disaster Recovery
- [ ] Recovery time objective (RTO) defined
- [ ] Recovery point objective (RPO) defined
- [ ] Disaster recovery procedures documented
- [ ] Infrastructure recreation procedures tested
- [ ] Data recovery procedures verified

## Compliance & Legal

### %s Domain Compliance
%s

### General Compliance
- [ ] Privacy policy updated and accessible
- [ ] Terms of service updated
- [ ] Data processing agreements in place
- [ ] Security incident response plan documented
- [ ] Audit logging implemented for compliance

## Launch Communication

### Internal Team
- [ ] Launch timeline communicated
- [ ] Support procedures documented
- [ ] Escalation paths defined
- [ ] Post-launch monitoring plan shared
- [ ] Rollback procedures documented

### External (if applicable)
- [ ] Customer communication prepared
- [ ] Support documentation updated
- [ ] Marketing materials ready
- [ ] Analytics tracking configured

## Post-Launch Tasks

### First 24 Hours
- [ ] Monitor all critical metrics
- [ ] Verify user registration/login flows
- [ ] Check error rates and response times
- [ ] Monitor resource utilization
- [ ] Review logs for any issues

### First Week
- [ ] Analyze user behavior and usage patterns
- [ ] Review performance metrics and optimize
- [ ] Address any user feedback or issues
- [ ] Fine-tune monitoring and alerting
- [ ] Document lessons learned

### First Month
- [ ] Capacity planning based on actual usage
- [ ] Security review and audit
- [ ] Performance optimization based on real data
- [ ] User feedback analysis and improvements
- [ ] Cost optimization review

---

**Generated with AI CodeKeeper v1.0.0**
**Complete Production Readiness Framework for %s**
`, g.spec.Name, g.spec.Domain, g.getDomainComplianceChecklist(), g.spec.CoreEntity, g.spec.Domain, g.getDomainComplianceRequirements(), g.spec.Name)
}

// Helper functions for final integration
func (g *Generator) getDomainSpecificFeatures() string {
	switch g.spec.Domain {
	case "fintech":
		return `- ✅ Decimal arithmetic for monetary calculations
- ✅ Comprehensive audit trails for transactions
- ✅ PCI compliance considerations
- ✅ Fraud detection patterns and rate limiting
- ✅ Multi-currency support and validation
- ✅ Regulatory compliance logging`
	case "healthcare":
		return `- ✅ HIPAA compliance for patient data
- ✅ Audit logging for data access
- ✅ Encryption for PHI (Personal Health Information)
- ✅ Patient consent management
- ✅ Medical record security and privacy
- ✅ Healthcare professional access controls`
	case "ecommerce":
		return `- ✅ Product catalog management
- ✅ Inventory tracking and management
- ✅ Shopping cart and checkout flows
- ✅ Payment processing integration
- ✅ Order management and tracking
- ✅ Customer review and rating systems`
	default:
		return `- ✅ Domain-specific business logic validation
- ✅ Entity lifecycle management
- ✅ User access controls and permissions
- ✅ Data integrity and validation rules
- ✅ Search and filtering capabilities
- ✅ Audit logging for critical operations`
	}
}

func (g *Generator) getGeneratedFileCount() string {
	return "50+" // Approximate count of all generated files
}

func (g *Generator) getDomainImplementationGuidelines() string {
	switch g.spec.Domain {
	case "fintech":
		return `   **Financial Domain Implementation:**
   - Use Decimal.js for all monetary calculations (never use floating point)
   - Implement comprehensive audit trails for all transactions
   - Add idempotency keys for transaction endpoints
   - Implement proper error handling for payment failures
   - Add rate limiting specific to transaction endpoints
   - Use encryption for sensitive financial data
   - Implement proper money transfer validation
   - Add compliance logging for regulatory requirements`
	case "healthcare":
		return `   **Healthcare Domain Implementation:**
   - Ensure HIPAA compliance for all patient data handling
   - Implement comprehensive audit logging for data access
   - Add proper consent management workflows
   - Use encryption for PHI (Personal Health Information)
   - Implement proper data retention policies
   - Add patient data anonymization capabilities
   - Ensure role-based access for medical staff
   - Add secure data transmission protocols`
	default:
		return `   **Domain-Specific Implementation:**
   - Follow established business logic patterns
   - Implement proper data validation rules
   - Add domain-specific search and filtering
   - Ensure proper audit logging
   - Implement entity lifecycle management
   - Add appropriate user access controls`
	}
}

func (g *Generator) getDomainProductionRequirements() string {
	switch g.spec.Domain {
	case "fintech":
		return `**Financial Applications:**
- PCI DSS compliance for payment data
- SOX compliance for financial reporting
- Anti-money laundering (AML) procedures
- Know Your Customer (KYC) verification
- Fraud detection and prevention
- Regulatory reporting capabilities`
	case "healthcare":
		return `**Healthcare Applications:**
- HIPAA compliance for patient data
- FDA regulations for medical devices
- HL7 FHIR standards for interoperability
- Patient consent and privacy controls
- Medical record retention policies
- Healthcare provider certification requirements`
	default:
		return `**General Applications:**
- Data privacy and protection
- User consent management
- Audit trail requirements
- Security best practices
- Performance and scalability
- Monitoring and observability`
	}
}

func (g *Generator) getDomainComplianceChecklist() string {
	switch g.spec.Domain {
	case "fintech":
		return `- [ ] PCI DSS compliance verified for payment processing
- [ ] SOX compliance for financial data integrity
- [ ] AML/KYC procedures implemented
- [ ] Regulatory reporting mechanisms in place
- [ ] Fraud detection algorithms active
- [ ] Financial audit trails complete
- [ ] Currency handling properly implemented
- [ ] Transaction idempotency verified`
	case "healthcare":
		return `- [ ] HIPAA compliance audit completed
- [ ] Patient data encryption verified
- [ ] Audit logging for PHI access implemented
- [ ] Patient consent mechanisms working
- [ ] Medical record retention policies active
- [ ] Healthcare provider access controls verified
- [ ] HL7 FHIR compliance (if applicable)
- [ ] Medical device regulations compliance (if applicable)`
	default:
		return `- [ ] Data privacy compliance verified
- [ ] User consent mechanisms implemented
- [ ] Audit trail completeness verified
- [ ] Security best practices followed
- [ ] Performance requirements met
- [ ] Monitoring and alerting active`
	}
}

func (g *Generator) getDomainComplianceRequirements() string {
	switch g.spec.Domain {
	case "fintech":
		return `- **PCI DSS**: Payment card data protection
- **SOX**: Financial data integrity and controls
- **AML/BSA**: Anti-money laundering compliance
- **KYC**: Customer identification requirements
- **GDPR/CCPA**: Data privacy and protection
- **Regional**: Local financial regulations`
	case "healthcare":
		return `- **HIPAA**: Patient data privacy and security
- **HITECH**: Health information technology compliance
- **FDA**: Medical device regulations (if applicable)
- **HL7 FHIR**: Healthcare interoperability standards
- **GDPR**: Patient data protection
- **Regional**: Local healthcare regulations`
	default:
		return `- **GDPR**: General data protection regulation
- **CCPA**: California consumer privacy act
- **SOC 2**: Security and availability controls
- **ISO 27001**: Information security management
- **Regional**: Local data protection laws`
	}
}

// Additional integration files
func (g *Generator) generateGitAttributes() string {
	return `# Git attributes for consistent line endings and file handling

# Auto detect text files and perform LF normalization
* text=auto

# Explicitly declare text files you want to always be normalized
*.md text
*.js text
*.ts text
*.tsx text
*.json text
*.yaml text
*.yml text
*.css text
*.html text
*.sql text

# Declare files that will always have CRLF line endings on checkout
*.bat text eol=crlf

# Denote all files that are truly binary and should not be modified
*.png binary
*.jpg binary
*.jpeg binary
*.gif binary
*.ico binary
*.mov binary
*.mp4 binary
*.mp3 binary
*.flv binary
*.fla binary
*.swf binary
*.gz binary
*.zip binary
*.7z binary
*.ttf binary
*.eot binary
*.woff binary
*.woff2 binary
*.pyc binary
*.pdf binary
*.ez binary
*.bz2 binary
*.swp binary
*.lz binary

# Docker files
Dockerfile* text
*.dockerignore text

# Git files
.gitignore text
.gitattributes text

# CI/CD files
.github/** text

# Terraform files
*.tf text
*.tfvars text
*.hcl text`
}

func (g *Generator) generateEditorConfig() string {
	return `# EditorConfig is awesome: https://EditorConfig.org

# top-most EditorConfig file
root = true

# Unix-style newlines with a newline ending every file
[*]
end_of_line = lf
insert_final_newline = true
charset = utf-8
trim_trailing_whitespace = true

# TypeScript, JavaScript, JSON
[*.{ts,tsx,js,jsx,json}]
indent_style = space
indent_size = 2

# CSS, SCSS, LESS
[*.{css,scss,less}]
indent_style = space
indent_size = 2

# HTML, XML
[*.{html,xml}]
indent_style = space
indent_size = 2

# YAML
[*.{yml,yaml}]
indent_style = space
indent_size = 2

# Markdown
[*.md]
trim_trailing_whitespace = false

# Python
[*.py]
indent_style = space
indent_size = 4

# Go
[*.go]
indent_style = tab
indent_size = 4

# Terraform
[*.{tf,tfvars}]
indent_style = space
indent_size = 2

# Docker
[Dockerfile*]
indent_style = space
indent_size = 2

# Shell scripts
[*.{sh,bash,zsh}]
indent_style = space
indent_size = 2

# Make files
[Makefile]
indent_style = tab

# Configuration files
[*.{toml,ini,cfg,conf}]
indent_style = space
indent_size = 2`
}

func (g *Generator) generateRootPackageJson() string {
	return fmt.Sprintf(`{
  "name": "%s-monorepo",
  "version": "1.0.0",
  "description": "%s - Complete full-stack application",
  "private": true,
  "workspaces": [
    "apps/*",
    "packages/*"
  ],
  "scripts": {
    "dev": "concurrently \"npm run dev --workspace=apps/backend\" \"npm run dev --workspace=apps/frontend\"",
    "build": "npm run build --workspaces",
    "test": "npm run test --workspaces",
    "lint": "npm run lint --workspaces",
    "type-check": "npm run type-check --workspaces",
    "clean": "npm run clean --workspaces && rimraf node_modules",
    "docker:up": "docker-compose up -d",
    "docker:down": "docker-compose down",
    "docker:logs": "docker-compose logs -f",
    "check": "codekeeper check --project-root ."
  },
  "devDependencies": {
    "concurrently": "^8.2.0",
    "rimraf": "^5.0.1"
  },
  "engines": {
    "node": ">=18.0.0",
    "npm": ">=9.0.0"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/your-username/%s.git"
  },
  "keywords": [
    "%s",
    "monorepo",
    "typescript",
    "react",
    "%s",
    "fullstack",
    "ai-generated"
  ],
  "author": "Generated with AI CodeKeeper",
  "license": "MIT"
}`, g.spec.Name, g.spec.Description, g.spec.Name, g.spec.Domain, g.spec.Backend)
}

func (g *Generator) generateTemplatesReadme() string {
	return `# AI CodeKeeper Templates

This directory contains templates and patterns for AI implementation.

## Purpose

These templates provide standardized patterns for AI systems to follow when implementing features in this codebase.

## Usage

AI systems should reference these templates when:
- Adding new API endpoints
- Creating new React components
- Implementing database models
- Writing tests
- Adding documentation

## Template Categories

- **Backend**: API endpoints, services, models
- **Frontend**: Components, hooks, pages
- **Database**: Migrations, seeds, models
- **Testing**: Unit tests, integration tests
- **Documentation**: API docs, README updates

Generated with AI CodeKeeper v1.0.0`
}

func (g *Generator) generateImplementationPrompts() string {
	return fmt.Sprintf(`# Implementation Prompts for AI Systems

## Context

This project: %s
Domain: %s
Core Entity: %s
Backend: %s

## Implementation Guidelines

When implementing features in this codebase:

1. **Follow established patterns** in generated code
2. **Use domain-specific validation** for %s applications
3. **Maintain type safety** with TypeScript throughout
4. **Follow 12-Factor App principles** for configuration and deployment
5. **Implement comprehensive testing** with >80%% coverage
6. **Use structured logging** for debugging and monitoring
7. **Follow security best practices** for authentication and data handling

## Code Standards

- Use the generated TypeScript types from shared-types package
- Follow the established project structure
- Implement proper error handling with user-friendly messages
- Use the configured development tools (ESLint, Prettier, etc.)
- Maintain consistent code style across the project

## Domain Requirements

%s

Generated with AI CodeKeeper v1.0.0`, g.spec.Name, g.spec.Domain, g.spec.CoreEntity, g.spec.Backend, g.spec.Domain, g.getDomainImplementationGuidelines())
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