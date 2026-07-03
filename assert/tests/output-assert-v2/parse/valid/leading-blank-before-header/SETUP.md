# Scenario

**Feature**: V2-P9 — leading blank lines before `---` header are trimmed

```
# Go raw literal may start with newline before ---
\n\n---\nversion: 2\n---\nbody  -> v2 parse
```

## Steps
1. Set template with blank lines before the opening `---`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = `
---
version: 2
---
hello`
	return nil
}
```