# Scenario

**Feature**: ParseAssertFrontmatter library contract

```
ParseAssertFrontmatter(content) -> Labels, Explanation | error
```

## Steps

1. Leaf sets FrontmatterContent and Op=parse_frontmatter.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse_frontmatter"
	return nil
}
```
