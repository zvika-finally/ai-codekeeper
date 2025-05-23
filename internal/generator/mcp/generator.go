package mcp

import (
	"encoding/json"
)

// MCPGenerator handles MCP server configuration generation
type MCPGenerator struct {
	spec *ProjectSpec
}

// ProjectSpec represents the project specification
type ProjectSpec struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CoreEntity  string   `json:"core_entity"`
	Backend     string   `json:"backend"`
	Databases   []string `json:"databases"`
	APIStyle    string   `json:"api_style"`
	UserRoles   string   `json:"user_roles"`
	Domain      string   `json:"domain"`
	ProjectPath string   `json:"project_path,omitempty"`
}

// MCPConfig represents MCP server configuration
type MCPConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// NewMCPGenerator creates a new MCP generator
func NewMCPGenerator(spec *ProjectSpec) *MCPGenerator {
	return &MCPGenerator{spec: spec}
}

// Generate creates MCP server configurations
func (mg *MCPGenerator) Generate() (map[string]string, error) {
	files := make(map[string]string)
	
	// Generate main MCP configuration for Cursor
	mg.generateCursorMCPConfig(files)
	
	// Generate individual MCP server configurations
	mg.generateGitMCPServer(files)
	mg.generateProjectMCPServer(files)
	
	// Generate domain-specific MCP servers
	mg.generateDomainMCPServers(files)
	
	// Generate MCP ecosystem documentation
	mg.generateMCPDocumentation(files)
	
	return files, nil
}

// generateCursorMCPConfig creates the main Cursor MCP configuration
func (mg *MCPGenerator) generateCursorMCPConfig(files map[string]string) {
	config := map[string]interface{}{
		"mcpServers": mg.getMCPServers(),
	}
	
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	files[".cursor/mcp.json"] = string(configJSON)
	
	// Also create VS Code settings for MCP
	files[".vscode/settings.json"] = mg.generateVSCodeMCPSettings()
}

// generateGitMCPServer creates a Git MCP server
func (mg *MCPGenerator) generateGitMCPServer(files map[string]string) {
	files["scripts/mcp/git-server.js"] = `#!/usr/bin/env node

// Git MCP Server for ` + mg.spec.Name + `
// Provides Git operations and project context to AI assistants

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const { CallToolRequestSchema, ListToolsRequestSchema } = require('@modelcontextprotocol/sdk/types.js');
const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

class GitMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: '` + mg.spec.Name + `-git-mcp-server',
        version: '0.1.0',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.setupHandlers();
  }

  setupHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          {
            name: 'git_status',
            description: 'Get git repository status',
            inputSchema: {
              type: 'object',
              properties: {},
            },
          },
          {
            name: 'git_log',
            description: 'Get git commit history',
            inputSchema: {
              type: 'object',
              properties: {
                limit: {
                  type: 'number',
                  description: 'Number of commits to show',
                  default: 10,
                },
              },
            },
          },
          {
            name: 'git_diff',
            description: 'Get git diff for staging area',
            inputSchema: {
              type: 'object',
              properties: {
                file: {
                  type: 'string',
                  description: 'Specific file to diff (optional)',
                },
              },
            },
          },
          {
            name: 'project_structure',
            description: 'Get current project structure and key files',
            inputSchema: {
              type: 'object',
              properties: {},
            },
          },
        ],
      };
    });

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case 'git_status':
            return await this.gitStatus();
          case 'git_log':
            return await this.gitLog(args.limit || 10);
          case 'git_diff':
            return await this.gitDiff(args.file);
          case 'project_structure':
            return await this.projectStructure();
          default:
            throw new Error(` + "`Unknown tool: ${name}`" + `);
        }
      } catch (error) {
        return {
          content: [
            {
              type: 'text',
              text: ` + "`Error: ${error.message}`" + `,
            },
          ],
        };
      }
    });
  }

  async gitStatus() {
    const status = execSync('git status --porcelain', { encoding: 'utf8' });
    const branch = execSync('git branch --show-current', { encoding: 'utf8' }).trim();
    
    return {
      content: [
        {
          type: 'text',
          text: ` + "`Current branch: ${branch}\\n\\nStatus:\\n${status || 'Working tree clean'}`" + `,
        },
      ],
    };
  }

  async gitLog(limit) {
    const log = execSync(` + "`git log --oneline -${limit}`" + `, { encoding: 'utf8' });
    
    return {
      content: [
        {
          type: 'text',
          text: ` + "`Recent commits (${limit}):\\n${log}`" + `,
        },
      ],
    };
  }

  async gitDiff(file) {
    const command = file ? ` + "`git diff ${file}`" + ` : 'git diff --staged';
    const diff = execSync(command, { encoding: 'utf8' });
    
    return {
      content: [
        {
          type: 'text',
          text: diff || 'No changes to show',
        },
      ],
    };
  }

  async projectStructure() {
    const structure = this.getProjectStructure('.');
    const packageJson = this.getPackageInfo();
    
    return {
      content: [
        {
          type: 'text',
          text: ` + "`Project: ${mg.spec.Name}\\nDomain: ${mg.spec.Domain}\\n\\n${packageJson}\\n\\nStructure:\\n${structure}`" + `,
        },
      ],
    };
  }

  getProjectStructure(dir, level = 0, maxLevel = 2) {
    if (level > maxLevel) return '';
    
    const items = fs.readdirSync(dir);
    let structure = '';
    
    for (const item of items) {
      if (item.startsWith('.') && item !== '.env.example') continue;
      
      const itemPath = path.join(dir, item);
      const stat = fs.statSync(itemPath);
      const indent = '  '.repeat(level);
      
      if (stat.isDirectory()) {
        structure += ` + "`${indent}${item}/\\n`" + `;
        if (['src', 'apps', 'packages', 'docs'].includes(item)) {
          structure += this.getProjectStructure(itemPath, level + 1, maxLevel);
        }
      } else if (this.isImportantFile(item)) {
        structure += ` + "`${indent}${item}\\n`" + `;
      }
    }
    
    return structure;
  }

  isImportantFile(filename) {
    const important = [
      'package.json', 'go.mod', 'requirements.txt', 'Dockerfile',
      'docker-compose.yml', 'README.md', '.env.example',
      'tsconfig.json', 'tailwind.config.js'
    ];
    return important.includes(filename);
  }

  getPackageInfo() {
    try {
      const pkg = JSON.parse(fs.readFileSync('package.json', 'utf8'));
      return ` + "`Package: ${pkg.name} v${pkg.version}\\nDescription: ${pkg.description || 'N/A'}`" + `;
    } catch {
      return 'No package.json found';
    }
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('` + mg.spec.Name + ` Git MCP Server running on stdio');
  }
}

if (require.main === module) {
  const server = new GitMCPServer();
  server.run().catch(console.error);
}

module.exports = { GitMCPServer };
`
}

// generateProjectMCPServer creates a project-specific MCP server
func (mg *MCPGenerator) generateProjectMCPServer(files map[string]string) {
	files["scripts/mcp/project-server.js"] = `#!/usr/bin/env node

// Project-specific MCP Server for ` + mg.spec.Name + `
// Provides domain context and project-specific operations

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const { CallToolRequestSchema, ListToolsRequestSchema } = require('@modelcontextprotocol/sdk/types.js');
const fs = require('fs');
const path = require('path');

class ProjectMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: '` + mg.spec.Name + `-project-mcp-server',
        version: '0.1.0',
      },
      {
        capabilities: {
          tools: {},
          resources: {},
        },
      }
    );

    this.projectInfo = {
      name: '` + mg.spec.Name + `',
      domain: '` + mg.spec.Domain + `',
      coreEntity: '` + mg.spec.CoreEntity + `',
      backend: '` + mg.spec.Backend + `',
      apiStyle: '` + mg.spec.APIStyle + `',
    };

    this.setupHandlers();
  }

  setupHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          {
            name: 'get_project_context',
            description: 'Get comprehensive project context and architecture',
            inputSchema: {
              type: 'object',
              properties: {},
            },
          },
          {
            name: 'get_domain_guidelines',
            description: 'Get domain-specific development guidelines',
            inputSchema: {
              type: 'object',
              properties: {},
            },
          },
          {
            name: 'get_api_documentation',
            description: 'Get API documentation and endpoints',
            inputSchema: {
              type: 'object',
              properties: {},
            },
          },
          {
            name: 'get_deployment_info',
            description: 'Get deployment and infrastructure information',
            inputSchema: {
              type: 'object',
              properties: {},
            },
          },
        ],
      };
    });

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name } = request.params;

      try {
        switch (name) {
          case 'get_project_context':
            return await this.getProjectContext();
          case 'get_domain_guidelines':
            return await this.getDomainGuidelines();
          case 'get_api_documentation':
            return await this.getAPIDocumentation();
          case 'get_deployment_info':
            return await this.getDeploymentInfo();
          default:
            throw new Error(` + "`Unknown tool: ${name}`" + `);
        }
      } catch (error) {
        return {
          content: [
            {
              type: 'text',
              text: ` + "`Error: ${error.message}`" + `,
            },
          ],
        };
      }
    });
  }

  async getProjectContext() {
    const context = ` + "`" + `
# ` + mg.spec.Name + ` Project Context

## Overview
- **Domain**: ` + mg.spec.Domain + `
- **Core Entity**: ` + mg.spec.CoreEntity + `
- **Backend**: ` + mg.spec.Backend + `
- **API Style**: ` + mg.spec.APIStyle + `

## Architecture
This project follows clean architecture principles with modular components:
- **apps/backend**: ` + mg.spec.Backend + ` API server
- **apps/frontend**: React TypeScript application
- **docs/**: Comprehensive documentation and guidelines
- **infra/**: Infrastructure as code configurations

## Domain-Specific Context (` + mg.spec.Domain + `)
` + mg.getDomainSpecificContext() + `

## Development Guidelines
- Follow the standards in docs/frontend/STANDARDS.md
- Refer to docs/API_DESIGN.md for API patterns
- Check docs/SECURITY.md for security requirements
- See docs/DEPLOYMENT.md for deployment procedures
` + "`" + `;

    return {
      content: [
        {
          type: 'text',
          text: context,
        },
      ],
    };
  }

  async getDomainGuidelines() {
    try {
      const guidelines = fs.readFileSync('docs/frontend/' + this.projectInfo.domain.toUpperCase() + '_PATTERNS.md', 'utf8');
      return {
        content: [
          {
            type: 'text',
            text: guidelines,
          },
        ],
      };
    } catch {
      return {
        content: [
          {
            type: 'text',
            text: 'Domain-specific guidelines not found. Check docs/frontend/ directory.',
          },
        ],
      };
    }
  }

  async getAPIDocumentation() {
    try {
      const apiDocs = fs.readFileSync('docs/API_DESIGN.md', 'utf8');
      return {
        content: [
          {
            type: 'text',
            text: apiDocs,
          },
        ],
      };
    } catch {
      return {
        content: [
          {
            type: 'text',
            text: 'API documentation not found. Check docs/API_DESIGN.md',
          },
        ],
      };
    }
  }

  async getDeploymentInfo() {
    try {
      const deployDocs = fs.readFileSync('docs/DEPLOYMENT.md', 'utf8');
      const infraDocs = fs.readFileSync('docs/infrastructure/DOCKER.md', 'utf8');
      
      return {
        content: [
          {
            type: 'text',
            text: ` + "`${deployDocs}\\n\\n---\\n\\n${infraDocs}`" + `,
          },
        ],
      };
    } catch {
      return {
        content: [
          {
            type: 'text',
            text: 'Deployment documentation not found. Check docs/DEPLOYMENT.md and docs/infrastructure/',
          },
        ],
      };
    }
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('` + mg.spec.Name + ` Project MCP Server running on stdio');
  }
}

if (require.main === module) {
  const server = new ProjectMCPServer();
  server.run().catch(console.error);
}

module.exports = { ProjectMCPServer };
`
}