# Scenario

**Feature**: go test invokes only the suite package (one binary per tree)

```
RunTest(2-leaf, unified=true, GenDir=tmp)
  -> stdout display: cd <gen> && go test … ./…/suite
  -> package args: single suite path (not ./a ./b)
```

## Preconditions

- Runner prints a `go test` display line (same as classic/ref path).
- Layout fill parses package args from that line into `Response.GoTestPackageArgs`.

## Steps

1. Run with unified flag on.
2. Assert exactly one package arg and it refers to `suite`.

## Context

- Locked product shape: one go test package/binary per DOCTEST tree.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = true
	return nil
}
```
