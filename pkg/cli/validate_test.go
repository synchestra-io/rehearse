package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateCommand_validScenario(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "good.test.md")
	content := `# Scenario: Good

**Description:** A valid scenario.

## step-a

` + "```bash\necho ok\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	cmd := ValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{f})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("no errors")) {
		t.Errorf("expected 'no errors' in output, got: %s", buf.String())
	}
}

func TestValidateCommand_invalidScenario(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad.test.md")
	// Missing language annotation.
	content := `# Scenario: Bad

**Description:** Invalid scenario.

## step-a

` + "```\necho hi\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	cmd := ValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{f})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid scenario")
	}
	// Check exit code is 1.
	if ce, ok := err.(*CommandError); ok {
		if ce.Code != 1 {
			t.Errorf("expected exit code 1, got %d", ce.Code)
		}
	} else {
		t.Errorf("expected CommandError, got %T: %v", err, err)
	}
}

func TestValidateCommand_pathNotFound(t *testing.T) {
	cmd := ValidateCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"/nonexistent/path/to/file.test.md"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
	if ce, ok := err.(*CommandError); ok {
		if ce.Code != 3 {
			t.Errorf("expected exit code 3, got %d", ce.Code)
		}
	} else {
		t.Errorf("expected CommandError, got %T: %v", err, err)
	}
}

func TestValidateCommand_directory(t *testing.T) {
	dir := t.TempDir()

	// Valid scenario.
	content := `# Scenario: Dir Test

**Description:** Directory validation.

## check-it

` + "```bash\necho ok\n```\n"
	os.WriteFile(filepath.Join(dir, "ok.test.md"), []byte(content), 0o644)

	cmd := ValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{dir})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateCommand_specRootFlag(t *testing.T) {
	dir := t.TempDir()

	// Create spec root with an AC.
	specRoot := filepath.Join(dir, "spec")
	acsDir := filepath.Join(specRoot, "features", "my-feature", "_acs")
	os.MkdirAll(acsDir, 0o755)
	ac := `# AC: my-ac

**Status:** planned
**Feature:** [my-feature](../../)

## Description

Test AC.

## Inputs

None.
`
	os.WriteFile(filepath.Join(acsDir, "my-ac.ac.md"), []byte(ac), 0o644)

	// Scenario referencing the AC.
	testsDir := filepath.Join(dir, "tests")
	os.MkdirAll(testsDir, 0o755)
	scenario := `# Scenario: Spec Root Test

**Description:** Tests spec-root flag.

## check-ac

**ACs:**

| Feature | ACs |
| --- | --- |
| [my-feature](spec/features/my-feature) | [my-ac](my-ac.ac.md) |

` + "```bash\necho ok\n```\n"
	os.WriteFile(filepath.Join(testsDir, "ref.test.md"), []byte(scenario), 0o644)

	cmd := ValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{testsDir, "--spec-root", specRoot})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("expected no error with valid AC ref, got: %v", err)
	}
}

func TestValidateCommand_validACFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "good-ac.ac.md")
	content := `# AC: good-ac

**Status:** implemented
**Feature:** [f](../../)

## Description

Good AC.

## Inputs

| Name | Required | Description |
| --- | --- | --- |
| STEP_STDOUT | yes | stdout |

## Verification

` + "```bash\ntest -n \"$STEP_STDOUT\"\n```\n"
	os.WriteFile(f, []byte(content), 0o644)

	cmd := ValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{f})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateCommand_invalidACFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad-ac.ac.md")
	content := `# AC: wrong-slug

**Status:** bogus
`
	os.WriteFile(f, []byte(content), 0o644)

	cmd := ValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{f})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid AC")
	}
	if ce, ok := err.(*CommandError); ok {
		if ce.Code != 1 {
			t.Errorf("expected exit code 1, got %d", ce.Code)
		}
	}
}
