# Scenario

**Feature**: vet rejects legacy placement of `Request`/`Response`/`Run` in root `SETUP.md`

```
# legacy layout forbidden at root
root SETUP.md -> Request/Response/Run -> vet error (must be in DOCTEST.md)
```

## Preconditions

- Root `SETUP.md` must not define `Request`, `Response`, or `func Run`.

## Steps

1. Create a tree with types in root `SETUP.md` instead of `DOCTEST.md`.
2. Run `doctest vet <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "DOCTEST_STRUCTURE_VET_RUN_IN_SETUP=1")
	return nil
}
```