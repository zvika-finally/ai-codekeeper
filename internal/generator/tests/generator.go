package tests

import (
	"strings"
)

// TestGenerator handles test file generation
type TestGenerator struct {
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

// NewTestGenerator creates a new test generator
func NewTestGenerator(spec *ProjectSpec) *TestGenerator {
	return &TestGenerator{spec: spec}
}

// Generate creates test files and patterns
func (tg *TestGenerator) Generate() (map[string]string, error) {
	files := make(map[string]string)
	
	// Generate backend tests
	tg.generateBackendTests(files)
	
	// Generate frontend tests
	tg.generateFrontendTests(files)
	
	// Generate integration tests
	tg.generateIntegrationTests(files)
	
	// Generate test documentation
	tg.generateTestDocumentation(files)
	
	// Generate test configuration
	tg.generateTestConfiguration(files)
	
	return files, nil
}

// generateBackendTests creates backend test patterns
func (tg *TestGenerator) generateBackendTests(files map[string]string) {
	switch tg.spec.Backend {
	case "javascript", "node", "nodejs":
		tg.generateNodeJSTests(files)
	case "python":
		tg.generatePythonTests(files)
	case "go":
		tg.generateGoTests(files)
	default:
		tg.generateGenericBackendTests(files)
	}
}

// generateNodeJSTests creates Node.js/JavaScript test patterns
func (tg *TestGenerator) generateNodeJSTests(files map[string]string) {
	// Jest configuration
	files["apps/backend/jest.config.js"] = `module.exports = {
  testEnvironment: 'node',
  collectCoverageFrom: [
    'src/**/*.{js,ts}',
    '!src/**/*.test.{js,ts}',
    '!src/**/__tests__/**',
  ],
  coverageDirectory: 'coverage',
  coverageReporters: ['text', 'lcov', 'html'],
  setupFilesAfterEnv: ['<rootDir>/tests/setup.js'],
  testMatch: [
    '<rootDir>/src/**/__tests__/**/*.{js,ts}',
    '<rootDir>/src/**/*.{test,spec}.{js,ts}',
    '<rootDir>/tests/**/*.{test,spec}.{js,ts}',
  ],
  transform: {
    '^.+\\.(js|ts)$': 'babel-jest',
  },
  collectCoverageFrom: [
    'src/**/*.{js,ts}',
    '!src/index.js',
    '!**/node_modules/**',
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
};`

	// Test setup file
	files["apps/backend/tests/setup.js"] = `// Test setup for ` + tg.spec.Name + ` backend
require('dotenv').config({ path: '.env.test' });

// Global test configuration
global.console = {
  ...console,
  // Silence console during tests unless debugging
  log: process.env.DEBUG_TESTS ? console.log : jest.fn(),
  debug: process.env.DEBUG_TESTS ? console.debug : jest.fn(),
  info: process.env.DEBUG_TESTS ? console.info : jest.fn(),
  warn: console.warn,
  error: console.error,
};

// Domain-specific test setup
` + tg.getDomainSpecificTestSetup() + `

// Global test helpers
global.testHelpers = {
  createMock` + strings.Title(tg.spec.CoreEntity) + `() {
    return ` + tg.generateMockEntity() + `;
  },
  
  generateTestToken() {
    // Generate JWT token for testing
    return 'test-jwt-token';
  },
  
  setupTestDatabase() {
    // Database setup for tests
    return Promise.resolve();
  },
  
  cleanupTestDatabase() {
    // Database cleanup after tests
    return Promise.resolve();
  },
};`

	// API route tests
	files["apps/backend/tests/routes/" + strings.ToLower(tg.spec.CoreEntity) + ".test.js"] = tg.generateAPIRouteTests()

	// Service layer tests
	files["apps/backend/tests/services/" + strings.ToLower(tg.spec.CoreEntity) + ".test.js"] = tg.generateServiceTests()

	// Middleware tests
	files["apps/backend/tests/middleware/auth.test.js"] = tg.generateMiddlewareTests()

	// Integration tests
	files["apps/backend/tests/integration/api.test.js"] = tg.generateBackendIntegrationTests()

	// Test environment file
	files["apps/backend/.env.test"] = tg.generateTestEnvFile()
}

// generateFrontendTests creates frontend test patterns
func (tg *TestGenerator) generateFrontendTests(files map[string]string) {
	// Jest configuration for React
	files["apps/frontend/jest.config.js"] = `module.exports = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.ts'],
  moduleNameMapping: {
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  collectCoverageFrom: [
    'src/**/*.{ts,tsx}',
    '!src/**/*.d.ts',
    '!src/main.tsx',
    '!src/vite-env.d.ts',
  ],
  coverageThreshold: {
    global: {
      branches: 75,
      functions: 75,
      lines: 75,
      statements: 75,
    },
  },
  transform: {
    '^.+\\.(ts|tsx)$': 'ts-jest',
  },
};`

	// Test setup
	files["apps/frontend/src/setupTests.ts"] = `import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach, beforeAll, afterAll } from 'vitest';

// Cleanup after each test
afterEach(() => {
  cleanup();
});

// Mock environment variables
process.env.VITE_API_URL = 'http://localhost:8080';

// Mock modules that don't work well in test environment
vi.mock('react-router-dom', () => ({
  ...vi.importActual('react-router-dom'),
  useNavigate: () => vi.fn(),
  useLocation: () => ({ pathname: '/' }),
}));

// Domain-specific test setup for ` + tg.spec.Domain + `
` + tg.getDomainSpecificFrontendTestSetup()

	// Component tests
	files["apps/frontend/src/components/__tests__/" + strings.Title(tg.spec.CoreEntity) + "List.test.tsx"] = tg.generateComponentTests()

	// Hook tests
	files["apps/frontend/src/hooks/__tests__/use" + strings.Title(tg.spec.CoreEntity) + ".test.ts"] = tg.generateHookTests()

	// Service tests
	files["apps/frontend/src/services/__tests__/api.test.ts"] = tg.generateFrontendServiceTests()

	// Utils tests
	files["apps/frontend/src/utils/__tests__/validation.test.ts"] = tg.generateUtilsTests()
}

// generateIntegrationTests creates integration test patterns
func (tg *TestGenerator) generateIntegrationTests(files map[string]string) {
	// Playwright configuration
	files["tests/playwright.config.ts"] = `import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
  webServer: {
    command: 'npm run preview',
    port: 3000,
    cwd: './apps/frontend',
  },
});`

	// E2E tests
	files["tests/e2e/" + strings.ToLower(tg.spec.CoreEntity) + ".spec.ts"] = tg.generateE2ETests()

	// API tests
	files["tests/api/" + strings.ToLower(tg.spec.CoreEntity) + ".spec.ts"] = tg.generateAPITests()

	// Performance tests
	files["tests/performance/load.spec.ts"] = tg.generatePerformanceTests()
}

// generateTestDocumentation creates test documentation
func (tg *TestGenerator) generateTestDocumentation(files map[string]string) {
	files["docs/TESTING.md"] = `# Testing Guide for ` + tg.spec.Name + `

## Overview

This project uses a comprehensive testing strategy with multiple levels of testing:

- **Unit Tests**: Test individual components and functions
- **Integration Tests**: Test component interactions
- **API Tests**: Test backend API endpoints
- **E2E Tests**: Test complete user workflows
- **Performance Tests**: Test system performance under load

## Domain-Specific Testing (` + tg.spec.Domain + `)

` + tg.getDomainSpecificTestingDocs() + `

## Running Tests

### Backend Tests (Node.js)
` + "```bash" + `
cd apps/backend

# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run specific test file
npm test -- tests/routes/` + strings.ToLower(tg.spec.CoreEntity) + `.test.js

# Run in watch mode
npm run test:watch
` + "```" + `

### Frontend Tests (React)
` + "```bash" + `
cd apps/frontend

# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run specific test
npm test -- src/components/__tests__/` + strings.Title(tg.spec.CoreEntity) + `List.test.tsx

# Run in UI mode
npm run test:ui
` + "```" + `

### E2E Tests (Playwright)
` + "```bash" + `
cd tests

# Install browsers
npx playwright install

# Run E2E tests
npm run test:e2e

# Run specific test
npx playwright test e2e/` + strings.ToLower(tg.spec.CoreEntity) + `.spec.ts

# Run with UI
npx playwright test --ui
` + "```" + `

### API Tests
` + "```bash" + `
cd tests

# Run API tests
npm run test:api

# Run performance tests
npm run test:performance
` + "```" + `

## Test Structure

### Backend Tests
` + "```" + `
apps/backend/tests/
├── setup.js              # Test configuration
├── routes/               # API route tests
├── services/             # Business logic tests
├── middleware/           # Middleware tests
├── models/               # Data model tests
└── integration/          # Integration tests
` + "```" + `

### Frontend Tests
` + "```" + `
apps/frontend/src/
├── components/__tests__/  # Component tests
├── hooks/__tests__/      # Custom hook tests
├── services/__tests__/   # Service layer tests
├── utils/__tests__/      # Utility function tests
└── setupTests.ts         # Test configuration
` + "```" + `

### E2E Tests
` + "```" + `
tests/
├── e2e/                  # End-to-end tests
├── api/                  # API integration tests
├── performance/          # Performance tests
└── playwright.config.ts  # Playwright configuration
` + "```" + `

## Testing Best Practices

### Unit Tests
1. Test behavior, not implementation
2. Use descriptive test names
3. Follow AAA pattern (Arrange, Act, Assert)
4. Mock external dependencies
5. Aim for 80%+ code coverage

### Integration Tests
1. Test realistic user scenarios
2. Use test databases/environments
3. Clean up after each test
4. Test error conditions
5. Validate data persistence

### E2E Tests
1. Test critical user journeys
2. Use page object patterns
3. Keep tests independent
4. Use stable selectors
5. Test on multiple browsers

### Domain-Specific Testing
` + tg.getDomainSpecificTestingBestPractices() + `

## Continuous Integration

Tests run automatically on:
- Pull request creation
- Push to main branch
- Scheduled daily runs

See ` + "`.github/workflows/ci.yml`" + ` for CI configuration.

## Test Data Management

- Use factories for test data creation
- Implement database seeding for integration tests
- Use realistic but anonymized data
- Follow data privacy guidelines for ` + tg.spec.Domain + `

## Debugging Tests

1. Use ` + "`--verbose`" + ` flag for detailed output
2. Add ` + "`console.log`" + ` statements temporarily
3. Use debugger breakpoints
4. Check test setup and teardown
5. Verify test environment configuration
`
}

// generateTestConfiguration creates test configuration files
func (tg *TestGenerator) generateTestConfiguration(files map[string]string) {
	// GitHub Actions test workflow
	files[".github/workflows/test.yml"] = `name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  backend-tests:
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
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: apps/backend/package-lock.json
      
      - name: Install dependencies
        run: |
          cd apps/backend
          npm ci
      
      - name: Run tests
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost:5432/test_db
          NODE_ENV: test
        run: |
          cd apps/backend
          npm run test:coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: apps/backend/coverage/lcov.info

  frontend-tests:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: apps/frontend/package-lock.json
      
      - name: Install dependencies
        run: |
          cd apps/frontend
          npm ci
      
      - name: Run tests
        run: |
          cd apps/frontend
          npm run test:coverage

  e2e-tests:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd tests
          npm ci
          npx playwright install
      
      - name: Start services
        run: |
          docker-compose up -d
          sleep 30
      
      - name: Run E2E tests
        run: |
          cd tests
          npm run test:e2e
      
      - name: Upload test results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: tests/playwright-report/
`

	// Test package.json for E2E tests
	files["tests/package.json"] = `{
  "name": "` + tg.spec.Name + `-tests",
  "version": "1.0.0",
  "description": "Integration and E2E tests for ` + tg.spec.Name + `",
  "scripts": {
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:api": "jest api/",
    "test:performance": "jest performance/"
  },
  "devDependencies": {
    "@playwright/test": "^1.40.0",
    "jest": "^29.7.0",
    "supertest": "^6.3.3"
  }
}`
}

// Helper methods for generating specific test content
func (tg *TestGenerator) generateMockEntity() string {
	return `{
      id: 1,
      name: 'Test ` + strings.Title(tg.spec.CoreEntity) + `',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      // Add domain-specific fields for ` + tg.spec.Domain + `
    }`
}

func (tg *TestGenerator) getDomainSpecificTestSetup() string {
	switch tg.spec.Domain {
	case "fintech":
		return `// Fintech test helpers
global.fintechHelpers = {
  generateTestTransaction() {
    return {
      amount: '100.00',
      currency: 'USD',
      type: 'credit',
      reference: 'TEST-' + Date.now(),
    };
  },
  
  mockComplianceCheck() {
    return Promise.resolve({ approved: true, reason: 'test' });
  },
};`
	case "healthcare":
		return `// Healthcare test helpers
global.healthcareHelpers = {
  generateTestPatient() {
    return {
      firstName: 'Test',
      lastName: 'Patient',
      dateOfBirth: '1990-01-01',
      mrn: 'TEST-' + Date.now(),
    };
  },
  
  mockHIPAAValidation() {
    return Promise.resolve({ compliant: true });
  },
};`
	case "ecommerce":
		return `// E-commerce test helpers
global.ecommerceHelpers = {
  generateTestProduct() {
    return {
      name: 'Test Product',
      price: 29.99,
      sku: 'TEST-' + Date.now(),
      inventory: 100,
    };
  },
  
  generateTestOrder() {
    return {
      items: [{ productId: 1, quantity: 1, price: 29.99 }],
      total: 29.99,
      status: 'pending',
    };
  },
};`
	default:
		return `// Generic test helpers
global.genericHelpers = {
  generateTestData() {
    return { id: Date.now(), name: 'Test Data' };
  },
};`
	}
}

func (tg *TestGenerator) getDomainSpecificFrontendTestSetup() string {
	switch tg.spec.Domain {
	case "fintech":
		return `// Mock financial data validation
vi.mock('@/utils/decimal', () => ({
  formatCurrency: vi.fn((amount) => '$' + amount),
  validateAmount: vi.fn(() => true),
}));`
	case "healthcare":
		return `// Mock HIPAA compliance utilities
vi.mock('@/utils/hipaa', () => ({
  encryptPHI: vi.fn((data) => 'encrypted_' + data),
  validateConsent: vi.fn(() => true),
}));`
	case "ecommerce":
		return `// Mock e-commerce utilities
vi.mock('@/utils/cart', () => ({
  calculateTotal: vi.fn((items) => items.reduce((sum, item) => sum + item.price, 0)),
  validateInventory: vi.fn(() => true),
}));`
	default:
		return `// Generic frontend test setup`
	}
}

// Additional helper methods that were referenced but not implemented
func (tg *TestGenerator) generatePythonTests(files map[string]string) {
	// Python test implementation would go here
	files["apps/backend/pytest.ini"] = `[tool:pytest]
testpaths = tests
python_files = test_*.py
python_classes = Test*
python_functions = test_*
addopts = -v --cov=src --cov-report=html --cov-report=xml`
}

func (tg *TestGenerator) generateGoTests(files map[string]string) {
	// Go test implementation would go here
	entity := strings.ToLower(tg.spec.CoreEntity)
	files["pkg/"+entity+"/"+entity+"_test.go"] = `package ` + entity + `

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreate` + strings.Title(tg.spec.CoreEntity) + `(t *testing.T) {
	// Go test implementation
	assert.True(t, true, "Test should pass")
}`
}

func (tg *TestGenerator) generateGenericBackendTests(files map[string]string) {
	// Generic backend test patterns
	files["tests/README.md"] = `# Backend Tests

This directory contains backend tests for ` + tg.spec.Name + `.

## Running Tests

Refer to your backend framework's testing documentation for specific instructions.`
}

// generateAPIRouteTests creates API route test patterns
func (tg *TestGenerator) generateAPIRouteTests() string {
	entity := strings.ToLower(tg.spec.CoreEntity)
	entityTitle := strings.Title(tg.spec.CoreEntity)
	return `const request = require('supertest');
const app = require('../../src/app');
const { ` + entityTitle + ` } = require('../../src/models');

describe('` + entityTitle + ` Routes', () => {
  let server;
  
  beforeAll(async () => {
    server = app.listen(0);
    await testHelpers.setupTestDatabase();
  });
  
  afterAll(async () => {
    await testHelpers.cleanupTestDatabase();
    await server.close();
  });
  
  beforeEach(async () => {
    await ` + entityTitle + `.deleteMany({});
  });
  
  describe('GET /api/` + entity + `', () => {
    it('should return all ` + entity + `s', async () => {
      const mock` + entityTitle + ` = testHelpers.createMock` + entityTitle + `();
      await ` + entityTitle + `.create(mock` + entityTitle + `);
      
      const response = await request(server)
        .get('/api/` + entity + `')
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .expect(200);
      
      expect(response.body).toHaveLength(1);
      expect(response.body[0]).toMatchObject({
        name: mock` + entityTitle + `.name,
      });
    });
    
    it('should handle empty results', async () => {
      const response = await request(server)
        .get('/api/` + entity + `')
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .expect(200);
      
      expect(response.body).toHaveLength(0);
    });
  });
  
  describe('POST /api/` + entity + `', () => {
    it('should create a new ` + entity + `', async () => {
      const new` + entityTitle + ` = testHelpers.createMock` + entityTitle + `();
      delete new` + entityTitle + `.id;
      
      const response = await request(server)
        .post('/api/` + entity + `')
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .send(new` + entityTitle + `)
        .expect(201);
      
      expect(response.body).toMatchObject({
        name: new` + entityTitle + `.name,
      });
      expect(response.body.id).toBeDefined();
    });
    
    it('should validate required fields', async () => {
      const response = await request(server)
        .post('/api/` + entity + `')
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .send({})
        .expect(400);
      
      expect(response.body.error).toBeDefined();
    });
  });
  
  describe('GET /api/` + entity + `/:id', () => {
    it('should return a specific ` + entity + `', async () => {
      const mock` + entityTitle + ` = await ` + entityTitle + `.create(testHelpers.createMock` + entityTitle + `());
      
      const response = await request(server)
        .get(` + "`/api/" + entity + "/${mock" + entityTitle + "._id}`" + `)
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .expect(200);
      
      expect(response.body.name).toBe(mock` + entityTitle + `.name);
    });
    
    it('should return 404 for non-existent ` + entity + `', async () => {
      const response = await request(server)
        .get('/api/` + entity + `/507f1f77bcf86cd799439011')
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .expect(404);
    });
  });
  
  describe('PUT /api/` + entity + `/:id', () => {
    it('should update a ` + entity + `', async () => {
      const mock` + entityTitle + ` = await ` + entityTitle + `.create(testHelpers.createMock` + entityTitle + `());
      const updateData = { name: 'Updated Name' };
      
      const response = await request(server)
        .put(` + "`/api/" + entity + "/${mock" + entityTitle + "._id}`" + `)
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .send(updateData)
        .expect(200);
      
      expect(response.body.name).toBe(updateData.name);
    });
  });
  
  describe('DELETE /api/` + entity + `/:id', () => {
    it('should delete a ` + entity + `', async () => {
      const mock` + entityTitle + ` = await ` + entityTitle + `.create(testHelpers.createMock` + entityTitle + `());
      
      await request(server)
        .delete(` + "`/api/" + entity + "/${mock" + entityTitle + "._id}`" + `)
        .set('Authorization', 'Bearer ' + testHelpers.generateTestToken())
        .expect(204);
      
      const found = await ` + entityTitle + `.findById(mock` + entityTitle + `._id);
      expect(found).toBeNull();
    });
  });
});`
}

// generateServiceTests creates service layer test patterns
func (tg *TestGenerator) generateServiceTests() string {
	entity := strings.ToLower(tg.spec.CoreEntity)
	entityTitle := strings.Title(tg.spec.CoreEntity)
	return `const ` + entityTitle + `Service = require('../../src/services/` + entity + `Service');
const { ` + entityTitle + ` } = require('../../src/models');

jest.mock('../../src/models');

describe('` + entityTitle + `Service', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  describe('getAll` + entityTitle + `s', () => {
    it('should return all ` + entity + `s', async () => {
      const mock` + entityTitle + `s = [testHelpers.createMock` + entityTitle + `()];
      ` + entityTitle + `.find.mockResolvedValue(mock` + entityTitle + `s);
      
      const result = await ` + entityTitle + `Service.getAll` + entityTitle + `s();
      
      expect(result).toEqual(mock` + entityTitle + `s);
      expect(` + entityTitle + `.find).toHaveBeenCalledWith({});
    });
    
    it('should handle database errors', async () => {
      ` + entityTitle + `.find.mockRejectedValue(new Error('Database error'));
      
      await expect(` + entityTitle + `Service.getAll` + entityTitle + `s())
        .rejects.toThrow('Database error');
    });
  });
  
  describe('create` + entityTitle + `', () => {
    it('should create a new ` + entity + `', async () => {
      const new` + entityTitle + `Data = testHelpers.createMock` + entityTitle + `();
      const created` + entityTitle + ` = { ...new` + entityTitle + `Data, _id: '507f1f77bcf86cd799439011' };
      
      ` + entityTitle + `.create.mockResolvedValue(created` + entityTitle + `);
      
      const result = await ` + entityTitle + `Service.create` + entityTitle + `(new` + entityTitle + `Data);
      
      expect(result).toEqual(created` + entityTitle + `);
      expect(` + entityTitle + `.create).toHaveBeenCalledWith(new` + entityTitle + `Data);
    });
    
    it('should validate required fields', async () => {
      await expect(` + entityTitle + `Service.create` + entityTitle + `({}))
        .rejects.toThrow('Validation error');
    });
  });
  
  describe('get` + entityTitle + `ById', () => {
    it('should return ` + entity + ` by id', async () => {
      const mock` + entityTitle + ` = testHelpers.createMock` + entityTitle + `();
      ` + entityTitle + `.findById.mockResolvedValue(mock` + entityTitle + `);
      
      const result = await ` + entityTitle + `Service.get` + entityTitle + `ById('507f1f77bcf86cd799439011');
      
      expect(result).toEqual(mock` + entityTitle + `);
      expect(` + entityTitle + `.findById).toHaveBeenCalledWith('507f1f77bcf86cd799439011');
    });
    
    it('should return null for non-existent ` + entity + `', async () => {
      ` + entityTitle + `.findById.mockResolvedValue(null);
      
      const result = await ` + entityTitle + `Service.get` + entityTitle + `ById('507f1f77bcf86cd799439011');
      
      expect(result).toBeNull();
    });
  });
  
  describe('update` + entityTitle + `', () => {
    it('should update ` + entity + `', async () => {
      const updateData = { name: 'Updated Name' };
      const updated` + entityTitle + ` = { ...testHelpers.createMock` + entityTitle + `(), ...updateData };
      
      ` + entityTitle + `.findByIdAndUpdate.mockResolvedValue(updated` + entityTitle + `);
      
      const result = await ` + entityTitle + `Service.update` + entityTitle + `('507f1f77bcf86cd799439011', updateData);
      
      expect(result).toEqual(updated` + entityTitle + `);
      expect(` + entityTitle + `.findByIdAndUpdate).toHaveBeenCalledWith(
        '507f1f77bcf86cd799439011',
        updateData,
        { new: true }
      );
    });
  });
  
  describe('delete` + entityTitle + `', () => {
    it('should delete ` + entity + `', async () => {
      ` + entityTitle + `.findByIdAndDelete.mockResolvedValue(testHelpers.createMock` + entityTitle + `());
      
      const result = await ` + entityTitle + `Service.delete` + entityTitle + `('507f1f77bcf86cd799439011');
      
      expect(result).toBe(true);
      expect(` + entityTitle + `.findByIdAndDelete).toHaveBeenCalledWith('507f1f77bcf86cd799439011');
    });
    
    it('should return false for non-existent ` + entity + `', async () => {
      ` + entityTitle + `.findByIdAndDelete.mockResolvedValue(null);
      
      const result = await ` + entityTitle + `Service.delete` + entityTitle + `('507f1f77bcf86cd799439011');
      
      expect(result).toBe(false);
    });
  });
});`
}

// generateBackendIntegrationTests creates backend integration test patterns
func (tg *TestGenerator) generateBackendIntegrationTests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `const request = require('supertest');
const app = require('../../src/app');

describe('API Integration Tests', () => {
  beforeAll(async () => {
    await testHelpers.setupTestDatabase();
  });
  
  afterAll(async () => {
    await testHelpers.cleanupTestDatabase();
  });
  
  describe('` + entityTitle + ` CRUD Integration', () => {
    it('should complete full CRUD workflow', async () => {
      const token = testHelpers.generateTestToken();
      
      // CREATE
      const createResponse = await request(app)
        .post('/api/` + entity + `')
        .set('Authorization', 'Bearer ' + token)
        .send({ name: 'Integration Test ` + entityTitle + `' })
        .expect(201);
      
      const created` + entityTitle + ` = createResponse.body;
      expect(created` + entityTitle + `.id).toBeDefined();
      
      // READ
      const readResponse = await request(app)
        .get(` + "`/api/" + entity + "/${created" + entityTitle + ".id}`" + `)
        .set('Authorization', 'Bearer ' + token)
        .expect(200);
      
      expect(readResponse.body.name).toBe('Integration Test ` + entityTitle + `');
      
      // UPDATE
      const updateResponse = await request(app)
        .put(` + "`/api/" + entity + "/${created" + entityTitle + ".id}`" + `)
        .set('Authorization', 'Bearer ' + token)
        .send({ name: 'Updated Integration Test ` + entityTitle + `' })
        .expect(200);
      
      expect(updateResponse.body.name).toBe('Updated Integration Test ` + entityTitle + `');
      
      // DELETE
      await request(app)
        .delete(` + "`/api/" + entity + "/${created" + entityTitle + ".id}`" + `)
        .set('Authorization', 'Bearer ' + token)
        .expect(204);
      
      // VERIFY DELETE
      await request(app)
        .get(` + "`/api/" + entity + "/${created" + entityTitle + ".id}`" + `)
        .set('Authorization', 'Bearer ' + token)
        .expect(404);
    });
  });
});`
}

// generateMiddlewareTests creates middleware test patterns
func (tg *TestGenerator) generateMiddlewareTests() string {
	return `const request = require('supertest');
const app = require('../../src/app');
const jwt = require('jsonwebtoken');

describe('Auth Middleware', () => {
  describe('Authentication', () => {
    it('should allow access with valid token', async () => {
      const token = testHelpers.generateTestToken();
      
      const response = await request(app)
        .get('/api/` + strings.ToLower(tg.spec.CoreEntity) + `')
        .set('Authorization', 'Bearer ' + token)
        .expect(200);
    });
    
    it('should reject requests without token', async () => {
      await request(app)
        .get('/api/` + strings.ToLower(tg.spec.CoreEntity) + `')
        .expect(401);
    });
    
    it('should reject requests with invalid token', async () => {
      await request(app)
        .get('/api/` + strings.ToLower(tg.spec.CoreEntity) + `')
        .set('Authorization', 'Bearer invalid-token')
        .expect(401);
    });
  });
  
  describe('Authorization', () => {
    it('should allow admin access to protected routes', async () => {
      const adminToken = jwt.sign(
        { userId: 1, role: 'admin' },
        process.env.JWT_SECRET || 'test-secret'
      );
      
      await request(app)
        .delete('/api/` + strings.ToLower(tg.spec.CoreEntity) + `/1')
        .set('Authorization', 'Bearer ' + adminToken)
        .expect(204);
    });
    
    it('should deny user access to admin routes', async () => {
      const userToken = jwt.sign(
        { userId: 2, role: 'user' },
        process.env.JWT_SECRET || 'test-secret'
      );
      
      await request(app)
        .delete('/api/` + strings.ToLower(tg.spec.CoreEntity) + `/1')
        .set('Authorization', 'Bearer ' + userToken)
        .expect(403);
    });
  });
});`
}// generateComponentTests creates React component test patterns
func (tg *TestGenerator) generateComponentTests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';
import ` + entityTitle + `List from '../` + entityTitle + `List';
import { ` + entity + `Service } from '@/services/api';

// Mock the API service
vi.mock('@/services/api');

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {children}
      </BrowserRouter>
    </QueryClientProvider>
  );
};

describe('` + entityTitle + `List', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  
  it('renders loading state', () => {
    (` + entity + `Service.getAll as vi.Mock).mockReturnValue(
      new Promise(() => {}) // Never resolves
    );
    
    render(<` + entityTitle + `List />, { wrapper: createWrapper() });
    
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
  
  it('renders ` + entity + ` list', async () => {
    const mock` + entityTitle + `s = [
      { id: 1, name: 'Test ` + entityTitle + ` 1' },
      { id: 2, name: 'Test ` + entityTitle + ` 2' },
    ];
    
    (` + entity + `Service.getAll as vi.Mock).mockResolvedValue(mock` + entityTitle + `s);
    
    render(<` + entityTitle + `List />, { wrapper: createWrapper() });
    
    await waitFor(() => {
      expect(screen.getByText('Test ` + entityTitle + ` 1')).toBeInTheDocument();
      expect(screen.getByText('Test ` + entityTitle + ` 2')).toBeInTheDocument();
    });
  });
  
  it('renders error state', async () => {
    (` + entity + `Service.getAll as vi.Mock).mockRejectedValue(
      new Error('Failed to fetch')
    );
    
    render(<` + entityTitle + `List />, { wrapper: createWrapper() });
    
    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });
});`
}

// generateHookTests creates custom hook test patterns
func (tg *TestGenerator) generateHookTests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { use` + entityTitle + ` } from '../use` + entityTitle + `';
import { ` + entity + `Service } from '@/services/api';

vi.mock('@/services/api');

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('use` + entityTitle + `', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  
  it('fetches ` + entity + ` data successfully', async () => {
    const mock` + entityTitle + ` = { id: 1, name: 'Test ` + entityTitle + `' };
    (` + entity + `Service.getById as vi.Mock).mockResolvedValue(mock` + entityTitle + `);
    
    const { result } = renderHook(() => use` + entityTitle + `(1), {
      wrapper: createWrapper(),
    });
    
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });
    
    expect(result.current.data).toEqual(mock` + entityTitle + `);
  });
  
  it('handles fetch error', async () => {
    (` + entity + `Service.getById as vi.Mock).mockRejectedValue(
      new Error('Failed to fetch')
    );
    
    const { result } = renderHook(() => use` + entityTitle + `(1), {
      wrapper: createWrapper(),
    });
    
    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });
    
    expect(result.current.error).toEqual(expect.any(Error));
  });
});`
}

// generateFrontendServiceTests creates frontend service test patterns
func (tg *TestGenerator) generateFrontendServiceTests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `import { ` + entity + `Service } from '../api';

// Mock fetch globally
global.fetch = vi.fn();

describe('` + entityTitle + ` API Service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    (fetch as vi.Mock).mockClear();
  });
  
  describe('getAll', () => {
    it('fetches all ` + entity + `s successfully', async () => {
      const mock` + entityTitle + `s = [
        { id: 1, name: 'Test ` + entityTitle + ` 1' },
        { id: 2, name: 'Test ` + entityTitle + ` 2' },
      ];
      
      (fetch as vi.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mock` + entityTitle + `s),
      });
      
      const result = await ` + entity + `Service.getAll();
      
      expect(fetch).toHaveBeenCalledWith('/api/` + entity + `', {
        headers: { 'Content-Type': 'application/json' },
      });
      expect(result).toEqual(mock` + entityTitle + `s);
    });
    
    it('handles API error', async () => {
      (fetch as vi.Mock).mockResolvedValue({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
      });
      
      await expect(` + entity + `Service.getAll()).rejects.toThrow(
        'HTTP error! status: 500'
      );
    });
  });
});`
}

// generateUtilsTests creates utility function test patterns
func (tg *TestGenerator) generateUtilsTests() string {
	return `import { validateEmail, formatDate, sanitizeInput } from '../validation';

describe('Validation Utils', () => {
  describe('validateEmail', () => {
    it('validates correct email addresses', () => {
      expect(validateEmail('test@example.com')).toBe(true);
      expect(validateEmail('user.name@domain.co.uk')).toBe(true);
      expect(validateEmail('user+tag@example.org')).toBe(true);
    });
    
    it('rejects invalid email addresses', () => {
      expect(validateEmail('invalid-email')).toBe(false);
      expect(validateEmail('@domain.com')).toBe(false);
      expect(validateEmail('user@')).toBe(false);
      expect(validateEmail('')).toBe(false);
    });
  });
  
  describe('formatDate', () => {
    it('formats dates correctly', () => {
      const date = new Date('2024-01-15T10:30:00Z');
      expect(formatDate(date)).toBe('2024-01-15');
    });
    
    it('handles invalid dates', () => {
      expect(formatDate(null)).toBe('');
      expect(formatDate(undefined)).toBe('');
      expect(formatDate('invalid')).toBe('');
    });
  });
});`
}

// generateE2ETests creates end-to-end test patterns
func (tg *TestGenerator) generateE2ETests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `import { test, expect } from '@playwright/test';

test.describe('` + entityTitle + ` Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login or set up authentication
    await page.goto('/login');
    await page.fill('[data-testid="email"]', 'test@example.com');
    await page.fill('[data-testid="password"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await expect(page).toHaveURL('/dashboard');
  });
  
  test('should display ` + entity + ` list', async ({ page }) => {
    await page.goto('/` + entity + `');
    
    await expect(page).toHaveTitle(/` + entityTitle + ` Management/);
    await expect(page.locator('[data-testid="` + entity + `-list"]')).toBeVisible();
  });
  
  test('should create new ` + entity + `', async ({ page }) => {
    await page.goto('/` + entity + `');
    
    // Click create button
    await page.click('[data-testid="create-` + entity + `-button"]');
    await expect(page).toHaveURL('/` + entity + `/new');
    
    // Fill form
    await page.fill('[data-testid="` + entity + `-name"]', 'Test ` + entityTitle + `');
    await page.fill('[data-testid="` + entity + `-description"]', 'Test description');
    
    // Submit form
    await page.click('[data-testid="save-` + entity + `-button"]');
    
    // Verify success
    await expect(page.locator('[data-testid="success-message"]')).toBeVisible();
    await expect(page).toHaveURL('/` + entity + `');
    await expect(page.locator('text=Test ` + entityTitle + `')).toBeVisible();
  });
});`
}

// generateAPITests creates API integration test patterns
func (tg *TestGenerator) generateAPITests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `import { test, expect } from '@playwright/test';

const API_BASE_URL = 'http://localhost:8080/api';

test.describe('` + entityTitle + ` API', () => {
  let authToken: string;
  
  test.beforeAll(async ({ request }) => {
    // Get authentication token
    const response = await request.post(API_BASE_URL + '/auth/login', {
      data: {
        email: 'test@example.com',
        password: 'password123',
      },
    });
    
    const data = await response.json();
    authToken = data.token;
  });
  
  test('GET /` + entity + ` - should return ` + entity + ` list', async ({ request }) => {
    const response = await request.get(API_BASE_URL + '/` + entity + `', {
      headers: {
        Authorization: 'Bearer ' + authToken,
      },
    });
    
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    expect(Array.isArray(data)).toBeTruthy();
  });
  
  test('POST /` + entity + ` - should create new ` + entity + `', async ({ request }) => {
    const new` + entityTitle + ` = {
      name: 'Test ` + entityTitle + `',
      description: 'Test description',
    };
    
    const response = await request.post(API_BASE_URL + '/` + entity + `', {
      headers: {
        Authorization: 'Bearer ' + authToken,
      },
      data: new` + entityTitle + `,
    });
    
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(201);
    
    const data = await response.json();
    expect(data.name).toBe(new` + entityTitle + `.name);
    expect(data.id).toBeDefined();
  });
});`
}

// generatePerformanceTests creates performance test patterns
func (tg *TestGenerator) generatePerformanceTests() string {
	entityTitle := strings.Title(tg.spec.CoreEntity)
	entity := strings.ToLower(tg.spec.CoreEntity)
	return `import { test, expect } from '@playwright/test';

test.describe('Performance Tests', () => {
  test('` + entityTitle + ` list page load performance', async ({ page }) => {
    // Start timing
    const startTime = Date.now();
    
    await page.goto('/` + entity + `');
    
    // Wait for content to load
    await page.waitForSelector('[data-testid="` + entity + `-list"]');
    
    const endTime = Date.now();
    const loadTime = endTime - startTime;
    
    // Page should load within 3 seconds
    expect(loadTime).toBeLessThan(3000);
  });
  
  test('API response time', async ({ request }) => {
    const startTime = Date.now();
    
    const response = await request.get('http://localhost:8080/api/` + entity + `');
    
    const endTime = Date.now();
    const responseTime = endTime - startTime;
    
    expect(response.ok()).toBeTruthy();
    // API should respond within 1 second
    expect(responseTime).toBeLessThan(1000);
  });
});`
}

// generateTestEnvFile creates test environment configuration
func (tg *TestGenerator) generateTestEnvFile() string {
	return `# Test Environment Configuration
NODE_ENV=test
PORT=8081

# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/` + strings.ToLower(tg.spec.Name) + `_test
REDIS_URL=redis://localhost:6379/1

# JWT
JWT_SECRET=test-jwt-secret-key
JWT_EXPIRES_IN=1h

# External APIs (use test/mock endpoints)
EXTERNAL_API_URL=http://localhost:3001/mock-api
EXTERNAL_API_KEY=test-api-key

# Feature flags for testing
ENABLE_DEBUG_LOGS=false
ENABLE_TEST_HELPERS=true

# Domain-specific test configuration
` + tg.getDomainSpecificTestEnv()
}

// getDomainSpecificTestingDocs returns domain-specific testing documentation
func (tg *TestGenerator) getDomainSpecificTestingDocs() string {
	switch tg.spec.Domain {
	case "fintech":
		return `### Financial Data Testing

**Compliance Testing**
- All financial calculations must be tested with decimal precision
- Test compliance validation workflows
- Verify audit trail generation
- Test regulatory reporting accuracy

**Security Testing**
- Validate PCI DSS compliance in payment flows
- Test data encryption for sensitive financial information
- Verify access controls for financial operations
- Test fraud detection mechanisms`

	case "healthcare":
		return `### Healthcare Data Testing

**HIPAA Compliance Testing**
- Verify PHI (Protected Health Information) encryption
- Test access logging and audit trails
- Validate consent management workflows
- Test data anonymization for reporting

**Clinical Data Testing**
- Use realistic but anonymized patient data
- Test medical record integrity
- Validate clinical decision support accuracy
- Test integration with health information exchanges`

	case "ecommerce":
		return `### E-commerce Testing

**Transaction Testing**
- Test payment processing workflows
- Validate inventory management accuracy
- Test order fulfillment processes
- Verify refund and return procedures

**Performance Testing**
- Test high-traffic scenarios (Black Friday, sales events)
- Validate cart abandonment recovery
- Test search and filtering performance
- Verify mobile responsiveness`

	default:
		return `### Domain-Specific Testing

- Implement tests specific to your ` + tg.spec.Domain + ` domain
- Focus on business logic validation
- Test integration points with external services
- Validate data integrity and security requirements`
	}
}

// getDomainSpecificTestingBestPractices returns domain-specific testing best practices
func (tg *TestGenerator) getDomainSpecificTestingBestPractices() string {
	switch tg.spec.Domain {
	case "fintech":
		return `
#### Financial Services Best Practices

1. **Decimal Precision**: Always use decimal libraries for monetary calculations
2. **Audit Trails**: Test that all financial operations are logged
3. **Compliance**: Validate regulatory requirements (PCI DSS, SOX, etc.)
4. **Security**: Test encryption of sensitive financial data
5. **Reconciliation**: Test that all transactions can be reconciled
6. **Fraud Detection**: Test fraud prevention mechanisms
7. **Rate Limiting**: Test API rate limits for financial operations`

	case "healthcare":
		return `
#### Healthcare Best Practices

1. **Data Privacy**: Ensure all PHI is properly protected in tests
2. **Consent Management**: Test patient consent workflows
3. **Clinical Accuracy**: Validate medical calculations and recommendations
4. **Interoperability**: Test HL7 FHIR compliance
5. **Audit Logging**: Test comprehensive audit trails
6. **Emergency Access**: Test break-glass access procedures
7. **Data Retention**: Test data lifecycle management`

	case "ecommerce":
		return `
#### E-commerce Best Practices

1. **Inventory Accuracy**: Test stock management and overselling prevention
2. **Cart Persistence**: Test shopping cart behavior across sessions
3. **Payment Security**: Test PCI compliance in payment flows
4. **Order Integrity**: Test order processing from cart to fulfillment
5. **Performance**: Test under high load conditions
6. **Mobile Experience**: Test responsive design and mobile payments
7. **Abandoned Cart**: Test cart recovery workflows`

	default:
		return `
#### General Best Practices

1. **Business Logic**: Focus on domain-specific business rules
2. **Data Integrity**: Test data validation and constraints
3. **Security**: Test authentication and authorization
4. **Performance**: Test under realistic load conditions
5. **Integration**: Test external service dependencies
6. **Error Handling**: Test error scenarios and recovery
7. **Monitoring**: Test logging and alerting mechanisms`
	}
}

// getDomainSpecificTestEnv returns domain-specific test environment variables
func (tg *TestGenerator) getDomainSpecificTestEnv() string {
	switch tg.spec.Domain {
	case "fintech":
		return `# Fintech-specific test configuration
PAYMENT_GATEWAY_URL=http://localhost:3001/mock-payment
COMPLIANCE_API_URL=http://localhost:3001/mock-compliance
FRAUD_DETECTION_API_URL=http://localhost:3001/mock-fraud
ENABLE_TRANSACTION_LOGGING=true
DECIMAL_PRECISION=2`

	case "healthcare":
		return `# Healthcare-specific test configuration
HIPAA_ENCRYPTION_KEY=test-hipaa-encryption-key
EHR_INTEGRATION_URL=http://localhost:3001/mock-ehr
CLINICAL_API_URL=http://localhost:3001/mock-clinical
PHI_AUDIT_ENABLED=true
CONSENT_MANAGEMENT_URL=http://localhost:3001/mock-consent`

	case "ecommerce":
		return `# E-commerce-specific test configuration
PAYMENT_PROCESSOR_URL=http://localhost:3001/mock-payments
INVENTORY_API_URL=http://localhost:3001/mock-inventory
SHIPPING_API_URL=http://localhost:3001/mock-shipping
TAX_CALCULATION_URL=http://localhost:3001/mock-tax
ENABLE_CART_PERSISTENCE=true`

	default:
		return `# Add domain-specific test environment variables here`
	}
}
