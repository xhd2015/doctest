# Scenario

**Bug**: vet must reject `DOCTEST.md` Go block missing `Request` or `Response`

```
# incomplete type definitions
DOCTEST.md -> Response + Run, no type Request -> vet error
```

## Preconditions

- `DOCTEST.md` Go block is missing `type Request`.

## Steps

1. Write a tree whose `DOCTEST.md` Go block omits `type Request`.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	treeDir := writeTree(t, treeOpts{
		withVersion:    true,
		doctestGoBlock: typesWithoutRequestGoBlock(),
	})
	setVetArgs(t, req, treeDir)
	return nil
}
```