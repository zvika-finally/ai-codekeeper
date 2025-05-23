package cli

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zvika-finally/ai-codekeeper/internal/generator"
	"github.com/zvika-finally/ai-codekeeper/internal/models"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [project-name]",
		Short: "Generate a new application following AI Master Prompt specifications",
		Long: `Creates a complete, production-ready application based on user requirements.
Follows the 8-question process defined in AI_MASTER_PROMPT.md:

1. Application Name
2. Description & Core Entity  
3. Backend Framework Choice
4. Database Needs
5. API Style
6. User Roles (RBAC)
7. Integration Points
8. Internationalization

Generates monorepo structure (apps/, packages/, infra/, docs/) with
DevContainer/Docker Compose setup for unified development environment.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runCreate,
	}

	cmd.Flags().String("domain", "", "Domain expertise to apply (fintech, ecommerce, etc.)")
	cmd.Flags().Bool("interactive", true, "Run in interactive mode")
	cmd.Flags().Bool("quick", false, "Use sensible defaults for rapid prototyping")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	color.Blue("🚀 AI Development Framework - Application Generator")
	color.Blue("Following AI Master Prompt specification\n")

	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	}

	// Initialize project specification
	spec := &generator.ProjectSpec{}

	// Get domain preference early to inform other questions
	domain, _ := cmd.Flags().GetString("domain")
	if domain == "" {
		domain = promptForDomain()
	}
	spec.Domain = domain

	// Ask the 8 core questions from AI_MASTER_PROMPT.md
	if err := askCoreQuestions(spec, projectName); err != nil {
		return fmt.Errorf("failed to gather requirements: %w", err)
	}

	// Get AI model recommendations for tech stack
	color.Yellow("🧠 Analyzing requirements and selecting optimal tech stack...")
	if err := getTechStackRecommendations(spec); err != nil {
		return fmt.Errorf("failed to get tech recommendations: %w", err)
	}

	// Generate the project
	color.Green("📁 Generating project structure...")
	
	// Test modular generator temporarily
	if os.Getenv("USE_MODULAR") == "true" {
		color.Cyan("🧪 Using new modular generator...")
		gen := generator.NewModular(spec)
		if err := gen.Generate(); err != nil {
			return fmt.Errorf("failed to generate project with modular generator: %w", err)
		}
	} else {
		gen := generator.New(spec)
		if err := gen.Generate(); err != nil {
			return fmt.Errorf("failed to generate project: %w", err)
		}
	}

	// Success message
	color.Green("\n✅ Project '%s' generated successfully!", spec.Name)
	fmt.Printf("\n📋 Next steps:\n")
	fmt.Printf("  cd %s\n", spec.Name)
	fmt.Printf("  codekeeper env setup    # Setup development environment\n")
	fmt.Printf("  codekeeper check        # Validate setup\n")
	
	if spec.DevEnvironment == "devcontainer" {
		fmt.Printf("  code .              # Open in VS Code with DevContainer\n")
	} else {
		fmt.Printf("  docker compose up   # Start development environment\n")
	}

	return nil
}

func promptForDomain() string {
	domains := []string{
		"generic - General web application",
		"fintech - Financial technology (payments, banking, compliance)",
		"ecommerce - E-commerce and retail",
		"healthcare - Healthcare and medical systems",
		"education - Educational platforms and tools",
	}

	var selected string
	prompt := &survey.Select{
		Message: "What domain expertise should I apply?",
		Options: domains,
		Help:    "Domain expertise provides specialized knowledge for tech stack recommendations and best practices",
	}
	survey.AskOne(prompt, &selected)

	// Extract domain name (before " - ")
	if idx := len(selected); idx > 0 {
		for i, r := range selected {
			if r == ' ' && i+3 < len(selected) && selected[i:i+3] == " - " {
				return selected[:i]
			}
		}
	}
	return "generic"
}

func askCoreQuestions(spec *generator.ProjectSpec, initialName string) error {
	color.Cyan("📝 Gathering application requirements (8 core questions):\n")

	// 1. Application Name
	if initialName == "" {
		prompt := &survey.Input{
			Message: "1. What is the name of your application?",
			Help:    "This will be used for directory names, package names, etc. (e.g., 'my-awesome-app')",
		}
		if err := survey.AskOne(prompt, &spec.Name, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	} else {
		spec.Name = initialName
		fmt.Printf("1. Application name: %s\n", spec.Name)
	}

	// 2. Description & Core Entity
	prompt2 := &survey.Multiline{
		Message: "2. Briefly describe your application and identify one core entity for CRUD operations:",
		Help:    "What are the 1-3 core features? What's the central entity? (e.g., 'Task management app for teams. Core entity: Task')",
	}
	if err := survey.AskOne(prompt2, &spec.Description, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Extract core entity from description or ask separately
	entityPrompt := &survey.Input{
		Message: "What is the main entity for initial CRUD operations?",
		Help:    "e.g., 'Task', 'Product', 'User', 'Transaction'",
	}
	if err := survey.AskOne(entityPrompt, &spec.CoreEntity, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// 3. Backend Framework Choice
	backendOptions := []string{
		"Node.js with Express",
		"Node.js with NestJS", 
		"Python with Django",
		"Python with Flask",
		"Go with Gin",
		"Go with standard library",
		"Let AI recommend based on requirements",
	}
	
	var backend string
	backendPrompt := &survey.Select{
		Message: "3. Choose your backend framework:",
		Options: backendOptions,
		Help:    "Consider your team's familiarity, performance needs, and application type",
	}
	if err := survey.AskOne(backendPrompt, &backend); err != nil {
		return err
	}
	spec.Backend = backend

	// 4. Database Needs
	var databases []string
	databasePrompt := &survey.MultiSelect{
		Message: "4. Database requirements (PostgreSQL is default):",
		Options: []string{
			"PostgreSQL (default - transactional data)",
			"Redis (caching, sessions)",
			"MongoDB (flexible schemas)",
			"Elasticsearch (full-text search)",
			"Vector DB for AI (Pinecone, Weaviate)",
			"Time-series DB (InfluxDB, TimescaleDB)",
			"Data warehouse (Snowflake, BigQuery)",
		},
		Help: "Select all databases your application needs",
	}
	if err := survey.AskOne(databasePrompt, &databases); err != nil {
		return err
	}
	spec.Databases = databases

	// 5. API Style
	apiOptions := []string{
		"RESTful APIs",
		"GraphQL",
		"gRPC",
		"Let AI recommend based on requirements",
	}
	
	var apiStyle string
	apiPrompt := &survey.Select{
		Message: "5. API design style:",
		Options: apiOptions,
		Help:    "RESTful for broad compatibility, GraphQL for complex queries, gRPC for high-performance microservices",
	}
	if err := survey.AskOne(apiPrompt, &apiStyle); err != nil {
		return err
	}
	spec.APIStyle = apiStyle

	// 6. User Roles
	rolesPrompt := &survey.Input{
		Message: "6. Key user roles for RBAC (comma-separated):",
		Help:    "e.g., 'user, admin, editor' or 'None' for single-user type",
		Default: "user, admin",
	}
	if err := survey.AskOne(rolesPrompt, &spec.UserRoles); err != nil {
		return err
	}

	// 7. Integration Points
	var integrations []string
	integrationsPrompt := &survey.MultiSelect{
		Message: "7. Integration points:",
		Options: []string{
			"JIRA (task tracking)",
			"Confluence (documentation)",
			"Figma (UI designs)",
			"GitHub Actions (CI/CD)",
			"AWS (cloud services)",
			"Stripe (payments)",
			"Auth0 (authentication)",
			"None",
		},
		Help: "Select external systems you plan to integrate with",
	}
	if err := survey.AskOne(integrationsPrompt, &integrations); err != nil {
		return err
	}
	spec.Integrations = integrations

	// 8. Internationalization
	var needsI18n bool
	i18nPrompt := &survey.Confirm{
		Message: "8. Do you need internationalization (i18n) support?",
		Help:    "Sets up basic i18n structure for multi-language support",
		Default: false,
	}
	if err := survey.AskOne(i18nPrompt, &needsI18n); err != nil {
		return err
	}
	spec.I18n = needsI18n

	return nil
}

func getTechStackRecommendations(spec *generator.ProjectSpec) error {
	// Initialize AI model selector
	selector := models.NewSelector()
	
	// Get recommendation for tech stack
	recommendation, err := selector.RecommendTechStack(spec)
	if err != nil {
		color.Yellow("⚠️  Could not get AI recommendations, using defaults")
		return nil
	}

	// Display recommendations with reasoning
	color.Cyan("\n🎯 AI Tech Stack Recommendations:")
	for _, rec := range recommendation.Recommendations {
		fmt.Printf("  %s %s: %s\n", rec.Icon, rec.Component, rec.Choice)
		fmt.Printf("    Reason: %s\n", rec.Reasoning)
	}

	// Ask user to confirm or modify
	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Accept these recommendations?",
		Default: true,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if confirm {
		spec.TechStack = recommendation
	} else {
		color.Yellow("You can manually adjust recommendations after project generation")
	}

	// Set development environment preference
	devEnvOptions := []string{
		"DevContainer (VS Code)",
		"Docker Compose (universal)",
		"Nix (reproducible)",
		"Let AI recommend",
	}
	
	var devEnv string
	devEnvPrompt := &survey.Select{
		Message: "Preferred development environment:",
		Options: devEnvOptions,
		Help:    "Ensures dev/prod parity and unified team environment",
	}
	if err := survey.AskOne(devEnvPrompt, &devEnv); err != nil {
		return err
	}
	
	switch devEnv {
	case "DevContainer (VS Code)":
		spec.DevEnvironment = "devcontainer"
	case "Docker Compose (universal)":
		spec.DevEnvironment = "docker-compose"
	case "Nix (reproducible)":
		spec.DevEnvironment = "nix"
	default:
		spec.DevEnvironment = "docker-compose" // Default recommendation
	}

	return nil
}