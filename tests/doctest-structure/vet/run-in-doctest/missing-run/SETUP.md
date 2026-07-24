# Scenario

**Bug**: vet must reject `DOCTEST.md` Go block without `func Run`

```
# incomplete Go block
DOCTEST.md -> Request + Response, no func Run -> vet error
```

## Preconditions

- `DOCTEST.md` has `## Version` and DSN but its Go block lacks `func Run`.

## Steps

1. Write a tree whose `DOCTEST.md` Go block omits `func Run`.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	treeDir := writeTree(t, treeOpts{
		withVersion:    true,
		doctestGoBlock: typesWithoutRunGoBlock(),
	})
	setVetArgs(t, req, treeDir)
	return nil
}
```