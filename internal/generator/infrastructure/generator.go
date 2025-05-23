package infrastructure

import (
	"strings"
)

// InfrastructureGenerator handles infrastructure as code generation
type InfrastructureGenerator struct {
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

// NewInfrastructureGenerator creates a new infrastructure generator
func NewInfrastructureGenerator(spec *ProjectSpec) *InfrastructureGenerator {
	return &InfrastructureGenerator{spec: spec}
}

// Generate creates infrastructure files
func (ig *InfrastructureGenerator) Generate() (map[string]string, error) {
	files := make(map[string]string)
	
	// Generate Docker infrastructure
	ig.generateDockerInfrastructure(files)
	
	// Generate AWS infrastructure
	ig.generateAWSInfrastructure(files)
	
	// Generate Render infrastructure
	ig.generateRenderInfrastructure(files)
	
	// Generate CI/CD configuration
	ig.generateCICDConfiguration(files)
	
	return files, nil
}

// generateDockerInfrastructure creates Docker configurations
func (ig *InfrastructureGenerator) generateDockerInfrastructure(files map[string]string) {
	// Docker Compose for development
	files["docker-compose.yml"] = ig.generateDockerCompose()
	files["docker-compose.prod.yml"] = ig.generateDockerComposeProd()
	
	// Docker configuration guide
	files["infra/docker/README.md"] = ig.generateDockerReadme()
}

// generateAWSInfrastructure creates AWS Terraform configurations
func (ig *InfrastructureGenerator) generateAWSInfrastructure(files map[string]string) {
	// Main Terraform configuration
	files["infra/aws/main.tf"] = ig.generateAWSMain()
	files["infra/aws/variables.tf"] = ig.generateAWSVariables()
	files["infra/aws/outputs.tf"] = ig.generateAWSOutputs()
	files["infra/aws/terraform.tfvars.example"] = ig.generateAWSTerraformVars()
	
	// ECS module
	files["infra/aws/modules/ecs/main.tf"] = ig.generateECSModule()
	files["infra/aws/modules/ecs/variables.tf"] = ig.generateECSVariables()
	files["infra/aws/modules/ecs/outputs.tf"] = ig.generateECSOutputs()
	
	// RDS module
	files["infra/aws/modules/rds/main.tf"] = ig.generateRDSModule()
	files["infra/aws/modules/rds/variables.tf"] = ig.generateRDSVariables()
	
	// VPC module
	files["infra/aws/modules/vpc/main.tf"] = ig.generateVPCModule()
	files["infra/aws/modules/vpc/variables.tf"] = ig.generateVPCVariables()
	files["infra/aws/modules/vpc/outputs.tf"] = ig.generateVPCOutputs()
	
	// Secrets module
	files["infra/aws/modules/secrets/main.tf"] = ig.generateSecretsModule()
	
	// AWS deployment guide
	files["infra/aws/README.md"] = ig.generateAWSReadme()
}

// generateRenderInfrastructure creates Render deployment configurations
func (ig *InfrastructureGenerator) generateRenderInfrastructure(files map[string]string) {
	files["infra/render/render.yaml"] = ig.generateRenderConfig()
	files["infra/render/README.md"] = ig.generateRenderReadme()
}

// generateCICDConfiguration creates CI/CD pipeline configurations
func (ig *InfrastructureGenerator) generateCICDConfiguration(files map[string]string) {
	// GitHub Actions
	files[".github/workflows/ci.yml"] = ig.generateGitHubActions()
	files[".github/workflows/deploy.yml"] = ig.generateDeployWorkflow()
	
	// Health check script
	files["scripts/health-check.sh"] = ig.generateHealthCheckScript()
}

// Docker generation methods
func (ig *InfrastructureGenerator) generateDockerCompose() string {
	return `version: '3.8'

services:
  # Backend service
  backend:
    build:
      context: ./apps/backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=development
      - DATABASE_URL=postgresql://postgres:password@db:5432/` + ig.spec.Name + `
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    volumes:
      - ./apps/backend:/app
      - /app/node_modules
    command: npm run dev

  # Frontend service
  frontend:
    build:
      context: ./apps/frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - VITE_API_URL=http://localhost:8080
    volumes:
      - ./apps/frontend:/app
      - /app/node_modules
    command: npm run dev

  # Database service
  db:
    image: postgres:15-alpine
    restart: unless-stopped
    environment:
      POSTGRES_DB: ` + ig.spec.Name + `
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init.sql

  # Redis for caching
  redis:
    image: redis:7-alpine
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # Nginx reverse proxy (production)
  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./infra/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./infra/nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - backend
      - frontend
    profiles:
      - production

volumes:
  postgres_data:
  redis_data:

networks:
  default:
    driver: bridge
`
}

func (ig *InfrastructureGenerator) generateDockerComposeProd() string {
	return `version: '3.8'

# Production Docker Compose configuration
# Use with: docker-compose -f docker-compose.yml -f docker-compose.prod.yml up

services:
  backend:
    environment:
      - NODE_ENV=production
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
    restart: unless-stopped
    command: npm start
    volumes: []

  frontend:
    environment:
      - VITE_API_URL=${API_URL}
    restart: unless-stopped
    command: npm run preview
    volumes: []

  db:
    restart: unless-stopped
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_prod_data:/var/lib/postgresql/data

  redis:
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD}

volumes:
  postgres_prod_data:
    external: true
`
}

func (ig *InfrastructureGenerator) generateDockerReadme() string {
	return `# Docker Infrastructure

## Development Environment

Start the full development stack:

` + "```bash" + `
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
` + "```" + `

## Production Environment

Deploy to production:

` + "```bash" + `
# Build production images
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Start production stack
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
` + "```" + `

## Environment Variables

Create ` + "`.env`" + ` file for local development:

` + "```" + `
DATABASE_URL=postgresql://postgres:password@localhost:5432/` + ig.spec.Name + `
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
API_URL=http://localhost:8080
` + "```" + `

## Services

- **Backend**: ` + strings.Title(ig.spec.Backend) + ` API server (port 8080)
- **Frontend**: React TypeScript app (port 3000)
- **Database**: PostgreSQL (port 5432)
- **Redis**: Caching and sessions (port 6379)
- **Nginx**: Reverse proxy (ports 80/443, production only)

## Domain-Specific Notes (` + ig.spec.Domain + `)

` + ig.getDomainSpecificDockerNotes() + `

## Health Checks

` + "```bash" + `
# Check backend health
curl http://localhost:8080/health

# Check frontend
curl http://localhost:3000

# Check database connection
docker-compose exec db psql -U postgres -d ` + ig.spec.Name + ` -c "SELECT 1;"
` + "```" + `
`
}

func (ig *InfrastructureGenerator) getDomainSpecificDockerNotes() string {
	switch ig.spec.Domain {
	case "fintech":
		return `### Fintech Security Requirements
- Database encryption at rest enabled
- Redis password protection required
- SSL/TLS certificates for HTTPS
- Separate secrets management (see AWS Secrets Manager)
- PCI-DSS compliant container security`
	case "healthcare":
		return `### Healthcare Compliance Requirements
- HIPAA-compliant container configuration
- Encrypted volumes for patient data
- Access logging enabled
- Network isolation for sensitive services
- Backup encryption required`
	case "ecommerce":
		return `### E-commerce Optimization
- Redis caching for product catalog
- CDN integration for static assets
- Session persistence for shopping carts
- Payment gateway health monitoring
- Inventory service scaling`
	default:
		return `### General Web Application
- Standard security practices applied
- Development/production environment separation
- Health monitoring and logging enabled`
	}
}