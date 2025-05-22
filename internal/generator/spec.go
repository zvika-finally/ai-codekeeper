package generator

import (
	"fmt"
	"strings"
	
	"github.com/tycoonlabs/ai-codekeeper/internal/models"
)

// ProjectSpec holds all user requirements for project generation
// Based on the 8 questions from AI_MASTER_PROMPT.md
type ProjectSpec struct {
	// Core requirements (from 8 questions)
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	CoreEntity   string   `json:"core_entity"`
	Backend      string   `json:"backend"`
	Databases    []string `json:"databases"`
	APIStyle     string   `json:"api_style"`
	UserRoles    string   `json:"user_roles"`
	Integrations []string `json:"integrations"`
	I18n         bool     `json:"i18n"`

	// Domain expertise
	Domain string `json:"domain"`

	// AI recommendations
	TechStack *models.TechStackRecommendation `json:"tech_stack,omitempty"`

	// Development environment
	DevEnvironment string `json:"dev_environment"`

	// Generated metadata
	ProjectPath string            `json:"project_path,omitempty"`
	GeneratedAt string            `json:"generated_at,omitempty"`
	Framework   *FrameworkConfig  `json:"framework,omitempty"`
}

// FrameworkConfig holds configuration for ongoing AI assistance
type FrameworkConfig struct {
	Version        string                    `json:"version"`
	Domain         DomainConfig              `json:"domain"`
	GuardRails     GuardRailsConfig          `json:"guard_rails"`
	AIModels       models.ModelConfiguration `json:"ai_models"`
	DevEnvironment DevEnvironmentConfig      `json:"dev_environment"`
	Integrations   IntegrationsConfig        `json:"integrations"`
}

// DomainConfig holds domain-specific knowledge and rules
type DomainConfig struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Entities    map[string]string `json:"entities"`
	Patterns    []string          `json:"patterns"`
	Compliance  []string          `json:"compliance"`
	TechPrefs   map[string][]string `json:"tech_preferences"`
}

// GuardRailsConfig defines enforcement rules
type GuardRailsConfig struct {
	Enforcement string   `json:"enforcement"` // "strict", "advisory", "off"
	PreCommit   bool     `json:"pre_commit"`
	CI          bool     `json:"ci"`
	IDE         string   `json:"ide"` // "vscode", "cursor", "terminal"
	Rules       []string `json:"rules"`
}

// DevEnvironmentConfig defines the development setup
type DevEnvironmentConfig struct {
	Type        string            `json:"type"` // "devcontainer", "docker-compose", "nix"
	Services    []string          `json:"services"`
	Ports       map[string]int    `json:"ports"`
	Volumes     []string          `json:"volumes"`
	Environment map[string]string `json:"environment"`
}

// IntegrationsConfig defines external system integrations
type IntegrationsConfig struct {
	JIRA       *JIRAConfig       `json:"jira,omitempty"`
	Confluence *ConfluenceConfig `json:"confluence,omitempty"`
	GitHub     *GitHubConfig     `json:"github,omitempty"`
	Cloud      *CloudConfig      `json:"cloud,omitempty"`
}

type JIRAConfig struct {
	URL        string `json:"url"`
	ProjectKey string `json:"project_key"`
	IssueTypes []string `json:"issue_types"`
}

type ConfluenceConfig struct {
	URL   string `json:"url"`
	Space string `json:"space"`
}

type GitHubConfig struct {
	Repository string   `json:"repository"`
	Workflows  []string `json:"workflows"`
}

type CloudConfig struct {
	Provider string            `json:"provider"` // "aws", "gcp", "azure"
	Services map[string]string `json:"services"`
}

// Validate ensures the project specification is complete and valid
func (spec *ProjectSpec) Validate() error {
	if spec.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if spec.Description == "" {
		return fmt.Errorf("project description is required")
	}
	if spec.CoreEntity == "" {
		return fmt.Errorf("core entity is required")
	}
	if spec.Backend == "" {
		return fmt.Errorf("backend framework is required")
	}
	if spec.APIStyle == "" {
		return fmt.Errorf("API style is required")
	}
	return nil
}

// GetBackendLanguage extracts the programming language from backend choice
func (spec *ProjectSpec) GetBackendLanguage() string {
	switch {
	case strings.Contains(strings.ToLower(spec.Backend), "node"):
		return "javascript"
	case strings.Contains(strings.ToLower(spec.Backend), "python"):
		return "python"
	case strings.Contains(strings.ToLower(spec.Backend), "go"):
		return "go"
	default:
		return "javascript" // Default fallback
	}
}

// GetBackendFramework extracts the framework name from backend choice
func (spec *ProjectSpec) GetBackendFramework() string {
	backend := strings.ToLower(spec.Backend)
	switch {
	case strings.Contains(backend, "express"):
		return "express"
	case strings.Contains(backend, "nestjs"):
		return "nestjs"
	case strings.Contains(backend, "django"):
		return "django"
	case strings.Contains(backend, "flask"):
		return "flask"
	case strings.Contains(backend, "gin"):
		return "gin"
	case strings.Contains(backend, "standard"):
		return "stdlib"
	default:
		return "express" // Default fallback
	}
}

// HasDatabase checks if a specific database is required
func (spec *ProjectSpec) HasDatabase(dbType string) bool {
	for _, db := range spec.Databases {
		if strings.Contains(strings.ToLower(db), strings.ToLower(dbType)) {
			return true
		}
	}
	return false
}

// GetUserRolesList returns user roles as a slice
func (spec *ProjectSpec) GetUserRolesList() []string {
	if spec.UserRoles == "" || strings.ToLower(spec.UserRoles) == "none" {
		return []string{}
	}
	
	roles := strings.Split(spec.UserRoles, ",")
	for i, role := range roles {
		roles[i] = strings.TrimSpace(role)
	}
	return roles
}