# AC: validates-ac-index-sync

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

Given an `_acs/` directory where every `.ac.md` file is listed in `_acs/README.md`
and every README entry has a corresponding file, validation passes. Given a phantom
entry (listed in the README but no file on disk) or an orphaned file (file exists but
not listed in the README), validation fails with an error naming the specific sync
issue.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| spec_root | Yes | Path to a spec root containing an `_acs/` directory with a sync issue |

## Verification

```bash
output=$("$binary_path" validate --spec-root "$spec_root" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for AC index sync issue, got $rc"
  echo "$output"
  exit 1
}

# Error should mention the nature of the sync problem
echo "$output" | grep -qiE '(orphan|phantom|not listed|missing file|not in readme|sync)' || {
  echo "Error output does not describe the sync issue"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
