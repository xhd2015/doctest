# Scenario

**Feature**: `--changed` narrows candidates before `--label` filter

```
git changed files -> leaf subset -> label expression -> run/skip
```

## Preconditions

- Git available on PATH.

## Steps

1. Create committed fixture mod inside a git repo.
2. Modify only the `slow` leaf ASSERT.md.
3. Run `doctest test <mod> --changed --label EXPR` from repo root.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	_ = t
	return nil
}
```