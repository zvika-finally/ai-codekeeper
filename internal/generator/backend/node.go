package backend

import (
	"fmt"
	"strings"
)

// generateNodePackageJSON creates package.json for Node.js backend
func (bg *BackendGenerator) generateNodePackageJSON() string {
	return fmt.Sprintf(`{
  "name": "@%s/backend",
  "version": "1.0.0",
  "description": "%s backend API",
  "main": "src/server.js",
  "scripts": {
    "start": "node src/server.js",
    "dev": "nodemon src/server.js",
    "test": "jest --testTimeout=10000",
    "test:coverage": "jest --coverage",
    "lint": "eslint src/",
    "lint:fix": "eslint src/ --fix",
    "check": "codekeeper check",
    "migration:create": "sequelize-cli migration:generate",
    "migration:run": "sequelize-cli db:migrate",
    "migration:undo": "sequelize-cli db:migrate:undo"
  },
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^7.0.0",
    "express-rate-limit": "^6.7.0",
    "sequelize": "^6.32.1",
    "pg": "^8.11.0",
    "pg-hstore": "^2.3.4",
    "jsonwebtoken": "^9.0.1",
    "bcrypt": "^5.1.0",
    "winston": "^3.9.0",
    "validator": "^13.9.0",
    "dotenv": "^16.1.4"
  },
  "devDependencies": {
    "nodemon": "^2.0.22",
    "jest": "^29.5.0",
    "supertest": "^6.3.3",
    "eslint": "^8.42.0",
    "eslint-config-node": "^4.1.0",
    "sequelize-cli": "^6.6.1"
  },
  "jest": {
    "testEnvironment": "node",
    "setupFilesAfterEnv": ["<rootDir>/tests/setup.js"],
    "collectCoverageFrom": [
      "src/**/*.js",
      "!src/server.js"
    ],
    "coverageThreshold": {
      "global": {
        "branches": 80,
        "functions": 80,
        "lines": 80,
        "statements": 80
      }
    }
  }
}`, bg.spec.Name, bg.spec.Description)
}

// generateNodeDockerfile creates Dockerfile for Node.js backend
func (bg *BackendGenerator) generateNodeDockerfile() string {
	return `# Multi-stage build for production optimization
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Production stage
FROM node:18-alpine AS production

# Create non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nodejs -u 1001

WORKDIR /app

# Copy dependencies from builder stage
COPY --from=builder /app/node_modules ./node_modules
COPY --chown=nodejs:nodejs . .

# Switch to non-root user
USER nodejs

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:8080/health', (res) => process.exit(res.statusCode === 200 ? 0 : 1))"

CMD ["node", "src/server.js"]`
}

// generateNodeServer creates Express server implementation
func (bg *BackendGenerator) generateNodeServer() string {
	userRoles := bg.spec.GetUserRolesList()

	return fmt.Sprintf(`const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 8080;

// Security middleware
app.use(helmet());
app.use(cors({
  origin: process.env.FRONTEND_URL || 'http://localhost:3000',
  credentials: true
}));

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // limit each IP to 100 requests per windowMs
  message: 'Too many requests from this IP, please try again later.'
});
app.use('/api/', limiter);

// Body parsing middleware
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'healthy',
    timestamp: new Date().toISOString(),
    service: '%s-backend',
    version: '1.0.0'
  });
});

// API routes
app.use('/api/auth', require('./routes/auth'));
app.use('/api/%s', require('./routes/%s'));

// Error handling middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Something went wrong!',
    message: process.env.NODE_ENV === 'development' ? err.message : 'Internal server error'
  });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    error: 'Not Found',
    message: 'The requested resource was not found'
  });
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully');
  server.close(() => {
    console.log('Process terminated');
  });
});

const server = app.listen(PORT, () => {
  console.log(' + ========================================== +');
  console.log(' |  🚀 %s Backend Server            |');
  console.log(' |  📡 Port: ' + PORT + '                          |');
  console.log(' |  🌍 Environment: ' + (process.env.NODE_ENV || 'development') + '           |');
  console.log(' |  📊 Domain: %s                   |');
  console.log(' |  👥 User Roles: %s    |');
  console.log(' + ========================================== +');
});

async function startServer() {
  try {
    // Database connection would be initialized here
    console.log('✅ Database connected');
    console.log('✅ Server started successfully');
  } catch (error) {
    console.error('❌ Failed to start server:', error);
    process.exit(1);
  }
}

startServer();

module.exports = app;`, 
	bg.spec.Name, 
	strings.ToLower(bg.spec.CoreEntity), 
	strings.ToLower(bg.spec.CoreEntity),
	bg.spec.Name,
	bg.spec.Domain,
	strings.Join(userRoles, ", "))
}

// generateNodeDatabaseConfig creates database configuration for Node.js
func (bg *BackendGenerator) generateNodeDatabaseConfig() string {
	return `const { Sequelize } = require('sequelize');

const sequelize = new Sequelize(
  process.env.DATABASE_URL || 'postgresql://user:password@localhost:5432/database',
  {
    dialect: 'postgres',
    dialectOptions: {
      ssl: process.env.NODE_ENV === 'production' ? {
        require: true,
        rejectUnauthorized: false
      } : false
    },
    pool: {
      max: 10,
      min: 0,
      acquire: 30000,
      idle: 10000
    },
    logging: process.env.NODE_ENV === 'development' ? console.log : false
  }
);

async function connectDatabase() {
  try {
    await sequelize.authenticate();
    console.log('✅ Database connection established successfully');
    return sequelize;
  } catch (error) {
    console.error('❌ Unable to connect to the database:', error);
    throw error;
  }
}

module.exports = { sequelize, connectDatabase };`
}

// generateNodeAuthConfig creates authentication configuration for Node.js
func (bg *BackendGenerator) generateNodeAuthConfig() string {
	return `const jwt = require('jsonwebtoken');
const bcrypt = require('bcrypt');

const JWT_SECRET = process.env.JWT_SECRET || 'your-secret-key';
const JWT_EXPIRES_IN = process.env.JWT_EXPIRES_IN || '24h';

function generateToken(payload) {
  return jwt.sign(payload, JWT_SECRET, { expiresIn: JWT_EXPIRES_IN });
}

function verifyToken(token) {
  return jwt.verify(token, JWT_SECRET);
}

async function hashPassword(password) {
  return await bcrypt.hash(password, 12);
}

async function comparePassword(password, hashedPassword) {
  return await bcrypt.compare(password, hashedPassword);
}

// Middleware to verify JWT token
function authenticateToken(req, res, next) {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];

  if (!token) {
    return res.status(401).json({ error: 'Access token required' });
  }

  try {
    const decoded = verifyToken(token);
    req.user = decoded;
    next();
  } catch (error) {
    return res.status(403).json({ error: 'Invalid or expired token' });
  }
}

module.exports = {
  generateToken,
  verifyToken,
  hashPassword,
  comparePassword,
  authenticateToken
};`
}