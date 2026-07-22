# Scenario

**Feature**: explanation without label parses with empty Labels

```
---
explanation: note only
---
→ Labels empty, Explanation set
```

## Steps

1. Op=parse_frontmatter with explanation-only YAML.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse_frontmatter"
	req.FrontmatterContent = "---\nexplanation: documentation note only\n---\n\n## Expected\n"
	return nil
}
```
