# Scenario

**Feature**: doctest vet validates the new root tree layout rules

```
# walk tree and check structure
doctest vet <dir> -> validate DOCTEST.md sections -> report violations

# new layout rules
## Version present | Request/Response/Run in DOCTEST.md | root SETUP optional
```

## Preconditions

- `doctest vet` enforces the restructured layout rules.
- Each leaf creates a temporary tree with a specific layout defect or valid configuration.

## Steps

1. Materialize a temporary doctest tree per leaf scenario.
2. Run `doctest vet <treeDir>`.

## Context

- Version checks validate presence only, not value against `VERSION.txt`.
- Valid trees include `## Version`, DSN, and a Go block with `Request`/`Response`/`Run` in `DOCTEST.md`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Env = append(req.Env, "DOCTEST_STRUCTURE_VET=1")
	return nil
}
```