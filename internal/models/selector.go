package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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

	// Convert spec to JSON for AI analysis
	specJSON, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec: %w", err)
	}

	// Call AI service to get real recommendations
	recommendation, err := ms.callAIService(service, model, specJSON)
	if err != nil {
		// Fallback to intelligent defaults if AI call fails
		return ms.getIntelligentFallback(service, model, reasoning, spec)
	}

	recommendation.Reasoning = fmt.Sprintf("Using %s/%s for analysis:\n%s\n\nAI Analysis:\n%s", 
		service.Name, model.Name, reasoning, recommendation.Reasoning)
	recommendation.ModelUsed = fmt.Sprintf("%s/%s", service.Name, model.Name)

	return recommendation, nil
}

// callAIService makes actual API call to AI service for tech stack recommendations
func (ms *ModelSelector) callAIService(service *ServiceConfig, model *Model, specJSON []byte) (*TechStackRecommendation, error) {
	prompt := fmt.Sprintf(`You are an expert software architect. Analyze the following project specification and recommend an optimal technology stack.

Project Specification:
%s

Please provide recommendations for:
1. Backend framework/language
2. Database (primary and caching if needed)  
3. Frontend framework
4. Additional tools/services if relevant

For each recommendation, provide:
- Specific choice (e.g., "Node.js with Express", "PostgreSQL")
- Clear reasoning for the choice
- Confidence level (0.0 to 1.0)
- Consider domain-specific needs (fintech needs compliance, healthcare needs security, etc.)

Respond in JSON format matching this structure:
{
  "reasoning": "Overall analysis and approach",
  "confidence": 0.85,
  "recommendations": [
    {
      "component": "Backend",
      "choice": "Specific technology choice",
      "reasoning": "Why this choice is optimal",
      "confidence": 0.9,
      "icon": "🚀"
    }
  ],
  "alternatives": [
    {
      "component": "Backend", 
      "choice": "Alternative option",
      "reasoning": "When this might be better",
      "confidence": 0.8,
      "icon": "⚡"
    }
  ],
  "warnings": ["Any important considerations or trade-offs"]
}`, string(specJSON))

	switch strings.ToLower(service.Name) {
	case "openai":
		return ms.callOpenAI(service, model, prompt)
	case "anthropic":
		return ms.callAnthropic(service, model, prompt)
	case "google gemini":
		return ms.callGemini(service, model, prompt)
	default:
		return nil, fmt.Errorf("unsupported AI service: %s", service.Name)
	}
}

// callOpenAI makes API call to OpenAI
func (ms *ModelSelector) callOpenAI(service *ServiceConfig, model *Model, prompt string) (*TechStackRecommendation, error) {
	url := "https://api.openai.com/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model": model.Name,
		"messages": []map[string]string{
			{"role": "system", "content": "You are an expert software architect who provides technology stack recommendations."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
		"max_tokens": 2000,
	}

	return ms.makeAPICall(url, service.APIKey, payload, "Bearer")
}

// callAnthropic makes API call to Anthropic Claude
func (ms *ModelSelector) callAnthropic(service *ServiceConfig, model *Model, prompt string) (*TechStackRecommendation, error) {
	url := "https://api.anthropic.com/v1/messages"
	
	payload := map[string]interface{}{
		"model": model.Name,
		"max_tokens": 2000,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": prompt,
			},
		},
	}

	return ms.makeAPICall(url, service.APIKey, payload, "x-api-key")
}

// callGemini makes API call to Google Gemini
func (ms *ModelSelector) callGemini(service *ServiceConfig, model *Model, prompt string) (*TechStackRecommendation, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model.Name, service.APIKey)
	
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.3,
			"maxOutputTokens": 2000,
		},
	}

	return ms.makeAPICall(url, "", payload, "")
}

// makeAPICall handles the actual HTTP request and response parsing
func (ms *ModelSelector) makeAPICall(url, apiKey string, payload map[string]interface{}, authType string) (*TechStackRecommendation, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		switch authType {
		case "Bearer":
			req.Header.Set("Authorization", "Bearer "+apiKey)
		case "x-api-key":
			req.Header.Set("x-api-key", apiKey)
			req.Header.Set("anthropic-version", "2023-06-01")
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response and extract content based on service
	content, err := ms.extractContentFromResponse(body, url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content: %w", err)
	}

	// Parse AI response as JSON
	var recommendation TechStackRecommendation
	if err := json.Unmarshal([]byte(content), &recommendation); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	return &recommendation, nil
}

// extractContentFromResponse extracts the actual content from different AI service responses
func (ms *ModelSelector) extractContentFromResponse(body []byte, url string) (string, error) {
	if strings.Contains(url, "openai.com") {
		var response struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}
		if len(response.Choices) == 0 {
			return "", fmt.Errorf("no choices in OpenAI response")
		}
		return response.Choices[0].Message.Content, nil
	}

	if strings.Contains(url, "anthropic.com") {
		var response struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}
		if len(response.Content) == 0 {
			return "", fmt.Errorf("no content in Anthropic response")
		}
		return response.Content[0].Text, nil
	}

	if strings.Contains(url, "googleapis.com") {
		var response struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}
		if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
			return "", fmt.Errorf("no content in Gemini response")
		}
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("unknown service response format")
}

// getIntelligentFallback provides smart defaults when AI service is unavailable
func (ms *ModelSelector) getIntelligentFallback(service *ServiceConfig, model *Model, reasoning string, spec interface{}) (*TechStackRecommendation, error) {
	// Analyze the spec to make intelligent decisions
	specStr := fmt.Sprintf("%+v", spec)
	
	// Detect domain for appropriate recommendations
	domain := "general"
	if strings.Contains(strings.ToLower(specStr), "fintech") || strings.Contains(strings.ToLower(specStr), "financial") {
		domain = "fintech"
	} else if strings.Contains(strings.ToLower(specStr), "healthcare") || strings.Contains(strings.ToLower(specStr), "medical") {
		domain = "healthcare"
	} else if strings.Contains(strings.ToLower(specStr), "ecommerce") || strings.Contains(strings.ToLower(specStr), "commerce") {
		domain = "ecommerce"
	}

	recommendations := ms.getDomainSpecificRecommendations(domain)
	
	return &TechStackRecommendation{
		Reasoning:       fmt.Sprintf("AI service unavailable. Using intelligent fallback based on domain: %s\n%s", domain, reasoning),
		Confidence:      0.75, // Lower confidence for fallback
		ModelUsed:       fmt.Sprintf("%s/%s (fallback)", service.Name, model.Name),
		Recommendations: recommendations,
		Warnings:        []string{"AI service was unavailable - using intelligent defaults", "Consider configuring API keys for AI-powered recommendations"},
	}, nil
}

// getDomainSpecificRecommendations returns appropriate tech stack for each domain
func (ms *ModelSelector) getDomainSpecificRecommendations(domain string) []ComponentChoice {
	switch domain {
	case "fintech":
		return []ComponentChoice{
			{
				Component:  "Backend",
				Choice:     "Go with Gin/Echo",
				Reasoning:  "High performance, strong typing, excellent for financial calculations",
				Confidence: 0.9,
				Icon:       "🏦",
			},
			{
				Component:  "Database",
				Choice:     "PostgreSQL with Redis",
				Reasoning:  "ACID compliance essential for financial data, Redis for fast lookups",
				Confidence: 0.95,
				Icon:       "🔒",
			},
			{
				Component:  "Frontend",
				Choice:     "React with TypeScript",
				Reasoning:  "Type safety critical for financial UIs, excellent ecosystem",
				Confidence: 0.88,
				Icon:       "💰",
			},
		}
	case "healthcare":
		return []ComponentChoice{
			{
				Component:  "Backend",
				Choice:     "Java with Spring Boot",
				Reasoning:  "Enterprise-grade security, HIPAA compliance features, mature ecosystem",
				Confidence: 0.92,
				Icon:       "🏥",
			},
			{
				Component:  "Database",
				Choice:     "PostgreSQL with encryption",
				Reasoning:  "Strong security features, encryption at rest, audit logging",
				Confidence: 0.93,
				Icon:       "🔐",
			},
			{
				Component:  "Frontend",
				Choice:     "React with TypeScript",
				Reasoning:  "Secure development patterns, accessibility features, type safety",
				Confidence: 0.85,
				Icon:       "⚕️",
			},
		}
	case "ecommerce":
		return []ComponentChoice{
			{
				Component:  "Backend",
				Choice:     "Node.js with Express",
				Reasoning:  "Fast development, great for APIs, excellent payment integrations",
				Confidence: 0.88,
				Icon:       "🛒",
			},
			{
				Component:  "Database",
				Choice:     "PostgreSQL with Redis",
				Reasoning:  "Handles complex queries, Redis for cart/session management",
				Confidence: 0.90,
				Icon:       "📦",
			},
			{
				Component:  "Frontend",
				Choice:     "Next.js with TypeScript",
				Reasoning:  "SEO optimization, server-side rendering, great performance",
				Confidence: 0.92,
				Icon:       "🚀",
			},
		}
	default:
		return []ComponentChoice{
			{
				Component:  "Backend",
				Choice:     "Node.js with Express",
				Reasoning:  "Fast development, large ecosystem, versatile for most use cases",
				Confidence: 0.85,
				Icon:       "🚀",
			},
			{
				Component:  "Database",
				Choice:     "PostgreSQL",
				Reasoning:  "Reliable, feature-rich, good performance for most applications",
				Confidence: 0.88,
				Icon:       "🗄️",
			},
			{
				Component:  "Frontend",
				Choice:     "React with TypeScript",
				Reasoning:  "Industry standard, great tooling, type safety",
				Confidence: 0.85,
				Icon:       "⚛️",
			},
		}
	}
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