# Scenario

**Feature**: vet passes when types are in `DOCTEST.md` and root `SETUP.md` has only `Setup`

```
# valid new layout with optional root setup
DOCTEST.md -> Request/Response/Run
root SETUP.md -> func Setup only -> vet pass
```

## Preconditions

- `DOCTEST.md` contains `## Version`, DSN, and the shared type definitions.
- Root `SETUP.md` defines only `func Setup`, not `Request`/`Response`/`Run`.

## Steps

1. Write a valid tree with optional root `SETUP.md` containing `Setup` only.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	treeDir := writeTree(t, treeOpts{withVersion: true, withRootSetup: true})
	setVetArgs(t, req, treeDir)
	return nil
}
```