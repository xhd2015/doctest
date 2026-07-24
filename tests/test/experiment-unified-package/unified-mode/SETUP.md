# Scenario

**Feature**: default unified gen layout + suite run for a 2-leaf fixture

```
# fixture: root DOCTEST defines ExperimentUnifiedRootMarker + Run; leaves a/, b/
RunTest(fixture, GenDir=tmp)
  -> suite-only go test; hierarchical packages under gen
```

## Preconditions

- Default generation is unified (no flags).
- Explicit GenDir under temp for layout inspection.

## Steps

1. Set `Op=run_gen`.
2. Run fills GenDir layout helpers on Response.
3. Leaf Assert checks success and/or layout.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "run_gen"
	return nil
}
```
