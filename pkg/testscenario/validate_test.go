package testscenario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"hello-world", true},
		{"step-1", true},
		{"a", true},
		{"hello", true},
		{"Hello", false},
		{"hello_world", false},
		{"hello world", false},
		{"HELLO", false},
		{"", false},
		{"hello--world", true},
		{"-leading", true},
		{"trailing-", true},
		{"123", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isKebabCase(tt.input); got != tt.want {
				t.Errorf("isKebabCase(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateScenarioFile_valid(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "valid.test.md")
	content := `# Scenario: Valid Test

**Description:** A valid scenario.
**Tags:** unit

## check-output

` + "```bash\necho hello\n```\n"

	os.WriteFile(f, []byte(content), 0o644)
	errs := ValidateScenarioFile(f, "")
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateScenarioFile_missingTitle(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "no-title.test.md")
	os.WriteFile(f, []byte("Just some text\n\nNo heading.\n"), 0o644)

	errs := ValidateScenarioFile(f, "")
	if len(errs) == 0 {
		t.Fatal("expected errors for missing title")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "missing '# Scenario:'") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected missing title error, got %v", errs)
	}
}

func TestValidateScenarioFile_missingDescription(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "no-desc.test.md")
	content := `# Scenario: No Desc

## step-a

` + "```bash\necho hi\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateScenarioFile(f, "")
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "Description") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected missing description error, got %v", errs)
	}
}

func TestValidateScenarioFile_nonKebabStepName(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad-name.test.md")
	content := `# Scenario: Bad Name

**Description:** Test bad name.

## Bad_Name

` + "```bash\necho hi\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateScenarioFile(f, "")
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "not kebab-case") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected kebab-case error, got %v", errs)
	}
}

func TestValidateScenarioFile_parseError(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad-parse.test.md")
	// Missing language annotation on code block.
	content := `# Scenario: Parse Error

**Description:** Test parse error.

## step-a

` + "```\necho hi\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateScenarioFile(f, "")
	if len(errs) == 0 {
		t.Fatal("expected errors for missing language annotation")
	}
}

func TestValidateScenarioFile_ACRefResolves(t *testing.T) {
	// Create a spec directory with an AC file.
	dir := t.TempDir()
	specRoot := filepath.Join(dir, "spec")
	acsDir := filepath.Join(specRoot, "features", "my-feature", "_acs")
	os.MkdirAll(acsDir, 0o755)
	os.WriteFile(filepath.Join(acsDir, "my-ac.ac.md"), []byte("# AC: my-ac\n"), 0o644)

	// Create a scenario referencing the AC.
	scenarioDir := filepath.Join(dir, "tests")
	os.MkdirAll(scenarioDir, 0o755)
	scenarioFile := filepath.Join(scenarioDir, "ref.test.md")
	content := `# Scenario: AC Ref Test

**Description:** Tests AC ref resolution.

## check-ac

**ACs:**

| Feature | ACs |
| --- | --- |
| [my-feature](spec/features/my-feature) | [my-ac](my-ac.ac.md) |

` + "```bash\necho ok\n```\n"
	os.WriteFile(scenarioFile, []byte(content), 0o644)

	errs := ValidateScenarioFile(scenarioFile, specRoot)
	for _, e := range errs {
		if strings.Contains(e.Message, "AC reference") {
			t.Errorf("unexpected AC ref error: %v", e)
		}
	}
}

func TestValidateScenarioFile_ACRefMissing(t *testing.T) {
	dir := t.TempDir()
	specRoot := filepath.Join(dir, "spec")
	// Create _acs dir but without the referenced file.
	os.MkdirAll(filepath.Join(specRoot, "features", "my-feature", "_acs"), 0o755)

	scenarioDir := filepath.Join(dir, "tests")
	os.MkdirAll(scenarioDir, 0o755)
	scenarioFile := filepath.Join(scenarioDir, "missing-ref.test.md")
	content := `# Scenario: Missing AC

**Description:** Tests missing AC ref.

## check-ac

**ACs:**

| Feature | ACs |
| --- | --- |
| [my-feature](spec/features/my-feature) | [nonexistent](nonexistent.ac.md) |

` + "```bash\necho ok\n```\n"
	os.WriteFile(scenarioFile, []byte(content), 0o644)

	errs := ValidateScenarioFile(scenarioFile, specRoot)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "AC reference") || strings.Contains(e.Message, "does not exist") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected missing AC ref error, got %v", errs)
	}
}

func TestValidateACFile_valid(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "my-check.ac.md")
	content := `# AC: my-check

**Status:** implemented
**Feature:** [my-feature](../../)

## Description

This checks something.

## Inputs

| Name | Required | Description |
| --- | --- | --- |
| STEP_STDOUT | yes | stdout |

## Verification

` + "```bash\ntest -n \"$STEP_STDOUT\"\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateACFile(f)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateACFile_slugMismatch(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "wrong-name.ac.md")
	content := `# AC: different-slug

**Status:** planned
**Feature:** [f](../../)

## Description

Mismatch test.

## Inputs

None.
`
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateACFile(f)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "does not match filename") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected slug mismatch error, got %v", errs)
	}
}

func TestValidateACFile_invalidStatus(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad-status.ac.md")
	content := `# AC: bad-status

**Status:** unknown
**Feature:** [f](../../)

## Description

Bad status value.

## Inputs

None.
`
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateACFile(f)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "invalid status") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected invalid status error, got %v", errs)
	}
}

func TestValidateACFile_missingSections(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "missing-sections.ac.md")
	content := `# AC: missing-sections

**Status:** implemented
**Feature:** [f](../../)
`
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateACFile(f)
	descErr, inputErr, verErr := false, false, false
	for _, e := range errs {
		if strings.Contains(e.Message, "Description") {
			descErr = true
		}
		if strings.Contains(e.Message, "Inputs") {
			inputErr = true
		}
		if strings.Contains(e.Message, "Verification") {
			verErr = true
		}
	}
	if !descErr {
		t.Error("expected missing Description error")
	}
	if !inputErr {
		t.Error("expected missing Inputs error")
	}
	if !verErr {
		t.Error("expected missing Verification error (status=implemented)")
	}
}

func TestValidateACFile_plannedNoVerification(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "planned-ok.ac.md")
	content := `# AC: planned-ok

**Status:** planned
**Feature:** [f](../../)

## Description

This is planned.

## Inputs

None.
`
	os.WriteFile(f, []byte(content), 0o644)

	errs := ValidateACFile(f)
	for _, e := range errs {
		if strings.Contains(e.Message, "Verification") {
			t.Errorf("planned AC should not require Verification, got: %v", e)
		}
	}
}

func TestValidateACIndex_synced(t *testing.T) {
	dir := t.TempDir()
	acsDir := filepath.Join(dir, "_acs")
	os.MkdirAll(acsDir, 0o755)
	os.WriteFile(filepath.Join(acsDir, "check-a.ac.md"), []byte("# AC: check-a\n"), 0o644)
	os.WriteFile(filepath.Join(acsDir, "check-b.ac.md"), []byte("# AC: check-b\n"), 0o644)
	readme := `# ACs

| AC | Status |
| --- | --- |
| [check-a](check-a.ac.md) | planned |
| [check-b](check-b.ac.md) | planned |
`
	os.WriteFile(filepath.Join(acsDir, "README.md"), []byte(readme), 0o644)

	errs := ValidateACIndex(acsDir)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateACIndex_orphan(t *testing.T) {
	dir := t.TempDir()
	acsDir := filepath.Join(dir, "_acs")
	os.MkdirAll(acsDir, 0o755)
	os.WriteFile(filepath.Join(acsDir, "check-a.ac.md"), []byte("# AC: check-a\n"), 0o644)
	os.WriteFile(filepath.Join(acsDir, "orphan.ac.md"), []byte("# AC: orphan\n"), 0o644)
	readme := `# ACs

| AC | Status |
| --- | --- |
| [check-a](check-a.ac.md) | planned |
`
	os.WriteFile(filepath.Join(acsDir, "README.md"), []byte(readme), 0o644)

	errs := ValidateACIndex(acsDir)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "orphaned") && strings.Contains(e.Message, "orphan.ac.md") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected orphan error, got %v", errs)
	}
}

func TestValidateACIndex_phantom(t *testing.T) {
	dir := t.TempDir()
	acsDir := filepath.Join(dir, "_acs")
	os.MkdirAll(acsDir, 0o755)
	os.WriteFile(filepath.Join(acsDir, "check-a.ac.md"), []byte("# AC: check-a\n"), 0o644)
	readme := `# ACs

| AC | Status |
| --- | --- |
| [check-a](check-a.ac.md) | planned |
| [phantom](phantom.ac.md) | planned |
`
	os.WriteFile(filepath.Join(acsDir, "README.md"), []byte(readme), 0o644)

	errs := ValidateACIndex(acsDir)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "phantom") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected phantom error, got %v", errs)
	}
}

func TestValidateAll_directory(t *testing.T) {
	dir := t.TempDir()

	// Create a valid scenario.
	scenario := `# Scenario: All Test

**Description:** Integration test.

## step-one

` + "```bash\necho ok\n```\n"
	os.WriteFile(filepath.Join(dir, "good.test.md"), []byte(scenario), 0o644)

	// Create a valid AC.
	acsDir := filepath.Join(dir, "_acs")
	os.MkdirAll(acsDir, 0o755)
	ac := `# AC: check-good

**Status:** planned
**Feature:** [f](../../)

## Description

Good AC.

## Inputs

None.
`
	os.WriteFile(filepath.Join(acsDir, "check-good.ac.md"), []byte(ac), 0o644)
	readme := `# ACs

| AC | Status |
| --- | --- |
| [check-good](check-good.ac.md) | planned |
`
	os.WriteFile(filepath.Join(acsDir, "README.md"), []byte(readme), 0o644)

	result, err := ValidateAll(dir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasErrors() {
		t.Errorf("expected no errors, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.ScenariosCount != 1 {
		t.Errorf("expected 1 scenario, got %d", result.ScenariosCount)
	}
	if result.ACsCount != 1 {
		t.Errorf("expected 1 AC, got %d", result.ACsCount)
	}
}

func TestValidateAll_singleFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "single.test.md")
	content := `# Scenario: Single

**Description:** Single file test.

## do-it

` + "```bash\necho ok\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	result, err := ValidateAll(f, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasErrors() {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
	if result.ScenariosCount != 1 {
		t.Errorf("expected 1 scenario, got %d", result.ScenariosCount)
	}
}

func TestValidateAll_pathNotFound(t *testing.T) {
	_, err := ValidateAll("/nonexistent/path", "", 0)
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestFormatValidationResult_noErrors(t *testing.T) {
	result := &ValidationResult{ScenariosCount: 3, ACsCount: 5}
	output := FormatValidationResult(result)
	if !strings.Contains(output, "3 scenarios") || !strings.Contains(output, "5 ACs") {
		t.Errorf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "no errors") {
		t.Errorf("expected 'no errors' in output: %s", output)
	}
}

func TestFormatValidationResult_withErrors(t *testing.T) {
	result := &ValidationResult{
		ScenariosCount: 1,
		Errors: []ValidationError{
			{File: "a.test.md", Line: 5, Message: "bad step"},
			{File: "a.test.md", Line: 10, Message: "another"},
			{File: "b.ac.md", Message: "missing section"},
		},
	}
	output := FormatValidationResult(result)
	if !strings.Contains(output, "a.test.md") || !strings.Contains(output, "b.ac.md") {
		t.Errorf("expected file names in output: %s", output)
	}
	if !strings.Contains(output, "3 errors") {
		t.Errorf("expected error count in output: %s", output)
	}
}

func TestValidateAll_failFast_stopsEarly(t *testing.T) {
	dir := t.TempDir()

	// Create multiple bad files so there are many potential errors.
	for _, name := range []string{"a.test.md", "b.test.md", "c.test.md"} {
		os.WriteFile(filepath.Join(dir, name), []byte("no title here\n"), 0o644)
	}

	// Without limit: collects all.
	all, _ := ValidateAll(dir, "", 0)
	if len(all.Errors) < 3 {
		t.Fatalf("expected at least 3 errors without limit, got %d", len(all.Errors))
	}
	if all.Truncated {
		t.Error("should not be truncated without limit")
	}

	// With limit=1: stops after 1 error.
	limited, _ := ValidateAll(dir, "", 1)
	if len(limited.Errors) != 1 {
		t.Errorf("expected 1 error with --fail-fast=1, got %d", len(limited.Errors))
	}
	if !limited.Truncated {
		t.Error("expected Truncated=true with --fail-fast=1")
	}

	// With limit=2: stops after 2.
	limited2, _ := ValidateAll(dir, "", 2)
	if len(limited2.Errors) != 2 {
		t.Errorf("expected 2 errors with --fail-fast=2, got %d", len(limited2.Errors))
	}
	if !limited2.Truncated {
		t.Error("expected Truncated=true with --fail-fast=2")
	}
}

func TestFormatValidationResult_truncated(t *testing.T) {
	result := &ValidationResult{
		ScenariosCount: 5,
		Truncated:      true,
		Errors: []ValidationError{
			{File: "a.test.md", Line: 1, Message: "error"},
		},
	}
	output := FormatValidationResult(result)
	if !strings.Contains(output, "--fail-fast") {
		t.Errorf("expected truncation note in output: %s", output)
	}
}
