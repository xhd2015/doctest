# Scenario

**Feature**: `--color` and `--no-color` together are rejected

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Both `--color` and `--no-color` are passed to `doctest test`.

## Steps
1. Run `doctest test --color --no-color .` (directory unused — parse fails first).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--color", "--no-color", "."}
	return nil
}
```