# Acceptance Criteria: cli/validate

Acceptance criteria for [`rehearse validate`](../README.md).

| AC | Description | Status |
|---|---|---|
| [validates-scenario-structure](validates-scenario-structure.ac.md) | Well-formed scenarios pass; malformed scenarios rejected with line-number errors | planned |
| [validates-ac-structure](validates-ac-structure.ac.md) | Well-formed ACs pass; malformed ACs rejected with specific errors | planned |
| [validates-ac-refs-resolve](validates-ac-refs-resolve.ac.md) | AC references in scenarios must resolve to actual files on disk | planned |
| [validates-ac-index-sync](validates-ac-index-sync.ac.md) | _acs/README.md entries and .ac.md files must be in sync | planned |
| [no-execution](no-execution.ac.md) | Validation does not execute any scripts or code blocks | planned |
| [exit-0-all-valid](exit-0-all-valid.ac.md) | Exits 0 when all checks pass | planned |
| [exit-1-validation-errors](exit-1-validation-errors.ac.md) | Exits 1 when validation errors are found | planned |
| [exit-2-invalid-args](exit-2-invalid-args.ac.md) | Exits 2 on invalid arguments | planned |

## Outstanding Questions

None at this time.
