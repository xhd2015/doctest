# Scenario

**Bug**: nested go test timeout panic is not surfaced as a clear Error line

```
# sleep leaf exceeds -timeout=2s
doctest test -timeout=2s --no-color <sleep-tree>
  -> go test -timeout=2s
  -> panic: test timed out after 2s

# expected fail path
doctest -> exit ≠ 0
doctest -> Error: go test timed out after 2s  (visible on stdout/stderr)
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture Run sleeps ≥3s so a 2s go test timeout always fires.

## Steps

1. Create a temp sleep tree (Run sleeps 5s).
2. Run `doctest test -timeout=2s --no-color <tree>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createSleepTree(t, 5)
	req.Args = []string{"test", "-timeout=2s", "--no-color", testDir}
	return nil
}
```
