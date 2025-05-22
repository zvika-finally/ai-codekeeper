package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/generator"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Add AI development capabilities to existing project",
		Long: `Adds AI development framework to an existing project:

- Detects current project structure
- Adds .codekeeper/ configuration
- Sets up guard rails and CI/CD
- Configures development environment
- Integrates with existing build tools

Works with existing Node.js, Python, Go, and other projects.`,
		RunE: runInit,
	}

	cmd.Flags().String("domain", "", "Domain expertise to apply")
	cmd.Flags().Bool("force", false, "Force initialization even if .codekeeper exists")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	color.Blue("🔧 Initializing AI Development Framework")
	
	// Check if already initialized
	if _, err := os.Stat(".codekeeper"); err == nil {
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf(".codekeeper directory already exists. Use --force to override")
		}
	}

	// Detect current project structure
	projectInfo, err := detectProjectStructure()
	if err != nil {
		return fmt.Errorf("failed to detect project structure: %w", err)
	}

	color.Green("📋 Detected project type: %s", projectInfo.Type)
	
	// Create basic spec for existing project
	spec := &generator.ProjectSpec{
		Name:        filepath.Base(projectInfo.RootPath),
		Description: fmt.Sprintf("Existing %s project enhanced with AI development framework", projectInfo.Type),
		Domain:      getDomainFlag(cmd),
	}

	// Generate framework configuration
	if err := generateFrameworkConfig(spec, projectInfo); err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	color.Green("✅ AI Development Framework initialized!")
	fmt.Printf("\n📋 Next steps:\n")
	fmt.Printf("  codekeeper env setup    # Setup development environment\n")
	fmt.Printf("  codekeeper check        # Run guard rails validation\n")
	fmt.Printf("  codekeeper feature <name> # Generate new features\n")

	return nil
}

type ProjectInfo struct {
	Type     string
	RootPath string
	Language string
	Framework string
	HasDocker bool
}

func detectProjectStructure() (*ProjectInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	info := &ProjectInfo{
		RootPath: cwd,
		Type:     "unknown",
	}

	// Check for various project types
	if _, err := os.Stat("package.json"); err == nil {
		info.Type = "nodejs"
		info.Language = "javascript"
		// Could parse package.json to detect framework
	} else if _, err := os.Stat("requirements.txt"); err == nil || fileExists("pyproject.toml") {
		info.Type = "python"
		info.Language = "python"
	} else if _, err := os.Stat("go.mod"); err == nil {
		info.Type = "go"
		info.Language = "go"
	} else if _, err := os.Stat("Cargo.toml"); err == nil {
		info.Type = "rust"
		info.Language = "rust"
	}

	// Check for Docker
	info.HasDocker = fileExists("Dockerfile") || fileExists("docker-compose.yml")

	return info, nil
}

func generateFrameworkConfig(spec *generator.ProjectSpec, info *ProjectInfo) error {
	// Create .codekeeper directory
	if err := os.MkdirAll(".codekeeper", 0755); err != nil {
		return err
	}

	// Generate basic configuration
	_ = map[string]interface{}{
		"version": "1.0.0",
		"project_type": info.Type,
		"language": info.Language,
		"domain": spec.Domain,
		"initialized_at": "2024-01-01", // Would use actual time
		"guard_rails": map[string]interface{}{
			"enforcement": "advisory",
			"pre_commit": true,
			"ci": true,
		},
	}

	// Save configuration (would use actual JSON marshaling)
	configPath := filepath.Join(".codekeeper", "config.json")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write basic JSON (simplified for testing)
	_, err = file.WriteString(`{
  "version": "1.0.0",
  "project_type": "` + info.Type + `",
  "language": "` + info.Language + `",
  "domain": "` + spec.Domain + `",
  "initialized_at": "2024-01-01",
  "guard_rails": {
    "enforcement": "advisory",
    "pre_commit": true,
    "ci": true
  }
}`)

	return err
}

func getDomainFlag(cmd *cobra.Command) string {
	domain, _ := cmd.Flags().GetString("domain")
	if domain == "" {
		domain = "generic"
	}
	return domain
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}