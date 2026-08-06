# Scenario

**Feature**: user overlay seed merges with internal layers into one `-overlay=` flag

```
user seed Replace
  + optional pre_test hook writes
  + optional vendor-gomod bridge
  -> single driver overlay file
  -> GoFlags: exactly one -overlay=<file> when non-empty
```

## Preconditions

- L2 library surface: `ApplyPreTestHooksWithUserOverlay` and/or
  `MaterializeUserVendorOverlay` (see root DOCTEST implementer surface).
- Fixture paths filled on Response for key assertions.

## Steps

1. Set `Mode=materialize`.
2. Leaves set user Replace, hooks, bridges, or materialize-helper flags.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = modeMaterialize
	return nil
}
```
