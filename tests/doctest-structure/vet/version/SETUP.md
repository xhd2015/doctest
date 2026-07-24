# Scenario

**Feature**: vet checks for `## Version` in root `DOCTEST.md`

```
# version section gate
DOCTEST.md -> ## Version heading -> vet pass or fail
```

## Preconditions

- Root `DOCTEST.md` must declare a `## Version` section.

## Steps

1. Create a tree with or without the version section.
2. Run `doctest vet`.

## Context

- Vet checks presence only; the version value is not compared to `VERSION.txt`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Env = append(req.Env, "DOCTEST_STRUCTURE_VET_VERSION=1")
	return nil
}
```