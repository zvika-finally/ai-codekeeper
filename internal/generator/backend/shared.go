package backend

import (
	"fmt"
	"strings"
)

// generateBackendReadme creates comprehensive README for backend
func (bg *BackendGenerator) generateBackendReadme() string {
	framework := bg.spec.GetBackendFramework()
	language := bg.spec.GetBackendLanguage()
	
	return fmt.Sprintf(`# %s Backend

%s backend API built with %s (%s)

## 🚀 Quick Start

### Prerequisites

%s

### Installation

%s

### Development

%s

### Production Deployment

%s

## 📡 API Endpoints

### Authentication
- "POST /api/auth/login" - User login
- "POST /api/auth/register" - User registration

### %s Management
- "GET /api/%s" - List all %s
- "POST /api/%s" - Create new %s
- "GET /api/%s/:id" - Get %s by ID
- "PUT /api/%s/:id" - Update %s
- "DELETE /api/%s/:id" - Delete %s

## 🔧 Configuration

Environment variables:

%s

## 🧪 Testing

%s

## 📊 Domain: %s

%s

## 🔐 Security

- JWT-based authentication
- Input validation and sanitization
- Rate limiting
- CORS configuration
- Environment-based secrets management

## 📝 Development Guidelines

%s

Generated with AI CodeKeeper v2.0.0 - Ready for AI implementation
`,
		bg.spec.Name,
		bg.spec.Description,
		framework,
		language,
		bg.getPrerequisites(),
		bg.getInstallationInstructions(),
		bg.getDevelopmentInstructions(),
		bg.getProductionInstructions(),
		bg.spec.CoreEntity,
		strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
		strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
		strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
		strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
		strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
		bg.getEnvironmentVariables(),
		bg.getTestingInstructions(),
		bg.spec.Domain,
		bg.getDomainSpecificGuidelines(),
		bg.getDevelopmentGuidelines())
}

// generateEnvExample creates .env.example file
func (bg *BackendGenerator) generateEnvExample() string {
	return fmt.Sprintf(`# Environment Configuration
NODE_ENV=development
PORT=8080

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/%s
REDIS_URL=redis://localhost:6379

# Authentication
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24h

# External Services
FRONTEND_URL=http://localhost:3000

# Logging
LOG_LEVEL=info

# Domain-specific configuration
%s

# Security
BCRYPT_ROUNDS=12
RATE_LIMIT_WINDOW_MS=900000
RATE_LIMIT_MAX_REQUESTS=100

# Development
DEBUG_MODE=true`, 
		strings.ToLower(bg.spec.Name),
		bg.getDomainEnvVars())
}

// Helper functions for content generation
func (bg *BackendGenerator) getPrerequisites() string {
	switch bg.spec.GetBackendLanguage() {
	case "javascript":
		return "- Node.js 18+ and npm\n- PostgreSQL 14+\n- Redis (optional)"
	case "python":
		return "- Python 3.9+\n- PostgreSQL 14+\n- Redis (optional)"
	case "go":
		return "- Go 1.21+\n- PostgreSQL 14+\n- Redis (optional)"
	default:
		return "- Runtime environment\n- Database\n- Cache (optional)"
	}
}

func (bg *BackendGenerator) getInstallationInstructions() string {
	switch bg.spec.GetBackendLanguage() {
	case "javascript":
		return "1. Install dependencies:\n   npm install\n\n2. Copy environment variables:\n   cp .env.example .env\n\n3. Configure your database connection in .env\n\n4. Run database migrations:\n   npm run migration:run"
	case "python":
		return "1. Create virtual environment:\n   python -m venv venv && source venv/bin/activate\n\n2. Install dependencies:\n   pip install -r requirements.txt\n\n3. Copy environment variables:\n   cp .env.example .env\n\n4. Run database migrations:\n   alembic upgrade head"
	case "go":
		return "1. Install dependencies:\n   go mod download\n\n2. Copy environment variables:\n   cp .env.example .env\n\n3. Build the application:\n   go build -o main ./src/"
	default:
		return "Follow language-specific installation instructions"
	}
}

func (bg *BackendGenerator) getDevelopmentInstructions() string {
	switch bg.spec.GetBackendLanguage() {
	case "javascript":
		return `1. Start development server:
   "npm run dev"

2. Run tests:
   "npm test"

3. Lint code:
   "npm run lint"`
	case "python":
		return `1. Start development server:
   "python src/main.py"

2. Run tests:
   "pytest"

3. Format code:
   "black . && isort ."`
	case "go":
		return `1. Start development server:
   "go run src/main.go"

2. Run tests:
   "go test ./..."

3. Format code:
   "go fmt ./..."`
	default:
		return "Follow language-specific development instructions"
	}
}

func (bg *BackendGenerator) getProductionInstructions() string {
	return `1. Build Docker image:
   "docker build -t "+bg.spec.Name+"-backend ."

2. Deploy using provided configurations:
   - Render: See "../infra/render/"
   - AWS: See "../infra/aws/"

3. Set production environment variables

4. Run database migrations in production`
}

func (bg *BackendGenerator) getEnvironmentVariables() string {
	envVars := `- "NODE_ENV" - Environment (development/production)
- "PORT" - Server port (default: 8080)
- "DATABASE_URL" - PostgreSQL connection string
- "REDIS_URL" - Redis connection string
- "JWT_SECRET" - JWT signing secret
- "FRONTEND_URL" - Frontend application URL`

	if bg.spec.Domain == "fintech" {
		envVars += `
- "STRIPE_SECRET_KEY" - Stripe payment processing
- "PLAID_CLIENT_ID" - Plaid banking integration
- "AUDIT_LOG_LEVEL" - Audit logging level`
	}

	return envVars
}

func (bg *BackendGenerator) getTestingInstructions() string {
	switch bg.spec.GetBackendLanguage() {
	case "javascript":
		return `- Unit tests: "npm test"
- Coverage: "npm run test:coverage"
- Integration tests included`
	case "python":
		return `- Unit tests: "pytest"
- Coverage: "pytest --cov"
- Async tests supported`
	case "go":
		return `- Unit tests: "go test ./..."
- Coverage: "go test -cover ./..."
- Benchmark tests: "go test -bench ."`
	default:
		return "Language-specific testing instructions"
	}
}

func (bg *BackendGenerator) getDomainSpecificGuidelines() string {
	switch bg.spec.Domain {
	case "fintech":
		return `**Financial Domain Implementation:**
- Use Decimal.js for all monetary calculations (never use floating point)
- Implement comprehensive audit trails for all transactions
- Add idempotency keys for transaction endpoints
- Implement proper error handling for payment failures
- Add rate limiting specific to transaction endpoints
- Use encryption for sensitive financial data`
	case "healthcare":
		return `**Healthcare Domain Implementation:**
- Ensure HIPAA compliance for all patient data handling
- Implement comprehensive audit logging for data access
- Add proper consent management workflows
- Use encryption for PHI (Personal Health Information)
- Implement proper data retention policies`
	default:
		return `**Domain-Specific Implementation:**
- Follow established business logic patterns
- Implement proper data validation rules
- Add domain-specific search and filtering
- Ensure proper audit logging`
	}
}

func (bg *BackendGenerator) getDevelopmentGuidelines() string {
	return `- Follow the established project structure
- Use TypeScript types from shared-types package
- Implement comprehensive error handling
- Add proper logging for debugging and monitoring
- Write tests for all new features
- Follow security best practices
- Use environment variables for configuration`
}

func (bg *BackendGenerator) getDomainEnvVars() string {
	switch bg.spec.Domain {
	case "fintech":
		return `# Financial domain
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
PLAID_CLIENT_ID=your_plaid_client_id
PLAID_SECRET=your_plaid_secret
AUDIT_LOG_LEVEL=DEBUG`
	case "healthcare":
		return `# Healthcare domain
HIPAA_COMPLIANCE_MODE=true
PHI_ENCRYPTION_KEY=your_phi_encryption_key
AUDIT_LOG_RETENTION_DAYS=2555`
	default:
		return `# Domain-specific environment variables
DOMAIN_FEATURE_FLAG=true`
	}
}