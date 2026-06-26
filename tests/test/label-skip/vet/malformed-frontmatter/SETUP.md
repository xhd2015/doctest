# Scenario

**Feature**: vet rejects malformed frontmatter YAML

```
# broken YAML in ASSERT.md
doctest vet <tree> -> non-zero exit
```

## Steps

1. Create temp tree with malformed frontmatter.
2. Run `doctest vet <tree>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeMalformedAssertTree(t)
	req.Args = []string{"vet", root}
	return nil
}
```