package frontend

// generateFrontendStandardsDoc creates the main frontend development standards
func (fg *FrontendGenerator) generateFrontendStandardsDoc() string {
	return `# Frontend Development Standards

## Architecture Principles

### Component Design
- Use functional components with TypeScript
- Implement single responsibility principle
- Follow composition over inheritance
- Use custom hooks for complex logic

### State Management
- Use React Query for server state
- Use Zustand for client state
- Avoid prop drilling - use context when needed
- Keep state as close to usage as possible

### Code Organization
` + "```" + `
src/
├── components/          # Reusable UI components
├── pages/              # Route components
├── hooks/              # Custom hooks
├── services/           # API calls
├── stores/             # Global state
├── utils/              # Helper functions
└── types/              # TypeScript definitions
` + "```" + `

### TypeScript Standards
- Use strict mode
- Define interfaces for all props and data
- Use generic types for reusable components
- Avoid 'any' - use 'unknown' when needed

### Performance Guidelines
- Use React.memo for expensive components
- Implement proper key props for lists
- Use useCallback for event handlers
- Use useMemo for expensive calculations

### Security Best Practices
- Sanitize user inputs
- Use HTTPS for all API calls
- Implement proper authentication headers
- Validate data at component boundaries
`
}

// generateComponentPatternsDoc creates component pattern guidelines
func (fg *FrontendGenerator) generateComponentPatternsDoc() string {
	return `# Component Patterns

## Standard Component Structure

` + "```typescript" + `
interface ComponentProps {
  // Define all props with types
}

const Component: React.FC<ComponentProps> = ({ prop1, prop2 }) => {
  // 1. State declarations
  // 2. Effect hooks
  // 3. Event handlers
  // 4. Render logic
  
  return (
    <div>
      {/* JSX content */}
    </div>
  );
};

export default Component;
` + "```" + `

## Form Components
- Use controlled components
- Implement proper validation
- Handle loading and error states
- Use React Hook Form for complex forms

## Data Display Components
- Handle loading states
- Implement error boundaries
- Use skeleton loading
- Paginate large datasets

## Domain-Specific Components
Based on domain: ` + fg.spec.Domain + `

### Required Component Types:
- List/Table components for data display
- Form components for data entry
- Modal/Dialog components for actions
- Navigation components for user flow
`
}

// generateStateManagementDoc creates state management guidelines
func (fg *FrontendGenerator) generateStateManagementDoc() string {
	return `# State Management Guidelines

## State Types

### Server State (React Query)
- Use for API data
- Implement caching strategies
- Handle error and loading states
- Use optimistic updates

### Client State (Zustand)
- Use for UI state
- Keep stores focused and small
- Implement persistence when needed
- Use devtools for debugging

### Local State (useState)
- Use for component-specific state
- Keep state close to usage
- Use useReducer for complex state logic

## Best Practices
- Normalize data structures
- Implement proper error handling
- Use TypeScript for type safety
- Test state interactions

## Domain-Specific State (` + fg.spec.Domain + `)
- Define domain entities clearly
- Implement proper relationships
- Handle business logic validation
- Use proper data normalization
`
}

// generateTestingGuidelinesDoc creates testing guidelines
func (fg *FrontendGenerator) generateTestingGuidelinesDoc() string {
	return `# Testing Guidelines

## Testing Strategy
- Unit tests for components
- Integration tests for user flows
- E2E tests for critical paths
- Visual regression tests

## Tools
- Jest for unit testing
- React Testing Library for component testing
- Playwright for E2E testing
- Storybook for component development

## Best Practices
- Test behavior, not implementation
- Use proper queries (getByRole, getByText)
- Mock external dependencies
- Test error states and edge cases

## Domain Testing (` + fg.spec.Domain + `)
- Test domain-specific business logic
- Validate data transformations
- Test security requirements
- Test accessibility compliance
`
}

// generateSecurityPracticesDoc creates security guidelines
func (fg *FrontendGenerator) generateSecurityPracticesDoc() string {
	return `# Frontend Security Practices

## Input Validation
- Sanitize all user inputs
- Use proper form validation
- Implement rate limiting
- Validate on both client and server

## Authentication
- Use secure token storage
- Implement proper logout
- Handle token expiration
- Use HTTPS only

## Data Protection
- Minimize sensitive data in frontend
- Use proper CORS configuration
- Implement CSP headers
- Avoid logging sensitive information

## Domain-Specific Security (` + fg.spec.Domain + `)
Based on ` + fg.spec.Domain + ` domain requirements:
- Implement domain-specific compliance
- Follow industry security standards
- Use proper data encryption
- Implement audit logging
`
}

// generateSetupGuideDoc creates setup instructions
func (fg *FrontendGenerator) generateSetupGuideDoc() string {
	return `# Frontend Setup Guide

## Initial Setup
` + "```" + `bash
# Create React + TypeScript project
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install

# Add required dependencies
npm install @tanstack/react-query zustand react-router-dom
npm install -D @types/node

# Setup development tools
npm install -D prettier eslint @typescript-eslint/parser
` + "```" + `

## Project Configuration
- Configure TypeScript strict mode
- Setup ESLint and Prettier
- Configure path aliases
- Setup environment variables

## Development Workflow
1. Create feature branch
2. Implement components with tests
3. Run linting and type checking
4. Submit for code review
5. Deploy to staging environment

## Build and Deployment
- Use production build optimizations
- Implement proper error monitoring
- Configure CDN for static assets
- Setup CI/CD pipeline
`
}

// generateProjectStructureDoc creates project structure guidelines
func (fg *FrontendGenerator) generateProjectStructureDoc() string {
	return `# Frontend Project Structure

` + "```" + `
apps/frontend/
├── public/
├── src/
│   ├── components/
│   │   ├── ui/           # Reusable UI components
│   │   ├── forms/        # Form components
│   │   └── domain/       # Domain-specific components
│   ├── pages/            # Route components
│   ├── hooks/            # Custom hooks
│   ├── services/         # API services
│   ├── stores/           # State management
│   ├── utils/            # Helper functions
│   ├── types/            # TypeScript definitions
│   └── styles/           # Global styles
├── tests/                # Test files
├── docs/                 # Component documentation
└── package.json
` + "```" + `

## File Naming Conventions
- Components: PascalCase (UserProfile.tsx)
- Hooks: camelCase with 'use' prefix (useUserData.ts)
- Utils: camelCase (formatDate.ts)
- Types: PascalCase (UserTypes.ts)

## Import Organization
1. React imports
2. Third-party libraries
3. Internal components
4. Relative imports
5. Type-only imports

## Domain-Specific Structure (` + fg.spec.Domain + `)
Add domain-specific folders:
- Business logic components
- Domain-specific types
- Industry-specific utilities
- Compliance-related components
`
}

// Domain-specific pattern generators
func (fg *FrontendGenerator) generateFintechPatternsDoc() string {
	return `# Fintech Frontend Patterns

## Security Requirements
- PCI DSS compliance for payment forms
- Secure handling of financial data
- Multi-factor authentication UI
- Session timeout handling

## Component Patterns
- Transaction tables with sorting/filtering
- Payment form components with validation
- Balance displays with formatting
- Account selection dropdowns

## Data Handling
- Decimal precision for monetary values
- Currency formatting
- Transaction categorization
- Real-time balance updates

## Regulatory Compliance
- Audit trail UI components
- Consent management interfaces
- Data export functionality
- Privacy control panels
`
}

func (fg *FrontendGenerator) generateHealthcarePatternsDoc() string {
	return `# Healthcare Frontend Patterns

## Compliance Requirements
- HIPAA compliance for patient data
- Accessibility standards (WCAG 2.1)
- Medical device integration
- Electronic health records UI

## Component Patterns
- Patient information forms
- Appointment scheduling interfaces
- Medical history displays
- Prescription management

## Data Security
- Patient data encryption
- Role-based access controls
- Audit logging interfaces
- Secure messaging components

## Medical Standards
- Medical terminology displays
- Clinical workflow interfaces
- Lab result presentations
- Treatment plan components
`
}

func (fg *FrontendGenerator) generateEcommercePatternsDoc() string {
	return `# E-commerce Frontend Patterns

## User Experience
- Product catalog with search/filter
- Shopping cart functionality
- Checkout flow optimization
- Order tracking interfaces

## Component Patterns
- Product cards and listings
- Shopping cart components
- Payment form integration
- Order history displays

## Performance
- Image optimization and lazy loading
- Search result pagination
- Real-time inventory updates
- Progressive web app features

## Business Features
- Recommendation engines UI
- Wishlist functionality
- Review and rating systems
- Promotional banner components
`
}