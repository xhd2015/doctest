# Scenario

**Feature**: test/build leaf selection via `FilterByChangedFiles`

```
# filter leaves by synthetic changed files
DiscoverTreeCases -> FilterByChangedFiles -> FilteredPaths + ChangedRunInfo
```

## Preconditions

- Fixture tree lives under a synthetic git root (`t.TempDir()`).

## Steps

1. Create a baseline fixture tree.
2. Set `ChangedFiles` per leaf scenario.
3. `Run` applies filter policy in-process.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Policy == "" {
		req.Policy = PolicyFilter
	}
	return nil
}
```
