# Scenario

**Feature**: `build --changed` shares test leaf selection via `FilterByChangedFiles`

```
# same filter API as test
DiscoverTreeCases -> FilterByChangedFiles -> subset of leaves
```

## Preconditions

- Selection policy is identical for test and build; this branch documents that.

## Steps

1. Create a baseline fixture tree.
2. Set synthetic changed paths.
3. Assert filtered leaf set (no compile step required for selection policy).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Policy == "" {
		req.Policy = PolicyFilter
	}
	return nil
}
```
