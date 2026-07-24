# Scenario

**Feature**: no doctest changes yields zero selection and silent announce

```
# empty or out-of-tree changes
FilterByChangedFiles -> [] ; ChangedCount 0 ; Announce false without -v
```

## Preconditions

- Fixture tree present; changed list does not hit doctest leaves.

## Steps

1. Prepare baseline tree (descendant sets ChangedFiles).
2. Assert zero selection and silent announce.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Policy == "" {
		req.Policy = PolicyFilter
	}
	return nil
}
```
