package infrastructure

import "strings"

// AWS Terraform generation methods

func (ig *InfrastructureGenerator) generateAWSMain() string {
	return `terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  backend "s3" {
    # Configure your S3 backend here
    # bucket = "your-terraform-state-bucket"
    # key    = "` + ig.spec.Name + `/terraform.tfstate"
    # region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = "` + ig.spec.Name + `"
      Environment = var.environment
      Domain      = "` + ig.spec.Domain + `"
      ManagedBy   = "Terraform"
    }
  }
}

# VPC and networking
module "vpc" {
  source = "./modules/vpc"
  
  name        = var.project_name
  environment = var.environment
  cidr_block  = var.vpc_cidr
}

# RDS Database
module "rds" {
  source = "./modules/rds"
  
  name               = var.project_name
  environment        = var.environment
  vpc_id            = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  database_name     = var.database_name
  master_username   = var.db_username
  master_password   = var.db_password
  instance_class    = var.db_instance_class
}

# ECS Cluster and Service
module "ecs" {
  source = "./modules/ecs"
  
  name               = var.project_name
  environment        = var.environment
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
  
  # Application configuration
  backend_image     = var.backend_image
  frontend_image    = var.frontend_image
  backend_port      = var.backend_port
  frontend_port     = var.frontend_port
  
  # Database connection
  database_url = module.rds.connection_string
  
  # Secrets
  jwt_secret_arn = module.secrets.jwt_secret_arn
}

# Secrets Management
module "secrets" {
  source = "./modules/secrets"
  
  name        = var.project_name
  environment = var.environment
  
  secrets = {
    jwt_secret = var.jwt_secret
    db_password = var.db_password` + ig.getDomainSpecificSecrets() + `
  }
}
`
}

func (ig *InfrastructureGenerator) generateAWSVariables() string {
	return `variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "` + ig.spec.Name + `"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "` + strings.ReplaceAll(ig.spec.Name, "-", "_") + `"
}

variable "db_username" {
  description = "Database master username"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "backend_image" {
  description = "Backend Docker image"
  type        = string
  default     = "` + ig.spec.Name + `-backend:latest"
}

variable "frontend_image" {
  description = "Frontend Docker image"
  type        = string
  default     = "` + ig.spec.Name + `-frontend:latest"
}

variable "backend_port" {
  description = "Backend service port"
  type        = number
  default     = 8080
}

variable "frontend_port" {
  description = "Frontend service port"
  type        = number
  default     = 3000
}

variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
}

# Domain-specific variables
` + ig.getDomainSpecificVariables()
}

func (ig *InfrastructureGenerator) generateAWSOutputs() string {
	return `output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "database_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.endpoint
  sensitive   = true
}

output "load_balancer_dns" {
  description = "Load balancer DNS name"
  value       = module.ecs.load_balancer_dns
}

output "backend_service_url" {
  description = "Backend service URL"
  value       = "https://${module.ecs.load_balancer_dns}"
}

output "frontend_service_url" {
  description = "Frontend service URL"
  value       = "https://${module.ecs.load_balancer_dns}"
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs.cluster_name
}
`
}

func (ig *InfrastructureGenerator) generateAWSTerraformVars() string {
	return `# AWS Infrastructure Configuration
# Copy this file to terraform.tfvars and fill in your values

aws_region  = "us-east-1"
environment = "dev"

# Database configuration
db_password = "change-me-in-production"
jwt_secret  = "change-me-in-production"

# Scaling configuration
db_instance_class = "db.t3.micro"  # Use db.t3.small+ for production

# Domain-specific configuration (` + ig.spec.Domain + `)
` + ig.getDomainSpecificTerraformVars()
}

func (ig *InfrastructureGenerator) generateECSModule() string {
	return `# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.name}-${var.environment}"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${var.name}-${var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets           = var.public_subnet_ids

  enable_deletion_protection = false
}

# ALB Security Group
resource "aws_security_group" "alb" {
  name_prefix = "${var.name}-alb"
  vpc_id      = var.vpc_id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    protocol    = "tcp"
    from_port   = 443
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Target Groups
resource "aws_lb_target_group" "backend" {
  name     = "${var.name}-backend-${var.environment}"
  port     = var.backend_port
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/health"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
}

resource "aws_lb_target_group" "frontend" {
  name     = "${var.name}-frontend-${var.environment}"
  port     = var.frontend_port
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
}

# ALB Listeners
resource "aws_lb_listener" "backend" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.backend.arn
  }
}

# ECS Task Definition
resource "aws_ecs_task_definition" "backend" {
  family                   = "${var.name}-backend-${var.environment}"
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn           = aws_iam_role.ecs_task.arn
  network_mode            = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                     = 256
  memory                  = 512

  container_definitions = jsonencode([
    {
      name  = "backend"
      image = var.backend_image
      
      portMappings = [
        {
          containerPort = var.backend_port
          protocol      = "tcp"
        }
      ]
      
      environment = [
        {
          name  = "NODE_ENV"
          value = "production"
        },
        {
          name  = "PORT"
          value = tostring(var.backend_port)
        }
      ]
      
      secrets = [
        {
          name      = "DATABASE_URL"
          valueFrom = var.database_url
        },
        {
          name      = "JWT_SECRET"
          valueFrom = var.jwt_secret_arn
        }
      ]
      
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.backend.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

# ECS Service
resource "aws_ecs_service" "backend" {
  name            = "${var.name}-backend-${var.environment}"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets         = var.private_subnet_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.backend.arn
    container_name   = "backend"
    container_port   = var.backend_port
  }

  depends_on = [aws_lb_listener.backend]
}

# ECS Security Group
resource "aws_security_group" "ecs_tasks" {
  name_prefix = "${var.name}-ecs-tasks"
  vpc_id      = var.vpc_id

  ingress {
    protocol        = "tcp"
    from_port       = var.backend_port
    to_port         = var.backend_port
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "backend" {
  name              = "/ecs/${var.name}-backend-${var.environment}"
  retention_in_days = 30
}

# IAM Roles
resource "aws_iam_role" "ecs_execution" {
  name = "${var.name}-ecs-execution-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution" {
  role       = aws_iam_role.ecs_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "ecs_task" {
  name = "${var.name}-ecs-task-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

data "aws_region" "current" {}
`
}

func (ig *InfrastructureGenerator) generateECSVariables() string {
	return `variable "name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private subnet IDs"
  type        = list(string)
}

variable "backend_image" {
  description = "Backend Docker image"
  type        = string
}

variable "frontend_image" {
  description = "Frontend Docker image"
  type        = string
}

variable "backend_port" {
  description = "Backend service port"
  type        = number
  default     = 8080
}

variable "frontend_port" {
  description = "Frontend service port"
  type        = number
  default     = 3000
}

variable "database_url" {
  description = "Database connection URL"
  type        = string
}

variable "jwt_secret_arn" {
  description = "JWT secret ARN from Secrets Manager"
  type        = string
}
`
}

func (ig *InfrastructureGenerator) generateECSOutputs() string {
	return `output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "load_balancer_dns" {
  description = "Load balancer DNS name"
  value       = aws_lb.main.dns_name
}

output "load_balancer_zone_id" {
  description = "Load balancer zone ID"
  value       = aws_lb.main.zone_id
}
`
}

func (ig *InfrastructureGenerator) getDomainSpecificSecrets() string {
	switch ig.spec.Domain {
	case "fintech":
		return `
    stripe_secret_key = var.stripe_secret_key
    plaid_client_id = var.plaid_client_id
    plaid_secret = var.plaid_secret`
	case "healthcare":
		return `
    fhir_auth_token = var.fhir_auth_token
    hipaa_encryption_key = var.hipaa_encryption_key`
	case "ecommerce":
		return `
    shopify_access_token = var.shopify_access_token
    stripe_secret_key = var.stripe_secret_key`
	default:
		return ""
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificVariables() string {
	switch ig.spec.Domain {
	case "fintech":
		return `# Fintech-specific variables
variable "stripe_secret_key" {
  description = "Stripe secret key for payment processing"
  type        = string
  sensitive   = true
}

variable "plaid_client_id" {
  description = "Plaid client ID for banking data"
  type        = string
}

variable "plaid_secret" {
  description = "Plaid secret for banking data"
  type        = string
  sensitive   = true
}`
	case "healthcare":
		return `# Healthcare-specific variables
variable "fhir_auth_token" {
  description = "FHIR API authentication token"
  type        = string
  sensitive   = true
}

variable "hipaa_encryption_key" {
  description = "HIPAA-compliant encryption key"
  type        = string
  sensitive   = true
}`
	case "ecommerce":
		return `# E-commerce-specific variables
variable "shopify_access_token" {
  description = "Shopify API access token"
  type        = string
  sensitive   = true
}

variable "stripe_secret_key" {
  description = "Stripe secret key for payment processing"
  type        = string
  sensitive   = true
}`
	default:
		return `# No domain-specific variables needed`
	}
}

func (ig *InfrastructureGenerator) getDomainSpecificTerraformVars() string {
	switch ig.spec.Domain {
	case "fintech":
		return `# Fintech payment processing
stripe_secret_key = "sk_test_..."
plaid_client_id   = "your_plaid_client_id"
plaid_secret      = "your_plaid_secret"`
	case "healthcare":
		return `# Healthcare FHIR integration
fhir_auth_token      = "your_fhir_token"
hipaa_encryption_key = "your_hipaa_key"`
	case "ecommerce":
		return `# E-commerce integrations
shopify_access_token = "your_shopify_token"
stripe_secret_key    = "sk_test_..."`
	default:
		return `# No domain-specific configuration needed`
	}
}