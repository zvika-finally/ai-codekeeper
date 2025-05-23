package generator

import "strings"

// generateProjectReadme creates the main project README
func (g *NewModularGenerator) generateProjectReadme() string {
	return `# ` + g.spec.Name + `

` + g.spec.Description + `

## Architecture

This project follows clean architecture principles with:
- **Domain**: ` + g.spec.Domain + ` business logic
- **Backend**: ` + g.spec.Backend + ` API with ` + g.spec.APIStyle + ` design
- **Frontend**: React TypeScript application
- **Database**: Configured for scalable data storage

## Development Guidelines

See the following documentation for development standards:
- [Frontend Standards](docs/frontend/STANDARDS.md)
- [API Design](docs/API_DESIGN.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Contributing](docs/CONTRIBUTING.md)

## Getting Started

1. Follow [deployment guide](docs/DEPLOYMENT.md)
2. Review [infrastructure setup](docs/infrastructure/)
3. Check domain-specific patterns in docs/

## Core Entity

This application centers around: **` + g.spec.CoreEntity + `**

## User Roles

` + strings.Join(g.spec.GetUserRolesList(), ", ") + `
`
}

// generateAPIDesignDoc creates API design guidelines
func (g *NewModularGenerator) generateAPIDesignDoc() string {
	return `# API Design Guidelines

## REST API Standards

### Endpoint Design
- Use nouns for resources: ` + "`/api/v1/" + strings.ToLower(g.spec.CoreEntity) + "s`" + `
- Use HTTP methods appropriately (GET, POST, PUT, DELETE)
- Include version in URL: ` + "`/api/v1/`" + `
- Use consistent naming conventions

### Request/Response Format
- JSON content type for all requests/responses
- Include pagination for list endpoints
- Standard error response format
- Consistent field naming (camelCase)

### Authentication
- JWT tokens for API authentication
- Include Bearer token in Authorization header
- Implement token refresh mechanism
- Rate limiting per user/endpoint

### Domain-Specific API Patterns (` + g.spec.Domain + `)
Based on ` + g.spec.Domain + ` domain requirements:
- Implement domain-specific validation rules
- Follow industry-specific API standards
- Include domain-specific security measures
- Implement required compliance endpoints

### Error Handling
- Use appropriate HTTP status codes
- Include error codes and messages
- Provide helpful error descriptions
- Log errors for debugging

### Documentation
- Use OpenAPI/Swagger specification
- Include request/response examples
- Document authentication requirements
- Provide testing guidelines
`
}

// generateDeploymentDoc creates deployment guidelines
func (g *NewModularGenerator) generateDeploymentDoc() string {
	return `# Deployment Guide

## Environment Setup

### Development
` + "```bash" + `
# Install dependencies
npm install  # Frontend
go mod download  # Backend

# Run locally
npm run dev  # Frontend (port 3000)
go run main.go  # Backend (port 8080)
` + "```" + `

### Production

#### Docker Deployment
` + "```bash" + `
# Build and run with Docker Compose
docker-compose up --build
` + "```" + `

#### Cloud Deployment
- **Frontend**: Deploy to Vercel/Netlify
- **Backend**: Deploy to Railway/Render/AWS
- **Database**: Use managed database service

### Environment Variables
` + "```" + `
# Backend
DATABASE_URL=postgresql://...
JWT_SECRET=your-secret-key
PORT=8080

# Frontend
VITE_API_URL=https://your-api.com
` + "```" + `

### Security Configuration
- Enable HTTPS in production
- Configure CORS properly
- Set secure headers
- Use environment-specific secrets

### Monitoring
- Set up health checks
- Configure logging
- Monitor performance metrics
- Set up alerting

### Domain-Specific Deployment (` + g.spec.Domain + `)
- Follow ` + g.spec.Domain + ` compliance requirements
- Configure domain-specific security
- Set up required monitoring
- Implement backup strategies
`
}

// generateContributingDoc creates contribution guidelines
func (g *NewModularGenerator) generateContributingDoc() string {
	return `# Contributing Guidelines

## Development Workflow

1. **Create Feature Branch**
   ` + "```bash" + `
   git checkout -b feature/your-feature-name
   ` + "```" + `

2. **Follow Code Standards**
   - Frontend: Follow [Frontend Standards](frontend/STANDARDS.md)
   - Backend: Use ` + g.spec.Backend + ` best practices
   - Database: Follow normalization principles

3. **Testing Requirements**
   - Write unit tests for new features
   - Ensure integration tests pass
   - Test edge cases and error scenarios

4. **Code Review Process**
   - Create pull request with clear description
   - Include testing instructions
   - Address review feedback promptly

## Code Standards

### Git Commit Messages
` + "```" + `
feat: add new feature
fix: resolve bug
docs: update documentation
refactor: improve code structure
test: add missing tests
` + "```" + `

### Code Quality
- Use TypeScript strict mode
- Follow linting rules
- Maintain code coverage above 80%
- Write clear, documented code

### Domain Guidelines (` + g.spec.Domain + `)
- Follow ` + g.spec.Domain + ` industry standards
- Implement required compliance measures
- Use domain-specific patterns
- Validate business rules properly

## Pull Request Template
1. **Description**: What changes are included?
2. **Testing**: How was this tested?
3. **Impact**: What areas are affected?
4. **Documentation**: Any docs updated?
`
}

// generateArchitectureDoc creates architecture guidelines
func (g *NewModularGenerator) generateArchitectureDoc() string {
	return `# Architecture Documentation

## System Overview

This ` + g.spec.Domain + ` application follows a modular architecture:

### Frontend (React TypeScript)
- **Components**: Reusable UI components
- **Pages**: Route-specific components  
- **Hooks**: Custom React hooks for logic
- **Services**: API communication layer
- **Stores**: Client-side state management

### Backend (` + strings.Title(g.spec.Backend) + `)
- **Routes**: API endpoint definitions
- **Controllers**: Request/response handling
- **Services**: Business logic layer
- **Models**: Data models and validation
- **Database**: Data persistence layer

### Infrastructure
- **Docker**: Containerized deployment
- **CI/CD**: Automated testing and deployment
- **Monitoring**: Application health and metrics
- **Security**: Authentication and authorization

## Design Principles

### Clean Architecture
- Separation of concerns
- Dependency inversion
- Domain-driven design
- SOLID principles

### Scalability
- Horizontal scaling capability
- Stateless service design
- Caching strategies
- Database optimization

### Security
- Input validation at all layers
- Authentication and authorization
- Data encryption in transit/rest
- Regular security audits

### Domain-Specific Architecture (` + g.spec.Domain + `)
- ` + strings.Title(g.spec.Domain) + `-specific business rules
- Industry compliance requirements
- Domain entity relationships
- Specialized data patterns

## Data Flow
1. User interaction in React frontend
2. API call to ` + g.spec.Backend + ` backend
3. Business logic processing
4. Database operations
5. Response back to frontend
6. UI state updates

## Technology Decisions
- **Frontend**: React for component-based UI
- **Backend**: ` + strings.Title(g.spec.Backend) + ` for ` + g.spec.APIStyle + ` APIs
- **Database**: Relational for structured ` + g.spec.Domain + ` data
- **Deployment**: Docker for consistency
`
}

// generateDockerGuideDoc creates Docker deployment guide
func (g *NewModularGenerator) generateDockerGuideDoc() string {
	return `# Docker Deployment Guide

## Container Strategy

### Multi-Stage Builds
- Optimize image sizes
- Separate build and runtime environments
- Cache dependencies effectively

### Service Architecture
` + "```yaml" + `
services:
  frontend:
    build: ./apps/frontend
    ports: ["3000:3000"]
    
  backend:
    build: ./apps/backend
    ports: ["8080:8080"]
    environment:
      - DATABASE_URL=postgresql://db:5432/app
      
  database:
    image: postgres:15
    environment:
      - POSTGRES_DB=app
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
` + "```" + `

### Production Considerations
- Use specific image tags, not 'latest'
- Set resource limits
- Configure health checks
- Use secrets management
- Enable logging drivers

### Security Best Practices
- Run as non-root user
- Scan images for vulnerabilities
- Use minimal base images
- Keep containers updated

### Domain-Specific Containers (` + g.spec.Domain + `)
- Include ` + g.spec.Domain + `-specific tools
- Configure compliance requirements
- Set up required monitoring agents
- Include security scanning tools
`
}

// generateCICDGuideDoc creates CI/CD pipeline guide
func (g *NewModularGenerator) generateCICDGuideDoc() string {
	return `# CI/CD Pipeline Guide

## Pipeline Stages

### 1. Code Quality
` + "```yaml" + `
test:
  script:
    - npm run lint        # Frontend linting
    - npm run type-check  # TypeScript checking
    - go vet ./...        # Go code analysis
    - go test ./...       # Backend tests
    - npm run test        # Frontend tests
` + "```" + `

### 2. Security Scanning
- Dependency vulnerability scanning
- Static code analysis
- Container image scanning
- Secret detection

### 3. Build & Package
- Build production artifacts
- Create Docker images
- Run integration tests
- Generate documentation

### 4. Deployment
- Deploy to staging environment
- Run smoke tests
- Deploy to production (manual approval)
- Monitor deployment health

## Environment Promotion

### Development → Staging → Production
- Automatic deployment to staging
- Manual approval for production
- Rollback capabilities
- Environment-specific configurations

### Database Migrations
- Version-controlled schema changes
- Automated migration testing
- Rollback scripts available
- Zero-downtime deployments

### Domain-Specific CI/CD (` + g.spec.Domain + `)
- ` + strings.Title(g.spec.Domain) + ` compliance testing
- Industry-specific security scans
- Regulatory approval workflows
- Audit trail generation

## Monitoring & Alerting
- Build success/failure notifications
- Performance regression detection
- Security vulnerability alerts
- Deployment status monitoring
`
}

// generateMonitoringDoc creates monitoring guidelines
func (g *NewModularGenerator) generateMonitoringDoc() string {
	return `# Monitoring & Observability

## Application Metrics

### Key Performance Indicators
- Response time (p95, p99)
- Error rates by endpoint
- Throughput (requests per second)
- Database query performance

### Business Metrics (` + g.spec.Domain + `)
- Domain-specific KPIs
- User engagement metrics
- Business process success rates
- Compliance monitoring

## Logging Strategy

### Structured Logging
` + "```json" + `
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "service": "` + g.spec.Name + `",
  "user_id": "123",
  "action": "` + strings.ToLower(g.spec.CoreEntity) + `_created",
  "duration_ms": 150
}
` + "```" + `

### Log Levels
- **ERROR**: System errors requiring attention
- **WARN**: Potential issues to monitor
- **INFO**: Important business events
- **DEBUG**: Detailed debugging information

## Health Checks

### Application Health
- Database connectivity
- External API availability
- Memory and CPU usage
- Disk space monitoring

### Custom Health Checks
- Domain-specific validations
- Business rule integrity
- Data consistency checks
- Integration point status

## Alerting Rules

### Critical Alerts
- Service downtime
- High error rates (>5%)
- Database connection failures
- Security incidents

### Warning Alerts
- Performance degradation
- High resource usage
- Failed business processes
- Compliance violations

## Dashboard Design
- Executive overview dashboard
- Technical operations dashboard
- Domain-specific metrics dashboard
- Real-time incident dashboard
`
}

// generateSecurityDoc creates security guidelines
func (g *NewModularGenerator) generateSecurityDoc() string {
	return `# Security Guidelines

## Authentication & Authorization

### JWT Implementation
- Use strong secret keys
- Implement token expiration
- Include refresh token mechanism
- Store tokens securely

### Role-Based Access Control
User roles: ` + strings.Join(g.spec.GetUserRolesList(), ", ") + `
- Define granular permissions
- Implement role hierarchy
- Regular access reviews
- Principle of least privilege

## Data Protection

### Encryption
- HTTPS for all communications
- Database field encryption for sensitive data
- File encryption for stored documents
- Key rotation policies

### Input Validation
- Sanitize all user inputs
- Use parameterized queries
- Validate file uploads
- Rate limiting on endpoints

## Domain-Specific Security (` + g.spec.Domain + `)

### Compliance Requirements
- Industry-specific regulations
- Data retention policies
- Audit trail requirements
- Privacy controls

### Threat Model
- Identify ` + g.spec.Domain + ` specific threats
- Implement appropriate countermeasures
- Regular security assessments
- Incident response procedures

## Security Monitoring

### Detection
- Failed authentication attempts
- Unusual access patterns
- Data export activities
- System configuration changes

### Response
- Automated threat response
- Incident escalation procedures
- Forensic data collection
- Recovery processes

## Secure Development

### Code Security
- Static code analysis
- Dependency vulnerability scanning
- Security code reviews
- Penetration testing

### Infrastructure Security
- Container security scanning
- Network segmentation
- Firewall configuration
- Regular security updates
`
}

// generateScalingDoc creates scaling guidelines
func (g *NewModularGenerator) generateScalingDoc() string {
	return `# Scaling Guidelines

## Horizontal Scaling

### Application Tier
- Stateless service design
- Load balancer configuration
- Session management strategy
- Auto-scaling policies

### Database Scaling
- Read replicas for query distribution
- Connection pooling
- Query optimization
- Caching strategies

## Performance Optimization

### Frontend Performance
- Code splitting and lazy loading
- CDN for static assets
- Browser caching strategies
- Progressive web app features

### Backend Performance
- Database indexing strategy
- API response caching
- Background job processing
- Resource pooling

## Caching Strategy

### Levels of Caching
1. **Browser Cache**: Static assets
2. **CDN Cache**: Global content delivery
3. **Application Cache**: Computed results
4. **Database Cache**: Query results

### Cache Invalidation
- Time-based expiration
- Event-driven invalidation
- Cache warming strategies
- Monitoring cache hit rates

## Domain-Specific Scaling (` + g.spec.Domain + `)

### Business Growth Patterns
- Anticipated ` + g.spec.Domain + ` growth
- Seasonal usage patterns
- Geographic expansion needs
- Compliance scaling requirements

### Data Growth Management
- Archiving strategies
- Data partitioning
- Backup and recovery scaling
- Analytics data handling

## Monitoring Scaling

### Key Metrics
- Response time under load
- Resource utilization trends
- Error rates during peak usage
- Business metrics correlation

### Capacity Planning
- Traffic growth projections
- Resource requirement forecasting
- Cost optimization strategies
- Performance testing protocols
`
}