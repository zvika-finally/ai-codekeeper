package backend

import (
	"strings"
)

// Generator handles backend code generation
type BackendGenerator struct {
	spec *ProjectSpec
}

// ProjectSpec represents the project specification (will be imported from parent package)
type ProjectSpec struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CoreEntity  string   `json:"core_entity"`
	Backend     string   `json:"backend"`
	Databases   []string `json:"databases"`
	APIStyle    string   `json:"api_style"`
	UserRoles   string   `json:"user_roles"`
	Domain      string   `json:"domain"`
	ProjectPath string   `json:"project_path,omitempty"`
}

// NewBackendGenerator creates a new backend generator
func NewBackendGenerator(spec *ProjectSpec) *BackendGenerator {
	return &BackendGenerator{spec: spec}
}

// Generate creates all backend-related files
func (bg *BackendGenerator) Generate() (map[string]string, error) {
	files := make(map[string]string)
	
	// Language-specific generation
	switch bg.spec.GetBackendLanguage() {
	case "javascript":
		bg.generateNodeJSFiles(files)
	case "python":
		bg.generatePythonFiles(files)
	case "go":
		bg.generateGoFiles(files)
	default:
		bg.generateNodeJSFiles(files) // Default to Node.js
	}
	
	// Common files for all backends
	files["apps/backend/README.md"] = bg.generateBackendReadme()
	files["apps/backend/.env.example"] = bg.generateEnvExample()
	
	return files, nil
}

// generateNodeJSFiles creates Node.js/Express specific files
func (bg *BackendGenerator) generateNodeJSFiles(files map[string]string) {
	files["apps/backend/package.json"] = bg.generateNodePackageJSON()
	files["apps/backend/Dockerfile"] = bg.generateNodeDockerfile()
	files["apps/backend/src/server.js"] = bg.generateNodeServer()
	files["apps/backend/src/config/database.js"] = bg.generateNodeDatabaseConfig()
	files["apps/backend/src/config/auth.js"] = bg.generateNodeAuthConfig()
}

// generatePythonFiles creates Python/Django or Flask specific files
func (bg *BackendGenerator) generatePythonFiles(files map[string]string) {
	files["apps/backend/requirements.txt"] = bg.generatePythonRequirements()
	files["apps/backend/pyproject.toml"] = bg.generatePyprojectToml()
	files["apps/backend/Dockerfile"] = bg.generatePythonDockerfile()
	files["apps/backend/src/main.py"] = bg.generatePythonMain()
	files["apps/backend/.flake8"] = bg.generateFlake8Config()
}

// generateGoFiles creates Go specific files
func (bg *BackendGenerator) generateGoFiles(files map[string]string) {
	files["apps/backend/go.mod"] = bg.generateGoMod()
	files["apps/backend/Dockerfile"] = bg.generateGoDockerfile()
	files["apps/backend/src/main.go"] = bg.generateGoMain()
	// IMPLEMENT: Add .golangci.yml for Go linting configuration
}

// GetBackendLanguage extracts the programming language from backend choice
func (spec *ProjectSpec) GetBackendLanguage() string {
	switch {
	case strings.Contains(strings.ToLower(spec.Backend), "node"):
		return "javascript"
	case strings.Contains(strings.ToLower(spec.Backend), "python"):
		return "python"
	case strings.Contains(strings.ToLower(spec.Backend), "go"):
		return "go"
	default:
		return "javascript" // Default fallback
	}
}

// GetBackendFramework extracts the framework name from backend choice
func (spec *ProjectSpec) GetBackendFramework() string {
	backend := strings.ToLower(spec.Backend)
	switch {
	case strings.Contains(backend, "express"):
		return "express"
	case strings.Contains(backend, "nestjs"):
		return "nestjs"
	case strings.Contains(backend, "django"):
		return "django"
	case strings.Contains(backend, "flask"):
		return "flask"
	case strings.Contains(backend, "gin"):
		return "gin"
	case strings.Contains(backend, "standard"):
		return "stdlib"
	default:
		return "express" // Default fallback
	}
}

// HasDatabase checks if a specific database is required
func (spec *ProjectSpec) HasDatabase(dbType string) bool {
	for _, db := range spec.Databases {
		if strings.Contains(strings.ToLower(db), strings.ToLower(dbType)) {
			return true
		}
	}
	return false
}

// GetUserRolesList returns user roles as a slice
func (spec *ProjectSpec) GetUserRolesList() []string {
	if spec.UserRoles == "" || strings.ToLower(spec.UserRoles) == "none" {
		return []string{}
	}
	
	roles := strings.Split(spec.UserRoles, ",")
	for i, role := range roles {
		roles[i] = strings.TrimSpace(role)
	}
	return roles
}