# Scenario

**Feature**: hierarchical ref layout under default unified generation

```
RunTest(2-leaf, GenDir=tmp)
  -> shared root marker once; thin leaves import ancestors
```

## Preconditions

- Default generation.
- Explicit GenDir for layout inspection.

## Steps

1. Set `Op=ref_gen`.
2. Assert run success and/or layout metrics.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "ref_gen"
	return nil
}
```
