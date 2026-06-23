# Scenario

**Feature**: `doctest test --changed` runs only affected leaves

```
# filter leaves by changed files
doctest test --changed <tree> -> build -> run subset -> summary (N Run, ...)
```

## Preconditions

- Fixture tree lives inside an initialized git repository.

## Steps

1. Create and commit a baseline fixture tree.
2. Modify paths per leaf scenario.
3. Run `doctest test <tree> --changed`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_SUBCMD=test")
	return nil
}
```