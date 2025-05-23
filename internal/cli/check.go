package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Run guard rails and compliance checks",
		Long: `Validates code against domain-specific guard rails:

- Security best practices
- Compliance requirements  
- Code quality standards
- Domain-specific patterns

Supports different enforcement levels:
- Advisory: warnings and suggestions
- Strict: blocks with violations
- CI: optimized for continuous integration`,
		RunE: runCheck,
	}

	cmd.Flags().Bool("enforce", false, "Strict enforcement mode")
	cmd.Flags().Bool("ci", false, "CI/CD mode")
	cmd.Flags().String("format", "text", "Output format (text, json)")

	return cmd
}

func runCheck(cmd *cobra.Command, args []string) error {
	enforce, _ := cmd.Flags().GetBool("enforce")
	ci, _ := cmd.Flags().GetBool("ci")
	format, _ := cmd.Flags().GetString("format")
	
	color.Blue("🛡️ Running guard rails validation...")
	
	if enforce {
		color.Yellow("⚠️ Strict enforcement mode - will block on violations")
	}
	
	if ci {
		color.Cyan("🔄 CI/CD mode - optimized output")
	}

	return performGuardRailsCheck(enforce, ci, format)
}

type CheckResult struct {
	Category string
	Rule     string
	File     string
	Line     int
	Severity string
	Message  string
	Passed   bool
}

func performGuardRailsCheck(enforce, ci bool, format string) error {
	var results []CheckResult
	
	// Security checks
	results = append(results, checkSecurityPatterns()...)
	
	// Code quality checks
	results = append(results, checkCodeQuality()...)
	
	// Domain-specific checks
	results = append(results, checkDomainPatterns()...)
	
	// Compliance checks
	results = append(results, checkCompliance()...)

	// Display results
	if format == "json" {
		return displayResultsJSON(results)
	}
	
	return displayResultsText(results, enforce, ci)
}

func checkSecurityPatterns() []CheckResult {
	var results []CheckResult
	
	// Check for common security anti-patterns
	patterns := map[string]string{
		`(?i)(password|secret|key|token)\s*[:=]\s*["'][^"']*["']`: "Hardcoded secrets detected",
		`(?i)http://`:                                              "Insecure HTTP protocol used",
		`eval\s*\(`:                                               "Dangerous eval() usage",
		`innerHTML\s*=`:                                           "Potential XSS vulnerability",
		`\.exec\s*\(`:                                             "Command execution detected",
	}
	
	for pattern, message := range patterns {
		matches := findInFiles(pattern, []string{"*.js", "*.ts", "*.go", "*.py", "*.java"})
		for _, match := range matches {
			results = append(results, CheckResult{
				Category: "Security",
				Rule:     "Security Pattern",
				File:     match.File,
				Line:     match.Line,
				Severity: "HIGH",
				Message:  message,
				Passed:   false,
			})
		}
	}
	
	// Add passed checks
	if len(results) == 0 {
		results = append(results, CheckResult{
			Category: "Security",
			Rule:     "Security Patterns",
			Severity: "INFO",
			Message:  "No security anti-patterns detected",
			Passed:   true,
		})
	}
	
	return results
}

func checkCodeQuality() []CheckResult {
	var results []CheckResult
	
	// Check for code quality issues
	patterns := map[string]string{
		`console\.log|print\(`:        "Debug statements left in code",
		`TODO|FIXME|HACK`:             "Technical debt markers found",
		`(?m)^.{120,}$`:               "Lines exceeding 120 characters",
		`(?m)^\s*$\n\s*$\n\s*$`:       "Multiple consecutive empty lines",
	}
	
	for pattern, message := range patterns {
		matches := findInFiles(pattern, []string{"*.js", "*.ts", "*.go", "*.py", "*.java"})
		for _, match := range matches {
			severity := "MEDIUM"
			if strings.Contains(message, "Debug statements") {
				severity = "LOW"
			}
			
			results = append(results, CheckResult{
				Category: "Quality",
				Rule:     "Code Quality",
				File:     match.File,
				Line:     match.Line,
				Severity: severity,
				Message:  message,
				Passed:   false,
			})
		}
	}
	
	return results
}

func checkDomainPatterns() []CheckResult {
	var results []CheckResult
	
	// Check for domain-specific patterns (example: fintech)
	domain := detectDomain()
	
	switch domain {
	case "fintech":
		// Fintech-specific checks
		patterns := map[string]string{
			`(?i)(balance|amount|price)\s*[:=]\s*["'][^"']*["']`: "Financial data as string (should be decimal)",
			`Math\.random\(\)`:                                  "Non-cryptographic random for financial data",
			`parseFloat|parseInt`:                               "Unsafe number parsing for financial values",
		}
		
		for pattern, message := range patterns {
			matches := findInFiles(pattern, []string{"*.js", "*.ts", "*.go", "*.py"})
			for _, match := range matches {
				results = append(results, CheckResult{
					Category: "Domain",
					Rule:     "Fintech Compliance",
					File:     match.File,
					Line:     match.Line,
					Severity: "HIGH",
					Message:  message,
					Passed:   false,
				})
			}
		}
	}
	
	if len(results) == 0 {
		results = append(results, CheckResult{
			Category: "Domain",
			Rule:     "Domain Patterns",
			Severity: "INFO",
			Message:  fmt.Sprintf("Domain-specific checks passed (%s)", domain),
			Passed:   true,
		})
	}
	
	return results
}

func checkCompliance() []CheckResult {
	var results []CheckResult
	
	// Check for compliance requirements
	requiredFiles := []string{
		"README.md",
		"LICENSE",
		".gitignore",
	}
	
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			results = append(results, CheckResult{
				Category: "Compliance",
				Rule:     "Required Files",
				File:     file,
				Severity: "MEDIUM",
				Message:  fmt.Sprintf("Required file %s is missing", file),
				Passed:   false,
			})
		}
	}
	
	// Check for proper license headers in source files
	if hasSourceFiles() && !hasLicenseHeaders() {
		results = append(results, CheckResult{
			Category: "Compliance",
			Rule:     "License Headers",
			Severity: "LOW",
			Message:  "Source files missing license headers",
			Passed:   false,
		})
	}
	
	if len(results) == 0 {
		results = append(results, CheckResult{
			Category: "Compliance",
			Rule:     "Compliance Check",
			Severity: "INFO",
			Message:  "All compliance requirements met",
			Passed:   true,
		})
	}
	
	return results
}

type FileMatch struct {
	File string
	Line int
}

func findInFiles(pattern string, extensions []string) []FileMatch {
	var matches []FileMatch
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return matches
	}
	
	for _, ext := range extensions {
		files, _ := filepath.Glob(ext)
		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			
			lines := strings.Split(string(content), "\n")
			for i, line := range lines {
				if regex.MatchString(line) {
					matches = append(matches, FileMatch{
						File: file,
						Line: i + 1,
					})
				}
			}
		}
	}
	
	return matches
}

func detectDomain() string {
	// Try to detect domain from config file
	if data, err := os.ReadFile(".codekeeper/env.local"); err == nil {
		content := string(data)
		if strings.Contains(content, "fintech") {
			return "fintech"
		}
		if strings.Contains(content, "healthcare") {
			return "healthcare"
		}
		if strings.Contains(content, "ecommerce") {
			return "ecommerce"
		}
	}
	
	// Try to detect from project structure
	if _, err := os.Stat("package.json"); err == nil {
		if data, err := os.ReadFile("package.json"); err == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "finance") || strings.Contains(content, "payment") {
				return "fintech"
			}
		}
	}
	
	return "general"
}

func hasSourceFiles() bool {
	extensions := []string{"*.js", "*.ts", "*.go", "*.py", "*.java"}
	for _, ext := range extensions {
		if files, _ := filepath.Glob(ext); len(files) > 0 {
			return true
		}
	}
	return false
}

func hasLicenseHeaders() bool {
	// Simple check for license headers in source files
	files, _ := filepath.Glob("*.go")
	if len(files) == 0 {
		return true // No Go files to check
	}
	
	for _, file := range files[:1] { // Check first file only
		if content, err := os.ReadFile(file); err == nil {
			firstLines := strings.Join(strings.Split(string(content), "\n")[:5], "\n")
			if strings.Contains(firstLines, "Copyright") || strings.Contains(firstLines, "License") {
				return true
			}
		}
	}
	return false
}

func displayResultsText(results []CheckResult, enforce, ci bool) error {
	passed := 0
	warnings := 0
	errors := 0
	
	for _, result := range results {
		if result.Passed {
			passed++
			if !ci {
				color.Green("✓ [%s] %s", result.Category, result.Message)
			}
		} else {
			switch result.Severity {
			case "HIGH":
				errors++
				color.Red("✗ [%s] %s", result.Category, result.Message)
				if result.File != "" {
					fmt.Printf("  → %s:%d\n", result.File, result.Line)
				}
			case "MEDIUM":
				warnings++
				color.Yellow("⚠ [%s] %s", result.Category, result.Message)
				if result.File != "" {
					fmt.Printf("  → %s:%d\n", result.File, result.Line)
				}
			default:
				warnings++
				color.Cyan("ℹ [%s] %s", result.Category, result.Message)
			}
		}
	}
	
	fmt.Println()
	color.Blue("📊 Summary:")
	color.Green("  ✓ Passed: %d", passed)
	color.Yellow("  ⚠ Warnings: %d", warnings)
	color.Red("  ✗ Errors: %d", errors)
	
	if errors > 0 && enforce {
		color.Red("🚫 Guard rails enforcement failed - %d errors found", errors)
		return fmt.Errorf("guard rails violations detected")
	}
	
	if errors == 0 && warnings == 0 {
		color.Green("🎉 All guard rails checks passed!")
	} else if errors == 0 {
		color.Yellow("⚠️  Some warnings found, but no blocking errors")
	}
	
	return nil
}

func displayResultsJSON(results []CheckResult) error {
	// Simple JSON output for CI/CD integration
	fmt.Println("[")
	for i, result := range results {
		fmt.Printf(`  {
    "category": "%s",
    "rule": "%s",
    "file": "%s",
    "line": %d,
    "severity": "%s",
    "message": "%s",
    "passed": %t
  }`, result.Category, result.Rule, result.File, result.Line, result.Severity, result.Message, result.Passed)
		
		if i < len(results)-1 {
			fmt.Println(",")
		} else {
			fmt.Println()
		}
	}
	fmt.Println("]")
	return nil
}