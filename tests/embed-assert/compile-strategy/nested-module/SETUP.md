# Scenario

**Feature**: legacy nested testcase module always wires assert via replace in generated go.mod

```
# public imports only, gen-dir outside module
# always-on assertImport for external modules (with or without author assert import)
doctest test --gen-dir <outside> -> module testcase -> replace assert => cache
```

## Preconditions

- Temp module uses public `pkg/greet` (no `internal/` import).
- `--gen-dir` is outside the parent module so nested `go.mod` is written.

## Steps

1. Descendant creates public-import project and runs doctest with outside gen-dir.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```