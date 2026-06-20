# Scenario

**Feature**: legacy nested-module behavior when no internal imports detected

```
# no internal import: legacy nested module testcase + replace
doctest test --gen-dir <outside> -> module testcase -> public import OK
```

## Preconditions

- Public package imports only (no `internal/` paths in assembled Go).
- Gen-dir outside parent module triggers legacy nested module.

## Steps

1. Each legacy leaf sets full doctest args including the `test` subcommand.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "GOWORK=off")
	return nil
}
```