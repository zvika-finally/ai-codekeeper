package generator

import "fmt"

// DEPRECATED: Old monolithic generator
// Use NewModular() for the new modular, standards-focused generator

// Generator represents the old monolithic generator (DEPRECATED)
type Generator struct {
	spec *ProjectSpec
}

// New creates the old monolithic generator (DEPRECATED)
// Use generator.NewModular() instead for modern, standards-focused generation
func New(spec *ProjectSpec) *Generator {
	return &Generator{spec: spec}
}

// Generate runs the old monolithic generator (DEPRECATED)
func (g *Generator) Generate() error {
	return fmt.Errorf("DEPRECATED: Old monolithic generator is no longer supported. Use generator.NewModular() instead")
}

// Stub methods for compatibility (DEPRECATED)
func (g *Generator) getDomainGuardRails() []string {
	return []string{}
}

