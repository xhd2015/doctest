# Scenario

**Feature**: build and test work with the new root layout

```
# assembler reads types from DOCTEST.md
minimal valid tree -> doctest build -> generated test compiles

# runner executes leaves
minimal valid tree -> doctest test -> leaf passes
```

## Preconditions

- Tree uses new layout: `## Version`, DSN, and `Request`/`Response`/`Run` in `DOCTEST.md`.
- Tree includes one runnable leaf (`leaf/SETUP.md` + `leaf/ASSERT.md`).

## Steps

1. Create a minimal valid tree in `t.TempDir()`.
2. Run `doctest build` or `doctest test` against that tree.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Env = append(req.Env, "DOCTEST_STRUCTURE_INTEGRATION=1")
	return nil
}
```