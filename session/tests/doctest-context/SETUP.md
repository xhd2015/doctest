# Scenario

**Feature**: public `session.Doctest` context type exposes three string fields only

```
# caller constructs or receives session.Doctest
Caller -> session.Doctest{ DOCTEST_ROOT, DOCTEST_CASE, DOCTEST_SESSION_ID }

# field reads return set values; zero value is empty strings
Caller <- d.DOCTEST_ROOT
Caller <- d.DOCTEST_CASE
Caller <- d.DOCTEST_SESSION_ID
```

## Preconditions

- Public package path is `github.com/xhd2015/doctest/session`.
- Type under test (Classic TDD — not implemented yet; tree must RED until it lands):

  ```
  type Doctest struct {
      DOCTEST_ROOT       string
      DOCTEST_CASE       string
      DOCTEST_SESSION_ID string
  }
  ```

- Fields only (no methods). Names are intentional ALL_CAPS inject-contract names.
- Fields are **not** environment variables; do not use `os.Getenv` / `syscall.Getenv`
  to read them.

## Steps

1. Leaf `Setup` sets `req.Mode` and optional `Want*` values.
2. Root `Run` constructs `session.Doctest` according to mode and snapshots fields
   into `Response`.
3. Leaf `Assert` checks empty strings, exact equality, or independence.

## Context

- This tree is a sibling of `session/tests/once/`; it does not exercise `session.Once`.
- Harness may use the injected `DOCTEST_SESSION_ID` variable only for unique
  fixture strings if needed; product fields under test are struct fields.
- P1 scope: data model only. Assemble, harness migration, and docs are out of scope.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Defaults; leaves override Mode and Want*.
	if req.Mode == "" {
		req.Mode = "construct"
	}
	return nil
}
```
