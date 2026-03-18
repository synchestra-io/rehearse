# Acceptance Criteria: cli/run

Acceptance criteria for [`rehearse run`](../README.md).

| AC | Description | Status |
|---|---|---|
| [runs-single-file](runs-single-file.ac.md) | Single scenario file executes and reports results | planned |
| [runs-directory](runs-directory.ac.md) | Directory scan discovers and runs all `.test.md` files | planned |
| [filters-by-tag](filters-by-tag.ac.md) | `--tag` flag filters scenarios to matching tags only | planned |
| [json-output](json-output.ac.md) | `--format json` produces valid JSON with required fields | planned |
| [skips-manual-in-scan](skips-manual-in-scan.ac.md) | Manual-tagged scenarios skipped during directory scans | planned |
| [runs-manual-when-direct](runs-manual-when-direct.ac.md) | Manual-tagged scenario runs when file path given directly | planned |
| [runs-manual-with-flag](runs-manual-with-flag.ac.md) | `--run-manual-tests` includes manual scenarios in scans | planned |
| [overrides-spec-root](overrides-spec-root.ac.md) | `--spec-root` changes AC resolution root directory | planned |
| [exit-0-on-pass](exit-0-on-pass.ac.md) | Exits 0 when all scenarios pass | planned |
| [exit-1-on-failure](exit-1-on-failure.ac.md) | Exits 1 when any scenario fails | planned |
| [exit-2-invalid-args](exit-2-invalid-args.ac.md) | Exits 2 on invalid arguments | planned |
| [exit-3-not-found](exit-3-not-found.ac.md) | Exits 3 when target path does not exist | planned |

## Outstanding Questions

None at this time.
