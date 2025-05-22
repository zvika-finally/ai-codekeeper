package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tycoonlabs/ai-codekeeper/internal/cli"
)

var version = "dev" // Set during build

func main() {
	rootCmd := &cobra.Command{
		Use:   "codekeeper",
		Short: "Finally AI CodeKeeper - AI development with domain-specific guard rails",
		Long: `Finally AI CodeKeeper - The intelligent development framework that ensures
AI assistants (Cursor, Cline) follow domain-specific guard rails and best practices.

Features:
• Domain expertise integration (fintech, healthcare, e-commerce)
• Real-time guard rails enforcement for AI code generation
• Comprehensive MCP ecosystem (Git, GitHub, JIRA, AWS, Figma)
• Compliance and security validation built-in
• Production-ready application scaffolding

Finally, AI development you can trust.`,
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(cli.NewCreateCmd())
	rootCmd.AddCommand(cli.NewInitCmd())
	rootCmd.AddCommand(cli.NewFeatureCmd())
	rootCmd.AddCommand(cli.NewEnvCmd())
	rootCmd.AddCommand(cli.NewCheckCmd())
	rootCmd.AddCommand(cli.NewConfigCmd())
	rootCmd.AddCommand(cli.NewCursorCmd())
	rootCmd.AddCommand(cli.NewMCPServerCmd())
	rootCmd.AddCommand(cli.NewDomainServerCmd())
	rootCmd.AddCommand(cli.NewMCPEcosystemCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}