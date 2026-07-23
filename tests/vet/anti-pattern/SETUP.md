# Scenario

**Feature**: the vet command detects anti-patterns in test file content (L2 in-process)

```
# inspect fixture tree for structural issues
runner.VetArgs([dir]) -> validate.RunWithOptions -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions

- Anti-pattern leaves write bad `SETUP.md` content (or place it under `testdata/`)
  and invoke in-process `VetArgs`.

## Steps

1. Child leaf builds a fixture tree under `t.TempDir()`.
2. Run vet via shared in-process `Run` (no binary).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
