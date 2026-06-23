# Scenario

**Feature**: `doctest build --changed` compiles only affected leaves

```
# filter before compile
doctest build --changed <tree> --gen-dir <dir> -> generate subset only
```

## Preconditions

- Fixture tree lives inside an initialized git repository.

## Steps

1. Create and commit a baseline fixture tree.
2. Modify paths per leaf scenario.
3. Run `doctest build <tree> --changed --gen-dir <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_SUBCMD=build")
	return nil
}
```