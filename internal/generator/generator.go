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
	files := map[string]string{
		"README.md": g.generateReadme(),
		"apps/backend/package.json": g.generateBackendPackageJson(),
		"apps/frontend/package.json": g.generateFrontendPackageJson(),
		"packages/shared-types/package.json": g.generateSharedTypesPackageJson(),
		"docker-compose.yml": g.generateDockerCompose(),
		".devcontainer/devcontainer.json": g.generateDevContainer(),
		"docs/00_OVERVIEW.md": g.generateOverviewDoc(),
		".github/workflows/ci.yml": g.generateCIWorkflow(),
		"LICENSE": g.generateLicense(),
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