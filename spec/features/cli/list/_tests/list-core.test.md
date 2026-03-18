# Scenario: List core behaviors

**Description:** Integration test of the `rehearse list` command — scenario discovery across `spec/tests/` and `spec/features/*/_tests/` directories, columnar output format, and tag-based filtering via `--tag`.
**Tags:** integration, cli, list

## Setup

````bash
BINARY_PATH="${BINARY_PATH:-$(go env GOPATH)/bin/rehearse}"
SPEC_ROOT="$(git rev-parse --show-toplevel)/spec"
FIXTURE_DIR=$(mktemp -d)

# Build the rehearse binary
cd "$(git rev-parse --show-toplevel)"
go build -o "$BINARY_PATH" .

# Create fixture spec root with two discovery paths
mkdir -p "$FIXTURE_DIR/spec/tests"
mkdir -p "$FIXTURE_DIR/spec/features/example/_tests"

# Cross-feature scenario (tags: e2e)
cat > "$FIXTURE_DIR/spec/tests/e2e-flow.test.md" << 'SCENARIO'
# Scenario: E2E flow

**Description:** End-to-end integration flow.
**Tags:** e2e

## placeholder

```bash
echo "e2e"
```
SCENARIO

# Feature-scoped scenario (tags: unit)
cat > "$FIXTURE_DIR/spec/features/example/_tests/unit.test.md" << 'SCENARIO'
# Scenario: Unit checks

**Description:** Basic unit-level checks.
**Tags:** unit

## placeholder

```bash
echo "unit"
```
SCENARIO

# Propagate vars to context
echo "BINARY_PATH=$BINARY_PATH"
echo "SPEC_ROOT=$SPEC_ROOT"
echo "FIXTURE_DIR=$FIXTURE_DIR"
````

## list-all

**Outputs:**

| Name | Store | Extract |
|---|---|---|
| binary_path | context | `echo $BINARY_PATH` |
| spec_root | context | `echo $FIXTURE_DIR` |

**ACs:**

| Feature | ACs |
|---|---|
| [cli/list]($SPEC_ROOT/features/cli/list/) | [lists-all-scenarios]($SPEC_ROOT/features/cli/list/_acs/lists-all-scenarios.ac.md) |

````bash
# Run rehearse list against the fixture spec root
output=$("$BINARY_PATH" list --spec-root "$FIXTURE_DIR" 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Verify columnar header
echo "$output" | head -1 | grep -q "SCENARIO" || { echo "Missing SCENARIO column header"; exit 1; }
echo "$output" | head -1 | grep -q "DESCRIPTION" || { echo "Missing DESCRIPTION column header"; exit 1; }
echo "$output" | head -1 | grep -q "TAGS" || { echo "Missing TAGS column header"; exit 1; }

# Verify both scenarios appear
echo "$output" | grep -q "E2E flow" || { echo "Scenario 'E2E flow' not found in output"; exit 1; }
echo "$output" | grep -q "Unit checks" || { echo "Scenario 'Unit checks' not found in output"; exit 1; }

# Verify at least header + 2 data rows
row_count=$(echo "$output" | wc -l | tr -d ' ')
test "$row_count" -ge 3 || { echo "Expected at least 3 lines (header + 2 data), got $row_count"; exit 1; }

echo "PASS: all scenarios listed"
echo "$output"
````

## filter-by-tag

**Outputs:**

| Name | Store | Extract |
|---|---|---|
| binary_path | context | `echo $BINARY_PATH` |
| spec_root | context | `echo $FIXTURE_DIR` |
| filter_tag | context | `echo e2e` |

**ACs:**

| Feature | ACs |
|---|---|
| [cli/list]($SPEC_ROOT/features/cli/list/) | [filters-by-tag]($SPEC_ROOT/features/cli/list/_acs/filters-by-tag.ac.md) |

````bash
# Run rehearse list with tag filter
filtered=$("$BINARY_PATH" list --spec-root "$FIXTURE_DIR" --tag e2e 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$filtered"; exit 1; }

# Verify the e2e scenario appears
echo "$filtered" | grep -q "E2E flow" || { echo "Scenario 'E2E flow' not found in filtered output"; exit 1; }

# Verify the unit scenario is absent
echo "$filtered" | grep -q "Unit checks" && { echo "Scenario 'Unit checks' should not appear with --tag e2e"; exit 1; }

echo "PASS: --tag e2e returns only matching scenarios"
echo "$filtered"
````

## Teardown

```bash
rm -rf "$FIXTURE_DIR"
```
