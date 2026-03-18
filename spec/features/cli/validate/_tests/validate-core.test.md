# Scenario: Validate core behaviors

**Description:** Integration test of `rehearse validate` — scenario structure validation, AC structure validation, cross-reference resolution, AC index synchronization, no-execution guarantee, exit code contract, and self-referential dogfooding against the real spec tree. This scenario is executed by the runner itself; the final step validates the spec tree that contains this very test file.
**Tags:** integration, cli, validate

## Setup

````bash
BINARY_PATH="${BINARY_PATH:-$(go env GOPATH)/bin/rehearse}"
SPEC_ROOT="$(git rev-parse --show-toplevel)/spec"
FIXTURE_DIR=$(mktemp -d)

# Build the binary
cd "$(git rev-parse --show-toplevel)"
go build -o "$BINARY_PATH" .

# ── Valid scenario fixture ──────────────────────────────────────────
cat > "$FIXTURE_DIR/valid-scenario.test.md" << 'SCENARIO'
# Scenario: Fixture valid

**Description:** A minimal valid scenario used as a validation fixture.
**Tags:** fixture

## do-something

```bash
echo "hello from valid scenario"
```
SCENARIO

# ── Malformed scenario fixture (bare code fence, no language annotation) ──
cat > "$FIXTURE_DIR/malformed-scenario.test.md" << 'SCENARIO'
# Scenario: Fixture malformed

**Description:** Has a bare code fence without a language annotation.

## bare-fence

```
echo "this fence has no language"
```
SCENARIO

# ── Valid AC fixture ────────────────────────────────────────────────
mkdir -p "$FIXTURE_DIR/features/fixture-feature/_acs"

cat > "$FIXTURE_DIR/features/fixture-feature/_acs/always-pass.ac.md" << 'AC'
# AC: always-pass

**Status:** implemented
**Feature:** [fixture-feature](../README.md)

## Description

A minimal valid AC used as a validation fixture.

## Inputs

| Name | Required | Description |
|---|---|---|

## Verification

```bash
exit 0
```

## Scenarios

(None yet.)
AC

# ── Malformed AC fixture (missing Status field) ────────────────────
cat > "$FIXTURE_DIR/features/fixture-feature/_acs/bad-ac.ac.md" << 'AC'
# AC: bad-ac

**Feature:** [fixture-feature](../README.md)

## Description

This AC is missing the required Status field.

## Inputs

| Name | Required | Description |
|---|---|---|

## Verification

```bash
exit 0
```

## Scenarios

(None yet.)
AC

# ── Synced AC index for the valid AC ───────────────────────────────
cat > "$FIXTURE_DIR/features/fixture-feature/_acs/README.md" << 'README'
# Acceptance Criteria: fixture-feature

| AC | Description | Status |
|---|---|---|
| [always-pass](always-pass.ac.md) | Always passes | implemented |
| [bad-ac](bad-ac.ac.md) | Missing status field | implemented |
README

cat > "$FIXTURE_DIR/features/fixture-feature/README.md" << 'README'
# Feature: fixture-feature

**Status:** In Progress

## Summary

Fixture feature for validation testing.

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [always-pass](_acs/always-pass.ac.md) | Always passes | implemented |
| [bad-ac](_acs/bad-ac.ac.md) | Missing status field | implemented |
README

# ── Phantom/orphan AC index fixture ────────────────────────────────
mkdir -p "$FIXTURE_DIR/features/sync-feature/_acs"

cat > "$FIXTURE_DIR/features/sync-feature/_acs/real-ac.ac.md" << 'AC'
# AC: real-ac

**Status:** implemented
**Feature:** [sync-feature](../README.md)

## Description

A real AC file that exists on disk.

## Inputs

| Name | Required | Description |
|---|---|---|

## Verification

```bash
exit 0
```

## Scenarios

(None yet.)
AC

# Orphaned file: exists on disk but NOT listed in the README
cat > "$FIXTURE_DIR/features/sync-feature/_acs/orphaned-ac.ac.md" << 'AC'
# AC: orphaned-ac

**Status:** planned
**Feature:** [sync-feature](../README.md)

## Description

This AC file exists on disk but is not listed in the _acs/README.md index.

## Inputs

| Name | Required | Description |
|---|---|---|

## Verification

```bash
exit 0
```

## Scenarios

(None yet.)
AC

# README lists real-ac + phantom-ac (no file), but omits orphaned-ac
cat > "$FIXTURE_DIR/features/sync-feature/_acs/README.md" << 'README'
# Acceptance Criteria: sync-feature

| AC | Description | Status |
|---|---|---|
| [real-ac](real-ac.ac.md) | A real AC | implemented |
| [phantom-ac](phantom-ac.ac.md) | Listed but no file on disk | planned |
README

cat > "$FIXTURE_DIR/features/sync-feature/README.md" << 'README'
# Feature: sync-feature

**Status:** In Progress

## Summary

Fixture feature for AC index sync testing.

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [real-ac](_acs/real-ac.ac.md) | A real AC | implemented |
README

# ── Scenario with unresolvable AC reference ────────────────────────
cat > "$FIXTURE_DIR/bad-ac-ref.test.md" << SCENARIO
# Scenario: Bad AC reference

**Description:** References an AC that does not exist on disk.
**Tags:** fixture

## step-with-bad-ref

**ACs:**

| Feature | ACs |
|---|---|
| [fixture-feature]($FIXTURE_DIR/features/fixture-feature/) | [nonexistent-ac]($FIXTURE_DIR/features/fixture-feature/_acs/nonexistent-ac.ac.md) |

\`\`\`bash
echo "this step references a missing AC"
\`\`\`
SCENARIO

# ── Scenario whose step writes a marker file (for no-execution test) ──
MARKER_FILE="$FIXTURE_DIR/execution.marker"
cat > "$FIXTURE_DIR/marker-scenario.test.md" << SCENARIO
# Scenario: Marker writer

**Description:** If executed, this scenario writes a marker file to disk.
**Tags:** fixture

## write-marker

\`\`\`bash
touch "$MARKER_FILE"
\`\`\`
SCENARIO

# Propagate vars to context
echo "BINARY_PATH=$BINARY_PATH"
echo "SPEC_ROOT=$SPEC_ROOT"
echo "FIXTURE_DIR=$FIXTURE_DIR"
echo "MARKER_FILE=$MARKER_FILE"
````

## validate-valid-scenario

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [validates-scenario-structure]($SPEC_ROOT/features/cli/validate/_acs/validates-scenario-structure.ac.md) |

````bash
output=$("$BINARY_PATH" validate "$FIXTURE_DIR/valid-scenario.test.md" 2>&1)
rc=$?
test $rc -eq 0 || {
  echo "Expected exit 0 for valid scenario, got $rc"
  echo "$output"
  exit 1
}
echo "Valid scenario passed validation (exit $rc)"
````

## validate-malformed-scenario

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [validates-scenario-structure]($SPEC_ROOT/features/cli/validate/_acs/validates-scenario-structure.ac.md) |

````bash
output=$("$BINARY_PATH" validate "$FIXTURE_DIR/malformed-scenario.test.md" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for malformed scenario, got $rc"
  echo "$output"
  exit 1
}

# Error must mention the file path
echo "$output" | grep -q "malformed-scenario.test.md" || {
  echo "Error output does not mention file path"
  echo "$output"
  exit 1
}

# Error must include a line number
echo "$output" | grep -qE 'line [0-9]+' || {
  echo "Error output does not include a line number"
  echo "$output"
  exit 1
}

echo "Malformed scenario correctly rejected (exit $rc)"
echo "$output"
````

## validate-valid-ac

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [validates-ac-structure]($SPEC_ROOT/features/cli/validate/_acs/validates-ac-structure.ac.md) |

````bash
output=$("$BINARY_PATH" validate "$FIXTURE_DIR/features/fixture-feature/_acs/always-pass.ac.md" 2>&1)
rc=$?
test $rc -eq 0 || {
  echo "Expected exit 0 for valid AC, got $rc"
  echo "$output"
  exit 1
}
echo "Valid AC passed validation (exit $rc)"
````

## validate-malformed-ac

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [validates-ac-structure]($SPEC_ROOT/features/cli/validate/_acs/validates-ac-structure.ac.md) |

````bash
output=$("$BINARY_PATH" validate "$FIXTURE_DIR/features/fixture-feature/_acs/bad-ac.ac.md" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for malformed AC, got $rc"
  echo "$output"
  exit 1
}

# Error must reference the file path
echo "$output" | grep -q "bad-ac.ac.md" || {
  echo "Error output does not mention file path"
  echo "$output"
  exit 1
}

# Error must describe the structural problem (missing Status)
echo "$output" | grep -qiE '(status|missing)' || {
  echo "Error output does not describe the structural issue"
  echo "$output"
  exit 1
}

echo "Malformed AC correctly rejected (exit $rc)"
echo "$output"
````

## validate-ac-refs

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [validates-ac-refs-resolve]($SPEC_ROOT/features/cli/validate/_acs/validates-ac-refs-resolve.ac.md) |

````bash
output=$("$BINARY_PATH" validate "$FIXTURE_DIR/bad-ac-ref.test.md" --spec-root "$FIXTURE_DIR" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for unresolvable AC reference, got $rc"
  echo "$output"
  exit 1
}

# Error must mention the missing AC
echo "$output" | grep -qiE '(unresolved|not found|missing|does not exist|nonexistent)' || {
  echo "Error output does not mention the unresolved reference"
  echo "$output"
  exit 1
}

# Error must name the scenario file
echo "$output" | grep -q "bad-ac-ref.test.md" || {
  echo "Error output does not mention the scenario file path"
  echo "$output"
  exit 1
}

echo "Unresolvable AC reference correctly rejected (exit $rc)"
echo "$output"
````

## validate-ac-index-sync

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [validates-ac-index-sync]($SPEC_ROOT/features/cli/validate/_acs/validates-ac-index-sync.ac.md) |

````bash
# The sync-feature _acs/ directory has:
#   - phantom entry: phantom-ac listed in README but no file on disk
#   - orphaned file: orphaned-ac.ac.md on disk but not in README
output=$("$BINARY_PATH" validate --spec-root "$FIXTURE_DIR" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for AC index sync issue, got $rc"
  echo "$output"
  exit 1
}

# Error must mention the sync problem
echo "$output" | grep -qiE '(orphan|phantom|not listed|missing file|not in readme|sync)' || {
  echo "Error output does not describe the sync issue"
  echo "$output"
  exit 1
}

echo "AC index sync issue correctly detected (exit $rc)"
echo "$output"
````

## validate-no-execution

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [no-execution]($SPEC_ROOT/features/cli/validate/_acs/no-execution.ac.md) |

````bash
# Ensure marker does not already exist
rm -f "$MARKER_FILE"

# Run validate (not run) on the marker scenario
"$BINARY_PATH" validate "$FIXTURE_DIR/marker-scenario.test.md" 2>&1 || true

# Marker file must NOT exist — validate must not execute code blocks
if [ -f "$MARKER_FILE" ]; then
  echo "FAIL: marker file '$MARKER_FILE' exists — validate executed step code"
  rm -f "$MARKER_FILE"
  exit 1
fi

echo "No-execution guarantee verified: marker file does not exist"
````

## validate-all-valid

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [exit-0-all-valid]($SPEC_ROOT/features/cli/validate/_acs/exit-0-all-valid.ac.md) |

````bash
# Self-referential dogfooding: validate the REAL spec tree (not the fixture dir).
# This spec tree contains the very test file that defines this step.
output=$("$BINARY_PATH" validate "$SPEC_ROOT" 2>&1)
rc=$?
test $rc -eq 0 || {
  echo "Expected exit 0 for real spec root, got $rc"
  echo "$output"
  exit 1
}

# Output should contain a summary indicating success
echo "$output" | grep -qiE '(no errors|valid|pass)' || {
  echo "Expected summary line indicating success"
  echo "$output"
  exit 1
}

echo "Real spec tree validated successfully (exit $rc)"
echo "$output"
````

## validate-exit-codes

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [exit-2-invalid-args]($SPEC_ROOT/features/cli/validate/_acs/exit-2-invalid-args.ac.md) |

````bash
# Pass an unknown flag to trigger an argument error
output=$("$BINARY_PATH" validate --no-such-flag 2>&1)
rc=$?
test $rc -eq 2 || {
  echo "Expected exit 2 for invalid arguments, got $rc"
  echo "$output"
  exit 1
}

echo "Invalid arguments correctly returned exit 2"
echo "$output"
````

## validate-fail-fast

**ACs:**

| Feature | ACs |
|---|---|
| [cli/validate]($SPEC_ROOT/features/cli/validate/) | [fail-fast-stops-early]($SPEC_ROOT/features/cli/validate/_acs/fail-fast-stops-early.ac.md) |

````bash
# The fixture dir contains multiple bad files (malformed-scenario.test.md, bad-ac.ac.md, etc.)
# Without --fail-fast, all errors are collected
all_output=$("$BINARY_PATH" validate "$FIXTURE_DIR" --spec-root "$FIXTURE_DIR" 2>&1)
all_rc=$?
test $all_rc -eq 1 || {
  echo "Expected exit 1 without --fail-fast, got $all_rc"
  echo "$all_output"
  exit 1
}

# Count errors without --fail-fast
all_error_lines=$(echo "$all_output" | grep -cE '^\s+(line [0-9]+:|[a-z])' || true)
echo "Errors without --fail-fast: $all_error_lines"

# With --fail-fast=1, output should be truncated
ff_output=$("$BINARY_PATH" validate "$FIXTURE_DIR" --spec-root "$FIXTURE_DIR" --fail-fast 2>&1)
ff_rc=$?
test $ff_rc -eq 1 || {
  echo "Expected exit 1 with --fail-fast, got $ff_rc"
  echo "$ff_output"
  exit 1
}

# Must contain truncation note
echo "$ff_output" | grep -q "fail-fast" || {
  echo "Expected truncation note in output"
  echo "$ff_output"
  exit 1
}

# Fewer error lines than full run
ff_error_lines=$(echo "$ff_output" | grep -cE '^\s+(line [0-9]+:|[a-z])' || true)
echo "Errors with --fail-fast: $ff_error_lines"

test "$ff_error_lines" -lt "$all_error_lines" || {
  echo "Expected fewer errors with --fail-fast ($ff_error_lines) than without ($all_error_lines)"
  exit 1
}

echo "Fail-fast correctly truncated output"
echo "$ff_output"
````

## Teardown

````bash
rm -rf "$FIXTURE_DIR"
````
