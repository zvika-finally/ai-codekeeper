package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tycoonlabs/ai-codekeeper/internal/cursor"
)

// generateCursorIntegration adds Cursor IDE integration to generated projects
func (g *Generator) generateCursorIntegration() error {
	// Generate Cursor configuration with MCP integration
	cursorConfig, err := cursor.GenerateCursorConfig(g.spec.ProjectPath, g.spec.Domain, g.getDomainGuardRails())
	if err != nil {
		return err
	}

	// Save Cursor configuration
	if err := cursor.SaveCursorConfig(g.spec.ProjectPath, cursorConfig); err != nil {
		return err
	}

	// Add MCP server scripts
	if err := g.generateMCPServers(); err != nil {
		return err
	}

	return nil
}

// generateMCPServers creates MCP server scripts for guard rails integration
func (g *Generator) generateMCPServers() error {
	
	files := map[string]string{
		"scripts/mcp/guardrails-server.js": g.generateGuardRailsServer(),
		"scripts/mcp/domain-server.js": g.generateDomainServer(),
		"scripts/mcp/package.json": g.generateMCPPackageJson(),
	}

	for filePath, content := range files {
		fullPath := filepath.Join(g.spec.ProjectPath, filePath)
		
		// Ensure directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		
		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateGuardRailsServer() string {
	return `#!/usr/bin/env node
/**
 * MCP Server for AI Development Framework Guard Rails
 * Provides real-time validation and domain expertise to Cursor IDE
 */

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');

class GuardRailsServer {
  constructor() {
    this.server = new Server(
      {
        name: 'codekeeper-guardrails',
        version: '1.0.0',
      },
      {
        capabilities: {
          tools: {},
          resources: {},
        },
      }
    );

    this.setupToolHandlers();
    this.setupResourceHandlers();
  }

  setupToolHandlers() {
    // Validate code against guard rails
    this.server.setRequestHandler('tools/call', async (request) => {
      if (request.params.name === 'validate_code') {
        return this.validateCode(request.params.arguments);
      }
      
      if (request.params.name === 'get_domain_patterns') {
        return this.getDomainPatterns(request.params.arguments);
      }

      throw new Error(` + "`" + `Unknown tool: ${request.params.name}` + "`" + `);
    });

    // List available tools
    this.server.setRequestHandler('tools/list', async () => {
      return {
        tools: [
          {
            name: 'validate_code',
            description: 'Validate code against domain-specific guard rails',
            inputSchema: {
              type: 'object',
              properties: {
                code: { type: 'string' },
                language: { type: 'string' },
                domain: { type: 'string' }
              },
              required: ['code']
            }
          },
          {
            name: 'get_domain_patterns',
            description: 'Get domain-specific code patterns and templates',
            inputSchema: {
              type: 'object',
              properties: {
                domain: { type: 'string' },
                pattern_type: { type: 'string' }
              },
              required: ['domain']
            }
          }
        ]
      };
    });
  }

  setupResourceHandlers() {
    // Provide guard rails documentation
    this.server.setRequestHandler('resources/list', async () => {
      return {
        resources: [
          {
            uri: 'guardrails://` + g.spec.Domain + `/rules',
            name: '` + g.spec.Domain + ` Guard Rails',
            description: 'Domain-specific development rules and patterns'
          }
        ]
      };
    });

    this.server.setRequestHandler('resources/read', async (request) => {
      if (request.params.uri === 'guardrails://` + g.spec.Domain + `/rules') {
        return {
          contents: [
            {
              uri: request.params.uri,
              mimeType: 'text/markdown',
              text: this.getGuardRailsDocumentation()
            }
          ]
        };
      }
      
      throw new Error(` + "`" + `Unknown resource: ${request.params.uri}` + "`" + `);
    });
  }

  async validateCode(args) {
    const { code, language = 'typescript', domain = '` + g.spec.Domain + `' } = args;
    
    // Implement domain-specific validation
    const violations = this.checkGuardRails(code, domain);
    
    return {
      content: [
        {
          type: 'text',
          text: JSON.stringify({
            valid: violations.length === 0,
            violations: violations,
            suggestions: this.getSuggestions(violations)
          }, null, 2)
        }
      ]
    };
  }

  async getDomainPatterns(args) {
    const { domain = '` + g.spec.Domain + `', pattern_type = 'all' } = args;
    
    const patterns = this.getDomainSpecificPatterns(domain, pattern_type);
    
    return {
      content: [
        {
          type: 'text',
          text: JSON.stringify(patterns, null, 2)
        }
      ]
    };
  }

  checkGuardRails(code, domain) {
    const violations = [];
    
    // Domain-specific checks
    if (domain === 'fintech') {
      // Check for floating point arithmetic with money
      if (code.match(/(amount|price|balance|fee).*[+\-*/].*\d*\.\d/)) {
        violations.push({
          type: 'error',
          rule: 'decimal_arithmetic',
          message: 'Use Decimal type for monetary calculations, not floating point',
          line: this.findLineNumber(code, /(amount|price|balance|fee).*[+\-*/].*\d*\.\d/)
        });
      }

      // Check for missing audit logging
      if (code.match(/(transaction|payment|transfer)/) && !code.match(/audit|log/)) {
        violations.push({
          type: 'warning',
          rule: 'audit_trails',
          message: 'Financial operations should include audit logging',
          line: this.findLineNumber(code, /(transaction|payment|transfer)/)
        });
      }
    }

    // General security checks
    if (code.match(/password.*=.*["'].*["']/)) {
      violations.push({
        type: 'error',
        rule: 'no_hardcoded_secrets',
        message: 'Never hardcode passwords or secrets',
        line: this.findLineNumber(code, /password.*=.*["'].*["']/)
      });
    }

    return violations;
  }

  getDomainSpecificPatterns(domain, patternType) {
    const patterns = {
      fintech: {
        monetary_calculation: {
          good: ` + "`" + `const total = new Decimal(amount).plus(fee);` + "`" + `,
          bad: ` + "`" + `const total = amount + fee; // Never use + for money` + "`" + `,
          explanation: 'Always use Decimal for monetary calculations'
        },
        audit_logging: {
          good: ` + "`" + `await auditLog.record('PAYMENT_CREATED', { transactionId, amount, userId });` + "`" + `,
          explanation: 'Log all financial operations for compliance'
        }
      }
    };

    return patterns[domain] || {};
  }

  findLineNumber(code, pattern) {
    const lines = code.split('\n');
    for (let i = 0; i < lines.length; i++) {
      if (pattern.test(lines[i])) {
        return i + 1;
      }
    }
    return 1;
  }

  getSuggestions(violations) {
    return violations.map(v => {
      switch (v.rule) {
        case 'decimal_arithmetic':
          return 'Consider using the Decimal.js library for precise monetary calculations';
        case 'audit_trails':
          return 'Add audit logging using the project audit logger';
        default:
          return 'Follow the project guard rails and best practices';
      }
    });
  }

  getGuardRailsDocumentation() {
    return ` + "`" + `# ` + g.spec.Domain + ` Development Guard Rails

## Core Rules
` + strings.Join(g.getDomainGuardRails(), "\n- ") + `

## Code Patterns
See domain-specific patterns for examples of correct implementations.

## Enforcement
These rules are enforced in real-time through Cursor IDE integration.
` + "`" + `;
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
  }
}

// Start the server
if (require.main === module) {
  const server = new GuardRailsServer();
  server.run().catch(console.error);
}

module.exports = { GuardRailsServer };`
}

func (g *Generator) generateDomainServer() string {
	return `#!/usr/bin/env node
/**
 * MCP Server for Domain-Specific Expertise
 * Provides domain knowledge and templates to Cursor IDE
 */

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');

class DomainExpertServer {
  constructor(domain = '` + g.spec.Domain + `') {
    this.domain = domain;
    this.server = new Server(
      {
        name: ` + "`" + `codekeeper-domain-${domain}` + "`" + `,
        version: '1.0.0',
      },
      {
        capabilities: {
          tools: {},
          resources: {},
        },
      }
    );

    this.setupToolHandlers();
  }

  setupToolHandlers() {
    this.server.setRequestHandler('tools/call', async (request) => {
      if (request.params.name === 'get_template') {
        return this.getTemplate(request.params.arguments);
      }
      
      if (request.params.name === 'get_recommendation') {
        return this.getRecommendation(request.params.arguments);
      }

      throw new Error(` + "`" + `Unknown tool: ${request.params.name}` + "`" + `);
    });

    this.server.setRequestHandler('tools/list', async () => {
      return {
        tools: [
          {
            name: 'get_template',
            description: ` + "`" + `Get ${this.domain}-specific code templates` + "`" + `,
            inputSchema: {
              type: 'object',
              properties: {
                template_type: { type: 'string' },
                entity: { type: 'string' }
              },
              required: ['template_type']
            }
          },
          {
            name: 'get_recommendation',
            description: ` + "`" + `Get ${this.domain} technology recommendations` + "`" + `,
            inputSchema: {
              type: 'object',
              properties: {
                requirement: { type: 'string' },
                context: { type: 'string' }
              },
              required: ['requirement']
            }
          }
        ]
      };
    });
  }

  async getTemplate(args) {
    const { template_type, entity } = args;
    
    const template = this.getDomainTemplate(template_type, entity);
    
    return {
      content: [
        {
          type: 'text',
          text: template
        }
      ]
    };
  }

  async getRecommendation(args) {
    const { requirement, context } = args;
    
    const recommendation = this.getDomainRecommendation(requirement, context);
    
    return {
      content: [
        {
          type: 'text',
          text: JSON.stringify(recommendation, null, 2)
        }
      ]
    };
  }

  getDomainTemplate(templateType, entity) {
    if (this.domain === 'fintech') {
      switch (templateType) {
        case 'model':
          return this.getFinancialModelTemplate(entity);
        case 'service':
          return this.getFinancialServiceTemplate(entity);
        case 'controller':
          return this.getFinancialControllerTemplate(entity);
        default:
          return 'Unknown template type';
      }
    }
    
    return 'Generic template';
  }

  getFinancialModelTemplate(entity) {
    return ` + "`" + `import { Decimal } from 'decimal.js';

export interface ${entity} {
  id: string;
  amount: Decimal;
  currency: string;
  status: '${entity.toLowerCase()}_status';
  createdAt: Date;
  updatedAt: Date;
  auditTrail: AuditEntry[];
  metadata: Record<string, any>;
}

export interface Create${entity}Request {
  amount: string; // String to avoid precision loss
  currency: string;
  description?: string;
  metadata?: Record<string, any>;
}

export type ${entity}Status = 'pending' | 'completed' | 'failed' | 'cancelled';

export interface AuditEntry {
  action: string;
  timestamp: Date;
  userId: string;
  changes: Record<string, any>;
}` + "`" + `;
  }

  getFinancialServiceTemplate(entity) {
    return ` + "`" + `import { Decimal } from 'decimal.js';
import { ${entity}, Create${entity}Request } from '../models/${entity}';
import { AuditLogger } from '../utils/audit-logger';

export class ${entity}Service {
  constructor(
    private auditLogger: AuditLogger
  ) {}

  async create${entity}(request: Create${entity}Request): Promise<${entity}> {
    // Validate amount
    const amount = new Decimal(request.amount);
    if (amount.lessThanOrEqualTo(0)) {
      throw new Error('Amount must be positive');
    }

    // Create ${entity.toLowerCase()}
    const ${entity.toLowerCase()} = await this.repository.create({
      ...request,
      amount,
      status: 'pending',
    });

    // Audit log
    await this.auditLogger.log({
      action: '${entity.toUpperCase()}_CREATED',
      entityId: ${entity.toLowerCase()}.id,
      userId: request.userId,
      data: { amount: amount.toString(), currency: request.currency }
    });

    return ${entity.toLowerCase()};
  }

  async process${entity}(id: string): Promise<${entity}> {
    const ${entity.toLowerCase()} = await this.repository.findById(id);
    if (!${entity.toLowerCase()}) {
      throw new Error('${entity} not found');
    }

    // Process with external service
    // Implement idempotency
    // Handle failures gracefully
    
    return ${entity.toLowerCase()};
  }
}` + "`" + `;
  }

  getDomainRecommendation(requirement, context) {
    if (this.domain === 'fintech') {
      return {
        technology: this.getFinancialTechRecommendation(requirement),
        security: this.getFinancialSecurityRecommendation(requirement),
        compliance: this.getComplianceRecommendation(requirement)
      };
    }
    
    return { recommendation: 'Generic recommendation' };
  }

  getFinancialTechRecommendation(requirement) {
    const recommendations = {
      'payment_processing': 'Use Stripe or Modern Treasury for payment processing',
      'database': 'PostgreSQL for ACID compliance, Redis for caching',
      'authentication': 'Auth0 or AWS Cognito with MFA',
      'monitoring': 'Datadog or New Relic for financial transaction monitoring'
    };
    
    return recommendations[requirement] || 'Consult fintech architecture guidelines';
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
  }
}

// Start the server
if (require.main === module) {
  const domain = process.argv[2] || '` + g.spec.Domain + `';
  const server = new DomainExpertServer(domain);
  server.run().catch(console.error);
}

module.exports = { DomainExpertServer };`
}

func (g *Generator) generateMCPPackageJson() string {
	return `{
  "name": "@` + g.spec.Name + `/mcp-servers",
  "version": "1.0.0",
  "description": "MCP servers for AI Development Framework",
  "main": "guardrails-server.js",
  "scripts": {
    "start:guardrails": "node guardrails-server.js",
    "start:domain": "node domain-server.js"
  },
  "dependencies": {
    "@modelcontextprotocol/sdk": "latest"
  }
}`
}