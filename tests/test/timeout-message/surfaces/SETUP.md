# Scenario

**Bug**: on go test timeout, FAIL still uses actual_run denom and omits cancelled; Error/hint order vs FAIL may be wrong

```
# multi-leaf sleep exceeds -timeout=2s
doctest test --timeout=2s --no-color <3-sleep-tree>
  -> go test -timeout=2s
  -> panic: test timed out after 2s

# expected fail path (plain text, no ANSI)
doctest -> exit ≠ 0
doctest -> Error: go test timed out after 2s
doctest -> hint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)
doctest -> progress (N Run, …) without "cancelled"
doctest -> FAIL (0/3, N cancelled) in …   # planned=3, N>0
# print order: Error/hint before FAIL summary on stdout
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture Run sleeps ≥3s so a 2s go test timeout always fires.
- Three discovery leaves so planned=3 and cancelled is forced > 0 even when
  package-level fail contributes 1 to the fail count.

## Steps

1. Create a temp 3-leaf sleep tree (Run sleeps 5s).
2. Run `doctest test --timeout=2s --no-color <tree>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createSleepTree(t, 3, 5)
	req.Args = []string{"test", "--timeout=2s", "--no-color", testDir}
	return nil
}
```
