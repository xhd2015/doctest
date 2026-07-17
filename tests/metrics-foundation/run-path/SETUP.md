# Scenario

**Feature**: run file paths under cache with UTC timestamps and exclusive create

```
# path shape
$CACHE/doctest/metrics/<project_id>/runs/YYYY-MM-DD-HH-MM-SS-NN-<suffix>.jsonl

# exclusive create
same second -> different NN or suffix -> two distinct files
```

## Preconditions

- Cache dir is a temp directory.
- Timestamps for path formatting are UTC.

## Steps

1. Leaf chooses pure path formatting or exclusive multi-create.
2. Run builds path(s) via metrics helpers.
3. Assert layout segments and uniqueness.

## Context

- Filename time is UTC even if the process local zone differs.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	// Group default project id segment used when a leaf does not override.
	if req.ProjectID == "" {
		req.ProjectID = "github.com_xhd2015_doctest"
	}
	return nil
}
```
