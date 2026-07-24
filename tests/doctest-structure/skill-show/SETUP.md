# Scenario

**Feature**: skill show commands inject the canonical spec version into embedded prompts

```
# embedded markdown templates carry a version placeholder
__DOCTEST_VERSION__ in template -> skill show -> literal 0.0.2 in stdout

# no placeholder residue
resolved output must not contain __DOCTEST_VERSION__
```

## Preconditions

- `cmd/doctest/VERSION.txt` is the single source of truth (currently `0.0.2`).
- `doctest skill tdd|tdd-lite|designer|implementer --show` resolves the placeholder at runtime.

## Steps

1. Run the leaf's `doctest skill <name> --show` command.
2. Assert stdout contains the canonical version and no placeholder.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Env = append(req.Env, "DOCTEST_STRUCTURE_SKILL=1")
	return nil
}
```