package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewFeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature [name]",
		Short: "Generate a new feature with domain expertise",
		Long: `Generates a new feature following domain best practices:

- Domain-specific code patterns
- Guard rails enforcement
- Complete CRUD operations
- Tests and documentation
- Integration with existing project structure

Examples:
  codekeeper feature payments --domain fintech
  codekeeper feature user-management --domain generic
  codekeeper feature inventory --domain ecommerce`,
		Args: cobra.ExactArgs(1),
		RunE: runFeature,
	}

	cmd.Flags().String("domain", "", "Domain expertise to apply")
	cmd.Flags().String("type", "", "Feature type (crud, api, service)")

	return cmd
}

func runFeature(cmd *cobra.Command, args []string) error {
	featureName := args[0]
	domain, _ := cmd.Flags().GetString("domain")
	featureType, _ := cmd.Flags().GetString("type")
	
	// Check if we're in a project directory
	if !isInProject() {
		return fmt.Errorf("not in a CodeKeeper project directory. Run this command from your project root")
	}
	
	// Auto-detect domain if not provided
	if domain == "" {
		domain = detectProjectDomain()
	}
	
	// Auto-detect feature type if not provided
	if featureType == "" {
		featureType = promptFeatureType()
	}
	
	color.Blue("🎯 Generating feature: %s", featureName)
	color.Yellow("📚 Domain: %s", domain)
	color.Cyan("🔧 Type: %s", featureType)
	
	generator := &FeatureGenerator{
		Name:   featureName,
		Domain: domain,
		Type:   featureType,
	}
	
	if err := generator.Generate(); err != nil {
		return fmt.Errorf("failed to generate feature: %w", err)
	}
	
	color.Green("✅ Feature '%s' generated successfully!", featureName)
	fmt.Printf("\n📋 Next steps:\n")
	fmt.Printf("  1. Review generated files in apps/backend/src/features/%s/\n", featureName)
	fmt.Printf("  2. Update routes in apps/backend/src/routes/index.js\n")
	fmt.Printf("  3. Run tests: npm test\n")
	fmt.Printf("  4. Check guard rails: codekeeper check\n")

	return nil
}

type FeatureGenerator struct {
	Name   string
	Domain string
	Type   string
}

func (g *FeatureGenerator) Generate() error {
	// Create feature directory structure
	featureDir := filepath.Join("apps", "backend", "src", "features", g.Name)
	if err := os.MkdirAll(featureDir, 0755); err != nil {
		return fmt.Errorf("failed to create feature directory: %w", err)
	}
	
	// Generate files based on feature type
	switch g.Type {
	case "crud":
		return g.generateCRUDFeature(featureDir)
	case "api":
		return g.generateAPIFeature(featureDir)
	case "service":
		return g.generateServiceFeature(featureDir)
	default:
		return g.generateCRUDFeature(featureDir) // Default to CRUD
	}
}

func (g *FeatureGenerator) generateCRUDFeature(baseDir string) error {
	files := map[string]string{
		"model.js":      g.generateModel(),
		"controller.js": g.generateController(),
		"routes.js":     g.generateRoutes(),
		"service.js":    g.generateService(),
		"validation.js": g.generateValidation(),
	}
	
	// Add test files
	testDir := filepath.Join(baseDir, "__tests__")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %w", err)
	}
	
	files[filepath.Join("__tests__", "controller.test.js")] = g.generateControllerTest()
	files[filepath.Join("__tests__", "service.test.js")] = g.generateServiceTest()
	
	// Write all files
	for relativePath, content := range files {
		fullPath := filepath.Join(baseDir, relativePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", relativePath, err)
		}
	}
	
	return nil
}

func (g *FeatureGenerator) generateAPIFeature(baseDir string) error {
	files := map[string]string{
		"endpoints.js":  g.generateEndpoints(),
		"middleware.js": g.generateMiddleware(),
		"validation.js": g.generateValidation(),
	}
	
	for relativePath, content := range files {
		fullPath := filepath.Join(baseDir, relativePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", relativePath, err)
		}
	}
	
	return nil
}

func (g *FeatureGenerator) generateServiceFeature(baseDir string) error {
	files := map[string]string{
		"service.js":    g.generateService(),
		"utils.js":      g.generateUtils(),
		"config.js":     g.generateConfig(),
	}
	
	for relativePath, content := range files {
		fullPath := filepath.Join(baseDir, relativePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", relativePath, err)
		}
	}
	
	return nil
}

func (g *FeatureGenerator) generateModel() string {
	entityName := strings.Title(g.Name)
	
	baseModel := fmt.Sprintf(`const { DataTypes } = require('sequelize');
const { sequelize } = require('../../../config/database');

const %s = sequelize.define('%s', {
  id: {
    type: DataTypes.UUID,
    defaultValue: DataTypes.UUIDV4,
    primaryKey: true
  },
  name: {
    type: DataTypes.STRING,
    allowNull: false,
    validate: {
      notEmpty: true,
      len: [1, 255]
    }
  },
  description: {
    type: DataTypes.TEXT,
    allowNull: true
  },
  status: {
    type: DataTypes.ENUM('active', 'inactive', 'pending'),
    defaultValue: 'active'
  },
  createdAt: {
    type: DataTypes.DATE,
    defaultValue: DataTypes.NOW
  },
  updatedAt: {
    type: DataTypes.DATE,
    defaultValue: DataTypes.NOW
  }
}, {
  tableName: '%s',
  timestamps: true,
  paranoid: true // Soft delete support
});

module.exports = %s;
`, entityName, entityName, strings.ToLower(g.Name)+"s", entityName)

	// Add domain-specific fields
	switch g.Domain {
	case "fintech":
		return strings.Replace(baseModel, `  status: {
    type: DataTypes.ENUM('active', 'inactive', 'pending'),
    defaultValue: 'active'
  },`, `  amount: {
    type: DataTypes.DECIMAL(15, 2),
    allowNull: false,
    validate: {
      min: 0
    }
  },
  currency: {
    type: DataTypes.STRING(3),
    defaultValue: 'USD',
    validate: {
      isIn: [['USD', 'EUR', 'GBP', 'JPY']]
    }
  },
  status: {
    type: DataTypes.ENUM('active', 'inactive', 'pending', 'processed'),
    defaultValue: 'pending'
  },`, 1)
	default:
		return baseModel
	}
}

func (g *FeatureGenerator) generateController() string {
	entityName := strings.Title(g.Name)
	
	return fmt.Sprintf(`const %sService = require('./service');
const { validationResult } = require('express-validator');

class %sController {
  async getAll(req, res) {
    try {
      const { page = 1, limit = 10, search } = req.query;
      const result = await %sService.getAll({ page, limit, search });
      
      res.json({
        success: true,
        data: result.data,
        pagination: result.pagination
      });
    } catch (error) {
      res.status(500).json({
        success: false,
        message: 'Failed to fetch %s',
        error: error.message
      });
    }
  }

  async getById(req, res) {
    try {
      const { id } = req.params;
      const %s = await %sService.getById(id);
      
      if (!%s) {
        return res.status(404).json({
          success: false,
          message: '%s not found'
        });
      }
      
      res.json({
        success: true,
        data: %s
      });
    } catch (error) {
      res.status(500).json({
        success: false,
        message: 'Failed to fetch %s',
        error: error.message
      });
    }
  }

  async create(req, res) {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          success: false,
          message: 'Validation failed',
          errors: errors.array()
        });
      }

      const %s = await %sService.create(req.body);
      
      res.status(201).json({
        success: true,
        data: %s,
        message: '%s created successfully'
      });
    } catch (error) {
      res.status(500).json({
        success: false,
        message: 'Failed to create %s',
        error: error.message
      });
    }
  }

  async update(req, res) {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          success: false,
          message: 'Validation failed',
          errors: errors.array()
        });
      }

      const { id } = req.params;
      const %s = await %sService.update(id, req.body);
      
      if (!%s) {
        return res.status(404).json({
          success: false,
          message: '%s not found'
        });
      }
      
      res.json({
        success: true,
        data: %s,
        message: '%s updated successfully'
      });
    } catch (error) {
      res.status(500).json({
        success: false,
        message: 'Failed to update %s',
        error: error.message
      });
    }
  }

  async delete(req, res) {
    try {
      const { id } = req.params;
      const deleted = await %sService.delete(id);
      
      if (!deleted) {
        return res.status(404).json({
          success: false,
          message: '%s not found'
        });
      }
      
      res.json({
        success: true,
        message: '%s deleted successfully'
      });
    } catch (error) {
      res.status(500).json({
        success: false,
        message: 'Failed to delete %s',
        error: error.message
      });
    }
  }
}

module.exports = new %sController();
`, entityName, entityName, entityName, g.Name, g.Name, entityName, g.Name, entityName, g.Name, g.Name, g.Name, entityName, g.Name, entityName, g.Name, g.Name, entityName, g.Name, entityName, g.Name, entityName, g.Name, entityName, g.Name, entityName, g.Name, entityName)
}

func (g *FeatureGenerator) generateRoutes() string {
	entityName := strings.Title(g.Name)
	
	return fmt.Sprintf(`const express = require('express');
const router = express.Router();
const %sController = require('./controller');
const %sValidation = require('./validation');
const auth = require('../../../middleware/auth');

// GET /%s - Get all %s with pagination
router.get('/', auth.optional, %sController.getAll);

// GET /%s/:id - Get %s by ID
router.get('/:id', auth.optional, %sController.getById);

// POST /%s - Create new %s
router.post('/', 
  auth.required, 
  %sValidation.create, 
  %sController.create
);

// PUT /%s/:id - Update %s
router.put('/:id', 
  auth.required, 
  %sValidation.update, 
  %sController.update
);

// DELETE /%s/:id - Delete %s
router.delete('/:id', 
  auth.required, 
  %sController.delete
);

module.exports = router;
`, entityName, g.Name, g.Name, g.Name, entityName, g.Name, g.Name, entityName, g.Name, g.Name, g.Name, entityName, g.Name, g.Name, g.Name, entityName, g.Name, g.Name, entityName)
}

func (g *FeatureGenerator) generateService() string {
	entityName := strings.Title(g.Name)
	
	baseService := fmt.Sprintf(`const %s = require('./model');
const { Op } = require('sequelize');

class %sService {
  async getAll(options = {}) {
    const { page = 1, limit = 10, search } = options;
    const offset = (page - 1) * limit;
    
    const whereClause = {};
    if (search) {
      whereClause[Op.or] = [
        { name: { [Op.iLike]: '%%' + search + '%%' } },
        { description: { [Op.iLike]: '%%' + search + '%%' } }
      ];
    }
    
    const { count, rows } = await %s.findAndCountAll({
      where: whereClause,
      limit: parseInt(limit),
      offset: parseInt(offset),
      order: [['createdAt', 'DESC']]
    });
    
    return {
      data: rows,
      pagination: {
        page: parseInt(page),
        limit: parseInt(limit),
        total: count,
        totalPages: Math.ceil(count / limit)
      }
    };
  }

  async getById(id) {
    return await %s.findByPk(id);
  }

  async create(data) {
    return await %s.create(data);
  }

  async update(id, data) {
    const %s = await %s.findByPk(id);
    if (!%s) return null;
    
    return await %s.update(data);
  }

  async delete(id) {
    const %s = await %s.findByPk(id);
    if (!%s) return false;
    
    await %s.destroy();
    return true;
  }
}

module.exports = new %sService();
`, entityName, entityName, entityName, entityName, entityName, g.Name, entityName, g.Name, g.Name, g.Name, entityName, g.Name, g.Name, entityName)

	// Add domain-specific methods
	switch g.Domain {
	case "fintech":
		return baseService + fmt.Sprintf(`

  // Fintech-specific methods
  async getByStatus(status) {
    return await %s.findAll({
      where: { status },
      order: [['createdAt', 'DESC']]
    });
  }

  async getTotalAmount(currency = 'USD') {
    const result = await %s.sum('amount', {
      where: { currency, status: 'processed' }
    });
    return result || 0;
  }
`, entityName, entityName)
	default:
		return baseService
	}
}

func (g *FeatureGenerator) generateValidation() string {
	return fmt.Sprintf(`const { body, param } = require('express-validator');

const createValidation = [
  body('name')
    .notEmpty()
    .withMessage('Name is required')
    .isLength({ min: 1, max: 255 })
    .withMessage('Name must be between 1 and 255 characters'),
  
  body('description')
    .optional()
    .isLength({ max: 1000 })
    .withMessage('Description must not exceed 1000 characters'),
  
  body('status')
    .optional()
    .isIn(['active', 'inactive', 'pending'])
    .withMessage('Status must be active, inactive, or pending')
];

const updateValidation = [
  param('id')
    .isUUID()
    .withMessage('Invalid ID format'),
  
  ...createValidation
];

module.exports = {
  create: createValidation,
  update: updateValidation
};
`)
}

func (g *FeatureGenerator) generateControllerTest() string {
	entityName := strings.Title(g.Name)
	
	return fmt.Sprintf(`const request = require('supertest');
const app = require('../../../app');
const %s = require('../model');

describe('%s Controller', () => {
  let test%s;

  beforeEach(async () => {
    test%s = await %s.create({
      name: 'Test %s',
      description: 'Test description'
    });
  });

  afterEach(async () => {
    await %s.destroy({ where: {}, truncate: true });
  });

  describe('GET /%s', () => {
    it('should return all %s', async () => {
      const response = await request(app)
        .get('/api/%s')
        .expect(200);

      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveLength(1);
      expect(response.body.data[0].name).toBe('Test %s');
    });
  });

  describe('GET /%s/:id', () => {
    it('should return %s by id', async () => {
      const response = await request(app)
        .get('/api/%s/' + test%s.id)
        .expect(200);

      expect(response.body.success).toBe(true);
      expect(response.body.data.name).toBe('Test %s');
    });

    it('should return 404 for non-existent %s', async () => {
      const response = await request(app)
        .get('/api/%s/00000000-0000-0000-0000-000000000000')
        .expect(404);

      expect(response.body.success).toBe(false);
    });
  });

  describe('POST /%s', () => {
    it('should create new %s', async () => {
      const new%s = {
        name: 'New %s',
        description: 'New description'
      };

      const response = await request(app)
        .post('/api/%s')
        .send(new%s)
        .expect(201);

      expect(response.body.success).toBe(true);
      expect(response.body.data.name).toBe('New %s');
    });

    it('should return 400 for invalid data', async () => {
      const response = await request(app)
        .post('/api/%s')
        .send({})
        .expect(400);

      expect(response.body.success).toBe(false);
    });
  });
});
`, entityName, entityName, entityName, entityName, entityName, entityName, entityName, g.Name, g.Name, g.Name, entityName, g.Name, g.Name, g.Name, entityName, entityName, g.Name, g.Name, g.Name, entityName, entityName, g.Name, entityName, entityName)
}

func (g *FeatureGenerator) generateServiceTest() string {
	entityName := strings.Title(g.Name)
	
	return fmt.Sprintf(`const %sService = require('../service');
const %s = require('../model');

describe('%s Service', () => {
  beforeEach(async () => {
    await %s.destroy({ where: {}, truncate: true });
  });

  describe('getAll', () => {
    it('should return paginated results', async () => {
      await %s.create({ name: 'Test 1' });
      await %s.create({ name: 'Test 2' });

      const result = await %sService.getAll({ page: 1, limit: 10 });

      expect(result.data).toHaveLength(2);
      expect(result.pagination.total).toBe(2);
    });

    it('should filter by search term', async () => {
      await %s.create({ name: 'Apple' });
      await %s.create({ name: 'Banana' });

      const result = await %sService.getAll({ search: 'Apple' });

      expect(result.data).toHaveLength(1);
      expect(result.data[0].name).toBe('Apple');
    });
  });

  describe('create', () => {
    it('should create new %s', async () => {
      const data = { name: 'New %s' };
      const %s = await %sService.create(data);

      expect(%s.name).toBe('New %s');
      expect(%s.id).toBeDefined();
    });
  });

  describe('update', () => {
    it('should update existing %s', async () => {
      const %s = await %s.create({ name: 'Original' });
      const updated = await %sService.update(%s.id, { name: 'Updated' });

      expect(updated.name).toBe('Updated');
    });

    it('should return null for non-existent %s', async () => {
      const result = await %sService.update('00000000-0000-0000-0000-000000000000', { name: 'Test' });
      expect(result).toBeNull();
    });
  });
});
`, entityName, entityName, entityName, entityName, entityName, entityName, entityName, entityName, entityName, entityName, g.Name, entityName, g.Name, entityName, g.Name, entityName, g.Name, g.Name, entityName, entityName, entityName, g.Name, g.Name, entityName)
}

func (g *FeatureGenerator) generateEndpoints() string {
	return fmt.Sprintf(`// API endpoints for %s feature
const endpoints = {
  GET_ALL: '/api/%s',
  GET_BY_ID: '/api/%s/:id',
  CREATE: '/api/%s',
  UPDATE: '/api/%s/:id',
  DELETE: '/api/%s/:id'
};

module.exports = endpoints;
`, g.Name, g.Name, g.Name, g.Name, g.Name, g.Name)
}

func (g *FeatureGenerator) generateMiddleware() string {
	return fmt.Sprintf(`// Middleware for %s feature
const rateLimit = require('express-rate-limit');

const %sLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
});

module.exports = {
  rateLimiter: %sLimiter
};
`, g.Name, g.Name, g.Name)
}

func (g *FeatureGenerator) generateUtils() string {
	return fmt.Sprintf(`// Utility functions for %s feature

class %sUtils {
  static formatResponse(data, message = 'Success') {
    return {
      success: true,
      message,
      data
    };
  }

  static formatError(error, message = 'Error occurred') {
    return {
      success: false,
      message,
      error: error.message
    };
  }
}

module.exports = %sUtils;
`, g.Name, strings.Title(g.Name), strings.Title(g.Name))
}

func (g *FeatureGenerator) generateConfig() string {
	return fmt.Sprintf(`// Configuration for %s feature

const config = {
  feature: '%s',
  version: '1.0.0',
  settings: {
    enableCaching: true,
    cacheTimeout: 300, // 5 minutes
    maxItems: 1000
  }
};

module.exports = config;
`, g.Name, g.Name)
}

func isInProject() bool {
	// Check for .codekeeper directory or apps directory
	if _, err := os.Stat(".codekeeper"); err == nil {
		return true
	}
	if _, err := os.Stat("apps"); err == nil {
		return true
	}
	return false
}

func detectProjectDomain() string {
	// Try to detect from .codekeeper config
	if data, err := os.ReadFile(".codekeeper/env.local"); err == nil {
		content := string(data)
		if strings.Contains(content, "CODEKEEPER_DOMAIN=fintech") {
			return "fintech"
		}
		if strings.Contains(content, "CODEKEEPER_DOMAIN=healthcare") {
			return "healthcare"
		}
		if strings.Contains(content, "CODEKEEPER_DOMAIN=ecommerce") {
			return "ecommerce"
		}
	}
	
	// Try to detect from package.json
	if data, err := os.ReadFile("package.json"); err == nil {
		content := strings.ToLower(string(data))
		if strings.Contains(content, "finance") || strings.Contains(content, "payment") {
			return "fintech"
		}
		if strings.Contains(content, "health") || strings.Contains(content, "medical") {
			return "healthcare"
		}
		if strings.Contains(content, "commerce") || strings.Contains(content, "shop") {
			return "ecommerce"
		}
	}
	
	return "general"
}

func promptFeatureType() string {
	types := []string{
		"crud - Complete CRUD operations with database model",
		"api - API endpoints and middleware only",
		"service - Business logic service with utilities",
	}
	
	var selected string
	prompt := &survey.Select{
		Message: "What type of feature would you like to generate?",
		Options: types,
		Default: "crud - Complete CRUD operations with database model",
	}
	survey.AskOne(prompt, &selected)
	
	// Extract type name (before " - ")
	if idx := strings.Index(selected, " - "); idx > 0 {
		return selected[:idx]
	}
	return "crud"
}