package models

import (
	"fmt"
	"os"
	"strings"
)

// ModelSelector handles AI model selection based on task type and requirements
type ModelSelector struct {
	config ModelConfiguration
}

// ModelConfiguration defines available AI models and their capabilities
type ModelConfiguration struct {
	Services map[string]ServiceConfig `json:"services"`
	Tasks    map[string]TaskConfig    `json:"tasks"`
}

// ServiceConfig defines an AI service (OpenAI, Anthropic, etc.)
type ServiceConfig struct {
	Name      string            `json:"name"`
	Available bool              `json:"available"`
	Models    map[string]Model  `json:"models"`
	APIKey    string            `json:"api_key,omitempty"`
	Endpoint  string            `json:"endpoint,omitempty"`
}

// Model defines a specific AI model and its capabilities
type Model struct {
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
	MaxTokens    int      `json:"max_tokens"`
	CostPer1K    float64  `json:"cost_per_1k"`
	Speed        string   `json:"speed"` // "fast", "medium", "slow"
	Quality      string   `json:"quality"` // "high", "medium", "basic"
}

// TaskConfig defines requirements for different types of tasks
type TaskConfig struct {
	Description     string   `json:"description"`
	RequiredCaps    []string `json:"required_capabilities"`
	PreferredSpeed  string   `json:"preferred_speed"`
	PreferredQuality string  `json:"preferred_quality"`
	MaxCost         float64  `json:"max_cost_per_1k"`
}

// TechStackRecommendation holds AI recommendations for technology choices
type TechStackRecommendation struct {
	Reasoning        string                 `json:"reasoning"`
	Confidence       float64               `json:"confidence"`
	Recommendations  []ComponentChoice     `json:"recommendations"`
	Alternatives     []ComponentChoice     `json:"alternatives,omitempty"`
	Warnings         []string              `json:"warnings,omitempty"`
	ModelUsed        string                `json:"model_used"`
}

// ComponentChoice represents a recommended technology choice
type ComponentChoice struct {
	Component  string  `json:"component"`  // "backend", "database", "frontend", etc.
	Choice     string  `json:"choice"`     // "Node.js with Express", "PostgreSQL", etc.
	Reasoning  string  `json:"reasoning"`  // Why this choice was made
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
	Icon       string  `json:"icon"`       // Emoji for display
}

// NewSelector creates a new model selector with default configuration
func NewSelector() *ModelSelector {
	selector := &ModelSelector{}
	selector.loadConfiguration()
	return selector
}

// loadConfiguration loads model configuration from environment and defaults
func (ms *ModelSelector) loadConfiguration() {
	// Default configuration
	ms.config = ModelConfiguration{
		Services: map[string]ServiceConfig{
			"openai": {
				Name:      "OpenAI",
				Available: os.Getenv("OPENAI_API_KEY") != "",
				APIKey:    os.Getenv("OPENAI_API_KEY"),
				Models: map[string]Model{
					"gpt-4o": {
						Name:         "gpt-4o",
						Capabilities: []string{"reasoning", "code_generation", "analysis", "planning"},
						MaxTokens:    128000,
						CostPer1K:    0.005,
						Speed:        "medium",
						Quality:      "high",
					},
					"gpt-3.5-turbo": {
						Name:         "gpt-3.5-turbo",
						Capabilities: []string{"code_generation", "analysis"},
						MaxTokens:    16385,
						CostPer1K:    0.001,
						Speed:        "fast",
						Quality:      "medium",
					},
				},
			},
			"anthropic": {
				Name:      "Anthropic",
				Available: os.Getenv("ANTHROPIC_API_KEY") != "",
				APIKey:    os.Getenv("ANTHROPIC_API_KEY"),
				Models: map[string]Model{
					"claude-3-opus": {
						Name:         "claude-3-opus-20240229",
						Capabilities: []string{"reasoning", "analysis", "safety", "complex_planning"},
						MaxTokens:    200000,
						CostPer1K:    0.015,
						Speed:        "slow",
						Quality:      "high",
					},
					"claude-3-sonnet": {
						Name:         "claude-3-sonnet-20240229",
						Capabilities: []string{"reasoning", "code_generation", "analysis"},
						MaxTokens:    200000,
						CostPer1K:    0.003,
						Speed:        "medium",
						Quality:      "high",
					},
				},
			},
			"gemini": {
				Name:      "Google Gemini",
				Available: os.Getenv("GOOGLE_API_KEY") != "",
				APIKey:    os.Getenv("GOOGLE_API_KEY"),
				Models: map[string]Model{
					"gemini-pro": {
						Name:         "gemini-pro",
						Capabilities: []string{"reasoning", "code_generation", "explanation"},
						MaxTokens:    32768,
						CostPer1K:    0.0005,
						Speed:        "fast",
						Quality:      "medium",
					},
				},
			},
		},
		Tasks: map[string]TaskConfig{
			"tech_recommendation": {
				Description:      "Recommend optimal technology stack based on requirements",
				RequiredCaps:     []string{"reasoning", "analysis", "planning"},
				PreferredSpeed:   "medium",
				PreferredQuality: "high",
				MaxCost:          0.01,
			},
			"code_generation": {
				Description:      "Generate application code and boilerplate",
				RequiredCaps:     []string{"code_generation"},
				PreferredSpeed:   "medium",
				PreferredQuality: "high",
				MaxCost:          0.005,
			},
			"security_analysis": {
				Description:      "Analyze code for security vulnerabilities",
				RequiredCaps:     []string{"analysis", "safety"},
				PreferredSpeed:   "slow",
				PreferredQuality: "high",
				MaxCost:          0.015,
			},
		},
	}
}

// SelectBestModel chooses the optimal model for a given task
func (ms *ModelSelector) SelectBestModel(taskType string) (*ServiceConfig, *Model, string, error) {
	taskConfig, exists := ms.config.Tasks[taskType]
	if !exists {
		return nil, nil, "", fmt.Errorf("unknown task type: %s", taskType)
	}

	var bestService *ServiceConfig
	var bestModel *Model
	var bestScore float64
	var reasoning strings.Builder

	reasoning.WriteString(fmt.Sprintf("Task: %s\n", taskConfig.Description))
	reasoning.WriteString(fmt.Sprintf("Requirements: %s\n", strings.Join(taskConfig.RequiredCaps, ", ")))

	for _, service := range ms.config.Services {
		if !service.Available {
			reasoning.WriteString(fmt.Sprintf("- %s: Unavailable (missing API key)\n", service.Name))
			continue
		}

		for _, model := range service.Models {
			score := ms.scoreModel(&model, &taskConfig)
			reasoning.WriteString(fmt.Sprintf("- %s/%s: Score %.2f\n", service.Name, model.Name, score))

			if score > bestScore {
				bestScore = score
				bestService = &service
				bestModel = &model
			}
		}
	}

	if bestModel == nil {
		return nil, nil, "", fmt.Errorf("no suitable model found for task: %s", taskType)
	}

	reasoning.WriteString(fmt.Sprintf("\nSelected: %s/%s (score: %.2f)\n", bestService.Name, bestModel.Name, bestScore))
	return bestService, bestModel, reasoning.String(), nil
}

// scoreModel calculates a score for how well a model fits a task
func (ms *ModelSelector) scoreModel(model *Model, task *TaskConfig) float64 {
	score := 0.0

	// Check required capabilities
	capScore := 0.0
	for _, reqCap := range task.RequiredCaps {
		for _, modelCap := range model.Capabilities {
			if reqCap == modelCap {
				capScore += 1.0
				break
			}
		}
	}
	capScore = capScore / float64(len(task.RequiredCaps))
	score += capScore * 0.4 // 40% weight on capabilities

	// Speed preference
	speedScore := 0.0
	switch {
	case task.PreferredSpeed == model.Speed:
		speedScore = 1.0
	case task.PreferredSpeed == "fast" && model.Speed == "medium":
		speedScore = 0.7
	case task.PreferredSpeed == "medium" && (model.Speed == "fast" || model.Speed == "slow"):
		speedScore = 0.8
	default:
		speedScore = 0.5
	}
	score += speedScore * 0.2 // 20% weight on speed

	// Quality preference  
	qualityScore := 0.0
	switch {
	case task.PreferredQuality == model.Quality:
		qualityScore = 1.0
	case task.PreferredQuality == "high" && model.Quality == "medium":
		qualityScore = 0.8
	case task.PreferredQuality == "medium" && model.Quality == "high":
		qualityScore = 0.9
	default:
		qualityScore = 0.6
	}
	score += qualityScore * 0.3 // 30% weight on quality

	// Cost consideration
	costScore := 1.0
	if model.CostPer1K > task.MaxCost {
		costScore = task.MaxCost / model.CostPer1K
	}
	score += costScore * 0.1 // 10% weight on cost

	return score
}

// RecommendTechStack uses AI to recommend optimal technology stack
func (ms *ModelSelector) RecommendTechStack(spec interface{}) (*TechStackRecommendation, error) {
	service, model, reasoning, err := ms.SelectBestModel("tech_recommendation")
	if err != nil {
		return nil, fmt.Errorf("no model available for tech recommendations: %w", err)
	}

	// For now, return a structured recommendation based on common patterns
	// In a full implementation, this would call the actual AI service
	
	recommendation := &TechStackRecommendation{
		Reasoning:  fmt.Sprintf("Using %s/%s for analysis:\n%s", service.Name, model.Name, reasoning),
		Confidence: 0.85,
		ModelUsed:  fmt.Sprintf("%s/%s", service.Name, model.Name),
		Recommendations: []ComponentChoice{
			{
				Component:  "Backend",
				Choice:     "Node.js with Express",
				Reasoning:  "Fast development, great ecosystem, excellent for APIs",
				Confidence: 0.9,
				Icon:       "🚀",
			},
			{
				Component:  "Database",
				Choice:     "PostgreSQL with Redis",
				Reasoning:  "Reliable ACID compliance with fast caching layer",
				Confidence: 0.95,
				Icon:       "🗄️",
			},
			{
				Component:  "Frontend",
				Choice:     "React with TypeScript",
				Reasoning:  "Industry standard, great tooling, type safety",
				Confidence: 0.88,
				Icon:       "⚛️",
			},
		},
	}

	return recommendation, nil
}

// GetAvailableServices returns list of services with API keys configured
func (ms *ModelSelector) GetAvailableServices() []string {
	var available []string
	for name, service := range ms.config.Services {
		if service.Available {
			available = append(available, name)
		}
	}
	return available
}

// GetConfiguration returns the current model configuration
func (ms *ModelSelector) GetConfiguration() ModelConfiguration {
	return ms.config
}