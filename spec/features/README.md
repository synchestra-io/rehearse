# Features: Rehearse

Feature specifications for the Rehearse test framework, managed by [Synchestra](https://github.com/synchestra-io/synchestra).

## Index

| Feature | Status | Description |
|---|---|---|
| [testing-framework](testing-framework/README.md) | Conceptual | Markdown-native testing framework: scenario format, Go-based runner, and CLI integration for composing ACs into E2E and integration test flows |
| [acceptance-criteria](acceptance-criteria/README.md) | Conceptual | First-class, individually addressable verification artifacts — full specification in [synchestra-io/synchestra](https://github.com/synchestra-io/synchestra/blob/main/spec/features/acceptance-criteria/) |
| [cli](cli/README.md) | In Progress | The `rehearse` CLI — primary interface for running and listing test scenarios |

## Feature Summaries

### [Testing Framework](testing-framework/README.md)

Turns specifications into executable verification — without leaving markdown. Composes acceptance criteria into multi-step test workflows that read as documentation and execute as test suites. Contains two sub-features: [test-scenario](testing-framework/test-scenario/README.md) defines the human-readable scenario format (named steps, data passing between steps, AC references, sub-flow includes, parallel groups), and [test-runner](testing-framework/test-runner/README.md) is the Go execution engine that parses scenarios, resolves AC verification scripts from `_acs/` directories, and produces structured pass/fail reports. The framework dogfoods itself — the runner's own test scenarios are executed by the runner it verifies.

### [Acceptance Criteria](acceptance-criteria/README.md)

The contract between what a feature promises and what the system delivers. Each AC is a standalone markdown file — readable by product owners, auditable by reviewers, executable by the test runner. This is a proxy feature — the full specification lives in the [synchestra-io/synchestra](https://github.com/synchestra-io/synchestra/blob/main/spec/features/acceptance-criteria/) repository.

### [CLI](cli/README.md)

The `rehearse` command-line interface. Follows a `rehearse <action>` pattern for running and listing test scenarios with consistent exit codes and structured output formats.

## Outstanding Questions

None at this time.
