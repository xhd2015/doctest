# Scenario

**Feature**: V3S-P9 — leading blank lines before `---` header are trimmed

```
# Go raw literal may start with newline before ---
\n\n---\nversion: 3\n---\nbody  -> v3 parse
```

## Steps
1. Set template with blank lines before the opening `---`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = `
---
version: 3
---
hello`
	return nil
}
```
