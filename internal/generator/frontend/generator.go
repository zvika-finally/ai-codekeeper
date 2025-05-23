package frontend

import (
	"strings"
)

// FrontendGenerator handles frontend code generation
type FrontendGenerator struct {
	spec *ProjectSpec
}

// ProjectSpec represents the project specification
type ProjectSpec struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CoreEntity  string   `json:"core_entity"`
	Backend     string   `json:"backend"`
	APIStyle    string   `json:"api_style"`
	UserRoles   string   `json:"user_roles"`
	Domain      string   `json:"domain"`
	ProjectPath string   `json:"project_path,omitempty"`
}

// NewFrontendGenerator creates a new frontend generator
func NewFrontendGenerator(spec *ProjectSpec) *FrontendGenerator {
	return &FrontendGenerator{spec: spec}
}

// Generate creates frontend development guidelines and standards
func (fg *FrontendGenerator) Generate() (map[string]string, error) {
	files := make(map[string]string)
	
	// Generate frontend development standards
	fg.generateFrontendStandards(files)
	
	// Generate minimal setup guide
	fg.generateSetupGuide(files)
	
	return files, nil
}

// generateFrontendStandards creates development standards and guidelines
func (fg *FrontendGenerator) generateFrontendStandards(files map[string]string) {
	files["docs/frontend/STANDARDS.md"] = fg.generateFrontendStandardsDoc()
	files["docs/frontend/COMPONENT_PATTERNS.md"] = fg.generateComponentPatternsDoc()
	files["docs/frontend/STATE_MANAGEMENT.md"] = fg.generateStateManagementDoc()
	files["docs/frontend/TESTING_GUIDELINES.md"] = fg.generateTestingGuidelinesDoc()
	files["docs/frontend/SECURITY_PRACTICES.md"] = fg.generateSecurityPracticesDoc()
	
	// Add domain-specific guidelines if applicable
	domainGuidelines := fg.getDomainSpecificGuidelines()
	for path, content := range domainGuidelines {
		files[path] = content
	}
}

// generateSetupGuide creates minimal setup instructions
func (fg *FrontendGenerator) generateSetupGuide(files map[string]string) {
	files["docs/frontend/SETUP.md"] = fg.generateSetupGuideDoc()
	files["docs/frontend/PROJECT_STRUCTURE.md"] = fg.generateProjectStructureDoc()
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

// getDomainSpecificGuidelines returns domain-specific development guidelines
func (fg *FrontendGenerator) getDomainSpecificGuidelines() map[string]string {
	guidelines := make(map[string]string)
	
	switch fg.spec.Domain {
	case "fintech":
		guidelines["docs/frontend/FINTECH_PATTERNS.md"] = fg.generateFintechPatternsDoc()
	case "healthcare":
		guidelines["docs/frontend/HEALTHCARE_PATTERNS.md"] = fg.generateHealthcarePatternsDoc()
	case "ecommerce":
		guidelines["docs/frontend/ECOMMERCE_PATTERNS.md"] = fg.generateEcommercePatternsDoc()
	}
	
	return guidelines
}

// Helper method to format component name
func (fg *FrontendGenerator) getComponentName() string {
	return strings.Title(strings.ToLower(fg.spec.CoreEntity))
}