package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/cursor"
)

// NewMCPServerCmd creates the mcp-server command
func NewMCPServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "Start the MCP server for guard rails validation",
		Long: `Starts the Model Context Protocol (MCP) server that provides
guard rails validation and domain expertise to Cursor IDE.

This server runs in the background and communicates with Cursor
through the MCP protocol to provide real-time code validation.`,
		RunE: runMCPServer,
	}

	cmd.Flags().StringP("domain", "d", "generic", "Domain expertise to provide")
	cmd.Flags().StringP("project", "p", ".", "Project path")

	return cmd
}

// NewDomainServerCmd creates the domain-server command
func NewDomainServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain-server",
		Short: "Start the MCP server for domain expertise",
		Long: `Starts the Model Context Protocol (MCP) server that provides
domain-specific templates and recommendations to Cursor IDE.

This server provides code templates, technology recommendations,
and domain-specific patterns.`,
		RunE: runDomainServer,
	}

	cmd.Flags().StringP("domain", "d", "generic", "Domain expertise to provide")

	return cmd
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	domain, _ := cmd.Flags().GetString("domain")
	projectPath, _ := cmd.Flags().GetString("project")

	// Set environment variables for the server
	os.Setenv("CODEKEEPER_DOMAIN", domain)
	os.Setenv("CODEKEEPER_PROJECT_PATH", projectPath)

	fmt.Fprintf(os.Stderr, "Starting MCP Guard Rails Server for domain: %s\n", domain)
	
	// Start the MCP server
	return cursor.StartGuardRailsServer(domain, projectPath)
}

func runDomainServer(cmd *cobra.Command, args []string) error {
	domain, _ := cmd.Flags().GetString("domain")

	// Set environment variables for the server
	os.Setenv("CODEKEEPER_DOMAIN", domain)

	fmt.Fprintf(os.Stderr, "Starting MCP Domain Server for: %s\n", domain)
	
	// Start the domain expert server
	return cursor.StartDomainServer(domain)
}