# Scenario

**Bug**: vet must reject `Request`/`Response`/`Run` in root `SETUP.md`

```
# legacy root SETUP placement
DOCTEST.md (no types) + root SETUP.md (Request/Response/Run) -> vet error
```

## Preconditions

- `DOCTEST.md` has version and DSN but no type definitions in its Go block.
- Root `SETUP.md` still declares `Request`/`Response`/`Run`.

## Steps

1. Write a tree with types only in root `SETUP.md`.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	treeDir := writeTree(t, treeOpts{
		withVersion:    true,
		withRootSetup:  true,
		rootSetupBody:  "types",
		doctestGoBlock: goBlock("// types intentionally omitted — legacy placement in root SETUP.md\n"),
	})
	setVetArgs(t, req, treeDir)
	return nil
}
```