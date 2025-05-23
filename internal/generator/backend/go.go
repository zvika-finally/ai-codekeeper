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
	return `package main

// STARTER: ` + bg.spec.Name + ` - ` + bg.spec.Domain + ` Go Backend
// IMPLEMENT: Follow docs/backend/ guidelines for implementation

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	
	// API routes for ` + bg.spec.CoreEntity + `
	api := r.Group("/api/v1")
	{
		// Authentication
		api.POST("/login", handleLogin)
		api.POST("/register", handleRegister)
		
		// CRUD for ` + bg.spec.CoreEntity + `
		api.GET("/` + strings.ToLower(bg.spec.CoreEntity) + `s", getResourcesHandler)
		api.POST("/` + strings.ToLower(bg.spec.CoreEntity) + `s", createResourceHandler)
		api.GET("/` + strings.ToLower(bg.spec.CoreEntity) + `s/:id", getResourceHandler)
		api.PUT("/` + strings.ToLower(bg.spec.CoreEntity) + `s/:id", updateResourceHandler)
		api.DELETE("/` + strings.ToLower(bg.spec.CoreEntity) + `s/:id", deleteResourceHandler)
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)
}

// Handler functions - IMPLEMENT according to your business logic

// Authentication handlers
func handleLogin(c *gin.Context) {
	// IMPLEMENT: Add JWT authentication, validate credentials against database
	// SECURITY: Hash passwords, implement rate limiting, add CORS
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - implement authentication"})
}

func handleRegister(c *gin.Context) {
	// IMPLEMENT: Validate input, hash password, save to database
	// SECURITY: Check password strength, prevent duplicate emails
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - implement user registration"})
}

// %s CRUD handlers
func get%sHandler(c *gin.Context) {
	// IMPLEMENT: Query database, add pagination, filtering, sorting
	// PERFORMANCE: Use database indexes, implement caching
	c.JSON(http.StatusOK, gin.H{
		"message": "Get all %s",
		"data":    []%s{},
	})
}

func create%sHandler(c *gin.Context) {
	// IMPLEMENT: Validate input, save to database, return created entity
	// VALIDATION: Check required fields, business rules
	c.JSON(http.StatusCreated, gin.H{"message": "Create %s endpoint"})
}

func get%sHandler(c *gin.Context) {
	id := c.Param("id")
	// IMPLEMENT: Query database by ID, handle not found
	// VALIDATION: Validate ID format, check permissions
	c.JSON(http.StatusOK, gin.H{
		"message": "Get %s by ID",
		"id":      id,
	})
}

func update%sHandler(c *gin.Context) {
	id := c.Param("id")
	// IMPLEMENT: Validate input, update database, return updated entity
	// VALIDATION: Check exists, validate partial updates
	c.JSON(http.StatusOK, gin.H{
		"message": "Update %s",
		"id":      id,
	})
}

func deleteResourceHandler(c *gin.Context) {
	id := c.Param("id")
	// IMPLEMENT: Soft delete, check dependencies, log action
	// VALIDATION: Check exists, verify permissions
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete ` + bg.spec.CoreEntity + `",
		"id":      id,
	})
}
`
}
