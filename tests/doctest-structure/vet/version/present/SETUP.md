# Scenario

**Feature**: vet accepts any `## Version` value when the section is present

```
# version section present (value not validated)
DOCTEST.md -> ## Version 9.9.9 -> doctest vet -> pass
```

## Preconditions

- A temporary tree whose `DOCTEST.md` includes `## Version` with an arbitrary value.

## Steps

1. Write a tree with `## Version` set to `9.9.9`.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	treeDir := writeTree(t, treeOpts{withVersion: true, version: "9.9.9"})
	setVetArgs(t, req, treeDir)
	return nil
}
```