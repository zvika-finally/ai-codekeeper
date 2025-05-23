package generator

import (
	"os"
	"path/filepath"
)

// StructureGenerator handles creation of the base directory structure
type StructureGenerator struct {
	spec *ProjectSpec
}

// NewStructureGenerator creates a new structure generator
func NewStructureGenerator(spec *ProjectSpec) *StructureGenerator {
	return &StructureGenerator{spec: spec}
}

// Generate creates the complete monorepo directory structure
func (sg *StructureGenerator) Generate() error {
	return sg.generateMonorepoStructure()
}

// generateMonorepoStructure creates the base directory structure
func (sg *StructureGenerator) generateMonorepoStructure() error {
	dirs := []string{
		"apps",
		"apps/backend",
		"apps/backend/src",
		"apps/frontend",
		"apps/frontend/src",
		"packages",
		"packages/shared-types",
		"packages/shared-types/src",
		"infra",
		"infra/aws",
		"infra/aws/modules",
		"infra/aws/modules/vpc",
		"infra/aws/modules/ecs",
		"infra/aws/modules/rds",
		"infra/aws/modules/secrets",
		"infra/render",
		"docs",
		".github",
		".github/workflows",
		".codekeeper", // Framework configuration
		".codekeeper/prompts",
		".codekeeper/templates",
		".devcontainer",
		".vscode",
		"scripts",
		"scripts/mcp",
	}

	for _, dir := range dirs {
		path := filepath.Join(sg.spec.ProjectPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}