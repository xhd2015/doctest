# Scenario

**Feature**: a slow leaf produces durations of at least one second

```
# slow leaf execution
doctest test <dir> -> leaf Setup sleeps ~1s -> inline + final DURATION >= 1s
```

## Preconditions

- A temp tree with one leaf whose Setup sleeps for one second exists.

## Steps

1. Create a slow-leaf temp tree.
2. Run `doctest test <dir>` (non-verbose).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createSlowLeafTree(t)
	req.Args = []string{"test", testDir}
	return nil
}
```