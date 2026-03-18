# Scenario: Version core behaviors

**Description:** Verifies that `rehearse version` prints version info in the expected format and exits with code 0.
**Tags:** integration, cli, version

## Setup

````bash
BINARY_PATH="${BINARY_PATH:-$(go env GOPATH)/bin/rehearse}"
SPEC_ROOT="$(git rev-parse --show-toplevel)/spec"

echo "BINARY_PATH=$BINARY_PATH"
echo "SPEC_ROOT=$SPEC_ROOT"
````

## build-binary

```bash
cd "$(git rev-parse --show-toplevel)"
go build -o "$BINARY_PATH" .
```

## prints-version

**Outputs:**

| Name | Store | Extract |
|---|---|---|
| version_output | context | `cat $STEP_STDOUT` |

**ACs:**

| Feature | ACs |
|---|---|
| [cli/version]($SPEC_ROOT/features/cli/version/) | [prints-version-info]($SPEC_ROOT/features/cli/version/_acs/prints-version-info.ac.md) |

````bash
output=$("$BINARY_PATH" version 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Verify output matches expected format: rehearse {version} ({commit}) {date}
pattern='^rehearse .+ \(.+\) .+$'
echo "$output" | grep -Eq "$pattern" || {
  echo "Output does not match expected format"
  echo "Expected pattern: rehearse {version} ({commit}) {date}"
  echo "Got: $output"
  exit 1
}

echo "$output"
````
