# Scenario

**Feature**: unified flag sets Options true and auto-enables ref

```
# parse with only the unified flag (no explicit ref flag)
parseTestOptions(["--experiment-unified-package-per-doctest-tree", "./tests"])
  -> ExperimentUnifiedPackagePerDoctestTree=true
  -> ExperimentRefInsteadOfInline=true   # auto-enabled
```

## Preconditions

- Flag may appear before the directory operand.
- Argv does **not** include `--experiment-ref-instead-of-inline`; ref must be forced by unified.

## Steps

1. Parse unified long flag plus a path.
2. Expect unified true and ref true; remain args still include the path.

## Context

- Product rule: unified automatically enables ref-instead-of-inline.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--experiment-unified-package-per-doctest-tree", "./tests"}
	return nil
}
```
