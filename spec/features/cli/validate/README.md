# Feature: CLI — validate

**Status:** In Progress

## Summary

`rehearse validate` checks the structural validity of test scenarios and acceptance criteria files without executing any code. It is the static analysis counterpart to `rehearse run` — catching format errors, missing references, and orphaned files before runtime.

## Problem

Test scenarios and AC files follow strict markdown conventions. Format errors — missing language annotations, malformed sections, unresolvable references — are only caught at runtime today, when `rehearse run` fails with a parse error. By then the feedback loop is long: you author a scenario, push to CI, wait for the run, and only then discover a typo in an AC reference or a bare code fence without a language annotation.

Validation should happen earlier: at authoring time, in CI as a pre-check, or as a pre-commit hook. A dedicated `validate` command enables all three without the overhead of executing scripts, standing up test fixtures, or waiting for real step execution. It answers one question: "Are these files structurally correct?" — and answers it in seconds.

## Behavior

### Usage

```
rehearse validate [path]              — validate scenario/AC files at path
rehearse validate                     — validate all scenarios and ACs under spec root
rehearse validate --spec-root ./spec  — override spec root directory
rehearse validate --fail-fast         — stop after the first error
rehearse validate --fail-fast=5       — stop after 5 errors
```

When `path` is a file, validate checks that single file. When `path` is a directory, validate recursively discovers and checks all `.test.md` and `.ac.md` files under it. When no path is given, validate scans the entire spec root.

### Flags

| Flag | Default | Description |
|---|---|---|
| `--spec-root` | `spec` | Override the spec root directory |
| `--fail-fast[=N]` | `0` (disabled) | Stop after N errors. When used without a value, defaults to 1. When set to 0 or omitted, all errors are collected. |

### Checks performed

The validate command performs four categories of structural checks. All checks are static — no scripts are executed, no interpreters are invoked.

#### 1. Scenario structure validation

Validates that each `.test.md` file conforms to the [test scenario](../../testing-framework/test-scenario/README.md) format:

- **Title present.** The file must start with `# Scenario: {name}`. Missing or malformed titles are rejected.
- **Description metadata present.** The `**Description:**` field must appear after the title. Scenarios without a description are rejected.
- **All steps have a code block or Include directive.** A step heading (`## {name}`) with no code block and no `**Include:**` directive is a validation error — every step must do something.
- **All code blocks have a language annotation.** Bare code fences (`` ``` `` without a language) are rejected. Valid annotations are `bash`, `python`, `sql`, `starlark`, and `http`.
- **Step names are kebab-case and unique within the scenario.** Duplicate step names or names with spaces/underscores/capitals are rejected.
- **`Depends on` references point to steps defined earlier in the file.** No forward references, no references to undefined steps, no cycles. Dependency chains are validated as a DAG.
- **Reserved step names (`Setup`, `Teardown`) do not use Outputs, ACs, Parallel, or Depends on.** These are lifecycle hooks, not test steps — metadata beyond a code block is a validation error.
- **No circular includes.** When a step uses `**Include:**`, the referenced file is checked recursively. If the include chain forms a cycle, it is rejected.

#### 2. AC file structure validation

Validates that each `.ac.md` file conforms to the [acceptance criteria](../../acceptance-criteria/README.md) format:

- **Title present.** The file must start with `# AC: {slug}`.
- **Status field present and valid.** The `**Status:**` field must be one of `planned`, `wip`, `implemented`, or `deprecated`.
- **Feature back-reference present.** The `**Feature:**` field must link back to the parent feature.
- **Description section present.** The `## Description` section must exist and be non-empty.
- **Inputs table present.** The `## Inputs` section must exist. The table may be empty (no required inputs), but the section itself is mandatory.
- **Verification section present with a code block.** For `wip` and `implemented` statuses, the `## Verification` section must contain a code block. For `planned` status, the section may exist without a code block. For `deprecated` status, the verification section is optional.
- **Verification code block has a language annotation.** Same rule as scenario code blocks — bare fences are rejected.
- **Slug matches filename.** The slug in `# AC: {slug}` must match the filename stem. For example, `creates-spec-config` must live in `creates-spec-config.ac.md`.

#### 3. Cross-reference validation

Validates that all references between scenarios and AC files resolve to real files on disk:

- **Every AC reference in a scenario's ACs table resolves to an actual `.ac.md` file.** A scenario step that references `cli/project/new/creates-spec-config` must have a corresponding file at `{spec_root}/features/cli/project/new/_acs/creates-spec-config.ac.md`.
- **Wildcard references (`*`) resolve to a non-empty `_acs/` directory.** A wildcard reference to a feature's ACs is valid only if the `_acs/` directory exists and contains at least one `.ac.md` file.
- **Include references resolve to actual `.md` files.** A step with `**Include:** [name](path)` is valid only if the referenced file exists on disk.

#### 4. AC index synchronization

Validates that `_acs/` directories and their `README.md` index files are in sync:

- **Every `.ac.md` file in an `_acs/` directory is listed in that directory's `README.md`.** An AC file that exists on disk but is missing from the index is an orphaned file — a validation error.
- **Every entry in an `_acs/README.md` table has a corresponding `.ac.md` file.** An index entry that references a non-existent file is a phantom entry — a validation error.

### What validate does NOT check

- **Does not execute verification scripts or test steps.** Validate is purely structural — it reads markdown, not bash.
- **Does not check feature READMEs for AC sections.** Enforcing `## Acceptance Criteria` in feature READMEs is Synchestra's responsibility, not the validate command's.
- **Does not validate the semantic correctness of scripts.** Whether a bash script would succeed, a SQL query would parse, or a Python expression would evaluate is out of scope. Validate checks that the code block exists and has a language annotation — not that its contents are correct.

### Fail-fast mode

By default, validate collects all errors across all files before reporting. This gives a complete picture but may be slow on large trees with many structural problems.

When `--fail-fast` is used, validate stops as soon as the error limit is reached:

- `--fail-fast` — stop after the first error (equivalent to `--fail-fast=1`).
- `--fail-fast=N` — stop after N errors.
- `--fail-fast=0` — same as omitting the flag; collect all errors.

When validation is truncated, the output includes a note: `(output truncated due to --fail-fast)`. The exit code is still `1` — truncation does not change the exit code contract.

Fail-fast applies at the file boundary: errors within a single file are always fully collected, but no further files are processed once the limit is reached.

### Output

Validate lists all errors with file path, line number (where possible), and error description. Errors are grouped by file for readability:

```
spec/features/cli/project/new/_acs/creates-spec-config.ac.md
  line 1: missing AC title (expected "# AC: {slug}")
  line 5: status field missing or invalid

spec/tests/project-lifecycle.test.md
  line 42: code block missing language annotation
  line 58: step "verify-configs" references undefined AC "cli/project/new/missing-ac"

2 files, 4 errors
```

When all checks pass, validate prints a summary and exits silently:

```
Validated 12 scenarios, 34 ACs — no errors.
```

### Exit code contract

| Exit code | Meaning |
|---|---|
| `0` | All checks passed |
| `1` | Validation errors found |
| `2` | Invalid arguments |
| `3` | Resource not found (path does not exist) |
| `10+` | Unexpected errors |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Testing Framework](../../testing-framework/README.md) | Validate checks files against the formats defined by the testing framework. |
| [Test Scenario](../../testing-framework/test-scenario/README.md) | Validate enforces the scenario structure spec — title, metadata, step format, code block annotations. |
| [Acceptance Criteria](../../acceptance-criteria/README.md) | Validate enforces the AC file format and cross-reference integrity between scenarios and `_acs/` directories. |
| [Test Runner](../../testing-framework/test-runner/README.md) | Validate shares parsing logic with the runner. The runner parses to execute; validate parses to check. |

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [validates-scenario-structure](_acs/validates-scenario-structure.ac.md) | Well-formed scenarios pass; malformed scenarios rejected with line-number errors | planned |
| [validates-ac-structure](_acs/validates-ac-structure.ac.md) | Well-formed ACs pass; malformed ACs rejected with specific errors | planned |
| [validates-ac-refs-resolve](_acs/validates-ac-refs-resolve.ac.md) | AC references in scenarios must resolve to actual files on disk | planned |
| [validates-ac-index-sync](_acs/validates-ac-index-sync.ac.md) | _acs/README.md entries and .ac.md files must be in sync | planned |
| [no-execution](_acs/no-execution.ac.md) | Validation does not execute any scripts or code blocks | planned |
| [exit-0-all-valid](_acs/exit-0-all-valid.ac.md) | Exits 0 when all checks pass | planned |
| [exit-1-validation-errors](_acs/exit-1-validation-errors.ac.md) | Exits 1 when validation errors are found | planned |
| [exit-2-invalid-args](_acs/exit-2-invalid-args.ac.md) | Exits 2 on invalid arguments | planned |
| [fail-fast-stops-early](_acs/fail-fast-stops-early.ac.md) | --fail-fast stops validation after the error limit is reached | planned |

## Outstanding Questions

None at this time.
