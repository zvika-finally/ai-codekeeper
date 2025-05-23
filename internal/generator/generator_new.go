package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zvika-finally/ai-codekeeper/internal/generator/backend"
	"github.com/zvika-finally/ai-codekeeper/internal/generator/frontend"
	"github.com/zvika-finally/ai-codekeeper/internal/generator/tooling"
)

// NewModularGenerator creates a streamlined generator using modular components
type NewModularGenerator struct {
	spec *ProjectSpec
}

// NewModular creates a new modular generator instance
func NewModular(spec *ProjectSpec) *NewModularGenerator {
	return &NewModularGenerator{spec: spec}
}

// Generate creates the complete project using modular components
func (g *NewModularGenerator) Generate() error {
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

	// Generate components using modular generators
	components := []ComponentGenerator{
		NewStructureGenerator(g.spec),
		NewConfigGenerator(g.spec),
	}

	for _, component := range components {
		if err := component.Generate(); err != nil {
			return fmt.Errorf("failed to generate component: %w", err)
		}
	}

	// Generate file-based components
	if err := g.generateFileComponents(); err != nil {
		return fmt.Errorf("failed to generate file components: %w", err)
	}

	// Initialize git repository last
	gitGen := tooling.NewGitGenerator(&tooling.ProjectSpec{
		Name:        g.spec.Name,
		Description: g.spec.Description,
		Domain:      g.spec.Domain,
		CoreEntity:  g.spec.CoreEntity,
		ProjectPath: g.spec.ProjectPath,
	})
	
	if err := gitGen.Generate(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return nil
}

// generateFileComponents generates all file-based components
func (g *NewModularGenerator) generateFileComponents() error {
	allFiles := make(map[string]string)

	// Backend files
	backendGen := backend.NewBackendGenerator(&backend.ProjectSpec{
		Name:        g.spec.Name,
		Description: g.spec.Description,
		CoreEntity:  g.spec.CoreEntity,
		Backend:     g.spec.Backend,
		Databases:   g.spec.Databases,
		APIStyle:    g.spec.APIStyle,
		UserRoles:   g.spec.UserRoles,
		Domain:      g.spec.Domain,
		ProjectPath: g.spec.ProjectPath,
	})

	backendFiles, err := backendGen.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate backend files: %w", err)
	}

	// Merge backend files
	for path, content := range backendFiles {
		allFiles[path] = content
	}

	// Frontend files
	frontendGen := frontend.NewFrontendGenerator(&frontend.ProjectSpec{
		Name:        g.spec.Name,
		Description: g.spec.Description,
		CoreEntity:  g.spec.CoreEntity,
		Backend:     g.spec.Backend,
		APIStyle:    g.spec.APIStyle,
		UserRoles:   g.spec.UserRoles,
		Domain:      g.spec.Domain,
		ProjectPath: g.spec.ProjectPath,
	})

	frontendFiles, err := frontendGen.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate frontend files: %w", err)
	}

	// Merge frontend files
	for path, content := range frontendFiles {
		allFiles[path] = content
	}

	// Generate documentation guidelines
	docsFiles, err := g.generateDocsGuidelines()
	if err != nil {
		return fmt.Errorf("failed to generate docs guidelines: %w", err)
	}
	for path, content := range docsFiles {
		allFiles[path] = content
	}

	// Generate infrastructure guidelines  
	infraFiles, err := g.generateInfrastructureGuidelines()
	if err != nil {
		return fmt.Errorf("failed to generate infrastructure guidelines: %w", err)
	}
	for path, content := range infraFiles {
		allFiles[path] = content
	}

	// Write all files
	for filePath, content := range allFiles {
		fullPath := filepath.Join(g.spec.ProjectPath, filePath)
		
		// Ensure directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}

	return nil
}

// generateDocsGuidelines creates documentation standards and guidelines
func (g *NewModularGenerator) generateDocsGuidelines() (map[string]string, error) {
	files := make(map[string]string)
	
	files["docs/README.md"] = g.generateProjectReadme()
	files["docs/API_DESIGN.md"] = g.generateAPIDesignDoc()
	files["docs/DEPLOYMENT.md"] = g.generateDeploymentDoc()
	files["docs/CONTRIBUTING.md"] = g.generateContributingDoc()
	files["docs/ARCHITECTURE.md"] = g.generateArchitectureDoc()
	
	return files, nil
}

// generateInfrastructureGuidelines creates infrastructure standards and guidelines
func (g *NewModularGenerator) generateInfrastructureGuidelines() (map[string]string, error) {
	files := make(map[string]string)
	
	files["docs/infrastructure/DOCKER.md"] = g.generateDockerGuideDoc()
	files["docs/infrastructure/CI_CD.md"] = g.generateCICDGuideDoc()
	files["docs/infrastructure/MONITORING.md"] = g.generateMonitoringDoc()
	files["docs/infrastructure/SECURITY.md"] = g.generateSecurityDoc()
	files["docs/infrastructure/SCALING.md"] = g.generateScalingDoc()
	
	return files, nil
}

// ConfigGenerator handles framework configuration generation
type ConfigGenerator struct {
	spec *ProjectSpec
}

// NewConfigGenerator creates a new config generator
func NewConfigGenerator(spec *ProjectSpec) *ConfigGenerator {
	return &ConfigGenerator{spec: spec}
}

// Generate creates framework configuration files
func (cg *ConfigGenerator) Generate() error {
	// Create .codekeeper configuration
	config := &FrameworkConfig{
		Version: "2.0.0",
		Domain: DomainConfig{
			Name:    cg.spec.Domain,
			Version: "1.0.0",
		},
		GuardRails: GuardRailsConfig{
			Enforcement: "advisory",
			PreCommit:   true,
			CI:          true,
			IDE:         "cursor",
		},
		DevEnvironment: DevEnvironmentConfig{
			Type:     "docker-compose",
			Services: []string{"postgres", "redis"},
			Ports: map[string]int{
				"backend":  8080,
				"frontend": 3000,
			},
		},
	}

	configPath := filepath.Join(cg.spec.ProjectPath, ".codekeeper", "config.json")
	return saveConfigJSON(configPath, config)
}

// saveConfigJSON saves data as JSON file
func saveConfigJSON(filename string, data interface{}) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}