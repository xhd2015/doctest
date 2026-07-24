# Scenario

**Feature**: label + explanation frontmatter parses

```
label: ui-automation
explanation: heavy ui
→ Labels [ui-automation], Explanation set
```

## Steps

1. Op=parse_frontmatter with label+explanation.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse_frontmatter"
	req.FrontmatterContent = "---\nlabel: ui-automation\nexplanation: heavy ui test\n---\n\n## Expected\n"
	return nil
}
```
