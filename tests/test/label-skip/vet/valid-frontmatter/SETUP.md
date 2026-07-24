# Scenario

**Feature**: vet accepts valid label and explanation frontmatter

```
# valid YAML on labeled leaf
doctest vet <tree> -> exit 0
```

## Steps

1. Create temp tree with valid label+explanation frontmatter.
2. Run `doctest vet <tree>`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "valid frontmatter")
	req.Args = []string{"vet", root}
	return nil
}
```