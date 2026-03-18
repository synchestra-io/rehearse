# Scenario: CLI run core behaviors

**Description:** Integration test of the `rehearse run` command — single-file execution, directory scanning, tag filtering, manual-scenario handling, JSON output, and exit codes.
**Tags:** integration, cli, run

## Setup

````bash
BINARY_PATH="${BINARY_PATH:-$(go env GOPATH)/bin/rehearse}"
SPEC_ROOT="$(git rev-parse --show-toplevel)/spec"
FIXTURE_DIR=$(mktemp -d)

# Create a passing scenario (tag: e2e)
cat > "$FIXTURE_DIR/pass.test.md" << 'SCENARIO'
# Scenario: Passing test

**Description:** A minimal scenario that passes.
**Tags:** e2e

## pass-step

```bash
echo "pass"
exit 0
```
SCENARIO

# Create a failing scenario (tag: smoke)
cat > "$FIXTURE_DIR/fail.test.md" << 'SCENARIO'
# Scenario: Failing test

**Description:** A minimal scenario that fails.
**Tags:** smoke

## fail-step

```bash
echo "fail"
exit 1
```
SCENARIO

# Create a manual scenario (tag: manual)
cat > "$FIXTURE_DIR/manual-demo.test.md" << 'SCENARIO'
# Scenario: Manual demo

**Description:** A scenario tagged as manual.
**Tags:** manual

## manual-step

```bash
echo "manual ran"
```
SCENARIO

# Propagate vars to context
echo "BINARY_PATH=$BINARY_PATH"
echo "SPEC_ROOT=$SPEC_ROOT"
echo "FIXTURE_DIR=$FIXTURE_DIR"
````

## build-binary

```bash
cd "$(git rev-parse --show-toplevel)"
go build -o "$BINARY_PATH" .
```

## run-single-file

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [runs-single-file]($SPEC_ROOT/features/cli/run/_acs/runs-single-file.ac.md), [exit-0-on-pass]($SPEC_ROOT/features/cli/run/_acs/exit-0-on-pass.ac.md) |

````bash
"$BINARY_PATH" run "$FIXTURE_DIR/pass.test.md" --format json
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; exit 1; }
echo "Single-file run passed with exit $rc"
````

## run-directory

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [runs-directory]($SPEC_ROOT/features/cli/run/_acs/runs-directory.ac.md) |

````bash
output=$("$BINARY_PATH" run "$FIXTURE_DIR" --format json 2>&1) || true
# Verify multiple scenarios were discovered (pass + fail at minimum)
echo "$output" | grep -q "pass.test.md\|Passing test" || { echo "Expected passing scenario in output"; exit 1; }
echo "$output" | grep -q "fail.test.md\|Failing test" || { echo "Expected failing scenario in output"; exit 1; }
echo "Directory scan found multiple scenarios"
````

## filter-by-tag

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [filters-by-tag]($SPEC_ROOT/features/cli/run/_acs/filters-by-tag.ac.md) |

````bash
output=$("$BINARY_PATH" run "$FIXTURE_DIR" --tag e2e --format json 2>&1)
rc=$?
# Only the e2e-tagged scenario should run (and pass)
test $rc -eq 0 || { echo "Expected exit 0 for e2e-only run, got $rc"; exit 1; }
# The smoke-tagged failing scenario should not appear
echo "$output" | grep -qi "Failing test\|fail.test.md" && { echo "Failing scenario should have been filtered out"; exit 1; }
echo "Tag filtering correctly ran only e2e scenarios"
````

## json-output

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [json-output]($SPEC_ROOT/features/cli/run/_acs/json-output.ac.md) |

````bash
output=$("$BINARY_PATH" run "$FIXTURE_DIR/pass.test.md" --format json)
# Verify the output is valid JSON
echo "$output" | python3 -c "import sys, json; json.load(sys.stdin)" || { echo "Output is not valid JSON"; exit 1; }
# Verify expected fields are present
echo "$output" | python3 -c "
import sys, json
data = json.load(sys.stdin)
# Accept top-level object or array
obj = data[0] if isinstance(data, list) else data
assert 'scenario' in obj or 'name' in obj, 'Missing scenario/name field'
assert 'steps' in obj, 'Missing steps field'
assert 'status' in obj or 'result' in obj, 'Missing status/result field'
print('JSON structure valid')
"
````

## manual-skipped-in-scan

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [skips-manual-in-scan]($SPEC_ROOT/features/cli/run/_acs/skips-manual-in-scan.ac.md) |

````bash
output=$("$BINARY_PATH" run "$FIXTURE_DIR" --format json 2>&1) || true
# Manual scenario should NOT appear in directory-scan results
echo "$output" | grep -qi "Manual demo\|manual-demo.test.md\|manual ran" && {
  echo "Manual scenario should be skipped in directory scan"
  exit 1
}
echo "Manual scenario correctly skipped in directory scan"
````

## manual-direct-path

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [runs-manual-when-direct]($SPEC_ROOT/features/cli/run/_acs/runs-manual-when-direct.ac.md) |

````bash
output=$("$BINARY_PATH" run "$FIXTURE_DIR/manual-demo.test.md" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0 for direct manual run, got $rc"; exit 1; }
echo "$output" | grep -qi "Manual demo\|manual" || { echo "Expected manual scenario in output"; exit 1; }
echo "Manual scenario executed when given by direct path"
````

## manual-with-flag

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [runs-manual-with-flag]($SPEC_ROOT/features/cli/run/_acs/runs-manual-with-flag.ac.md) |

````bash
output=$("$BINARY_PATH" run "$FIXTURE_DIR" --run-manual-tests --format json 2>&1) || true
# Manual scenario SHOULD now appear in results
echo "$output" | grep -qi "Manual demo\|manual-demo.test.md\|manual ran" || {
  echo "Manual scenario should be included with --run-manual-tests flag"
  exit 1
}
echo "Manual scenario included with --run-manual-tests flag"
````

## exit-on-failure

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [exit-1-on-failure]($SPEC_ROOT/features/cli/run/_acs/exit-1-on-failure.ac.md) |

````bash
"$BINARY_PATH" run "$FIXTURE_DIR/fail.test.md" --format json; rc=$?
test $rc -eq 1 || { echo "Expected exit 1 for failing scenario, got $rc"; exit 1; }
echo "Failing scenario correctly exited with code 1"
````

## exit-not-found

**ACs:**

| Feature | ACs |
|---|---|
| [cli/run]($SPEC_ROOT/features/cli/run/) | [exit-3-not-found]($SPEC_ROOT/features/cli/run/_acs/exit-3-not-found.ac.md) |

````bash
"$BINARY_PATH" run "/nonexistent/path/does-not-exist.test.md" --format json 2>&1; rc=$?
test $rc -eq 3 || { echo "Expected exit 3 for not-found path, got $rc"; exit 1; }
echo "Not-found path correctly exited with code 3"
````

## Teardown

```bash
rm -rf "$FIXTURE_DIR"
```
