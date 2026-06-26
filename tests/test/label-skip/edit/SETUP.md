# Scenario

**Feature**: `doctest edit` updates single-leaf ASSERT.md frontmatter

```
# concrete leaf only
doctest edit <leaf> --add-label/--add-explanation -> mutate ASSERT.md
```

## Steps

1. Prepare temp tree and invoke `doctest edit`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```