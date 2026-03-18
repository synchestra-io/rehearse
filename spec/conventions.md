# Spec Conventions

Rules and naming conventions for files in the `spec/` directory tree.

## File naming

| File type | Extension | Location | Example |
|---|---|---|---|
| Feature README | `README.md` | `spec/features/{feature}/` | `spec/features/cli/project/new/README.md` |
| Acceptance criterion | `{ac-slug}.ac.md` | `spec/features/{feature}/_acs/` | `spec/features/cli/project/new/_acs/creates-spec-config.ac.md` |
| AC index | `README.md` | `spec/features/{feature}/_acs/` | `spec/features/cli/project/new/_acs/README.md` |
| Test scenario | `{scenario-slug}.test.md` | `spec/features/{feature}/_tests/` | `spec/features/testing-framework/test-runner/_tests/runner-core.test.md` |
| Test index | `README.md` | `spec/features/{feature}/_tests/` | `spec/features/testing-framework/test-runner/_tests/README.md` |

### Slugs

All slugs (AC and scenario) are lowercase, hyphen-separated, and unique within their parent directory.

### Reserved directories

| Directory | Purpose |
|---|---|
| `_acs/` | Acceptance criteria for the parent feature |
| `_tests/` | Test scenarios for the parent feature |

The `_` prefix signals "not a sub-feature" — these directories are excluded from the feature index and Contents table.

## Outstanding Questions

None at this time.
