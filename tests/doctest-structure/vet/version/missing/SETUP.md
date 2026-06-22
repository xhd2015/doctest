# Scenario

**Bug**: vet must reject `DOCTEST.md` without `## Version`

```
# missing version section
DOCTEST.md (no ## Version) -> doctest vet -> error mentioning version
```

## Preconditions

- A temporary tree whose `DOCTEST.md` lacks `## Version`.

## Steps

1. Write a tree with DSN and valid Go block but no version section.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	treeDir := writeTree(t, treeOpts{withVersion: false})
	setVetArgs(t, req, treeDir)
	return nil
}
```