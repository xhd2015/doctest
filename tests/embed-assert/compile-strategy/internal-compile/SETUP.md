# Scenario

**Feature**: internal-compile path resolves assert via temp -modfile instead of nested go.mod

```
# internal/ import detected
doctest -> .doctest_run_* under moduleRoot -> -modfile with parent go.mod + assert replace
```

## Preconditions

- Fixture module imports `example.com/app/internal/greet` in harness `Run()`.
- No nested `go.mod` is written in compile root or gen-dir dump.

## Steps

1. Descendant copies internal fixture and runs doctest with strategy-specific args.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```