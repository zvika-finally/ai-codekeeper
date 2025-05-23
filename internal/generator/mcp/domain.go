package mcp

// Domain-specific MCP server generation

// generateDomainMCPServers creates domain-specific MCP servers
func (mg *MCPGenerator) generateDomainMCPServers(files map[string]string) {
	switch mg.spec.Domain {
	case "fintech":
		mg.generateFintechMCPServers(files)
	case "healthcare":
		mg.generateHealthcareMCPServers(files)
	case "ecommerce":
		mg.generateEcommerceMCPServers(files)
	}
}

// generateFintechMCPServers creates fintech-specific MCP servers
func (mg *MCPGenerator) generateFintechMCPServers(files map[string]string) {
	// Compliance MCP Server
	files["scripts/mcp/compliance-server.js"] = `#!/usr/bin/env node

// Fintech Compliance MCP Server
// Provides compliance checks and financial regulations guidance

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const { CallToolRequestSchema, ListToolsRequestSchema } = require('@modelcontextprotocol/sdk/types.js');

class ComplianceMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: 'fintech-compliance-mcp-server',
        version: '0.1.0',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.complianceRules = {
      'PCI-DSS': 'Payment Card Industry Data Security Standard',
      'SOX': 'Sarbanes-Oxley Act',
      'GDPR': 'General Data Protection Regulation',
      'KYC': 'Know Your Customer',
      'AML': 'Anti-Money Laundering',
      'FFIEC': 'Federal Financial Institutions Examination Council',
    };

    this.setupHandlers();
  }

  setupHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          {
            name: 'check_compliance',
            description: 'Check code or configuration for compliance issues',
            inputSchema: {
              type: 'object',
              properties: {
                type: {
                  type: 'string',
                  enum: ['code', 'config', 'data'],
                  description: 'Type of compliance check',
                },
                content: {
                  type: 'string',
                  description: 'Content to check',
                },
              },
              required: ['type', 'content'],
            },
          },
          {
            name: 'get_compliance_requirements',
            description: 'Get compliance requirements for specific regulation',
            inputSchema: {
              type: 'object',
              properties: {
                regulation: {
                  type: 'string',
                  enum: ['PCI-DSS', 'SOX', 'GDPR', 'KYC', 'AML', 'FFIEC'],
                  description: 'Compliance regulation',
                },
              },
              required: ['regulation'],
            },
          },
          {
            name: 'generate_audit_trail',
            description: 'Generate audit trail template for financial operations',
            inputSchema: {
              type: 'object',
              properties: {
                operation: {
                  type: 'string',
                  description: 'Financial operation type',
                },
              },
              required: ['operation'],
            },
          },
        ],
      };
    });

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case 'check_compliance':
            return await this.checkCompliance(args.type, args.content);
          case 'get_compliance_requirements':
            return await this.getComplianceRequirements(args.regulation);
          case 'generate_audit_trail':
            return await this.generateAuditTrail(args.operation);
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

  async checkCompliance(type, content) {
    const issues = [];

    if (type === 'code') {
      // Check for common fintech compliance issues
      if (content.includes('console.log') && content.includes('password')) {
        issues.push('❌ Sensitive data logging detected - PCI-DSS violation');
      }
      if (content.includes('Math.random()') && content.includes('transaction')) {
        issues.push('❌ Insecure random number generation for transactions');
      }
      if (!content.includes('decrypt') && content.includes('encrypt')) {
        issues.push('⚠️ Consider implementing proper encryption/decryption patterns');
      }
      if (content.includes('parseFloat') && (content.includes('amount') || content.includes('price'))) {
        issues.push('❌ Use Decimal/BigNumber for monetary calculations, not parseFloat');
      }
    }

    const result = issues.length > 0 ? 
      ` + "`Compliance Issues Found:\\n${issues.join('\\n')}`" + ` :
      '✅ No compliance issues detected';

    return {
      content: [
        {
          type: 'text',
          text: result,
        },
      ],
    };
  }

  async getComplianceRequirements(regulation) {
    const requirements = {
      'PCI-DSS': ` + "`" + `
# PCI DSS Requirements for Fintech

## Key Requirements:
1. **Secure Network**: Implement firewalls and secure configurations
2. **Protect Data**: Never store sensitive authentication data
3. **Vulnerability Management**: Maintain secure systems and applications
4. **Strong Access Controls**: Restrict access on business need-to-know
5. **Monitor Networks**: Regularly monitor and test networks
6. **Information Security Policy**: Maintain policy that addresses information security

## Implementation Guidelines:
- Use strong encryption for card data transmission
- Implement multi-factor authentication
- Regular security testing and vulnerability scans
- Secure coding practices for payment processing
- Audit logging for all card data access
      ` + "`" + `,
      'SOX': ` + "`" + `
# SOX Compliance for Financial Systems

## Key Requirements:
1. **Internal Controls**: Establish and maintain internal controls over financial reporting
2. **Documentation**: Document all financial processes and controls
3. **Testing**: Regular testing of internal controls effectiveness
4. **Segregation of Duties**: Separate authorization, recording, and custody functions
5. **IT General Controls**: Secure IT systems that support financial reporting

## Implementation Guidelines:
- Automated audit trails for all financial transactions
- Role-based access controls with approval workflows
- Change management processes for financial systems
- Data backup and recovery procedures
- Regular compliance testing and reporting
      ` + "`" + `,
      // Add other regulations...
    };

    return {
      content: [
        {
          type: 'text',
          text: requirements[regulation] || 'Regulation not found',
        },
      ],
    };
  }

  async generateAuditTrail(operation) {
    const template = ` + "`" + `
// Audit Trail Template for ${operation}

interface AuditTrail {
  id: string;
  timestamp: Date;
  operation: '${operation}';
  userId: string;
  userRole: string;
  entityId: string;
  beforeState?: any;
  afterState?: any;
  metadata: {
    ipAddress: string;
    userAgent: string;
    sessionId: string;
    requestId: string;
  };
  complianceFlags: {
    pciRelevant: boolean;
    soxRelevant: boolean;
    regulatoryReporting: boolean;
  };
}

// Implementation example:
async function create${operation}AuditTrail(context: AuditContext) {
  const auditEntry: AuditTrail = {
    id: generateUUID(),
    timestamp: new Date(),
    operation: '${operation}',
    userId: context.user.id,
    userRole: context.user.role,
    entityId: context.entity.id,
    beforeState: context.beforeState,
    afterState: context.afterState,
    metadata: {
      ipAddress: context.request.ip,
      userAgent: context.request.userAgent,
      sessionId: context.session.id,
      requestId: context.request.id,
    },
    complianceFlags: {
      pciRelevant: isPCIRelevant('${operation}'),
      soxRelevant: isSOXRelevant('${operation}'),
      regulatoryReporting: requiresReporting('${operation}'),
    },
  };

  await auditLogger.log(auditEntry);
  
  // For SOX compliance, also log to immutable storage
  if (auditEntry.complianceFlags.soxRelevant) {
    await immutableAuditStore.append(auditEntry);
  }
}
    ` + "`" + `;

    return {
      content: [
        {
          type: 'text',
          text: template,
        },
      ],
    };
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('Fintech Compliance MCP Server running on stdio');
  }
}

if (require.main === module) {
  const server = new ComplianceMCPServer();
  server.run().catch(console.error);
}

module.exports = { ComplianceMCPServer };
`
}

// generateHealthcareMCPServers creates healthcare-specific MCP servers
func (mg *MCPGenerator) generateHealthcareMCPServers(files map[string]string) {
	files["scripts/mcp/hipaa-server.js"] = `#!/usr/bin/env node

// Healthcare HIPAA Compliance MCP Server
// Provides HIPAA compliance checks and healthcare data handling guidance

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const { CallToolRequestSchema, ListToolsRequestSchema } = require('@modelcontextprotocol/sdk/types.js');

class HIPAAMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: 'healthcare-hipaa-mcp-server',
        version: '0.1.0',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.phiElements = [
      'name', 'email', 'phone', 'ssn', 'medical_record_number',
      'patient_id', 'birth_date', 'address', 'diagnosis'
    ];

    this.setupHandlers();
  }

  setupHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          {
            name: 'check_phi_exposure',
            description: 'Check code for potential PHI exposure',
            inputSchema: {
              type: 'object',
              properties: {
                code: {
                  type: 'string',
                  description: 'Code to analyze for PHI exposure',
                },
              },
              required: ['code'],
            },
          },
          {
            name: 'generate_consent_form',
            description: 'Generate HIPAA consent form template',
            inputSchema: {
              type: 'object',
              properties: {
                purpose: {
                  type: 'string',
                  description: 'Purpose of data collection',
                },
              },
              required: ['purpose'],
            },
          },
          {
            name: 'validate_access_controls',
            description: 'Validate access control implementation',
            inputSchema: {
              type: 'object',
              properties: {
                implementation: {
                  type: 'string',
                  description: 'Access control code to validate',
                },
              },
              required: ['implementation'],
            },
          },
        ],
      };
    });

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case 'check_phi_exposure':
            return await this.checkPHIExposure(args.code);
          case 'generate_consent_form':
            return await this.generateConsentForm(args.purpose);
          case 'validate_access_controls':
            return await this.validateAccessControls(args.implementation);
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

  async checkPHIExposure(code) {
    const issues = [];

    // Check for PHI in logs
    if (code.includes('console.log') || code.includes('logger.')) {
      for (const phi of this.phiElements) {
        if (code.includes(phi)) {
          issues.push(` + "`❌ Potential PHI logging: ${phi}`" + `);
        }
      }
    }

    // Check for unencrypted PHI storage
    if (code.includes('localStorage') || code.includes('sessionStorage')) {
      issues.push('❌ PHI should not be stored in browser storage');
    }

    // Check for proper encryption
    if (code.includes('patient') && !code.includes('encrypt')) {
      issues.push('⚠️ Consider encryption for patient data');
    }

    const result = issues.length > 0 ? 
      ` + "`HIPAA Compliance Issues:\\n${issues.join('\\n')}`" + ` :
      '✅ No HIPAA compliance issues detected';

    return {
      content: [
        {
          type: 'text',
          text: result,
        },
      ],
    };
  }

  async generateConsentForm(purpose) {
    const form = ` + "`" + `
# HIPAA Authorization Form

## Purpose
This authorization allows [HEALTHCARE_ENTITY] to use and disclose your protected health information (PHI) for: ${purpose}

## Information to be Used/Disclosed
- Medical records and treatment information
- Billing and insurance information
- ${purpose === 'research' ? 'Research data and outcomes' : 'Relevant health information'}

## Recipient(s)
The information will be disclosed to:
- Authorized healthcare providers
- ${purpose === 'treatment' ? 'Consulting specialists' : 'Authorized personnel'}
- Third-party service providers under BAA

## Right to Revoke
You have the right to revoke this authorization at any time by submitting a written request.

## Expiration
This authorization expires on [DATE] or upon completion of ${purpose}.

## Required Elements
- Patient signature and date
- Witness signature (if required)
- Copy provided to patient

## Electronic Consent Implementation
` + "```typescript" + `
interface HIPAAConsent {
  patientId: string;
  consentType: '${purpose}';
  consentDate: Date;
  expirationDate: Date;
  digitalSignature: string;
  ipAddress: string;
  revokedDate?: Date;
  witnessSignature?: string;
}

async function recordConsent(consent: HIPAAConsent) {
  // Encrypt consent record
  const encryptedConsent = await encrypt(consent);
  
  // Store with audit trail
  await consentDatabase.create({
    ...encryptedConsent,
    auditTrail: {
      action: 'consent_granted',
      timestamp: new Date(),
      source: 'patient_portal'
    }
  });
}
` + "```" + `
    ` + "`" + `;

    return {
      content: [
        {
          type: 'text',
          text: form,
        },
      ],
    };
  }

  async validateAccessControls(implementation) {
    const feedback = [];

    if (implementation.includes('role') || implementation.includes('Role')) {
      feedback.push('✅ Role-based access control detected');
    } else {
      feedback.push('❌ Implement role-based access controls');
    }

    if (implementation.includes('audit') || implementation.includes('log')) {
      feedback.push('✅ Audit logging present');
    } else {
      feedback.push('❌ Add audit logging for access attempts');
    }

    if (implementation.includes('encrypt') || implementation.includes('hash')) {
      feedback.push('✅ Encryption/hashing detected');
    } else {
      feedback.push('⚠️ Consider encryption for sensitive operations');
    }

    return {
      content: [
        {
          type: 'text',
          text: ` + "`HIPAA Access Control Validation:\\n${feedback.join('\\n')}`" + `,
        },
      ],
    };
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('Healthcare HIPAA MCP Server running on stdio');
  }
}

if (require.main === module) {
  const server = new HIPAAMCPServer();
  server.run().catch(console.error);
}

module.exports = { HIPAAMCPServer };
`
}

// generateEcommerceMCPServers creates e-commerce-specific MCP servers
func (mg *MCPGenerator) generateEcommerceMCPServers(files map[string]string) {
	files["scripts/mcp/ecommerce-server.js"] = `#!/usr/bin/env node

// E-commerce MCP Server
// Provides e-commerce patterns and validation

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const { CallToolRequestSchema, ListToolsRequestSchema } = require('@modelcontextprotocol/sdk/types.js');

class EcommerceMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: 'ecommerce-mcp-server',
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
            name: 'validate_cart_logic',
            description: 'Validate shopping cart implementation',
            inputSchema: {
              type: 'object',
              properties: {
                code: {
                  type: 'string',
                  description: 'Shopping cart code to validate',
                },
              },
              required: ['code'],
            },
          },
          {
            name: 'generate_payment_flow',
            description: 'Generate secure payment flow template',
            inputSchema: {
              type: 'object',
              properties: {
                provider: {
                  type: 'string',
                  enum: ['stripe', 'paypal', 'square'],
                  description: 'Payment provider',
                },
              },
              required: ['provider'],
            },
          },
          {
            name: 'check_inventory_management',
            description: 'Check inventory management implementation',
            inputSchema: {
              type: 'object',
              properties: {
                implementation: {
                  type: 'string',
                  description: 'Inventory management code',
                },
              },
              required: ['implementation'],
            },
          },
        ],
      };
    });

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case 'validate_cart_logic':
            return await this.validateCartLogic(args.code);
          case 'generate_payment_flow':
            return await this.generatePaymentFlow(args.provider);
          case 'check_inventory_management':
            return await this.checkInventoryManagement(args.implementation);
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

  async validateCartLogic(code) {
    const issues = [];

    // Check for proper price calculations
    if ((code.includes('price') || code.includes('total')) && code.includes('parseFloat')) {
      issues.push('❌ Use Decimal/BigNumber for price calculations, not parseFloat');
    }

    // Check for inventory validation
    if (code.includes('addToCart') && !code.includes('inventory')) {
      issues.push('⚠️ Consider inventory validation before adding to cart');
    }

    // Check for session persistence
    if (code.includes('cart') && !code.includes('localStorage') && !code.includes('session')) {
      issues.push('⚠️ Consider cart persistence across sessions');
    }

    const result = issues.length > 0 ? 
      ` + "`Shopping Cart Issues:\\n${issues.join('\\n')}`" + ` :
      '✅ Cart logic looks good';

    return {
      content: [
        {
          type: 'text',
          text: result,
        },
      ],
    };
  }

  async generatePaymentFlow(provider) {
    const flows = {
      stripe: ` + "`" + `
# Stripe Payment Flow

## Frontend Implementation
` + "```" + `typescript
import { loadStripe } from '@stripe/stripe-js';

const stripePromise = loadStripe(process.env.STRIPE_PUBLISHABLE_KEY);

async function handlePayment(paymentData: PaymentRequest) {
  const stripe = await stripePromise;
  
  // Create payment intent on backend
  const response = await fetch('/api/payments/create-intent', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(paymentData)
  });
  
  const { clientSecret } = await response.json();
  
  // Confirm payment
  const result = await stripe.confirmCardPayment(clientSecret, {
    payment_method: {
      card: cardElement,
      billing_details: { name: customerName }
    }
  });
  
  if (result.error) {
    // Handle payment error
    showError(result.error.message);
  } else {
    // Payment succeeded
    redirectToSuccess(result.paymentIntent.id);
  }
}
` + "```" + `

## Backend Implementation
` + "```" + `javascript
const stripe = require('stripe')(process.env.STRIPE_SECRET_KEY);

app.post('/api/payments/create-intent', async (req, res) => {
  try {
    const { amount, currency, metadata } = req.body;
    
    const paymentIntent = await stripe.paymentIntents.create({
      amount: amount * 100, // Convert to cents
      currency: currency || 'usd',
      metadata,
      automatic_payment_methods: { enabled: true }
    });
    
    res.json({ clientSecret: paymentIntent.client_secret });
  } catch (error) {
    res.status(400).json({ error: error.message });
  }
});
` + "```" + `
      ` + "`" + `,
      // Add other providers...
    };

    return {
      content: [
        {
          type: 'text',
          text: flows[provider] || 'Payment provider not supported',
        },
      ],
    };
  }

  async checkInventoryManagement(implementation) {
    const feedback = [];

    if (implementation.includes('stock') || implementation.includes('quantity')) {
      feedback.push('✅ Stock tracking detected');
    } else {
      feedback.push('❌ Implement stock/quantity tracking');
    }

    if (implementation.includes('reserve') || implementation.includes('hold')) {
      feedback.push('✅ Inventory reservation logic present');
    } else {
      feedback.push('⚠️ Consider inventory reservation during checkout');
    }

    if (implementation.includes('low_stock') || implementation.includes('reorder')) {
      feedback.push('✅ Low stock handling detected');
    } else {
      feedback.push('⚠️ Add low stock alerts and reorder points');
    }

    return {
      content: [
        {
          type: 'text',
          text: ` + "`Inventory Management Analysis:\\n${feedback.join('\\n')}`" + `,
        },
      ],
    };
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('E-commerce MCP Server running on stdio');
  }
}

if (require.main === module) {
  const server = new EcommerceMCPServer();
  server.run().catch(console.error);
}

module.exports = { EcommerceMCPServer };
`
}