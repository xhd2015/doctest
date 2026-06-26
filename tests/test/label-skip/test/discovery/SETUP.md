# Scenario

**Feature**: discovery-mode `doctest test` skips labeled leaves

```
# tree root / grouping / ./... pattern
doctest test <discovery-target> -> omit labeled ASSERT.md leaves
```

## Steps

1. Build a temp tree with labeled and/or fast leaves.
2. Run `doctest test` against a discovery target (not a concrete leaf).

```go
func Setup(t *testing.T, req *Request) error {
	if req.WorkDir == "" {
		req.WorkDir = "."
	}
	return nil
}
```