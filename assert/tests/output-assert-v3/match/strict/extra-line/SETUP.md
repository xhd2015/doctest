# Scenario

**Feature**: V3-M10 — strict extra actual line rejected

```
# single-line template cannot match two-line actual
Matcher rejects trailing extra line
```

## Steps
1. Set single literal line template and two-line actual.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Escape the space? No — "only one line" has no RE metachars except... none of . * etc
	// Spaces and letters are fine as raw RE.
	req.Template = v3Template("", "only one line")
	req.Actual = "only one line\nextra line"
	return nil
}
```
