# Features: Rehearse

Feature specifications for the Rehearse test framework, managed by [Synchestra](https://github.com/synchestra-io/synchestra).

## Index

| Feature | Status | Description |
|---|---|---|
| [testing-framework](testing-framework/README.md) | Conceptual | Markdown-native testing framework: scenario format, Go-based runner, and CLI integration for composing ACs into E2E and integration test flows |
| [acceptance-criteria](acceptance-criteria/README.md) | Conceptual | First-class, individually addressable verification artifacts with status lifecycle, typed inputs, and executable verification scripts |
| [cli](cli/README.md) | In Progress | The `rehearse` CLI — primary interface for running, listing, and validating test scenarios |

## Feature Summaries

### [Testing Framework](testing-framework/README.md)

Turns specifications into executable verification — without leaving markdown. Composes acceptance criteria into multi-step test workflows that read as documentation and execute as test suites. Contains two sub-features: [test-scenario](testing-framework/test-scenario/README.md) defines the human-readable scenario format (named steps, data passing between steps, AC references, sub-flow includes, parallel groups), and [test-runner](testing-framework/test-runner/README.md) is the Go execution engine that parses scenarios, resolves AC verification scripts from `_acs/` directories, and produces structured pass/fail reports. The framework dogfoods itself — the runner's own test scenarios are executed by the runner it verifies.

### [Acceptance Criteria](acceptance-criteria/README.md)

The contract between what a feature promises and what the system delivers. Each AC is a standalone markdown file — readable by product owners, auditable by reviewers, executable by the test runner. Defines the AC file format, supported verification languages (Bash, Python, SQL, Starlark), status lifecycle, identification scheme, and validation rules.

### [CLI](cli/README.md)

The `rehearse` command-line interface. Follows a `rehearse <action>` pattern with consistent exit codes and structured output formats. Contains four sub-commands: [run](cli/run/README.md) executes test scenarios, [list](cli/list/README.md) discovers and lists available scenarios, [validate](cli/validate/README.md) checks structural validity of scenarios and ACs without execution, and [version](cli/version/README.md) prints build information.

## Outstanding Questions

None at this time.
