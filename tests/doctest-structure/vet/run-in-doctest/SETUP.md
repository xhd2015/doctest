# Scenario

**Feature**: vet requires `Request`/`Response`/`Run` in root `DOCTEST.md` Go block

```
# types and Run live in DOCTEST.md
DOCTEST.md Go block -> type Request, type Response, func Run -> vet pass

# root SETUP.md optional
optional root SETUP.md -> Setup only (no type redefinition)
```

## Preconditions

- `Request`, `Response`, and `func Run` must be defined in the root `DOCTEST.md` Go block.
- Root `SETUP.md` is optional and must not redefine those symbols.

## Steps

1. Create a tree matching the leaf's layout configuration.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "DOCTEST_STRUCTURE_VET_RUN_IN_DOCTEST=1")
	return nil
}
```