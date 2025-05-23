package backend

import (
	"fmt"
	"strings"
)

// generateGoMod creates go.mod file for Go backend
func (bg *BackendGenerator) generateGoMod() string {
	return fmt.Sprintf(`module %s-backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/lib/pq v1.10.9
	github.com/golang-jwt/jwt/v5 v5.1.0
	github.com/joho/godotenv v1.4.0
	github.com/golang-migrate/migrate/v4 v4.16.2
	golang.org/x/crypto v0.15.0
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)`, bg.spec.Name)
}

// generateGoDockerfile creates Dockerfile for Go backend
func (bg *BackendGenerator) generateGoDockerfile() string {
	return fmt.Sprintf(`# Multi-stage build for production optimization
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for Go modules
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/

# Production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]`)
}

// generateGoMain creates main.go for Go backend
func (bg *BackendGenerator) generateGoMain() string {
	userRoles := bg.spec.GetUserRolesList()
	rolesStr := strings.Join(userRoles, ", ")
	if rolesStr == "" {
		rolesStr = "user"
	}

	return fmt.Sprintf(`package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// %s represents the core entity
type %s struct {
	ID        uint      ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	// Add your domain-specific fields here
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "%s-backend",
			"version":   "1.0.0",
			"domain":    "%s",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Authentication routes
		api.POST("/auth/login", handleLogin)
		api.POST("/auth/register", handleRegister)
		
		// %s routes
		api.GET("/%s", get%sHandler)
		api.POST("/%s", create%sHandler)
		api.GET("/%s/:id", get%sHandler)
		api.PUT("/%s/:id", update%sHandler)
		api.DELETE("/%s/:id", delete%sHandler)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 %s Backend Server")
	log.Printf("📡 Port: %%s", port)
	log.Printf("🌍 Environment: %%s", os.Getenv("GIN_MODE"))
	log.Printf("📊 Domain: %s")
	log.Printf("👥 User Roles: %s")
	log.Printf("✅ Server started successfully")

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Authentication handlers
func handleLogin(c *gin.Context) {
	// TODO: Implement login logic
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - implement authentication"})
}

func handleRegister(c *gin.Context) {
	// TODO: Implement registration logic
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - implement user registration"})
}

// %s CRUD handlers
func get%sHandler(c *gin.Context) {
	// TODO: Implement get all %s logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Get all %s",
		"data":    []%s{},
	})
}

func create%sHandler(c *gin.Context) {
	// TODO: Implement create %s logic
	c.JSON(http.StatusCreated, gin.H{"message": "Create %s endpoint"})
}

func get%sHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement get single %s logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Get %s by ID",
		"id":      id,
	})
}

func update%sHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement update %s logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Update %s",
		"id":      id,
	})
}

func delete%sHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement delete %s logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete %s",
		"id":      id,
	})
}
`, 
	bg.spec.CoreEntity, bg.spec.CoreEntity,
	bg.spec.Name, bg.spec.Domain,
	bg.spec.CoreEntity, strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
	strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
	strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
	strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
	strings.ToLower(bg.spec.CoreEntity), bg.spec.CoreEntity,
	bg.spec.Name, bg.spec.Domain, rolesStr,
	bg.spec.CoreEntity, bg.spec.CoreEntity, bg.spec.CoreEntity,
	bg.spec.CoreEntity, bg.spec.CoreEntity,
	bg.spec.CoreEntity, bg.spec.CoreEntity, bg.spec.CoreEntity,
	bg.spec.CoreEntity, bg.spec.CoreEntity, bg.spec.CoreEntity,
	bg.spec.CoreEntity, bg.spec.CoreEntity, bg.spec.CoreEntity,
	bg.spec.CoreEntity, bg.spec.CoreEntity, bg.spec.CoreEntity)
}

// generateGolangCIConfig creates .golangci.yml configuration
func (bg *BackendGenerator) generateGolangCIConfig() string {
	return `run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - dupl
    - goconst
    - gocyclo
    - gosec
    - depguard
    - prealloc
    - exportloopref
    - nolintlint

linters-settings:
  gocyclo:
    min-complexity: 15
  
  goconst:
    min-len: 2
    min-occurrences: 2
  
  dupl:
    threshold: 100
  
  gosec:
    severity: "low"
    confidence: "low"
    excludes:
      - G104 # Errors unhandled (we handle this with errcheck)
  
  depguard:
    list-type: blacklist
    packages:
      - github.com/pkg/errors
    packages-with-error-message:
      - github.com/pkg/errors: "use standard library errors package"

issues:
  exclude-use-default: false
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - dupl
    - path: main\.go
      linters:
        - gocyclo

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true`
}