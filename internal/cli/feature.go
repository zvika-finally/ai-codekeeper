package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewFeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature [name]",
		Short: "Generate a new feature with domain expertise",
		Long: `Generates a new feature following domain best practices:

- Domain-specific code patterns
- Guard rails enforcement
- Complete CRUD operations
- Tests and documentation
- Integration with existing project structure

Examples:
  codekeeper feature payments --domain fintech
  codekeeper feature user-management --domain generic
  codekeeper feature inventory --domain ecommerce`,
		Args: cobra.ExactArgs(1),
		RunE: runFeature,
	}

	cmd.Flags().String("domain", "", "Domain expertise to apply")
	cmd.Flags().String("type", "", "Feature type (crud, api, service)")

	return cmd
}

func runFeature(cmd *cobra.Command, args []string) error {
	featureName := args[0]
	domain, _ := cmd.Flags().GetString("domain")
	
	color.Blue("🎯 Generating feature: %s", featureName)
	
	if domain != "" {
		color.Yellow("📚 Applying %s domain expertise", domain)
	}

	// TODO: Implement actual feature generation
	fmt.Printf("Feature generation not yet implemented\n")
	fmt.Printf("Would generate feature '%s' with domain '%s'\n", featureName, domain)

	return nil
}