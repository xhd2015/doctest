# Scenario

**Feature**: go test invocation is suite-only under default unified gen

```
RunTest(2-leaf, GenDir=tmp)
  -> displayed go test has single package containing "suite"
```

## Preconditions

- Default generation.
- Run captures stdout/stderr go test display line.

## Steps

1. `Op=run_gen`.
2. Assert package args.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	return nil
}
```
