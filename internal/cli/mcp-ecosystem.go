package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/cursor"
)

// NewMCPEcosystemCmd creates the mcp-ecosystem command
func NewMCPEcosystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-ecosystem",
		Short: "Setup complete MCP ecosystem for development",
		Long: `Configures a comprehensive MCP ecosystem including:

Development Tools:
- Git integration for version control
- GitHub/GitLab for repository management
- JIRA for project management
- Figma for design collaboration

Cloud & Infrastructure:
- AWS for cloud resources
- Terraform for infrastructure
- Docker for containerization
- Kubernetes for orchestration

Monitoring & Observability:
- DataDog/New Relic for monitoring
- Sentry for error tracking
- LogRocket for user sessions

This creates a complete development environment where AI assistants
can interact with your entire tech stack while following guard rails.`,
		RunE: setupMCPEcosystem,
	}

	cmd.Flags().StringP("domain", "d", "generic", "Domain expertise to apply")
	cmd.Flags().StringSliceP("services", "s", []string{"git", "github", "aws"}, "MCP services to enable")
	cmd.Flags().BoolP("all", "a", false, "Enable all available MCP services")
	cmd.Flags().StringP("output", "o", ".", "Output directory for configuration")

	return cmd
}

func setupMCPEcosystem(cmd *cobra.Command, args []string) error {
	domain, _ := cmd.Flags().GetString("domain")
	services, _ := cmd.Flags().GetStringSlice("services")
	enableAll, _ := cmd.Flags().GetBool("all")
	outputDir, _ := cmd.Flags().GetString("output")

	if enableAll {
		services = []string{
			"git", "github", "gitlab", "jira", "figma", 
			"aws", "terraform", "docker", "kubernetes",
			"datadog", "sentry", "logrocket",
		}
	}

	fmt.Printf("🌐 Setting up MCP ecosystem for %s domain...\n", domain)
	fmt.Printf("📦 Services: %v\n", services)

	// Generate comprehensive MCP configuration
	if err := cursor.GenerateComprehensiveMCPConfig(outputDir, domain, services); err != nil {
		return fmt.Errorf("failed to generate MCP ecosystem: %w", err)
	}

	// Generate Cline integration
	if err := cursor.GenerateClineIntegration(outputDir, domain, services); err != nil {
		return fmt.Errorf("failed to generate Cline integration: %w", err)
	}

	fmt.Println("\n✅ MCP ecosystem configured successfully!")
	fmt.Println("\n📋 Next steps:")
	fmt.Println("1. Configure API keys for enabled services")
	fmt.Println("2. Restart Cursor IDE and Cline")
	fmt.Println("3. Verify MCP servers in settings")
	fmt.Println("4. Test AI assistant with ecosystem integration")

	return nil
}