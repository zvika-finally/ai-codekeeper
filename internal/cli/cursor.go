package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/cursor"
)

func NewCursorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cursor",
		Short: "Cursor IDE integration with AI guard rails",
		Long: `Integrates AI Development Framework with Cursor IDE:

- Configures MCP (Model Context Protocol) servers
- Sets up domain-specific AI rules  
- Enables guard rails in Cursor AI assistant
- Provides real-time code validation
- Integrates domain expertise into code completion

This ensures Cursor's AI follows your project's guard rails and domain expertise.`,
	}

	cmd.AddCommand(newCursorSetupCmd())
	cmd.AddCommand(newCursorValidateCmd())
	cmd.AddCommand(newCursorRulesCmd())

	return cmd
}

func newCursorSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Setup Cursor IDE integration",
		Long: `Sets up Cursor IDE with AI Development Framework integration:

- Creates .cursor/ configuration directory
- Configures MCP servers for guard rails
- Sets up domain-specific AI rules
- Generates .cursorrules file
- Configures AI assistant settings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Blue("🎯 Setting up Cursor IDE integration...")

			// Check if .codekeeper exists
			config, err := loadProjectConfig()
			if err != nil {
				return fmt.Errorf("no AI dev project found. Run 'codekeeper init' first")
			}

			// Generate Cursor configuration
			cursorConfig, err := cursor.GenerateCursorConfig(".", config.Domain, config.GuardRails)
			if err != nil {
				return fmt.Errorf("failed to generate Cursor config: %w", err)
			}

			// Save Cursor configuration
			if err := cursor.SaveCursorConfig(".", cursorConfig); err != nil {
				return fmt.Errorf("failed to save Cursor config: %w", err)
			}

			color.Green("✅ Cursor IDE integration configured!")
			fmt.Printf("\n📋 Next steps:\n")
			fmt.Printf("1. Restart Cursor IDE\n")
			fmt.Printf("2. Verify MCP servers are loaded in Cursor settings\n")
			fmt.Printf("3. Test AI assistant with guard rails: Cmd+K\n")
			fmt.Printf("4. Check .cursorrules is being applied\n")

			return nil
		},
	}
}

func newCursorValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate Cursor IDE configuration",
		Long:  "Validates that Cursor IDE is properly configured with guard rails",
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Blue("🔍 Validating Cursor IDE configuration...")

			// TODO: Implement validation logic
			fmt.Printf("Cursor validation not yet implemented\n")
			fmt.Printf("Would check:\n")
			fmt.Printf("- .cursor/config.json exists and is valid\n")
			fmt.Printf("- MCP servers are accessible\n")
			fmt.Printf("- .cursorrules file is present\n")
			fmt.Printf("- Guard rails are active\n")

			return nil
		},
	}
}

func newCursorRulesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: "Manage Cursor AI rules",
		Long:  "View and manage AI rules for Cursor IDE",
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Blue("📋 Cursor AI Rules:")

			// TODO: Load and display current rules
			fmt.Printf("Rules management not yet implemented\n")
			fmt.Printf("Would show:\n")
			fmt.Printf("- Active guard rails\n")
			fmt.Printf("- Domain-specific rules\n")
			fmt.Printf("- Custom project rules\n")
			fmt.Printf("- Rule enforcement status\n")

			return nil
		},
	}
}

type ProjectConfig struct {
	Domain     string   `json:"domain"`
	GuardRails []string `json:"guard_rails"`
}

func loadProjectConfig() (*ProjectConfig, error) {
	// TODO: Load actual configuration from .codekeeper/config.json
	// For now, return mock data
	return &ProjectConfig{
		Domain:     "fintech",
		GuardRails: []string{"decimal_arithmetic", "audit_trails", "input_validation"},
	}, nil
}