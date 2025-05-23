package infrastructure

import "strings"

// Additional AWS module generation methods

func (ig *InfrastructureGenerator) generateRDSModule() string {
	return `resource "aws_db_subnet_group" "main" {
  name       = "${var.name}-${var.environment}"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.name}-${var.environment}"
  }
}

resource "aws_security_group" "rds" {
  name_prefix = "${var.name}-rds"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [data.aws_vpc.main.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "main" {
  identifier = "${var.name}-${var.environment}"

  engine         = "postgres"
  engine_version = "15.4"
  instance_class = var.instance_class

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_type          = "gp2"
  storage_encrypted     = true

  db_name  = var.database_name
  username = var.master_username
  password = var.master_password

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name

  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"

  skip_final_snapshot = var.environment == "dev"
  deletion_protection = var.environment == "prod"

  # Domain-specific configuration
` + ig.getDomainSpecificRDSConfig() + `

  tags = {
    Name = "${var.name}-${var.environment}"
  }
}

data "aws_vpc" "main" {
  id = var.vpc_id
}
`
}

func (ig *InfrastructureGenerator) generateRDSVariables() string {
	return `variable "name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs"
  type        = list(string)
}

variable "database_name" {
  description = "Database name"
  type        = string
}

variable "master_username" {
  description = "Database master username"
  type        = string
}

variable "master_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

variable "instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}
`
}

func (ig *InfrastructureGenerator) generateVPCModule() string {
	return `resource "aws_vpc" "main" {
  cidr_block           = var.cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.name}-${var.environment}"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.name}-${var.environment}"
  }
}

resource "aws_subnet" "public" {
  count = length(var.availability_zones)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.cidr_block, 8, count.index)
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.name}-${var.environment}-public-${count.index + 1}"
    Type = "public"
  }
}

resource "aws_subnet" "private" {
  count = length(var.availability_zones)

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, count.index + length(var.availability_zones))
  availability_zone = var.availability_zones[count.index]

  tags = {
    Name = "${var.name}-${var.environment}-private-${count.index + 1}"
    Type = "private"
  }
}

resource "aws_eip" "nat" {
  count = length(var.availability_zones)

  domain = "vpc"
  depends_on = [aws_internet_gateway.main]

  tags = {
    Name = "${var.name}-${var.environment}-nat-${count.index + 1}"
  }
}

resource "aws_nat_gateway" "main" {
  count = length(var.availability_zones)

  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id

  tags = {
    Name = "${var.name}-${var.environment}-nat-${count.index + 1}"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.name}-${var.environment}-public"
  }
}

resource "aws_route_table" "private" {
  count = length(var.availability_zones)

  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }

  tags = {
    Name = "${var.name}-${var.environment}-private-${count.index + 1}"
  }
}

resource "aws_route_table_association" "public" {
  count = length(var.availability_zones)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count = length(var.availability_zones)

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}
`
}

func (ig *InfrastructureGenerator) generateVPCVariables() string {
	return `variable "name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "cidr_block" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}
`
}

func (ig *InfrastructureGenerator) generateVPCOutputs() string {
	return `output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = aws_subnet.private[*].id
}

output "internet_gateway_id" {
  description = "ID of the Internet Gateway"
  value       = aws_internet_gateway.main.id
}
`
}

func (ig *InfrastructureGenerator) generateSecretsModule() string {
	return `resource "aws_secretsmanager_secret" "secrets" {
  for_each = var.secrets

  name                    = "${var.name}-${var.environment}-${each.key}"
  description             = "${each.key} for ${var.name} ${var.environment}"
  recovery_window_in_days = var.environment == "prod" ? 30 : 0
}

resource "aws_secretsmanager_secret_version" "secrets" {
  for_each = var.secrets

  secret_id     = aws_secretsmanager_secret.secrets[each.key].id
  secret_string = each.value
}

# IAM policy for ECS tasks to access secrets
resource "aws_iam_policy" "secrets_access" {
  name        = "${var.name}-${var.environment}-secrets-access"
  description = "Policy to access secrets from Secrets Manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [
          for secret in aws_secretsmanager_secret.secrets : secret.arn
        ]
      }
    ]
  })
}
`
}

func (ig *InfrastructureGenerator) generateRenderConfig() string {
	return `# Render deployment configuration
services:
  # Backend service
  - type: web
    name: ` + ig.spec.Name + `-backend
    env: node
    buildCommand: npm install && npm run build
    startCommand: npm start
    envVars:
      - key: NODE_ENV
        value: production
      - key: DATABASE_URL
        fromDatabase:
          name: ` + ig.spec.Name + `-db
          property: connectionString
      - key: JWT_SECRET
        generateValue: true
      - key: PORT
        value: 8080
` + ig.getDomainSpecificRenderEnvVars() + `

  # Frontend service
  - type: static
    name: ` + ig.spec.Name + `-frontend
    buildCommand: npm install && npm run build
    staticPublishPath: ./dist
    envVars:
      - key: VITE_API_URL
        value: https://` + ig.spec.Name + `-backend.onrender.com

# Database
databases:
  - name: ` + ig.spec.Name + `-db
    databaseName: ` + strings.ReplaceAll(ig.spec.Name, "-", "_") + `
    user: ` + ig.spec.Name + `_user
` + ig.getDomainSpecificRenderDBConfig()
}

func (ig *InfrastructureGenerator) generateRenderReadme() string {
	return `# Render Deployment

Simple cloud deployment using Render.com platform.

## Setup

1. Connect your GitHub repository to Render
2. Use the ` + "`render.yaml`" + ` configuration for automatic deployment
3. Set up environment variables in Render dashboard

## Environment Variables

Set these in the Render dashboard:

### Backend Service
- ` + "`JWT_SECRET`" + `: Auto-generated secure secret
- ` + "`DATABASE_URL`" + `: Auto-configured from database
` + ig.getDomainSpecificRenderEnvDocs() + `

### Frontend Service
- ` + "`VITE_API_URL`" + `: Backend service URL

## Domain-Specific Configuration (` + ig.spec.Domain + `)

` + ig.getDomainSpecificRenderNotes() + `

## Deployment

1. Push to main branch
2. Render automatically builds and deploys
3. Database migrations run automatically
4. Health checks verify deployment

## Monitoring

- View logs in Render dashboard
- Set up uptime monitoring
- Configure alert notifications
`
}

func (ig *InfrastructureGenerator) getDomainSpecificRDSConfig() string {
	switch ig.spec.Domain {
	case "fintech":
		return `  # Enhanced security for financial data
  enabled_cloudwatch_logs_exports = ["postgresql"]
  monitoring_interval             = 60
  monitoring_role_arn            = aws_iam_role.rds_monitoring.arn
  performance_insights_enabled   = true`
	case "healthcare":
		return `  # HIPAA compliance requirements
  enabled_cloudwatch_logs_exports = ["postgresql"]
  monitoring_interval             = 60
  performance_insights_enabled   = true
  performance_insights_kms_key_id = aws_kms_key.rds.arn`
	default:
		return `  # Standard configuration
  enabled_cloudwatch_logs_exports = ["postgresql"]`
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificRenderEnvVars() string {
	switch ig.spec.Domain {
	case "fintech":
		return `      - key: STRIPE_SECRET_KEY
        sync: false
      - key: PLAID_CLIENT_ID
        sync: false
      - key: PLAID_SECRET
        sync: false`
	case "healthcare":
		return `      - key: FHIR_BASE_URL
        sync: false
      - key: FHIR_AUTH_TOKEN
        sync: false`
	case "ecommerce":
		return `      - key: SHOPIFY_ACCESS_TOKEN
        sync: false
      - key: STRIPE_SECRET_KEY
        sync: false`
	default:
		return ""
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificRenderDBConfig() string {
	switch ig.spec.Domain {
	case "fintech", "healthcare":
		return `    plan: starter  # Use standard+ for production with backups`
	default:
		return `    plan: free     # Upgrade for production use`
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificRenderEnvDocs() string {
	switch ig.spec.Domain {
	case "fintech":
		return `- ` + "`STRIPE_SECRET_KEY`" + `: Stripe API secret key
- ` + "`PLAID_CLIENT_ID`" + `: Plaid client ID for banking
- ` + "`PLAID_SECRET`" + `: Plaid secret key`
	case "healthcare":
		return `- ` + "`FHIR_BASE_URL`" + `: FHIR server base URL
- ` + "`FHIR_AUTH_TOKEN`" + `: FHIR authentication token`
	case "ecommerce":
		return `- ` + "`SHOPIFY_ACCESS_TOKEN`" + `: Shopify API token
- ` + "`STRIPE_SECRET_KEY`" + `: Stripe payment processing`
	default:
		return ""
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificRenderNotes() string {
	switch ig.spec.Domain {
	case "fintech":
		return `### Security Requirements
- Use Render's SOC 2 compliant infrastructure
- Enable database encryption
- Configure proper backup policies
- Set up monitoring and alerting`
	case "healthcare":
		return `### HIPAA Compliance
- Render provides BAA (Business Associate Agreement)
- Enable audit logging
- Configure data retention policies
- Implement proper access controls`
	case "ecommerce":
		return `### Performance Optimization
- Use CDN for static assets
- Configure caching strategies
- Set up auto-scaling policies
- Monitor transaction performance`
	default:
		return `Standard deployment configuration with basic monitoring and security.`
	}
}