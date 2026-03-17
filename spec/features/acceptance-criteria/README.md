# Feature: Acceptance Criteria

**Status:** Conceptual

## Summary

Acceptance criteria are the contract between what a feature promises and what the system actually delivers. Each AC is a standalone markdown file — readable by product owners, auditable by reviewers, and executable by the [test runner](../testing-framework/test-runner/README.md). ACs live alongside the features they verify, carry their own lifecycle, and compose into [test scenarios](../testing-framework/test-scenario/README.md) for end-to-end validation. Write an AC once; reference it from any number of test flows.

The full specification for this feature lives in the [synchestra-io/synchestra](https://github.com/synchestra-io/synchestra/blob/main/spec/features/acceptance-criteria/) repository, where it is defined as a core Synchestra concept. Rehearse implements the execution side — parsing AC files and running their verification scripts.

## Key Concepts

- **AC file location:** `spec/features/{feature}/_acs/{ac-slug}.md`
- **AC identification:** Path-based — e.g., `cli/project/new/creates-spec-config`
- **AC statuses:** `planned` → `wip` → `implemented` → `deprecated`
- **Supported languages:** Bash, Python, SQL, Starlark
- **Mandatory language annotation:** Code blocks without annotations are validation errors

## Full Specification

For the complete behavior, file format, validation rules, and interaction details, see:

- [Acceptance Criteria](https://github.com/synchestra-io/synchestra/blob/main/spec/features/acceptance-criteria/) — full feature specification in synchestra-io/synchestra

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Acceptance criteria not yet defined for this feature.
