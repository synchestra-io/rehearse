package testscenario

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationError represents a single validation problem.
type ValidationError struct {
	File    string
	Line    int // 0 means unknown
	Message string
}

func (e ValidationError) String() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s: line %d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.File, e.Message)
}

// ValidationResult holds all errors from a validation run.
type ValidationResult struct {
	Errors         []ValidationError
	ScenariosCount int
	ACsCount       int
}

// HasErrors returns true if any validation errors were found.
func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// addError appends a validation error.
func (r *ValidationResult) addError(file string, line int, msg string) {
	r.Errors = append(r.Errors, ValidationError{File: file, Line: line, Message: msg})
}

// ValidateAll discovers and validates all scenario and AC files under path.
// If path is a file, validates that single file.
// If path is a directory, recursively discovers .test.md and .ac.md files.
func ValidateAll(path, specRoot string) (*ValidationResult, error) {
	result := &ValidationResult{}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path not found: %s", path)
	}

	if !info.IsDir() {
		// Single file — detect type by extension.
		if strings.HasSuffix(path, ".test.md") {
			result.ScenariosCount++
			errs := ValidateScenarioFile(path, specRoot)
			result.Errors = append(result.Errors, errs...)
		} else if strings.HasSuffix(path, ".ac.md") {
			result.ACsCount++
			errs := ValidateACFile(path)
			result.Errors = append(result.Errors, errs...)
		}
		return result, nil
	}

	// Directory — walk and discover files.
	_ = filepath.Walk(path, func(p string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil || fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".test.md") {
			result.ScenariosCount++
			errs := ValidateScenarioFile(p, specRoot)
			result.Errors = append(result.Errors, errs...)
		} else if strings.HasSuffix(p, ".ac.md") {
			result.ACsCount++
			errs := ValidateACFile(p)
			result.Errors = append(result.Errors, errs...)
		}
		return nil
	})

	// Validate AC index synchronization for every _acs/ directory found.
	_ = filepath.Walk(path, func(p string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil || !fi.IsDir() || fi.Name() != "_acs" {
			return nil
		}
		errs := ValidateACIndex(p)
		result.Errors = append(result.Errors, errs...)
		return nil
	})

	return result, nil
}

// ValidateScenarioFile validates a single .test.md scenario file.
func ValidateScenarioFile(path, specRoot string) []ValidationError {
	var errs []ValidationError

	data, err := os.ReadFile(path)
	if err != nil {
		return []ValidationError{{File: path, Message: fmt.Sprintf("cannot read file: %v", err)}}
	}

	text := string(data)
	lines := strings.Split(text, "\n")

	// Check title.
	titleLine := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "# Scenario:") {
			titleLine = i
			title := strings.TrimSpace(strings.TrimPrefix(line, "# Scenario:"))
			if title == "" {
				errs = append(errs, ValidationError{File: path, Line: i + 1, Message: "scenario title is empty"})
			}
			break
		}
	}
	if titleLine < 0 {
		errs = append(errs, ValidationError{File: path, Line: 1, Message: "missing '# Scenario:' heading"})
		return errs // Can't continue without a title
	}

	// Check description metadata.
	hasDescription := false
	for i := titleLine + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		if strings.HasPrefix(trimmed, "**Description:**") {
			hasDescription = true
			desc := strings.TrimSpace(strings.TrimPrefix(trimmed, "**Description:**"))
			if desc == "" {
				errs = append(errs, ValidationError{File: path, Line: i + 1, Message: "description is empty"})
			}
			break
		}
	}
	if !hasDescription {
		errs = append(errs, ValidationError{File: path, Line: titleLine + 1, Message: "missing '**Description:**' metadata"})
	}

	// Try full parse — this catches duplicate step names, missing code blocks,
	// missing language annotations, bad depends-on refs, etc.
	scenario, parseErr := ParseScenario(data)
	if parseErr != nil {
		line := findErrorLine(lines, parseErr.Error())
		errs = append(errs, ValidationError{File: path, Line: line, Message: parseErr.Error()})
		return errs // Can't do further validation without a parse
	}

	// Validate step names are kebab-case.
	for _, step := range scenario.Steps {
		if !isKebabCase(step.Name) {
			line := findHeadingLine(lines, step.Name)
			errs = append(errs, ValidationError{
				File:    path,
				Line:    line,
				Message: fmt.Sprintf("step name %q is not kebab-case", step.Name),
			})
		}
	}

	// Validate reserved steps don't have disallowed metadata.
	for _, step := range scenario.Steps {
		if step.Name == "Setup" || step.Name == "Teardown" {
			continue // Setup/Teardown are parsed separately, not in Steps
		}
	}

	// Cross-reference validation: AC refs resolve to files.
	if specRoot != "" {
		for _, step := range scenario.Steps {
			for _, acRef := range step.ACs {
				refErrs := validateACRef(path, specRoot, step.Name, acRef, lines)
				errs = append(errs, refErrs...)
			}
		}
	}

	// Cross-reference validation: Include refs resolve to files.
	for _, step := range scenario.Steps {
		if step.Include != "" {
			includePath := step.Include
			if !filepath.IsAbs(includePath) {
				includePath = filepath.Join(filepath.Dir(path), includePath)
			}
			if _, err := os.Stat(includePath); os.IsNotExist(err) {
				line := findHeadingLine(lines, step.Name)
				errs = append(errs, ValidationError{
					File:    path,
					Line:    line,
					Message: fmt.Sprintf("include reference %q does not exist", step.Include),
				})
			}
		}
	}

	return errs
}

// validateACRef checks that an AC reference in a scenario resolves to actual files.
func validateACRef(scenarioPath, specRoot, stepName string, acRef ACRef, lines []string) []ValidationError {
	var errs []ValidationError
	featurePath := acRef.FeaturePath
	acsDir := filepath.Join(specRoot, "features", filepath.FromSlash(featurePath), "_acs")

	if acRef.ACs == "*" {
		// Wildcard: _acs/ directory must exist and contain at least one .ac.md file.
		entries, err := os.ReadDir(acsDir)
		if err != nil {
			line := findHeadingLine(lines, stepName)
			errs = append(errs, ValidationError{
				File:    scenarioPath,
				Line:    line,
				Message: fmt.Sprintf("wildcard AC reference to %q: _acs/ directory does not exist at %s", featurePath, acsDir),
			})
			return errs
		}
		hasAC := false
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".ac.md") {
				hasAC = true
				break
			}
		}
		if !hasAC {
			line := findHeadingLine(lines, stepName)
			errs = append(errs, ValidationError{
				File:    scenarioPath,
				Line:    line,
				Message: fmt.Sprintf("wildcard AC reference to %q: _acs/ directory at %s contains no .ac.md files", featurePath, acsDir),
			})
		}
		return errs
	}

	// Specific ACs: each must resolve to a file.
	slugs := strings.Split(acRef.ACs, ",")
	for _, slug := range slugs {
		slug = strings.TrimSpace(slug)
		// Strip markdown link syntax: [slug](path) → slug
		if idx := strings.Index(slug, "]"); idx > 0 && slug[0] == '[' {
			slug = slug[1:idx]
		}
		acPath := filepath.Join(acsDir, slug+".ac.md")
		if _, err := os.Stat(acPath); os.IsNotExist(err) {
			line := findHeadingLine(lines, stepName)
			errs = append(errs, ValidationError{
				File:    scenarioPath,
				Line:    line,
				Message: fmt.Sprintf("AC reference %q does not exist: %s", slug, acPath),
			})
		}
	}
	return errs
}

// ValidateACFile validates a single .ac.md acceptance criteria file.
func ValidateACFile(path string) []ValidationError {
	var errs []ValidationError

	data, err := os.ReadFile(path)
	if err != nil {
		return []ValidationError{{File: path, Message: fmt.Sprintf("cannot read file: %v", err)}}
	}

	text := string(data)
	lines := strings.Split(text, "\n")

	// Extract expected slug from filename.
	base := filepath.Base(path)
	expectedSlug := strings.TrimSuffix(base, ".ac.md")

	// Check title.
	titleLine := -1
	actualSlug := ""
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# AC:") {
			titleLine = i
			actualSlug = strings.TrimSpace(strings.TrimPrefix(trimmed, "# AC:"))
			if actualSlug == "" {
				errs = append(errs, ValidationError{File: path, Line: i + 1, Message: "AC title is empty"})
			}
			break
		}
	}
	if titleLine < 0 {
		errs = append(errs, ValidationError{File: path, Line: 1, Message: "missing '# AC:' heading"})
		return errs
	}

	// Check slug matches filename.
	if actualSlug != "" && actualSlug != expectedSlug {
		errs = append(errs, ValidationError{
			File:    path,
			Line:    titleLine + 1,
			Message: fmt.Sprintf("slug %q does not match filename %q", actualSlug, base),
		})
	}

	// Check Status field.
	hasStatus := false
	status := ""
	hasFeature := false
	for i := titleLine + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		if strings.HasPrefix(trimmed, "**Status:**") {
			hasStatus = true
			status = strings.TrimSpace(strings.TrimPrefix(trimmed, "**Status:**"))
			switch status {
			case "planned", "wip", "implemented", "deprecated":
				// Valid
			default:
				errs = append(errs, ValidationError{
					File:    path,
					Line:    i + 1,
					Message: fmt.Sprintf("invalid status %q (expected: planned, wip, implemented, deprecated)", status),
				})
			}
		}
		if strings.HasPrefix(trimmed, "**Feature:**") {
			hasFeature = true
		}
	}
	if !hasStatus {
		errs = append(errs, ValidationError{File: path, Line: titleLine + 2, Message: "missing '**Status:**' field"})
	}
	if !hasFeature {
		errs = append(errs, ValidationError{File: path, Line: titleLine + 2, Message: "missing '**Feature:**' field"})
	}

	// Check required sections by scanning for ## headings.
	hasDescription := false
	hasInputs := false
	hasVerification := false
	verificationHasCode := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## Description" {
			hasDescription = true
			// Check that description section is non-empty.
			isEmpty := true
			for j := i + 1; j < len(lines); j++ {
				t := strings.TrimSpace(lines[j])
				if strings.HasPrefix(t, "## ") {
					break
				}
				if t != "" {
					isEmpty = false
					break
				}
			}
			if isEmpty {
				errs = append(errs, ValidationError{File: path, Line: i + 1, Message: "Description section is empty"})
			}
		}
		if trimmed == "## Inputs" {
			hasInputs = true
		}
		if trimmed == "## Verification" {
			hasVerification = true
			// Check for code block in verification section.
			for j := i + 1; j < len(lines); j++ {
				t := strings.TrimSpace(lines[j])
				if strings.HasPrefix(t, "## ") {
					break
				}
				if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "````") {
					verificationHasCode = true
					break
				}
			}
		}
	}

	if !hasDescription {
		errs = append(errs, ValidationError{File: path, Message: "missing '## Description' section"})
	}
	if !hasInputs {
		errs = append(errs, ValidationError{File: path, Message: "missing '## Inputs' section"})
	}
	if !hasVerification && (status == "wip" || status == "implemented") {
		errs = append(errs, ValidationError{File: path, Message: fmt.Sprintf("missing '## Verification' section (required for status %q)", status)})
	}
	if hasVerification && !verificationHasCode && (status == "wip" || status == "implemented") {
		errs = append(errs, ValidationError{File: path, Message: fmt.Sprintf("Verification section has no code block (required for status %q)", status)})
	}

	// Try full parse to catch code block annotation issues.
	_, parseErr := ParseACFile(data, expectedSlug)
	if parseErr != nil {
		line := findErrorLine(lines, parseErr.Error())
		errs = append(errs, ValidationError{File: path, Line: line, Message: parseErr.Error()})
	}

	return errs
}

// ValidateACIndex checks that an _acs/ directory's README.md and .ac.md files are in sync.
func ValidateACIndex(acsDir string) []ValidationError {
	var errs []ValidationError
	readmePath := filepath.Join(acsDir, "README.md")

	// Collect actual .ac.md files on disk.
	filesOnDisk := make(map[string]bool)
	entries, err := os.ReadDir(acsDir)
	if err != nil {
		return nil // Directory doesn't exist or unreadable — not an index sync error
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".ac.md") {
			filesOnDisk[e.Name()] = true
		}
	}

	// Parse README.md for AC table entries.
	readmeData, err := os.ReadFile(readmePath)
	if err != nil {
		if len(filesOnDisk) > 0 {
			errs = append(errs, ValidationError{
				File:    acsDir,
				Message: fmt.Sprintf("_acs/ directory contains %d .ac.md files but has no README.md", len(filesOnDisk)),
			})
		}
		return errs
	}

	// Extract AC filenames referenced in the README table.
	readmeLines := strings.Split(string(readmeData), "\n")
	referencedFiles := make(map[string]bool)
	for _, line := range readmeLines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") || strings.HasPrefix(trimmed, "| AC") || strings.Contains(trimmed, "---") {
			continue
		}
		// Extract filename from markdown link in table cell: [slug](filename.ac.md)
		cells := splitTableRow(trimmed)
		if len(cells) < 1 {
			continue
		}
		cell := cells[0]
		_, url := parseMarkdownLink(cell)
		if url != "" && strings.HasSuffix(url, ".ac.md") {
			referencedFiles[url] = true
		}
	}

	// Check for orphaned files (on disk but not in README).
	for file := range filesOnDisk {
		if !referencedFiles[file] {
			errs = append(errs, ValidationError{
				File:    readmePath,
				Message: fmt.Sprintf("orphaned AC file: %s exists on disk but is not listed in README.md", file),
			})
		}
	}

	// Check for phantom entries (in README but no file on disk).
	for file := range referencedFiles {
		if !filesOnDisk[file] {
			errs = append(errs, ValidationError{
				File:    readmePath,
				Message: fmt.Sprintf("phantom AC entry: %s listed in README.md but file does not exist", file),
			})
		}
	}

	return errs
}

// isKebabCase returns true if s is lowercase, hyphen-separated, with no spaces/underscores/capitals.
func isKebabCase(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return false
		}
		if c == '_' || c == ' ' {
			return false
		}
		if !(c >= 'a' && c <= 'z') && !(c >= '0' && c <= '9') && c != '-' {
			return false
		}
	}
	return true
}

// findHeadingLine finds the line number (1-indexed) of a ## heading with the given name.
func findHeadingLine(lines []string, name string) int {
	target := "## " + name
	for i, line := range lines {
		if strings.TrimSpace(line) == target {
			return i + 1
		}
	}
	return 0
}

// findErrorLine tries to extract a line number from a parse error message,
// or returns 0 if none found.
func findErrorLine(lines []string, errMsg string) int {
	// Look for "line N" in the error message.
	idx := strings.Index(errMsg, "line ")
	if idx < 0 {
		return 0
	}
	numStr := ""
	for _, c := range errMsg[idx+5:] {
		if c >= '0' && c <= '9' {
			numStr += string(c)
		} else {
			break
		}
	}
	if numStr == "" {
		return 0
	}
	n := 0
	for _, c := range numStr {
		n = n*10 + int(c-'0')
	}
	return n
}

// FormatValidationResult formats validation results as human-readable text.
func FormatValidationResult(result *ValidationResult) string {
	if !result.HasErrors() {
		return fmt.Sprintf("Validated %d scenarios, %d ACs — no errors.\n",
			result.ScenariosCount, result.ACsCount)
	}

	var sb strings.Builder

	// Group errors by file.
	fileOrder := make([]string, 0)
	byFile := make(map[string][]ValidationError)
	for _, e := range result.Errors {
		if _, ok := byFile[e.File]; !ok {
			fileOrder = append(fileOrder, e.File)
		}
		byFile[e.File] = append(byFile[e.File], e)
	}

	for _, file := range fileOrder {
		sb.WriteString(file)
		sb.WriteString("\n")
		for _, e := range byFile[file] {
			if e.Line > 0 {
				sb.WriteString(fmt.Sprintf("  line %d: %s\n", e.Line, e.Message))
			} else {
				sb.WriteString(fmt.Sprintf("  %s\n", e.Message))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("%d files, %d errors\n", len(fileOrder), len(result.Errors)))
	return sb.String()
}
