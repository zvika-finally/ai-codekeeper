package generator

import (
	"encoding/json"
	"fmt"
	"os"
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
	files := map[string]string{
		"README.md": g.generateReadme(),
		"apps/backend/package.json": g.generateBackendPackageJson(),
		"apps/frontend/package.json": g.generateFrontendPackageJson(),
		"packages/shared-types/package.json": g.generateSharedTypesPackageJson(),
		"docker-compose.yml": g.generateDockerCompose(),
		".devcontainer/devcontainer.json": g.generateDevContainer(),
		"docs/00_OVERVIEW.md": g.generateOverviewDoc(),
		".github/workflows/ci.yml": g.generateCIWorkflow(),
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
	framework := strings.ToLower(g.spec.GetBackendFramework())
	
	return fmt.Sprintf(`{
  "name": "@%s/backend",
  "version": "1.0.0",
  "description": "%s backend API",
  "main": "src/server.js",
  "scripts": {
    "start": "node src/server.js",
    "dev": "nodemon src/server.js",
    "test": "jest",
    "lint": "eslint src/",
    "check": "codekeeper check"
  },
  "dependencies": {
    "%s": "latest"
  },
  "devDependencies": {
    "nodemon": "latest",
    "jest": "latest",
    "eslint": "latest"
  }
}`, g.spec.Name, g.spec.Description, framework)
}

func (g *Generator) generateFrontendPackageJson() string {
	return fmt.Sprintf(`{
  "name": "@%s/frontend",
  "version": "1.0.0",
  "description": "%s frontend application",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "check": "codekeeper check"
  },
  "dependencies": {
    "next": "latest",
    "react": "latest",
    "react-dom": "latest"
  },
  "devDependencies": {
    "typescript": "latest",
    "@types/react": "latest",
    "eslint": "latest"
  }
}`, g.spec.Name, g.spec.Description)
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
	services := `version: '3.8'

services:
  backend:
    build: ./apps/backend
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=development
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
	return `name: CI

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
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run tests
      run: npm test
    
    - name: Run guard rails
      run: codekeeper check --ci
    
    - name: Build
      run: npm run build`
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