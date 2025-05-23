package infrastructure

// CI/CD generation methods

func (ig *InfrastructureGenerator) generateGitHubActions() string {
	return `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  NODE_VERSION: '18'
  DOCKER_REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: |
          cd apps/backend && npm ci
          cd ../frontend && npm ci

      - name: Run linting
        run: |
          cd apps/backend && npm run lint
          cd ../frontend && npm run lint

      - name: Run type checking
        run: |
          cd apps/frontend && npm run type-check

      - name: Run backend tests
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost:5432/test_db
          JWT_SECRET: test-secret
          NODE_ENV: test
        run: |
          cd apps/backend && npm test

      - name: Run frontend tests
        run: |
          cd apps/frontend && npm test

      - name: Build applications
        run: |
          cd apps/backend && npm run build
          cd ../frontend && npm run build

  security:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run security audit
        run: |
          cd apps/backend && npm audit --audit-level high
          cd ../frontend && npm audit --audit-level high

      - name: Run CodeQL Analysis
        uses: github/codeql-action/init@v3
        with:
          languages: javascript

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v3

` + ig.getDomainSpecificCISteps() + `

  build-images:
    runs-on: ubuntu-latest
    needs: [test, security]
    if: github.ref == 'refs/heads/main'
    
    permissions:
      contents: read
      packages: write
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.DOCKER_REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.DOCKER_REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-

      - name: Build and push backend image
        uses: docker/build-push-action@v5
        with:
          context: ./apps/backend
          push: true
          tags: ${{ env.DOCKER_REGISTRY }}/${{ env.IMAGE_NAME }}-backend:${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}

      - name: Build and push frontend image
        uses: docker/build-push-action@v5
        with:
          context: ./apps/frontend
          push: true
          tags: ${{ env.DOCKER_REGISTRY }}/${{ env.IMAGE_NAME }}-frontend:${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
`
}

func (ig *InfrastructureGenerator) generateDeployWorkflow() string {
	return `name: Deploy

on:
  push:
    branches: [ main ]
  workflow_dispatch:
    inputs:
      environment:
        description: 'Environment to deploy to'
        required: true
        default: 'staging'
        type: choice
        options:
        - staging
        - production

env:
  AWS_REGION: us-east-1

jobs:
  deploy-staging:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' || github.event.inputs.environment == 'staging'
    environment: staging
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Deploy to ECS
        run: |
          # Update ECS service with new task definition
          aws ecs update-service \
            --cluster ` + ig.spec.Name + `-staging \
            --service ` + ig.spec.Name + `-backend-staging \
            --force-new-deployment

      - name: Wait for deployment
        run: |
          aws ecs wait services-stable \
            --cluster ` + ig.spec.Name + `-staging \
            --services ` + ig.spec.Name + `-backend-staging

      - name: Run health check
        run: |
          chmod +x ./scripts/health-check.sh
          ./scripts/health-check.sh ${{ secrets.STAGING_URL }}

  deploy-production:
    runs-on: ubuntu-latest
    if: github.event.inputs.environment == 'production'
    environment: production
    needs: [deploy-staging]
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Deploy to production ECS
        run: |
          aws ecs update-service \
            --cluster ` + ig.spec.Name + `-production \
            --service ` + ig.spec.Name + `-backend-production \
            --force-new-deployment

      - name: Wait for deployment
        run: |
          aws ecs wait services-stable \
            --cluster ` + ig.spec.Name + `-production \
            --services ` + ig.spec.Name + `-backend-production

      - name: Run production health check
        run: |
          chmod +x ./scripts/health-check.sh
          ./scripts/health-check.sh ${{ secrets.PRODUCTION_URL }}

      - name: Notify deployment
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: '` + ig.spec.Name + ` deployed to production'
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
`
}

func (ig *InfrastructureGenerator) generateHealthCheckScript() string {
	return `#!/bin/bash

# Health check script for ` + ig.spec.Name + `
# Usage: ./health-check.sh <base_url>

set -e

BASE_URL=${1:-"http://localhost:8080"}
MAX_ATTEMPTS=30
DELAY=10

echo "🏥 Running health checks for ` + ig.spec.Name + `..."
echo "Base URL: $BASE_URL"

# Function to check HTTP endpoint
check_endpoint() {
    local endpoint=$1
    local expected_status=${2:-200}
    local url="$BASE_URL$endpoint"
    
    echo "Checking $url..."
    
    for i in $(seq 1 $MAX_ATTEMPTS); do
        if response=$(curl -s -o /dev/null -w "%{http_code}" "$url"); then
            if [ "$response" = "$expected_status" ]; then
                echo "✅ $endpoint is healthy (HTTP $response)"
                return 0
            fi
        fi
        
        echo "⏳ Attempt $i/$MAX_ATTEMPTS failed, retrying in ${DELAY}s..."
        sleep $DELAY
    done
    
    echo "❌ $endpoint health check failed after $MAX_ATTEMPTS attempts"
    return 1
}

# Basic health checks
check_endpoint "/health"

# API endpoint checks
check_endpoint "/api/health"

# Domain-specific health checks
` + ig.getDomainSpecificHealthChecks() + `

echo "🎉 All health checks passed!"
`
}

func (ig *InfrastructureGenerator) generateAWSReadme() string {
	return `# AWS Infrastructure

This directory contains Terraform configurations for deploying ` + ig.spec.Name + ` to AWS.

## Prerequisites

1. AWS CLI configured with appropriate permissions
2. Terraform >= 1.0 installed
3. Docker for building images

## Setup

1. **Initialize Terraform**
   ` + "```bash" + `
   cd infra/aws
   terraform init
   ` + "```" + `

2. **Configure Variables**
   ` + "```bash" + `
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your values
   ` + "```" + `

3. **Plan Deployment**
   ` + "```bash" + `
   terraform plan
   ` + "```" + `

4. **Deploy Infrastructure**
   ` + "```bash" + `
   terraform apply
   ` + "```" + `

## Architecture

- **VPC**: Multi-AZ setup with public and private subnets
- **ECS**: Fargate containers for backend and frontend
- **RDS**: PostgreSQL database with encryption
- **ALB**: Application Load Balancer with SSL termination
- **Secrets Manager**: Secure storage for sensitive data

## Domain-Specific Configuration (` + ig.spec.Domain + `)

` + ig.getDomainSpecificAWSNotes() + `

## Environments

- **dev**: Single AZ, smaller instances
- **staging**: Multi-AZ, production-like
- **prod**: Multi-AZ, high availability, backups

## Monitoring

- CloudWatch logs and metrics
- Application performance monitoring
- Database performance insights
- Custom alarms and notifications

## Security

- WAF protection
- Security groups with least privilege
- Encrypted storage and transit
- Secrets rotation
- Compliance logging

## Cost Optimization

- Right-sized instances
- Scheduled scaling
- Reserved instances for production
- S3 lifecycle policies
`
}

func (ig *InfrastructureGenerator) getDomainSpecificCISteps() string {
	switch ig.spec.Domain {
	case "fintech":
		return `  compliance:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: PCI DSS Compliance Check
        run: |
          # Add PCI DSS specific compliance checks
          echo "Running PCI DSS compliance validation..."

      - name: SOX Compliance Check
        run: |
          # Add SOX compliance checks
          echo "Running SOX compliance validation..."`
	case "healthcare":
		return `  compliance:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: HIPAA Compliance Check
        run: |
          # Add HIPAA specific compliance checks
          echo "Running HIPAA compliance validation..."`
	default:
		return ""
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificHealthChecks() string {
	switch ig.spec.Domain {
	case "fintech":
		return `# Fintech-specific checks
check_endpoint "/api/payments/health"
check_endpoint "/api/accounts/health"
check_endpoint "/api/compliance/status"`
	case "healthcare":
		return `# Healthcare-specific checks
check_endpoint "/api/patients/health"
check_endpoint "/api/fhir/health"
check_endpoint "/api/compliance/hipaa"`
	case "ecommerce":
		return `# E-commerce-specific checks
check_endpoint "/api/products/health"
check_endpoint "/api/orders/health"
check_endpoint "/api/payments/health"`
	default:
		return `# Standard API checks
check_endpoint "/api/status"`
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificAWSNotes() string {
	switch ig.spec.Domain {
	case "fintech":
		return `### Financial Services Requirements
- Enhanced monitoring and alerting
- Audit logging to CloudTrail
- Encryption at rest and in transit
- SOC 2 and PCI DSS compliance
- Multi-region disaster recovery
- Real-time fraud detection integration`
	case "healthcare":
		return `### Healthcare Compliance
- HIPAA-compliant infrastructure
- BAA with AWS (Business Associate Agreement)
- Audit logging for all data access
- Encryption with customer-managed keys
- Data retention policies
- Access controls and role-based permissions`
	case "ecommerce":
		return `### E-commerce Optimization
- CDN integration for global performance
- Auto-scaling for traffic spikes
- Payment processing compliance
- Inventory management integration
- Customer data protection (GDPR)
- Real-time analytics and reporting`
	default:
		return `### Standard Web Application
- Basic security and monitoring
- Cost-optimized configuration
- Standard compliance practices`
	}
}