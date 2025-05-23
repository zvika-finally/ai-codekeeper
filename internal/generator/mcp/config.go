package mcp

import (
	"encoding/json"
)

// getMCPServers returns the MCP server configuration for Cursor
func (mg *MCPGenerator) getMCPServers() map[string]MCPConfig {
	servers := make(map[string]MCPConfig)

	// Core project servers
	servers["git"] = MCPConfig{
		Command: "node",
		Args:    []string{"./scripts/mcp/git-server.js"},
		Env: map[string]string{
			"PROJECT_ROOT": ".",
		},
	}

	servers["project"] = MCPConfig{
		Command: "node",
		Args:    []string{"./scripts/mcp/project-server.js"},
	}

	// Domain-specific servers
	switch mg.spec.Domain {
	case "fintech":
		servers["compliance"] = MCPConfig{
			Command: "node",
			Args:    []string{"./scripts/mcp/compliance-server.js"},
		}
	case "healthcare":
		servers["hipaa"] = MCPConfig{
			Command: "node",
			Args:    []string{"./scripts/mcp/hipaa-server.js"},
		}
	case "ecommerce":
		servers["ecommerce"] = MCPConfig{
			Command: "node",
			Args:    []string{"./scripts/mcp/ecommerce-server.js"},
		}
	}

	// External MCP servers (optional)
	servers["github"] = MCPConfig{
		Command: "npx",
		Args:    []string{"-y", "@modelcontextprotocol/server-github"},
		Env: map[string]string{
			"GITHUB_PERSONAL_ACCESS_TOKEN": "${GITHUB_TOKEN}",
		},
	}

	return servers
}

// generateVSCodeMCPSettings creates VS Code settings for MCP
func (mg *MCPGenerator) generateVSCodeMCPSettings() string {
	settings := map[string]interface{}{
		"mcp.servers": mg.getMCPServers(),
		"mcp.enabled": true,
	}

	settingsJSON, _ := json.MarshalIndent(settings, "", "  ")
	return string(settingsJSON)
}

// generateMCPDocumentation creates documentation for MCP integration
func (mg *MCPGenerator) generateMCPDocumentation(files map[string]string) {
	files["docs/MCP_INTEGRATION.md"] = mg.generateMCPIntegrationDoc()
	files["scripts/mcp/README.md"] = mg.generateMCPReadme()
	files["scripts/mcp/package.json"] = mg.generateMCPPackageJSON()
}

// generateMCPIntegrationDoc creates MCP integration documentation
func (mg *MCPGenerator) generateMCPIntegrationDoc() string {
	return `# MCP (Model Context Protocol) Integration

This project includes MCP servers that provide AI assistants with project context and domain-specific guidance.

## Available MCP Servers

### Core Servers
- **Git Server**: Provides Git operations and project structure context
- **Project Server**: Provides project-specific context and documentation

### Domain-Specific Servers (` + mg.spec.Domain + `)
` + mg.getDomainSpecificMCPDocs() + `

### External Servers
- **GitHub Server**: Provides GitHub API integration (requires token)

## Setup

1. **Install Dependencies**
   ` + "```bash" + `
   cd scripts/mcp
   npm install
   ` + "```" + `

2. **Configure Cursor IDE**
   - The MCP configuration is automatically added to ` + "`.cursor/mcp.json`" + `
   - Restart Cursor to load the MCP servers

3. **Configure VS Code**
   - MCP settings are in ` + "`.vscode/settings.json`" + `
   - Install the MCP extension for VS Code

## Environment Variables

Create a ` + "`.env`" + ` file in the project root:

` + "```" + `
# GitHub integration (optional)
GITHUB_TOKEN=your_github_personal_access_token

# Domain-specific tokens
` + mg.getDomainSpecificEnvVars() + `
` + "```" + `

## Usage

Once configured, AI assistants can:
- Get project context and architecture information
- Access Git history and status
- Receive domain-specific guidance and compliance checks
- Get real-time documentation and API information

## Custom MCP Servers

You can extend the MCP ecosystem by creating custom servers in ` + "`scripts/mcp/`" + `.

### Creating a Custom Server

` + "```javascript" + `
const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');

class CustomMCPServer {
  constructor() {
    this.server = new Server({
      name: 'custom-server',
      version: '0.1.0',
    }, {
      capabilities: { tools: {} },
    });
    
    this.setupHandlers();
  }
  
  setupHandlers() {
    // Implement your handlers here
  }
  
  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
  }
}
` + "```" + `

## Troubleshooting

- Ensure Node.js is installed and accessible
- Check that all required dependencies are installed
- Verify environment variables are set correctly
- Check MCP server logs for errors
`
}

// generateMCPReadme creates README for MCP scripts directory
func (mg *MCPGenerator) generateMCPReadme() string {
	return `# MCP Servers

This directory contains Model Context Protocol (MCP) servers that provide AI assistants with project-specific context and capabilities.

## Available Servers

| Server | Purpose | Domain |
|--------|---------|---------|
| git-server.js | Git operations and project structure | All |
| project-server.js | Project context and documentation | All |
` + mg.getDomainSpecificMCPTable() + `

## Installation

` + "```bash" + `
npm install
` + "```" + `

## Running Servers

Each server can be run independently:

` + "```bash" + `
# Git server
node git-server.js

# Project server  
node project-server.js

# Domain-specific servers
` + mg.getDomainSpecificRunCommands() + `
` + "```" + `

## Development

To create a new MCP server:

1. Create a new JavaScript file in this directory
2. Implement the MCP server interface
3. Add configuration to ` + "`.cursor/mcp.json`" + `
4. Test the server with your AI assistant

## Dependencies

- @modelcontextprotocol/sdk: Core MCP SDK
- Domain-specific dependencies as needed
`
}

// generateMCPPackageJSON creates package.json for MCP servers
func (mg *MCPGenerator) generateMCPPackageJSON() string {
	packageJSON := map[string]interface{}{
		"name":        mg.spec.Name + "-mcp-servers",
		"version":     "0.1.0",
		"description": "MCP servers for " + mg.spec.Name + " project",
		"main":        "git-server.js",
		"scripts": map[string]string{
			"start:git":     "node git-server.js",
			"start:project": "node project-server.js",
		},
		"dependencies": map[string]string{
			"@modelcontextprotocol/sdk": "^0.4.0",
		},
		"keywords": []string{"mcp", "ai", "assistant", mg.spec.Domain},
		"author":   "AI CodeKeeper",
		"license":  "MIT",
	}

	// Add domain-specific scripts
	scripts := packageJSON["scripts"].(map[string]string)
	switch mg.spec.Domain {
	case "fintech":
		scripts["start:compliance"] = "node compliance-server.js"
	case "healthcare":
		scripts["start:hipaa"] = "node hipaa-server.js"
	case "ecommerce":
		scripts["start:ecommerce"] = "node ecommerce-server.js"
	}

	jsonBytes, _ := json.MarshalIndent(packageJSON, "", "  ")
	return string(jsonBytes)
}

// getDomainSpecificContext returns context for the project MCP server
func (mg *MCPGenerator) getDomainSpecificContext() string {
	switch mg.spec.Domain {
	case "fintech":
		return `This is a financial technology application focusing on secure transaction processing, compliance with financial regulations (PCI-DSS, SOX, GDPR), and audit trail management. Key areas include payment processing, account management, and regulatory compliance.`
	case "healthcare":
		return `This is a healthcare application requiring HIPAA compliance, secure handling of Protected Health Information (PHI), and integration with healthcare standards like FHIR. Key areas include patient data management, consent handling, and clinical workflows.`
	case "ecommerce":
		return `This is an e-commerce application focusing on product catalog management, secure payment processing, inventory management, and customer experience optimization. Key areas include shopping cart functionality, order processing, and payment integration.`
	default:
		return `This is a general web application following modern development practices and clean architecture principles.`
	}
}

// Helper methods for domain-specific documentation
func (mg *MCPGenerator) getDomainSpecificMCPDocs() string {
	switch mg.spec.Domain {
	case "fintech":
		return `- **Compliance Server**: Provides PCI-DSS, SOX, and financial regulation compliance checks`
	case "healthcare":
		return `- **HIPAA Server**: Provides HIPAA compliance validation and PHI handling guidance`
	case "ecommerce":
		return `- **E-commerce Server**: Provides e-commerce patterns and payment flow validation`
	default:
		return `- No domain-specific servers for generic applications`
	}
}

func (mg *MCPGenerator) getDomainSpecificEnvVars() string {
	switch mg.spec.Domain {
	case "fintech":
		return `STRIPE_SECRET_KEY=your_stripe_secret_key
PLAID_CLIENT_ID=your_plaid_client_id
PLAID_SECRET=your_plaid_secret`
	case "healthcare":
		return `FHIR_BASE_URL=your_fhir_server_url
FHIR_AUTH_TOKEN=your_fhir_auth_token`
	case "ecommerce":
		return `SHOPIFY_ACCESS_TOKEN=your_shopify_token
STRIPE_SECRET_KEY=your_stripe_secret_key`
	default:
		return `# No domain-specific environment variables needed`
	}
}

func (mg *MCPGenerator) getDomainSpecificMCPTable() string {
	switch mg.spec.Domain {
	case "fintech":
		return `| compliance-server.js | Financial compliance and audit checks | Fintech |`
	case "healthcare":
		return `| hipaa-server.js | HIPAA compliance and PHI validation | Healthcare |`
	case "ecommerce":
		return `| ecommerce-server.js | E-commerce patterns and validation | E-commerce |`
	default:
		return ``
	}
}

func (mg *MCPGenerator) getDomainSpecificRunCommands() string {
	switch mg.spec.Domain {
	case "fintech":
		return `node compliance-server.js`
	case "healthcare":
		return `node hipaa-server.js`
	case "ecommerce":
		return `node ecommerce-server.js`
	default:
		return `# No domain-specific servers`
	}
}