package generator

// ComponentGenerator defines the interface for all component generators
type ComponentGenerator interface {
	Generate() error
}

// FileGenerator defines the interface for generators that return file maps
type FileGenerator interface {
	Generate() (map[string]string, error)
}

// BackendGeneratorInterface defines the interface for backend generation
type BackendGeneratorInterface interface {
	FileGenerator
	GeneratePackageJSON() string
	GenerateDockerfile() string
	GenerateReadme() string
	GenerateEnvExample() string
}

// FrontendGeneratorInterface defines the interface for frontend generation
type FrontendGeneratorInterface interface {
	FileGenerator
	GeneratePackageJSON() string
	GenerateComponents() map[string]string
	GenerateDockerfile() string
}

// DocsGeneratorInterface defines the interface for documentation generation
type DocsGeneratorInterface interface {
	FileGenerator
	GenerateReadme() string
	GenerateArchitecture() string
	GenerateAPI() string
	GenerateProduction() string
}

// InfrastructureGeneratorInterface defines the interface for infrastructure generation
type InfrastructureGeneratorInterface interface {
	FileGenerator
	GenerateDockerCompose() string
	GenerateRenderConfig() string
	GenerateAWSConfig() string
}

// ToolingGeneratorInterface defines the interface for development tooling
type ToolingGeneratorInterface interface {
	FileGenerator
	GenerateGitFiles() map[string]string
	GenerateVSCodeConfig() map[string]string
	GenerateCIConfig() string
}