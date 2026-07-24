# Scenario

**Feature**: timeout Error/hint/cancelled use color accents when `--color` is on

```
# multi-leaf sleep exceeds -timeout=2s with forced color
doctest test --timeout=2s --color <3-sleep-tree>
  -> go test -timeout=2s
  -> panic: test timed out after 2s

# expected colored fail path
doctest -> exit ≠ 0
doctest -> red  Error: go test timed out after 2s
doctest -> gray hint: increase with -timeout=DURATION …
doctest -> FAIL (0/3, N cancelled): FAIL token red; "N cancelled" orange (38;5;208)
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Same multi-sleep fixture as `surfaces` (planned=3) so cancelled > 0.
- `--color` forces ANSI even when stdout is a pipe (harness capture).

## Steps

1. Create a temp 3-leaf sleep tree (Run sleeps 5s).
2. Run `doctest test --timeout=2s --color <tree>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createSleepTree(t, 3, 5)
	req.Args = []string{"test", "--timeout=2s", "--color", testDir}
	return nil
}
```
