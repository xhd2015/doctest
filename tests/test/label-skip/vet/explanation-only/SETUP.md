# Scenario

**Feature**: vet accepts explanation-only frontmatter without label

```
# explanation without label
doctest vet <tree> -> exit 0
```

## Steps

1. Create temp tree with explanation-only frontmatter.
2. Run `doctest vet <tree>`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeExplanationOnlyTree(t, "docs only")
	req.Args = []string{"vet", root}
	return nil
}
```