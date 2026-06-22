# Scenario

**Feature**: vet passes when root `SETUP.md` is absent

```
# minimal valid tree without root SETUP.md
DOCTEST.md only -> Request/Response/Run -> vet pass
```

## Preconditions

- Tree has `DOCTEST.md` with version, DSN, and Go block.
- No root `SETUP.md` file exists.

## Steps

1. Write a valid tree without root `SETUP.md`.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	treeDir := writeTree(t, treeOpts{withVersion: true})
	setVetArgs(t, req, treeDir)
	return nil
}
```