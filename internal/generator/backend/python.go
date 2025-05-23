package backend

import (
	"fmt"
	"strings"
)

// generatePythonRequirements creates requirements.txt for Python backend
func (bg *BackendGenerator) generatePythonRequirements() string {
	framework := bg.spec.GetBackendFramework()
	requirements := `# Core framework
fastapi==0.104.1
uvicorn[standard]==0.24.0

# Database
sqlalchemy==2.0.23
alembic==1.12.1
psycopg2-binary==2.9.9

# Authentication & Security
python-jose[cryptography]==3.3.0
passlib[bcrypt]==1.7.4
python-multipart==0.0.6

# Environment & Configuration
python-dotenv==1.0.0
pydantic==2.5.1
pydantic-settings==2.1.0

# HTTP & API
httpx==0.25.2
requests==2.31.0

# Development & Testing
pytest==7.4.3
pytest-asyncio==0.21.1
pytest-cov==4.1.0
black==23.11.0
isort==5.12.0
flake8==6.1.0
mypy==1.7.1

# Production
gunicorn==21.2.0
redis==5.0.1`

	if framework == "django" {
		requirements = `# Django framework
Django==4.2.7
djangorestframework==3.14.0
django-cors-headers==4.3.1
django-environ==0.11.2

# Database
psycopg2-binary==2.9.9
redis==5.0.1

# Authentication & Security
djangorestframework-simplejwt==5.3.0
django-oauth-toolkit==1.7.1

# Testing & Development
pytest-django==4.7.0
pytest-cov==4.1.0
factory-boy==3.3.0
black==23.11.0
isort==5.12.0
flake8==6.1.0
mypy==1.7.1

# Production
gunicorn==21.2.0
whitenoise==6.6.0`
	}

	return requirements
}

// generatePyprojectToml creates pyproject.toml for Python backend
func (bg *BackendGenerator) generatePyprojectToml() string {
	return fmt.Sprintf(`[build-system]
requires = ["setuptools>=45", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "%s-backend"
version = "1.0.0"
description = "%s backend API"
authors = [
    {name = "Generated with AI CodeKeeper", email = "noreply@example.com"}
]
license = {text = "MIT"}
readme = "README.md"
requires-python = ">=3.9"
classifiers = [
    "Development Status :: 4 - Beta",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
]

[tool.black]
line-length = 88
target-version = ['py39']
include = '\.pyi?$'
extend-exclude = '''
/(
  # directories
  \.eggs
  | \.git
  | \.hg
  | \.mypy_cache
  | \.tox
  | \.venv
  | build
  | dist
)/
'''

[tool.isort]
profile = "black"
multi_line_output = 3
include_trailing_comma = true
force_grid_wrap = 0
use_parentheses = true
ensure_newline_before_comments = true
line_length = 88

[tool.mypy]
python_version = "3.9"
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true
disallow_incomplete_defs = true
check_untyped_defs = true
disallow_untyped_decorators = true
no_implicit_optional = true
warn_redundant_casts = true
warn_unused_ignores = true
warn_no_return = true
warn_unreachable = true
strict_equality = true

[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py", "*_test.py"]
addopts = "-v --tb=short --strict-markers"
markers = [
    "slow: marks tests as slow (deselect with '-m \"not slow\"')",
    "integration: marks tests as integration tests",
]

[tool.coverage.run]
source = ["src"]
omit = [
    "*/tests/*",
    "*/venv/*",
    "*/.venv/*",
    "*/migrations/*",
]

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "if self.debug:",
    "if settings.DEBUG",
    "raise AssertionError",
    "raise NotImplementedError",
    "if 0:",
    "if __name__ == .__main__.:",
    "class .*\bProtocol\):",
    "@(abc\.)?abstractmethod",
]
show_missing = true
skip_covered = false
`, bg.spec.Name, bg.spec.Description)
}

// generatePythonDockerfile creates Dockerfile for Python backend
func (bg *BackendGenerator) generatePythonDockerfile() string {
	return `# Multi-stage build for production optimization
FROM python:3.11-slim AS builder

WORKDIR /app

# Install system dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir --user -r requirements.txt

# Production stage
FROM python:3.11-slim AS production

# Create non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

WORKDIR /app

# Install runtime dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        libpq5 \
    && rm -rf /var/lib/apt/lists/*

# Copy Python dependencies from builder stage
COPY --from=builder /root/.local /home/appuser/.local

# Copy application code
COPY --chown=appuser:appuser . .

# Switch to non-root user
USER appuser

# Add local user's bin to PATH
ENV PATH=/home/appuser/.local/bin:$PATH

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD python -c "import requests; requests.get('http://localhost:8080/health')"

CMD ["python", "src/main.py"]`
}

// generatePythonMain creates main.py for Python backend
func (bg *BackendGenerator) generatePythonMain() string {
	framework := bg.spec.GetBackendFramework()
	userRoles := bg.spec.GetUserRolesList()
	rolesStr := strings.Join(userRoles, ", ")
	if rolesStr == "" {
		rolesStr = "user"
	}

	if framework == "django" {
		return bg.generateDjangoMain()
	}

	return fmt.Sprintf(`"""
%s Backend API
Generated with AI CodeKeeper v2.0.0
Domain: %s
Core Entity: %s
"""

import os
import logging
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.trustedhost import TrustedHostMiddleware
from fastapi.responses import JSONResponse
import uvicorn
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%%(asctime)s - %%(name)s - %%(levelname)s - %%(message)s"
)
logger = logging.getLogger(__name__)

@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager"""
    logger.info("🚀 Starting %s Backend Server")
    logger.info("📊 Domain: %s")
    logger.info("👥 User Roles: %s")
    yield
    logger.info("💾 Shutting down gracefully")

# Create FastAPI application
app = FastAPI(
    title="%s API",
    description="%s",
    version="1.0.0",
    docs_url="/docs" if os.getenv("ENVIRONMENT") != "production" else None,
    redoc_url="/redoc" if os.getenv("ENVIRONMENT") != "production" else None,
    lifespan=lifespan
)

# Security middleware
app.add_middleware(
    TrustedHostMiddleware,
    allowed_hosts=["localhost", "127.0.0.1", "*.render.com", "*.amazonaws.com"]
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
        "http://localhost:5173",
        os.getenv("FRONTEND_URL", "")
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "timestamp": "2024-01-01T00:00:00Z",
        "service": "%s-backend",
        "version": "1.0.0",
        "domain": "%s"
    }

@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "Welcome to %s API",
        "docs": "/docs",
        "health": "/health"
    }

# API Routes would be added here
# app.include_router(auth_router, prefix="/api/auth", tags=["authentication"])
# app.include_router(%s_router, prefix="/api/%s", tags=["%s"])

@app.exception_handler(Exception)
async def global_exception_handler(request, exc):
    """Global exception handler"""
    logger.error(f"Global exception: {exc}")
    return JSONResponse(
        status_code=500,
        content={
            "error": "Internal server error",
            "message": str(exc) if os.getenv("ENVIRONMENT") == "development" else "Something went wrong"
        }
    )

if __name__ == "__main__":
    port = int(os.getenv("PORT", 8080))
    host = os.getenv("HOST", "0.0.0.0")
    
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        reload=os.getenv("ENVIRONMENT") == "development",
        log_level="info"
    )
`, 
	bg.spec.Name, 
	bg.spec.Domain, 
	bg.spec.CoreEntity,
	bg.spec.Name,
	bg.spec.Domain,
	rolesStr,
	bg.spec.Name,
	bg.spec.Description,
	bg.spec.Name,
	bg.spec.Domain,
	bg.spec.Name,
	strings.ToLower(bg.spec.CoreEntity),
	strings.ToLower(bg.spec.CoreEntity),
	bg.spec.CoreEntity)
}

// generateDjangoMain creates Django-specific main configuration
func (bg *BackendGenerator) generateDjangoMain() string {
	return fmt.Sprintf(`"""
Django WSGI config for %s project.
Generated with AI CodeKeeper v2.0.0
"""

import os
from django.core.wsgi import application

os.environ.setdefault('DJANGO_SETTINGS_MODULE', '%s.settings')

application = get_wsgi_application()
`, bg.spec.Name, bg.spec.Name)
}

// generateFlake8Config creates .flake8 configuration
func (bg *BackendGenerator) generateFlake8Config() string {
	return `[flake8]
max-line-length = 88
extend-ignore = E203, E266, E501, W503, F403, F401
max-complexity = 18
select = B,C,E,F,W,T4,B9
exclude = 
    .git,
    __pycache__,
    venv,
    .venv,
    migrations,
    .mypy_cache,
    .pytest_cache,
    build,
    dist

per-file-ignores =
    # imported but unused
    __init__.py: F401
    # module level import not at top of file
    settings.py: E402`
}