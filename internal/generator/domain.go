package generator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DomainConfigYAML represents domain-specific configuration from YAML files
type DomainConfigYAML struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	
	Entities map[string]string `yaml:"entities"`
	
	TechPreferences struct {
		HighSecurity        []string `yaml:"high_security"`
		RealTimeProcessing  []string `yaml:"real_time_processing"`
		ComplianceFriendly  []string `yaml:"compliance_friendly"`
		AuditLogging        []string `yaml:"audit_logging"`
		PaymentProcessing   []string `yaml:"payment_processing"`
	} `yaml:"tech_preferences"`
	
	Compliance       []string `yaml:"compliance"`
	SecurityPatterns []string `yaml:"security_patterns"`
	GuardRails       []string `yaml:"guard_rails"`
	
	RecommendedIntegrations struct {
		PaymentProcessors     []string `yaml:"payment_processors"`
		BankingApis          []string `yaml:"banking_apis"`
		IdentityVerification []string `yaml:"identity_verification"`
		FraudPrevention      []string `yaml:"fraud_prevention"`
		Compliance           []string `yaml:"compliance"`
	} `yaml:"recommended_integrations"`
	
	BestPractices struct {
		DataHandling []string `yaml:"data_handling"`
		Security     []string `yaml:"security"`
		Compliance   []string `yaml:"compliance"`
		Operations   []string `yaml:"operations"`
	} `yaml:"best_practices"`
	
	EnvironmentVariables struct {
		Development map[string]string `yaml:"development"`
		Production  map[string]string `yaml:"production"`
	} `yaml:"environment_variables"`
	
	CodeTemplates map[string]string `yaml:"code_templates"`
}

// LoadDomainConfig loads domain configuration from YAML file
func LoadDomainConfig(domain string) (*DomainConfigYAML, error) {
	if domain == "" || domain == "generic" {
		return &DomainConfigYAML{
			Name:        "generic",
			Description: "Generic web application",
		}, nil
	}
	
	configPath := filepath.Join("domains", domain, "config.yaml")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		// Fallback to generic if domain config not found
		return &DomainConfigYAML{
			Name:        domain,
			Description: fmt.Sprintf("%s application", domain),
		}, nil
	}
	
	var config DomainConfigYAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse domain config %s: %w", configPath, err)
	}
	
	return &config, nil
}

// GetRecommendedBackend returns recommended backend tech for domain
func (dc *DomainConfigYAML) GetRecommendedBackend() string {
	if len(dc.TechPreferences.HighSecurity) > 0 {
		return dc.TechPreferences.HighSecurity[0]
	}
	if len(dc.TechPreferences.RealTimeProcessing) > 0 {
		return dc.TechPreferences.RealTimeProcessing[0]
	}
	return "javascript" // Default
}

// GetRecommendedDatabase returns recommended database for domain
func (dc *DomainConfigYAML) GetRecommendedDatabase() string {
	if len(dc.TechPreferences.ComplianceFriendly) > 0 {
		return dc.TechPreferences.ComplianceFriendly[0]
	}
	return "postgresql" // Default
}

// GetDomainSpecificEntity returns a domain-appropriate core entity
func (dc *DomainConfigYAML) GetDomainSpecificEntity() string {
	// Return first entity or default
	for entity := range dc.Entities {
		return entity
	}
	return "Item"
}

// EnhanceProjectSpec enhances project spec with domain-specific recommendations
func (dc *DomainConfigYAML) EnhanceProjectSpec(spec *ProjectSpec) {
	// Only enhance if not already set
	if spec.Backend == "" {
		spec.Backend = dc.GetRecommendedBackend()
	}
	
	if len(spec.Databases) == 0 {
		spec.Databases = []string{dc.GetRecommendedDatabase()}
	}
	
	// Add domain-specific description enhancement
	if spec.Description == "" {
		spec.Description = fmt.Sprintf("A %s application", dc.Name)
	}
}