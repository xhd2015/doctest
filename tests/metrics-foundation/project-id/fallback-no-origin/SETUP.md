# Scenario

**Feature**: missing origin falls back to nogit_ + short hash of abs root

```
# no origin
abs root path -> sha256 -> nogit_<12 hex chars>
```

## Preconditions

- Empty origin path; absolute root is a stable temp path string.

## Steps

1. Set abs root to a fixed absolute path.
2. Call `ProjectIDFallback`.

## Context

- Prefix is `nogit_` (not `local_`) per P1 exit criteria wording in the requirement.
- Hash is SHA-256 of the absolute root string; first 12 hex characters (lowercase).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "project_id_fallback"
	req.AbsRoot = "/tmp/metrics-foundation-fixture-root"
	return nil
}
```
