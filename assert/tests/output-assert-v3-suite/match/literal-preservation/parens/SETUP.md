# Scenario

**Feature**: V3S-M12 — escaped parens are literals under v3 raw RE

```
# \(1 Cached\) — parens are not a capturing group
Matcher <- exact (1 Cached)
```

## Steps
1. Set escaped pattern line `\(1 Cached\)` and identical actual.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("", "\\(1 Cached\\)")
	req.Actual = "(1 Cached)"
	return nil
}
```
